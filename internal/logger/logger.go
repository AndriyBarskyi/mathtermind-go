package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"runtime"
	"strings"
	"time"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorCyan   = "\033[36m"
	colorGray   = "\033[90m"
)

type devHandler struct {
	handler slog.Handler
	out     io.Writer
}

func (h *devHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

func (h *devHandler) Handle(ctx context.Context, r slog.Record) error {
	var level string
	var color string

	switch r.Level {
	case slog.LevelDebug:
		level = "DBG"
		color = colorGray
	case slog.LevelInfo:
		level = "INF"
		color = colorBlue
	case slog.LevelWarn:
		level = "WRN"
		color = colorYellow
	case slog.LevelError:
		level = "ERR"
		color = colorRed
	default:
		level = "???"
	}

	// Format: [TIME] LEVEL message key=value key=value ...
	fmt.Fprintf(h.out, "%s[%s]%s %s%-5s%s %s",
		colorGray, r.Time.Format(time.TimeOnly), colorReset,
		color, level, colorReset,
		r.Message)

	if r.PC != 0 {
		fs := runtime.CallersFrames([]uintptr{r.PC})
		if f, _ := fs.Next(); f.File != "" {
			file := f.File
			if idx := strings.LastIndex(file, "/"); idx >= 0 {
				file = file[idx+1:]
			}
			fmt.Fprintf(h.out, " %s(%s:%d)%s", colorGray, file, f.Line, colorReset)
		}
	}

	r.Attrs(func(attr slog.Attr) bool {
		fmt.Fprintf(h.out, " %s%s%s=%v", colorCyan, attr.Key, colorReset, attr.Value)
		return true
	})

	fmt.Fprintln(h.out)
	return nil
}

func (h *devHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &devHandler{handler: h.handler.WithAttrs(attrs), out: h.out}
}

func (h *devHandler) WithGroup(name string) slog.Handler {
	return &devHandler{handler: h.handler.WithGroup(name), out: h.out}
}

func NewDevelopment() *slog.Logger {
	return slog.New(&devHandler{
		handler: slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}),
		out: os.Stdout,
	})
}
