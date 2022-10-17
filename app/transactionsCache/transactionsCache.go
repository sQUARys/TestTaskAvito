package transactionsCache

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/sQUARys/TestTaskAvito/app/users"
	"time"
)

type Cache struct {
	Client *redis.Client
}

type FormatOfTransaction struct {
	UserId      int
	Description string
	Date        string
}

func New() (*Cache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})

	_, err := client.Ping().Result()
	if err != nil {
		err = fmt.Errorf("%w", err)
	}

	return &Cache{
		Client: client,
	}, err
}

func (c *Cache) AddTransaction(userId int, description string, date string) error {
	transaction := FormatOfTransaction{
		UserId:      userId,
		Description: description,
		Date:        date,
	}

	userJSON, err := json.Marshal(transaction)
	if err != nil {
		return err
	}

	err = c.Client.Set("transaction"+description, userJSON, time.Minute).Err()
	if err != nil {
		return err
	}
	return nil
}

func (c *Cache) GetUserTransaction(key string) (users.Transaction, error) {
	transactionJSON, err := c.Client.Get(key).Result()
	if err != nil {
		return users.Transaction{}, err
	}
	var transaction users.Transaction
	err = json.Unmarshal([]byte(transactionJSON), &transaction)
	return transaction, err
}

func (c *Cache) GetUserTransactions(id int) ([]users.Transaction, error) {
	iter := c.Client.Scan(0, "", 0).Iterator()
	var transactions []users.Transaction

	for iter.Next() {
		transaction, err := c.GetUserTransaction(iter.Val())
		if err != nil {
			return nil, err
		}
		if transaction.UserId == id {
			transactions = append(transactions, transaction)
		}
	}

	if err := iter.Err(); err != nil {
		return nil, err
	}

	return transactions, nil

}
