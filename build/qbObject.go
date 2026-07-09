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
	/*
		Allows for storage of additional data
		that isn`t as rooted as Data

		IMPORTANT!!!
		Extra should never be validated/used in any way by the generic qb handlers
		and should be used with caution
	*/
	extra	map[string]json.RawMessage
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

	obj.extra = make(map[string]json.RawMessage)
	obj.Type = _type
	return obj, true
}

func (_obj QB_Object) MarshalJSON() ([]byte, error){
	return json.Marshal(&struct {
        	Type QB_ObjectType 	`json:"type"`
        	Data any 		`json:"data"`
        	Extra map[string]json.RawMessage `json:"extra,omitempty"`
	}{
		Type: _obj.Type,
		Data: _obj.Data,
		Extra: _obj.extra,
	})
}
func (_obj *QB_Object) UnmarshalJSON(data []byte) error{
    	var raw struct {
        	Type QB_ObjectType  `json:"type"`
        	Data json.RawMessage `json:"data"`
        	Extra map[string]json.RawMessage `json:"extra,omitempty"`
    	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	_obj.Type = raw.Type
	_obj.extra = raw.Extra

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

func QBGetObjectExtra[T any](_obj *QB_Object, _slot string) (res bool, data T){
	raw, res := _obj.extra[_slot]

	if err := json.Unmarshal(raw, &data); err != nil{
		return 
	}

	return res, data
}
func QBSetObjectExtra[T any](_obj *QB_Object, _slot string, _data T) (res bool){
	val, err := json.Marshal(_data)
	if err != nil{
		return
	}

	_obj.extra[_slot] = val
	return true
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

type QB_ObjectSet map[string]QB_Object
func (_set *QB_ObjectSet) Update(_obj QB_Object) QB_ObjectSet{
	if *_set == nil{
		*_set = make(QB_ObjectSet)
	}
	(*_set)[_obj.String()] = _obj
	return *_set
}
func (_set *QB_ObjectSet) Merge(_objs QB_ObjectSet){
	if *_set == nil{
		*_set = make(QB_ObjectSet)
	}

	for _, v := range _objs{
		_set.Update(v)
	}
}

