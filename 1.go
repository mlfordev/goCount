package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	ch := make(chan int, 5)
	total := 0

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		go func() {
			counter := 0
			URL := scanner.Text()

			res, err := http.Get(URL)
			if err != nil {
				fmt.Println(err)
				return
			}

			defer res.Body.Close()

			date, err := ioutil.ReadAll(res.Body)
			if err != nil {
				fmt.Println(err)
				return
			}

			counter = bytes.Count(date, []byte("Go"))

			fmt.Printf("Count for %s: %d\n", URL, counter)

			ch <- counter
		}()
		total += <-ch
	}

	fmt.Printf("Total: %d\n", total)
}
