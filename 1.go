package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type Target struct {
	URL     string
	counter int
	error   string
}

func countWords(t Target, client *http.Client) Target {
	response, err := client.Get(t.URL)

	if err != nil {
		t.error = "Error: " + err.Error()
	} else {
		defer response.Body.Close()
		date, err := ioutil.ReadAll(response.Body)
		if err != nil {
			t.error = "Error: " + err.Error()
		} else {
			t.counter = bytes.Count(date, []byte("Go"))
		}
	}
	return t
}

func main() {
	client := &http.Client{}
	targetCh := make(chan Target)
	resultCh := make(chan Target)
	goCounter := 0
	goMax := 5
	goDone := make(chan bool)
	done := make(chan bool)
	total := 0

	// Горутина выводит результаты подсчета слов
	go func() {
		// Принимаем из канала структуры с результатами подсчета слов с заданных URL-ов
		for t := range resultCh {
			fmt.Printf("Count for %s: %d\t%s\n", t.URL, t.counter, t.error)
			total += t.counter
		}
		fmt.Printf("Total: %d\n", total)

		// Отправляем сообщение о завершении работы горутины, выводящей результаты подсчета слов
		done <- true
	}()

	// Построчно считываем URL
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {

		// Позволяем запуститься 5-ти горутинам-иполнителям
		if goCounter >= goMax {
			// блокирум запуск последующих горутин-исполнителей
			// до получения сообщения о завершении работы одной из ранее запущеных горутин-исполнителей
			// и ждем сообщения о завершении работы горутин-иполнителей
			<-goDone
		}

		// запускаем горутины-исполнители
		go func() {
			// принимаем из канала структуру с URL-ом для подсчета слов
			t := <-targetCh
			// Делаем http-запрос, считаем количество вхождений слова, записываем полученные данные в структуру
			// отправляем структуру через канал в горутину для вывода результатов подсчета вхождений слова
			resultCh <- countWords(t, client)
			// отправляем сообщение в канал о завершении работы горутины-иполнителя
			goDone <- true
		}()

		// Записываем в структуру URL для подсчета слов
		t := Target{URL: scanner.Text()}
		// Отправляем структуру через канал в горутину-исполнитель
		targetCh <- t
		goCounter++
	}

	if goCounter > goMax {
		goCounter = goMax
	}

	// Принимаем недополученные сообщения о завершении работы 1-5 горутин-исполнителей
	for i := 0; i < goCounter; i++ {
		<-goDone
	}
	// Закрываем канал передачи структур с результатами работы горутин-исполнителей,
	// тем самым завершаем цикл в горутине, выводящей результаты подсчета слов
	close(resultCh)

	// Получаем сообщение о завершении работы горутины, выводящей результаты подсчета слов
	// и позволаем завершиться приложению
	<-done
}
