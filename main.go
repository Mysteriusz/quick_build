package main

import(
	"fmt"

	"qb/configs"
	"qb/build"
	"qb/build/runner"
	//"qb/policies/vc"
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
	_ = cfg

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
	/*o1, _ := qb.InitObject("D:/fld/build/file0.o", qb.TYPE_FILE)
	o2, _ := qb.InitObject("D:/fld/build/file1.o", qb.TYPE_FILE)

	s1 := qb.ObjectSet{}
	s1.Update(o1)
	s1.Update(o2)

	s2 := qb.ObjectSet{}
	s2.Update(o1)

	diff := vc.DiffObjects(s1, s2)
	println("m")
	for _, f:= range diff.Modified{
		println(f.String())
	}
	println("r")
	for _, f:= range diff.Removed{
		println(f.String())
	}*/


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


