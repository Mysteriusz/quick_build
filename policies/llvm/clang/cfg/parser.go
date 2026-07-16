package clang

import(
	"os"
	"bufio"
	"bytes"
	"strings"

	. "qb/io"
)

type ClangVC_D struct{
	/*
		All dependencies

		[0] == .o
		[1] == .c
		[2:] == .h
	*/
	Deps 	QB_FileArray
}

/*
	Trim character set for clang generated .d files

	Contains:
		- Tab (0x09) 
		- Space (0x20)
		- Slash-Back (0x5c)
		- Colon (0x3a)
*/
const CLANG_VC_D_TRIM = "\x09\x20\x5c\x3a"
func ClangVCParseD(_file QB_File) (res bool, dep ClangVC_D){
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

		path := strings.Trim(buf.String(), CLANG_VC_D_TRIM)

		// Ignore the path if it`s not correct
		stat, err := os.Stat(path)
		if err != nil || stat.IsDir(){
			continue
		}

		dep.Deps = append(dep.Deps, QBInitFile(path))
		buf.Reset()
	}
	if len(dep.Deps) < 2{
		return
	}

	return true, dep
}

