package configs

import(
	"fmt"
	"path/filepath"
	
	"github.com/pelletier/go-toml/v2"

	. "qb/io"
)

type QB_PipeIdx = uint8
type QB_PipeEntry struct{
	Command 		string	 	`toml:"command"`
	CommandPolicyAlias 	string	 	`toml:"command_policy_alias"`
	CommandPolicyName 	string	 	`toml:"command_policy_name"`
	Flags 			[]string	`toml:"flags"`
	Definitions 		[]string	`toml:"definitions"`
	Hooks 			[]string	`toml:"hooks"`
	AlwaysRebuild		bool 		`toml:"always_rebuild"`
}

type QB_ConfigEntry struct{
	Name 			string	 `toml:"name"`
	SourceDirectory 	string	 `toml:"source_directory"`
	OutputDirectory 	string	 `toml:"output_directory"`
	Pipeline 		[]QB_PipeEntry `toml:"Pipeline"`
}
type QB_Config struct{
	Entries 		[]QB_ConfigEntry `toml:"Entry"`
}

func QBConfigLoad(_path string) (cfg QB_Config, res bool){
	file := QBInitFile(_path)

	// Config has to be .toml
	if filepath.Ext(file.FullPath) != ".toml"{
		fmt.Printf("Invalid config file extendsion: %s\n", file.FullPath)
		return
	}

	return QBConfigDecode(file)
}
func QBConfigDecode(_file QB_File) (cfg QB_Config, res bool){
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

