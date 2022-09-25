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
//перевод средств от пользователя к пользователю DONE
// метод получения баланса пользователя DONE
// в любой валюте DONE
//	Сервис должен предоставлять HTTP API и принимать/отдавать
//	запросы/ответы в формате JSON.

//Задача: добавить к методу получения баланса доп. параметр.
//Пример: ?currency=USD. Если этот параметр присутствует, DONE
//то мы должны конвертировать баланс пользователя с рубля на указанную валюту.
//Данные по текущему курсу валют можно взять отсюда https://exchangeratesapi.io/ или из любого другого открытого источника.
//Примечание: напоминаем, что базовая валюта которая хранится на балансе у нас всегда рубль. В рамках этой задачи конвертация всегда происходит с базовой валюты.

//все отправлять/возвращать в Json
// тесты
// mutex не везде есть

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
