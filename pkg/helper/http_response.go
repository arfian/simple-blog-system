package helper

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type Response struct {
	Data      interface{} `json:"data"`
	Success   bool        `json:"success"`
	Message   string      `json:"message"`
	RequestId any         `json:"request_id"`
}

type ginHands struct {
	Path       string
	Latency    time.Duration
	Method     string
	StatusCode int
	ClientIP   string
	MsgStr     string
	RequestId  string
}

type ResponseErrorData struct {
	Type string `json:"data"`
	Code int64  `json:"success"`
}

func ResponseData(c *gin.Context, res *Response) {
	requestID, _ := c.Get("requestID")
	res.Success = true
	res.RequestId = requestID

	SaveAuditLog(c, res.Message)
	c.JSON(200, res)
}

func GetIpAddress(c *gin.Context) string {
	ip := c.GetHeader("X-Real-IP")
	if ip == "" {
		ip = c.GetHeader("X-Forwarded-For")
	}

	if ip == "" {
		ip = c.ClientIP()
	}
	return ip
}

func SaveAuditLog(c *gin.Context, msg string) {
	timeStart, _ := c.Get("timeStart")
	timeParsed, _ := time.Parse(time.RFC3339, timeStart.(string))

	requestID, _ := c.Get("requestID")

	path := c.Request.URL.Path
	raw := c.Request.URL.RawQuery
	if raw != "" {
		path = path + "?" + raw
	}
	cData := &ginHands{
		Path:       path,
		Latency:    time.Since(timeParsed),
		Method:     c.Request.Method,
		StatusCode: c.Writer.Status(),
		ClientIP:   GetIpAddress(c),
		MsgStr:     msg,
		RequestId:  requestID.(string),
	}

	logSwitch(cData)
}

func logSwitch(data *ginHands) {
	switch {
	case data.StatusCode >= 400 && data.StatusCode < 500:
		log.Warn().Str("request_id", data.RequestId).Str("method", data.Method).Str("path", data.Path).Dur("resp_time", data.Latency).Int("status", data.StatusCode).Str("client_ip", data.ClientIP).Msg(data.MsgStr)
	case data.StatusCode >= 500:
		log.Error().Str("request_id", data.RequestId).Str("method", data.Method).Str("path", data.Path).Dur("resp_time", data.Latency).Int("status", data.StatusCode).Str("client_ip", data.ClientIP).Msg(data.MsgStr)
	default:
		log.Info().Str("request_id", data.RequestId).Str("method", data.Method).Str("path", data.Path).Dur("resp_time", data.Latency).Int("status", data.StatusCode).Str("client_ip", data.ClientIP).Msg(data.MsgStr)
	}
}

func ResponseError(c *gin.Context, err error, opts ...interface{}) {
	t := "InternalServerError"
	d := err.Error()

	code := http.StatusInternalServerError

	// if request cancelled
	if c.Request.Context().Err() == context.Canceled {
		c.AbortWithStatus(http.StatusNoContent)
		return
	}

	for _, v := range opts {
		if tErr, ok := v.(string); ok {
			if strings.Contains(tErr, " ") {
				d = tErr
			} else {
				t = tErr
			}
		}
		if cErr, ok := v.(int); ok && cErr >= 100 && cErr <= 599 {
			code = cErr
		}
	}

	if errors.Is(err, gorm.ErrRecordNotFound) || strings.Contains(err.Error(), "not found") {
		code = http.StatusNotFound
		t = "NotFound"
	}

	requestID, _ := c.Get("requestID")
	SaveAuditLog(c, d)
	c.AbortWithStatusJSON(code, &Response{
		Success: false,
		Message: d,
		Data: &ResponseErrorData{
			Type: t,
			Code: int64(code),
		},
		RequestId: requestID,
	})
}
