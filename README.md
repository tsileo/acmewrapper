# ACMEWrapper

Add Let's Encrypt support to your golang server in 10 lines of code.

[![GoDoc](https://godoc.org/github.com/dkumor/acmewrapper?status.svg)](https://godoc.org/github.com/dkumor/acmewrapper)

```go
w, err := acmewrapper.New(acmewrapper.Config{
	Domains: []string{"example.com","www.example.com"},
	Address: ":443",

	TLSCertFile: "cert.pem",
	TLSKeyFile:  "key.pem",

	// Let's Encrypt stuff
	RegistrationFile: "user.reg",
	PrivateKeyFile:   "user.pem",

	TOSCallback: acmewrapper.TOSAgree,
})


if err!=nil {
	log.Fatal("Let's Encrypt: ", err)
}

listener, err := tls.Listen("tcp", ":443", w.TLSConfig())
```

Acmewrapper is built upon https://github.com/xenolf/lego, and handles all certificate generation, renewal
and replacement automatically. After the above code snippet, your certificate will automatically be renewed 30 days before expiring without downtime. Any files that don't exist will be created, and your "cert.pem" and "key.pem" will be kept up to date.

Since Let's Encrypt is usually an option that can be turned off, the wrapper allows disabling ACME support and just using normal certificates, with the bonus of allowing live reload (ie: change your certificates during runtime).

And finally, *technically*, none of the file names shown above are actually necessary. The only needed fields are Domains and TOSCallback. Without the given file names, acmewrapper runs in-memory. Beware, though: if you do that, you might run into rate limiting from Let's Encrypt if you restart too often!

**WARNING:** This code literally JUST started working. Do NOT use it anywhere important before it has some time to mature.

## How It Works

Let's Encrypt has SNI support for domain validation. That means we can update our certificate if we control the TLS configuration of a server. That is exactly what acmewrapper does. Not only does it transparently update your server's certificate, but it uses its control of SNI to pass validation tests.

This means that *no other changes* are needed to your code. You don't need any special handlers or hidden directories. So long as acmewrapper is able to set your TLS configuration, and your TLS server is running on port 443, you can instantly have a working Let's Encrypt certificate.

**NOTE:** The ability to manage certificates in this way was enabled in go 1.5 - acmewrapper will not work with older go versions.

## Example

### Old Code

This is sample code before adding Let's Encrypt support:

```go
package main

import (
    "io"
    "net/http"
    "log"
)

func HelloServer(w http.ResponseWriter, req *http.Request) {
    io.WriteString(w, "hello, world!\n")
}

func main() {
    http.HandleFunc("/hello", HelloServer)
    err := http.ListenAndServeTLS(":443", "cert.pem", "key.pem", nil)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}
```

### New Code

```go
package main

import (
    "io"
    "net/http"
    "log"
	"crypto/tls"

	"github.com/dkumor/acmewrapper"
)

func HelloServer(w http.ResponseWriter, req *http.Request) {
    io.WriteString(w, "hello, world!\n")
}

func main() {
    http.HandleFunc("/hello", HelloServer)

	w, err := acmewrapper.New(acmewrapper.Config{
		Domains: []string{"example.com","www.example.com"},
		Address: ":443",

		TLSCertFile: "cert.pem",
		TLSKeyFile:  "key.pem",

		RegistrationFile: "user.reg",
		PrivateKeyFile:   "user.pem",

		TOSCallback: acmewrapper.TOSAgree,
	})


	if err!=nil {
		log.Fatal("Let's Encrypt: ", err)
	}

	listener, err := tls.Listen("tcp", ":443", w.TLSConfig())
    if err != nil {
        log.Fatal("Listener: ", err)
    }
	http.Serve(listener, nil)
}
```

## Testing

Running the tests is a bit of a chore, since it requires a valid domain name, and access to port 443.
This is because ACMEWrapper uses the Let's Encrypt staging server to make sure the code is working.

To test on your own server, you need to change the domain name to your domain, and set a custom testing port
that will be routed to 443:

```bash
sudo iptables -t nat -A PREROUTING -p tcp --dport 443 -j REDIRECT --to-port 8443
export TLSADDRESS=":8443"
export DOMAIN_NAME="example.com"
go test
```

To delete the port forwarding rule, find it in your tables
```bash
sudo iptables -t nat --line-numbers -n -L
```

And delete the number that it was at
```bash
iptables -t nat -D PREROUTING 2
```

(This is based on http://serverfault.com/questions/112795/how-can-i-run-a-server-on-linux-on-port-80-as-a-normal-user)
