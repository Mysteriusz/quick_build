package policies

import(
	"fmt"

	"github.com/pelletier/go-toml/v2"
)

type QB_PolicyFile struct{
	Policies 	map[string]any  	`toml:"Policies"`
}

func QBLoadPolicyFile(_policy QB_PolicyInfo) (file QB_PolicyFile, res bool){
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
func QBDecodePolicy[T any](_file QB_PolicyFile, _policy string) (dec T, res bool){
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

