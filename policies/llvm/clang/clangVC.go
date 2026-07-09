package policies

import(
	"fmt"
	"bufio"
	"strings"

	"qb/misc"
	. "qb/policies/vc"
	. "qb/io"
	. "qb/build"
)

type ClangVC_D struct{
	obj 	QB_File // '.o' file
	src 	QB_File // '.c' file
	deps 	QB_FileArray // All dependencies
}

/*
	Gather all diffs based on the vc state
*/
func ClangVCDiff(_qb_state *QB_BuildState, _vc_state *VC_FileState)(src_diff QB_FileArray){
	if _vc_state == nil{
		return
	}

	for _, val := range _vc_state.Pipe().OutWorkingSet{
		if val.Type != TYPE_FILE{
			continue
		}

		obj := val.Data.(QB_FileObject).File
		_,dep := QBGetObjectExtra[QB_File](&val, CLANG_OUT_DEP)
		_,src := QBGetObjectExtra[QB_File](&val, CLANG_OUT_SRC)

		/*
			Validate the Dependency file
		*/
		if !dep.IsValid() || !(dep.ComputeHash() == QBInitFile(dep.FullPath).ComputeHash()){
			fmt.Printf("Dependency file changed:\n %s\n", dep.FullPath)
			src_diff = append(src_diff, src)
			continue
		}

		/*
			Validate the Object file
		*/
		if !obj.IsValid() || !(obj.ComputeHash() == QBInitFile(obj.FullPath).ComputeHash()){
			fmt.Printf("Object file changed:\n %s\n", obj.FullPath)
			src_diff = append(src_diff, src)
			continue
		}

		/*
			Validate the Source file
		*/
		if !src.IsValid() || !(src.ComputeHash() == QBInitFile(src.FullPath).ComputeHash()){
			fmt.Printf("Object file changed:\n %s\n", src.FullPath)
			src_diff = append(src_diff, src)
			continue
		}

		/*
			Parse the dependency file
		*/
		res, file := ClangVCParseD(dep)
		if !res{
			fmt.Printf("Unable to parse the dependency file:\n %s\n", dep.FullPath)
			src_diff = append(src_diff, src)
			continue
		}

		/*
			Intersect headers from the dependency file and gathered by the 'VC_FileState'
			If any are matched that means the file has to be re-compiled
		*/
		shared_diffs := misc.Intersect(file.deps.AllPaths(), _vc_state.DiffHeaders.AllPaths())
		if len(shared_diffs) != 0{
			src_diff = append(src_diff, src)
			continue
		}
	}

	return src_diff
}

/*
	Trim character set for clang generated .d files

	Contains:
		- Tab (0x09) 
		- Space (0x20)
		- Slash-Back (0x5c)
		- Colon (0x3a)
*/
const CLANG_VC_D_TRIM = "\x09\x20\x5c\x3a"
func ClangVCParseD(_file QB_File) (res bool, dep ClangVC_D){
	scanner := bufio.NewScanner(_file.GetFile())

	/*
		Read the object file
	*/
	if !scanner.Scan(){
		return
	}
	dep.obj = QBInitFile(strings.Trim(scanner.Text(), CLANG_VC_D_TRIM))

	/*
		Read the source file
	*/
	if !scanner.Scan(){
		return
	}
	dep.src = QBInitFile(strings.Trim(scanner.Text(), CLANG_VC_D_TRIM))

	/*
		Read the source file
	*/
	for scanner.Scan(){
		path := strings.Trim(scanner.Text(), CLANG_VC_D_TRIM)
		dep.deps = append(dep.deps, QBInitFile(path))
	}

	return true, dep
}

