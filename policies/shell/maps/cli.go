package shell

import(
	"fmt"

	. "qb/build"
	. "qb/policies/shell/cfg"
	. "qb/policies/shell/pwsh"
)

type Shell_ExecFunc func(*Shell_PolicyConfig, *QB_BuildState)(bool)
var SHELL_FUNC_MAP = map[string]Shell_ExecFunc{
	"powershell": PwshExecFromState,
}

func ShellCliFuncLookup(cli string)(exec Shell_ExecFunc, res bool){
	exec, res = SHELL_FUNC_MAP[cli]
	if !res{
		fmt.Printf("CLI function not found: '%s'", cli)
		return
	}

	return exec, true
}

