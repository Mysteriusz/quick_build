package build

import(
	"os"
	"fmt"
	"slices"

	"os/exec"
)

type QB_Command struct{
	/*
		Command line ordering
		[Args] [Input] [Output]

		Order of output and input may change depending
		on the 'input_idx' and 'output_idx' fields
	*/
	Args		[]string
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

func (_qb_cmd QB_Command) Run() (res bool){
	var io_args []string = _qb_cmd.Args[1:]

	/*
		Insert/Concat input
	*/
	if _qb_cmd.input_idx > 0{
		io_args = slices.Insert(io_args, _qb_cmd.input_idx, _qb_cmd.input...)
	}else{
		io_args = slices.Concat(io_args, _qb_cmd.input)
	}

	/*
		Insert/Concat output
	*/
	if _qb_cmd.output_idx > 0{
		io_args = slices.Insert(io_args, _qb_cmd.output_idx, _qb_cmd.output...)
	}else{
		io_args = slices.Concat(io_args, _qb_cmd.output)
	}


	// Execute the command with 'io_args'
	cmd := exec.Command(_qb_cmd.Args[0], io_args...)

	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	for _,arg := range cmd.Args{
		fmt.Println(arg)
	}

	err := cmd.Run()
	if err != nil{
		fmt.Println("Error occurend when execution specified command.")
		fmt.Println(err)
		return
	}

	return true
}

