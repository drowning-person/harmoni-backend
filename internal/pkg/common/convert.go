package common

import (
	"reflect"
	"unsafe"
)

func StringToBytes(s string) []byte {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := reflect.SliceHeader{
		Data: sh.Data,
		Len:  sh.Len,
		Cap:  sh.Len,
	}
	return *(*[]byte)(unsafe.Pointer(&bh))
}

func BytesToString(b []byte) string {
	sh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	bh := reflect.StringHeader{
		Data: sh.Data,
		Len:  sh.Len,
	}
	return *(*string)(unsafe.Pointer(&bh))
}
