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

//1. В некоторых местах забыты отладочные принты (например инициализация кэша) DONE
//2. Ошибки не оборачиваются через враппинг, придется искать по подстроке место возникновения ошибки DONE
//3. Контексты не везде прокидываются явно
//4. Используется дефолтный http клиент в конверторе без таймаутов DONE
//5. В запросах к базе данных не используются плейсхолдеры ВОПРОС
//6. Отсутствует файл миграции для базы данных, не ясна структура используемых таблиц
//7. Локи на уровне сервиса - пока ищем баланс одного пользователя, никто другой не может свой баланс получить DONE
//8. Обработка ошибок в контроллере - запись данных в ответ до записи заголовков это ошибка (как и формирование http статуса
//из количества записанных байт из Write) DONE
//9. Сделать оповещение в виде JSON о том что успешно сняли/положили деньги DONE
//10. Сделать с помощью regexp'a роутеры

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
