package handler

import (
	"net/http"
	v1 "nunu-eth/api/v1"
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

func (h *CommonHandler) TestConnectClient(ctx *gin.Context) {
	var req v1.ETHConnectRequestData
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	status, err := h.commonService.ConnectTest(ctx, &req)
	if err != nil {
		v1.HandleError(ctx, http.StatusUnauthorized, v1.ErrUnauthorized, nil)
		return
	}
	if status != 0 {
		v1.HandleError(ctx, http.StatusUnauthorized, v1.ErrUnauthorized, nil)
		return
	}
	v1.HandleSuccess(ctx, status)
}

func (h *CommonHandler) AccountFormt(ctx *gin.Context) {
	var req v1.AccountAddress
	if err := ctx.ShouldBind(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	accountInfo, err := h.commonService.AccountFormatInfo(ctx, &req)
	if err != nil {
		v1.HandleError(ctx, http.StatusUnauthorized, v1.ErrUnauthorized, nil)
		return
	}
	v1.HandleSuccess(ctx, accountInfo)
}
