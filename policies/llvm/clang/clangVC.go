package policies

import(
	. "qb/io"
	. "qb/build"
	. "qb/policies/vc"
	. "qb/policies/llvm/clang/vc"

	"qb/policies/llvm/clang/maps"
)

func ClangVCDiff(_policy *Clang_Policy, _qb_state *QB_BuildState, _vc_state *VC_FileState){
	if _policy == nil || _qb_state == nil || _vc_state == nil{
		return
	}

	cfg := _policy.GetConfig(_qb_state)
	output_diff_func, res := clang.OUTPUT_DIFF_FUNCS[cfg.Function]
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

	src_diff := ClangVCDiffSources(_qb_state, _vc_state)
	_vc_state.DiffSources.Modified = QBFileArrayUnion(_vc_state.DiffSources.Modified, src_diff.Modified)
}

