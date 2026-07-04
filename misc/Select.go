package misc

import(
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

