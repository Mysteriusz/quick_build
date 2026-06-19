package types

import(
	"fmt"
	"errors"

	"qb/qberr"
)

type CfgType uint8
const(
	CFG_KM CfgType = iota
	CFG_UM
)

type FileType uint8
const(
	FILE_EXE FileType = iota 	// Windows Executable
	FILE_LIB  			// Static library
)
func (_type *FileType) UnmarshalText(_text []byte) error{
	switch string(_text){
	case "executable":
		*_type = FILE_EXE
	case "static_library":
		*_type = FILE_LIB
	default:
		qberr.ERR(fmt.Sprintf("Unknown file type: %s", string(_text)))
		return errors.New("")
	}

	return nil
}


