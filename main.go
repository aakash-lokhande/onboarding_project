package main

import (
	"fmt"
	"net/http"
	"sync"
)

var websites = map[string]int{}

func main() {
	websites["https://www.google.com"] = 0
	websites["https://www.facebook.com"] = 0

	fmt.Println("starting server...")
	http.HandleFunc("/", handler)
	http.HandleFunc("/GET/websites", check_websites)
	go add()
	run_server()
	fmt.Println("nedt...")

}
func add() {
	websites["https://www.youtube.com"] = 0
}
func run_server() {
	http.ListenAndServe(":8080", nil)
}
func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}
func check_websites(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	name, present := query["name"]
	if !present || len(name) == 0 {
		check_all_websites(w, r)
	} else {
		check_individual_website_now(w, r, "https://"+name[0])
	}
}
func check_all_websites(w http.ResponseWriter, r *http.Request) {
	var wg sync.WaitGroup
	wg.Add(len(websites))

	for url, _ := range websites {
		go check_individual_website(w, r, url, &wg)

	}
	wg.Wait()
	for url, _ := range websites {
		fmt.Fprint(w, url, " : ", websites[url], "\n")
	}

}

func check_individual_website(w http.ResponseWriter, r *http.Request, url string, wg *sync.WaitGroup) {
	defer wg.Done()
	resp, _ := http.Get(url)
	websites[url] = resp.StatusCode

}
func check_individual_website_now(w http.ResponseWriter, r *http.Request, url string) {
	resp, _ := http.Get(url)
	websites[url] = resp.StatusCode

	fmt.Fprint(w, url, " : ", websites[url], "\n")
}
