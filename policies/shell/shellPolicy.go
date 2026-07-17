package shell

import(
	"path/filepath"

	"qb/qbio"
	"qb/build"
	"qb/policies"
	"qb/policies/vc"

	"qb/policies/shell/cfg"
	"qb/policies/shell/maps"
)

type Policy struct{
	/*
		Should never be modified
	*/
	PATH 		string // Has to start with '.' character
	CAPS 		policies.Capabilities

	file 		qbio.File
}

func (_policy Policy) GetCapabilities() policies.Capabilities{
	return _policy.CAPS
}
func (_policy *Policy) GetFile() *qbio.File{
	if _policy.file.IsValid(){
		return &_policy.file
	}

	if _policy.PATH[0] != '.'{
		panic("CLANG_POLICY: invalid corrupted path.")
	}

	abs, err := filepath.Abs(_policy.PATH)
	if err != nil{
		panic("CLANG_POLICY: Unable to resolve relative path.")
	}	

	_policy.file = qbio.InitFile(abs)
	return &_policy.file
}

func (_policy *Policy) Run(_state *qb.BuildState) (res bool){
	if _state == nil{
		return
	}

	// Get policy name
	policy_name := _state.CurrentPipe().CommandPolicyName

	// Load the policy file
	file, res := policies.LoadPolicyFile(_policy)
	if !res{
		return
	}

	// Find and decode the policy by name
	cfg, res := policies.DecodePolicy[shell.PolicyConfig](file, policy_name)
	if !res{
		return
	}

	// Lookup cli entry function
	exec, res := maps.CliFuncLookup(cfg.Cli)
	if !res{
		return
	}

	// Execute cli entry function
	return exec(&cfg, _state)
}

/*

	================ VERSION CONTROL ================

*/

func (_policy *Policy) BeginVersionControl(_state *qb.BuildState) (not_first_build bool, not_updated bool, _vc_state vc.FileState){
	return
}
func (_policy *Policy) EndVersionControl(_state *qb.BuildState, _vc_state *vc.FileState){
}

