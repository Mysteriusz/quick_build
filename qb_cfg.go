package main

import (
	"os"
)

type CfgType int
const (
	CFG_KM = iota
	CFG_UM
)

type Entry struct{
	Label 			string 		`toml:"label"`
	BuildDirectory 		string		`toml:"build_directory"`
	BaseDirectory 		string		`toml:"base_directory"`
	VersionControl 		bool		`toml:"version_control"`
	LinkHooks 		[]string	`toml:"link_hooks"`
	Dependencies 		[]string	`toml:"dependencies"`
	Definitions 		[]string	`toml:"definitions"`
	Flags 			[]string	`toml:"flags"`
}
type Doc struct{
	Build struct{
		Entries 	map[string]Entry `toml:"Entries"`
	} `toml:"Build"`
	Compiler struct{
		SourceOutExtension	string		`toml:"source_out_extension"`
		SourceExtensions 	[]string	`toml:"source_extensions"`
		PrefixGroup 		string		`toml:"prefix_group"`
		Cmd 			string		`toml:"cmd"`
	} `toml:"Compiler"`
}

func cfg_load(_type CfgType, _file string) (*Doc, bool){
	file, err := os.Open(_file)
	if err != nil{
		ERR("Unable to open the config file.", err)
		return nil, false 
	}
	defer file.Close()


	// Document buffer
	doc := &Doc{}
	res := false

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
func cfg_gather_deps(_doc *Doc, _root string, _deps map[string]bool, _deps_stack []*Entry) ([]*Entry, bool){
	if _doc == nil{
		ERR("Unresolved Document.")
		return nil, false
	}
	entry, exists := _doc.Build.Entries[_root]
	if exists == false{
		ERR("Unresolved dependency.")
		return nil, false
	}

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
	_deps_stack = append(_deps_stack, &entry)

	return _deps_stack, true
}
	
func cfg_build(_doc *Doc, _entry string) bool{
	stack, res := cfg_gather_deps(
		_doc, _entry,
		make(map[string]bool, 20),
		[]*Entry{})
	if !res{ 
		return false
	}

	// Build all dependencies
	for _, entry := range stack{
		res = cfg_entry_build(_doc, entry)
		if !res{
			return false
		}
	}

	return true
}

