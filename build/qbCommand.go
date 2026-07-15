package build

import(
	"os"
	"fmt"
	"slices"
	"strings"

	"os/exec"

	"qb/misc"
)

type QB_Command struct{
	/*
		Command line ordering
		[Exec] [Args] [Input] [Output]

		Order of output and input may change depending
		on the 'input_idx' and 'output_idx' fields
	*/
	Args		[]string
	Exec 		string
	Directory 	string
	input_idx	int
	output_idx	int
	input		[]string
	output		[]string
}

func QBInitCommand(_args []string, _input_idx int, _output_idx int) (cmd QB_Command){
	if _args == nil{
		return
	}

	cmd.Args = _args

	// _args[0] is the actual command hence the - 1 to normalize
	cmd.input_idx = _input_idx - 1
	cmd.output_idx = _output_idx - 1

	return cmd
}

func (_cmd *QB_Command) SetInput(_input []string) *QB_Command{
	_cmd.input = _input
	return _cmd
}
func (_cmd *QB_Command) SetOutput(_output []string) *QB_Command{
	_cmd.output = _output
	return _cmd
}

func (_cmd *QB_Command) AttachArgs(_args []string) *QB_Command{
	_cmd.Args = slices.Concat(_cmd.Args, _args)
	return _cmd
}
func (_cmd *QB_Command) AttachArgsFromIndex(_args []string, _idx int) *QB_Command{
	_cmd.Args = slices.Insert(_cmd.Args, _idx, _args...)
	return _cmd
}

func (_qb_cmd *QB_Command) GetCmd() (exec string, args []string){
	/*
		Check if executable string is set

		If not assume that 'QB_Command.Args[0]' is the executable
	*/
	offset := 0
	exec = _qb_cmd.Exec
	if _qb_cmd.Exec == ""{
		offset = 1
		exec = _qb_cmd.Args[0]
	}

	args = _qb_cmd.Args[offset:]

	/*
		Insert/Concat input
	*/
	if _qb_cmd.input_idx > 0{
		args = slices.Insert(args, _qb_cmd.input_idx, _qb_cmd.input...)
	}else{
		args = slices.Concat(args, _qb_cmd.input)
	}

	/*
		Insert/Concat output
	*/
	if _qb_cmd.output_idx > 0{
		args = slices.Insert(args, _qb_cmd.output_idx, _qb_cmd.output...)
	}else{
		args = slices.Concat(args, _qb_cmd.output)
	}

	return exec, args
}
func (_qb_cmd QB_Command) Run() (res bool){
	cmdline, args := _qb_cmd.GetCmd()

	cmd := exec.Command(cmdline, args...)
	cmd.Dir = _qb_cmd.Directory

	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	println(cmdline)
	misc.PrintArray(args)

	// Execute the command
	err := cmd.Run()
	if err != nil{
		fmt.Println("Error occurend when execution specified command.")
		fmt.Println(err)
		return
	}

	return true
}

/*
	TODO

	Change to something like RunAs(QB_CommandShellKind)
	and use some predefined shells

	IMPORTANT!!!
	Args have to be compliant with the powershell syntax
*/
func (_qb_cmd QB_Command) RunPowershell() (res bool){
	exe, args := _qb_cmd.GetCmd()
	cmdline := exe + " " + strings.Join(args, " ")

	// Create a powershell command object
	cmd := exec.Command("powershell", "-nologo", "-noprofile", "-command ", cmdline)

	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	misc.PrintArray(args)

	// Execute the command
	err := cmd.Run()
	if err != nil{
		fmt.Println("Error occured when execution specified command.")
		fmt.Println(err)
		return
	}

	return true
}

