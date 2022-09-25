package transactionsCache

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/sQUARys/TestTaskAvito/app/users"
	"math/rand"
	"time"
)

type Cache struct {
	Client *redis.Client
}

type Transaction struct {
	UserId        int    `json:"userId"`
	OrdinalNumber int    `json:"ordinalNumber"`
	Date          string `json:"date"`
	Description   string `json:"description"`
}

func New() *Cache {
	client := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)
	return &Cache{
		Client: client,
	}
}

func (c *Cache) AddTransaction(user users.User, description string, date string) error {
	trans := Transaction{
		UserId:        user.Id,
		OrdinalNumber: rand.Int() % 10, // поправить
		Date:          date,
		Description:   description,
	}

	userJSON, err := json.Marshal(trans)
	if err != nil {
		return err
	}

	err = c.Client.Set("transaction", userJSON, 10*time.Minute).Err()
	if err != nil {
		return err
	}
	return nil
}

func (c *Cache) GetUserTransaction(key string) (Transaction, error) {
	transactionJSON, err := c.Client.Get(key).Result()
	if err != nil {
		return Transaction{}, err
	}
	var transaction Transaction
	json.Unmarshal([]byte(transactionJSON), &transaction)
	return transaction, nil
}

func (c *Cache) GetUserTransactions() ([]Transaction, error) {
	iter := c.Client.Scan(0, "", 0).Iterator()
	var transactions []Transaction

	for iter.Next() {
		transaction, err := c.GetUserTransaction(iter.Val())
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}

	if err := iter.Err(); err != nil {
		return nil, err
	}

	return transactions, nil

}
