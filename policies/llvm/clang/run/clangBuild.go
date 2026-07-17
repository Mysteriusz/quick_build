package run

import(
	"path/filepath"

	"qb/build"
	"qb/qbio"
)
func WriteArgs(_prefix string, _args[] string) (res_args []string){
	res_args = make([]string, len(_args))

	for idx,arg := range _args{
		res_args[idx] = _prefix + arg
	}

	return res_args
}

func ToFileObject(_state *qb.BuildState, _path string) (obj_path string){
	if _state == nil{
		return
	}

	outname := filepath.Base(qbio.ChangeExtension(_path, ".o"))
	outfile := _state.Config.OutputDirectory + outname
	return outfile
}

