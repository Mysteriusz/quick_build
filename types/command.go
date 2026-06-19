package types

import(
	"qb/qberr"
)

type CommandMode uint32
const(
	Source = iota
	Header
)

type Command struct{
	Cmd string
	Args []string
	Output string
	Input string
	Mode CommandMode
}

func CreateCommand(_builder *EntryBuilder, _exec string, _capacity int)(*Command, bool){
	if _builder == nil{
		qberr.ERR("Invalid argument.")
		return nil, false
	}

	command := Command{}
	command.Cmd = _exec
	command.Args = make([]string, 0, _capacity)

	return &command, true
}
func (_command *Command) WritePaths(_builder *EntryBuilder, _prefix string, _paths []string) bool{
	if _builder == nil || _paths == nil{
		qberr.ERR("Invalid argument.")
		return false
	}

	var res bool
	doc := _builder.GetDoc()

	for _, arg := range _paths{
		arg, res = doc.PathFromKey(_builder.EntryKey, arg)
		if !res{
			return false
		}
		_command.Args = append(_command.Args, _prefix + arg)
	}

	return true
}
func (_command *Command) WriteArgs(_builder *EntryBuilder, _prefix string, _args []string) bool{
	if _builder == nil || _args == nil{
		qberr.ERR("Invalid argument.")
		return false
	}

	for _, arg := range _args{
		_command.Args = append(_command.Args, _prefix + arg)
	}

	return true
}

