package types

import(
	"strings"
	"path/filepath"
	"regexp"
	"os"

	"qb/qberr"
)

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
		SourceOutExtension	string			`toml:"source_out_extension"`
		SourceExtensions 	[]string		`toml:"source_extensions"`
		HeaderExtensions 	[]string		`toml:"header_extensions"`
		CompilerPrefixGroup 	CompilerGroup		`toml:"compiler_prefix_group"`
		LinkerPrefixGroup 	LinkerGroup		`toml:"linker_prefix_group"`
		LibraryPrefixGroup 	LibraryGroup		`toml:"library_prefix_group"`
		CompilerCmd 		string			`toml:"compiler_cmd"`
		LinkerCmd 		string			`toml:"linker_cmd"`
		LibraryCmd 		string			`toml:"library_cmd"`
	} `toml:"Compiler"`
}
func (_doc *Doc)FromKey(_key string) *Entry{
	entry, exists := _doc.Entries[_key]
	if !exists{
		panic("Invalid entry key.")
	}

	return &entry
}
func (_doc *Doc)PathFromKey(_key string, _relative_path string) (string, bool){
	if len(_relative_path) == 0{
		qberr.ERR("Path is too short.")
		return "", false
	}

	entry := _doc.FromKey(_key)
	var builder strings.Builder

	switch _relative_path[0]{
	/*
		If hook starts with '.' character it should be relative to _entry.BuildDirectory field.
	*/
	case '.':
		builder.WriteString(filepath.Join(
			entry.BaseDirectory,
			_relative_path))
	/*
		If hook starts with '$' character it means the string is a reference to a dependency build_directory field.
	*/
	case '$':
		regex := regexp.MustCompile(`^\${(.+)}`)
		exp := regex.FindStringSubmatch(_relative_path)
		if len(exp) == 2{ // ONLY A SINGLE CASE FOR NOW 
			dep_entry := _doc.FromKey(exp[1])

			// Resolve the field of the dependency
			dep_str, res := _doc.PathFromKey(exp[1], dep_entry.BaseDirectory)
			if !res{
				return "", false
			}

			// Join dependency path and the subpath of _rel_path
			builder.WriteString(filepath.Join(
				dep_str,
				_relative_path[len(exp[0]):]))
		}else{
			qberr.ERR("Invalid dependency reference path.")
			return "", false
		}
	default:
		builder.WriteString(_relative_path)
	}
	fullpath := builder.String()
	if _, err := os.Stat(fullpath); err != nil{
		qberr.ERR("Invalid path.", err)
		return "", false
	}

	return builder.String(), true
}

