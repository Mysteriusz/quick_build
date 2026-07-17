package pwsh

import(
	"fmt"
	"slices"

	"qb/build"
	"qb/policies/shell/cfg"
)


func ExecFromState(_policy *shell.PolicyConfig, _state *qb.BuildState) (res bool){
	if _policy == nil || _state == nil{
		return false
	}

	var refs []qb.RefVar
	for _, str := range _policy.Args{
		temp, err := qb.RefResolve(_state, str) 
		if !err{
			fmt.Printf("Invalid string argument reference:\n '%s'\n", str)
			return
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
		fmt.Println("Unable to merge reference variables.")
		return
	}
	
	cmd := qb.InitCommand(args, 0, 0)
	cmd.Exec = _state.CurrentPipe().Command
	cmd.RunPowershell()

	return true
}

