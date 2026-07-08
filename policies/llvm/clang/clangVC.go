package policies

import(
	//"fmt"
	"bufio"
	"strings"

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

	for _, obj := range _vc_state.Pipe().OutWorkingSet{
		if obj.Type != TYPE_FILE{
			continue
		}

		_, v := obj.GetExtra(CLANG_EXTRA_FIELD_DEP)
		println(v.(string))
		_ = obj
	}

	return
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

