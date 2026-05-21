package email

import (
	"context"
	"fmt"
	"net/smtp"

	"github.com/samber/do/v2"
	"github.com/xtsank/mypills-super-service/src/internal/infra/postgres/config"
	"github.com/xtsank/mypills-super-service/src/internal/service"
)

type SMTPSender struct {
	host     string
	port     string
	user     string
	password string
	from     string
}

func NewSMTPSender(i do.Injector) (service.EmailSender, error) {
	cfg := do.MustInvoke[*config.Config](i)
	return &SMTPSender{
		host:     cfg.SMTPHost,
		port:     cfg.SMTPPort,
		user:     cfg.SMTPUser,
		password: cfg.SMTPPassword,
		from:     cfg.SMTPFrom,
	}, nil
}

func (s *SMTPSender) Send(ctx context.Context, to string, subject string, body string) error {
	_ = ctx

	if s.host == "" || s.port == "" || s.from == "" {
		return fmt.Errorf("smtp config is incomplete")
	}

	addr := fmt.Sprintf("%s:%s", s.host, s.port)
	auth := smtp.PlainAuth("", s.user, s.password, s.host)

	message := "Subject: " + subject + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/plain; charset=UTF-8\r\n" +
		"\r\n" + body

	return smtp.SendMail(addr, auth, s.from, []string{to}, []byte(message))
}
