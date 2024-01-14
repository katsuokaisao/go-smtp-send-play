package main

import (
	"fmt"
	"io"
	"net/smtp"
	"strings"

	"github.com/spf13/viper"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

type Config struct {
	Smtp Smtp
}

type Smtp struct {
	Server   string
	Port     int
	Address  string
	Password string
}

func main() {
	cofig, err := loadConfig("config", "smtp", "toml")
	if err != nil {
		panic(fmt.Errorf("error in loadConfigFromToml: %w", err))
	}

	toMailAddress := "hogehoge@test.com"
	if err := sendMail(toMailAddress, cofig.Smtp); err != nil {
		panic(fmt.Errorf("error in sendMail: %w", err))
	}
	fmt.Println("[info]: success send mail!!")
}

func loadConfig(dir, fileName, fileType string) (*Config, error) {
	v := viper.New()
	v.SetConfigName(fileName)
	v.SetConfigType(fileType)
	v.AddConfigPath(dir)
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error in viper ReadInConfig: %w", err)
	}

	var cofig Config
	if err := v.Unmarshal(&cofig); err != nil {
		return nil, fmt.Errorf("error in viper Unmarshal: %w", err)
	}

	return &cofig, nil
}

func sendMail(toMailAddress string, smtpConn Smtp) error {
	toMailAddressList := []string{toMailAddress}
	mailSubject := "テストメール"
	mailBody := "テストメールです。"

	// Construct the email message
	msgStr := "To:" + toMailAddress + "\r\n" +
		"reply-to: hogehoge@test.com" + "\r\n" +
		"Subject:" + mailSubject + "\r\n" +
		"\r\n" + mailBody

	reader := strings.NewReader(msgStr)
	transformer := japanese.ISO2022JP.NewEncoder()
	msgISO2022JP, err := io.ReadAll(transform.NewReader(reader, transformer))
	if err != nil {
		return fmt.Errorf("unable to convert to ISO2022JP: %w", err)
	}
	msg := []byte(msgISO2022JP)

	// Connect to the SMTP server
	auth := smtp.PlainAuth("", smtpConn.Address, smtpConn.Password, smtpConn.Server)

	// Send the email
	if err := smtp.SendMail(fmt.Sprintf("%s:%d", smtpConn.Server, smtpConn.Port), auth, smtpConn.Address, toMailAddressList, msg); err != nil {
		return fmt.Errorf("error in smtp SendMail: %w", err)
	}

	return nil
}
