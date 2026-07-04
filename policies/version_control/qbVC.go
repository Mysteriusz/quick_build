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

type VC_PipeStructure struct{
	Hash 		string 		`json:"hash"`
	Iteration 	uint32 		`json:"iteration"`
	FirstBuild 	string 		`json:"first_build"`
	LastBuild 	string 		`json:"last_build"`
	InWorkingSet 	[]QB_Object 	`json:"input_working_set"`
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
*/
func VCStateUniqueHash(_state *QB_BuildState) string{
	if _state == nil{
		return ""
	}

	hash := sha256.New()
	for _, hh := range _state.GetHeaders().AllHashes() {
		io.WriteString(hash, hh)
	}
	
	for _, sh := range _state.GetSources().AllHashes() {
		io.WriteString(hash, sh)
	}

	io.WriteString(hash, _state.CurrentPipe().Command)
	io.WriteString(hash, _state.CurrentPipe().CommandPolicyAlias)
	io.WriteString(hash, _state.CurrentPipe().CommandPolicyName)

	return hex.EncodeToString(hash.Sum(nil))
}

/*
	Find/Create the version control file
*/
func VCFindOrCreateFile(_state *QB_BuildState) (vc_file VC_File){
	if _state == nil{
		return
	}

	qb_file := QBInitFile(ChangeDirectory(VC_FILE_NAME, _state.Config.OutputDirectory))
	if !qb_file.Save(){
		return
	}

	// Create the version control file object 
	vc_file = VC_File{
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

	return vc_file
}

/*
	Write a new pipe log
*/
func VCNewPipeLog(_state *QB_BuildState, _vc_file *VC_File) (res bool){
	if _state == nil{
		return
	}

	entry := VC_PipeStructure{
		Hash: VCStateUniqueHash(_state),
		Iteration: 1,
		FirstBuild: VCTimeToFormat(time.Now()),
		LastBuild: VCTimeToFormat(time.Now()),
		InWorkingSet: _state.WorkingSet,
		HeaderFiles: _state.GetHeaders(),
		SourceFiles: _state.GetSources(),
	}

	pipes := &_vc_file.Structure.Pipes
	*pipes = append(*pipes, entry)

	return true
}

/*
	Search and get reference to the pipe log 
*/
func VCSearchPipeLog(_state *QB_BuildState, _vc_file *VC_File) (found bool, _entry *VC_PipeStructure){
	if _state == nil || _vc_file == nil{
		return
	}

	hash := VCStateUniqueHash(_state)
	for _, pipe := range _vc_file.Structure.Pipes{
		if pipe.Hash == hash{
			return true, &pipe
		}
	}

	return false, nil
}

func VCTimeToFormat(_time time.Time) string{
	return _time.Format(VC_TIME_FORMAT)
}

