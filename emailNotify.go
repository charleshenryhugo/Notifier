package main

import (
	"bytes"
	"crypto/tls"
	"log"
	"net/smtp"
	"strings"
)

type Mail struct {
	senderID string
	toIds    []string
	subject  string
	body     string
}

type SmtpServer struct {
	host      string
	port      string
	tlsconfig *tls.Config
}

func (s *SmtpServer) ServerName() string {
	return s.host + ":" + s.port
}

//Email has specific form so we need to build it
func (mail *Mail) BuildMessage() string {
	message := bytes.NewBufferString("From: ")
	message.WriteString(mail.senderID)
	message.WriteString("\r\n")
	message.WriteString("To: ")
	if len(mail.toIds) > 0 {
		message.WriteString(strings.Join(mail.toIds, ";"))
		message.WriteString("\r\n")
	}
	message.WriteString("Subject: ")
	message.WriteString(mail.subject)
	message.WriteString("\r\n\r\n")
	message.WriteString(mail.body)

	return message.String()
}

//initialize a mail struct
//trim the subject if it is longer than MAX_EMAIL_SUBJECT_LEN
func newMail(from string, to []string, subject string, body string) *Mail {
	if len(subject) > MAX_EMAIL_SUBJECT_LEN {
		subject = subject[0:255]
	}
	mail := new(Mail)
	mail.senderID = from
	mail.toIds = to
	mail.subject = subject
	mail.body = body

	return mail
}
func newSMTPServer(host string, port string) *SmtpServer {
	smtpServer := new(SmtpServer)
	smtpServer.host = host
	smtpServer.port = port
	smtpServer.tlsconfig = &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         smtpServer.host,
	}

	return smtpServer
}

//send email using SMTP with a specific SMTP server and an account
func smtpEmail(mail *Mail, smtpServer *SmtpServer, pwd string) ERR {
	msgBody := mail.BuildMessage()
	log.Println("connecting smtpserver", smtpServer.ServerName())

	//build an authentication
	auth := smtp.PlainAuth("", mail.senderID, pwd, smtpServer.host)

	conn, err := tls.Dial("tcp", smtpServer.ServerName(), smtpServer.tlsconfig)
	if err != nil { //no such host
		log.Println(err)
		return SMTPM_SVR_CONN_ERR
	}

	client, err := smtp.NewClient(conn, smtpServer.host)
	if err != nil {
		log.Println(err)
		return SMTPM_CLT_BLD_ERR
	}

	//Use Auth
	if err = client.Auth(auth); err != nil { //authentication failed
		log.Println(err)
		return SMTPM_AUTH_ERR
	}
	//add sender and receivers
	if err = client.Mail(mail.senderID); err != nil {
		log.Println(err)
		return SMTPM_SENDER_ERR
	}
	for _, k := range mail.toIds {
		//no need to verify target addresses
		//Many servers will not verify addresses for security reasons.
		log.Println("receiver address: ", k, " added successfully")
		if err = client.Rcpt(k); err != nil {
			log.Println(err)
			return SMTPM_RCVR_ERR
		}
	}

	//Data
	w, err := client.Data()
	if err != nil {
		log.Println(err)
		return SMTPM_CLT_IO_ERR
	}

	_, err = w.Write([]byte(msgBody))
	if err != nil {
		log.Println(err)
		return SMTPM_CLT_DATA_ERR
	}

	err = w.Close()
	if err != nil {
		log.Println(err)
		return SMTPM_CLT_IO_ERR
	}
	err = client.Quit()
	if err != nil {
		log.Println(err)
		return SMTPM_CLT_CLOSE_ERR
	}

	return SUCCESS
}

func emailNotifyHelp(from string, to []string, subject string, msg string, SMTPHost string, SMTPPort string, pwd string) ERR {
	mail := newMail(from, to, subject, msg)
	smtpServer := newSMTPServer(SMTPHost, SMTPPort)
	return smtpEmail(mail, smtpServer, pwd)
}

//EmailNotify (to []string, subject, msg string, ntfs Notifiers)
//send an email with subject and message provided with parameters
//to the email address stored in(to []string)
func EmailNotify(to []string, subject, msg string, ntfs Notifiers) ERR {
	if len(to) == 0 {
		return SMTPM_NOTGT
	}
	ntf := ntfs.SMTPEmailNotifier

	//check the notification type "smtpemail" and find if the state is "on"
	//if no type of "smtpemail" or the state is "off", do nothing and return directly
	if ntf.Type == "smtpemail" && (ntf.State == true) {
		return emailNotifyHelp(ntf.Account, to, subject, msg,
			ntf.SMTPHost, ntf.SMTPPort, ntf.Pwd)
	}

	return SMTPM_INVAL
}
