package policies

import(
	. "qb/build"
	. "qb/policies/vc"
)

func ArVCDiff(_qb_state *QB_BuildState, _vc_state *VC_FileState)(obj_diff QB_ObjectSet){
	if _qb_state == nil || _vc_state == nil{
		return
	}

	for _, obj := range _vc_state.Pipe().OutWorkingSet{
		if !obj.Exists(){
			obj_diff.Update(obj)
		}
	}

	return obj_diff
}

