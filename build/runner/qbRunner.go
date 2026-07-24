package runner

import(
	"fmt"
	"time"

	"qb/build"
	"qb/policies"
	"qb/policies/vc"
	"qb/policies/maps"
)

/*
	IMPORTANT!
	Diff compute happens in the following order

	1) vc.HeaderDiffProvider
	2) vc.SourceDiffProvider
	3) vc.InputDiffProvider
	4) vc.OutputDiffProvider

	Executes all providers that mutate 'vc.FileState'
	based on how they operate	
*/
func RunVCProviders(
	_policy policies.PolicyInfoInt,
	_qb_state *qb.BuildState,
	_vc_state *vc.FileState,
) qb.BuildError{
	if _policy == nil || _vc_state == nil || _qb_state == nil{
		return qb.BuildError{}.NilArgument(_qb_state)
	}	

	/*
		Execute the header diff provider
	*/
	if prov, use := _policy.(vc.HeaderDiffProvider); use{
		diff, err := prov.ComputeHeaderDiff(_qb_state, _vc_state)
		if err.Check(){
			return err
		}
		_vc_state.DiffHeaders = diff
	}

	/*
		Execute the source diff provider
	*/
	if prov, use := _policy.(vc.SourceDiffProvider); use{
		diff, err := prov.ComputeSourceDiff(_qb_state, _vc_state)
		if err.Check(){
			return err
		}
		_vc_state.DiffSources = diff
	}

	/*
		Execute the input diff provider
	*/
	if prov, use := _policy.(vc.InputDiffProvider); use{
		diff, err := prov.ComputeInputDiff(_qb_state, _vc_state)
		if err.Check(){
			return err
		}
		_vc_state.DiffInput = diff
	}
	
	/*
		Execute the output diff provider
	*/
	if prov, use := _policy.(vc.OutputDiffProvider); use{
		diff, err := prov.ComputeOutputDiff(_qb_state, _vc_state)
		if err.Check(){
			return err
		}
		_vc_state.DiffOutput = diff
	}
	println("Header Diff")
	println("REMOVED")
	println(len(_vc_state.DiffHeaders.Removed))
	println("MODIFIED")
	println(len(_vc_state.DiffHeaders.Modified))
	println()
	println("Source Diff")
	println("REMOVED")
	println(len(_vc_state.DiffSources.Removed))
	println("MODIFIED")
	println(len(_vc_state.DiffSources.Modified))
	println()
	println("Object Diff")
	println("REMOVED")
	println(len(_vc_state.DiffOutput.Removed))
	println("MODIFIED")
	println(len(_vc_state.DiffOutput.Modified))

	return qb.BuildError{}.None()
}

func ExecutePolicy(_state *qb.BuildState, _data any) qb.BuildError{
	if _state == nil{
		return qb.BuildError{}.NilArgument(_state)
	}

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
	policy, found := maps.POLICY_INFO_LOOKUP[pipe.CommandPolicyAlias]
	if !found{
		return qb.BuildError{}.New(
			_state,
			fmt.Sprintf("Unsupported policy file alias: %s", pipe.CommandPolicyAlias))
	}

	fmt.Println("==================================")
	fmt.Println("POLICY INFO")
	fmt.Println("Policy file path: ", policy.GetFile().File.FullPath)
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
	var err qb.BuildError
	var vc_enabled bool = policy.GetCapabilities().VersionControl
	var vc_state vc.FileState

	var not_first_build bool

	if vc_enabled{
		/*
			Initialize the FileState
		*/
		not_first_build, vc_state = vc.InitState(_state)

		/*
			Set as first build if rebuild requested
			or state hash has changed
		*/
		if pipe.AlwaysRebuild || vc.StateUniqueHash(_state) != vc_state.Pipe().StateHash{
			goto rebuild
		}

		err := RunVCProviders(policy, _state, &vc_state)
		if err.Check(){
			return err
		}

		diff_len := vc_state.DiffSources.Len() +
			vc_state.DiffHeaders.Len() +
			vc_state.DiffInput.Len() +
			vc_state.DiffOutput.Len()
		/*
			If there are no diffs nor is the build first
			act as if the build had finished
		*/
		if diff_len == 0 && not_first_build{
			// Clear working set since the policy was not ran
			_state.ClearWorkingSet()
			goto timer_end
		}

		vc.LinkToBuildState(_state, &vc_state)
	}

rebuild:
	err = policy.Run(_state)
	if err.Check(){
		return err
	}

timer_end:
	if vc_enabled{
		vc_state.SaveWithState(_state)
	}

	// End execution timer
	end := time.Now()

	fmt.Println("Output working set entries: ", len(_state.WorkingSet))
	fmt.Println("Policy execution took: ", end.Sub(start))
	println()

	return qb.BuildError{}.None()
}

func ExecuteFromState(_state *qb.BuildState) qb.BuildError{
	return _state.IterPipes(ExecutePolicy, nil)
}

