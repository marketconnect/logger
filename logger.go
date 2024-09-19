package logger

import (
	"fmt"

	"net/http"
	"net/url"
	"os"
)

type Logger interface {
	Info(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
}

type TelegramLogger struct {
	serviceName string
	logger      Logger
	botToken    string
	chatID      string
}

func NewTelegramLogger(botToken string, chatID string, serviceName string, logger Logger) *TelegramLogger {
	return &TelegramLogger{
		serviceName: serviceName,
		logger:      logger,
		botToken:    botToken,
		chatID:      chatID,
	}
}

func (t *TelegramLogger) Info(args ...interface{}) {
	message := fmt.Sprint(args...)
	t.logger.Info("INFO:", message)
	if err := t.SendToTelegram("ℹ️ " + "<b>" + t.serviceName + "</b>: " + message); err != nil {
		t.logger.Error("ERROR: could not send to Telegram:", err)
	}
}

func (t *TelegramLogger) Error(args ...interface{}) {
	message := fmt.Sprint(args...)
	t.logger.Error("ERROR:", message)
	if err := t.SendToTelegram("❗ " + "<b>" + t.serviceName + "</b>: " + message); err != nil {
		t.logger.Error("ERROR: could not send to Telegram:", err)
	}
}

func (t *TelegramLogger) Fatal(args ...interface{}) {
	message := fmt.Sprint(args...)
	t.logger.Error("FATAL:", message)
	if err := t.SendToTelegram("🚨 " + "<b>" + t.serviceName + "</b>: " + message); err != nil {
		t.logger.Error("ERROR: could not send to Telegram:", err)
	}
	os.Exit(1)
}

func (t *TelegramLogger) SendToTelegram(message string) error {
	telegramAPI := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", t.botToken)

	data := url.Values{}
	data.Set("chat_id", t.chatID)
	data.Set("text", message)
	data.Set("parse_mode", "HTML")

	resp, err := http.PostForm(telegramAPI, data)
	if err != nil {
		return fmt.Errorf("could nod  Telegram: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("неожиданный статус код от Telegram: %d", resp.StatusCode)
	}

	return nil
}
