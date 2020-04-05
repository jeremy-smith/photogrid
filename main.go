package main

import (
	"github.com/julienschmidt/httprouter"
	"html/template"
	"log"
	"net/http"
	"time"
)

type sessionData struct {
	time time.Time
	vars map[string]string
}

var temps *template.Template
var sessions map[string]sessionData

func init() {
	Conf = ReadConf(ConfigFile)
	temps = template.Must(template.ParseGlob("templates/*"))
	dbConnect()
	sessions = make(map[string]sessionData)
}

func main() {
	r := httprouter.New()
	r.GET("/account", account)
	r.POST("/account/create", createAccount)
	r.GET("/admin", admin)
	r.POST("/account/login", login)

	go cleanSessions()

	log.Fatalln(http.ListenAndServe(":8081", r))
}
