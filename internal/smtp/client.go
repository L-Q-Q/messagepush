package smtp

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
)

// SMTPClient SMTP 客户端
type SMTPClient struct{}

// EmailConfig 邮件配置
type EmailConfig struct {
	Server   string
	Port     int
	Username string
	Password string
}

// Email 邮件内容
type Email struct {
	From    string
	To      string
	Subject string
	Body    string
}

// NewSMTPClient 创建 SMTP 客户端实例
func NewSMTPClient() *SMTPClient {
	return &SMTPClient{}
}

// Send 发送邮件
func (c *SMTPClient) Send(config *EmailConfig, email *Email) error {
	// 构建邮件地址
	addr := fmt.Sprintf("%s:%d", config.Server, config.Port)

	// 构建认证信息
	auth := smtp.PlainAuth("", config.Username, config.Password, config.Server)

	// 构建邮件内容（添加 MIME 头）
	message := fmt.Sprintf("MIME-Version: 1.0\r\n"+
		"Content-Type: text/plain; charset=UTF-8\r\n"+
		"From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"\r\n"+
		"%s\r\n",
		email.From, email.To, email.Subject, email.Body)

	// 根据端口判断是否使用 TLS
	if config.Port == 465 {
		// 使用 SSL/TLS
		return c.sendWithTLS(addr, auth, email.From, []string{email.To}, []byte(message), config.Server)
	} else if config.Port == 587 {
		// 使用 STARTTLS
		return c.sendWithSTARTTLS(addr, auth, email.From, []string{email.To}, []byte(message), config.Server)
	}

	// 使用明文（不推荐）
	return smtp.SendMail(addr, auth, email.From, []string{email.To}, []byte(message))
}

// sendWithTLS 使用 TLS 发送邮件
func (c *SMTPClient) sendWithTLS(addr string, auth smtp.Auth, from string, to []string, msg []byte, serverName string) error {
	// 创建 TLS 配置
	tlsConfig := &tls.Config{
		ServerName:         serverName,
		InsecureSkipVerify: false,
	}

	// 建立 TLS 连接
	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return fmt.Errorf("failed to dial: %w", err)
	}
	defer conn.Close()

	// 创建 SMTP 客户端
	client, err := smtp.NewClient(conn, serverName)
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}
	defer client.Close()

	// 认证
	if auth != nil {
		if err = client.Auth(auth); err != nil {
			return fmt.Errorf("failed to authenticate: %w", err)
		}
	}

	// 设置发件人
	if err = client.Mail(from); err != nil {
		return fmt.Errorf("failed to set sender: %w", err)
	}

	// 设置收件人
	for _, addr := range to {
		if err = client.Rcpt(addr); err != nil {
			return fmt.Errorf("failed to set recipient: %w", err)
		}
	}

	// 发送邮件内容
	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to get data writer: %w", err)
	}

	_, err = w.Write(msg)
	if err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	err = w.Close()
	if err != nil {
		return fmt.Errorf("failed to close data writer: %w", err)
	}

	return client.Quit()
}

// sendWithSTARTTLS 使用 STARTTLS 发送邮件
func (c *SMTPClient) sendWithSTARTTLS(addr string, auth smtp.Auth, from string, to []string, msg []byte, serverName string) error {
	// 建立普通连接
	client, err := smtp.Dial(addr)
	if err != nil {
		return fmt.Errorf("failed to dial: %w", err)
	}
	defer client.Close()

	// 启动 TLS
	tlsConfig := &tls.Config{
		ServerName:         serverName,
		InsecureSkipVerify: false,
	}

	if err = client.StartTLS(tlsConfig); err != nil {
		return fmt.Errorf("failed to start TLS: %w", err)
	}

	// 认证
	if auth != nil {
		if err = client.Auth(auth); err != nil {
			return fmt.Errorf("failed to authenticate: %w", err)
		}
	}

	// 设置发件人
	if err = client.Mail(from); err != nil {
		return fmt.Errorf("failed to set sender: %w", err)
	}

	// 设置收件人
	for _, addr := range to {
		if err = client.Rcpt(addr); err != nil {
			return fmt.Errorf("failed to set recipient: %w", err)
		}
	}

	// 发送邮件内容
	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to get data writer: %w", err)
	}

	_, err = w.Write(msg)
	if err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	err = w.Close()
	if err != nil {
		return fmt.Errorf("failed to close data writer: %w", err)
	}

	return client.Quit()
}
