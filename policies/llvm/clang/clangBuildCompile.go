package policies

import(
	"fmt"
	
	. "qb/io"
	. "qb/build"
)

func ClangCompileFromState(_state *QB_BuildState) (out_set QB_ObjectSet, res bool){
	if _state == nil{
		return
	}

	pipe := _state.CurrentPipe()

	args := make([]string, 1)
	args[0] = pipe.Command

	// Require dependency file
	args = append(args, "-MMD")

	// Write all Flags
	args = append(args,
		ClangWriteArgs("-", pipe.Flags)...
	)
	
	// Write all Definitions
	args = append(args,
		ClangWriteArgs("-D", pipe.Definitions)...
	)

	// Write all Hooks
	args = append(args,
		ClangWriteArgs("-I", pipe.Hooks)...
	)

	cmd := QBInitCommand(args, 0, 0)

	/*
		Command execution loop for every source file
	*/
	for _,src := range _state.GetSources().AllPaths(){
		// Resolve output path for the current state
		output_obj := ChangeDirectory(
			ChangeExtension(src, ".o"),
			_state.Config.OutputDirectory)

		output_dep := ChangeExtension(output_obj, ".d")

		// Attach input
		cmd.SetInput([]string{"-c", src})
		
		// Attach output
		cmd.SetOutput([]string{"-o", output_obj})

		// Run the command
		if !cmd.Run(){
			return
		}

		/*
			Initialize the QB_Object as the compiled file
		*/
		
		obj,_ := QBInitObject(output_obj, TYPE_FILE)

		/*
			Write an extra field as the generated dependency and source files
		*/

		if !QBSetObjectExtra(&obj, CLANG_OUT_DEP, QBInitFile(output_dep)){
			fmt.Printf("Failed to write dependency file to object: %s", output_dep)
		}
		if !QBSetObjectExtra(&obj, CLANG_OUT_SRC, QBInitFile(src)){
			fmt.Printf("Failed to write source file to object: %s", src)
		}

		println(src)
		out_set.Update(obj)
	}

	return out_set, true
}
