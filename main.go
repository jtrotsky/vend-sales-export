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
	// Get store info from command line flags.
	flag.StringVar(&domainPrefix, "d", "", "Vend store name (prefix of xxxx.vendhq.com)")
	flag.StringVar(&authToken, "t", "", "API Access Token, from Setup -> Personal Tokens.")
	flag.StringVar(&timeZone, "z", "Local", "Timezone in zoneinfo format."+
		"Default is computer's local timezone.")
	flag.Parse()

	// Check all required arguments are given.
	if authToken == "" {
		log.Println("Authentication token not given." +
			"Expected like: oe1R9xoQeJRUdyVkz6trbcf9GnUTBovJWKRSBCEf")
		os.Exit(0)
	}
	if domainPrefix == "" {
		log.Println("Domain prefix not given. Expected like: store-name.vendhq.com")
		os.Exit(0)
	}
	if timeZone == "" {
		log.Println("Timezone not given. Expected like: Australia/Melbourne")
		os.Exit(0)
	}

	// To save people who write DomainPrefix.vendhq.com.
	// Split DomainPrefix on the "." period character then grab the first part.
	parts := strings.Split(domainPrefix, ".")
	domainPrefix = parts[0]
}
