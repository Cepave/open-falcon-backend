package nqmDemo

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/mail"
	"net/smtp"
	"strings"

	"github.com/spf13/viper"
)

func encodeRFC2047(maddr string) string {
	// use mail's rfc2047 to encode any string
	addr := mail.Address{maddr, ""}
	return strings.Trim(addr.String(), " <>")
}

//http://11.11.11.11:9000/mail

func sendMail(Tos string, Subject string, Content string) {
	// Set up authentication information.

	smtpServer := viper.GetString("email_conf.smtp_server")
	stmpPort := viper.GetInt("email_conf.smtp_port")
	account := viper.GetString("email_conf.account")
	password := viper.GetString("email_conf.password")
	fromEmail := viper.GetString("email_conf.from_email")

	auth := smtp.PlainAuth(
		"",
		account,
		password,
		smtpServer,
	)

	from := mail.Address{"監控中心", fromEmail}
	to := mail.Address{"", Tos}

	header := make(map[string]string)
	header["From"] = from.String()
	header["To"] = to.String()
	header["Subject"] = encodeRFC2047(Subject)
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/plain; charset=\"utf-8\""
	header["Content-Transfer-Encoding"] = "base64"

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + base64.StdEncoding.EncodeToString([]byte(Content))

	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	err := smtp.SendMail(
		fmt.Sprintf("%s:%d", smtpServer, stmpPort),
		auth,
		from.Address,
		[]string{to.Address},
		[]byte(message),
		//[]byte("This is the email body."),
	)
	if err != nil {
		log.Fatal(err)
	}
}
