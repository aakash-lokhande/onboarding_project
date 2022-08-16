package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

type site_info struct {
	url        string
	statusCode int
	err        error
	runStatus  string
}

var websites = map[string]site_info{}

func main() {
	google := site_info{"https://www.google.com", 0, nil, "DOWN"}
	facebook := site_info{"https://www.facebook.com", 0, nil, "DOWN"}
	localhost := site_info{"http://localhost:8081/", 0, nil, "DOWN"}
	websites["www.google.com"] = google
	websites["www.facebook.com"] = facebook
	websites["localhost:8081"] = localhost

	fmt.Println("starting server...")
	http.HandleFunc("/GET/websites", check_websites)
	http.HandleFunc("/POST/websites", post_websites)
	//go check_status()
	run_server()

}
func check_status() {
	for {
		fmt.Println("checking now...")
		check_all_websites()
		time.Sleep(5 * time.Second)
	}
}
func run_server() {
	http.ListenAndServe(":8080", nil)
}

func post_websites(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println(err)
	}
	x := r.Form.Get("websites")
	fmt.Printf("%T", x)
	fmt.Println(" ", x, "\n")

}
func check_websites(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	name, present := query["name"]
	if !present || len(name) == 0 {
		check_all_websites()
		display_all_websites(w, r)
	} else {
		fmt.Fprint(w, name[0], " : ", websites[name[0]].runStatus, "\n")
	}
}
func check_all_websites() {
	var wg sync.WaitGroup
	wg.Add(len(websites))

	for name, _ := range websites {
		go conc_site_check(name, &wg)
	}
	wg.Wait()

}
func conc_site_check(name string, wg *sync.WaitGroup) {
	defer wg.Done()
	site := websites[name]
	site.update_running_status()
	websites[name] = site

}
func (site *site_info) update_running_status() {

	resp, err := http.Get(site.url)
	if err != nil {
		site.statusCode = 404
		site.err = err
		site.runStatus = "DOWN"
	} else if resp.StatusCode == 200 {
		site.statusCode = 200
		site.err = nil
		site.runStatus = "UP"
	} else {
		site.statusCode = resp.StatusCode
		site.err = nil
		site.runStatus = "DOWN"
	}
}

func display_all_websites(w http.ResponseWriter, r *http.Request) {
	for name, site := range websites {
		if site.statusCode == 200 {
			fmt.Fprint(w, name, " : ", site.runStatus, "\n")
		} else {
			fmt.Fprint(w, name, " : ", site.runStatus, "\n", "Error : ", site.err)
		}

	}
}
