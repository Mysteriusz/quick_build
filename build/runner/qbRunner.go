package runner

import(
	"fmt"
	"time"

	. "qb/build"
	. "qb/policies/maps/lookups"
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
	println()

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
		/*
 			Perform the version control scan

			Check for diffs, and their importance
			and return the information and the VC_FileState

			Load the diffs to the 'VC_FileState' object
			to allow further vc setups before running
		*/
		not_first_build, not_updated, vc_state = policy.BeginVersionControl(_state)

		/*
 			Force a full rebuild
		*/
		if pipe.AlwaysRebuild{
			not_first_build = false
		}

		/*
			Move finish the current build
			(nothing has changed)
		*/
		if not_first_build && not_updated{
			vc_state.File.Save()
			goto timer_end
		}

		// Link 'QB_BuildState' to 'VC_FileState' variables
		VCLinkState(_state, &vc_state)

		/*
			TODO
			Im not sure this will work correctly when inputs change
		*/
		vc_state.SetInputWorkingSet(_state)
	}
	
	// Exexcute the policy
	res = policy.Run(_state)
	if !res{
		fmt.Println("Policy execution has failed")
		return
	}

	/*
		Save the version control state
	*/
	if vc_enabled{
		/*
			Merge expected output with output that was produced

			This ensures that even when 2 new files were added
			and 30 were already in the expected output
			32 are actually saved
		*/
		vc_state.SetOutputWorkingSet(_state)
	}

	// End execution timer
timer_end:
	end := time.Now()

	/*
		Load the working set with expected output
	*/
	if vc_enabled{
		/*
			Set 'QB_BuildState.WorkingSet'
			as the expected output of the pipe
		*/
		_state.LoadWorkingSet(vc_state.Pipe().OutWorkingSet)

		/*
			Finish the version control by saving policy-specific variables
			and the version control file
		*/
		policy.EndVersionControl(_state, &vc_state)
	}

	fmt.Println("Output working set entries: ", len(_state.WorkingSet))
	fmt.Println("Policy execution took: ", end.Sub(start))
	println()

	return true
}

func ExecuteFromState(_state *QB_BuildState) (res bool){
	return _state.IterPipes(ExecutePolicy, nil)
}

