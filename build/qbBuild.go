package build

import(
	"path/filepath"
	"io/fs"
	"fmt"
	"os"

	. "qb/io"
	. "qb/configs"
)

type QB_FileGetFunc func()(QB_FileArray)
type QB_BuildState struct{
	Config 		*QB_ConfigEntry
	WorkingSet	QB_ObjectSet
	GetSources	QB_FileGetFunc
	GetHeaders	QB_FileGetFunc
	source_files 	QB_FileArray
	header_files 	QB_FileArray
	pipe_idx	QB_PipeIdx // Currently processed pipe index
}

/*
	Restores overridable functions back to default 
*/
func (_state *QB_BuildState) Restore(){
	_state.GetSources = _state.GatherAllSources
	_state.GetHeaders = _state.GatherAllHeaders
}

func QBInitBuild(_cfg *QB_ConfigEntry) (state QB_BuildState, res bool){
	if _cfg == nil{
		return
	}

	state.Restore()

	state.Config = _cfg
	return state, true
}

func (_state *QB_BuildState) PipeCount() uint8{
	return uint8(len(_state.Config.Pipeline))
}
func (_state *QB_BuildState) NextPipe() {
	_state.Restore()
	_state.pipe_idx++
}

func (_state *QB_BuildState) CurrentPipe() *QB_PipeEntry{
	return &_state.Config.Pipeline[_state.pipe_idx]
}
func (_state *QB_BuildState) CurrentPipeIdx() QB_PipeIdx{
	return _state.pipe_idx
}

func (_state *QB_BuildState) LoadWorkingSet(_objects QB_ObjectSet){
	_state.WorkingSet = _objects
}
func (_state *QB_BuildState) ClearWorkingSet(){
	_state.WorkingSet = nil
}

func (_state *QB_BuildState) GatherAllSources() QB_FileArray{
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
func (_state *QB_BuildState) GatherAllHeaders() QB_FileArray{
	if _state.header_files != nil{
		return _state.header_files
	}
	headers, res := gather_file_type(_state.Config.SourceDirectory, ".h")
	if !res{
		fmt.Println("Failed to gather header files.")
		panic("Assertion failed!!!")
	}

	_state.header_files = headers
	return headers
}

func gather_file_type(_base string, _ext string) (buf []QB_File, res bool){
	buf = make([]QB_File, 0)
	err := filepath.WalkDir(_base, func(path string, d fs.DirEntry, err error) error{
		info, err := os.Stat(path)
		if err != nil || info.IsDir(){
			return nil
		}
		
		file := QBInitFile(path)
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

type QB_IterFunction func(_state *QB_BuildState, _data any)(res bool)

func (_state *QB_BuildState) IterPipes(_func QB_IterFunction, _data any) (res bool){

	// Preserve the original pipe index
	org_idx := _state.CurrentPipeIdx()

	/*
		Iterate over all pipes
	*/
	for ; _state.CurrentPipeIdx() < _state.PipeCount(); _state.NextPipe(){
		if !_func(_state, _data){
			return
		}
	}

	_state.pipe_idx = org_idx

	return true
}
func (_state *QB_BuildState) IterPipesIdx(_from_idx QB_PipeIdx, _to_idx QB_PipeIdx,
	_func QB_IterFunction, _data any,
) (res bool){

	// Preserve the original pipe index
	org_idx := _state.CurrentPipeIdx()

	/*
		Iterate over the pipe range
	*/
	for idx := QB_PipeIdx(0); _from_idx < _to_idx; idx++{
		_state.pipe_idx = idx
		if !_func(_state, _data){
			return
		}
	}

	// Go back to the original pipe index
	_state.pipe_idx = org_idx

	return true
}

