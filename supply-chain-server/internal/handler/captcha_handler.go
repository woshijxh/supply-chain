package handler

import (
	"net/http"

	"supply-chain-server/pkg/captcha"

	"github.com/gin-gonic/gin"
)

type CaptchaHandler struct{}

func NewCaptchaHandler() *CaptchaHandler {
	return &CaptchaHandler{}
}

func (h *CaptchaHandler) Get(c *gin.Context) {
	captchaData := captcha.Generate()
	img, err := captcha.GenerateImage(captchaData.Expr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成验证码失败"})
		return
	}

	c.Header("X-Captcha-ID", captchaData.ID)
	c.Header("Cache-Control", "no-store")
	c.Data(http.StatusOK, "image/png", img)
}
