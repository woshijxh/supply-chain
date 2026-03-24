package handler

import (
	"supply-chain-server/internal/service"
	"supply-chain-server/pkg/response"

	"github.com/gin-gonic/gin"
)

type TraceHandler struct {
	service *service.TraceService
}

func NewTraceHandler(s *service.TraceService) *TraceHandler {
	return &TraceHandler{service: s}
}

// Trace 追溯商品
// GET /api/trace?code=xxx
func (h *TraceHandler) Trace(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		response.BadRequest(c, "请输入追溯编码")
		return
	}

	result, err := h.service.Trace(code)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, result)
}