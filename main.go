package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	response, err := http.Get("http://srv.msk01.gigacorp.local")
	if err != nil {
		fmt.Println(err)
	}
	if response.StatusCode == 200 {
		fmt.Printf("Status Code: %d\r\n", response.StatusCode)
		body, err := io.ReadAll(response.Body)
		defer response.Body.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(body))
	} else {
		fmt.Println("Unable to fetch server statistic")
	}
}
