package auth

import (
	"crypto/tls"
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gomail.v2"
)

type (
	EmailInfo struct{
		From		string
		To			string
		Subject		string
		Body		string
	}
)

func HashPassword(password string) (string,error){
	bytes, err := bcrypt.GenerateFromPassword([]byte(password),14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err:= bcrypt.CompareHashAndPassword([]byte(hash),[]byte(password))
	return err == nil
}

func SendEmail(emailInfo EmailInfo){

	m := gomail.NewMessage()

	m.SetHeader("From", emailInfo.From)
  
	m.SetHeader("To", emailInfo.To)
  
	m.SetHeader("Subject", emailInfo.Subject)
  
	m.SetBody("text/html", emailInfo.Body)
	
	d := gomail.NewDialer("smtp.gmail.com", 587, "onlibraryid@gmail.com", "Dev12345")
  

	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
  
	if err := d.DialAndSend(m); err != nil {
	  fmt.Println(err)
	  panic(err)
	}
}