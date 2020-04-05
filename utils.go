package main

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net"
	"net/http"
	"time"
)

func GetClientIP(r *http.Request) string {
	ipAndPort := _readClientIP(r)
	ip, _, err := net.SplitHostPort(ipAndPort)
	if err != nil {
		log.Fatalln(err)
	}
	return ip
}

func _readClientIP(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}
	return IPAddress
}

func isAdminIP(ip string) bool {
	for _, v := range Conf.AdminIps {
		if v == ip {
			return true
		}
	}
	return false
}

func HashPassword(s string) string {
	bs, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalln(err)
	}
	return string(bs)
}

func ValidatePassword(h string, p string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(h), []byte(p)); err != nil {
		return false
	}
	return true
}

func cleanSessions() {
	for {
		for k, v := range sessions {
			if v.time.Before(time.Now().Add(-20 * time.Second)) {
				delete(sessions, k)
			}
		}
		time.Sleep(5 * time.Second)
	}
}

func getSessionVars(req *http.Request) (sessionData, error) {
	c, err := req.Cookie("session")
	if err != nil {
		return sessionData{}, http.ErrNoCookie
	}
	v, ok := sessions[c.Value]
	if !ok {
		return sessionData{}, errors.New("session not found")
	}
	return v, nil
}

func setSessionVar(k string, v string, req *http.Request) error {
	c, err := req.Cookie("session")
	if err != nil {
		return http.ErrNoCookie
	}
	_, ok := sessions[c.Value]
	if !ok {
		return errors.New("session not found")
	}
	sessions[c.Value].vars[k] = v
	return nil
}
