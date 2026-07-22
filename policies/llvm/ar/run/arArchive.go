package run

import(
	"fmt"
	"path/filepath"

	"qb/qbio"
	"qb/build"
	"qb/policies/llvm/ar/cfg"
)

/*	

	Execute archive creation only for the qb.BuildState object`s

INPUT:
	_state.WorkingSet with types: TYPE_FILE
	_cfg.Mode
	_cfg.OutputExt
	_cfg.OutputName

OUTPUT:
	[]qb.Object with types: TYPE_FILE

*/
func ArchiveFromState(_config *ar.PolicyConfig, _state *qb.BuildState) (out_set qb.ObjectSet, err qb.BuildError){
	if _config == nil || _state == nil{
		return
	}

	args := make([]string, 0)

	/*
		ar requires a mode in the following example format:
			- "rcs"
	*/
	args = append(args, _config.Mode)

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

	output := filepath.Join(_state.Config.OutputDirectory, qbio.ChangeExtension(_config.OutputName, _config.OutputExt))
	
	// Output has to be at index after _cfg.Mode (required by ar)
	cmd := qb.InitCommand(args, 0, 2)

	// Set execution process to pipe defined command
	cmd.Exec = _state.CurrentPipe().Command
	// Set command working directory as output directory
	cmd.Directory = _state.Config.OutputDirectory

	// Attach input
	cmd.SetInput(inputs)
	// Attach output
	cmd.SetOutput([]string{output})

	res := cmd.Run()
	if !res{
		return qb.ObjectSet{}, qb.BuildError{}.New(_state,
			"Error occured when executing the command.")
	}

	output_obj, res := qb.InitObject(output, qb.TYPE_FILE)
	if !res{
		return
	}	

	return out_set.Update(output_obj), qb.BuildError{}.None()
}

