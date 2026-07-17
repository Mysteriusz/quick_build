package shell

type PolicyConfig struct{
	Cli 	string		`toml:"cli"`
	Args 	[]string	`toml:"args"`
}

