package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/satori/go.uuid"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"time"
)

func account(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	isAdminIP := isAdminIP(GetClientIP(req))

	qStr := req.URL.Query()

	emsg := qStr.Get("error")
	msg := qStr.Get("msg")

	type pageData struct {
		IsAdminIP bool
		Emsg      string
		Msg       string
	}

	pData := pageData{IsAdminIP: isAdminIP, Emsg: emsg, Msg: msg}

	err := temps.ExecuteTemplate(w, "account.gohtml", pData)
	if err != nil {
		log.Fatalln(err)
	}
}

func createAccount(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	isAdminIP := isAdminIP(GetClientIP(req))
	if !isAdminIP {
		emsg := "You dont have permission to access this page"
		http.Redirect(w, req, "/account?error="+url.QueryEscape(emsg), http.StatusSeeOther)
		return
	}

	email := req.PostFormValue("email")
	password := req.PostFormValue("password")
	password2 := req.PostFormValue("password2")

	// Do some simple validation
	m, err := regexp.MatchString(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`, email)
	if !m || err != nil {
		emsg := "Email did not validate."
		http.Redirect(w, req, "/account?error="+url.QueryEscape(emsg), http.StatusSeeOther)
		return
	}

	if password != password2 {
		emsg := "Passwords do not match"
		http.Redirect(w, req, "/account?error="+url.QueryEscape(emsg), http.StatusSeeOther)
		return
	}

	err = dbCreateAccount(email, password)
	if err != nil {
		emsg := "Could not create account."
		http.Redirect(w, req, "/account?error="+url.QueryEscape(emsg), http.StatusSeeOther)
		return
	}

	msg := "Account Created!"
	http.Redirect(w, req, "/account?msg="+url.QueryEscape(msg), http.StatusSeeOther)
	return
}

func login(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	email := req.PostFormValue("email")
	password := req.PostFormValue("password")

	if validateLogin(email, password) {

		// Setup session
		id := uuid.NewV4()

		sessions[id.String()] = sessionData{
			time: time.Now(),
			vars: map[string]string{"username": email},
		}

		http.SetCookie(w, &http.Cookie{
			Name:   "session",
			Value:  id.String(),
			MaxAge: 20,
			Path:   "/",
		})

		http.Redirect(w, req, "/admin", http.StatusSeeOther)
		return

	} else {
		emsg := "Login failed."
		http.Redirect(w, req, "/account?error="+url.QueryEscape(emsg), http.StatusSeeOther)
		return
	}
}
