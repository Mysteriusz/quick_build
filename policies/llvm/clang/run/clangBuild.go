package run

import(
	"path/filepath"

	. "qb/build"
	. "qb/io"
)
func ClangWriteArgs(_prefix string, _args[] string) (res_args []string){
	res_args = make([]string, len(_args))

	for idx,arg := range _args{
		res_args[idx] = _prefix + arg
	}

	return res_args
}

func ClangToFileObject(_state *QB_BuildState, _path string) (obj_path string){
	if _state == nil{
		return
	}

	outname := filepath.Base(ChangeExtension(_path, ".o"))
	outfile := _state.Config.OutputDirectory + outname
	return outfile
}

