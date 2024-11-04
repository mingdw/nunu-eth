package handler

import (
	"nunu-eth/internal/service"

	"github.com/gin-gonic/gin"
)

type CommonHandler struct {
	*Handler
	commonService service.CommonService
}

func NewCommonHandler(
	handler *Handler,
	commonService service.CommonService,
) *CommonHandler {
	return &CommonHandler{
		Handler:       handler,
		commonService: commonService,
	}
}

// GetProfile godoc
// @Summary 获取用户信息
// @Schemes
// @Description
// @Tags 用户模块
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} response.Response
// @Router /user [get]
func (h *CommonHandler) GetCommon(ctx *gin.Context) {

}

func (h *CommonHandler) Test(ctx *gin.Context) {

}
