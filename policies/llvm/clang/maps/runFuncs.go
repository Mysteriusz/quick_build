package clang

import(
	. "qb/build"
	. "qb/policies/llvm/clang/run"
	. "qb/policies/llvm/clang/cfg"
)

var EXECUTE_FUNCS = map[string]func(*Clang_PolicyConfig, *QB_BuildState)(bool){
	"Compile": CompileExec,
	"Link": LinkExec,
}

/*	

	Execute compilation only for the QB_BuildState object`s
	of the current QB_PipeEntry

INPUT:
	NONE

OUTPUT:
	[]QB_Object with types: TYPE_FILE

*/
func CompileExec(_config *Clang_PolicyConfig, _state *QB_BuildState) (res bool){
	_state.ClearWorkingSet()

	objects, res := ClangCompileFromState(_state)
	if !res{
		return false
	}

	_state.LoadWorkingSet(objects)

	return true
}

/*	

	Execute linking only for the QB_BuildState object`s
	of the current QB_PipeEntry

INPUT:
	_state.WorkingSet types: TYPE_FILE

OUTPUT:
	_state.WorkingSet types: TYPE_FILE

*/
func LinkExec(_config *Clang_PolicyConfig, _state *QB_BuildState) (res bool){
	if _state == nil{
		return
	}
	
	out_set, res := ClangLinkFromState(_state, _config)
	if !res{
		return
	}

	_state.ClearWorkingSet()
	_state.LoadWorkingSet(out_set)

	return true
}



