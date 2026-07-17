package configs

import(
	"fmt"
	"path/filepath"
	
	"github.com/pelletier/go-toml/v2"

	"qb/qbio"
)

type PipeIdx = uint8
type PipeEntry struct{
	Command 		string	 	`toml:"command"`
	CommandPolicyAlias 	string	 	`toml:"command_policy_alias"`
	CommandPolicyName 	string	 	`toml:"command_policy_name"`
	Flags 			[]string	`toml:"flags"`
	Definitions 		[]string	`toml:"definitions"`
	Hooks 			[]string	`toml:"hooks"`
	AlwaysRebuild		bool 		`toml:"always_rebuild"`
}

type ConfigEntry struct{
	Name 			string	 `toml:"name"`
	HeaderDirectory 	string	 `toml:"header_directory"`
	SourceDirectory 	string	 `toml:"source_directory"`
	OutputDirectory 	string	 `toml:"output_directory"`
	Pipeline 		[]PipeEntry `toml:"Pipeline"`
}
type Config struct{
	Entries 		[]ConfigEntry `toml:"Entry"`
}

func ConfigLoad(_path string) (cfg Config, res bool){
	file := qbio.InitFile(_path)

	// Config has to be .toml
	if filepath.Ext(file.FullPath) != ".toml"{
		fmt.Printf("Invalid config file extendsion: %s\n", file.FullPath)
		return
	}

	return ConfigDecode(file)
}
func ConfigDecode(_file qbio.File) (cfg Config, res bool){
	if !_file.IsValid(){
		return
	}

	err := toml.NewDecoder(_file.GetFile()).Decode(&cfg)
	if err != nil{
		fmt.Printf("Failed to decode: %s\n", _file.FullPath)
		return
	}

	return cfg, true
}

