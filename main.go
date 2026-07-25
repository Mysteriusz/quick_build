package main

import(
	"os"
	"fmt"

	"qb/policies/llvm/clang/cfg"
	"qb/qbio"
	"qb/misc"
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
	for idx := 1; idx < len(args); idx++{
		arg := args[idx]
		if (arg == "--config-file" || arg == "-f") && desc.ConfigFile == ""{
			desc.ConfigFile = check_arg_value(args, idx)
			idx++
		}else{ goto err_config_file }
	}

err_config_file:
	if !qbio.InitFile(desc.ConfigFile).IsValid(){
		println("Argument parser error.")
		println("Invalid config file provided.")
		println("'"+ desc.ConfigFile + "'")
		os.Exit(1)
	}

	return desc
}

func main(){
	r, d := clang.ParseD(qbio.InitFile("D:\\ax_project\\ax_virt_layer\\win64\\user\\build\\i64_cpu.d"))
	_ = r
	misc.PrintArray(d.Deps.AllPaths())
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
			println(err.Message())
			return
		}
	}
}


