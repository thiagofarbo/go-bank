package model

import (
	"go-bank/src/helper"
	"gopkg.in/gomail.v2"
	"log"
	"os"
	"strings"
)

func SendEmail(accountFrom Account, accountTo Account) {

	host := "sandbox.smtp.mailtrap.io"
	port := 587
	user := "62a153a9cec90b"
	password := "c9fe45faea3489"

	// Ler o conteúdo do arquivo HTML
	htmlContent, err := os.ReadFile("email_template.html")
	if err != nil {
		log.Fatalf("Erro ao ler o arquivo HTML: %v", err)
	}

	// Converter o conteúdo para string
	htmlStr := string(htmlContent)

	htmlStr = strings.Replace(htmlStr, "{{ContaRemetente}}", accountFrom.Number, -1)
	htmlStr = strings.Replace(htmlStr, "{{ContaDestinatario}}", accountTo.Number, -1)
	htmlStr = strings.Replace(htmlStr, "{{Valor}}", helper.ToString(accountFrom.Balance), -1)

	msg := gomail.NewMessage()
	msg.SetHeader("From", "thiagofarbo@gmail.com")
	msg.SetHeader("To", "webworkportal@gmail.com")
	msg.SetHeader("Subject", "test message")
	msg.SetBody("text/html", htmlStr)

	dialer := gomail.NewDialer(host, port, user, password)

	if err := dialer.DialAndSend(msg); err != nil {
		log.Fatalf("Erro ao enviar o e-mail: %v", err)
	} else {
		log.Println("E-mail enviado com sucesso!")
	}

}
