Программа читает из stdin строки, содержащие URL. На каждый URL нужно отправить HTTP-запрос методом GET
и посчитать кол-во вхождений строки "Go" в теле ответа. В конце работы приложение выводит на экран общее количество найденных строк "Go" во всех переданных URL, например:

```cassandraql
$ echo -e 'https://golang.org\nhttps://golang.org' | go run 1.go
Count for https://golang.org: 9
Count for https://golang.org: 9
Total: 18
```

Каждый URL должен начать обрабатываться сразу после считывания и параллельно со считыванием следующего. URL должны обрабатываться параллельно, но не более k=5 одновременно. Обработчики URL не должны порождать лишних горутин, т.е. если k=5, а обрабатываемых URL-ов всего 2, не должно создаваться 5 горутин.

Нужно обойтись без глобальных переменных и использовать только стандартную библиотеку.

Код необходимо залить в публичный репозиторий github или на gist.

тестовые команды для запуска
```
echo -e 'https://golang.org\nhttps://golang.org\nhttps://golang.org\nhttps://golang.org\nhttps://golang.org\nhttps://golang.org' | go run 1.go
```

с одним недоступным URL
```
echo -e 'https://golang.org\nhttps://golang.org\nhttps://golang.org\nhttps://golang.org\nhttps://golang.org\nhttps://golang.org\nhttps://golang.org\nhttps://golang.org\nhttps://golang.org\nhttps://golang.org.ru' | go run 1.go
```

все URL-ы несуществующие
```
echo -e 'https://golang9.org\nhttps://golang8.org\nhttps://golang7.org\nhttps://golang6.org\nhttps://golang5.org\nhttps://golang4.org\nhttps://golang3.org\nhttps://golang2.org\nhttps://golang1.org\nhttps://golang.org.ru' | go run 1.go
```