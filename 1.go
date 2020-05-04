package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type result struct {
	URL     string
	counter int
	error   string
}

func main() {
	client := &http.Client{}
	ch := make(chan int, 5)
	total := 0

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		go func() {
			r := result{URL: scanner.Text()}
			response, err := client.Get(r.URL)

			if err != nil {
				r.error = "Error: " + err.Error()
			} else {
				defer response.Body.Close()
				date, err := ioutil.ReadAll(response.Body)
				if err != nil {
					r.error = "Error: " + err.Error()
				} else {
					r.counter = bytes.Count(date, []byte("Go"))
				}
			}
			fmt.Printf("Count for %s: %d\t%s\n", r.URL, r.counter, r.error)
			ch <- r.counter
		}()
		total += <-ch
	}
	fmt.Printf("Total: %d\n", total)
}
