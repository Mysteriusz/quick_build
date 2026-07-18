package main

import(
	"fmt"

	"qb/configs"
	"qb/build"
	"qb/build/runner"
	/*"qb/misc"
	. "qb/io"
	. "qb/policies/vc"*/
)

func main(){
	//cfg, res := QBConfigLoad("D:/ax_project/quick_build_new/test.toml")
	cfg, res := configs.ConfigLoad("D:/ax_project/quick_build_new/cases.toml")
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
			return
		}
	}

	/*v := VCIntersectFiles(
		QB_FileArray{
			QBInitFile("D:\\ax_project\\ax_virt_layer_utils\\utils\\ax_utility_lib\\src\\io\\structures\\murmur.h"),
			QBInitFile("D:/ax_project/ax_virt_layer_utils/utils/ax_utility_lib/src/io/ax_memory_state.h"),
		},
		QB_FileArray{
			QBInitFile("D:\\ax_project\\ax_virt_layer_utils\\utils\\ax_utility_lib\\src\\io\\structures\\murmur.h"),
			QBInitFile("D:/ax_project/ax_virt_layer_utils/utils/ax_utility_lib/src/io/ax_memory_state.h"),
		},
	)
	misc.PrintArray(v.AllPaths())*/
	//misc.PrintArray(misc.Union([]string{"a", "b", "c"}, []string{"c", "b", "d"}))
}


