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
	dbInsertJSON        = `INSERT INTO "order_table"( "id", "balance") VALUES `
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

	dbInsertRequest := fmt.Sprintf(format, user.Id, user.Balance+user.Deposit)

	_, err := repo.DbStruct.ExecContext(
		ctx,
		dbInsertRequest,
	)

	if err != nil {
		return err
	}
	return nil
}
