package acmewrapper

import (
	"fmt"
	"os"
	"testing"
)

const TESTAPI = "https://acme-staging.api.letsencrypt.org/directory"

var TESTDOMAINS []string
var TLSADDRESS string

func TestMain(m *testing.M) {

	tlsaddr := os.Getenv("TLSADDRESS")
	if tlsaddr == "" {
		tlsaddr = ":443"
	}
	// Set up the domain to use for tests
	dom := os.Getenv("DOMAIN_NAME")
	if dom == "" {
		fmt.Printf("NO DOMAIN SET\n\tSet a valid testing domain name\n\tin the DOMAIN_NAME environmental variable:\n\n\t\texport DOMAIN_NAME=\"example.com\"\n\n")
		os.Exit(-1)
	}
	fmt.Printf("USING DOMAIN_NAME='%s'\nUSING TLSADDRESS='%s'\n", dom, tlsaddr)
	TESTDOMAINS = []string{dom}
	TLSADDRESS = tlsaddr

	retCode := m.Run()

	// call with result of m.Run()
	os.Exit(retCode)
}
