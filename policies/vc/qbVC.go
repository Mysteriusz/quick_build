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

type VC_FileState struct{
	File		VC_File
	PipeIdx		VC_PipeIdx
	InMissing	[]QB_Object
	OutMissing	[]QB_Object
	DiffSources 	QB_FileArray
	DiffHeaders	QB_FileArray
}
func (_file *VC_FileState)Pipe() *VC_PipeStructure{
	return _file.File.PipeFromIdx(_file.PipeIdx)
}

func (_vc_file *VC_File) Save() (res bool){
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

	_vc_file.QB_File.Save()

	return true
}

/*
	Link the following values to the 'VC_FileState'
		1) QB_BuildState.GetSources -> VC_FileState.DiffSources
		2) QB_BuildState.GetHeaders -> VC_FileState.DiffHeaders
*/
func VCLinkState(_vc_state *VC_FileState, _qb_state *QB_BuildState){
	_qb_state.GetHeaders = func ()(QB_FileArray){return _vc_state.DiffHeaders}
	_qb_state.GetSources = func ()(QB_FileArray){return _vc_state.DiffSources}
	_qb_state.WorkingSet = _vc_state.Pipe().InWorkingSet
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
	if !qb_file.Save(){
		return
	}

	// Create the version control file object 
	vc_file := VC_File{
		Structure: VC_Structure{},
		QB_File: qb_file,
	}

	err := json.NewDecoder(qb_file.GetFile()).Decode(&vc_file.Structure)

	// Means the file is empty and the initial header is to be written
	if errors.Is(err, io.EOF){
		vc_file.Save()
		err = json.NewDecoder(qb_file.GetFile()).Decode(&vc_file.Structure)
	}

	if err != nil{
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
		TODO:
		Hash only verifies if ANYTHING changed,
		that doesn`t include builds where a single file changed,
		which should just update the file entry (There should be something that identifies the pipe, maybe a name?)
	*/
	id := VCPipeUniqueId(_state)
	for idx, pipe := range _vc_file.Structure.Pipes{
		if pipe.Id == id{
			return true, uint32(idx)
		}
	}
	/*hash := VCStateUniqueHash(_state)
	for idx, pipe := range _vc_file.Structure.Pipes{
		// Return if found pipe by state hash
		if pipe.Hash == hash{
			return true, uint32(idx)
		}
	}*/


	res, idx := VCNewPipeLog(_state, _vc_file)
	if !res{
		panic("CLANG_POLICY: Version control failed to save the pipe log!!!")
	}

	return false, idx
}

