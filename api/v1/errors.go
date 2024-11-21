package v1

var (
	// common errors
	ErrSuccess             = newError(0, "成功")
	ErrBadRequest          = newError(400, "请求参数有误")
	ErrUnauthorized        = newError(401, "未授权的请求")
	ErrNotFound            = newError(404, "找不到对应服务")
	ErrInternalServerError = newError(500, "系统错误")
	ErrServerDealError     = newError(500, "业务处理出错")

	// more biz errors
	ErrEmailAlreadyUse = newError(1001, "The email is already in use.")
)
