package run

import(
	"fmt"
	
	"qb/qbio"
	"qb/build"
	"qb/policies/llvm/clang/cfg"
)

/*
	Clang compile command structure (in order)

	[cmd] `
		-MMD ` -> (Required for version control)
		[Flags] `
		[Definitions] `
		[Hooks] `
		[Input] `
		[Output] `

	Example:
		clang `
			-MMD `
			-myflag `
			-Dmydef `
			-c D:/my/hook/dir ` 
			D:/source.c `
			-o D:/my_ouptut.o `
*/
func CompileFromState(_state *qb.BuildState) (out_set qb.ObjectSet, err qb.BuildError){
	if _state == nil{
		return qb.ObjectSet{}, qb.BuildError{}.NilArgument(_state)
	}

	pipe := _state.CurrentPipe()

	args := make([]string, 0)

	// Require dependency file
	args = append(args, "-MMD")

	// Write all Flags
	args = append(args,
		WriteArgs("-", pipe.Flags)...
	)
	
	// Write all Definitions
	args = append(args,
		WriteArgs("-D", pipe.Definitions)...
	)

	// Write all Hooks
	args = append(args,
		WriteArgs("-I", pipe.Hooks)...
	)

	cmd := qb.InitCommand(args, 0, 0)

	/*
		Command execution loop for every source file
	*/
	for _,src := range _state.GetSources().AllPaths(){
		// Resolve output path for the current state
		output_path := qbio.ChangeDirectory(
			qbio.ChangeExtension(src, ".o"),
			_state.Config.OutputDirectory)

		output_dep := qbio.ChangeExtension(output_path, ".d")

		/*
			Configure the command
		*/

		// Set execution process to pipe defined command
		cmd.Exec = _state.CurrentPipe().Command
		// Set command working directory as output directory
		cmd.Directory = _state.Config.OutputDirectory

		// Attach input
		cmd.SetInput([]string{"-c", src})
		// Attach output
		cmd.SetOutput([]string{"-o", output_path})

		// Run the command
		if !cmd.Run(){
			return qb.ObjectSet{}, qb.BuildError{}.New(_state,
				"Error occured when executing the command.")
		}

		/*
			Initialize the QB_Object as the compiled file
		*/
		
		obj,_ := qb.InitObject(output_path, qb.TYPE_FILE)

		/*
			Write an extra field as the generated dependency and source files
		*/

		if !qb.SetObjectExtra(&obj, clang.OUT_DEP, qbio.InitFile(output_dep)){
			fmt.Printf("Failed to write dependency file to object: %s", output_dep)
		}
		if !qb.SetObjectExtra(&obj, clang.OUT_SRC, qbio.InitFile(src)){
			fmt.Printf("Failed to write source file to object: %s", src)
		}

		out_set.Update(obj)
	}

	return out_set, qb.BuildError{}.None()
}
