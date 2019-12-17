package sendmail

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/smtp"
	"os"
	"path/filepath"
)

var (
	host       = os.Getenv("EMAIL_HOST")
	username   = os.Getenv("EMAiL_USERNAME")
	password   = os.Getenv("EMAIL_PASSWORD")
	portNumber = os.Getenv("EMAIL_PORT")
)

// Sender : smtp struct
type Sender struct {
	auth smtp.Auth
}

// Message : email details
type Message struct {
	To          []string
	Subject     string
	Body        string
	Attachments map[string][]byte
}

// NewSender : new sender object
func NewSender() *Sender {
	auth := smtp.PlainAuth("", username, password, host)
	return &Sender{auth}
}

// Send : send message
func (s *Sender) Send(m *Message) error {
	return smtp.SendMail(fmt.Sprintf("%s:%s", host, portNumber), s.auth, username, m.To, m.ToBytes())
}

// NewMessage : new message
func NewMessage(s, b string) *Message {
	return &Message{Subject: s, Body: b, Attachments: make(map[string][]byte)}
}

// AttachFile : attach file to mail
func (m *Message) AttachFile(src string) error {
	b, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}

	_, fileName := filepath.Split(src)
	m.Attachments[fileName] = b
	return nil
}

// ToBytes : convert to byte
func (m *Message) ToBytes() []byte {
	buf := bytes.NewBuffer(nil)
	withAttachments := len(m.Attachments) > 0

	buf.WriteString(fmt.Sprintf("Subject: %s\n", m.Subject))
	buf.WriteString("MIME-Version: 1.0\n")
	writer := multipart.NewWriter(buf)
	boundary := writer.Boundary()

	if withAttachments {
		buf.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s\n", boundary))
		buf.WriteString(fmt.Sprintf("--%s\n", boundary))
	}

	buf.WriteString("Content-Type: text/plain; charset=utf-8\n")
	buf.WriteString(m.Body)

	if withAttachments {
		for k, v := range m.Attachments {
			buf.WriteString(fmt.Sprintf("\n\n--%s\n", boundary))
			buf.WriteString("Content-Type: application/octet-stream\n")
			buf.WriteString("Content-Transfer-Encoding: base64\n")
			buf.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=%s\n", k))

			b := make([]byte, base64.StdEncoding.EncodedLen(len(v)))
			base64.StdEncoding.Encode(b, v)
			buf.Write(b)
			buf.WriteString(fmt.Sprintf("\n--%s", boundary))
		}

		buf.WriteString("--")
	}

	return buf.Bytes()
}

// to send
// sender := NewSender()
// m := NewMessage("Test", "Body message.")
// m.To = []string{"to@gmail.com"}
// m.AttachFile("/path/to/file")
// fmt.Println(s.Send(m))
