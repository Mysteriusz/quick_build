package main

import(
	"strings"
	"path/filepath"
	"io/fs"
	"slices"
)

type PrefixDesc struct{
	SRC string
	OUT string
	INC string
	DEF string
	FLG string
}
var PREFIX_GROUPS = map[string]PrefixDesc{
	"clang": {SRC: "-c", OUT: "-o", INC: "-I", DEF: "-D", FLG: "-"},
}

func cfg_write_definitions(_doc *Doc, _entry *Entry, _out *strings.Builder) bool{
	var pref PrefixDesc = PREFIX_GROUPS[_doc.Compiler.PrefixGroup]

	// Write definitions
	for _, def := range _entry.Definitions{
		_out.WriteString(pref.DEF)
		_out.WriteString(def)
		_out.WriteString(" ")
	}

	return true
}
func cfg_write_includes(_doc *Doc, _entry *Entry, _out *strings.Builder) bool{
	var pref PrefixDesc = PREFIX_GROUPS[_doc.Compiler.PrefixGroup]

	// Write include (hook) directories
	for _, hook := range _entry.LinkHooks{
		fullpath, res := cfg_path_resolve(_doc, _entry, hook)
		if !res{
			return false
		}

		_out.WriteString(pref.INC)
		_out.WriteString(fullpath)
		_out.WriteString(" ")
	}

	return true
}
func cfg_write_flags(_doc *Doc, _entry *Entry, _out *strings.Builder) bool{
	var pref PrefixDesc = PREFIX_GROUPS[_doc.Compiler.PrefixGroup]

	for _, flag := range _entry.Flags{
		_out.WriteString(pref.FLG)
		_out.WriteString(flag)
		_out.WriteString(" ")
	}

	return true
}

func cfg_entry_build(_doc *Doc, _entry *Entry) bool{
	var pref PrefixDesc = PREFIX_GROUPS[_doc.Compiler.PrefixGroup]
	var f_builder strings.Builder
	var e_builder strings.Builder

	// Write the compiler executable
	f_builder.WriteString(_doc.Compiler.Cmd)
	f_builder.WriteString(" ")

	var res bool

	res = cfg_write_definitions(_doc, _entry, &f_builder)
	if !res{
		return false
	}

	res = cfg_write_includes(_doc, _entry, &f_builder)
	if !res{
		return false
	}

	res = cfg_write_flags(_doc, _entry, &f_builder)
	if !res{
		return false
	}

	// Write the source file prefix to FRONT_BUILDER
	f_builder.WriteString(pref.SRC)
	f_builder.WriteString(" ")
	FRONT_CMD := f_builder.String()

	// Write the output file prefix to END_BUILDER
	e_builder.WriteString(" ")
	e_builder.WriteString(pref.OUT)
	END_CMD := e_builder.String()

	/*
		Start compiling all source files
	*/
	filepath.Walk(_entry.BaseDirectory, func(path string, info fs.FileInfo, err error) error{
		if err != nil{
			ERR("Unknown error.", err)
			return err
		}

		// Skip file that`s not a source
		if !slices.Contains(_doc.Compiler.SourceExtensions, filepath.Ext(path)){
			return nil
		}

		out_path := filepath.Join(
			_entry.BuildDirectory,
			info.Name() + _doc.Compiler.SourceOutExtension)
		/*
			Form a command from builders and paths
		*/
		println(FRONT_CMD + path + END_CMD + out_path)

		return nil
	})

	return true
}
func cfg_entry_compile() bool{
	return true
}

