package logger

import (
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

func TestGinMiddleware(t *testing.T) {
	Logger := New("logs", "test_gin_logger.log", true,3, 30*1024*1024, 30)
	r := gin.Default()
	r.Use(GinMiddleware("test_gin"))
	r.GET("/", func(c *gin.Context) {
		Logger.Info().Msg("gin view log")
		c.String(200, "Hello World")
	})
	go func() {
		r.Run(":8080")
	}()
	time.Sleep(1)
	res, err := http.Get("http://127.0.0.1:8080/")
	if err != nil {
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			Logger.Error().Msgf("%+v",err)
		}
	}(res.Body)
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	Logger.Info().Msg("请求结果" + string(b))
}

func TestNew(t *testing.T) {
	Logger := New("logs", "test_logger.log", true,3, 30*1024*1024, 30)
	Logger.Info().Msg("Hello")
}
