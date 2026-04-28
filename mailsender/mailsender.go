package mailsender

import (
	"fmt"
	"net/smtp"
	"strings"
)

// Config はSMTP接続設定を保持する。
type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	From     string
}

// Sender はメール送信クライアント。
type Sender struct {
	cfg Config
}

// New はSenderを生成する。
func New(cfg Config) *Sender {
	return &Sender{cfg: cfg}
}

// Send は指定した宛先にメールを送信する。
// STARTTLS（ポート587）を使用する。
func (s *Sender) Send(to []string, subject, body string) error {
	auth := smtp.PlainAuth("", s.cfg.User, s.cfg.Password, s.cfg.Host)
	msg := buildMessage(s.cfg.From, to, subject, body)
	addr := s.cfg.Host + ":" + s.cfg.Port
	return smtp.SendMail(addr, auth, s.cfg.From, to, []byte(msg))
}

func buildMessage(from string, to []string, subject, body string) string {
	return fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s",
		from, strings.Join(to, ", "), subject, body,
	)
}
