package email

import (
	"net/smtp"
	"strings"

	"github.com/itfantasy/gonode/components/common"
)

type Email struct {
	user string
	pass string
	host string
	auth smtp.Auth
	opts *common.CompOptions
}

func NewEmail() *Email {
	e := new(Email)

	e.opts = common.NewCompOptions()
	return e
}

func (e *Email) Conn(url string, host string) error {
	ip := strings.Split(url, ":")
	auth := smtp.PlainAuth("", e.user, e.pass, ip[0])
	e.auth = auth
	e.host = url
	return nil
}

func (e *Email) SetAuthor(user string, pass string) {
	e.user = strings.Replace(user, "#", "@", -1)
	e.pass = pass
}

func (e *Email) SetOption(key string, val interface{}) {

}

func (e *Email) Close() {

}

func (e *Email) SendTo(to string, title string, content string) error {
	content_type := "Content-Type: text/html; charset=UTF-8"
	msg := []byte("To: " + to + "\r\nFrom: " + e.user + ">\r\nSubject: " + title + "\r\n" + content_type + "\r\n\r\n" + content)
	send_to := strings.Split(to, ",")
	err := smtp.SendMail(e.host, e.auth, e.user, send_to, msg)
	return err
}
