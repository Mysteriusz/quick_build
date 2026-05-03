package main

import(
	"strings"
	"path/filepath"
	"io/fs"
	"slices"
	"os"
	"os/exec"
)

func cfg_write_args(_doc *Doc, _entry *Entry, _buf *[]string, _args *[]string, _pref *string, _resolve bool) bool{
	if _buf == nil || _args == nil || _pref == nil{
		ERR("Invalid argument.")
		return false
	}

	var res bool
	for _, arg := range *_args{
		if _resolve{
			arg, res = cfg_path_resolve(_doc, _entry, arg)
			if !res{
				return false
			}
		}
		*_buf = append(*_buf, *_pref + arg)
	}

	return true
}

/*func cfg_write_definitions(_doc *Doc, _entry *Entry, _args *[]string) bool{
	// Write definitions
	for _, def := range _entry.Definitions{
		*_args = append(*_args, _doc.Compiler.CompilePrefixGroup.DEF + def)
	}

	return true
}
func cfg_write_includes(_doc *Doc, _entry *Entry, _args *[]string) bool{
	// Write include (hook) directories
	for _, hook := range _entry.LinkHooks{
		fullpath, res := cfg_path_resolve(_doc, _entry, hook)
		if !res{
			return false
		}

		*_args = append(*_args, _doc.Compiler.CompilePrefixGroup.INC + fullpath)
	}

	return true
}
func cfg_write_flags(_doc *Doc, _entry *Entry, _args *[]string) bool{
	for _, flag := range _entry.Flags{
		*_args = append(*_args, _doc.Compiler.CompilePrefixGroup.FLG + flag)
	}

	return true
}*/

func cfg_entry_compile(_doc *Doc, _entry *Entry) ([]string, bool){
	var res bool

	/*
		Estimate amount of arguments used to minimize extending the slice
	*/
	args_count := len(_entry.CompilerFlags) +
		len(_entry.Definitions) +
		len(_entry.LinkHooks) +
		1 + 	// SRC_PREF
		2	// OUT_PREF + OUT_PATH
	args := make([]string, 0, args_count)

	/*
		Write prefixed command arguments
	*/
	res = cfg_write_args(_doc, _entry, &args, &_entry.CompilerFlags, &_doc.Compiler.CompilerPrefixGroup.FLG, false)
	if !res{
		return nil, false
	}
	res = cfg_write_args(_doc, _entry, &args, &_entry.Definitions, &_doc.Compiler.CompilerPrefixGroup.DEF, false)
	if !res{
		return nil, false
	}
	res = cfg_write_args(_doc, _entry, &args, &_entry.LinkHooks, &_doc.Compiler.CompilerPrefixGroup.INC, true)
	if !res{
		return nil, false
	}

	args = append(args, _doc.Compiler.CompilerPrefixGroup.SRC)
	args = append(args, "")
	args = append(args, _doc.Compiler.CompilerPrefixGroup.OUT)
	args = append(args, "")

	/*
		{CMD} \
			{CompilePrefixGroup.SRC} 	{src_field}
			{CompilePrefixGroup.OUT} 	{out_field}
	*/
	var src_field *string = &args[len(args) - 3]
	var out_field *string = &args[len(args) - 1]

	var compiled []string
	var out_dir string

	out_dir, res = cfg_path_resolve(_doc, _entry, _entry.BuildDirectory)
	if !res{
		return nil, false
	}

	/*
		Start compiling all source files
	*/
	filepath.Walk(_entry.BaseDirectory, func(path string, info fs.FileInfo, err error) error{
		if err != nil{
			ERR("An error occured.", err)
			return err
		}

		// Skip file that`s not a source
		if !slices.Contains(_doc.Compiler.SourceExtensions, filepath.Ext(path)){
			return nil
		}

		out_path := filepath.Join(
			out_dir,
			info.Name() + _doc.Compiler.SourceOutExtension)

		*src_field = path
		*out_field = out_path
		println(strings.Join(args, " "))

		/*
			Form a command from builders and paths
		*/
		cmd := exec.Command(_doc.Compiler.CompilerCmd, args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()

		compiled = append(compiled, out_path)

		return nil
	})

	return compiled, true
}
func cfg_entry_link(_doc *Doc, _entry *Entry, _sources []string) (string, bool){
	if _doc == nil || _entry == nil || _sources == nil{
		ERR("Invalid arguments.")
		return "", false
	}

	var res bool
	var cmd_exec string
	var cmd_args []string

	switch _entry.OutputType{
	case FILE_EXE:
	case FILE_LIB:
		cmd_exec = _doc.Compiler.LibraryCmd
		cmd_args, res = cfg_lib_cmd(_doc, _entry, _sources)
		if !res{
			return "", false
		}
	default:
		ERR("Unknown output type.")
		return "", false
	}

	/*
		Form a command from builders and paths
	*/
	cmd := exec.Command(cmd_exec, cmd_args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()

	return "", true
}
func cfg_entry_build(_doc *Doc, _entry *Entry) bool{
	if _doc == nil || _entry == nil{
		ERR("Invalid arguments.")
		return false
	}

	sources, res := cfg_entry_compile(_doc, _entry)
	if !res{
		return false
	}

	linked, res := cfg_entry_link(_doc, _entry, sources)
	if !res{
		return false
	}
	_ = linked

	return true
}

