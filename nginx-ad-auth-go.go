// Copyright (c) 2022-2024 Ronan LE MEILLAT
// This program is licensed under the AGPLv3 license.
package main

import (
	"crypto/tls"
	_ "embed"
	"flag"
	"fmt"
	"log"
	"net"
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

//go:embed "not-found.html"
var notFoundHTML []byte

// init initializes the application by parsing command line flags and checking environment variables.
// It sets the values for port, ldapURI, ldapBase, adDomain, mailServer, and mailServerPort.
// If any of these values are missing or invalid, it logs a fatal error.
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

// If user hits any other endpoint, return a 404 error with the content of the file not-found.html
func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Not found: %s, IP: %s", r.URL.Path, r.RemoteAddr)
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusNotFound)
	w.Write(notFoundHTML)
}

func main() {
	http.HandleFunc("/auth", authHandler)
	http.HandleFunc("/", notFoundHandler)
	log.Printf("Starting server on port %d", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

// Get hostname from IP address like dig -x
func getHostnameFromIP(ip string) string {
	hostname, err := net.LookupAddr(ip)
	if err != nil {
		return ip
	}
	return hostname[0]
}

// authHandler is a function that handles authentication requests.
// It takes in an http.ResponseWriter and an http.Request as parameters.
// The function retrieves the user, password, and protocol from the request headers.
// If either the user or password is empty, it returns an "Auth-Status: No login or password" error response.
// If the user authentication fails, it returns an "Auth-Status: Invalid login or password" error response.
// The function determines the authentication port based on the protocol and sets the appropriate headers.
// Finally, it sets the "Auth-Status", "Auth-Server", and "Auth-Port" headers and writes a successful response.
func authHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Header.Get("Auth-User")
	pass := r.Header.Get("Auth-Pass")
	client_ip := r.Header.Get("Client-IP")
	client_hostname := r.Header.Get("Client-Host")
	protocol := r.Header.Get("Auth-Protocol")

	if client_hostname == "" {
		client_hostname = getHostnameFromIP(client_ip)
	}

	if user == "" || pass == "" {
		log.Printf("No login or password, nginx IP: %s, client IP: %s, client hostname: %s", r.RemoteAddr, client_ip, client_hostname)
		http.Error(w, "Auth-Status: No login or password", http.StatusOK)
		return
	}

	authenticated, err := authenticateUser(user, pass)
	if err != nil {
		log.Printf("Error authenticating user: %v", err)
		http.Error(w, "Auth-Status: Internal server error", http.StatusInternalServerError)
		return
	}
	if !authenticated {
		log.Printf("Invalid login or password, nginx IP: %s, client IP: %s, client hostname: %s", r.RemoteAddr, client_ip, client_hostname)
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
	log.Printf("Authenticated user: %s, nginx IP: %s, client IP: %s, client hostname: %s", user, r.RemoteAddr, client_ip, client_hostname)
}

// authenticateUser is a function that authenticates a user against an Active Directory server.
// It takes a username and password as parameters and returns a boolean value indicating whether the authentication was successful or not, along with an error if any.
// The function establishes a connection with the LDAP server using the specified ldapURI and binds the user's credentials.
// It then performs a search in the LDAP directory to check if the user exists.
// The search is based on the sAMAccountName attribute, which is the username attribute in Active Directory.
// If the search returns exactly one entry, it means the user exists and the function returns true.
// Otherwise, it returns false.
// If there is any error during the authentication process, it is returned as an error.
func authenticateUser(username, password string) (bool, error) {
	l, err := ldap.DialURL(ldapURI)
	if err != nil {
		return false, fmt.Errorf("failed to connect to LDAP server: %w", err)
	}
	defer l.Close()

	// Reconnect with TLS
	err = l.StartTLS(&tls.Config{InsecureSkipVerify: true})
	if err != nil {
		// If the server does not support StartTLS, return an error
		return false, fmt.Errorf("failed to start TLS: %w", err)
	}

	err = l.Bind(fmt.Sprintf("%s\\%s", adDomain, username), password)
	if err != nil {
		log.Printf("Failed to bind to LDAP: %v", err)
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
		return false, fmt.Errorf("failed to search LDAP: %w", err)
	}

	return len(sr.Entries) == 1, nil
}
