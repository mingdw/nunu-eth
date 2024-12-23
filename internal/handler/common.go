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
	accountAddress := ctx.Query("accountAddress")
	if accountAddress == "" {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	req.AccountAddress = accountAddress
	accountInfo, err := h.commonService.AccountFormatInfo(ctx, &req)
	if err != nil {
		v1.HandleError(ctx, http.StatusUnauthorized, v1.ErrUnauthorized, nil)
		return
	}
	v1.HandleSuccess(ctx, accountInfo)
}

func (h *CommonHandler) AccountBalance(ctx *gin.Context) {
	var req v1.AccountBalanceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	accountInfo, err := h.commonService.AccountBalance(ctx, &req)
	if err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	v1.HandleSuccess(ctx, accountInfo)
}

func (h *CommonHandler) BlockQuery(ctx *gin.Context) {
	var req v1.BlockQueryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	accountInfo, err := h.commonService.BlockQuery(ctx, &req)
	if err != nil {
		v1.HandleError(ctx, http.StatusOK, err, nil)
		return
	}
	v1.HandleSuccess(ctx, accountInfo)
}

func (h *CommonHandler) TransactionQuery(ctx *gin.Context) {
	var req v1.TransactionsQueryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	accountInfo, err := h.commonService.TransactionQuery(ctx, &req)
	if err != nil {
		v1.HandleError(ctx, http.StatusOK, err, nil)
		return
	}
	v1.HandleSuccess(ctx, accountInfo)
}

func (h *CommonHandler) CreateAccount(ctx *gin.Context) {
	accountInfo, err := h.commonService.CreateAccount(ctx)
	if err != nil {
		v1.HandleError(ctx, http.StatusOK, err, nil)
		return
	}
	v1.HandleSuccess(ctx, accountInfo)
}

func (h *CommonHandler) TxQuery(ctx *gin.Context) {
	hash := ctx.Query("txHash")
	accountInfo, err := h.commonService.TxQuery(ctx, hash)
	if err != nil {
		v1.HandleError(ctx, http.StatusOK, err, nil)
		return
	}
	v1.HandleSuccess(ctx, accountInfo)
}

func (h *CommonHandler) ETHTransfer(ctx *gin.Context) {
	var req v1.ETHTransferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}
	accountInfo, err := h.commonService.ETHTransfer(ctx, &req)
	if err != nil {
		v1.HandleError(ctx, http.StatusOK, err, nil)
		return
	}
	v1.HandleSuccess(ctx, accountInfo)
}
