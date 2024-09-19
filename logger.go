package logger

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strings"
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

func getCallerInfo() string {
	for i := 2; i < 10; i++ {
		if pc, file, line, ok := runtime.Caller(i); ok {
			fn := runtime.FuncForPC(pc)

			if !strings.Contains(file, "logger") {
				return fmt.Sprintf("%s:%d (%s)", file, line, fn.Name())
			}
		}
	}
	return "неизвестный файл:0"
}

func (t *TelegramLogger) Info(args ...interface{}) {
	callerInfo := getCallerInfo()
	message := fmt.Sprint(args...)
	logMessage := fmt.Sprintf("INFO: %s [%s]", message, callerInfo)
	t.logger.Info(logMessage)
	if err := t.SendToTelegram("ℹ️ " + "<b>" + t.serviceName + "</b>: " + message + " [" + callerInfo + "]"); err != nil {
		t.logger.Error("ERROR: could not send to Telegram:", err)
	}
}

func (t *TelegramLogger) Error(args ...interface{}) {
	callerInfo := getCallerInfo()
	message := fmt.Sprint(args...)
	logMessage := fmt.Sprintf("ERROR: %s [%s]", message, callerInfo)
	t.logger.Error(logMessage)
	if err := t.SendToTelegram("❗ " + "<b>" + t.serviceName + "</b>: " + message + " [" + callerInfo + "]"); err != nil {
		t.logger.Error("ERROR: could not send to Telegram:", err)
	}
}

func (t *TelegramLogger) Fatal(args ...interface{}) {
	callerInfo := getCallerInfo()
	message := fmt.Sprint(args...)
	logMessage := fmt.Sprintf("FATAL: %s [%s]", message, callerInfo)
	t.logger.Error(logMessage)
	if err := t.SendToTelegram("🚨 " + "<b>" + t.serviceName + "</b>: " + message + " [" + callerInfo + "]"); err != nil {
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
		return fmt.Errorf("could not send message to Telegram: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("неожиданный статус код от Telegram: %d", resp.StatusCode)
	}

	return nil
}
