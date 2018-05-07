package parse

import (
	"crypto/tls"
	"strconv"
	"fmt"
	"strings"
	"log"
	"net/smtp"
)

type EmailConfig struct {
	User   string
	Passwd string
	Server   string
	Port   int
}

type Mail struct {
	To 		[]string
	Cc 		[]string
	Bcc 	[]string
	Subject	string
	Body 	string
	Email 	*EmailConfig
}

func (m *Mail) dial() (*tls.Conn,error) {
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName: m.Email.Server,
	}
	server := m.Email.Server + ":" + strconv.Itoa(m.Email.Port)
	//fmt.Println(tlsconfig.ServerName)
	return tls.Dial("tcp",server,tlsconfig)
}

func (m *Mail) buildMessage() string {
	header := ""
	header += fmt.Sprintf("From: %s\r\n",m.Email.User)
	if len(m.To) > 0 {
		header += fmt.Sprintf("To: %s\r\n",strings.Join(m.To,";"))
	}
	if len(m.Cc) > 0 {
		header += fmt.Sprintf("Cc: %s\r\n",strings.Join(m.Cc,";"))
	}
	header += fmt.Sprintf("Subject: %s\r\n",m.Subject)
	header += "\r\n" + m.Body
	return header
}

func (m *Mail) SendEmail() (err error) {
	conn,err := m.dial()
	if err != nil {
		log.Printf("[Faild]: connection err: %v\n",err)
		//return err
		log.Panic(err)
	}
	client,err := smtp.NewClient(conn,m.Email.Server)
	//fmt.Println(client)
	if err != nil {
		log.Printf("[Faild]: smtp new client err: %v\n",err)
		//return err
		log.Panic(err)
	}

	auth := smtp.PlainAuth("", m.Email.User, m.Email.Passwd, m.Email.Server)
	if err := client.Auth(auth);err != nil {
		log.Printf("[Faild]: auth err: %v\n",err)
		log.Panic(err)
		//return err
	}

	if err := client.Mail(m.Email.User);err != nil {
		log.Printf("[Faild]: client mail err: %v\n",err)
		//return err
		log.Panic(err)
	}

	recevies := append(m.To,m.Cc...)
	recevies = append(recevies,m.Bcc...)
	for _,k := range recevies {
		log.Println("Sending to: ",k)
		if err := client.Rcpt(k);err != nil {
			log.Printf("[Faild]: client rcpt err: %v\n",err)
			log.Panic(err)
		}
	}

	writer,err := client.Data()
	if err != nil {
		log.Printf("[Faild]: client data err: %v\n",err)
		//return err
		log.Panic(err)
	}

	_,err = writer.Write([]byte(m.buildMessage()))
	if err != nil {
		log.Printf("[Faild]: write err: %v\n",err)
		//return err
		log.Panic(err)
	}

	return client.Quit()
}

