package main

import (
	controller "github.com/sQUARys/TestTaskAvito/app/controllers"
	"github.com/sQUARys/TestTaskAvito/app/repositories"
	"github.com/sQUARys/TestTaskAvito/app/routers"
	"github.com/sQUARys/TestTaskAvito/app/services"
	"log"
	"net/http"
)

//Необходимо реализовать микросервис для работы с балансом пользователей
//зачисление средств, DONE
//списание средств, DONE
// сделать mutex чтобы избежать гонки DONE
//перевод средств от пользователя к пользователю DONE
// метод получения баланса пользователя DONE
//написать тесты
// в любой валюте
//	Сервис должен предоставлять HTTP API и принимать/отдавать
//	запросы/ответы в формате JSON.

func main() {
	repo := repositories.New()
	service := services.New(repo)
	ctr := controller.New(service)
	routers := routers.New(ctr)

	routers.SetRoutes()
	err := http.ListenAndServe(":8081", routers.Router)

	if err != nil {
		log.Println("Error in main : ", err)
		return
	}
}
