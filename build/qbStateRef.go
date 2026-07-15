package build

/*
	TODO:

	Rethink the design and flow of reference values
*/

import(
	"slices"
	"strings"
)

type QB_RefVarKind uint8
const(
	REF_STRING QB_RefVarKind = iota
	REF_PATHS 
	REF_OBJECT
)
type QB_RefVar struct{
	Kind 	QB_RefVarKind
	Value 	[]string
}
var qb_EMPTY_REF = QB_RefVar{Kind: 0, Value: nil}

func QB_RefMapVar(_state *QB_BuildState, fld string) QB_RefVar{
	switch(fld){
	case "STATE.WORKING_SET":
		return QB_RefVar{
			Kind: REF_STRING,
			Value: _state.WorkingSet.StringArray(),
		}
	case "STATE.GET_SOURCES":
		return QB_RefVar{
			Kind: REF_PATHS,
			Value: _state.GatherAllSources().AllPaths(),
		}
	case "STATE.GET_HEADERS":
		return QB_RefVar{
			Kind: REF_PATHS,
			Value: _state.GatherAllHeaders().AllPaths(),
		}
	case "STATE.SOURCE_DIRECTORY":
		return QB_RefVar{
			Kind: REF_STRING,
			Value: []string{_state.Config.SourceDirectory},
		}
	case "STATE.OUTPUT_DIRECTORY":
		return QB_RefVar{
			Kind: REF_STRING,
			Value: []string{_state.Config.OutputDirectory},
		}
	case "STATE.NAME":
		return QB_RefVar{
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
func QBRefConvertFld(_state *QB_BuildState, fld string) QB_RefVar{
	if fld[0] != '{'{
		return qb_EMPTY_REF
	}

	idx := strings.Index(fld, "}")
	if idx == -1{
		return qb_EMPTY_REF
	}

	return QB_RefMapVar(_state, fld[1:idx])
}

/*
	Convert reference variable to string array

	Delimiter-Kind mapping: 
		- REF_STRING -> "" (no delimiter)
		- REF_PATHS -> " "
		- REF_OBJECT -> " "
*/
func QB_RefDelim(Kind QB_RefVarKind) string{
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
func QBRefResolve(_state *QB_BuildState, val string) (arr []QB_RefVar, res bool){
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
			ref := QBRefConvertFld(_state, fld)
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
			arr = append(arr, QB_RefVar{Kind: REF_STRING, Value: []string{buf}})
		}else{
			// Append the gathered buffer
			arr[len(arr) - 1].Value = append(arr[len(arr) - 1].Value, buf)
		}
	}

	return arr, true
}

/*
	Merge all strings based on their kind

	For example if 'QB_RefVar.Kind' == REF_STRING
	the entire 'QB_RefVar.Value' will be 'joined'
	and separated with kind-based delimiter

	The delimiter it defined by the function 'QB_RefDelim()'

	Example:
		QB_RefVar.Value = []string{"D:/path", "/to/something"}

		the joined version would be:

		QB_RefVar.Value = []string{"D:/path/to/something"}
*/
func QBRefMergeByKind(refs []QB_RefVar) (buf []string, res bool){
	if refs == nil{
		return
	}

	for _, ref := range refs{
		if ref.Kind == REF_STRING{
			temp := strings.Join(ref.Value, QB_RefDelim(ref.Kind))
			buf = append(buf, temp)
		}else{
			buf = slices.Concat(buf, ref.Value)
		}
	}
	return buf, true
}

