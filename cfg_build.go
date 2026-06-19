package main

import(
	"qb/qberr"
	"qb/types"
)
/*

func cfg_entry_compile(_doc *Doc, _key string, _ver *BuildVersion) ([]string, bool){
	if _doc == nil{
		return nil, false
	}

	var res bool
	var entry *Entry = entry_from_key(_doc, _key)



	var src_field *string = &args[len(args) - 3]
	var out_field *string = &args[len(args) - 1]

	var compiled []string
	var out_dir string

	out_dir, res = cfg_path_resolve(_doc, _key, entry.BuildDirectory)
	if !res{
		return nil, false
	}

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
}*/

func cfg_entry_build(_doc *types.Doc, _key string) bool{
	if _doc == nil{
		qberr.ERR("Invalid arguments.")
		return false
	}

	var res bool
	var builder *types.EntryBuilder

	builder, res = types.CreateEntryBuilder(_doc, _key)
	if !res{
		return false
	}

	builder.Compile()
	builder.Process()
	builder.Finish()

	return true
}

// Create a stack of entries to build
func GatherConfigDeps(_doc *types.Doc, _root string, _deps map[string]bool, _deps_stack []string) ([]string, bool){
	if _doc == nil{
		qberr.ERR("Invalid Document.")
		return nil, false
	}

	entry := _doc.FromKey(_root)

	if _deps[_root]{ // Dependency already exists
		qberr.ERR("Dependencies are nested.")
		return nil, false
	}

	_deps[_root] = true

	// End recursion
	if len(entry.Dependencies) == 0{
		goto end
	}

	// Validate dependencies for this entry
	for _, dep_name := range entry.Dependencies{
		temp, res := GatherConfigDeps(_doc, dep_name, _deps, _deps_stack)
		if !res{
			return nil, false
		}
		_deps_stack = temp
	}
end:

	// Add this entry after processing all entry depenedencies
	_deps_stack = append(_deps_stack, _root)

	return _deps_stack, true
}

func BuildWithDeps(_doc *types.Doc, _key string) bool{
	dep_stack, res := GatherConfigDeps(
		_doc, _key,
		make(map[string]bool, 20),
		[]string{})
	if !res{ 
		return false
	}

	// Build all dependencies
	for _, dep_entry := range dep_stack{
		res = cfg_entry_build(_doc, dep_entry)
		if !res{
			return false
		}
	}

	return true
}
func BuildNoDeps(_doc *types.Doc, _key string) bool{
	return true
}

