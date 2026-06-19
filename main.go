package main

import (
	"qb/types"
)

func main(){
	doc, res := LoadConfig(types.CFG_UM, "D:/ax_project/quick_build/test.toml")
	if !res{
		return
	}

	res = BuildWithDeps(doc, "ax_utility_lib")
	if !res{
		return
	}
}

