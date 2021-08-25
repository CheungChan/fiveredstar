package logger

import (
    "github.com/gin-gonic/gin"
    "github.com/gogf/gf/frame/g"
    "github.com/gogf/gf/os/glog"
    "time"
)

func InitLog(logDirectory string, file string, rotateSize string, writerColorEnable bool) {
    if file == "" {
        file = "{Y-m-d}.log"
    }
    err := glog.SetConfigWithMap(g.Map{
        "path":              logDirectory,
        "file":              file, // 日志文件格式。默认为"{Y-m-d}.log"
        "level":             "all",
        "stdout":            true,
        "StStatus":          1,
        "rotateSize":        rotateSize,
        "rotateExpire":      24 * time.Hour,
        "writerColorEnable": writerColorEnable,
    })
    glog.SetFlags(glog.F_TIME_STD | glog.F_FILE_SHORT)
    if err != nil {
        glog.Error("日志配置错误")
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
        glog.Warningf("%s|%d|%s|%s|%s|%s|%s", data.SerName, data.StatusCode, data.Method, data.Latency, data.ClientIP, data.Path, data.MsgStr)
    case data.StatusCode >= 500:
        glog.Errorf("%s|%d|%s|%s|%s|%s|%s", data.SerName, data.StatusCode, data.Method, data.Latency, data.ClientIP, data.Path, data.MsgStr)
    default:
        glog.Infof("%s|%d|%s|%s|%s|%s|%s", data.SerName, data.StatusCode, data.Method, data.Latency, data.ClientIP, data.Path, data.MsgStr)
    }
}
