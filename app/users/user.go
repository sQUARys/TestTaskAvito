package users

type User struct {
	Id          int `json:"id"`
	Balance     int `json:"balance,omitempty"`
	UpdateValue int `json:"updateValue,omitempty"`
}
