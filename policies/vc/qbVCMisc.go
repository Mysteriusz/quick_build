package policies

import(
	"time"
	"io"

	"strconv"
	"encoding/hex"
	"crypto/sha256"

	. "qb/io"
	. "qb/build"
)

/*
	Compute pipe log hash based on the state

	Computed fields of 'QB_BuildState' are:
		- CurrentPipe().Definitions
		- CurrentPipe().Hooks
		- CurrentPipe().Flags
		- GetHeaders().AllHashes()
		- GetSources().AllHashes()
*/
func VCStateUniqueHash(_state *QB_BuildState) string{
	if _state == nil{
		return ""
	}

	hash := sha256.New()
	for _, str := range _state.CurrentPipe().Definitions {
		io.WriteString(hash, str)
	}

	for _, str := range _state.CurrentPipe().Hooks {
		io.WriteString(hash, str)
	}

	for _, str := range _state.CurrentPipe().Flags {
		io.WriteString(hash, str)
	}

	for _, str := range _state.GetHeaders().AllHashes() {
		io.WriteString(hash, str)
	}
	
	for _, str := range _state.GetSources().AllHashes() {
		io.WriteString(hash, str)
	}

	return hex.EncodeToString(hash.Sum(nil))
}

/*
	Compute pipe log identifier based on the QB_BuildState.CurrentPipeIdx()

	Computed fields of 'QB_BuildState' are:
		- CurrentPipe().Command
		- CurrentPipe().CommandPolicyAlias
		- CurrentPipe().CommandPolicyName
		- CurrentPipeIdx()
*/
func VCPipeUniqueId(_state *QB_BuildState) string{
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
	Gather all changed objects to 'VC_FileState' object
*/
func VCDiff(_qb_state *QB_BuildState, _vc_state *VC_FileState) (not_diff bool){
	if _qb_state == nil || _vc_state == nil{
		return
	}

	d1 := VCDiffFiles(_qb_state.GetSources(), _vc_state.Pipe().SourceFiles)
	_vc_state.DiffSources = d1

	d2 := VCDiffFiles(_qb_state.GetHeaders(), _vc_state.Pipe().HeaderFiles)
	_vc_state.DiffHeaders = d2

	return len(d1) == 0 && len(d2) == 0
}

/*
	Check and return what files differ
*/
func VCDiffFiles(_f1 QB_FileArray, _f2 QB_FileArray) (diff QB_FileArray){
	if _f1 == nil || _f2 == nil{
		return
	}

	count := make(map[string]uint32)
	files := make(map[string]QB_File)
	s1 := make(map[string]bool)
	s2 := make(map[string]bool)
	for _, e := range _f1{
		if !s1[e.FullPath]{
			if count[e.FullPath] == 0{
				files[e.FullPath] = e
			}
			count[e.FullPath]++
			s1[e.FullPath] = true
		}
	}
	for _, e := range _f2{
		if !s2[e.FullPath]{
			if count[e.FullPath] == 0{
				files[e.FullPath] = e
			}

			/*
				Compare hashes to determine file content change
			*/
			f := files[e.FullPath]
			if f.ComputeHash() == e.ComputeHash(){
				count[e.FullPath]++
			}
			s2[e.FullPath] = true
		}
	}

	for k, v := range count{
		if v == 1{
			diff = append(diff, files[k])
		}
	}

	return diff
}
func VCIntersectFiles(_f1 QB_FileArray, _f2 QB_FileArray) (intersect QB_FileArray){
	if _f1 == nil || _f2 == nil{
		return
	}

	count := make(map[string]uint32)
	files := make(map[string]QB_File)
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
				continue
			}

			if f.ComputeHash() == e.ComputeHash(){
				count[e.FullPath]++
			}
			s2[e.FullPath] = true
		}
	}

	for k, v := range count{
		if v == 2{
			intersect = append(intersect, files[k])
		}
	}

	return intersect
}

func VCTimeToFormat(_time time.Time) string{
	return _time.Format(VC_TIME_FORMAT)
}

