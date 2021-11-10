package web

import (
	"bytes"
	"errors"
	"himawari8Capturer/core"
	"io"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func Run(addr ...string) error {
	r := gin.Default()
	configure(r)
	return r.Run(addr...)
}

// 请求管道配置
func configure(r *gin.Engine) {
	r.GET("/earth", func(c *gin.Context) {
		q, t, s := c.Query("quality"), c.Query("time"), c.Query("shorelines")
		var quality core.Quality
		var _time time.Time
		var shorelines core.Shoreline
		if q == "" {
			quality = core.HD
		} else if qq, err := strconv.Atoi(q); err != nil {
			log.Printf("%s", err)
		} else {
			quality = core.Quality(qq)
		}
		if t == "" {
			tt, err := core.GetLatestEarthTime()
			if err != nil {
				c.AbortWithError(500, errors.New("服务器异常"))
			}
			_time = tt
		} else if tt, err := strconv.ParseInt(t, 10, 64); err != nil {
			log.Printf("%s", err)
		} else {
			_time = time.UnixMilli(tt)
		}
		switch ss := strings.ToUpper(s); ss {
		case "RED":
			shorelines = core.Red
		case "GREEN":
			shorelines = core.Green
		case "YELLOW":
			shorelines = core.Yellow
		}
		earth, err := core.GetEarthWithShorelines(quality, _time, shorelines)
		if err != nil {
			c.AbortWithError(500, errors.New("服务器异常"))
		}
		var buf bytes.Buffer
		_, err = io.Copy(&buf, earth)
		if err != nil {
			c.AbortWithError(500, errors.New("服务器异常"))
		}
		c.Data(200, "image/png", buf.Bytes())
		c.Abort()
	})
}
