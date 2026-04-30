package mailsender

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"mime/quotedprintable"
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
	// Encoding はメール本文のエンコーディング。"base64" / "quoted-printable" / "" (未設定時はそのまま送信)
	Encoding string
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
func (s *Sender) Send(to []string, subject, body string) error {
	auth := smtp.PlainAuth("", s.cfg.User, s.cfg.Password, s.cfg.Host)
	msg := buildMessage(s.cfg.From, to, subject, body, s.cfg.Encoding)
	addr := s.cfg.Host + ":" + s.cfg.Port
	return smtp.SendMail(addr, auth, s.cfg.From, to, msg)
}

func buildMessage(from string, to []string, subject, body, encoding string) []byte {
	cte, encodedBody := encodeBody(body, encoding)

	header := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=UTF-8\r\n",
		from, strings.Join(to, ", "), subject,
	)
	if cte != "" {
		header += "Content-Transfer-Encoding: " + cte + "\r\n"
	}
	return append([]byte(header+"\r\n"), encodedBody...)
}

func encodeBody(body, encoding string) (cte string, encoded []byte) {
	switch strings.ToLower(encoding) {
	case "base64":
		raw := base64.StdEncoding.EncodeToString([]byte(body))
		var sb strings.Builder
		for i := 0; i < len(raw); i += 76 {
			end := i + 76
			if end > len(raw) {
				end = len(raw)
			}
			sb.WriteString(raw[i:end])
			sb.WriteString("\r\n")
		}
		return "base64", []byte(sb.String())
	case "quoted-printable":
		var buf bytes.Buffer
		w := quotedprintable.NewWriter(&buf)
		w.Write([]byte(body))
		w.Close()
		return "quoted-printable", buf.Bytes()
	default:
		return "", []byte(body)
	}
}
