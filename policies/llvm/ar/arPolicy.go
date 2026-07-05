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

type Ar_Policy struct{
	/*
		Should never be modified
	*/
	PATH 		string // Has to start with '.' character
	CAPS 		QB_Capabilities

	file 		QB_File
}

func (_policy *Ar_Policy) Run(_state *QB_BuildState) (res bool){
	if _state == nil{
		return false
	}
	fmt.Println("==================================")
	fmt.Println("CLANG POLICY INFO")
	fmt.Println("Policy file path: ", _policy.GetFile().FullPath)
	fmt.Println("Policy file alias: ", _state.CurrentPipe().CommandPolicyAlias)
	fmt.Println("Policy name: ", _state.CurrentPipe().CommandPolicyName)
	fmt.Println("==================================")

	res = ArRunFromState(_policy, _state)
	if !res{
		return
	}

	return true
}

func (_policy *Ar_Policy) BeginVersionControl(_state *QB_BuildState) (not_updated bool, vc_state VC_FileState){
	if _state == nil{
		return
	}
	return
}
func (_policy *Ar_Policy) EndVersionControl(_state *QB_BuildState, _vc_state *VC_FileState){
	if _state == nil || _vc_state == nil{
		return
	}
}

func (_policy *Ar_Policy) GetFile() *QB_File{
	if _policy.file.IsValid(){
		return &_policy.file
	}

	if _policy.PATH[0] != '.'{
		panic("AR_POLICY: invalid corrupted path.")
	}

	abs, err := filepath.Abs(_policy.PATH)
	if err != nil{
		panic("AR_POLICY: Unable to resolve relative path.")
	}	

	_policy.file = QBInitFile(abs)
	return &_policy.file
}
func (_policy *Ar_Policy) GetCapabilities() QB_Capabilities{
	return _policy.CAPS
}

/*	

FIELDS:
	'Mode' -> ar-compatbile mode ex: rcs
	'OutputExt' -> ar-compatbile extension of the output archive
	'OutputName' -> ar-compatbile name of the output archive

*/
type Ar_PolicyConfig struct{
	Mode		string 	`toml:"mode"`
	OutputExt	string 	`toml:"output_ext"`
	OutputName	string 	`toml:"output_name"`
}

func (_cfg *Ar_PolicyConfig)Execute(_state *QB_BuildState) (res bool){

/*	

	Execute archive creation only for the QB_BuildState object`s

INPUT:
	_state.WorkingSet types: TYPE_FILE
	_cfg.Mode
	_cfg.OutputExt
	_cfg.OutputName

OUTPUT:
	[]QB_Object with types: TYPE_FILE

*/

	objects, res := ArArchiveFromState(_cfg, _state)
	if !res{
		return
	}

	_state.ClearWorkingSet()
	_state.LoadWorkingSet(objects)

	return true
}

type Ar_PolicyFile struct{
	Policies 	map[string]*Ar_PolicyConfig 	`toml:"Policies"`
}

func ArInitPolicyFile(_policy *Ar_Policy) (cfg Ar_PolicyFile, res bool){
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

