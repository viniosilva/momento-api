package adapters

import (
	"context"
	"fmt"
	"net/smtp"
)

const SendResetPasswordEmailBody = `Olá,

Você solicitou a recuperação de senha da sua conta no Momento.

Clique no link abaixo para redefinir sua senha (válido por 1 hora):

%s

Se você não solicitou esta recuperação, ignore este email.

Atenciosamente,
Equipe Momento`

type emailService struct {
	host         string
	port         string
	user         string
	pass         string
	from         string
	resetURLBase string
}

func NewEmailService(host, user, pass, from, resetURLBase, port string) *emailService {
	return &emailService{
		host:         host,
		port:         port,
		user:         user,
		pass:         pass,
		from:         from,
		resetURLBase: resetURLBase,
	}
}

func (s *emailService) SendResetPasswordEmail(ctx context.Context, to, token string) error {
	resetLink := fmt.Sprintf("%s?token=%s", s.resetURLBase, token)

	subject := "Recuperação de Senha - Momento"
	body := fmt.Sprintf(SendResetPasswordEmailBody, resetLink)
	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s", s.from, to, subject, body)
	addr := fmt.Sprintf("%s:%s", s.host, s.port)
	var auth smtp.Auth
	if s.user != "" && s.pass != "" {
		auth = smtp.PlainAuth("", s.user, s.pass, s.host)
	}

	if err := smtp.SendMail(addr, auth, s.from, []string{to}, []byte(msg)); err != nil {
		return fmt.Errorf("smtp.SendMail: %w", err)
	}

	return nil
}
