package main

import(
	. "qb/configs"
	. "qb/build"
	. "qb/build/runner"
	/*"qb/misc"
	. "qb/io"
	. "qb/policies/vc"*/
)

func main(){
	cfg, res := QBConfigLoad("D:/ax_project/quick_build_new/test.toml")
	if !res{
		return
	}
	state, res := QBInitBuild(&cfg.Entries[0])
	if !res{
		return
	}
	if !ExecuteFromState(&state){
		return
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


