package build

import(
	"fmt"
	"encoding/json"

	. "qb/io"
)

type QB_FileObject struct{
	File 	QB_File		`json:"file"`
}

type QB_ObjectType uint8
const(
	/*
		QB_Object.Data == QB_File
	*/
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

func (_obj *QB_Object) UnmarshalJSON(data []byte) error{
    	var raw struct {
        	Type QB_ObjectType  `json:"type"`
        	Data json.RawMessage `json:"data"`
    	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	_obj.Type = raw.Type

	switch _obj.Type{
	case TYPE_FILE:
		var file QB_FileObject
		if err := json.Unmarshal(raw.Data, &file); err != nil {
			return err
		}
		_obj.Data = file
	default:
		return fmt.Errorf("Error when unmarshaling object.")
	}
	return nil
}
func (_obj *QB_Object) Exists() bool{
	switch _obj.Type{
	case TYPE_FILE:
		return _obj.Data.(QB_FileObject).File.IsValid()
	default:
		return false
	}
}
func (_obj *QB_Object) String() string{
	switch _obj.Type{
	case TYPE_FILE:
		return _obj.Data.(QB_FileObject).File.FullPath
	default:
		return ""
	}
}

