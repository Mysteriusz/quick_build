package clang_vc

import(
	"fmt"
	"qb/misc"

	"qb/qbio"
	"qb/build"
	"qb/policies/vc"
	"qb/policies/llvm/clang/cfg"
)

/*
	Clang source diff works by checking expected output set
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
func DiffSources(_qb_state *qb.BuildState, _vc_state *vc.FileState)(src_diff vc.FileDiff){
	if _vc_state == nil{
		return
	}

	for _, val := range _vc_state.Pipe().OutWorkingSet{
		if val.Type != qb.TYPE_FILE{
			continue
		}

		/*
			Resolve extra objects for the file
		*/

		obj := val.Data.(qb.FileObject).File
		_,dep := qb.GetObjectExtra[qbio.File](&val, clang.OUT_DEP)
		_,src := qb.GetObjectExtra[qbio.File](&val, clang.OUT_SRC)

		/*
			Validate the Dependency file
		*/

		v := !dep.IsValid() 
		h := !dep.InvalidateHash()
		if v || h{
			fmt.Printf("Dependency file changed:\n %s\n", dep.FullPath)
			src_diff.Modified = append(src_diff.Modified, src)
			continue
		}

		/*
			Validate the Object file
		*/

		v = !obj.IsValid() 
		h = !obj.InvalidateHash()
		if v || h{
			fmt.Printf("Object file changed:\n %s\n", obj.FullPath)
			src_diff.Modified = append(src_diff.Modified, src)
			continue
		}

		/*
			Validate the Source file
		*/

		v = !src.IsValid() 
		h = !src.InvalidateHash()
		if v || h{
			fmt.Printf("Source file changed:\n %s\n", src.FullPath)
			if v{
				src_diff.Removed = append(src_diff.Removed, src)
			}else{
				src_diff.Modified = append(src_diff.Modified, src)
			}
			continue
		}

		/*
			Parse the dependency file
		*/

		res, file := clang.ParseD(dep)
		if !res{
			fmt.Printf("Unable to parse the dependency file:\n %s\n", dep.FullPath)
			src_diff.Modified = append(src_diff.Modified, src)
			continue
		}

		if !src.Compare(file.Deps[1]){
			fmt.Printf("Dependency source files do not match:\n %s\n %s\n", src.FullPath, file.Deps[1].FullPath)
			src_diff.Modified = append(src_diff.Modified, src)
			continue
		}

		/*
			Intersect headers from the dependency file and gathered by the 'vc.FileState'
			If any are matched that means the file has to be re-compiled
		*/
		shared_diffs := misc.Intersect(file.Deps[2:].AllPaths(), _vc_state.DiffHeaders.Modified.AllPaths())
		if len(shared_diffs) != 0{
			fmt.Println("Header files do not match")
			src_diff.Modified = append(src_diff.Modified, src)
			continue
		}
	}

	return src_diff
}
func DiffOutForAll(_qb_state *qb.BuildState, _vc_state *vc.FileState) (out_diff vc.ObjectDiff){
	for _, val := range _vc_state.Pipe().OutWorkingSet{
		if val.Type != qb.TYPE_FILE{
			continue
		}

		_,src := qb.GetObjectExtra[qbio.File](&val, clang.OUT_SRC)

		if !src.IsValid(){
			out_diff.Removed.Update(val)
			continue
		}
		
		if !src.InvalidateHash(){
			out_diff.Modified.Update(val)
			continue
		}
	}
	
	return out_diff
}
func DiffOutForAny(_qb_state *qb.BuildState, _vc_state *vc.FileState) (out_diff vc.ObjectDiff){
	return
}

