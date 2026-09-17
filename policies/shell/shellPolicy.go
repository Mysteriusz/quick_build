package shell

import(
	"path/filepath"

	"qb/build"
	"qb/policies"

	"qb/policies/shell/cfg"
	"qb/policies/shell/maps"
)

type PolicyInfo struct{
	base 		policies.PolicyFile
	config		*shell.PolicyConfig
}

const POLICY_FILE_PATH string = "/shell/shell.toml"

func (_policy *PolicyInfo) GetCapabilities() policies.Capabilities{
	return policies.Capabilities{
		VersionControl: false,
	}
}
func (_policy *PolicyInfo) GetFile(_search_directory string) *policies.PolicyFile{
	if _policy.base.File.IsOpen(){
		return &_policy.base
	}

	file, res := policies.LoadPolicyFile(filepath.Join(_search_directory, POLICY_FILE_PATH))
	if res.Check(){
		panic(res.Message())
	}

	_policy.base = file
	return &file
}
func (_policy *PolicyInfo) Run(_state *qb.BuildState) qb.BuildError{
	if _state == nil{
		return qb.BuildError{}.NilArgument(_state)
	}

	policy_name := _state.CurrentPipe().CommandPolicyName

	// Find and decode the policy by name
	cfg, res := policies.DecodeConfig[shell.PolicyConfig](_policy.base, policy_name)
	if !res{
		return qb.BuildError{}.New(_state,
			"Unable to find policy by name: '%s'", policy_name)
	}

	// Lookup cli entry function
	exec, res := maps.CliFuncLookup(cfg.Cli)
	if !res{
		return qb.BuildError{}.New(_state,
			"Unsupported cli function: '%s'", cfg.Cli)
	}

	// Execute cli entry function
	return exec(&cfg, _state)
}

/*

	================ VERSION CONTROL ================

*/

