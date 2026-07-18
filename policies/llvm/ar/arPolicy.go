package ar

import(
	"qb/build"
	"qb/policies"
)

type PolicyInfo struct{
	base 		policies.PolicyFile
	config		*PolicyConfig
}

const POLICY_FILE_PATH string = "./policies/llvm/ar.toml"

func (_policy *PolicyInfo) GetCapabilities() policies.Capabilities{
	return policies.Capabilities{
		VersionControl: true,
	}
}
func (_policy *PolicyInfo) GetFile() *policies.PolicyFile{
	file, res := policies.LoadPolicyFile(POLICY_FILE_PATH)
	if !res{
		return nil
	}
	return &file
}
func (_policy *PolicyInfo) Run(_state *qb.BuildState) qb.BuildError{
	if _state == nil{
		return qb.BuildError{}.NilArgument(_state)
	}

	/*res = RunFromState(_policy, _state)
	if !res{
		return
	}*/

	return qb.BuildError{}.None()
}

/*	

FIELDS:
	'Mode' -> ar-compatbile mode ex: rcs
	'OutputExt' -> ar-compatbile extension of the output archive
	'OutputName' -> ar-compatbile name of the output archive

*/
type PolicyConfig struct{
	Mode		string 	`toml:"mode"`
	OutputExt	string 	`toml:"output_ext"`
	OutputName	string 	`toml:"output_name"`
}

func (_cfg *PolicyConfig)Execute(_state *qb.BuildState) (res bool){

/*	

	Execute archive creation only for the qb.BuildState object`s

INPUT:
	_state.WorkingSet with types: TYPE_FILE
	_cfg.Mode
	_cfg.OutputExt
	_cfg.OutputName

OUTPUT:
	[]qb.Object with types: TYPE_FILE

*/

	/*objects, res := ArchiveFromState(_cfg, _state)
	if !res{
		return
	}

	_state.ClearWorkingSet()
	_state.LoadWorkingSet(objects)*/

	return true
}

/*

	================ VERSION CONTROL ================

*/

