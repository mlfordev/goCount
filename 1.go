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
	goDone := make(chan bool)
	done := make(chan bool)
	total := 0

	go func() {
		for t := range resultCh {
			fmt.Printf("Count for %s: %d\t%s\n", t.URL, t.counter, t.error)
			total += t.counter
		}
		fmt.Printf("Total: %d\n", total)
		done <- true
	}()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		go func() {
			t := <-targetCh
			resultCh <- countWords(t, client)
			goDone <- true
		}()

		t := Target{URL: scanner.Text()}
		targetCh <- t
		goCounter++
	}

	for i := 0; i < goCounter; i++ {
		<-goDone
	}
	close(resultCh)

	<-done
}
