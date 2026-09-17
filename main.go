package main

import(
	"os"
	"fmt"

	"qb/qbio"
	"qb/build"
	"qb/configs"
	"qb/build/runner"
)

type args_desc struct{
	ConfigFile string
}

func check_arg_value(args []string, idx int) (string){
	// Index is the last/oob so there is no value
	if idx >= len(args) - 1{
		return ""
	}
	return args[idx + 1]
}

func parse_args(args []string) args_desc {
	var desc args_desc 
	var message string

	for idx := 1; idx < len(args); idx += 2{
		arg := args[idx]
		if arg == "--config-file" || arg == "-f"{
			if desc.ConfigFile != ""{
				message = arg
				goto err_arg_repeat
			}

			desc.ConfigFile = check_arg_value(args, idx)
		}else{
			message = "Invalid config file provided."
			goto err_arg_invalid 
		}
	}

	if desc.ConfigFile == ""{
		message = "--config-file"
		goto err_arg_missing
	}
	if !qbio.InitFile(desc.ConfigFile).IsValid(){
		message = desc.ConfigFile
		goto err_arg_invalid
	}

	return desc

err_arg_invalid:
	fmt.Println("Argument parser error.")
	fmt.Println("Invalid argument provided:")
	fmt.Println("'" + message + "'")
	os.Exit(1)

err_arg_repeat:
	fmt.Println("Argument parser error.")
	fmt.Println("Repeated argument provided: ")
	fmt.Println("'" + message + "'")
	os.Exit(1)

err_arg_missing:
	fmt.Println("Argument parser error.")
	fmt.Println("Missing required argument: ")
	fmt.Println("'" + message + "'")
	os.Exit(1)

	// Should never return here 
	return args_desc{};
}

func main(){
	fmt.Println("Quick build.")

	args := parse_args(os.Args)

	cfg, res := configs.ConfigLoad(args.ConfigFile)
	if !res{
		return
	}

	for _, entry := range cfg.Entries{
		fmt.Println("=================================================")
		fmt.Printf("Starting entry build\n")
		fmt.Printf("Build name: %s\n", entry.Name)
		fmt.Printf("Input directory: %s\n", entry.SourceDirectory)
		fmt.Printf("Output directory: %s\n", entry.OutputDirectory)
		fmt.Printf("Pipeline length: %d\n", len(entry.Pipeline))
		fmt.Println("=================================================")

		state, err := qb.InitBuild(entry)
		if err.Check(){
			return
		}

		err = runner.ExecuteFromState(&state)
		if err.Check(){
			fmt.Println(err.Message())
			return
		}
	}
}


