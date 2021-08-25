package logger

import (
    "github.com/gin-gonic/gin"
    "github.com/gogf/gf/os/glog"
    "github.com/pkg/errors"
    "io/ioutil"
    "net/http"
    "runtime"
    "testing"
)

func TestNew(t *testing.T) {
    InitLog("logs/test_logger", "10M", false)
    glog.Info("Hello")
    glog.Error(errors.New("这是一个错误,TestNew"))
}
func TestGinMiddleware(t *testing.T) {
    InitLog("logs/test_gin_logger", "10M", true)
    r := gin.New()
    r.Use(GinMiddleware("test_gin"))
    r.GET("/", func(c *gin.Context) {
        glog.Info("gin view log")
        glog.Error(errors.New("这里有一个错误在handler里"))
        c.String(200, "Hello World")
    })
    go func() {
        r.Run(":8080")
    }()
    runtime.Gosched()
    res, err := http.Get("http://127.0.0.1:8080/")
    if err != nil {
        glog.Error(err)
    }
    defer res.Body.Close()
    _, err = ioutil.ReadAll(res.Body)
    if err != nil {
        glog.Error(err)
    }

}
