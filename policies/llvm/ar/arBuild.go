package policies

import(
	"fmt"
	"path/filepath"

	. "qb/policies"
	. "qb/build"
	. "qb/io"
)

func ArRunFromState(_policy *Ar_Policy, _state *QB_BuildState) (res bool){
	if _state == nil || _policy == nil{
		return
	}

	// Get the policy name to execute
	policy_name := _state.CurrentPipe().CommandPolicyName

	// Load the policy file
	file, res := QBLoadPolicyFile(_policy)
	if !res{
		return
	}

	// Lookup named policy and execute
	cfg, res := QBDecodePolicy[Ar_PolicyConfig](file, policy_name)
	if !res{
		return
	}
	return cfg.Execute(_state)
}

func ArArchiveFromState(_cfg *Ar_PolicyConfig, _state *QB_BuildState) (out_set QB_ObjectSet, res bool){
	if _cfg == nil || _state == nil{
		return
	}

	args := make([]string, 0)

	args = append(args, _state.CurrentPipe().Command)

	/*
		ar requires a mode in the following example format:
			- "rcs"
	*/
	args = append(args, _cfg.Mode)

	inputs := make([]string, 0, len(_state.WorkingSet))
	for _,obj := range _state.WorkingSet{
		if obj.Type != TYPE_FILE{
			continue
		}

		file := obj.Data.(QB_FileObject).File
		if !file.IsValid(){
			fmt.Printf("Failed to validate: %s\n", file.FullPath)
		}
		
		args = append(args, file.FullPath)
	}

	output := filepath.Join(_state.Config.OutputDirectory, ChangeExtension(_cfg.OutputName, _cfg.OutputExt))
	
	// Output has to be at index after _cfg.Mode (required by ar)
	cmd := QBInitCommand(args, 0, 2)

	// Attach output
	cmd.SetOutput([]string{output})

	// Attach input
	cmd.SetInput(inputs)

	res = cmd.Run()
	if !res{
		return
	}

	output_obj, res := QBInitObject(output, TYPE_FILE)
	if !res{
		return
	}	

	return out_set.Update(output_obj), true
}

