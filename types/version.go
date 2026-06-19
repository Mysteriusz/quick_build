package types

import(
	"encoding/json"
	"encoding/hex"
	"path/filepath"
	"crypto/sha256"
	"time"
	"os"
	"io"
	"errors"

	"qb/qberr"
)
type BuildVersion struct{
	Version 	string 		`json:"version"`
	Iteration 	uint32 		`json:"iteration"`
	Timestamp 	time.Time	`json:"timestamp"`
	Files 		map[string]string `json:"files"`
}
const BUILD_VER_BASENAME = ".qb_build_version.json"

func build_version_load(_builder *EntryBuilder) bool{
	if _builder == nil{
		return false
	}

	/*
		Version control disabled
	*/
	if !_builder.entry.VersionControl{
		return true
	}

	fullpath := filepath.Join(_builder.out_dir, BUILD_VER_BASENAME)

	file, err := os.Open(fullpath)
	if err != nil{
		// Edge case where the file doesnt yet exist but should be created
		if errors.Is(err, os.ErrNotExist){
			qberr.ERR("Build version file not found, building from all files.")
			goto ret
		}

		qberr.ERR("Failed to open the version control file.", err)
		return false
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&_builder.ver)
	if err != nil{
		qberr.ERR("Failed to decode the build version file.", err)
		return false
	}

ret:
	if _builder.ver.Files == nil{
		_builder.ver.Files = make(map[string]string)
	}
	return true
}
func build_version_save(_builder *EntryBuilder) bool{
	if _builder == nil{
		return false
	}

	/*
		Version control disabled
	*/
	if !_builder.entry.VersionControl{
		return true
	}
	
	fullpath := filepath.Join(_builder.out_dir, BUILD_VER_BASENAME)

	file,err := os.Create(fullpath)
	if err != nil{
		qberr.ERR("Unable to save the build version file.", err)
		return false
	}
	defer file.Close()

	_builder.ver.Timestamp = time.Now()
	_builder.ver.Iteration += 1

	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")

	err = enc.Encode(_builder.ver)
	if err != nil{
		qberr.ERR("Failed to write the build version file.", err)
		return false
	}

	return true
}
func build_version_check_and_update(_builder *EntryBuilder, _fullpath string) (_updated bool, _result bool){
	if _builder == nil{
		return false, false
	}

	/*
		Version control disabled
	*/
	if !_builder.entry.VersionControl{
		return true, true
	}

	file, err := os.Open(_fullpath)
	if err != nil{
		qberr.ERR("Unable to update file entry for the version control.", err)
		return false, false
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil{
		qberr.ERR("Failed to compute SHA256 of a file for the version control.", err)
		return false, false
	}

	hash_hex := hex.EncodeToString(hash.Sum(nil))
	/*
		Hash is the same so there was no update
	*/
	if _builder.ver.Files[_fullpath] == hash_hex{
		return false, true
	}
	_builder.ver.Files[_fullpath] = hash_hex

	return true, true
}

