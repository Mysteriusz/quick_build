package types

import(
	"fmt"
	"errors"

	"qb/qberr"
)

/*
	Processing symbol group
*/
type AnyGroup struct{
	NAME string
	SRC string
	OUT string
	INC string
	DEF string
	FLG string
	defaultCmd func(*AnyGroup, *EntryBuilder) (*Command, bool)
	fromCmd func(*AnyGroup, *Command) bool
}

/*
	Create a default command for the builder
*/
func (_group *AnyGroup) DefaultCmd(_builder *EntryBuilder) (*Command, bool){
	if _group.defaultCmd == nil{
		qberr.ERR(fmt.Sprintf("Invalid default command for a group: %s.", _group.NAME))
		return nil, false
	}
	return _group.defaultCmd(_group, _builder)
}

/*
	Execute command
*/
func (_group *AnyGroup) FromCmd(_command *Command) bool{
	if _group.defaultCmd == nil{
		qberr.ERR(fmt.Sprintf("Invalid default command for a group: %s.", _group.NAME))
		return false
	}
	return _group.fromCmd(_group, _command)
}

type CompilerGroup struct{
	AnyGroup
	GatherHeaders func(*Doc, string)([]string, bool)
}
type LinkerGroup struct{
	AnyGroup
}
type LibraryGroup struct{
	AnyGroup
}

func NoneDefaultCmd(*AnyGroup, *EntryBuilder) (*Command, bool){
	return nil, true
}
func NoneFromCmd(*AnyGroup, *Command) bool{
	return true
}

var ANY_GROUP_MAP = map[string]AnyGroup{
	"clang":{
		NAME: "clang",
		SRC: "-c",
		OUT: "-o",
		INC: "-I",
		DEF: "-D",
		FLG: "-",
		defaultCmd: ClangDefaultCmd,
		fromCmd: ClangFromCmd,
	},
	"ar":{
		NAME: "ar",
		SRC: "",
		OUT: "",
		INC: "",
		DEF: "",
		FLG: "-",
		defaultCmd: NoneDefaultCmd,
		fromCmd: NoneFromCmd,
	},
}

func (_group *CompilerGroup) UnmarshalText(_text []byte) error{
	switch string(_text){
	/*
		Default clang prefix group
	*/
	case "clang": 
		*_group = CompilerGroup{
			AnyGroup: ANY_GROUP_MAP[string(_text)],
			GatherHeaders: ClangGatherHeaders,
		}
	default:
		qberr.ERR(fmt.Sprintf("Unknown prefix group: %s", string(_text)))
		return errors.New("")
	}

	return nil
}
func (_group *LinkerGroup) UnmarshalText(_text []byte) error{
	switch string(_text){
	case "ar": 
		*_group = LinkerGroup{
			AnyGroup: ANY_GROUP_MAP[string(_text)],
		}
	default:
		qberr.ERR(fmt.Sprintf("Unknown prefix group: %s", string(_text)))
		return errors.New("")
	}

	return nil
}
func (_group *LibraryGroup) UnmarshalText(_text []byte) error{
	switch string(_text){
	case "ar": 
		*_group = LibraryGroup{
			AnyGroup: ANY_GROUP_MAP[string(_text)],
		}
	default:
		qberr.ERR(fmt.Sprintf("Unknown prefix group: %s", string(_text)))
		return errors.New("")
	}

	return nil
}

