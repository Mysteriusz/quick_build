package io

import(
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"

	"io"
	"os"
	"fmt"
	"path"
	"path/filepath"

	. "qb/misc"
)

/*
	TODO: (VERY IMPORTANT)
	Make the fullpath be validated whenever you acquire the file
	and if 'QB_File.FullPath' reload the 'QB_File.file' field!!!
	(Use a private prev_fullpath field or something) 

	JSON representation example:
        {
          "fullpath": "path//to//my//file",
          "hash": "3f90c0d3ad1c7fcfeba79b2ae4c600e6610a8874e6d4656355ecbbb2dcc03bed"
        }
*/
type QB_File struct{
	FullPath 	string 
	hash		string
	file		*os.File
}

func (_file QB_File) T() string{
	return _file.hash
}
/*
			Hashing
*/
func(_file *QB_File) recomputeHash() string{
	hash := sha256.New()
	if _, err := io.Copy(hash, _file.GetFileReadOnly()); err != nil{
		return ""
	}

	_file.hash = hex.EncodeToString(hash.Sum(nil))
	_file.Save()

	return _file.hash
}
func (_file *QB_File) ComputeHash() string{
	if !_file.IsValid(){
		return ""
	}
	if _file.hash != ""{
		return _file.hash
	}
	return _file.recomputeHash()
}
func (_file *QB_File) InvalidateHash() bool{
	temp := QBInitFile(_file.FullPath)
	return _file.ComputeHash() == temp.ComputeHash()
}

func (_file QB_File) Compare(file QB_File) bool{
	return _file.ComputeHash() == file.ComputeHash() && _file.FullPath == file.FullPath
}

func (_file QB_File) MarshalJSON() ([]byte, error){
	type alias QB_File
	return json.Marshal(struct{
		FullPath string `json:"fullpath"`
		Hash string `json:"hash"`
	}{
		FullPath: _file.FullPath,
		Hash: _file.ComputeHash(),
	})
}
func (f *QB_File) UnmarshalJSON(data []byte) error {
 	var raw struct {
		FullPath string `json:"fullpath"`
        	Hash     string `json:"hash"`
	}

	if err := json.Unmarshal(data, &raw); err != nil{
		return err
	}

	f.FullPath = raw.FullPath
	f.hash = raw.Hash

	return nil
}

func QBInitFile(_path string) (qb_file QB_File){
	qb_file.FullPath = NormalizePath(_path)
	qb_file.file = nil

	return qb_file
}

func (_file *QB_File) GetFile() *os.File{
	// If already opened just pass the file
	if _file.file != nil{
		return _file.file
	}

	// Else open the file
	os_file, err := os.OpenFile(_file.FullPath, os.O_RDWR | os.O_CREATE, 0644)
	if err != nil{
		fmt.Printf("Failed to open:\n %s\n", _file.FullPath)
		panic("Assertion failed!!!")
	}

	_file.file = os_file
	return os_file
}
func (_file *QB_File) GetFileReadOnly() *os.File{
	// If already opened just pass the file
	if _file.file != nil{
		return _file.file
	}

	// Else open the file
	os_file, err := os.OpenFile(_file.FullPath, os.O_RDONLY, 0644)
	if err != nil{
		fmt.Printf("Failed to open:\n %s\n", _file.FullPath)
	}

	_file.file = os_file
	return os_file
}

/*
	Create/Close the file
*/
func (_file *QB_File) Save() (res bool){
	// State where the file exists but isn`t open
	if _file.IsValid() && _file.file == nil{
		return true
	}

	// State where the file exists and is open
	if _file.file != nil{
		err := _file.file.Close()
		_file.file = nil

		if err != nil{
			fmt.Printf("Unable to save file:\n %s\n", _file.FullPath)
			return
		}
	}else{
		// State where the doesn`t exist
		os_file,err := os.Create(_file.FullPath)
		_file.file = nil

		if err != nil{
			fmt.Printf("Unable to save file:\n %s\n", _file.FullPath)
			return
		}
		os_file.Close()
	}

	return true
}
func (_file *QB_File) Clear(){
	_file.GetFile()
	_file.file.Truncate(0)
	_file.file.Seek(0, 0)
}
func (_file QB_File) IsValid() (res bool){
	stat, err := os.Lstat(_file.FullPath)
	if err != nil{
		return
	}

	// Config cannot be a directory
	if stat.IsDir(){
		return
	}

	return true
}

func ChangeExtension(_path string, _ext string) string{
	// Slice out the extension
	no_ext := _path[:len(_path) - len(path.Ext(_path))]

	if _ext[0] == '.'{
		return no_ext + _ext
	}else{
		return no_ext + "." + _ext
	}
}
func ChangeDirectory(_path string, _dir string) string{
	// Slice out the directory
	no_dir := _path[len(_path) - len(filepath.Base(_path)):]

	return filepath.Join(_dir, no_dir)
}
func ComparePath(_p1 string, _p2 string) bool{
	return filepath.Clean(_p1) == filepath.Clean(_p2)
}
func NormalizePath(_path string) string{
	abs, err := filepath.Abs(_path)
	if err != nil{
		return _path
	}

	abs = filepath.Clean(abs)
	abs = filepath.FromSlash(abs)

	return abs
}

type QB_FileArray []QB_File

func (_farray QB_FileArray) AllPaths() []string{
	return Select[QB_File, string](
		_farray,
		"FullPath")
}
func (_farray QB_FileArray) AllHashes() []string{
	buf := make([]string, len(_farray))
	for _,file := range _farray{
		buf = append(buf, file.ComputeHash())
	}
	return buf
}
func (_farray QB_FileArray) AllInvalid() (invalid []QB_File){
	for _,file := range _farray{
		if !file.IsValid(){
			invalid = append(invalid, file)
		}
	}
	return
}


func QBFileArrayUnion(_a ...QB_FileArray) (array QB_FileArray){
	seen := make(map[string]bool)
	for _, slice := range _a{
		for _, e := range slice{
			if _, res := seen[e.FullPath]; res{
				continue
			}
                	seen[e.FullPath] = true
			array = append(array, e)
		}
	}

	return array
}
func QBFileArrayIntersect(_a ...QB_FileArray) (array QB_FileArray){
	count := make(map[string]uint32)
	files := make(map[string]QB_File)
	for _, slice := range _a{
		for _, e := range slice{
			if count[e.FullPath] == 0{
				files[e.FullPath] = e
			}
			count[e.FullPath]++
		}
	}

	for k, v := range count{
		if v == uint32(len(_a)){
			array = append(array, files[k])
		}
	}

	return array
}

