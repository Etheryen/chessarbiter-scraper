package utils

import (
	"fmt"
	"reflect"
)

func PrintStruct[T any](s T) {
	v := reflect.ValueOf(s)
	typeOfS := v.Type()

	for i := 0; i < v.NumField(); i++ {
		fmt.Printf(
			"%v: %v\n",
			typeOfS.Field(i).Name,
			v.Field(i).Interface(),
		)
	}
}

func PrintStructSlice[T any](slice []T) {
	fmt.Println("{")
	for i, s := range slice {
		if i != 0 {
			fmt.Println()
		}
		PrintStruct(s)
	}
	fmt.Println("}")
}
