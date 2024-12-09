package handler

import (
	"net/http"
	"personal-api/internal/entity"
	"personal-api/internal/service"
	"personal-api/pkg/response"

	"github.com/gin-gonic/gin"
)

type ContactInfoHandler struct {
	service service.ContactInfoService
}

func NewContactInfoHandler(service service.ContactInfoService) *ContactInfoHandler {
	return &ContactInfoHandler{
		service: service,
	}
}

type UpsertContactInfoRequest struct {
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
}

func (h *ContactInfoHandler) UpsertContactInfo(c *gin.Context) {
	var req UpsertContactInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(err.Error()))
		return
	}

	info := &entity.ContactInfo{
		Email:   req.Email,
		Phone:   req.Phone,
		Address: req.Address,
	}

	if err := h.service.UpsertContactInfo(info); err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success("Info retrieved successfully", info))
}

func (h *ContactInfoHandler) GetContactInfo(c *gin.Context) {
	info, err := h.service.GetContactInfo()
	if err != nil {
		c.JSON(http.StatusNotFound, response.Error("Contact info not found"))
		return
	}

	c.JSON(http.StatusOK, response.Success("Info retrieved successfully", info))
}
