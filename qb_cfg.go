package main

import(
	"os"
)

type CfgType uint8
const(
	CFG_KM CfgType = iota
	CFG_UM
)

type FileType uint8
const(
	FILE_EXE FileType = iota 	// Windows Executable
	FILE_LIB  			// Static library
)

type PrefixDesc struct{
	SRC string
	OUT string
	INC string
	DEF string
	FLG string
}

type Entry struct{
	Label 			string 		`toml:"label"`
	BuildDirectory 		string		`toml:"build_directory"`
	BaseDirectory 		string		`toml:"base_directory"`
	SrcDirectory 		string		`toml:"src_directory"`
	VersionControl 		bool		`toml:"version_control"`
	LinkHooks 		[]string	`toml:"link_hooks"`
	Dependencies 		[]string	`toml:"dependencies"`
	Definitions 		[]string	`toml:"definitions"`
	LinkerFlags 		[]string	`toml:"linker_flags"`
	LibraryFlags 		[]string	`toml:"library_flags"`
	CompilerFlags 		[]string	`toml:"compiler_flags"`
	OutputType 		FileType	`toml:"output_type"`
	OutputBasename 		string		`toml:"output_basename"`
}
type Doc struct{
	Entries 			map[string]Entry `toml:"Entries"`
	Compiler struct{
		SourceOutExtension	string		`toml:"source_out_extension"`
		SourceExtensions 	[]string	`toml:"source_extensions"`
		CompilerPrefixGroup 	PrefixDesc	`toml:"compiler_prefix_group"`
		LinkerPrefixGroup 	PrefixDesc	`toml:"linker_prefix_group"`
		LibraryPrefixGroup 	PrefixDesc	`toml:"library_prefix_group"`
		CompilerCmd 		string		`toml:"compiler_cmd"`
		LinkerCmd 		string		`toml:"linker_cmd"`
		LibraryCmd 		string		`toml:"library_cmd"`
	} `toml:"Compiler"`
}
func entry_from_key(_doc *Doc, _key string) *Entry{
	if _doc == nil{
		panic("Invalid entry key.")
	}
	entry, exists := _doc.Entries[_key]
	if !exists{
		panic("Invalid entry key.")
	}

	return &entry
}

func cfg_load(_type CfgType, _file string) (*Doc, bool){
	file, err := os.Open(_file)
	if err != nil{
		ERR("Unable to open the config file.", err)
		return nil, false 
	}
	defer file.Close()

	var res bool
	var doc *Doc

	switch _type{
	case CFG_UM:
		doc,res = cfg_um_decode(file)
	/*
	* 		Add CFG_KM
	*/
	default:
		ERR("Unsupported config type.")
		return nil, false
	}
	
	if !res{
		return nil, false
	}
	return doc, true
}

// Create a stack of entries to build
func cfg_gather_deps(_doc *Doc, _root string, _deps map[string]bool, _deps_stack []string) ([]string, bool){
	if _doc == nil{
		ERR("Invalid Document.")
		return nil, false
	}

	entry := entry_from_key(_doc, _root)

	if _deps[_root]{ // Dependency already exists
		ERR("Dependencies are nested.")
		return nil, false
	}

	_deps[_root] = true

	// End recursion
	if len(entry.Dependencies) == 0{
		goto end
	}

	// Validate dependencies for this entry
	for _, dep_name := range entry.Dependencies{
		temp, res := cfg_gather_deps(_doc, dep_name, _deps, _deps_stack)
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
	
func cfg_build(_doc *Doc, _key string) bool{
	dep_stack, res := cfg_gather_deps(
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

