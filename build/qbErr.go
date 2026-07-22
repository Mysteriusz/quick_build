package qb

import(
	"fmt"
	"errors"
)

type BuildError struct{
	Objects	[]string
	Pipe	string
	Err 	error
}

func(BuildError) None() BuildError{
	return BuildError{}
}
func(_err BuildError) New(_state *BuildState, _msg string, _objs ...string) BuildError{
	if _state == nil{
		_err.Err = errors.New("Nil state pointer.")
	}else{
		_err.Err = errors.New(_msg)
		_err.Pipe = _state.CurrentPipe().CommandPolicyName
	}

	_err.Objects = _objs
	return _err
}

func(_err BuildError) NilArgument(_state *BuildState) BuildError{
	return BuildError{}.New(_state, "Nil arugment provided.")
}

func(err BuildError) Message() string{
	if len(err.Objects) == 0{
		return fmt.Sprintf("Error in pipe: %s\n" + err.Err.Error(), err.Pipe)
	}else{
		return fmt.Sprintf("Error in pipe: %s\n" + err.Err.Error(), err.Pipe, err.Objects)
	}
}
func(err BuildError) Check() bool{
	return err.Err != nil
}

