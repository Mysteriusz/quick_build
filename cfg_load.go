package main

import(
	"github.com/pelletier/go-toml/v2"

	"os"

	"qb/types"
	"qb/qberr"
)

func cfg_um_decode(file *os.File) (*types.Doc, bool){
	if file == nil{
		qberr.ERR("Invalid config file.")
		return nil, false
	}

	var doc = &types.Doc{}

	err := toml.NewDecoder(file).Decode(doc)
	if err != nil{
		qberr.ERR("Config decoding failed.", err)
		return nil, false
	}

	return doc, true
}

func LoadConfig(_type types.CfgType, _file string) (*types.Doc, bool){
	file, err := os.Open(_file)
	if err != nil{
		qberr.ERR("Unable to open the config file.", err)
		return nil, false 
	}
	defer file.Close()

	var res bool
	var doc *types.Doc

	switch _type{
	case types.CFG_UM:
		doc,res = cfg_um_decode(file)
	/*
	* 		Add CFG_KM
	*/
	default:
		qberr.ERR("Unsupported config type.")
		return nil, false
	}
	
	if !res{
		return nil, false
	}
	return doc, true
}

