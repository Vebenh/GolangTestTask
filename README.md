# Тестовое задание на позицию Golang разработчика

## Задание
A) Реализовать приложение на golang 
#### Методы
1. localhost:8080/update \
   Импорт / обновление необходимых данных из https://www.treasury.gov/ofac/downloads/sdn.xml в локальную базу 
    PostgreSQL 14. в базу должны попадать записи с sdnType=Individual \
    Результат:
    - success:
   {"result": true, "info": "", "code": 200}
    - fail:
    {"result": false, "info": "service unavailable", "code": 503}


2. localhost:8080/state \
   Получение текущего состояния данных
   - нет данных: \
   {"result": false, "info": "empty"}
   - в процессе обновления:\
   {"result": false, "info": "updating"}
   - данные готовы к использованию:\
   {"result": true, "info": "ok"}


3.  localhost:8080/get_names?name={SOME_VALUE}&type={strong|weak} \
    Получение списка всех возможных имён человека из локальной базы данных с указанием основного uid в виде JSON.
    Если параметр type не указан / указан ошибочно, то выдаём список состоящий из всех типов. 
    Параметр type независим от регистра. strong - это точное совпадение имени и фамилии, 
    weak - должно найти любое совпадение в имени либо фамилии
    - Запрос: localhost:8080/get_names?name=MUZONZINI&type=strong \
    Результат: \
    [{uid:7535, first_name:"Elisha", last_name:"Muzonzini"}] 
    - Запрос: localhost:8080/get_names?name=Elisha Muzonzini \
    Результат: \
    [{uid:7535, first_name:"Elisha", last_name:"Muzonzini"}]
    - Запрос: localhost:8080/get_names?name=Mohammed Musa&type=weak \
    Результат: \
    [{uid:15582, first_name:"Musa", last_name:"Kalim"}, {uid:15582, first_name:"Barich", last_name:"Musa Kalim"}, {uid:15582, first_name:"Mohammed Musa", last_name:"Kalim"}, {uid:15582, first_name:"Musa Khalim", last_name:"Alizari"}, {uid:15582, first_name:"Qualem", last_name:"Musa"}, {uid:15582, first_name:"Qualim", last_name:"Musa"}, {uid:15582, first_name:"Khaleem", last_name:"Musa"}, {uid:15582, first_name:"Kaleem", last_name:"Musa"}]


B) Написать инструкции Docker Compose для разворачивания реализованного приложения на порту 8080 с использованием Postgresql14.

C) Описать алгоритм для более эффективного обновления данных при повторном вызове метода localhost:8080/update ( можно реализовать, но не обязательно )

##  Эффективное обновление

#### Сделано:
- Чтение и парсинг XML в потоке, чтобы не хранить его целиком в памяти;
- Параллельное чтение из xml и запись в базу в горутинах;
- Проверка существующих записей по uid. Для этого в таблице sdn_entries добавлен индекс для быстрого поиска по iud.

#### Не сделано
- Проверка необходимости обновления. В начале xml есть тег `<publshInformation>`.
  Чтобы не читать весь файл целиком, можно в потоке считать этот тег, и если с прошлого апдейта не было изменений,
  сразу возвращать соответвующий ответ и прекращать чтение;
- Использование временных таблиц. Можно по мере чтения xml записывать данные во временные таблицы, а после сделать вставку из временных таблиц в основные таблицы.


## Деплой
1. После клонирования проекта создаем файлы `config.yaml` в дирректории `/config` и `.env` в корне проекта. Соответствующие sample файлы лежат в нужных местах. Пароль должен содержать цифры и буквы в верхнем и нижнем регистре, иначе postgres контейнер не запустится.
2. Выполняем из корня проекта команду
```bash
make build
```
3. Готово! Теперь можно выполнить запрос 
```bash
curl http://localhost:8080/state
```
Должен прийти ответ:
```bash
{"result":false,"info":"empty"}%                                                                                                                                                                                test/GolangTestTask [main●] » 
```
## Использованные библиотеки

- [viper](https://github.com/spf13/viper) - библиотека для работы с конфигами
- [chi](https://github.com/go-chi/chi) - маршрутизатор для создания HTTP-сервисов
- [gorm](https://github.com/go-gorm/gorm) - ORM для Golang