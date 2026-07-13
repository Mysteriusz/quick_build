package runner

import(
	"fmt"
	"time"

	. "qb/build"
	. "qb/policies/map"
	. "qb/policies/vc"
)

func ExecutePolicy(_state *QB_BuildState, _data any) (res bool){
	/*
		Get the pipe
	*/
	pipe := _state.CurrentPipe()
	if pipe == nil{
		fmt.Println("Unable to get pipe at index: %i", _state.CurrentPipeIdx())
		return
	}

	/*
		Lookup and execute the policy
	*/
	policy, res := QBPolicyLookup(pipe.CommandPolicyAlias)
	if !res{
		return
	}

	fmt.Println("==================================")
	fmt.Println("POLICY INFO")
	fmt.Println("Policy file path: ", policy.GetFile().FullPath)
	fmt.Println("Policy file alias: ", _state.CurrentPipe().CommandPolicyAlias)
	fmt.Println("Policy name: ", _state.CurrentPipe().CommandPolicyName)
	fmt.Println("==================================")

	fmt.Println("Input working set entries: ", len(_state.WorkingSet))
	
	// Start execution timer
	start := time.Now()

	/*
		Version control variables
	*/
	var vc_enabled bool = policy.GetCapabilities().VersionControl
	var vc_state VC_FileState = VC_FileState{}
	var not_first_build, not_updated bool

	/*
		Execute a version control check of the policy info
		and ignore if it was successfull
	*/
	if vc_enabled{
		not_first_build, not_updated, vc_state = policy.BeginVersionControl(_state)
		if not_first_build && not_updated{
			vc_state.File.Save()
			goto timer_end
		}

		if not_first_build{
			VCLinkState(_state, &vc_state)
		}

		vc_state.SetInputWorkingSet(_state)
	}
	
	policy.Run(_state)

	/*
		Save the version control state
	*/
	if vc_enabled{
		vc_state.SetOutputWorkingSet(_state)

		policy.EndVersionControl(_state, &vc_state)
	}


	// End execution timer
timer_end:
	end := time.Now()

	/*
		Load the working set with expected output
	*/
	if vc_enabled{
		_state.LoadWorkingSet(vc_state.Pipe().OutWorkingSet)
	}

	fmt.Println("Output working set entries: ", len(_state.WorkingSet))
	fmt.Println("Policy execution took: ", end.Sub(start))

	return true
}

func ExecuteFromState(_state *QB_BuildState) (res bool){
	return _state.IterPipes(ExecutePolicy, nil)
}

