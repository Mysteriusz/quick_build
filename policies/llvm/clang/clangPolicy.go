package policies

import(
	"fmt"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"

	. "qb/policies/version_control"
	. "qb/policies"
	. "qb/build"
	. "qb/io"
)

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

func (_policy *Clang_Policy) RunVersionControl(_state *QB_BuildState) (not_updated bool){
	if _state == nil{
		return
	}
	if !_policy.GetCapabilities().VersionControl{
		return
	}

	vc_file := VCFindOrCreateFile(_state)
	found, log := VCSearchPipeLog(_state, &vc_file)

	if found && log != nil{
		println("Found entry")
		println(log.Hash)

		// Return as not updated
		vc_file.Save()

		return true
	}else{
		// Create a version control pipe log and return as updated
		VCNewPipeLog(_state, &vc_file)
		vc_file.Save()

		return false
	}
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

