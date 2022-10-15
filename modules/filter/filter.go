package filter

import (
	"fmt"
	"reflect"
)

func RemoveElementinInterface(t interface{}, keyname string, object string) error {
	intType := reflect.TypeOf(t).Elem()
	resultsVal := reflect.New(intType)
	rv := (reflect.ValueOf(t))
	if resultsVal.Kind() != reflect.Ptr || rv.Kind() != reflect.Ptr {
		return fmt.Errorf("results argument must be a pointer to a slice, but was a %s", resultsVal.Kind())
	}
	rvElem := rv.Elem()
	interfaceVal := resultsVal.Elem()
	Len := (rvElem).Len()
	for i := 0; i < Len; i++ {
		value := rvElem.Index(i)
		if ((rvElem.Index(i)).FieldByName(keyname)).String() == object {
			continue
		}
		interfaceVal = reflect.Append(interfaceVal, value)
	}
	rv.Elem().Set(interfaceVal)
	return nil
}
