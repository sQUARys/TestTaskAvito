package repositories

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/sQUARys/TestTaskAvito/app/users"
	"log"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "myUser"
	password = "myPassword"
	dbname   = "myDb"

	format                 = "(%d , '%s' , %d , '%s'),"
	connectionStringFormat = "host=%s port=%d user=%s password=%s dbname=%s sslmode=disable"

	dbOrdersByIdRequest = "SELECT * FROM user_table WHERE id = $1"
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
