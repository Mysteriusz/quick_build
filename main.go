package main

import(
	. "qb/configs"
	. "qb/build"
	. "qb/build/runner"
	//"qb/misc"
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
	ExecuteFromState(&state)

	//misc.PrintArray(misc.Union([]string{"a", "b", "c"}, []string{"c", "b", "d"}))
}


