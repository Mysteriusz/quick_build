package clang

import(
	"qb/qbio"
	"qb/build"
	"qb/policies/vc"

	"qb/policies/llvm/clang/vc"
	"qb/policies/llvm/clang/maps"
)

func Diff(_policy *Policy, _qb_state *qb.BuildState, _vc_state *vc.FileState){
	if _policy == nil || _qb_state == nil || _vc_state == nil{
		return
	}

	cfg := _policy.GetConfig(_qb_state)
	output_diff_func, res := maps.OUTPUT_DIFF_FUNCS[cfg.Function]
	if !res{
		return
	}

	/*
		Output diff
	*/
	
	_vc_state.DiffOutput = output_diff_func(_qb_state, _vc_state)	

	/*
		Source diff
	*/

	src_diff := clang_vc.DiffSources(_qb_state, _vc_state)
	_vc_state.DiffSources.Modified = qbio.FileArrayUnion(_vc_state.DiffSources.Modified, src_diff.Modified)
}

