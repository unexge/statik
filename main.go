package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	fmt.Println("statik - v0.0.1")

	b, err := ioutil.ReadFile("statik.yml")
	if err != nil {
		log.Fatalf("cannot read config file: %q\n", err)
	}

	config, err := NewConfigFromYaml(b)
	if err != nil {
		log.Fatalf("cannot create config: %q\n", err)
	}

	handler, err := NewHandler(config)
	if err != nil {
		log.Fatalf("cannot create handler: %q\n", err)
	}

	if config.IsHTTPS() {
		fmt.Printf("listening https://%s\n", config.Listen)
		http.ListenAndServeTLS(config.Listen, config.HTTPS.Cert, config.HTTPS.Key, handler)
	} else {
		fmt.Printf("listening http://%s\n", config.Listen)
		http.ListenAndServe(config.Listen, handler)
	}
}
