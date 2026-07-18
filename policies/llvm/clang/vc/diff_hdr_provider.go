package clang_vc

import(
	"qb/build"
	"qb/policies/vc"
)

func DiffHeaders(_qb_state *qb.BuildState, _vc_state *vc.FileState) (vc.FileDiff, qb.BuildError){
	return vc.DiffFiles(_qb_state.GatherAllHeaders(), _vc_state.Pipe().HeaderFiles), qb.BuildError{}.None()
}

