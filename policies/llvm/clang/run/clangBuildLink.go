package run

import(
	"fmt"

	"qb/qbio"
	"qb/build"
	"qb/policies/llvm/clang/cfg"
)

/*
	Clang link command structure (in order)

	[cmd] `
		[Hooks/Headers] `
		[Definitions] `
		[Input] `
		[Flags] `
		[Output] ` -> state.config.output_directory/output_name.output_ext

	Example:
		clang `
		D:/my/hook/dir `
		-Dmydef D:/object0.o D:/object1.o `
		-myflag `
		-o D:/output.exe `
*/
func LinkFromState(_state *qb.BuildState, _cfg *clang.PolicyConfig) (out_set qb.ObjectSet, err qb.BuildError){
	if _state == nil || _cfg == nil{
		return qb.ObjectSet{}, qb.BuildError{}.NilArgument(_state)
	}

	/*
		Custom link specific Variables

		They are required to continue
	*/

	output_name, res := _cfg.Vars["output_name"]
	if !res{
		return qb.ObjectSet{}, qb.BuildError{}.New(_state,
			"Policy variable called 'output_name' not found but required.")
	}

	output_ext, res := _cfg.Vars["output_ext"]
	if !res{
		fmt.Println()
		return qb.ObjectSet{}, qb.BuildError{}.New(_state,
			"Policy variable called 'output_ext' not found but required.")
	}


	pipe := _state.CurrentPipe()
	args := make([]string, 0)

	// Write all hooks
	if len(pipe.Hooks) == 0{
		// All files are included directly
		args = append(args,
			WriteArgs("-I", _state.GatherAllHeaders().AllPaths())...
		)
	}else{
		// Write only provided hooks
		args = append(args,
			WriteArgs("-I", pipe.Hooks)...
		)
	}

	args = append(args,
		WriteArgs("-D", pipe.Definitions)...
	)
	args = append(args,
		WriteArgs("-", pipe.Flags)...
	)

	/*
		Iterate the working set and gather objects to link
	*/

	inputs := make([]string, 0, len(_state.WorkingSet))
	for _,obj := range _state.WorkingSet{
		if obj.Type != qb.TYPE_FILE{
			continue
		}

		file := obj.Data.(qb.FileObject).File
		if !file.IsValid(){
			return qb.ObjectSet{}, qb.BuildError{}.New(_state,
				"Failed to validate file:\n %s\n", file.FullPath)
		}

		inputs = append(inputs, file.FullPath)
	}

	/*
		Resolve fullpath of the output
	*/

	output_path := qbio.ChangeDirectory(qbio.ChangeExtension(output_name.(string), output_ext.(string)), _state.Config.OutputDirectory)

	/*
		Initialize and configure the command 
	*/

	cmd := qb.InitCommand(args, len(pipe.Hooks) + len(pipe.Definitions) + 1, 0)

	// Set execution process to pipe defined command
	cmd.Exec = _state.CurrentPipe().Command
	// Set command working directory as output directory
	cmd.Directory = _state.Config.OutputDirectory

	// Attach input
	cmd.SetInput(inputs)
	// Attach output
	cmd.SetOutput([]string{"-o", output_path})

	res = cmd.Run()
	if !res{
		return qb.ObjectSet{}, qb.BuildError{}.New(_state,
			"Error occured when executing the command.")
	}

	/*
		Initialize the qb.Object as the linked output
	*/

	obj, res := qb.InitObject(output_path, qb.TYPE_FILE)
	if !res{
		return
	}

	out_set.Update(obj)

	return out_set, qb.BuildError{}.None()
}

