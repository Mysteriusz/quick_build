package pwsh

import(
	"fmt"
	"slices"

	. "qb/build"
	. "qb/policies/shell/cfg"
)


func PwshExecFromState(_policy *Shell_PolicyConfig, _state *QB_BuildState) (res bool){
	if _policy == nil || _state == nil{
		return false
	}

	var refs []QB_RefVar
	for _, str := range _policy.Args{
		temp, err := QBRefResolve(_state, str) 
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
	args, err := QBRefMergeByKind(refs)
	if !err{
		fmt.Println("Unable to merge reference variables.")
		return
	}
	
	cmd := QBInitCommand(args, 0, 0)
	cmd.Exec = _state.CurrentPipe().Command
	cmd.RunPowershell()

	return true
}

