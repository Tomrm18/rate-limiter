package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	for {
		res, err := http.Get("http://localhost:8080")
		if err != nil {
			fmt.Println(err)
		}
		if res != nil {
			fmt.Println("request status = ", res.Status)
			fmt.Println("res body = ", res.Body)
			res.Body.Close()
		}

		time.Sleep(time.Millisecond * 500)
	}
}
