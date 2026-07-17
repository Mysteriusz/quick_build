package ar

import(
	"fmt"
	"path/filepath"

	"qb/policies"
	"qb/build"
	"qb/qbio"
)

func RunFromState(_policy *Policy, _state *qb.BuildState) (res bool){
	if _state == nil || _policy == nil{
		return
	}

	// Get the policy name to execute
	policy_name := _state.CurrentPipe().CommandPolicyName

	// Load the policy file
	file, res := policies.LoadPolicyFile(_policy)
	if !res{
		return
	}

	// Lookup named policy and execute
	cfg, res := policies.DecodePolicy[PolicyConfig](file, policy_name)
	if !res{
		return
	}
	return cfg.Execute(_state)
}

func ArchiveFromState(_cfg *PolicyConfig, _state *qb.BuildState) (out_set qb.ObjectSet, res bool){
	if _cfg == nil || _state == nil{
		return
	}

	args := make([]string, 0)

	/*
		ar requires a mode in the following example format:
			- "rcs"
	*/
	args = append(args, _cfg.Mode)

	inputs := make([]string, 0, len(_state.WorkingSet))
	for _,obj := range _state.WorkingSet{
		if obj.Type != qb.TYPE_FILE{
			continue
		}

		file := obj.Data.(qb.FileObject).File
		if !file.IsValid(){
			fmt.Printf("Failed to validate: %s\n", file.FullPath)
		}
		
		args = append(args, file.FullPath)
	}

	output := filepath.Join(_state.Config.OutputDirectory, qbio.ChangeExtension(_cfg.OutputName, _cfg.OutputExt))
	
	// Output has to be at index after _cfg.Mode (required by ar)
	cmd := qb.InitCommand(args, 0, 2)

	// Attach output
	cmd.SetOutput([]string{output})

	// Attach input
	cmd.SetInput(inputs)

	// Set execution process to pipe defined command
	cmd.Exec = _state.CurrentPipe().Command
	// Set command working directory as output directory
	cmd.Directory = _state.Config.OutputDirectory

	res = cmd.Run()
	if !res{
		return
	}

	output_obj, res := qb.InitObject(output, qb.TYPE_FILE)
	if !res{
		return
	}	

	return out_set.Update(output_obj), true
}

