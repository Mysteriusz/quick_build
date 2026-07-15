package policies

import(
	"fmt"
	"bufio"
	"strings"
	"path/filepath"

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
	Clang diff works by checking expected output set
	and reading it`s metadata

	Since the expected output to be 'diffed' are always
	object '.o' files any other output should be ignored

	JSON representation example:
        "D:\\ax_project\\ax_virt_layer\\win64\\user\\build\\vrow_thread.o": {
          "type": 0,
          "data": {
            "file": {
              "fullpath": "D:\\ax_project\\ax_virt_layer\\win64\\user\\build\\vrow_thread.o",
              "hash": "b5d85e5a23374cb2aba0b974bd3a17f3b4418f5e9bc8a9e2c8a89d24010f6e44"
            }
          },
          "extra": {
            "dependency_file": {
              "fullpath": "D:\\ax_project\\ax_virt_layer\\win64\\user\\build\\vrow_thread.d",
              "hash": "36581a02116138edd0f2fe671e6cc01602cb1cf8c2006926f6862d99559a9b41"
            },
            "source_file": {
              "fullpath": "D:\\ax_project\\ax_virt_layer\\win64\\user\\src\\mte\\pipe\\vrow_thread.c",
              "hash": "1b6af9a686420f01466a394758ee337ffde72072f0a7ce4316dbca46cb4934c3"
            }
          }
        }
*/

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

		// The file should only work for object files
		if filepath.Ext(obj.FullPath) != ".o"{
			continue
		}

		/*
			Resolve extra object for the file
		*/

		_,dep := QBGetObjectExtra[QB_File](&val, CLANG_OUT_DEP)
		_,src := QBGetObjectExtra[QB_File](&val, CLANG_OUT_SRC)

		/*
			Validate the Dependency file
		*/

		if !dep.IsValid() || !dep.InvalidateHash(){
			fmt.Printf("Dependency file changed:\n %s\n", dep.FullPath)
			src_diff = append(src_diff, src)
			continue
		}

		/*
			Validate the Object file
		*/
		if !obj.IsValid() || !obj.InvalidateHash(){
			fmt.Printf("Object file changed:\n %s\n", obj.FullPath)
			src_diff = append(src_diff, src)
			continue
		}

		/*
			Validate the Source file
		*/
		if !src.IsValid() || !src.InvalidateHash(){
			fmt.Printf("Source file changed:\n %s\n", src.FullPath)
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
	defer _file.Save()

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

