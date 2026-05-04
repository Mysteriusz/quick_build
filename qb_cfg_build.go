package main

import(
	"strings"
	"path/filepath"
	"io/fs"
	"slices"
	"os"
	"errors"
	"os/exec"
)

func cfg_write_args(_doc *Doc, _key string, _buf *[]string, _args *[]string, _pref *string, _resolve bool) bool{
	if _doc == nil || _buf == nil || _args == nil || _pref == nil{
		ERR("Invalid argument.")
		return false
	}

	var res bool
	for _, arg := range *_args{
		if _resolve{
			arg, res = cfg_path_resolve(_doc, _key, arg)
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

func cfg_entry_compile(_doc *Doc, _key string, _ver *BuildVersion) ([]string, bool){
	if _doc == nil{
		return nil, false
	}

	var res bool
	var entry *Entry = entry_from_key(_doc, _key)

	/*
		Estimate amount of arguments used to minimize extending the slice
	*/
	args_count := len(entry.CompilerFlags) +
		len(entry.Definitions) +
		len(entry.LinkHooks) +
		1 + 	// SRC_PREF
		2	// OUT_PREF + OUT_PATH
	args := make([]string, 0, args_count)

	/*
		Write prefixed command arguments
	*/
	res = cfg_write_args(_doc, _key, &args, &entry.CompilerFlags, &_doc.Compiler.CompilerPrefixGroup.FLG, false)
	if !res{
		return nil, false
	}
	res = cfg_write_args(_doc, _key, &args, &entry.Definitions, &_doc.Compiler.CompilerPrefixGroup.DEF, false)
	if !res{
		return nil, false
	}
	res = cfg_write_args(_doc, _key, &args, &entry.LinkHooks, &_doc.Compiler.CompilerPrefixGroup.INC, true)
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

	out_dir, res = cfg_path_resolve(_doc, _key, entry.BuildDirectory)
	if !res{
		return nil, false
	}

	/*
		Start compiling all source files
	*/
	filepath.Walk(entry.BaseDirectory, func(path string, info fs.FileInfo, err error) error{
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

		// Does the compiled file exist
		_, err = os.Stat(out_path)
		out_path_exists := err == nil

		// Check if file was updated
		updated, res := build_version_check_and_update(_doc, _key, _ver, path)

		// No matter the outcome add file to compiled
		compiled = append(compiled, out_path)

		if !updated && out_path_exists{ // File wasnt updated and is compiled
			return nil
		}
		if !res{
			return errors.New("")
		}

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

		return nil
	})

	return compiled, true
}
func cfg_entry_process(_doc *Doc, _key string, _sources []string) (string, bool){
	if _doc == nil ||  _sources == nil{
		ERR("Invalid arguments.")
		return "", false
	}

	var res bool
	var cmd_exec string
	var cmd_args []string
	var entry *Entry = entry_from_key(_doc, _key)

	switch entry.OutputType{
	case FILE_EXE:
	case FILE_LIB:
		cmd_exec = _doc.Compiler.LibraryCmd
		cmd_args, res = cfg_lib_cmd(_doc, _key, _sources)
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
func cfg_lib_cmd(_doc *Doc, _key string, _sources []string) ([]string, bool){
	if _doc == nil || _sources == nil{
		ERR("Invalid arguments.")
		return nil, false
	}

	var res bool
	var entry *Entry = entry_from_key(_doc, _key)

	args_count := len(entry.LibraryFlags) + 
		len(_sources) * 2 + 	// len(_sources) * (SRC_PREF + SRC_PATH)
		2			// OUT_PREF + OUT_PATH
	args := make([]string, 0, args_count)

	// Create the path
	out_dir, res := cfg_path_resolve(_doc, _key, entry.BuildDirectory)
	if !res{
		return nil, false
	}
	out_path := filepath.Join(out_dir, entry.OutputBasename)

	args = append(args, _doc.Compiler.LibraryPrefixGroup.OUT)
	args = append(args, out_path)

	res = cfg_write_args(_doc, _key, &args, &entry.LibraryFlags, &_doc.Compiler.LibraryPrefixGroup.FLG, false)
	if !res{
		return nil, false
	}
	res = cfg_write_args(_doc, _key, &args, &_sources, &_doc.Compiler.LibraryPrefixGroup.SRC, false)
	if !res{
		return nil, false
	}

	println(strings.Join(args, " "))

	return args, true
}
func cfg_entry_build(_doc *Doc, _key string) bool{
	if _doc == nil{
		ERR("Invalid arguments.")
		return false
	}

	var res bool
	var ver *BuildVersion

	ver, res = build_version_load(_doc, _key)
	if !res{
		return false
	}

	sources, res := cfg_entry_compile(_doc, _key, ver)
	if !res{
		return false
	}
	_ = sources

	output, res := cfg_entry_process(_doc, _key, sources)
	if !res{
		return false
	}
	_ = output

	res = build_version_save(_doc, _key, ver)
	if !res{
		return false
	}

	return true
}

