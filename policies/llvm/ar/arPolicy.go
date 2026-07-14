package policies

import(
	"path/filepath"

	. "qb/policies/vc"
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

	res = ArRunFromState(_policy, _state)
	if !res{
		return
	}

	return true
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
	_state.WorkingSet with types: TYPE_FILE
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

/*

	================ VERSION CONTROL ================

*/

func (_policy *Ar_Policy) BeginVersionControl(_state *QB_BuildState) (not_first_build bool, not_updated bool, vc_state VC_FileState){
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

	/*
		AR EXCLUSIVE

		Ar ignores source/header diffs
		and only calculates diff for input/output objects
	*/
	in_diff := vc_state.DiffInput
	out_diff := vc_state.DiffOutput
	if not_first_build{
		in_diff = VCDiffObjects(vc_state.Pipe().InWorkingSet, _state.WorkingSet)
		vc_state.DiffInput = in_diff

		out_diff = ArVCDiff(_state, &vc_state)
		vc_state.DiffOutput = out_diff
	}
	no_hash_diff := (VCStateUniqueHash(_state) == vc_state.Pipe().StateHash)

	/*
	println(len(in_diff))
	println(len(out_diff))
	*/

	return not_first_build, no_hash_diff && len(in_diff) == 0 && len(out_diff) == 0, vc_state
}
func (_policy *Ar_Policy) EndVersionControl(_qb_state *QB_BuildState, _vc_state *VC_FileState){
	if _qb_state == nil || _vc_state == nil{
		return
	}

	_vc_state.Pipe().StateHash = VCStateUniqueHash(_qb_state)

	// Save to file
	_vc_state.File.Save()
}


