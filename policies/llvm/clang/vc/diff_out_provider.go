package clang_vc

import(
	"qb/qbio"
	"qb/build"
	"qb/policies/vc"
	"qb/policies/llvm/clang/cfg"
)

func DiffOutForAll(_qb_state *qb.BuildState, _vc_state *vc.FileState) (vc.ObjectDiff, qb.BuildError){
	var out_diff vc.ObjectDiff
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
	
	return out_diff, qb.BuildError{}.None()
}
func DiffOutForAny(_qb_state *qb.BuildState, _vc_state *vc.FileState) (vc.ObjectDiff, qb.BuildError){
	var out_diff vc.ObjectDiff
	return out_diff, qb.BuildError{}.None()
}

