package main

import(
	"os"
	"strings"
	"regexp"
	"path/filepath"
)

func cfg_path_resolve(_doc *Doc, _key string, _rel_path string) (string, bool){
	if _doc == nil{
		ERR("Invalid Document.")
		return "", false
	}
	if len(_rel_path) == 0{
		ERR("Path is too short.")
		return "", false
	}

	var entry *Entry = entry_from_key(_doc, _key)
	var builder strings.Builder

	switch _rel_path[0]{
	/*
		If hook starts with '.' character it should be relative to _entry.BuildDirectory field.
	*/
	case '.':
		builder.WriteString(filepath.Join(
			entry.BaseDirectory,
			_rel_path))
	/*
		If hook starts with '$' character it means the string is a reference to a dependency build_directory field.
	*/
	case '$':
		regex := regexp.MustCompile(`^\${(.+)}`)
		exp := regex.FindStringSubmatch(_rel_path)
		if len(exp) == 2{ // ONLY A SINGLE CASE FOR NOW 
			var dep_entry *Entry = entry_from_key(_doc, exp[1])

			// Resolve the field of the dependency
			dep_str, res := cfg_path_resolve(_doc, exp[1], dep_entry.BaseDirectory)
			if !res{
				return "", false
			}

			// Join dependency path and the subpath of _rel_path
			builder.WriteString(filepath.Join(
				dep_str,
				_rel_path[len(exp[0]):]))
		}else{
			ERR("Invalid dependency reference path.")
			return "", false
		}
	default:
		builder.WriteString(_rel_path)
	}
	fullpath := builder.String()
	if _, err := os.Stat(fullpath); err != nil{
		ERR("Failed to resolve path.", err)
		return "", false
	}

	return builder.String(), true
}

