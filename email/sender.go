package email

import "github.com/go-mail/mail"

func Send() {
	m := mail.NewMessage()
	m.SetHeader("From", "thiagoemidio37@yahoo.com.br")

	m.SetHeader("To", "thiagoemidio37@yahoo.com.br")

	//m.SetAddressHeader("Cc", "oliver.doe@example.com", "Oliver")

	m.SetHeader("Subject", "Hello!")

	m.SetBody("text/html", "Hello <b>Kate</b> and <i>Noah</i>!")

	m.Attach("lolcat.jpg")

	d := mail.NewDialer("smtp.mail.yahoo.com", 465, "thiagoemidio37@yahoo.com.br", "Tgo4471@!")

	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}
