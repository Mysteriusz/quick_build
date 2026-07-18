package vc

import(
	"time"
	"fmt"
	"errors"
	"io"
	"encoding/json"

	"qb/qbio"
	"qb/build"
)

type SourceDiffProvider interface{
	ComputeSourceDiff(_qb_state *qb.BuildState, _vc_state *FileState) (FileDiff, qb.BuildError)
}
type HeaderDiffProvider interface{
	ComputeHeaderDiff(_qb_state *qb.BuildState, _vc_state *FileState) (FileDiff, qb.BuildError)
}
type InputDiffProvider interface{
	OutputInputDiff(_qb_state *qb.BuildState, _vc_state *FileState) (ObjectDiff, qb.BuildError)
}
type OutputDiffProvider interface{
	ComputeOutputDiff(_qb_state *qb.BuildState, _vc_state *FileState) (ObjectDiff, qb.BuildError)
}

var TIME_FORMAT string = "15:04:05 PM MST January 02/2006"
var FILE_NAME string = "VERSION_CONTROL.JSON"

type PipeIdx = uint32
type PipeStructure struct{
	/*
		Computed by 'PipeUniqueId'
	*/
	Id 		string 		`json:"id"`

	/*
		Computed by 'StateUniqueHash'
	*/
	StateHash 	string 		`json:"state_hash"`
	TimeStamp 	string 		`json:"timestamp"`
	InWorkingSet 	qb.ObjectSet 	`json:"input_working_set"`
	OutWorkingSet 	qb.ObjectSet 	`json:"output_working_set"`
	SourceFiles	qbio.FileArray	`json:"source_files"`
	HeaderFiles	qbio.FileArray	`json:"header_files"`
}

type Structure struct{
	Iteration 	uint32 	`json:"iteration"`
	FirstBuild 	string 	`json:"first_build"`
	LastBuild 	string 	`json:"last_build"`
	Pipes 		[]PipeStructure	`json:"pipes"`
}

// Version control file type
type VCFile struct{
	Structure 	Structure 
	qbio.File
}
func (_file *VCFile)PipeFromIdx(_idx PipeIdx) *PipeStructure{
	return &_file.Structure.Pipes[_idx]
}

type FileDiff struct{
	Modified qbio.FileArray // Includes both added and modified
	Removed qbio.FileArray
}
func (diff FileDiff) Len() int{
	return len(diff.Modified) + len(diff.Removed)
}

type ObjectDiff struct{
	Modified qb.ObjectSet // Includes both added and modified
	Removed qb.ObjectSet
}
func (diff ObjectDiff) Len() int{
	return len(diff.Modified) + len(diff.Removed)
}

type FileState struct{
	File		VCFile
	PipeIdx		PipeIdx
	DiffInput	ObjectDiff
	DiffOutput	ObjectDiff
	DiffSources 	FileDiff
	DiffHeaders	FileDiff
}

func (_vc_state *FileState)Pipe() *PipeStructure{
	return _vc_state.File.PipeFromIdx(_vc_state.PipeIdx)
}
func (_vc_state *FileState) Save() (res bool){
	/*
		Update source files
	*/
	_vc_state.Pipe().SourceFiles = qbio.FileArrayUnion(_vc_state.Pipe().SourceFiles, _vc_state.DiffSources.Modified)
	for _, f := range _vc_state.DiffSources.Removed{
		_vc_state.Pipe().SourceFiles.Remove(f)
	}

	/*
		Update header files
	*/
	_vc_state.Pipe().HeaderFiles = qbio.FileArrayUnion(_vc_state.Pipe().HeaderFiles, _vc_state.DiffHeaders.Modified)
	for _, f := range _vc_state.DiffHeaders.Removed{
		_vc_state.Pipe().HeaderFiles.Remove(f)
	}

	return _vc_state.File.Save()
}


func (_vc_file *VCFile) Save() (res bool){
	defer _vc_file.File.Save()
	_vc_file.Clear()

	/*
		Update version file metadata
	*/
	if _vc_file.Structure.Iteration == 0{
		_vc_file.Structure.FirstBuild = TimeToFormat(time.Now())
	}
	_vc_file.Structure.LastBuild = TimeToFormat(time.Now())
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
	Link the following values to the 'FileState'
		1) qb.BuildState.GetSources -> FileState.DiffSources.Modified
		2) qb.BuildState.GetHeaders -> FileState.DiffHeaders.Modified
*/
func LinkToBuildState(_qb_state *qb.BuildState, _vc_state *FileState){
	_qb_state.GetHeaders = func ()(qbio.FileArray){return _vc_state.DiffHeaders.Modified}
	_qb_state.GetSources = func ()(qbio.FileArray){return _vc_state.DiffSources.Modified}
}

/*
	Merge both sets so that input working set of the 
	version control has both added and updated objects stored 

	TODO:
	(Need some type of deletion control,
	so that objects that were deleted are not hanging)
*/
func (_vc_state *FileState)SetInputWorkingSet(_qb_state *qb.BuildState){
	_vc_state.Pipe().InWorkingSet.Merge(_qb_state.WorkingSet)
}

/*
	Merge both sets so that output working set of the 
	version control has both added and updated objects stored 

	TODO:
	(Need some type of deletion control,
	so that objects that were deleted are not hanging)
*/
func (_vc_state *FileState)SetOutputWorkingSet(_qb_state *qb.BuildState){
	_vc_state.Pipe().OutWorkingSet.Merge(_qb_state.WorkingSet)
	for _, o := range _vc_state.DiffOutput.Removed{
		_vc_state.Pipe().OutWorkingSet.Remove(o)
	}
}

/*
	Find/Create the version control file,
	and initialize the version control state object
*/
func InitState(_state *qb.BuildState) (not_first_build bool, vc_state FileState){
	if _state == nil{
		return
	}

	qb_file := qbio.InitFile(qbio.ChangeDirectory(FILE_NAME, _state.Config.OutputDirectory))

	// Create the version control file object 
	vc_file := VCFile{
		Structure: Structure{},
		File: qb_file,
	}

	err := json.NewDecoder(qb_file.GetFile()).Decode(&vc_file.Structure)

	// EOF means the file is empty and the initial header is to be written
	if err != nil && !errors.Is(err, io.EOF){
		fmt.Printf("Failed to decode/create the version control file:\n %s\n", qb_file.FullPath)
		fmt.Println("Error message:\n", err)
		return
	}

	existed, pipe_idx := LoadPipeLog(_state, &vc_file)
	return existed, FileState{
		File: vc_file,
		PipeIdx: pipe_idx,
	}
}

/*
	Gather all changed objects to 'FileState' object
*/
func Diff(_qb_state *qb.BuildState, _vc_state *FileState) (not_diff bool, not_crit_diff bool){
	if _qb_state == nil || _vc_state == nil{
		return
	}

	// Source files diff
	d1 := DiffFiles(_qb_state.GatherAllSources(), _vc_state.Pipe().SourceFiles)
	_vc_state.DiffSources = d1

	// Header files diff
	d2 := DiffFiles(_qb_state.GatherAllHeaders(), _vc_state.Pipe().HeaderFiles)
	_vc_state.DiffHeaders = d2

	// Input objects diff
	d3 := DiffObjects(_qb_state.WorkingSet, _vc_state.Pipe().InWorkingSet)
	_vc_state.DiffInput = d3

	/*
		Unique hash diff
		(if this if false the rebuild has to happen for the entire build)
	*/
	d4 := (StateUniqueHash(_qb_state) == _vc_state.Pipe().StateHash)

	/*
		TODO:
		Add input and output set validation via hash
	*/

	return d1.Len() == 0 && d2.Len() == 0 && d3.Len() == 0, d4
}

/*
	Write a new pipe log
*/
func NewPipeLog(_state *qb.BuildState, _vc_file *VCFile) (res bool, pipe_idx PipeIdx){
	if _state == nil{
		return
	}

	entry := PipeStructure{
		Id: PipeUniqueId(_state),
		StateHash: "",
		TimeStamp: TimeToFormat(time.Now()),
		InWorkingSet: make(qb.ObjectSet),
		OutWorkingSet: make(qb.ObjectSet),
		HeaderFiles: make(qbio.FileArray, 0),
		SourceFiles: make(qbio.FileArray, 0),
	}

	pipes := &_vc_file.Structure.Pipes
	*pipes = append(*pipes, entry)

	return true, uint32(len(*pipes) - 1)
}

/*
	Search/Create and get reference to the pipe log 
*/
func LoadPipeLog(_state *qb.BuildState, _vc_file *VCFile) (found bool, pipe_idx PipeIdx){
	if _state == nil || _vc_file == nil{
		return
	}

	/*
		Find pipe by unique state based ID
	*/
	id := PipeUniqueId(_state)
	for idx, pipe := range _vc_file.Structure.Pipes{
		if pipe.Id == id{
			return true, uint32(idx)
		}
	}

	res, idx := NewPipeLog(_state, _vc_file)
	if !res{
		panic("CLANG_POLICY: Version control failed to save the pipe log!!!")
	}

	return false, idx
}

