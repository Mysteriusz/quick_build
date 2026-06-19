package qberr

import(
	"fmt"
	"github.com/fatih/color"
)

func ERR(_msg string, _err ...error){
	if len(_err) > 0{
		fmt.Println(color.YellowString(_msg))
		fmt.Println(fmt.Errorf(color.YellowString("Error message: ") + "%s", color.RedString(_err[0].Error())))
	}else{
		fmt.Println(color.YellowString(_msg))
	}
}

