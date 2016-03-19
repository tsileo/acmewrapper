package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/dkumor/acmewrapper"
)

/*
	This example is the equivalent of SimpleHTTPServer in python, with the twist that it does
	Let's Encrypt certificates.
*/

var (
	address = flag.String("host", ":443", "The address at which to run server")
	cert    = flag.String("cert", "cert.crt", "Location of TLS certificate")
	key     = flag.String("key", "key.pem", "Location of TLS key")
	reg     = flag.String("reg", "user.reg", "Location to write user registration")
	priv    = flag.String("priv", "userkey.pem", "Location to store user's private key")
	test    = flag.Bool("test", false, "Use the Let's Encrypt staging server")
	acme    = flag.String("server", acmewrapper.DefaultServer, "The ACME server to use")
	accept  = flag.Bool("accept", false, "Accept the ACME server's TOS?")
	email   = flag.String("email", "", "The email to use when registering")
	help    = flag.Bool("help", false, "Show help message")
)

func main() {
	flag.Parse()
	if *help || flag.NArg() < 2 {
		fmt.Printf("Usage: example -agree mywebsite.com www.mywebsite.com ./www\n will serve the ./www directory with TLS certs for mywebsite.com and www.mywebsite.com\n\n")
		flag.Usage()
	}
	if !*accept {
		fmt.Printf("To run the server, you must accept the Let's Encrypt TOS with -accept")
		os.Exit(1)
	}
	if *test {
		*acme = "https://acme-staging.api.letsencrypt.org/directory"
	}

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(flag.Arg(flag.NArg()-1))))

	w, err := acmewrapper.New(acmewrapper.Config{
		Address: *address,

		Domains: flag.Args()[:flag.NArg()-1],

		Email: *email,

		TLSCertFile: *cert,
		TLSKeyFile:  *key,

		RegistrationFile: *reg,
		PrivateKeyFile:   *priv,

		Server: *acme,

		TOSCallback: acmewrapper.TOSAgree,
	})
	if err != nil {
		fmt.Printf("ERROR: %s", err.Error())
		os.Exit(1)
	}

	tlsconfig := w.TLSConfig()

	listener, err := tls.Listen("tcp", *address, tlsconfig)
	if err != nil {
		fmt.Printf("ERROR: %s", err.Error())
		os.Exit(1)
	}

	fmt.Printf("\n\nRunning server at %s\n\n", *address)

	// In order to enable http2, we can't just use http.Serve in go1.6, so we need
	// to create a manual http.Server, since it needs the tlsconfig
	// https://github.com/golang/go/issues/14374
	server := &http.Server{
		Addr:      *address,
		Handler:   mux,
		TLSConfig: tlsconfig,
	}
	server.Serve(listener)
}
