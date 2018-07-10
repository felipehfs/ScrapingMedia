package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func scrapingImages(url string) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	html := string(body)
	imgPattern := regexp.MustCompile(`<img[\w\W]*?src="(https://[^"]+?)"[\w\W]*?>`)
	var urls []string
	for _, image := range imgPattern.FindAllStringSubmatch(html, -1) {
		urls = append(urls, image[1])
	}
	return urls, nil
}

func GetResources(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	search := vars["search"]
	url := fmt.Sprintf("https://pixabay.com/en/photos/?q=%s&hp=&image_type=all&order=&cat=&min_width=&min_height=", search)
	images, err := scrapingImages(url)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(images)
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/api/scrape/pixaby/{search}", GetResources).Methods("GET")
	allowedOrigins := handlers.AllowedOrigins([]string{"*"})
	allowedHeaders := handlers.AllowedHeaders([]string{"Content-Type", "Origin", "X-Requested-With", "Accept"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"})

	http.ListenAndServe(":8081", handlers.CORS(allowedOrigins, allowedMethods, allowedHeaders)(router))
}
