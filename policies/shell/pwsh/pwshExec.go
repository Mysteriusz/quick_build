package pwsh

import(
	"slices"

	"qb/build"
	"qb/policies/shell/cfg"
)


func ExecFromState(_policy *shell.PolicyConfig, _state *qb.BuildState) qb.BuildError{
	if _policy == nil || _state == nil{
		return qb.BuildError{}.NilArgument(_state)
	}

	var refs []qb.RefVar
	for _, str := range _policy.Args{
		temp, res := qb.RefResolve(_state, str) 
		if !res{
			return qb.BuildError{}.New(_state,
				"Invalid string argument reference:\n '%s'", str)
		}

		refs = slices.Concat(refs, temp)
	}

	for idx := range refs{
		PwshFormatRef(&refs[idx])
	}

	/*
		Merge all strings based on their kind
	*/
	args, err := qb.RefMergeByKind(refs)
	if !err{
		return qb.BuildError{}.New(_state,
			"Unable to merge reference variables.")
	}
	
	cmd := qb.InitCommand(args, 0, 0)
	cmd.Exec = _state.CurrentPipe().Command
	cmd.RunPowershell()

	return qb.BuildError{}.None()
}

