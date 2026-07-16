package policies

import(
	"time"
	"fmt"
	"errors"
	"io"
	"encoding/json"

	. "qb/io"
	. "qb/build"
)

type VC_PolicyInt interface{
	/*
		Begin policy defined version check 
		transaction on the build state object

		The function should but isn`t required to check
		for Version control capability of it`s policy info

		IMPORTANT!
		This should require 'GetCapabilities' to return an object
		with field 'QB_Capabilities.VersionControl' == true

		Else the not_updated value should always be 0 (false)
	*/
	BeginVersionControl(_state *QB_BuildState) (not_first_build bool, not_updated bool, _vc_state VC_FileState)
	/*
		Finish and save the version check transaction

		The function should but isn`t required to check
		for Version control capability of it`s policy info

		IMPORTANT!
		This should require 'GetCapabilities' to return an object
		with field 'QB_Capabilities.VersionControl' == true
	*/
	EndVersionControl(_state *QB_BuildState, _vc_state *VC_FileState)
}

var VC_TIME_FORMAT string = "15:04:05 PM MST January 02/2006"
var VC_FILE_NAME string = "VERSION_CONTROL.JSON"

type VC_PipeIdx = uint32
type VC_PipeStructure struct{
	/*
		Computed by 'VCPipeUniqueId'
	*/
	Id 		string 		`json:"id"`

	/*
		Computed by 'VCStateUniqueHash'
	*/
	StateHash 	string 		`json:"state_hash"`
	TimeStamp 	string 		`json:"timestamp"`
	InWorkingSet 	QB_ObjectSet 	`json:"input_working_set"`
	OutWorkingSet 	QB_ObjectSet 	`json:"output_working_set"`
	SourceFiles	QB_FileArray	`json:"source_files"`
	HeaderFiles	QB_FileArray	`json:"header_files"`
}

type VC_Structure struct{
	Iteration 	uint32 	`json:"iteration"`
	FirstBuild 	string 	`json:"first_build"`
	LastBuild 	string 	`json:"last_build"`
	Pipes 		[]VC_PipeStructure	`json:"pipes"`
}

// Version control file type
type VC_File struct{
	Structure 	VC_Structure 
	QB_File
}
func (_file *VC_File)PipeFromIdx(_idx VC_PipeIdx) *VC_PipeStructure{
	return &_file.Structure.Pipes[_idx]
}

type VC_FileDiff struct{
	Modified QB_FileArray // Includes both added and modified
	Removed QB_FileArray
}
func (diff VC_FileDiff) Len() int{
	return len(diff.Modified) + len(diff.Removed)
}

type VC_ObjectDiff struct{
	Modified QB_ObjectSet // Includes both added and modified
	Removed QB_ObjectSet
}
func (diff VC_ObjectDiff) Len() int{
	return len(diff.Modified) + len(diff.Removed)
}

type VC_FileState struct{
	File		VC_File
	PipeIdx		VC_PipeIdx
	DiffInput	VC_ObjectDiff
	DiffOutput	VC_ObjectDiff
	DiffSources 	VC_FileDiff
	DiffHeaders	VC_FileDiff
}

func (_file *VC_FileState)Pipe() *VC_PipeStructure{
	return _file.File.PipeFromIdx(_file.PipeIdx)
}

func (_vc_file *VC_File) Save() (res bool){
	defer _vc_file.QB_File.Save()

	_vc_file.Clear()

	/*
		Update version file metadata
	*/
	if _vc_file.Structure.Iteration == 0{
		_vc_file.Structure.FirstBuild = VCTimeToFormat(time.Now())
	}
	_vc_file.Structure.LastBuild = VCTimeToFormat(time.Now())
	_vc_file.Structure.Iteration++

	enc := json.NewEncoder(_vc_file.GetFile())
	enc.SetIndent("", "  ")
	err := enc.Encode(&_vc_file.Structure)

	if err != nil{
		fmt.Printf("Failed to encode the version control file:\n %s\n", _vc_file.FullPath)
		fmt.Println("Error message:\n", err)
		return
	}

	return true
}

/*
	Link the following values to the 'VC_FileState'
		1) QB_BuildState.GetSources -> VC_FileState.DiffSources.Modified
		2) QB_BuildState.GetHeaders -> VC_FileState.DiffHeaders.Modified
*/
func VCLinkState(_qb_state *QB_BuildState, _vc_state *VC_FileState){
	_qb_state.GetHeaders = func ()(QB_FileArray){return _vc_state.DiffHeaders.Modified}
	_qb_state.GetSources = func ()(QB_FileArray){return _vc_state.DiffSources.Modified}
}

/*
	Merge both sets so that input working set of the 
	version control has both added and updated objects stored 

	TODO:
	(Need some type of deletion control,
	so that objects that were deleted are not hanging)
*/
func (_vc_state *VC_FileState)SetInputWorkingSet(_qb_state *QB_BuildState){
	_vc_state.Pipe().InWorkingSet.Merge(_qb_state.WorkingSet)
}

/*
	Merge both sets so that output working set of the 
	version control has both added and updated objects stored 

	TODO:
	(Need some type of deletion control,
	so that objects that were deleted are not hanging)
*/
func (_vc_state *VC_FileState)SetOutputWorkingSet(_qb_state *QB_BuildState){
	_vc_state.Pipe().OutWorkingSet.Merge(_qb_state.WorkingSet)
	for _, o := range _vc_state.DiffOutput.Removed{
		_vc_state.Pipe().OutWorkingSet.Remove(o)
	}
}

/*
	Find/Create the version control file,
	and initialize the version control state object
*/
func VCFindOrCreateState(_state *QB_BuildState) (not_first_build bool, vc_state VC_FileState){
	if _state == nil{
		return
	}

	qb_file := QBInitFile(ChangeDirectory(VC_FILE_NAME, _state.Config.OutputDirectory))

	// Create the version control file object 
	vc_file := VC_File{
		Structure: VC_Structure{},
		QB_File: qb_file,
	}

	err := json.NewDecoder(qb_file.GetFile()).Decode(&vc_file.Structure)

	// EOF means the file is empty and the initial header is to be written
	if err != nil && !errors.Is(err, io.EOF){
		fmt.Printf("Failed to decode/create the version control file:\n %s\n", qb_file.FullPath)
		fmt.Println("Error message:\n", err)
		return
	}

	existed, pipe_idx := VCLoadPipeLog(_state, &vc_file)
	return existed, VC_FileState{
		File: vc_file,
		PipeIdx: pipe_idx,
	}
}

/*
	Gather all changed objects to 'VC_FileState' object
*/
func VCDiff(_qb_state *QB_BuildState, _vc_state *VC_FileState) (not_diff bool, not_crit_diff bool){
	if _qb_state == nil || _vc_state == nil{
		return
	}

	// Source files diff
	d1 := VCDiffFiles(_qb_state.GatherAllSources(), _vc_state.Pipe().SourceFiles)
	_vc_state.DiffSources = d1

	// Header files diff
	d2 := VCDiffFiles(_qb_state.GatherAllHeaders(), _vc_state.Pipe().HeaderFiles)
	_vc_state.DiffHeaders = d2

	// Input objects diff
	d3 := VCDiffObjects(_qb_state.WorkingSet, _vc_state.Pipe().InWorkingSet)
	_vc_state.DiffInput = d3

	/*
		Unique hash diff
		(if this if false the rebuild has to happen for the entire build)
	*/
	d4 := (VCStateUniqueHash(_qb_state) == _vc_state.Pipe().StateHash)

	/*
		TODO:
		Add input and output set validation via hash
	*/

	return d1.Len() == 0 && d2.Len() == 0 && d3.Len() == 0, d4
}

/*
	Write a new pipe log
*/
func VCNewPipeLog(_state *QB_BuildState, _vc_file *VC_File) (res bool, pipe_idx VC_PipeIdx){
	if _state == nil{
		return
	}

	entry := VC_PipeStructure{
		Id: VCPipeUniqueId(_state),
		StateHash: "",
		TimeStamp: VCTimeToFormat(time.Now()),
		InWorkingSet: make(QB_ObjectSet),
		OutWorkingSet: make(QB_ObjectSet),
		HeaderFiles: make(QB_FileArray, 0),
		SourceFiles: make(QB_FileArray, 0),
	}

	pipes := &_vc_file.Structure.Pipes
	*pipes = append(*pipes, entry)

	return true, uint32(len(*pipes) - 1)
}

/*
	Search/Create and get reference to the pipe log 
*/
func VCLoadPipeLog(_state *QB_BuildState, _vc_file *VC_File) (found bool, pipe_idx VC_PipeIdx){
	if _state == nil || _vc_file == nil{
		return
	}

	/*
		Find pipe by unique state based ID
	*/
	id := VCPipeUniqueId(_state)
	for idx, pipe := range _vc_file.Structure.Pipes{
		if pipe.Id == id{
			return true, uint32(idx)
		}
	}

	res, idx := VCNewPipeLog(_state, _vc_file)
	if !res{
		panic("CLANG_POLICY: Version control failed to save the pipe log!!!")
	}

	return false, idx
}

