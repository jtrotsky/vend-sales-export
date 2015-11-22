package main

import (
	"flag"
	"strings"

	"github.com/jtrotsky/govend/vend"
	"github.com/jtrotsky/spate/manager"
)

var (
	domainPrefix string
	tz           string
	token        string
)

func main() {

	v := vend.NewClient(token, domainPrefix, tz)
	manager := manager.NewManager(v)

	manager.Run()
}

func init() {

	// Get store info from command line flags.
	flag.StringVar(&domainPrefix, "d", "",
		"The Vend store name (prefix of xxxx.vendhq.com)")
	flag.StringVar(&token, "t", "",
		"Personal API Access Token for the store, generated from Setup -> API Access.")
	flag.StringVar(&tz, "z", "Local",
		"Timezone of the store in zoneinfo format. The default is to try and use the computer's local timezone.")
	flag.Parse()

	// To save people who write DomainPrefix.vendhq.com.
	// Split DomainPrefix on the "." period character then grab the first part.
	parts := strings.Split(domainPrefix, ".")
	domainPrefix = parts[0]
}
