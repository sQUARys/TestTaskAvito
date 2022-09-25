package users

type User struct {
	Id          int `json:"id"`
	Balance     int `json:"balance,omitempty"`
	UpdateValue int `json:"updateValue,omitempty"`
}

type TransferMoney struct {
	IdOfSenderUser    int `json:"id-sender"`
	IdOfRecipientUser int `json:"id-recipient"`
	SendingAmount     int `json:"sending-amount"`
}
