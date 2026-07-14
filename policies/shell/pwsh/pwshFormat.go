package pwsh

import(
	. "qb/build"
)

func PwshFormatAsString(ref *QB_RefVar){
	if ref == nil{
		return
	}

	if len(ref.Value) == 0{
		return
	}

	for idx := range ref.Value{
		ref.Value[idx] = "\"" + ref.Value[idx] + "\""
	}
}
func PwshFormatAsArray(ref *QB_RefVar){
	if ref == nil{
		return
	}

	if len(ref.Value) == 0{
		return
	}

	/*
		Format the 'QB_RefVar.Value' 'edges' into 
		powershell compliant array

		Example:
			[]string{"@(some string" "other string)")}
	*/
	ref.Value[0] = "@(" + ref.Value[0]
	ref.Value[len(ref.Value) - 1] = ref.Value[len(ref.Value) - 1] + ")"

	/*
		Add commas
		
		Example:
			@("some string")
			@("some string", "other string")
	*/
	for idx := 0; idx < len(ref.Value) - 1; idx++{
		ref.Value[idx] = ref.Value[idx] + ","
	}
}

func PwshFormatRef(ref *QB_RefVar){
	if ref == nil{
		return
	}

	switch(ref.Kind){
	case REF_PATHS:
		PwshFormatAsString(ref)
		PwshFormatAsArray(ref)
	case REF_STRING:
		fallthrough
	case REF_OBJECT:
		fallthrough
	default:
		return
	}
}


