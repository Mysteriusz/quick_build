package policies

import(
	"fmt"
	"time"
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
	cfg.Execute(_state)

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

func ClangCompileFromState(_state *QB_BuildState) (_out_objects []QB_Object, res bool){
	if _state == nil{
		return
	}

	start := time.Now()
	pipe := _state.CurrentPipe()

	args := make([]string, 1)
	args[0] = pipe.Command

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

	objects := make([]QB_Object, 0)

	/*
		Command execution loop for every source file
	*/
	for _,src := range _state.GetSources(){
		// Resolve output path for the current state
		output := ChangeDirectory(
			ChangeExtension(src.FullPath, ".o"),
			_state.Config.OutputDirectory)

		// Attach input
		cmd.SetInput([]string{"-c", src.FullPath})
		
		// Attach output
		cmd.SetOutput([]string{"-o", output})

		// Run the command
		if !cmd.Run(){
			return
		}

		// Initialize the QB_Object as the compiled file
		obj,_ := QBInitObject(output, TYPE_FILE)

		// Append to _state.WorkingSet
		objects = append(objects, obj)
	}
	end := time.Now()

	fmt.Println("Compiled in: ", end.Sub(start))

	return objects, true
}

func ClangLinkFromState(_state *QB_BuildState) (_out_objects []QB_Object, res bool){
	if _state == nil{
		return
	}

	pipe := _state.CurrentPipe()
	args := make([]string, 0)

	// Write all hooks
	if len(pipe.Hooks) == 0{
		// All files are included directly
		args = append(args,
			ClangWriteArgs("-I", _state.GetHeaders().AllPaths())...
		)
	}else{
		// Write only provided hooks
		args = append(args,
			ClangWriteArgs("-I", pipe.Hooks)...
		)
	}

	inputs := make([]string, len(_state.WorkingSet))
	for _,obj := range _state.WorkingSet{
		if obj.Type != TYPE_FILE{
			continue
		}

		file := obj.Data.(QB_FileObject).File
		if !file.IsValid(){
			fmt.Printf("Failed to validate: %s\n", file.FullPath)
		}

		inputs = append(args, file.FullPath)
	}

	/*
		TODO
	*/	
	cmd := QBInitCommand(args, 0, 0)
	_ = cmd
	_ = inputs

	// TODO
	return []QB_Object{}, true
}

