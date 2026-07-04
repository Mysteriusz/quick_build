package build

import(
	. "qb/io"
)

type QB_FileObject struct{
	File 	QB_File
}

type QB_ObjectType uint8
const(
	TYPE_FILE QB_ObjectType = iota
)

type QB_Object struct{
	Type 	QB_ObjectType 	`json:"type"`
	Data	any		`json:"data"`
}

func QBInitObject(_data any, _type QB_ObjectType) (obj QB_Object, res bool){
	switch _type{
	case TYPE_FILE:
		obj.Data = QB_FileObject{
			File: QBInitFile(_data.(string)),
		}
	default:
		return
	}

	obj.Type = _type
	return obj, true
}

