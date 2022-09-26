package main

import (
	controller "github.com/sQUARys/TestTaskAvito/app/controllers"
	"github.com/sQUARys/TestTaskAvito/app/repositories"
	"github.com/sQUARys/TestTaskAvito/app/routers"
	"github.com/sQUARys/TestTaskAvito/app/services"
	"github.com/sQUARys/TestTaskAvito/app/transactionsCache"
	"log"
	"net/http"
)

func main() {
	repoUsers := repositories.New()
	repoTransactions := transactionsCache.New()

	service := services.New(repoUsers, repoTransactions)

	ctr := controller.New(service)

	routers := routers.New(ctr)

	routers.SetRoutes()
	err := http.ListenAndServe(":8083", routers.Router)

	if err != nil {
		log.Println("Error in main : ", err)
		return
	}
}
