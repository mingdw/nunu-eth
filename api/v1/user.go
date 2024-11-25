package v1

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email" example:"1234@gmail.com"`
	Password string `json:"password" binding:"required" example:"123456"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"1234@gmail.com"`
	Password string `json:"password" binding:"required" example:"123456"`
}
type LoginResponseData struct {
	AccessToken string `json:"accessToken"`
}
type LoginResponse struct {
	Response
	Data LoginResponseData
}

type UpdateProfileRequest struct {
	Nickname string `json:"nickname" example:"alan"`
	Email    string `json:"email" binding:"required,email" example:"1234@gmail.com"`
}
type GetProfileResponseData struct {
	UserId   string `json:"userId"`
	Nickname string `json:"nickname" example:"alan"`
}
type GetProfileResponse struct {
	Response
	Data GetProfileResponseData
}

type ETHConnectRequestData struct {
	Url  string `json:"url" example:"url"`
	Port string `json:"port"`
}

type AccountAddress struct {
	AccountAddress string `json:"accountAddress" example:"accountAddress"`
}

type AccountBalanceRequest struct {
	Url     string `json:"url" example:"url"`
	Address string `json:"address"`
	Block   string `json:"block"`
}

type BlockQueryRequest struct {
	Url      string `json:"url" example:"url"`
	BlockNum string `json:"blockNum"`
}

type TransactionsQueryRequest struct {
	Url       string `json:"url" example:"url"`
	BlockHash string `json:"blockHash"`
}

type ETHTransferRequest struct {
	From           string `json:"from" example:"from"`
	FromPrivateKey string `json:"fromPrivateKey"`
	To             string `json:"to"`
	Value          string `json:"value"`
}
