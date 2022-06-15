package util

import (
	"fmt"
	"testing"
)

func TestCreateCode(t *testing.T) {
	code, s, err := CreateCode()
	if err != nil {
		return
	}
	var _ = s
	fmt.Println("id =", code)
	fmt.Println("正确答案 =", Result.Get(code, true))
	result1 := VerifyCaptcha(code, Result.Get(code, true))
	result2 := VerifyCaptcha(code, "1")
	fmt.Println("result1 =", result1)
	fmt.Println("result2 =", result2)
}
