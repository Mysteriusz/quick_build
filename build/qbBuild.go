package qb

import(
	"path/filepath"
	"io/fs"
	"fmt"
	"os"

	"qb/qbio"
	"qb/configs"
)

type FileGetFunc func()(qbio.FileArray)
type BuildState struct{
	Config 		configs.ConfigEntry
	WorkingSet	ObjectSet
	GetSources	FileGetFunc
	GetHeaders	FileGetFunc
	source_files 	qbio.FileArray
	header_files 	qbio.FileArray
	pipe_idx	configs.PipeIdx // Currently processed pipe index
}

/*
	Restores overridable functions back to default 
*/
func (_state *BuildState) Restore(){
	_state.GetSources = _state.GatherAllSources
	_state.GetHeaders = _state.GatherAllHeaders
}

func InitBuild(_cfg configs.ConfigEntry) (BuildState, BuildError){
	state := BuildState{}
	state.Restore()

	state.Config = _cfg
	return state, BuildError{}.None()
}

func (_state *BuildState) PipeCount() uint8{
	return uint8(len(_state.Config.Pipeline))
}
func (_state *BuildState) NextPipe() {
	_state.Restore()
	_state.pipe_idx++
}

func (_state *BuildState) CurrentPipe() *configs.PipeEntry{
	return &_state.Config.Pipeline[_state.pipe_idx]
}
func (_state *BuildState) CurrentPipeIdx() configs.PipeIdx{
	return _state.pipe_idx
}

func (_state *BuildState) LoadWorkingSet(_objects ObjectSet){
	_state.WorkingSet = _objects
}
func (_state *BuildState) ClearWorkingSet(){
	_state.WorkingSet = nil
}

func (_state *BuildState) GatherAllSources() qbio.FileArray{
	if _state.source_files != nil{
		return _state.source_files
	}
	sources, res := gather_file_type(_state.Config.SourceDirectory, ".c")
	if !res{
		fmt.Println("Failed to gather source files.")
		panic("Assertion failed!!!")
	}

	_state.source_files = sources
	return sources
}
func (_state *BuildState) GatherAllHeaders() qbio.FileArray{
	if _state.header_files != nil{
		return _state.header_files
	}
	headers, res := gather_file_type(_state.Config.HeaderDirectory, ".h")
	if !res{
		fmt.Println("Failed to gather header files.")
		panic("Assertion failed!!!")
	}

	_state.header_files = headers
	return headers
}

func gather_file_type(_base string, _ext string) (buf []qbio.File, res bool){
	buf = make([]qbio.File, 0)
	err := filepath.WalkDir(_base, func(path string, d fs.DirEntry, err error) error{
		info, err := os.Stat(path)
		if err != nil || info.IsDir(){
			return nil
		}
		
		file := qbio.InitFile(path)
		if filepath.Ext(path) != _ext{
			return nil
		}
		buf = append(buf, file)

		return nil
	})
	if err != nil{
		fmt.Println("Unknown error.")
		return
	}

	return buf, true
}

type IterFunction func(_state *BuildState, _data any)(res BuildError)

func (_state *BuildState) IterPipes(_func IterFunction, _data any) BuildError{

	// Preserve the original pipe index
	org_idx := _state.CurrentPipeIdx()

	/*
		Iterate over all pipes
	*/
	for ; _state.CurrentPipeIdx() < _state.PipeCount(); _state.NextPipe(){
		err := _func(_state, _data)
		if err.Check(){
			return err
		}
	}

	_state.pipe_idx = org_idx

	return BuildError{}.None()
}
func (_state *BuildState) IterPipesIdx(_from_idx configs.PipeIdx, _to_idx configs.PipeIdx,
	_func IterFunction, _data any,
) (res BuildError){

	// Preserve the original pipe index
	org_idx := _state.CurrentPipeIdx()

	/*
		Iterate over the pipe range
	*/
	for idx := configs.PipeIdx(0); _from_idx < _to_idx; idx++{
		_state.pipe_idx = idx
		err := _func(_state, _data)
		if err.Check(){
			return err
		}
	}

	// Go back to the original pipe index
	_state.pipe_idx = org_idx

	return 
}

