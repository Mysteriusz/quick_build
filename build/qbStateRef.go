package qb

/*
	TODO:

	Rethink the design and flow of reference values
*/

import(
	"slices"
	"strings"
)

type RefVarKind uint8
const(
	REF_STRING RefVarKind = iota
	REF_PATHS 
	REF_OBJECT
)
type RefVar struct{
	Kind 	RefVarKind
	Value 	[]string
}
var qb_EMPTY_REF = RefVar{Kind: 0, Value: nil}

func RefMapVar(_state *BuildState, fld string) RefVar{

	switch(fld){
	case "STATE.WORKING_SET":
		return RefVar{
			Kind: REF_STRING,
			Value: _state.WorkingSet.StringArray(),
		}
	case "STATE.GET_SOURCES":
		return RefVar{
			Kind: REF_PATHS,
			Value: _state.GatherAllSources().AllPaths(),
		}
	case "STATE.GET_HEADERS":
		return RefVar{
			Kind: REF_PATHS,
			Value: _state.GatherAllHeaders().AllPaths(),
		}
	case "STATE.SOURCE_DIRECTORY":
		return RefVar{
			Kind: REF_STRING,
			Value: []string{_state.Config.SourceDirectory},
		}
	case "STATE.OUTPUT_DIRECTORY":
		return RefVar{
			Kind: REF_STRING,
			Value: []string{_state.Config.OutputDirectory},
		}
	case "STATE.NAME":
		return RefVar{
			Kind: REF_STRING,
			Value: []string{_state.Config.Name},
		}
	default:
		return qb_EMPTY_REF
	}
}

/*
	Resolve a single field
*/
func RefConvertFld(_state *BuildState, fld string) RefVar{
	if fld[0] != '{'{
		return qb_EMPTY_REF
	}

	idx := strings.Index(fld, "}")
	if idx == -1{
		return qb_EMPTY_REF
	}

	return RefMapVar(_state, fld[1:idx])
}

/*
	Convert reference variable to string array

	Delimiter-Kind mapping: 
		- REF_STRING -> "" (no delimiter)
		- REF_PATHS -> " "
		- REF_OBJECT -> " "
*/
func RefDelim(Kind RefVarKind) string{
	switch(Kind){
	case REF_STRING:
		return ""
	case REF_PATHS:
		// no modifiers
		fallthrough
	case REF_OBJECT:
		// no modifiers
		fallthrough
	default:
		return " "
	}
}

/*
	Resolve all references across the entire string
*/
func RefResolve(_state *BuildState, val string) (arr []RefVar, res bool){
	idx := 0
	skip := false
	buf := ""

	for idx < len(val){
		// Skip the current character
		if skip{
			buf += string(val[idx])
			idx++

			skip = false
			continue
		}

		switch val[idx]{
		case '\\':
			skip = true
		case '{':
			// Make a slice from '{' index forward
			fld := val[idx:]

			// Convert field in place
			ref := RefConvertFld(_state, fld)
			if ref.Value == nil{
				return
			}

			// Append and reset the gathered buffer
			if buf != ""{
				ref.Value = slices.Insert(ref.Value, 0, buf)
				buf = ""
			}

			arr = append(arr, ref)

			// Skip to trailing character
			idx += strings.Index(fld, "}")
		default:
			buf += string(val[idx])
		}
		idx++
	}

	if buf != ""{
		if len(arr) == 0{
			arr = append(arr, RefVar{Kind: REF_STRING, Value: []string{buf}})
		}else{
			// Append the gathered buffer
			arr[len(arr) - 1].Value = append(arr[len(arr) - 1].Value, buf)
		}
	}

	return arr, true
}

/*
	Merge all strings based on their kind

	For example if 'RefVar.Kind' == REF_STRING
	the entire 'RefVar.Value' will be 'joined'
	and separated with kind-based delimiter

	The delimiter it defined by the function 'RefDelim()'

	Example:
		RefVar.Value = []string{"D:/path", "/to/something"}

		the joined version would be:

		RefVar.Value = []string{"D:/path/to/something"}
*/
func RefMergeByKind(refs []RefVar) (buf []string, res bool){
	if refs == nil{
		return
	}

	for _, ref := range refs{
		if ref.Kind == REF_STRING{
			temp := strings.Join(ref.Value, RefDelim(ref.Kind))
			buf = append(buf, temp)
		}else{
			buf = slices.Concat(buf, ref.Value)
		}
	}
	return buf, true
}

