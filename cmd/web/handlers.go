package main

import (
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/google/uuid"
	"github.com/sumeshmurali/urlshortner/internal/database"
)

var repo database.Repository

func isValidUrl(rawUrl string) bool {
	if rawUrl == "" {
		return false
	}
	u, err := url.Parse(rawUrl)
	if err != nil {
		return false
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return false
	}
	if u.Host == "" {
		return false
	}
	return true
}

func index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	t, err := template.ParseFiles(
		"ui/html/base.html",
		"ui/html/index.html",
	)
	if err != nil {
		log.Printf("web.handlers.index: template.ParseFiles failed with %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	t.ExecuteTemplate(w, "base", nil)
}

func shortenUrlHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	rawUrl := r.Form.Get("url")
	if !isValidUrl(rawUrl) {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	id, err := uuid.NewRandom()
	if err != nil {
		log.Printf("web.handlers.shortenUrlHandle: Failed generating uuid with %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	_, err = repo.SetLongUrl(id.String(), rawUrl)
	if err != nil {
		log.Printf("web.handlers.shortenUrlHandle: repo.SetLongUrl failed with %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	t, err := template.ParseFiles(
		"ui/html/base.html",
		"ui/html/result.html",
	)
	if err != nil {
		log.Printf("web.handlers.shortenUrlHandle: template.ParseFiles failed with %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	shortUrl := "http://localhost:8080/s/" + id.String()
	w.WriteHeader(http.StatusCreated)
	t.ExecuteTemplate(w, "base", shortUrl)
}

func shortUrlHandle(w http.ResponseWriter, r *http.Request) {
	urlSplit := strings.Split(r.URL.Path, "/s/")
	if len(urlSplit) != 2 {
		http.Error(w, "Url is not found on our server. It is either deleted or expired. We don't have any other information", http.StatusNotFound)
		return
	}
	id := urlSplit[1]
	if id == "" {
		http.Error(w, "Url is not found on our server. It is either deleted or expired. We don't have any other information", http.StatusNotFound)
		return
	}
	urlDetail, err := repo.GetUrlDetail(id)
	if err != nil {
		if strings.Contains(err.Error(), "no record with url") {
			http.Error(w, "Url is not found on our server. It is either deleted or expired. We don't have any other information", http.StatusNotFound)
		} else {
			log.Printf("web.handlers.shortUrlHandle: failed with %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}
	// lets send this to background
	go func(u database.UrlDetail) {
		repo.RecordMeta(u.ID, r.RemoteAddr, "", "")
	}(urlDetail)

	w.Header().Set("location", urlDetail.LongUrl)
	w.WriteHeader(http.StatusPermanentRedirect)
	// w.Write([]byte(fmt.Sprintf("Redirecting you to the target url in a moment - %v", url)))
}

func GetRouter() *http.ServeMux {
	repo = database.NewRepository()
	err := repo.Init()
	if err != nil {
		panic(err)
	}
	mux := http.ServeMux{}
	mux.HandleFunc("/", index)
	mux.HandleFunc("/shorten", shortenUrlHandle)
	mux.HandleFunc("/s/", shortUrlHandle)
	return &mux
}
