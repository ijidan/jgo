package jutils

import "golang.org/x/crypto/bcrypt"

//密码工具
type PasswordUtil struct {
}

//密码加密
func (u *PasswordUtil) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 16)
	return string(bytes), err
}

//校验密码
func (u *PasswordUtil) CheckPasswordHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
