package main

import (
	"crypto/hmac"
	"crypto/sha512"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	const (
		_pass string = "1234"
	)

	_salt := []byte{}
	for i := 1; i <= 64; i++ {
		_salt = append(_salt, byte(i))
	}

	encryptedPass, err := encrypt(_pass)
	if err != nil {
		log.Fatalf("failed to encrypt, %v", err.Error())
		return
	}
	fmt.Printf("Encryted pass: %v\n", encryptedPass)
	err = decryptAndVerify(encryptedPass, _pass)
	if err != nil {
		fmt.Println("Verification failed")
	}
	fmt.Println("Verification success")

	fmt.Println("Sign password")
	var signed []byte
	signed, err = sign(encryptedPass, _salt)
	if err != nil {
		log.Fatalf("failed to sign, %v", err.Error())
		return
	}
	fmt.Printf("Signed pass: %q\n", string(signed))
	var isSame bool
	isSame, err = verifySign(encryptedPass, signed, _salt)
	if err != nil {
		log.Fatalf("failed to sign, %v", err.Error())
		return
	}
	if isSame {
		fmt.Println("Signature verification failed")
	} else {
		fmt.Println("Signature verification success")
	}
}

func encrypt(password string) (string, error) {
	val, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(val), nil
}

func decryptAndVerify(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func sign(password string, key []byte) ([]byte, error) {
	mac := hmac.New(sha512.New, key)
	_, err := mac.Write([]byte(password))
	if err != nil {
		return nil, fmt.Errorf("failed to sign: %v", err.Error())
	}
	return mac.Sum(nil), nil
}

func verifySign(password string, signedPass, key []byte) (bool, error) {

	newSign, err := sign(password, key)
	if err != nil {
		return false, fmt.Errorf("failed to verify: %v", err.Error())
	}
	fmt.Printf("Original: %q\n", signedPass)
	fmt.Printf("New: %q\n", newSign)
	// mac.Sum(nil)
	return hmac.Equal(newSign, signedPass), nil
}
