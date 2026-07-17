package policies

import(
	"fmt"

	"github.com/pelletier/go-toml/v2"
)

type PolicyFile struct{
	Policies 	map[string]any  	`toml:"Policies"`
}

func LoadPolicyFile(_policy PolicyInfo) (file PolicyFile, res bool){
	if _policy == nil{
		return 
	}

	qb_file := _policy.GetFile()
	defer qb_file.Save()

	err := toml.NewDecoder(qb_file.GetFile()).Decode(&file)
	if err != nil{
		fmt.Printf("Unable to decode the policy file:\n %s\n", qb_file.FullPath)
		return
	}

	return file, true
}
func DecodePolicy[T any](_file PolicyFile, _policy string) (dec T, res bool){
	payload, exists := _file.Policies[_policy]
	if !exists{
		fmt.Printf("Policy doesn`t exist: '%s'", _policy)
		return
	}

	// Marshal the policy data into bytes
	data, err := toml.Marshal(payload)
	if err != nil {
		return
	}

	// Unmarshal the policy into the desired interface
	err = toml.Unmarshal(data, &dec)
	if err != nil{
		fmt.Printf("Unable to decode named policy: '%s'", _policy)
		return
	}

	return dec, true
}

