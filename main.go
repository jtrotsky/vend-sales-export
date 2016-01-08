package main

import (
	"flag"
	"log"
	"os"
	"strings"

	"github.com/jtrotsky/govend/vend"
	"github.com/jtrotsky/spate/manager"
)

var (
	domainPrefix string
	timeZone     string
	authToken    string
)

func main() {

	vendClient := vend.NewClient(authToken, domainPrefix, timeZone)
	manager := manager.NewManager(vendClient)

	manager.Run()
}

func init() {
	// TODO: Better logging? File? Other?
	log.SetOutput(os.Stderr)

	// Get store info from command line flags.
	flag.StringVar(&domainPrefix, "d", "",
		"The Vend store name (prefix of xxxx.vendhq.com)")
	flag.StringVar(&authToken, "t", "",
		"Personal API Access Token for the store, generated from Setup -> API Access.")
	flag.StringVar(&timeZone, "z", "Local",
		"Timezone of the store in zoneinfo format. The default is to try and use the computer's local timezone.")
	flag.Parse()

	// Check all arguments are given.
	if authToken == "" {
		log.Println(
			"Authentication token not given. Expected like: oe1R9xoQeJRUdyVkz6trbcf9GnUTBovJWKRSBCEf")
		os.Exit(0)
	}
	if domainPrefix == "" {
		log.Println(
			"Domain prefix not given. Expected like: store-name.vendhq.com")
		os.Exit(0)
	}
	if timeZone == "" {
		log.Println(
			"Timezone not given. Expected like: Australia/Melbourne")
		os.Exit(0)
	}

	// To save people who write DomainPrefix.vendhq.com.
	// Split DomainPrefix on the "." period character then grab the first part.
	parts := strings.Split(domainPrefix, ".")
	domainPrefix = parts[0]
}
