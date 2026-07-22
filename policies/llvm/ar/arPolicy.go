package ar

import(
	"qb/build"
	"qb/policies"
	"qb/policies/vc"

	"qb/policies/llvm/ar/vc"
	"qb/policies/llvm/ar/run"
	"qb/policies/llvm/ar/cfg"
)

type PolicyInfo struct{
	base 		policies.PolicyFile
	config		*ar.PolicyConfig

	vc.InputDiffProvider
	vc.OutputDiffProvider
}

const POLICY_FILE_PATH string = "./policies/llvm/ar.toml"

func (_policy *PolicyInfo) GetCapabilities() policies.Capabilities{
	return policies.Capabilities{
		VersionControl: true,
	}
}
func (_policy *PolicyInfo) GetFile() *policies.PolicyFile{
	file, res := policies.LoadPolicyFile(POLICY_FILE_PATH)
	if !res{
		return nil
	}

	_policy.base = file
	return &file
}
func (_policy *PolicyInfo) Run(_state *qb.BuildState) qb.BuildError{
	if _state == nil{
		return qb.BuildError{}.NilArgument(_state)
	}

	/*
		Access the policy config
	*/

	policy_name := _state.CurrentPipe().CommandPolicyName
	cfg, res := policies.DecodeConfig[ar.PolicyConfig](_policy.base, policy_name)
	if !res{
		return qb.BuildError{}.New(_state,
			"Unable to find policy by name: '%s'", policy_name)
	}

	/*
		Execute the ar archive function
	*/

	out_set, err := run.ArchiveFromState(&cfg, _state)
	if err.Check(){
		return err
	}

	/*
		Save output to working set
	*/

	_state.ClearWorkingSet()
	_state.LoadWorkingSet(out_set)

	return qb.BuildError{}.None()
}


/*

	================ VERSION CONTROL ================

*/

func (_policy *PolicyInfo) ComputeInputDiff(_qb_state *qb.BuildState, _vc_state *vc.FileState) (vc.ObjectDiff, qb.BuildError){
	if _qb_state == nil || _vc_state == nil{
		return vc.ObjectDiff{}, qb.BuildError{}.NilArgument(_qb_state)
	}

	return ar_vc.DiffIn(_qb_state, _vc_state)
}
func (_policy *PolicyInfo) ComputeOutputDiff(_qb_state *qb.BuildState, _vc_state *vc.FileState) (vc.ObjectDiff, qb.BuildError){
	if _qb_state == nil || _vc_state == nil{
		return vc.ObjectDiff{}, qb.BuildError{}.NilArgument(_qb_state)
	}

	return ar_vc.DiffOut(_qb_state, _vc_state)
}

