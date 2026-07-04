package runner

import(
	"fmt"

	. "qb/build"
	. "qb/policies/map"
)

func ExecutePolicy(_state *QB_BuildState, _data any) (res bool){
	/*
		Get the pipe
	*/
	pipe := _state.CurrentPipe()
	if pipe == nil{
		fmt.Println("Unable to get pipe at index: %i", _state.CurrentPipeIndex())
		return
	}

	/*
		Lookup and execute the policy
	*/
	policy, res := QBPolicyLookup(pipe.CommandPolicyAlias)
	if !res{
		return
	}

	/*
		Execute a version control check of the policy info
		and ignore if it was successfull
	*/
	if policy.GetCapabilities().VersionControl{
		if policy.RunVersionControl(_state){
			return true
		}
	}
	//policy.Run(_state)

	fmt.Println("Working set entries: ", len(_state.WorkingSet))
	//fmt.Println("Working set entry 0: ", _state.WorkingSet[0].Data)

	return true
}

func ExecuteFromState(_state *QB_BuildState) (res bool){
	return _state.IterPipes(ExecutePolicy, nil)
}


