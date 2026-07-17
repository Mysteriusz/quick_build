package maps

import(
	"fmt"

	"qb/build"
	"qb/policies/shell/cfg"
	"qb/policies/shell/pwsh"
)

type ExecFunc func(*shell.PolicyConfig, *qb.BuildState)(bool)
var SHELL_FUNC_MAP = map[string]ExecFunc{
	"powershell": pwsh.ExecFromState,
}

func CliFuncLookup(cli string)(exec ExecFunc, res bool){
	exec, res = SHELL_FUNC_MAP[cli]
	if !res{
		fmt.Printf("CLI function not found: '%s'", cli)
		return
	}

	return exec, true
}

