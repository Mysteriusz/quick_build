package policies

import(
	"path/filepath"

	"qb/misc"
	"qb/policies/llvm/clang/maps"

	. "qb/policies/llvm/clang/cfg"
	. "qb/policies/vc"
	. "qb/policies"
	. "qb/build"
	. "qb/io"
)

type Clang_Policy struct{
	/*
		Should never be modified
	*/
	PATH 		string // Has to start with '.' character
	CAPS 		QB_Capabilities

	file 		QB_File
	config		*Clang_PolicyConfig
}

func (_policy *Clang_Policy) GetConfig(_state *QB_BuildState) *Clang_PolicyConfig{
	if _policy.config != nil{
		return _policy.config
	}

	// Get the policy name to execute
	policy_name := _state.CurrentPipe().CommandPolicyName

	// Load the policy file
	file, res := QBLoadPolicyFile(_policy)
	if !res{
		return nil
	}

	// Lookup named policy and execute
	cfg, res := QBDecodePolicy[Clang_PolicyConfig](file, policy_name)
	if !res{
		return nil
	}
	
	_policy.config = &cfg
	return _policy.config
}

func (_policy *Clang_Policy) Run(_state *QB_BuildState) (res bool){
	if _state == nil{
		return false
	}

	cfg := _policy.GetConfig(_state)
	if cfg == nil{
		return
	}

	exec, res := clang.EXECUTE_FUNCS[cfg.Function]
	if !res{
		return
	}

	res = exec(cfg, _state)
	if !res{
		return
	}

	return true
}

func (_policy *Clang_Policy) GetCapabilities() QB_Capabilities{
	return _policy.CAPS
}

/*

	================ VERSION CONTROL ================

*/

func (_policy *Clang_Policy) BeginVersionControl(_state *QB_BuildState) (not_first_build bool, not_updated bool, vc_state VC_FileState){
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
	not_first_build, vc_state = VCFindOrCreateState(_state)
	no_diff, no_crit_diff := VCDiff(_state, &vc_state)
	
	/*
		If critical diff is present,
		then act as a complete rebuild
	*/
	not_first_build = not_first_build && no_crit_diff

	/*
		CLANG EXCLUSIVE

		Update diff source files for clang 
		(Only when the build already exist since it requires VC_FileState.OutWorkingSet)
	*/
	if not_first_build{
		ClangVCDiff(_policy, _state, &vc_state)
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
func (_policy *Clang_Policy) EndVersionControl(_qb_state *QB_BuildState, _vc_state *VC_FileState){
	if _qb_state == nil{
		return
	}
	
	_vc_state.Pipe().SourceFiles = _qb_state.GatherAllSources()
	_vc_state.Pipe().HeaderFiles = _qb_state.GatherAllHeaders()
	_vc_state.Pipe().StateHash = VCStateUniqueHash(_qb_state)

	// Save to file
	_vc_state.File.Save()
}

func (_policy *Clang_Policy) GetFile() *QB_File{
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
