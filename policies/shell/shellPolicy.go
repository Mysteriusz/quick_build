package policies

import(
	"path/filepath"

	. "qb/io"
	. "qb/build"
	. "qb/policies"
	. "qb/policies/vc"

	. "qb/policies/shell/cfg"
	. "qb/policies/shell/maps"
)

type Shell_Policy struct{
	/*
		Should never be modified
	*/
	PATH 		string // Has to start with '.' character
	CAPS 		QB_Capabilities

	file 		QB_File
}

func (_policy Shell_Policy) GetCapabilities() QB_Capabilities{
	return _policy.CAPS
}
func (_policy *Shell_Policy) GetFile() *QB_File{
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

	_policy.file = QBInitFile(abs)
	return &_policy.file
}

func (_policy *Shell_Policy) Run(_state *QB_BuildState) (res bool){
	if _state == nil{
		return
	}

	// Get policy name
	policy_name := _state.CurrentPipe().CommandPolicyName

	// Load the policy file
	file, res := QBLoadPolicyFile(_policy)
	if !res{
		return
	}

	// Find and decode the policy by name
	cfg, res := QBDecodePolicy[Shell_PolicyConfig](file, policy_name)
	if !res{
		return
	}

	// Lookup cli entry function
	exec, res := ShellCliFuncLookup(cfg.Cli)
	if !res{
		return
	}

	// Execute cli entry function
	return exec(&cfg, _state)
}

/*

	================ VERSION CONTROL ================

*/

func (_policy *Shell_Policy) BeginVersionControl(_state *QB_BuildState) (not_first_build bool, not_updated bool, _vc_state VC_FileState){
	return
}
func (_policy *Shell_Policy) EndVersionControl(_state *QB_BuildState, _vc_state *VC_FileState){
}

