package logger

import (
	"fmt"

	"net/http"
	"net/url"
	"os"
)

// Logger определяет интерфейс для логирования.
type Logger interface {
	Info(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
}

// TelegramLogger реализует интерфейс Logger и отправляет сообщения в Telegram.
type TelegramLogger struct {
	logger   Logger
	botToken string
	chatID   string
}

// NewTelegramLogger создает новый экземпляр TelegramLogger.
func NewTelegramLogger(botToken string, chatID string, logger Logger) *TelegramLogger {
	return &TelegramLogger{
		logger:   logger,
		botToken: botToken,
		chatID:   chatID,
	}
}

// Info выводит информационное сообщение и отправляет его в Telegram.
func (t *TelegramLogger) Info(args ...interface{}) {
	message := fmt.Sprint(args...)
	t.logger.Info("INFO:", message)
	if err := t.SendToTelegram("ℹ️ " + message); err != nil {
		t.logger.Error("ERROR: could not send to Telegram:", err)
	}
}

// Error выводит сообщение об ошибке и отправляет его в Telegram.
func (t *TelegramLogger) Error(args ...interface{}) {
	message := fmt.Sprint(args...)
	t.logger.Error("ERROR:", message)
	if err := t.SendToTelegram("❗ " + message); err != nil {
		t.logger.Error("ERROR: could not send to Telegram:", err)
	}
}

// Fatal выводит фатальное сообщение, отправляет его в Telegram и завершает приложение.
func (t *TelegramLogger) Fatal(args ...interface{}) {
	message := fmt.Sprint(args...)
	t.logger.Error("FATAL:", message)
	if err := t.SendToTelegram("🚨 " + message); err != nil {
		t.logger.Error("ERROR: could not send to Telegram:", err)
	}
	os.Exit(1)
}

// SendToTelegram отправляет сообщение в указанный чат Telegram.
func (t *TelegramLogger) SendToTelegram(message string) error {
	telegramAPI := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", t.botToken)

	// Подготовка данных для отправки сообщения
	data := url.Values{}
	data.Set("chat_id", t.chatID)
	data.Set("text", message)
	data.Set("parse_mode", "HTML") // Используем HTML для форматирования

	// Выполнение POST-запроса к API Telegram
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
