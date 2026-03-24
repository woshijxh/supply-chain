package handler

import (
	"net/http"
	"strconv"
	"supply-chain-server/internal/model"
	"supply-chain-server/internal/service"
	"supply-chain-server/pkg/response"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	service *service.ProductService
}

func NewProductHandler(s *service.ProductService) *ProductHandler {
	return &ProductHandler{service: s}
}

func (h *ProductHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	keyword := c.Query("keyword")

	products, total, err := h.service.List(page, pageSize, keyword)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.PageSuccess(c, products, total, page, pageSize)
}

func (h *ProductHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的产品ID")
		return
	}

	product, err := h.service.GetByID(uint(id))
	if err != nil {
		response.Error(c, http.StatusNotFound, "产品不存在")
		return
	}
	response.Success(c, product)
}

func (h *ProductHandler) Create(c *gin.Context) {
	var product model.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误")
		return
	}

	if err := h.service.Create(&product); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, product)
}

func (h *ProductHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的产品ID")
		return
	}

	var product model.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误")
		return
	}

	product.ID = uint(id)
	if err := h.service.Update(&product); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, product)
}

func (h *ProductHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的产品ID")
		return
	}

	if err := h.service.Delete(uint(id)); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, nil)
}
