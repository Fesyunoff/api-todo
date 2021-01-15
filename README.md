## Необходимо спроектировать API сервис TODO.
 
Функционал:
- Несколько аккаунтов с правами "Админ" и "Пользователь".
- Создание/Просмотр/Изменение/Удаление TODO.
- Админ может управлять TODO всех и создавать, удалять юзеров.
- SQL база данных(мы используем MySQL, но для данного задания можно воспользоваться SQLite).
- После создания TODO, должна быть возможность уведомить внешний сервис об этом, например сообщение в Telegram или WebHook(реализовать интерфейс для этого, сама реализация отсылки не требуется).

Восстановления пароля и работа с email, в данном задании, не требуется.


Что хочется увидеть:
- Project Layout.
- Какие используются внешние зависимости.
- Как организованно внедрение зависимостей(Dependency Injection).

Если для выполнения работы используются какие-то статьи\репозитории, просьба приложить на них ссылки.  

Build Service:

$ go build -o ./bin/service ./cmd/service/

Start Service:

$ ./bin/service

Test Service:

$ chmod +x ./test_script
$ ./test_script

List of sources used:

https://godoc.org/github.com/mattn/go-sqlite3
https://www.jsonrpc.org/specification
https://github.com/swipe-io/swipe
https://habr.com/ru/post/430300/
https://habr.com/ru/company/skillbox/blog/446454/
