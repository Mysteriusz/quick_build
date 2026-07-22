package clang_vc

import(
	"qb/build"
	"qb/policies/vc"
)

func DiffHeaders(_qb_state *qb.BuildState, _vc_state *vc.FileState) (vc.FileDiff, qb.BuildError){
	if _qb_state == nil || _vc_state == nil{
		return vc.FileDiff{}, qb.BuildError{}.NilArgument(_qb_state)
	}
	return vc.DiffFiles(_qb_state.GatherAllHeaders(), _vc_state.Pipe().HeaderFiles), qb.BuildError{}.None()
}

