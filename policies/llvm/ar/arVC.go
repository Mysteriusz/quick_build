package ar

import(
	"qb/build"
	"qb/policies/vc"
)

func Diff(_qb_state *qb.BuildState, _vc_state *vc.FileState)(obj_diff vc.ObjectDiff){
	if _qb_state == nil || _vc_state == nil{
		return
	}

	for _, obj := range _vc_state.Pipe().OutWorkingSet{
		if !obj.Exists(){
			obj_diff.Removed.Update(obj)
		}
		if !obj.InvalidateHash(){
			obj_diff.Modified.Update(obj)
		}
	}

	return obj_diff
}

