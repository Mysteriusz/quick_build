package policies

import(
	"fmt"
	"path/filepath"

	"qb/misc"
	. "qb/policies/vc"
	. "qb/policies"
	. "qb/build"
	. "qb/io"
)

const CLANG_OUT_DEP string = "dependency_file"
const CLANG_OUT_SRC string = "source_file"

type Clang_Policy struct{
	/*
		Should never be modified
	*/
	PATH 		string // Has to start with '.' character
	CAPS 		QB_Capabilities

	file 		QB_File
}

func (_policy *Clang_Policy) Run(_state *QB_BuildState) (res bool){
	if _state == nil{
		return false
	}

	res = ClangRunFromState(_policy, _state)
	if !res{
		return
	}

	return true
}

func (_policy *Clang_Policy) GetCapabilities() QB_Capabilities{
	return _policy.CAPS
}

/*	

FIELDS:
	'Function' -> determines what execution flow the policy uses
	Acceptable 'Function' for Clang are:
		- Compile 			-> ClangCompileFromState
		- Link 				-> ClangLinkFromState

*/

type Clang_PolicyConfig struct{
	Function	string 	`toml:"function"`
	Vars		map[string]any `toml:"vars"`
}

func (_cfg *Clang_PolicyConfig)Execute(_state *QB_BuildState) bool{
	switch(_cfg.Function){
/*	

	Execute compilation only for the QB_BuildState object`s
	of the current QB_PipeEntry

INPUT:
	NONE

OUTPUT:
	[]QB_Object with types: TYPE_FILE

*/
	case "Compile":
		_state.ClearWorkingSet()

		objects, res := ClangCompileFromState(_state)
		if !res{
			return false
		}

		_state.LoadWorkingSet(objects)

		return true

/*	

	Execute linking only for the QB_BuildState object`s
	of the current QB_PipeEntry

INPUT:
	_state.WorkingSet types: TYPE_FILE

OUTPUT:
	_state.WorkingSet types: TYPE_FILE

*/
	case "Link":	
		out_set, res := ClangLinkFromState(_state, _cfg)
		if !res{
			return false
		}

		_state.ClearWorkingSet()
		_state.LoadWorkingSet(out_set)

		return true
	}

	fmt.Printf("Invalid policy function: %s\n", _cfg.Function)
	panic("")
}

/*

	================ VERSION CONTROL ================

*/

func (_policy *Clang_Policy) BeginVersionControl(_state *QB_BuildState) (not_first_build bool, not_updated bool, vc_state VC_FileState){
	if _state == nil{
		return
	}
	if !_policy.GetCapabilities().VersionControl{
		return
	}

	/*
		Load/Create version control object
		and load it`s diff
	*/
	not_first_build, vc_state = VCFindOrCreateState(_state)
	no_diff, no_crit_diff := VCDiff(_state, &vc_state)
	
	/*
		If critical diff is present,
		then act as a complete rebuild
	*/
	not_first_build = not_first_build && no_crit_diff

	src_diff := vc_state.DiffSources

	/*
		CLANG EXCLUSIVE

		Update diff source files for clang 
		(Only when the build already exist since it requires VC_FileState.OutWorkingSet)
	*/
	if not_first_build{
		src_diff = ClangVCDiff(_state, &vc_state)
		vc_state.DiffSources = QBFileArrayUnion(vc_state.DiffSources, src_diff)
	}

	if len(vc_state.DiffHeaders) > 0{
		println()
		println("==================================HDR DIFF==================================")
		misc.PrintArray(vc_state.DiffHeaders.AllPaths())
		println()
	}
	if len(vc_state.DiffSources) > 0{
		println()
		println("==================================SRC DIFF==================================")
		misc.PrintArray(vc_state.DiffSources.AllPaths())
		println()
	}

	return not_first_build, no_diff && len(src_diff) == 0, vc_state
}
func (_policy *Clang_Policy) EndVersionControl(_qb_state *QB_BuildState, _vc_state *VC_FileState){
	if _qb_state == nil{
		return
	}
	
	_vc_state.Pipe().SourceFiles = _qb_state.GatherAllSources()
	_vc_state.Pipe().HeaderFiles = _qb_state.GatherAllHeaders()
	_vc_state.Pipe().StateHash = VCStateUniqueHash(_qb_state)

	// Save to file
	_vc_state.File.Save()
}

func (_policy *Clang_Policy) GetFile() *QB_File{
	if _policy.file.IsValid(){
		return &_policy.file
	}

	if _policy.PATH[0] != '.'{
		panic("CLANG_POLICY: invalid corrupted path.")
	}

	abs, err := filepath.Abs(_policy.PATH)
	if err != nil{
		panic("CLANG_POLICY: Unable to resolve relative path.")
	}	

	_policy.file = QBInitFile(abs)
	return &_policy.file
}
