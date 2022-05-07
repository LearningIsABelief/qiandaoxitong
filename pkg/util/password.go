package util

import "golang.org/x/crypto/bcrypt"

// Encrypt 密码加密
func Encrypt(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedBytes), err
}

// 密码解密
