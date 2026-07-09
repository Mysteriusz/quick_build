package policies

import(
	"fmt"
	. "qb/build"
)

func ClangLinkFromState(_state *QB_BuildState) (out_set QB_ObjectSet, res bool){
	if _state == nil{
		return
	}

	pipe := _state.CurrentPipe()
	args := make([]string, 0)

	// Write all hooks
	if len(pipe.Hooks) == 0{
		// All files are included directly
		args = append(args,
			ClangWriteArgs("-I", _state.GatherAllHeaders().AllPaths())...
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
	return nil, true
}

