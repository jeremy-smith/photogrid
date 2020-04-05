package main

import (
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"net/url"
)

func admin(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	c, err := getSessionVars(req)
	if err != nil {
		emsg := "You have been logged out"
		http.Redirect(w, req, "/account?error="+url.QueryEscape(emsg), http.StatusSeeOther)
		return
	}

	err = temps.ExecuteTemplate(w, "admin.gohtml", c.vars)
	if err != nil {
		log.Fatalln(err)
	}
}
