package main

import (
	controller "github.com/sQUARys/TestTaskAvito/app/controllers"
	"github.com/sQUARys/TestTaskAvito/app/repositories"
	routers "github.com/sQUARys/TestTaskAvito/app/routers"
	"github.com/sQUARys/TestTaskAvito/app/services"
	"net/http"
)

//Необходимо реализовать микросервис для работы с балансом пользователей
//(зачисление средств, списание средств,
//перевод средств от пользователя к пользователю,
//а также метод получения баланса пользователя).
//	Сервис должен предоставлять HTTP API и принимать/отдавать
//	запросы/ответы в формате JSON.

func main() {
	repo := repositories.New()
	service := services.New(repo)
	ctr := controller.New(service)
	routers := routers.New(ctr)

	routers.SetRoutes()
	http.ListenAndServe(":8080", routers.Router)
}
