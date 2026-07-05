package runner

import(
	"fmt"

	. "qb/build"
	. "qb/policies/map"
	. "qb/policies/version_control"
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
		Version control variables
	*/
	var vc_enabled bool = policy.GetCapabilities().VersionControl
	var vc_state VC_FileState = VC_FileState{}
	var not_updated bool = false

	/*
		Execute a version control check of the policy info
		and ignore if it was successfull
	*/
	if vc_enabled{
		not_updated, vc_state = policy.BeginVersionControl(_state)
		if not_updated{
			vc_state.File.Save()
			/*
				Load the working set with expected output
			*/
			_state.LoadWorkingSet(vc_state.Pipe().OutWorkingSet)
			return true
		}
	}
	
	policy.Run(_state)

	/*
		Save the version control state
	*/
	if vc_enabled{
		policy.EndVersionControl(_state, &vc_state)
	}

	fmt.Println("Working set entries: ", len(_state.WorkingSet))
	//fmt.Println("Working set entry 0: ", _state.WorkingSet[0].Data)

	return true
}

func ExecuteFromState(_state *QB_BuildState) (res bool){
	return _state.IterPipes(ExecutePolicy, nil)
}


