# Тестовое задание на позицию стажера-бекендера

## Микросервис для работы с балансом пользователей.

**Задача:**

Необходимо реализовать микросервис для работы с балансом пользователей (зачисление средств, списание средств, перевод средств от пользователя к пользователю, а также метод получения баланса пользователя). Сервис должен предоставлять HTTP API и принимать/отдавать запросы/ответы в формате JSON. 

**Стек используемых в сервисе технологий**
1.Golang
2.PostgreSQL(для хранения данных о балансах, id и сумм снятия пользователей)
3.Redis(для хранения и отображения всех транзакций)

Реализовано:

1.Метод начисления средств на баланс. Принимает id пользователя и сколько средств зачислить.

2. Метод списания средств с баланса. Принимает id пользователя и сколько средств списать. 

3. Метод перевода средств от пользователя к пользователю. Принимает id пользователя с которого нужно списать средства, id пользователя которому должны зачислить средства, а также сумму.

4.Метод получения текущего баланса пользователя. Принимает id пользователя. Баланс всегда в рублях. (по умолчанию сервис не содержит в себе никаких данных о балансах (пустая табличка в БД). Данные о балансе появляются при первом зачислении денег)

5. В методе получения баланса сделан доп. параметр. Пример: ?currency=USD. 
Если этот параметр присутствует, то мы должны конвертировать баланс пользователя с рубля на указанную валюту. Данные по текущему курсу валют я беру из https://exchangeratesapi.io/.(базовая валюта хранится на балансе в рублях)

6. Метод получения списка транзакция каждого пользователя.

**Принцип работы**
