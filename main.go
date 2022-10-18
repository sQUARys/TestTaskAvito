package main

import (
	controller "github.com/sQUARys/TestTaskAvito/app/controllers"
	"github.com/sQUARys/TestTaskAvito/app/providers"
	"github.com/sQUARys/TestTaskAvito/app/repositories"
	"github.com/sQUARys/TestTaskAvito/app/routers"
	"github.com/sQUARys/TestTaskAvito/app/services"
	"github.com/sQUARys/TestTaskAvito/app/transactionsCache"
	"log"
	"net/http"
	"time"
)

//3. Контексты не везде прокидываются явно
//5. В запросах к базе данных не используются плейсхолдеры ВОПРОС
//6. Отсутствует файл миграции для базы данных, не ясна структура используемых таблиц

func main() {
	currencyProvider := providers.New()
	repoUsers := repositories.New()
	repoTransactions, err := transactionsCache.New()
	if err != nil {
		log.Println(err)
	}
	service := services.New(repoUsers, repoTransactions, currencyProvider)

	ctr := controller.New(service)

	routers := routers.New(ctr)

	routers.SetRoutes()

	server := http.Server{
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 3 * time.Second,
		Addr:         ":8081",
		Handler:      routers.Router,
	}

	err = server.ListenAndServe()

	if err != nil {
		log.Println("Error in main : ", err)
		return
	}
}
