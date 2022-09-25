package main

import (
	"fmt"
	controller "github.com/sQUARys/TestTaskAvito/app/controllers"
	"github.com/sQUARys/TestTaskAvito/app/repositories"
	"github.com/sQUARys/TestTaskAvito/app/routers"
	"github.com/sQUARys/TestTaskAvito/app/services"
	"github.com/sQUARys/TestTaskAvito/app/transactionsCache"
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

// тесты
// mutex не везде есть done

//Задача: необходимо предоставить метод получения списка транзакций
//с комментариями откуда и зачем были начислены/списаны средства с баланса.
//Необходимо предусмотреть пагинацию и сортировку по сумме и дате.

func main() {
	repoUsers := repositories.New()

	repoTransactions := transactionsCache.New()
	//repoTransactions.AddTransaction("1", users.User{Id: 1}, "Hello 1", time.Now().Format("01-02-2006 15:04:05"))
	//repoTransactions.AddTransaction("2", users.User{Id: 2}, "Hello 2", time.Now().Format("01-02-2006 15:04:05"))

	fmt.Println(repoTransactions.GetUserTransactions())
	//time.Now().Format("01-02-2006 15:04:05")
	service := services.New(repoUsers)
	ctr := controller.New(service)

	routers := routers.New(ctr)

	routers.SetRoutes()
	err := http.ListenAndServe(":8082", routers.Router)

	if err != nil {
		log.Println("Error in main : ", err)
		return
	}
}
