// Copyright (c) 2022-2024 Ronan LE MEILLAT
// This program is licensed under the AGPLv3 license.
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/go-ldap/ldap/v3"
)

var (
	port           int
	ldapURI        string
	ldapBase       string
	adDomain       string
	mailServer     string
	mailServerPort int
)

func init() {
	flag.IntVar(&port, "port", 8080, "Port to listen on")
	flag.StringVar(&ldapURI, "ldap-uri", "", "LDAP URI")
	flag.StringVar(&ldapBase, "ldap-base", "", "LDAP base")
	flag.StringVar(&adDomain, "ad-domain", "", "AD domain")
	flag.StringVar(&mailServer, "mail-server", "", "Mail server")
	flag.IntVar(&mailServerPort, "mail-server-port", 0, "Mail server port")
	flag.Parse()

	// Check environment variables
	if envPort := os.Getenv("NGINX_AUTH_PORT"); envPort != "" {
		port, _ = strconv.Atoi(envPort)
	}
	if envLDAPURI := os.Getenv("NGINX_AUTH_LDAP_URI"); envLDAPURI != "" {
		ldapURI = envLDAPURI
	}
	if envLDAPBase := os.Getenv("NGINX_AUTH_LDAP_BASE"); envLDAPBase != "" {
		ldapBase = envLDAPBase
	}
	if envADDomain := os.Getenv("NGINX_AUTH_AD_DOMAIN"); envADDomain != "" {
		adDomain = envADDomain
	}
	if envMailServer := os.Getenv("NGINX_AUTH_MAIL_SERVER"); envMailServer != "" {
		mailServer = envMailServer
	}
	if envMailServerPort := os.Getenv("NGINX_AUTH_MAIL_SERVER_PORT"); envMailServerPort != "" {
		mailServerPort, _ = strconv.Atoi(envMailServerPort)
	}
	if ldapURI == "" || ldapBase == "" || adDomain == "" || mailServer == "" || mailServerPort == 0 {
		log.Fatal("LDAP URI, LDAP base, AD domain, mail server and mail server port are required")
	}
}

func main() {
	http.HandleFunc("/auth", authHandler)
	log.Printf("Starting server on port %d", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Header.Get("Auth-User")
	pass := r.Header.Get("Auth-Pass")
	protocol := r.Header.Get("Auth-Protocol")

	if user == "" || pass == "" {
		http.Error(w, "Auth-Status: No login or password", http.StatusOK)
		return
	}

	if authenticated, err := authenticateUser(user, pass); err != nil || !authenticated {
		http.Error(w, "Auth-Status: Invalid login or password", http.StatusOK)
		return
	}

	authPort := mailServerPort
	switch protocol {
	case "imap":
		authPort = 143
	case "imaps":
		authPort = 993
	case "pop3":
		authPort = 110
	case "pop3s":
		authPort = 995
	case "smtp":
		authPort = 25
	case "smtps":
		authPort = 465
	}

	w.Header().Set("Auth-Status", "OK")
	w.Header().Set("Auth-Server", mailServer)
	w.Header().Set("Auth-Port", strconv.Itoa(authPort))
	w.WriteHeader(http.StatusOK)
}

func authenticateUser(username, password string) (bool, error) {
	l, err := ldap.DialURL(ldapURI)
	if err != nil {
		return false, err
	}
	defer l.Close()

	err = l.Bind(fmt.Sprintf("%s\\%s", adDomain, username), password)
	if err != nil {
		return false, nil
	}

	searchRequest := ldap.NewSearchRequest(
		ldapBase,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(sAMAccountName=%s)", ldap.EscapeFilter(username)),
		[]string{"dn"},
		nil,
	)

	sr, err := l.Search(searchRequest)
	if err != nil {
		return false, err
	}

	return len(sr.Entries) == 1, nil
}
