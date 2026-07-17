package ar

import(
	"path/filepath"

	"qb/policies/vc"
	"qb/policies"
	"qb/build"
	"qb/qbio"
)

type Policy struct{
	/*
		Should never be modified
	*/
	PATH 		string // Has to start with '.' character
	CAPS 		policies.Capabilities

	file 		qbio.File
}

func (_policy *Policy) Run(_state *qb.BuildState) (res bool){
	if _state == nil{
		return false
	}

	res = RunFromState(_policy, _state)
	if !res{
		return
	}

	return true
}

func (_policy *Policy) GetFile() *qbio.File{
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

	_policy.file = qbio.InitFile(abs)
	return &_policy.file
}
func (_policy *Policy) GetCapabilities() policies.Capabilities{
	return _policy.CAPS
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

	objects, res := ArchiveFromState(_cfg, _state)
	if !res{
		return
	}

	_state.ClearWorkingSet()
	_state.LoadWorkingSet(objects)

	return true
}

/*

	================ VERSION CONTROL ================

*/

func (_policy *Policy) BeginVersionControl(_state *qb.BuildState) (not_first_build bool, not_updated bool, vc_state vc.FileState){
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
	not_first_build, vc_state = vc.FindOrCreateState(_state)

	/*
		AR EXCLUSIVE

		Ar ignores source/header diffs
		and only calculates diff for input/output objects
	*/
	if not_first_build{
		vc_state.DiffInput = vc.DiffObjects(vc_state.Pipe().InWorkingSet, _state.WorkingSet)
		vc_state.DiffOutput = Diff(_state, &vc_state)
	}
	no_hash_diff := (vc.StateUniqueHash(_state) == vc_state.Pipe().StateHash)

	return not_first_build, no_hash_diff && len(vc_state.DiffInput.Modified) == 0 && len(vc_state.DiffOutput.Modified) == 0, vc_state
}
func (_policy *Policy) EndVersionControl(_qb_state *qb.BuildState, _vc_state *vc.FileState){
	if _qb_state == nil || _vc_state == nil{
		return
	}

	_vc_state.Pipe().StateHash = vc.StateUniqueHash(_qb_state)

	// Save to file
	_vc_state.File.Save()
}


