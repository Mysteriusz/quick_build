package clang_vc

import(
	"qb/qbio"
	"qb/build"
	"qb/policies/vc"
	"qb/policies/llvm/clang/cfg"
)

func DiffOutForAll(_qb_state *qb.BuildState, _vc_state *vc.FileState) (vc.ObjectDiff, qb.BuildError){
	if _qb_state == nil || _vc_state == nil{
		return vc.ObjectDiff{}, qb.BuildError{}.NilArgument(_qb_state)
	}

	var out_diff vc.ObjectDiff
	for _, obj := range _vc_state.Pipe().OutWorkingSet{
		if obj.Type != qb.TYPE_FILE{
			continue
		}

		_,src := qb.GetObjectExtra[qbio.File](&obj, clang.OUT_SRC)

		if !src.IsValid(){
			out_diff.Removed.Update(obj)
			continue
		}
		
		if !src.InvalidateHash(){
			out_diff.Modified.Update(obj)
			continue
		}
	}
	
	return out_diff, qb.BuildError{}.None()
}
func DiffOutForAny(_qb_state *qb.BuildState, _vc_state *vc.FileState) (vc.ObjectDiff, qb.BuildError){
	if _qb_state == nil || _vc_state == nil{
		return vc.ObjectDiff{}, qb.BuildError{}.NilArgument(_qb_state)
	}

	var out_diff vc.ObjectDiff
	return out_diff, qb.BuildError{}.None()
}

