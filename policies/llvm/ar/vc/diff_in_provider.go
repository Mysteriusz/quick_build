package ar_vc

import(
	"qb/build"
	"qb/policies/vc"
)

func DiffIn(_qb_state *qb.BuildState, _vc_state *vc.FileState)(vc.ObjectDiff, qb.BuildError){
	if _qb_state == nil || _vc_state == nil{
		return vc.ObjectDiff{}, qb.BuildError{}.NilArgument(_qb_state)
	}
	return vc.DiffObjects(_qb_state.WorkingSet, _vc_state.Pipe().InWorkingSet), qb.BuildError{}.None()
}

