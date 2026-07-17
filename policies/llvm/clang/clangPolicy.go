package clang

import(
	"path/filepath"

	"qb/misc"
	"qb/qbio"
	"qb/build"
	"qb/policies"
	"qb/policies/vc"

	"qb/policies/llvm/clang/maps"
	"qb/policies/llvm/clang/cfg"
)

type Policy struct{
	/*
		Should never be modified
	*/
	PATH 		string // Has to start with '.' character
	CAPS 		policies.Capabilities

	file 		qbio.File
	config		*clang.PolicyConfig
}

func (_policy *Policy) GetConfig(_state *qb.BuildState) *clang.PolicyConfig{
	if _policy.config != nil{
		return _policy.config
	}

	// Get the policy name to execute
	policy_name := _state.CurrentPipe().CommandPolicyName

	// Load the policy file
	file, res := policies.LoadPolicyFile(_policy)
	if !res{
		return nil
	}

	// Lookup named policy and execute
	cfg, res := policies.DecodePolicy[clang.PolicyConfig](file, policy_name)
	if !res{
		return nil
	}
	
	_policy.config = &cfg
	return _policy.config
}

func (_policy *Policy) Run(_state *qb.BuildState) (res bool){
	if _state == nil{
		return false
	}

	cfg := _policy.GetConfig(_state)
	if cfg == nil{
		return
	}

	exec, res := maps.EXECUTE_FUNCS[cfg.Function]
	if !res{
		return
	}

	res = exec(cfg, _state)
	if !res{
		return
	}

	return true
}

func (_policy *Policy) GetCapabilities() policies.Capabilities{
	return _policy.CAPS
}

/*

	================ VERSION CONTROL ================

*/

func (_policy *Policy) BeginVersionControl(_state *qb.BuildState) (not_first_build bool, not_updated bool, vc_state vc.FileState){
	if _state == nil{
		return
	}
	if !_policy.GetCapabilities().VersionControl{
		return
	}

	/*
		Load/Create version control object
		and load it`s diff
	*/
	not_first_build, vc_state = vc.FindOrCreateState(_state)
	no_diff, no_crit_diff := vc.Diff(_state, &vc_state)
	
	/*
		If critical diff is present,
		then act as a complete rebuild
	*/
	not_first_build = not_first_build && no_crit_diff

	/*
		CLANG EXCLUSIVE

		Update diff source files for clang 
		(Only when the build already exist since it requires vc.FileState.OutWorkingSet)
	*/
	if not_first_build{
		Diff(_policy, _state, &vc_state)
	}

	if vc_state.DiffHeaders.Len() > 0{
		println()
		println("==================================HDR DIFF==================================")
		misc.PrintArray(vc_state.DiffHeaders.Modified.AllPaths())
		misc.PrintArray(vc_state.DiffHeaders.Removed.AllPaths())
		println()
	}
	if vc_state.DiffSources.Len() > 0{
		println()
		println("==================================SRC DIFF==================================")
		misc.PrintArray(vc_state.DiffSources.Modified.AllPaths())
		misc.PrintArray(vc_state.DiffSources.Removed.AllPaths())
		println()
	}

	return not_first_build, no_diff && vc_state.DiffSources.Len() == 0, vc_state
}
func (_policy *Policy) EndVersionControl(_qb_state *qb.BuildState, _vc_state *vc.FileState){
	if _qb_state == nil{
		return
	}
	
	_vc_state.Pipe().SourceFiles = _qb_state.GatherAllSources()
	_vc_state.Pipe().HeaderFiles = _qb_state.GatherAllHeaders()
	_vc_state.Pipe().StateHash = vc.StateUniqueHash(_qb_state)

	// Save to file
	_vc_state.File.Save()
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
