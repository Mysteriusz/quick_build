package qb

import(
	"fmt"
	"encoding/json"

	"qb/qbio"
)

type FileObject struct{
	File 	qbio.File		`json:"file"`
}

type ObjectType uint8
const(
	/*
		Object.Data == File
	*/
	TYPE_FILE ObjectType = iota
)

type Object struct{
	Type 	ObjectType 	`json:"type"`
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

func InitObject(_data any, _type ObjectType) (obj Object, res bool){
	switch _type{
	case TYPE_FILE:
		obj.Data = FileObject{
			File: qbio.InitFile(_data.(string)),
		}
	default:
		return
	}

	obj.extra = make(map[string]json.RawMessage)
	obj.Type = _type
	return obj, true
}

func (_obj Object) MarshalJSON() ([]byte, error){
	return json.Marshal(&struct {
        	Type ObjectType 	`json:"type"`
        	Data any 		`json:"data"`
        	Extra map[string]json.RawMessage `json:"extra,omitempty"`
	}{
		Type: _obj.Type,
		Data: _obj.Data,
		Extra: _obj.extra,
	})
}
func (_obj *Object) UnmarshalJSON(data []byte) error{
    	var raw struct {
        	Type ObjectType  `json:"type"`
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
		var file FileObject
		if err := json.Unmarshal(raw.Data, &file); err != nil {
			return err
		}
		_obj.Data = file
	default:
		return fmt.Errorf("Error when unmarshaling object.")
	}
	return nil
}

func GetObjectExtra[T any](_obj *Object, _slot string) (res bool, data T){
	raw, res := _obj.extra[_slot]

	if err := json.Unmarshal(raw, &data); err != nil{
		return 
	}

	return res, data
}
func SetObjectExtra[T any](_obj *Object, _slot string, _data T) (res bool){
	val, err := json.Marshal(_data)
	if err != nil{
		return
	}

	_obj.extra[_slot] = val
	return true
}

func (_obj Object) Exists() bool{
	switch _obj.Type{
	case TYPE_FILE:
		file := _obj.Data.(FileObject).File
		return file.IsValid() && file.InvalidateHash()
	default:
		return false
	}
}
func (_obj Object) String() string{
	switch _obj.Type{
	case TYPE_FILE:
		return _obj.Data.(FileObject).File.FullPath
	default:
		return ""
	}
}

/*
			Hashing
*/
func(_obj *Object) ComputeHash() string{
	switch _obj.Type{
	case TYPE_FILE:
		f := _obj.Data.(FileObject).File
		return f.ComputeHash()
	default:
		return ""
	}
}
func (_obj *Object) InvalidateHash() bool{
	switch _obj.Type{
	case TYPE_FILE:
		f := _obj.Data.(FileObject).File
		return f.InvalidateHash()
	default:
		return false
	}
}

func (obj Object)Key() string{
	return obj.String()
}
func (obj Object)CheckKey(key string) bool{
	return obj.String() == key
}

type ObjectSet map[string]Object

func (_set ObjectSet) Has(_obj Object) bool{
	if _set == nil{
		return false
	}

	_, r := _set[_obj.String()]
	return r
}

func (_set *ObjectSet) Update(_obj Object) ObjectSet{
	if *_set == nil{
		*_set = make(ObjectSet)
	}
	(*_set)[_obj.Key()] = _obj
	return *_set
}
func (_set *ObjectSet) Remove(_obj Object) ObjectSet{
	if *_set == nil{
		*_set = make(ObjectSet)
		return *_set
	}
	delete(*_set, _obj.Key())
	return *_set
}

func (_set *ObjectSet) Intersect(_objs ObjectSet){
	if *_set == nil{
		*_set = make(ObjectSet)
	}

	for k, v := range _objs{
		if _set.Has(v){
			_set.Update(v)
		}else{
			delete(*_set, k)
		}
	}
}
func (_set *ObjectSet) Merge(_objs ObjectSet){
	if *_set == nil{
		*_set = make(ObjectSet)
	}

	for _, v := range _objs{
		_set.Update(v)
	}
}

func (_set *ObjectSet) StringArray() (arr []string){
	for k := range *_set{
		arr = append(arr, k)
	}
	return arr
}

