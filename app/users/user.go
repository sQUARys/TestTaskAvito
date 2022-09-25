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

type Transaction struct {
	UserId        int    `json:"userId"`
	OrdinalNumber int    `json:"ordinalNumber"`
	Date          string `json:"date"`
	Description   string `json:"description"`
}
