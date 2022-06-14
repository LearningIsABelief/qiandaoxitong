package util

import "unsafe"

// StringToByteSlice
// @Description: string转byte
// @Author: YangXuZheng
// @Date: 2022-06-13 11:28
func StringToByteSlice(s string) []byte {
	tmp1 := (*[2]uintptr)(unsafe.Pointer(&s))
	tmp2 := [3]uintptr{tmp1[0], tmp1[1], tmp1[1]}
	return *(*[]byte)(unsafe.Pointer(&tmp2))
}

// ByteSliceToString
// @Description: byte转string
// @Author: YangXuZheng
// @Date: 2022-06-13 11:36
func ByteSliceToString(bytes []byte) string {
	return *(*string)(unsafe.Pointer(&bytes))
}
