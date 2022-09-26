package repositories

import (
	"github.com/sQUARys/TestTaskAvito/app/users"
	"testing"
)

func TestRepository_CreateUser(t *testing.T) {
	type args struct {
		id int
	}

	tests := []struct {
		name    string
		repo    *Repository
		args    args
		wantErr bool
	}{
		{name: "Test1", repo: New(), args: args{id: 1}, wantErr: false},
		{name: "Test2", repo: New(), args: args{id: -100}, wantErr: true},
		{name: "Test3", repo: New(), args: args{id: 0}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &Repository{
				DbStruct: tt.repo.DbStruct,
			}
			if err := repo.CreateUser(tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("CreateUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRepository_WithdrawMoney(t *testing.T) {
	user1 := users.User{
		Id:          1,
		UpdateValue: 10,
	}
	user2 := users.User{
		Id:          2,
		UpdateValue: -500, // в данном случае мы не крашимся потому что проверка на отрицательные значения идет выше в уровне контроллера, а дб просто игнорирует этот запрос даже если он дойдет до данного уровня
	}
	user3 := users.User{
		Id:          3,
		UpdateValue: 1500,
	}
	repo := New()

	tests := []struct {
		name    string
		fields  *Repository
		args    users.User
		wantErr bool
	}{
		{name: "Test1", fields: repo, args: user1, wantErr: true},
		{name: "Test2", fields: repo, args: user2, wantErr: false},
		{name: "Test3", fields: repo, args: user3, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &Repository{
				DbStruct: tt.fields.DbStruct,
			}
			if err := repo.WithdrawMoney(tt.args); (err != nil) != tt.wantErr {
				t.Errorf("WithdrawMoney() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRepository_DepositMoney(t *testing.T) {
	user1 := users.User{
		Id: 1,
	}
	user2 := users.User{
		Id:          2,
		UpdateValue: -500,
	}
	user3 := users.User{
		Id:          3,
		UpdateValue: 1500,
	}
	repo := New()

	tests := []struct {
		name    string
		fields  *Repository
		args    users.User
		wantErr bool
	}{
		{name: "Test1", fields: repo, args: user1, wantErr: false},
		{name: "Test2", fields: repo, args: user2, wantErr: false},
		{name: "Test3", fields: repo, args: user3, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &Repository{
				DbStruct: tt.fields.DbStruct,
			}
			if err := repo.DepositMoney(tt.args); (err != nil) != tt.wantErr {
				t.Errorf("DepositMoney() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRepository_GetUserBalance(t *testing.T) {
	user1 := users.User{
		Id: 1,
	}
	user2 := users.User{ // в данном тесте в бд уже есть юзер с id = 3, и его баланс = 3000
		Id: 3,
	}
	repo := New()

	tests := []struct {
		name    string
		fields  *Repository
		args    users.User
		want    float64
		wantErr bool
	}{
		{name: "Test1", fields: repo, args: user1, want: 0, wantErr: false},
		{name: "Test2", fields: repo, args: user2, want: 3000, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &Repository{
				DbStruct: tt.fields.DbStruct,
			}
			got, err := repo.GetUserBalance(tt.args.Id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserBalance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetUserBalance() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRepository_TransferMoney(t *testing.T) {
	//данные в бд
	//id, balance
	//1,0
	//2,0
	//3,3000

	user1 := users.TransferMoney{
		IdOfSenderUser:    3,
		IdOfRecipientUser: 1,
		SendingAmount:     1500,
	}
	user2 := users.TransferMoney{ // в данном случае ждем ошибки потому что мы уже перевели 1500 и баланс 3 юзера стал 1500, а переводить больше чем у тебя есть нельзя
		IdOfSenderUser:    3,
		IdOfRecipientUser: 1,
		SendingAmount:     1600,
	}
	repo := New()

	tests := []struct {
		name    string
		fields  *Repository
		args    users.TransferMoney
		wantErr bool
	}{
		{name: "Test1", fields: repo, args: user1, wantErr: false},
		{name: "Test2", fields: repo, args: user2, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &Repository{
				DbStruct: tt.fields.DbStruct,
			}
			if err := repo.TransferMoney(tt.args); (err != nil) != tt.wantErr {
				t.Errorf("TransferMoney() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
