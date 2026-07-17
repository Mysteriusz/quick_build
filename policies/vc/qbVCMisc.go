package vc

import(
	"io"
	"fmt"
	"time"

	"strconv"
	"encoding/hex"
	"crypto/sha256"

	"qb/qbio"
	"qb/build"
)

/*
	Compute pipe log hash based on the state

	Computed fields of 'qb.BuildState' are:
		- CurrentPipe().PolicyName
		- CurrentPipe().Definitions
		- CurrentPipe().Hooks
		- CurrentPipe().Flags
*/
func StateUniqueHash(_state *qb.BuildState) string{
	if _state == nil{
		return ""
	}

	hash := sha256.New()

	io.WriteString(hash, _state.CurrentPipe().CommandPolicyName)

	for _, str := range _state.CurrentPipe().Definitions {
		io.WriteString(hash, str)
	}

	for _, str := range _state.CurrentPipe().Hooks {
		io.WriteString(hash, str)
	}

	for _, str := range _state.CurrentPipe().Flags {
		io.WriteString(hash, str)
	}

	return hex.EncodeToString(hash.Sum(nil))
}

/*
	Compute pipe log identifier based on the qb.BuildState.CurrentPipeIdx()

	Computed fields of 'qb.BuildState' are:
		- CurrentPipe().Command
		- CurrentPipe().CommandPolicyAlias
		- CurrentPipe().CommandPolicyName
		- CurrentPipeIdx()
*/
func PipeUniqueId(_state *qb.BuildState) string{
	if _state == nil{
		return ""
	}

	hash := sha256.New()

	io.WriteString(hash, _state.CurrentPipe().Command)
	io.WriteString(hash, _state.CurrentPipe().CommandPolicyAlias)
	io.WriteString(hash, _state.CurrentPipe().CommandPolicyName)
	io.WriteString(hash, strconv.Itoa(int(_state.CurrentPipeIdx())))

	return hex.EncodeToString(hash.Sum(nil))
}

/*
	Check and return what files differ
*/
func DiffFiles(_f1 qbio.FileArray, _f2 qbio.FileArray) (diff FileDiff){
	if _f1 == nil || _f2 == nil{
		return
	}

	count := make(map[string]uint32)
	files := make(map[string]qbio.File)
	s1 := make(map[string]bool)
	s2 := make(map[string]bool)
	for _, e := range _f1{
		if !s1[e.FullPath]{
			files[e.FullPath] = e
			count[e.FullPath]++
			s1[e.FullPath] = true
		}
	}
	for _, e := range _f2{
		if !s2[e.FullPath]{
			f, r := files[e.FullPath]
			if !r{
				// Didn`t exist in the '_f1' array
				files[e.FullPath] = e
				count[e.FullPath] = 1
				continue
			}

			/*
				Compare hashes to determine file content change
			*/
			if f.ComputeHash() == e.ComputeHash(){
				count[e.FullPath] = 0
			}
			s2[e.FullPath] = true
		}
	}

	for k, v := range count{
		if !files[k].IsValid(){
			diff.Removed = append(diff.Removed, files[k])
		}

		if v == 1{
			diff.Modified = append(diff.Modified, files[k])
		}
		// v == 0 means file has no diff
	}

	return diff
}
func DiffObjects(_f1 qb.ObjectSet, _f2 qb.ObjectSet) (diff ObjectDiff){
	if _f1 == nil || _f2 == nil{
		return
	}

	count := make(map[string]uint32)
	objs := make(map[string]qb.Object)
	s1 := make(map[string]bool)
	s2 := make(map[string]bool)
	for k, e := range _f1{
		if !e.CheckKey(k){
			fmt.Printf("Invalid key in object set.\n Valid key: %s\n Specified key: %s\n", e.Key(), k)
			delete(_f1, k)
		}

		if !s1[k]{
			objs[k] = e
			count[k]++
			s1[k] = true
		}
	}
	for k, e := range _f2{
		if !e.CheckKey(k){
			fmt.Printf("Invalid key in object set.\n Valid key: %s\n Specified key: %s\n", e.Key(), k)
			delete(_f1, k)
		}

		if !s2[k]{
			f, r := objs[k]
			if !r{
				// Didn`t exist in the '_f1' array
				objs[k] = e
				count[k] = 1
				continue
			}

			/*
				Compare hashes to determine file content change
			*/
			if f.ComputeHash() == e.ComputeHash(){
				count[k]++
			}
			s2[k] = true
		}
	}

	for k, v := range count{
		if !objs[k].Exists(){
			diff.Removed.Update(objs[k])
		}
		if v == 1{
			diff.Modified.Update(objs[k])
		}
	}
	return
}


func TimeToFormat(_time time.Time) string{
	return _time.Format(TIME_FORMAT)
}

