package repositories

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/sQUARys/TestTaskAvito/app/users"
	"log"
	"time"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "myUser"
	password = "myPassword"
	dbname   = "myDb"

	connectionStringFormat = "host=%s port=%d user=%s password=%s dbname=%s sslmode=disable"

	dbCreateUserRequest = `INSERT INTO "user_table"( "id", "balance") VALUES (%d, %d)`
	//"INSERT INTO user_table VALUES (%1 , %2)"
	dbOrdersByIdRequest = "SELECT * FROM user_table WHERE id = $1"
	dbUpdateJSON        = "UPDATE user_table SET balance=%d WHERE id=%d"
)

type Repository struct {
	DbStruct *sql.DB
}

func New() *Repository {

	connectionString := fmt.Sprintf(connectionStringFormat, host, port, user, password, dbname)

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatalln(err)
	}

	repo := Repository{
		DbStruct: db,
	}

	err = db.Ping()
	if err != nil {
		log.Fatalln(err)
	}

	return &repo
}

func (repo *Repository) isUserExisting(id int) bool {
	row := repo.DbStruct.QueryRow(dbOrdersByIdRequest, id)

	var user users.User

	err := row.Scan(&user.Id, &user.Balance)

	return err != sql.ErrNoRows
}

func (repo *Repository) GetUserBalance(id int) (int, error) {
	row := repo.DbStruct.QueryRow(dbOrdersByIdRequest, id)

	var user users.User

	err := row.Scan(&user.Id, &user.Balance)

	if err != nil {
		return user.Balance, err
	}

	if err := row.Err(); err != nil {
		return user.Balance, err
	}

	return user.Balance, nil
}

func (repo *Repository) DepositMoney(user users.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	currentBalance, err := repo.GetUserBalance(user.Id)

	if err != nil {
		return err
	}

	dbInsertRequest := fmt.Sprintf(dbUpdateJSON, currentBalance+user.UpdateValue, user.Id)

	_, err = repo.DbStruct.ExecContext(
		ctx,
		dbInsertRequest,
	)

	if err != nil {
		return err
	}
	return nil
}

func (repo *Repository) WithdrawMoney(user users.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	currentBalance, err := repo.GetUserBalance(user.Id)

	if err != nil {
		return err
	}

	if currentBalance-user.UpdateValue < 0 {
		return fmt.Errorf("You can't withdraw more than you have in a balance. Please repeat operation ")
	}

	dbInsertRequest := fmt.Sprintf(dbUpdateJSON, currentBalance-user.UpdateValue, user.Id)

	_, err = repo.DbStruct.ExecContext(
		ctx,
		dbInsertRequest,
	)

	if err != nil {
		return err
	}
	return nil
}

func (repo *Repository) TransferMoney(usersTransfer users.TransferMoney) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	senderBalance, err := repo.GetUserBalance(usersTransfer.IdOfSenderUser)
	if err != nil {
		return err
	}

	recipientBalance, err := repo.GetUserBalance(usersTransfer.IdOfRecipientUser)
	if err != nil {
		return err
	}

	if senderBalance-usersTransfer.SendingAmount < 0 {
		return fmt.Errorf("You can't transfer more than you have in a balance. Please repeat operation. ")
	}

	dbInsertRequestSender := fmt.Sprintf(dbUpdateJSON, senderBalance-usersTransfer.SendingAmount, usersTransfer.IdOfSenderUser)

	_, err = repo.DbStruct.ExecContext(
		ctx,
		dbInsertRequestSender,
	)

	if err != nil {
		return err
	}

	dbInsertRequestRecipient := fmt.Sprintf(dbUpdateJSON, recipientBalance+usersTransfer.SendingAmount, usersTransfer.IdOfRecipientUser)

	_, err = repo.DbStruct.ExecContext(
		ctx,
		dbInsertRequestRecipient,
	)

	if err != nil {
		return err
	}

	return nil
}
