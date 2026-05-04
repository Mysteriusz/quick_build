package main

import(
	"encoding/json"
	"encoding/hex"
	"path/filepath"
	"crypto/sha256"
	"time"
	"os"
	"io"
	"errors"
)
type BuildVersion struct{
	Version 	string 		`json:"version"`
	Iteration 	uint32 		`json:"iteration"`
	Timestamp 	time.Time	`json:"timestamp"`
	Files 		map[string]string `json:"files"`
}
const BUILD_VER_BASENAME = ".qb_build_version.json"

func build_version_load(_doc *Doc, _key string) (*BuildVersion, bool){
	if _doc == nil{
		return nil, false
	}

	var res bool
	var ver = &BuildVersion{}
	var entry *Entry = entry_from_key(_doc, _key)
	
	/*
		Version control disabled
	*/
	if !entry.VersionControl{
		return nil, true
	}

	build_dir,res := cfg_path_resolve(_doc, _key, entry.BuildDirectory)
	if !res{
		return nil, false
	}
	fullpath := filepath.Join(build_dir, BUILD_VER_BASENAME)

	file, err := os.Open(fullpath)
	if err != nil{
		// Edge case where the file doesnt yet exist but should be created
		if errors.Is(err, os.ErrNotExist){
			ERR("Build version file not found, building from all files.")
			goto ret
		}

		ERR("Failed to open the version control file.", err)
		return nil, false
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(ver)
	if err != nil{
		ERR("Failed to decode the build version file.", err)
		return nil, false
	}

ret:
	if ver.Files == nil{
		ver.Files = make(map[string]string)
	}
	return ver, true
}
func build_version_save(_doc *Doc, _key string, _ver *BuildVersion) bool{
	if _doc == nil{
		return false
	}
	var entry *Entry = entry_from_key(_doc, _key)

	/*
		Version control disabled
	*/
	if !entry.VersionControl{
		return true
	}
	
	build_dir,res := cfg_path_resolve(_doc, _key, entry.BuildDirectory)
	if !res{
		return false
	}
	fullpath := filepath.Join(build_dir, BUILD_VER_BASENAME)

	file,err := os.Create(fullpath)
	if err != nil{
		ERR("Unable to save the build version file.", err)
		return false
	}
	defer file.Close()

	_ver.Timestamp = time.Now()
	_ver.Iteration += 1

	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")

	err = enc.Encode(_ver)
	if err != nil{
		ERR("Failed to write the build version file.", err)
		return false
	}

	return true
}
func build_version_check_and_update(_doc *Doc, _key string, _ver *BuildVersion, _fullpath string) (_updated bool, _result bool){
	if _doc == nil{
		return false, false
	}
	var entry *Entry = entry_from_key(_doc, _key)

	/*
		Version control disabled
	*/
	if !entry.VersionControl{
		return false, true
	}

	file, err := os.Open(_fullpath)
	if err != nil{
		ERR("Unable to update file entry for the version control.", err)
		return false, false
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil{
		ERR("Failed to compute SHA256 of a file for the version control.", err)
		return false, false
	}

	if _ver.Files[_fullpath] == hex.EncodeToString(hash.Sum(nil)){
		return false, true
	}
	_ver.Files[_fullpath] = hex.EncodeToString(hash.Sum(nil))

	return true, true
}

