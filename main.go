package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

var websites = map[string]int{}

func main() {
	websites["https://www.google.com"] = 0
	websites["https://www.facebook.com"] = 0

	fmt.Println("starting server...")
	http.HandleFunc("/GET/websites", check_websites)
	go check_status()
	run_server()

}
func check_status() {
	for {
		fmt.Println("checking now...")
		check_all_websites()
		time.Sleep(60 * time.Second)
	}
}
func run_server() {
	http.ListenAndServe(":8080", nil)
}

func check_websites(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	name, present := query["name"]
	if !present || len(name) == 0 {
		check_all_websites()
		display_all_websites(w, r)
	} else {
		check_individual_website_now(w, r, "https://"+name[0])
	}
}
func check_all_websites() {
	var wg sync.WaitGroup
	wg.Add(len(websites))

	for url, _ := range websites {
		go check_individual_website(url, &wg)

	}
	wg.Wait()

}

func check_individual_website(url string, wg *sync.WaitGroup) {
	defer wg.Done()
	resp, _ := http.Get(url)
	websites[url] = resp.StatusCode

}
func check_individual_website_now(w http.ResponseWriter, r *http.Request, url string) {
	resp, _ := http.Get(url)
	websites[url] = resp.StatusCode

	fmt.Fprint(w, url, " : ", websites[url], "\n")
}

func display_all_websites(w http.ResponseWriter, r *http.Request) {
	for url, _ := range websites {
		fmt.Fprint(w, url, " : ", websites[url], "\n")
	}
}
