package clang

import(
	"os"
	"bufio"
	"bytes"
	"strings"

	"qb/qbio"
)

type DFile struct{
	/*
		All dependencies

		[0] == .o
		[1] == .c
		[2:] == .h
	*/
	Deps 	qbio.FileArray
}

/*
	Trim character set for clang generated .d files

	Contains:
		- Tab (0x09) 
		- Space (0x20)
		- Slash-Back (0x5c)
		- Colon (0x3a)
*/
const D_TRIM = "\x09\x20\x5c\x3a"
func ParseD(_file qbio.File) (res bool, dep DFile){
	reader := bufio.NewReader(_file.GetFile())
	defer _file.Save()

	var buf bytes.Buffer

	/*
		Read each dependency file path
	*/
	for{
		val, err := reader.ReadByte()
		if err != nil{
			break
		}

		if val != ' ' && val != '\n' && val != '\r'{
			buf.WriteByte(val)
			continue
		}

		path := strings.Trim(buf.String(), D_TRIM)

		// Ignore the path if it`s not correct
		stat, err := os.Stat(path)
		if err != nil || stat.IsDir(){
			continue
		}

		dep.Deps = append(dep.Deps, qbio.InitFile(path))
		buf.Reset()
	}
	if len(dep.Deps) < 2{
		return
	}

	return true, dep
}

