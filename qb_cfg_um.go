package main

import (
	"github.com/pelletier/go-toml/v2"
	"os"
)

func cfg_um_decode(file *os.File) (*Doc, bool){
	if file == nil{
		ERR("Invalid config file.")
		return nil, false
	}

	doc := Doc{}

	err := toml.NewDecoder(file).Decode(&doc)
	if err != nil{
		ERR("Config decoding failed.", err)
		return nil, false
	}

	return &doc, true
}

