package misc

import(
	"fmt"
	"reflect"
)

func Select[A, F any](_array []A, _field string) (array []F){
	for i := range _array{
		rv := reflect.ValueOf(_array[i])
		value := rv.FieldByName(_field).Interface()

		array = append(array, value.(F))
	}
	return array
}



/*
	TODO: Optimize this atrocity xd
*/
func Union[A comparable](_a ...[]A) (array []A){
	seen := make(map[A]bool)
	for _, slice := range _a{
		for _, e := range slice{
			if _, res := seen[e]; res{
				continue
			}
                	seen[e] = true
			array = append(array, e)
		}
	}

	return array
}
func Intersect[A comparable](_a ...[]A) (array []A){
	count := make(map[A]int)
	for _, slice := range _a{
		for _, e := range slice{
			count[e]++
		}
	}

	for k, v := range count{
		if v == len(_a){
			array = append(array, k)
		}
	}

	return array
}
func Diff[A comparable](_a ...[]A) (array []A){
	exclude := make(map[A]int)
	for _, slice := range _a{
		seen := make(map[A]bool)
		for _, e := range slice{
			if !seen[e]{
				exclude[e]++
			}
		}
	}

	for k, v := range exclude{
		if v == 1{
			array = append(array, k)
		}
	}

	return 
}

func PrintArray(_a ...[]string){
	for _, slice := range _a{
		for _, e := range slice{
			fmt.Printf("%s\n", e)
		}
	}
}

