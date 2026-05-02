package main

func main(){
	doc, res := cfg_load(CFG_UM, "D:/ax_project/quick_build/test.toml")
	if !res{
		return
	}
	_ = doc

	res = cfg_build(doc, "ax_virt_layer")
	if !res{
		return
	}
}

