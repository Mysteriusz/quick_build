package runner

import(
	"fmt"
	"time"

	"qb/build"
	"qb/policies/vc"
	"qb/policies/maps/lookups"
)

func ExecutePolicy(_state *qb.BuildState, _data any) qb.BuildError{
	/*
		Get the pipe
	*/
	pipe := _state.CurrentPipe()
	if pipe == nil{
		return qb.BuildError{}.New(
			_state,
			fmt.Sprintf("Unable to get pipe at index: %d", _state.CurrentPipeIdx()))
	}

	/*
		Lookup and execute the policy
	*/
	policy, found := lookup.PolicyLookup(pipe.CommandPolicyAlias)
	if !found{
		return qb.BuildError{}.New(
			_state,
			fmt.Sprintf("Unsupported policy file alias: %s", pipe.CommandPolicyAlias))
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
	var vc_state vc.FileState = vc.FileState{}
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

		// Link '_BuildState' to 'VC_FileState' variables
		vc.LinkState(_state, &vc_state)

		/*
			TODO
			Im not sure this will work correctly when inputs change
		*/
		vc_state.SetInputWorkingSet(_state)
	}
	
	// Exexcute the policy
	/*err = policy.Run(_state)
	if err.Check(){
		return err
	}*/

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
			Finish the version control by saving policy-specific variables
			and the version control file
		*/
		policy.EndVersionControl(_state, &vc_state)

		/*
			Set '_BuildState.WorkingSet'
			as the expected output of the pipe
		*/
		_state.LoadWorkingSet(vc_state.Pipe().OutWorkingSet)
	}

	fmt.Println("Output working set entries: ", len(_state.WorkingSet))
	fmt.Println("Policy execution took: ", end.Sub(start))
	println()

	return qb.BuildError{}.None()
}

func ExecuteFromState(_state *qb.BuildState) qb.BuildError{
	return _state.IterPipes(ExecutePolicy, nil)
}

