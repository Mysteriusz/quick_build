package policies

import(
	"fmt"
	"bufio"
	"strings"

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
func ClangVCDiff(_vc_state *VC_FileState)(src_diff QB_FileArray){
	if _vc_state == nil{
		return
	}

	// NOT WORKING
	for _,obj := range _vc_state.Pipe().OutWorkingSet{
		if obj.Type != TYPE_FILE{
			continue
		}

		// Object file
		obj := obj.Data.(QB_FileObject).File

		// Source file
		src := QBInitFile(ChangeExtension(obj.FullPath, ".c"))

		// Dependency file
		dep := QBInitFile(ChangeExtension(obj.FullPath, ".d"))

		// Validate dependency file
		if !dep.IsValid(){
			fmt.Printf("Missing definition file:\n %s\n", dep.FullPath)
			src_diff = append(src_diff, src)
			continue
		}

		// Parse dependency file
		r, d := ClangVCParseD(dep)
		if !r{
			panic("Unable to parse dependecy file.")
		}

		/*
			Dependency file doesn`t match object file using it
			Dependency file doesn`t match source file using it
		*/
		if !ComparePath(d.obj.FullPath, obj.FullPath) || !ComparePath(d.src.FullPath, src.FullPath){
			println(obj.FullPath)
			println(d.obj.FullPath)

			println(src.FullPath)
			println(d.src.FullPath)
			fmt.Println("Corrupted dependency files.")
			continue
		}

		// Check for headers in diff array
		header_diff := VCIntersectFiles(d.deps, _vc_state.DiffHeaders)

		// Check if source file should be updated (Header had changed)
		if len(header_diff) == 0{
			continue
		}

		src_diff = append(src_diff, src)
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

