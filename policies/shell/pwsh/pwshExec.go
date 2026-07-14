package policies

import(
	. "qb/misc"
	. "qb/build"

	. "qb/policies/shell/cfg"
)

func PwshExecFromState(_policy *Shell_PolicyConfig, _state *QB_BuildState) (res bool){
	if _policy == nil || _state == nil{
		return false
	}

	PrintArray(QBRefConvertFld(_state, "$STATE.WORKING_SET"))

	return true
}

