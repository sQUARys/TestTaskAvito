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

	format                 = "(%d , '%d'),"
	connectionStringFormat = "host=%s port=%d user=%s password=%s dbname=%s sslmode=disable"

	dbOrdersByIdRequest = "SELECT * FROM user_table WHERE id = $1"
	dbInsertJSON        = `INSERT INTO user_table( "id", "balance" ) VALUES `
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

func (repo *Repository) GetUserBalance(id int) (int, error) {
	row := repo.DbStruct.QueryRow(dbOrdersByIdRequest, id)

	var user users.User

	if err := row.Scan(&user.Id, &user.Balance); err != nil {
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
