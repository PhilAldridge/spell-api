package utils

import "reflect"

func ExtractIDsReflect[T any](items []*T) []int {
	ids := make([]int, 0, len(items))

	for _, item := range items {
		v := reflect.ValueOf(*item)
		field := v.FieldByName("ID")
		if field.IsValid() && field.Kind() == reflect.Int {
			ids = append(ids, int(field.Int()))
		}
	}

	return ids
}