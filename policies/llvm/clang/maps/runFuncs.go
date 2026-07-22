package maps

import(
	"qb/build"
	"qb/policies/llvm/clang/run"
	"qb/policies/llvm/clang/cfg"
)

var EXECUTE_FUNCS = map[string]func(*qb.BuildState, *clang.PolicyConfig)(qb.BuildError){
	"Compile": CompileExec,
	"Link": LinkExec,
}

/*	

	Execute compilation only for the qb.BuildState object`s
	of the current qb.PipeEntry

INPUT:
	NONE

OUTPUT:
	[]qb.Object with types: TYPE_FILE

*/
func CompileExec(_state *qb.BuildState, _config *clang.PolicyConfig) qb.BuildError{
	if _state == nil || _config == nil{
		return qb.BuildError{}.NilArgument(_state)
	}

	_state.ClearWorkingSet()

	objects, err := run.CompileFromState(_state)
	if err.Check(){
		return err
	}

	_state.LoadWorkingSet(objects)

	return qb.BuildError{}.None()
}

/*	

	Execute linking only for the qb.BuildState object`s
	of the current qb.PipeEntry

INPUT:
	_state.WorkingSet types: TYPE_FILE

OUTPUT:
	_state.WorkingSet types: TYPE_FILE

*/
func LinkExec(_state *qb.BuildState, _config *clang.PolicyConfig) qb.BuildError{
	if _state == nil || _config == nil{
		return qb.BuildError{}.NilArgument(_state)
	}
	
	out_set, err := run.LinkFromState(_state, _config)
	if err.Check(){
		return err
	}

	_state.ClearWorkingSet()
	_state.LoadWorkingSet(out_set)

	return qb.BuildError{}.None()
}



