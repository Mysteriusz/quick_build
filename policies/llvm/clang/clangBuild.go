package policies

import(
	"fmt"
	"path/filepath"

	. "qb/build"
	. "qb/io"
)

func ClangRunFromState(_policy *Clang_Policy, _state *QB_BuildState) (res bool){
	if _state == nil || _policy == nil{
		return
	}

	// Get the policy name to execute
	policy_name := _state.CurrentPipe().CommandPolicyName

	// Load the policy file
	file, res := ClangInitPolicyFile(_policy)
	if !res{
		return
	}

	// Lookup named policy and execute
	cfg, res := file.Policies[policy_name]
	if !res{
		fmt.Printf("Policy file: %s,\ndoes not contain policy called: %s\n",
			_policy.GetFile().FullPath,
			policy_name)
		return
	}

	res = cfg.Execute(_state)
	if !res{
		return
	}

	return true
}

func ClangWriteArgs(_prefix string, _args[] string) (res_args []string){
	res_args = make([]string, len(_args))

	for idx,arg := range _args{
		res_args[idx] = _prefix + arg
	}

	return res_args
}

func ClangToFileObject(_state *QB_BuildState, _path string) (obj_path string){
	if _state == nil{
		return
	}

	outname := filepath.Base(ChangeExtension(_path, ".o"))
	outfile := _state.Config.OutputDirectory + outname
	return outfile
}



