package main

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"time"
	"encoding/json"
)

type URL struct {
	ID string `json:"id"`
	OriginalURL string `json:"original_url"`
	Shorturl string `json:"short_url"`
	CreationDate time.Time `json:"creation_date"`
}

var urlDB = make(map[string]URL)

func generateShortURL(OriginalURL string) string {
	hasher:= md5.New()
	hasher.Write([]byte(OriginalURL))
	fmt.Println("hasker:", hasher)
	data:= hasher.Sum(nil);
	hash := hex.EncodeToString(data)
	fmt.Println("EncodedToString:",hash)
	fmt.Println("hasher data:",data)
	fmt.Println("Final string:",hash[:8])
	return hash[:8];
}

func createURL(originalURL string) string{
	shortURL := generateShortURL(originalURL)
	id := shortURL
	urlDB[id] = URL{
		ID: id,
		OriginalURL: originalURL,
		Shorturl: shortURL,
		CreationDate: time.Now(),
	}
	return shortURL
}
	
func getURL(id string) (URL, error){
	url,ok:=urlDB[id]
	if !ok {
		return URL{},errors.New("URL not found")
	}
	return url,nil
}

func RootPageURL(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w , "HELOO Workld")
}

func ShortURLHandler(w http.ResponseWriter,r *http.Request){
	var data struct{
		URL  string `json:"url"`
	}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Invalid request body",http.StatusBadRequest)
		return
	}

	shortURL_ := createURL(data.URL);
	//fmt.Fprintf(w, shortURL)
	response := struct{
		ShortURL string `json:"short_url"`
	}{ShortURL: shortURL_}

	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func redirectURLHandler(w http.ResponseWriter,r *http.Request){
	id := r.URL.Path[len("/redirect/"):]
	url,err := getURL(id);
	if err!=nil {
		http.Error(w,"Invalid request",http.StatusNotFound)
	}
	http.Redirect(w, r, url.OriginalURL, http.StatusFound)
}

func main(){
	//fmt.Println("starting URL shortner....")
	//OriginalURL := "https://www.google.com"
	// generateShortURL(OriginalURL)

	http.HandleFunc("/",RootPageURL)
	http.HandleFunc("/shorten", ShortURLHandler)
	http.HandleFunc("/redirect/",redirectURLHandler)

	fmt.Println("Starting server on port 3000....")
	err := http.ListenAndServe(":3000",nil)
	if err!= nil {
		fmt.Println("Error on starting server:",err)
	}
	
}