package util

import "golang.org/x/crypto/bcrypt"

func If(cond bool, yes interface{}, no interface{}) interface{} {
	var result interface{}
	if cond {
		result = yes
	} else {
		result = no
	}
	switch v := result.(type) {
	case func():
		v()
		return nil
	case func() error:
		return v()
	case func() interface{}:
		return v()
	default:
		return v
	}
}

func HashPassword(password []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
}

func ComparePassword(hashedPassword []byte, password []byte) bool {
	return bcrypt.CompareHashAndPassword(hashedPassword, password) == nil
}
