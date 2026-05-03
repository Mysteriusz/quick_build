package main

import(
	"strings"
	"path/filepath"
)

func cfg_slib_cmd(_doc *Doc, _entry *Entry, _sources []string) ([]string, bool){
	if _doc == nil || _entry == nil || _sources == nil{
		ERR("Invalid arguments.")
		return nil, false
	}

	out_dir, res := cfg_path_resolve(_doc, _entry, _entry.BuildDirectory)
	if !res{
		return nil, false
	}

	out_path := filepath.Join(out_dir, _entry.OutputBasename)

	args_count := len(_entry.Flags) +
		len(_entry.LinkFlags) + 
		len(_sources) * 2 + 	// len(_sources) * (SRC_PREF + SRC_PATH)
		2			// OUT_PREF + OUT_PATH
	args := make([]string, 0, args_count)

	for _, source := range _sources{
		args = append(args, _doc.Compiler.LibraryPrefixGroup.SRC)
		args = append(args, source)
	}
	args = append(args, _doc.Compiler.LibraryPrefixGroup.OUT)
	args = append(args, out_path)

	println(strings.Join(args, " "))

	return args, true
	//return _doc.Compiler.PrefixGroup.SLIB_CMD, f_builder.String(), true
}

