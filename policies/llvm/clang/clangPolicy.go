package policies

import(
	"fmt"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"

	"qb/misc"
	. "qb/policies/vc"
	. "qb/policies"
	. "qb/build"
	. "qb/io"
)

const CLANG_EXTRA_FIELD_DEP string = "dependency_path"

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

	fmt.Println("==================================")
	fmt.Println("CLANG POLICY INFO")
	fmt.Println("Policy file path: ", _policy.GetFile().FullPath)
	fmt.Println("Policy file alias: ", _state.CurrentPipe().CommandPolicyAlias)
	fmt.Println("Policy name: ", _state.CurrentPipe().CommandPolicyName)
	fmt.Println("==================================")

	res = ClangRunFromState(_policy, _state)
	if !res{
		return
	}

	return true
}

func (_policy *Clang_Policy) BeginVersionControl(_state *QB_BuildState) (not_updated bool, vc_state VC_FileState){
	if _state == nil{
		return
	}
	if !_policy.GetCapabilities().VersionControl{
		return
	}

	// Load/Create version control objects required
	not_first_build, vc_state := VCFindOrCreateState(_state)
	no_diff := VCDiff(_state, &vc_state)

	println(not_first_build)
	println(no_diff)

	println("==================================SRC DIFF==================================")
	misc.PrintArray(vc_state.DiffSources.AllPaths())
	println("==================================HDR DIFF==================================")
	misc.PrintArray(vc_state.DiffHeaders.AllPaths())

	// Update diff source files for clang 
	src_diff := ClangVCDiff(_state, &vc_state)

	println("==================================CLANG SRC DIFF==================================")
	misc.PrintArray(src_diff.AllPaths())

	return not_first_build && no_diff, vc_state
}
func (_policy *Clang_Policy) EndVersionControl(_qb_state *QB_BuildState, _vc_state *VC_FileState){
	if _qb_state == nil{
		return
	}
	
	// Set the output set as the current working set
	_vc_state.Pipe().OutWorkingSet = _qb_state.WorkingSet
	_vc_state.Pipe().SourceFiles = _qb_state.GetSources()
	_vc_state.Pipe().HeaderFiles = _qb_state.GetHeaders()
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
		objects, res := ClangLinkFromState(_state)
		if !res{
			return false
		}

		_state.ClearWorkingSet()
		_state.LoadWorkingSet(objects)

		return true
	}

	fmt.Printf("Invalid policy function: %s\n", _cfg.Function)
	panic("")
}

type Clang_PolicyFile struct{
	Policies 	map[string]*Clang_PolicyConfig 	`toml:"Policies"`
}

func ClangInitPolicyFile(_policy *Clang_Policy) (cfg Clang_PolicyFile, res bool){
	if _policy == nil{
		fmt.Println("Invalid clang policy.")
		return
	}

	file := _policy.GetFile()

	err := toml.NewDecoder(file.GetFile()).Decode(&cfg)
	if err != nil{
		fmt.Printf("Failed to decode: %s\n", file.FullPath)
		return
	}

	return cfg, true
}

