package maps

import(
	"qb/build"
	"qb/policies/llvm/clang/run"
	"qb/policies/llvm/clang/cfg"
)

var EXECUTE_FUNCS = map[string]func(*clang.PolicyConfig, *qb.BuildState)(bool){
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
func CompileExec(_config *clang.PolicyConfig, _state *qb.BuildState) (res bool){
	_state.ClearWorkingSet()

	objects, res := run.CompileFromState(_state)
	if !res{
		return false
	}

	_state.LoadWorkingSet(objects)

	return true
}

/*	

	Execute linking only for the qb.BuildState object`s
	of the current qb.PipeEntry

INPUT:
	_state.WorkingSet types: TYPE_FILE

OUTPUT:
	_state.WorkingSet types: TYPE_FILE

*/
func LinkExec(_config *clang.PolicyConfig, _state *qb.BuildState) (res bool){
	if _state == nil{
		return
	}
	
	out_set, res := run.LinkFromState(_state, _config)
	if !res{
		return
	}

	_state.ClearWorkingSet()
	_state.LoadWorkingSet(out_set)

	return true
}



