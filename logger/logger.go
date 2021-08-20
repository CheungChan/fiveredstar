package logger

import (
    "github.com/gin-gonic/gin"
    "github.com/rs/zerolog"
    "github.com/rs/zerolog/log"
    "gopkg.in/natefinch/lumberjack.v2"
    "io"
    "os"
    "path"
    "time"
)

var Logger zerolog.Logger

func New(directory string, fileName string, withCaller bool, maxBackups int, maxSize int, maxAge int, ) zerolog.Logger {
    Logger = log.Output(zerolog.MultiLevelWriter(
        zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339},
        zerolog.ConsoleWriter{Out: newRollingFile(directory, fileName, maxBackups, maxSize, maxAge), NoColor: true, TimeFormat: time.RFC3339})).With().Timestamp().Logger()
    if withCaller {
        Logger = Logger.With().Caller().Logger()
    }
    return Logger
}

func newRollingFile(directory string, fileName string, maxBackups int, maxSize int, maxAge int) io.Writer {
    if err := os.MkdirAll(directory, 0744); err != nil {
        log.Error().Err(err).Str("path", directory).Msg("can't create log directory")
        return nil
    }

    return &lumberjack.Logger{
        Filename:   path.Join(directory, fileName),
        MaxBackups: maxBackups, // files
        MaxSize:    maxSize,    // megabytes
        MaxAge:     maxAge,     // days
    }
}

type ginHands struct {
    SerName    string
    Path       string
    Latency    time.Duration
    Method     string
    StatusCode int
    ClientIP   string
    MsgStr     string
}

func GinMiddleware(serName string) gin.HandlerFunc {
    return func(c *gin.Context) {
        t := time.Now()
        // before request
        p := c.Request.URL.Path
        raw := c.Request.URL.RawQuery
        c.Next()
        // after request
        // latency := time.Since(t)
        // clientIP := c.ClientIP()
        // method := c.Request.Method
        // statusCode := c.Writer.Status()
        if raw != "" {
            p = p + "?" + raw
        }
        msg := c.Errors.String()
        if msg == "" {
            msg = "Request"
        }
        cData := &ginHands{
            SerName:    serName,
            Path:       p,
            Latency:    time.Since(t),
            Method:     c.Request.Method,
            StatusCode: c.Writer.Status(),
            ClientIP:   c.ClientIP(),
            MsgStr:     msg,
        }

        logSwitch(cData)
    }
}

func logSwitch(data *ginHands) {
    switch {
    case data.StatusCode >= 400 && data.StatusCode < 500:
        {
            Logger.Warn().Str("ser_name", data.SerName).Str("method", data.Method).Str("path", data.Path).Dur("resp_time", data.Latency).Int("status", data.StatusCode).Str("client_ip", data.ClientIP).Msg(data.MsgStr)
        }
    case data.StatusCode >= 500:
        {
            Logger.Error().Str("ser_name", data.SerName).Str("method", data.Method).Str("path", data.Path).Dur("resp_time", data.Latency).Int("status", data.StatusCode).Str("client_ip", data.ClientIP).Msg(data.MsgStr)
        }
    default:
        Logger.Info().Str("ser_name", data.SerName).Str("method", data.Method).Str("path", data.Path).Dur("resp_time", data.Latency).Int("status", data.StatusCode).Str("client_ip", data.ClientIP).Msg(data.MsgStr)
    }
}
