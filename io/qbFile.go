package io

import(
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
*/
type QB_File struct{
	FullPath 	string
	file		*os.File
}

func QBInitFile(_path string) (qb_file QB_File){
	qb_file.FullPath = _path
	qb_file.file = nil

	return qb_file
}

func (_file QB_File) GetFile() *os.File{
	// If already opened just pass the file
	if _file.file != nil{
		return _file.file
	}

	// Else open the file
	os_file, err := os.OpenFile(_file.FullPath, os.O_RDWR | os.O_CREATE, 0644)
	if err != nil{
		fmt.Printf("Failed to open: %s\n", _file.FullPath)
		panic("Assertion failed!!!")
	}

	_file.file = os_file
	return os_file
}

/*
	Create/Close the file and save
*/
func (_file QB_File) Save() (res bool){
	// State where the file exists but isn`t open
	if _file.IsValid() && _file.file == nil{
		return true
	}

	// State where the file exists and is open
	if _file.file != nil{
		err := _file.file.Close()
		if err != nil{
			fmt.Printf("Unable to save file:\n %s\n", _file.FullPath)
			return
		}
	}else{
		// State where the doesn`t exist
		os_file,err := os.Create(_file.FullPath)
		if err != nil{
			fmt.Printf("Unable to save file:\n %s\n", _file.FullPath)
			return
		}
		os_file.Close()
	}

	return true
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

type QB_FileArray []QB_File

func (_farray QB_FileArray) AllPaths() []string{
	return Select[QB_File, string](
		_farray,
		"FullPath")
}

