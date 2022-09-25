package users

type User struct {
	Id          int     `json:"id"`
	Balance     float64 `json:"balance,omitempty"`
	UpdateValue float64 `json:"updateValue,omitempty"`
}

type TransferMoney struct {
	IdOfSenderUser    int     `json:"id-sender"`
	IdOfRecipientUser int     `json:"id-recipient"`
	SendingAmount     float64 `json:"sending-amount"`
}
