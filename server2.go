package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("starting server 2")
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8081", nil)
}
func handler(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get("https://github.com/deasfse")
	fmt.Fprint(w, resp.StatusCode, err)
}
