package dtos

type InitTransactionParam struct {
	PaymentMethod string `json:"payment_method"`
	UserID        int64  `json:"-"`
}

type InitTransactionResponse struct {
	VANumber string `json:"va_number"`
}

type AcceptTransactionParam struct {
	PaymentMethod string `json:"payment_method"`
	VaNumber      string `json:"va_number"`
	UserID        int64  `json:"-"`
}
