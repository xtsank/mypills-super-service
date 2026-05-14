package middleware

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/samber/do/v2"
	apperrors "github.com/xtsank/mypills-super-service/src/internal/errors"
)

func NewLogger(i do.Injector) (*slog.Logger, error) {
	_ = godotenv.Load()

	logPath := os.Getenv("LOG_FILE")
	if logPath == "" {
		logPath = "logs/app.log"
	}
	_ = os.MkdirAll(filepath.Dir(logPath), 0755)
	file, _ := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	logLevel, logLevelValid, rawLevel := parseLogLevelEnv("LOG_LEVEL", slog.LevelInfo)

	opts := &slog.HandlerOptions{
		AddSource: false,
		Level:     logLevel,
	}

	fileHandler := slog.NewJSONHandler(file, opts)
	logger := slog.New(fileHandler)
	if !logLevelValid && rawLevel != "" {
		logger.Warn("Invalid log level, using default", slog.String("value", rawLevel), slog.String("default", logLevel.String()))
	}

	slog.SetDefault(logger)

	return logger, nil
}

func parseLogLevelEnv(key string, fallback slog.Level) (slog.Level, bool, string) {
	val := strings.ToLower(strings.TrimSpace(os.Getenv(key)))
	if val == "" {
		return fallback, true, ""
	}

	switch val {
	case "debug":
		return slog.LevelDebug, true, val
	case "info":
		return slog.LevelInfo, true, val
	case "warn", "warning":
		return slog.LevelWarn, true, val
	case "error":
		return slog.LevelError, true, val
	default:
		return fallback, false, val
	}
}

func Logger(i do.Injector) gin.HandlerFunc {
	baseLogger := do.MustInvoke[*slog.Logger](i)

	return func(c *gin.Context) {
		start := time.Now()
		eventID := uuid.New().String()

		c.Next()

		latency := time.Since(start)
		action := fmt.Sprintf("%s %s", c.Request.Method, c.Request.URL.Path)

		requestLogger := baseLogger.With(
			slog.String("event_id", eventID),
			slog.String("user_action", action),
			slog.Int("status", c.Writer.Status()),
			slog.Int64("latency_ms", latency.Milliseconds()),
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
		)

		if baseLogger.Enabled(c.Request.Context(), slog.LevelDebug) {
			requestLogger = requestLogger.With(
				slog.String("client_ip", c.ClientIP()),
				slog.String("user_agent", c.Request.UserAgent()),
				slog.String("query", c.Request.URL.RawQuery),
			)
		}

		if len(c.Errors) == 0 {
			requestLogger.Info("request_success")
			return
		}

		err := c.Errors.Last().Err
		requestLogger = requestLogger.With(
			slog.String("error_type", fmt.Sprintf("%T", err)),
		)
		if baseLogger.Enabled(c.Request.Context(), slog.LevelDebug) {
			requestLogger = requestLogger.With(
				slog.String("error_stack", string(debug.Stack())),
			)
		}

		var appErr *apperrors.AppError
		if errors.As(err, &appErr) {
			requestLogger = requestLogger.With(
				slog.String("error_code", appErr.Code),
				slog.String("error_message", appErr.Message),
				slog.Int("error_status", appErr.HTTPStatus),
			)
			if appErr.SourceFile != "" {
				requestLogger = requestLogger.With(
					slog.String("error_source_file", appErr.SourceFile),
					slog.Int("error_source_line", appErr.SourceLine),
					slog.String("error_source_func", appErr.SourceFunc),
				)
			}
			if cause := errors.Unwrap(appErr); cause != nil {
				requestLogger = requestLogger.With(
					slog.String("error_cause", cause.Error()),
					slog.String("error_cause_type", fmt.Sprintf("%T", cause)),
				)
			}
		} else {
			requestLogger = requestLogger.With(
				slog.String("error_message", err.Error()),
			)
		}

 		requestLogger.Error("request_failed")
	}
}

