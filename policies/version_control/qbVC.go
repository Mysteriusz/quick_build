package policies

import(
	"time"
	"fmt"
	"errors"
	"io"
	"encoding/json"
	"encoding/hex"
	"crypto/sha256"

	. "qb/io"
	. "qb/build"
)

var VC_TIME_FORMAT string = "15:04:05 PM MST January 02/2006"
var VC_FILE_NAME string = "VERSION_CONTROL.JSON"

type VC_PipeIdx = uint32
type VC_PipeStructure struct{
	Hash 		string 		`json:"hash"`
	TimeStamp 	string 		`json:"timestamp"`
	InWorkingSet 	[]QB_Object 	`json:"input_working_set"`
	OutWorkingSet 	[]QB_Object 	`json:"output_working_set"`
	SourceFiles	[]QB_File	`json:"source_files"`
	HeaderFiles	[]QB_File	`json:"header_files"`
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
	SourceMissing	[]QB_File
	HeaderMissing	[]QB_File
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
	Compute pipe log hash based on the state

	Computed field of 'QB_BuildState' are:
		- CurrentPipe().Command
		- CurrentPipe().CommandPolicyAlias
		- CurrentPipe().CommandPolicyName
		- GetHeaders().AllHashes()
		- GetSources().AllHashes()
		- WorkingSet
*/
func VCStateUniqueHash(_state *QB_BuildState) string{
	if _state == nil{
		return ""
	}

	hash := sha256.New()
	for _, hdr := range _state.GetHeaders().AllHashes() {
		io.WriteString(hash, hdr)
	}
	
	for _, src := range _state.GetSources().AllHashes() {
		io.WriteString(hash, src)
	}

	for _, set := range _state.WorkingSet {
		io.WriteString(hash, set.String())
	}

	io.WriteString(hash, _state.CurrentPipe().Command)
	io.WriteString(hash, _state.CurrentPipe().CommandPolicyAlias)
	io.WriteString(hash, _state.CurrentPipe().CommandPolicyName)

	return hex.EncodeToString(hash.Sum(nil))
}

/*
	Find/Create the version control file,
	and initialize the version control state object
*/
func VCFindOrCreateState(_state *QB_BuildState) (not_updated bool, vc_state VC_FileState){
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
func VCMissingCheck(_qb_state *QB_BuildState, _vc_state *VC_FileState) (not_missing bool){
	if _qb_state == nil || _vc_state == nil{
		return
	}

	pipe := _vc_state.Pipe()

	/*
		Verify input working set
	*/
	v1, m1 := VCVerifyObjects(pipe.InWorkingSet)

	/*
		Verify output working set
	*/
	v2, m2 := VCVerifyObjects(pipe.OutWorkingSet)

	/*
		Verify headers
	*/
	v3, m3 := VCVerifyFiles(pipe.HeaderFiles)

	/*
		Verify sources
	*/
	v4, m4 := VCVerifyFiles(pipe.SourceFiles)

	// Set state missing files
	_vc_state.InMissing = m1
	_vc_state.OutMissing = m2

	_vc_state.HeaderMissing = m3
	_vc_state.SourceMissing = m4

	return v1 && v2 && v3 && v4
}

/*
	Write a new pipe log
*/
func VCNewPipeLog(_state *QB_BuildState, _vc_file *VC_File,
	_in_set []QB_Object, _out_set []QB_Object,
) (res bool, pipe_idx VC_PipeIdx){
	if _state == nil{
		return
	}

	entry := VC_PipeStructure{
		Hash: VCStateUniqueHash(_state),
		TimeStamp: VCTimeToFormat(time.Now()),
		InWorkingSet: _in_set,
		OutWorkingSet: _out_set,
		HeaderFiles: _state.GetHeaders(),
		SourceFiles: _state.GetSources(),
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
	hash := VCStateUniqueHash(_state)
	for idx, pipe := range _vc_file.Structure.Pipes{
		// Return if found pipe by state hash
		if pipe.Hash == hash{
			return true, uint32(idx)
		}
	}

	res, idx := VCNewPipeLog(_state, _vc_file, nil, nil)
	if !res{
		panic("CLANG_POLICY: Version control failed to save the pipe log!!!")
	}

	return false, idx
}

/*
	Check and notify if any object doesn`t is invalid
*/
func VCVerifyObjects(_objects []QB_Object) (not_updated bool, missing []QB_Object){
	missing = make([]QB_Object, 0)
	for _,obj := range _objects{
		if !obj.Exists(){
			fmt.Printf("Object doesn`t exist:\n %s\n", obj.String())
			missing = append(missing, obj)
		}
	}
	return len(missing) == 0, missing
}
func VCVerifyFiles(_files []QB_File) (not_updated bool, missing []QB_File){
	missing = make([]QB_File, 0)
	for _,file := range _files{
		if !file.IsValid(){
			fmt.Printf("File doesn`t exist:\n %s\n", file.FullPath)
			missing = append(missing, file)
		}
	}
	return len(missing) == 0, missing
}

func VCTimeToFormat(_time time.Time) string{
	return _time.Format(VC_TIME_FORMAT)
}

