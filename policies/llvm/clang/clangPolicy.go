package clang

import(
	"qb/build"
	"qb/policies"
	"qb/policies/vc"

	"qb/policies/llvm/clang/maps"
	"qb/policies/llvm/clang/cfg"
)

type PolicyInfo struct{
	base 		policies.PolicyFile
	vc.HeaderDiffProvider
	vc.SourceDiffProvider
	vc.InputDiffProvider
	vc.OutputDiffProvider
}

const POLICY_FILE_PATH string = "./policies/llvm/clang.toml"
func (_policy *PolicyInfo) GetFile() *policies.PolicyFile{
	if _policy.base.File.IsOpen(){
		return &_policy.base
	}

	file, res := policies.LoadPolicyFile(POLICY_FILE_PATH)
	if !res{
		println(POLICY_FILE_PATH)
		panic("Corrupted policy file")
	}

	_policy.base = file
	return &_policy.base
}
func (_policy PolicyInfo) GetCapabilities() policies.Capabilities{
	return policies.Capabilities{
		VersionControl: true,
	}
}

func (_policy *PolicyInfo) Run(_state *qb.BuildState) qb.BuildError{
	if _state == nil{
		return qb.BuildError{}.NilArgument(_state)
	}

	policy_name := _state.CurrentPipe().CommandPolicyName

	cfg, res := policies.DecodeConfig[clang.PolicyConfig](_policy.base, policy_name)
	if !res{
		return qb.BuildError{}.New(_state,
			"Unable to find policy by name: '%s'", policy_name)
	}

	// Lookup execute function for the function
	exec, res := maps.EXECUTE_FUNCS[cfg.Function]
	if !res{
		return qb.BuildError{}.New(_state,
			"Unsupported execution function: '%s'", cfg.Function)
	}

	err := exec(_state, &cfg)
	
	if err.Check(){
		return err
	}

	return qb.BuildError{}.None()
}

/*

	================ VERSION CONTROL ================

*/

func (_policy *PolicyInfo) ComputeHeaderDiff(_qb_state *qb.BuildState, _vc_state *vc.FileState) (vc.FileDiff, qb.BuildError){
	if _qb_state == nil || _vc_state == nil{
		return vc.FileDiff{}, qb.BuildError{}.NilArgument(_qb_state)
	}

	/*
		Access the policy config
	*/

	policy_name := _qb_state.CurrentPipe().CommandPolicyName
	cfg, res := policies.DecodeConfig[clang.PolicyConfig](_policy.base, policy_name)
	if !res{
		return vc.FileDiff{}, qb.BuildError{}.New(_qb_state,
			"Unable to find policy by name: '%s'", policy_name)
	}

	/*
		Use the config data
	*/

	diff_func, found := maps.DIFF_HDR_PROVIDERS[cfg.Function]
	if !found{
		return vc.FileDiff{}, qb.BuildError{}.New(_qb_state,
			"Unsupported header provider for the function: '%s'", cfg.Function)
	}

	if diff_func == nil{
		return vc.FileDiff{}, qb.BuildError{}.None()
	}

	return diff_func(_qb_state, _vc_state)
}

func (_policy *PolicyInfo) ComputeSourceDiff(_qb_state *qb.BuildState, _vc_state *vc.FileState) (vc.FileDiff, qb.BuildError){
	if _qb_state == nil || _vc_state == nil{
		return vc.FileDiff{}, qb.BuildError{}.NilArgument(_qb_state)
	}

	/*
		Access the policy config
	*/

	policy_name := _qb_state.CurrentPipe().CommandPolicyName
	cfg, res := policies.DecodeConfig[clang.PolicyConfig](_policy.base, policy_name)
	if !res{
		return vc.FileDiff{}, qb.BuildError{}.New(_qb_state,
			"Unable to find policy by name: '%s'", policy_name)
	}

	/*
		Use the config data
	*/

	diff_func, found := maps.DIFF_SRC_PROVIDERS[cfg.Function]
	if !found{
		return vc.FileDiff{}, qb.BuildError{}.New(_qb_state,
			"Unsupported source provider for the function: '%s'", cfg.Function)
	}
	if diff_func == nil{
		return vc.FileDiff{}, qb.BuildError{}.None()
	}

	return diff_func(_qb_state, _vc_state)
}

func (_policy *PolicyInfo) ComputeInputDiff(_qb_state *qb.BuildState, _vc_state *vc.FileState) (vc.ObjectDiff, qb.BuildError){
	if _qb_state == nil || _vc_state == nil{
		return vc.ObjectDiff{}, qb.BuildError{}.NilArgument(_qb_state)
	}

	/*
		Access the policy config
	*/

	policy_name := _qb_state.CurrentPipe().CommandPolicyName
	cfg, res := policies.DecodeConfig[clang.PolicyConfig](_policy.base, policy_name)
	if !res{
		return vc.ObjectDiff{}, qb.BuildError{}.New(_qb_state,
			"Unable to find policy by name: '%s'", policy_name)
	}

	/*
		Use the config data
	*/

	diff_func, found := maps.DIFF_IN_PROVIDERS[cfg.Function]
	if !found{
		return vc.ObjectDiff{}, qb.BuildError{}.New(_qb_state,
			"Unsupported input provider for the function: '%s'", cfg.Function)
	}
	if diff_func == nil{
		return vc.ObjectDiff{}, qb.BuildError{}.None()
	}

	return diff_func(_qb_state, _vc_state)
}

func (_policy *PolicyInfo) ComputeOutputDiff(_qb_state *qb.BuildState, _vc_state *vc.FileState) (vc.ObjectDiff, qb.BuildError){
	if _qb_state == nil || _vc_state == nil{
		return vc.ObjectDiff{}, qb.BuildError{}.NilArgument(_qb_state)
	}

	/*
		Access the policy config
	*/

	policy_name := _qb_state.CurrentPipe().CommandPolicyName
	cfg, res := policies.DecodeConfig[clang.PolicyConfig](_policy.base, policy_name)
	if !res{
		return vc.ObjectDiff{}, qb.BuildError{}.New(_qb_state,
			"Unable to find policy by name: '%s'", policy_name)
	}

	/*
		Use the config data
	*/

	diff_func, found := maps.DIFF_OUT_PROVIDERS[cfg.Function]
	if !found{
		return vc.ObjectDiff{}, qb.BuildError{}.New(_qb_state,
			"Unsupported output provider for the function: '%s'", cfg.Function)
	}
	if diff_func == nil{
		return vc.ObjectDiff{}, qb.BuildError{}.None()
	}

	return diff_func(_qb_state, _vc_state)
}
