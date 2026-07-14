package shell

type Shell_PolicyConfig struct{
	Cli 	string		`toml:"cli"`
	Args 	[]string	`toml:"args"`
}

