package types

import(
	"os/exec"
	"os"
	"strings"
	"bytes"

	"qb/qberr"
)

/*
	Create a CLANG-compatible command object.
	The version of clang shouldn`t be a concern. (at least for now)
*/
func ClangDefaultCmd(_group *AnyGroup, _builder *EntryBuilder) (*Command, bool){
	if _group == nil || _builder == nil{
		qberr.ERR("Invalid argument.")
		return nil, false
	}

	var res bool
	doc := _builder.GetDoc()
	entry := doc.FromKey(_builder.EntryKey)

	/*
		Estimate amount of arguments used to minimize extending the slice
	*/
	args_count := len(entry.CompilerFlags) +
		len(entry.Definitions) +
		len(entry.LinkHooks) +
		2 + 	// SRC_PREF + SRC_PATH
		2	// OUT_PREF + OUT_PATH

	command, res := CreateCommand(_builder, doc.Compiler.CompilerCmd, args_count)
	if !res{
		return nil, false
	}

	/*
		Example:
			-x
			-std=c++17
	*/
	res = command.WriteArgs(_builder, _group.FLG, entry.CompilerFlags)
	if !res{
		return nil, false
	}

	/*
		Example:
			-DMY_DEF
	*/
	res = command.WriteArgs(_builder, _group.DEF, entry.Definitions)
	if !res{
		return nil, false
	}

	/*
		Example:
			-ID:/path/to/file
	*/
	res = command.WritePaths(_builder, _group.INC, entry.LinkHooks)
	if !res{
		return nil, false
	}

	return command, true
}
func ClangFromCmd(_group *AnyGroup, _command *Command) bool{
	if _group == nil || _command == nil{
		qberr.ERR("Invalid argument.")
		return false
	}

	// Ignore headers (temporary)
	if _command.Mode == Header{
		return true
	}

	// Form a command string
	args := make([]string, 0, len(_command.Args)+4)
	args = append(args, _command.Args...)
	args = append(args, _group.SRC, _command.Input, _group.OUT, _command.Output)

	println(strings.Join(args, " "))

	cmd := exec.Command(_command.Cmd, args...)
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil{
		qberr.ERR("Error during compilation.", err)
		return false
	}

	return true
}

func ClangGatherHeaders(_doc *Doc, _source string) ([]string, bool){
	if _doc == nil{
		qberr.ERR("Invalid document.")
		return nil, false
	}

	var stdout bytes.Buffer

	cmd := exec.Command(_doc.Compiler.CompilerCmd, "-MM", _source)
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = &stdout

	if err := cmd.Run(); err != nil{
		qberr.ERR("Error during header gathering.", err)
		return nil, false
	}

	qberr.ERR(stdout.ReadString(0))

	return nil, true
}

