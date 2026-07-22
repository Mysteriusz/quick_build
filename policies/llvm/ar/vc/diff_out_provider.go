package ar_vc

import(
	"qb/build"
	"qb/policies/vc"
)

func DiffOut(_qb_state *qb.BuildState, _vc_state *vc.FileState)(vc.ObjectDiff, qb.BuildError){
	if _qb_state == nil || _vc_state == nil{
		return vc.ObjectDiff{}, qb.BuildError{}.NilArgument(_qb_state)
	}

	var out_diff vc.ObjectDiff
	for _, obj := range _vc_state.Pipe().OutWorkingSet{
		if !obj.Exists(){
			out_diff.Removed.Update(obj)
		}
		if !obj.InvalidateHash(){
			out_diff.Modified.Update(obj)
		}
	}

	return out_diff, qb.BuildError{}.None()
}


