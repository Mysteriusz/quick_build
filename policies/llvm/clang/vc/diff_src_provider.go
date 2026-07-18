package clang_vc

import(
	"qb/qbio"
	"qb/build"
	"qb/policies/vc"
	"qb/policies/llvm/clang/cfg"
)

/*
	IMPORTANT!
	Required 'vc.FileState.DiffHeaders' to be computed

	Clang source diff works by checking expected output set
	and reading it`s metadata

	Since the expected output to be 'diffed' are always
	object files any other output should be ignored

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
func DiffSources(_qb_state *qb.BuildState, _vc_state *vc.FileState)(vc.FileDiff, qb.BuildError){
	if _qb_state == nil || _vc_state == nil{
		return vc.FileDiff{}, qb.BuildError{}.NilArgument(_qb_state)
	}

	var src_diff vc.FileDiff

	/*
		Catch all new files
	*/
	src_diff.Modified = vc.DiffNewFiles(_vc_state.Pipe().SourceFiles, _qb_state.GatherAllSources())
	println(len(src_diff.Modified))

	/*
		Catch all modified/removed files
	*/
	for _, obj := range _vc_state.Pipe().OutWorkingSet{
		if obj.Type != qb.TYPE_FILE{
			continue
		}

		/*
			Verify the output object data
		*/

		_,src_file := qb.GetObjectExtra[qbio.File](&obj, clang.OUT_SRC)

		removed, no_change := OutObjectVerify(obj, _vc_state)
		if no_change{
			continue
		}

		if removed{
			src_diff.Removed = append(src_diff.Removed, src_file)
		}else{
			src_diff.Modified = append(src_diff.Modified, src_file)
		}
	}

	return src_diff, qb.BuildError{}.None()
}

