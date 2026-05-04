package main

import (
	"github.com/pelletier/go-toml/v2"
	"os"
	"fmt"
	"errors"
)

func (_type *FileType) UnmarshalText(_text []byte) error{
	switch string(_text){
	case "executable":
		*_type = FILE_EXE
	case "static_library":
		*_type = FILE_LIB
	default:
		ERR(fmt.Sprintf("Unknown file type: %s", string(_text)))
		return errors.New("")
	}

	return nil
}
func (_group *PrefixDesc) UnmarshalText(_text []byte) error{
	switch string(_text){
	/*
		Default clang prefix group
	*/
	case "clang": 
		*_group = PrefixDesc{
			SRC: "-c",
			OUT: "-o",
			INC: "-I",
			DEF: "-D",
			FLG: "-"}
	case "ar": 
		*_group = PrefixDesc{
			SRC: "",
			OUT: "",
			INC: "",
			DEF: "",
			FLG: "-"}
	default:
		ERR(fmt.Sprintf("Prefix group: %s", string(_text)))
		return errors.New("")
	}

	return nil
}

func cfg_um_decode(file *os.File) (*Doc, bool){
	if file == nil{
		ERR("Invalid config file.")
		return nil, false
	}

	var doc = &Doc{}

	err := toml.NewDecoder(file).Decode(doc)
	if err != nil{
		ERR("Config decoding failed.", err)
		return nil, false
	}

	return doc, true
}

