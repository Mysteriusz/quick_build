package policies

import(
	"fmt"

	"qb/qbio"

	"github.com/pelletier/go-toml/v2"
)

type PolicyFile struct{
	Policies 	map[string]any  	`toml:"Policies"`
	qbio.File
}

func LoadPolicyFile(_path string) (policy PolicyFile, res bool){
	policy.File = qbio.InitFile(_path)
	defer policy.File.Save()

	err := toml.NewDecoder(policy.GetFileReadOnly()).Decode(&policy)
	if err != nil{
		fmt.Printf("Unable to decode the policy file:\n %s\n", policy.FullPath)
		return
	}

	return policy, true
}
func DecodeConfig[CFG_T any](_file PolicyFile, _name string) (dec CFG_T, res bool){
	payload, exists := _file.Policies[_name]
	if !exists{
		fmt.Printf("Policy doesn`t exist: '%s'", _name)
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
		fmt.Printf("Unable to decode named policy: '%s'", _name)
		return
	}

	return dec, true
}

