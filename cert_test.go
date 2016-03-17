package acmewrapper

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCert(t *testing.T) {
	os.Remove("cert.crt")
	os.Remove("key.pem")

	w, err := New(Config{
		Server:           TESTAPI,
		TOSCallback:      TOSAgree,
		Domains:          TESTDOMAINS,
		PrivateKeyFile:   "testinguser.key",
		RegistrationFile: "testinguser.reg",
		Address:          TLSADDRESS,

		TLSCertFile: "cert.crt",
		TLSKeyFile:  "key.pem",
	})

	require.NoError(t, err)

	c := w.GetCertificate().Certificate[0]

	// Make sure the files were written
	_, err = os.Stat("cert.crt")
	require.False(t, os.IsNotExist(err))
	_, err = os.Stat("key.pem")
	require.False(t, os.IsNotExist(err))

	// Make sure that the key and cert were generated correctly - make the TOS fail,
	// and the renew callback fail, since ACME shouldn't need to be set at all
	hadfailure := false
	w, err = New(Config{
		Server:           TESTAPI,
		TOSCallback:      TOSDecline,
		Domains:          TESTDOMAINS,
		PrivateKeyFile:   "testinguser.key",
		RegistrationFile: "testinguser.reg",
		Address:          TLSADDRESS,

		TLSCertFile: "cert.crt",
		TLSKeyFile:  "key.pem",

		RenewCallback: func() {
			hadfailure = true
		},
	})

	require.NoError(t, err)
	require.False(t, hadfailure)
	require.True(t, bytes.Equal(w.GetCertificate().Certificate[0], c))

	// Now make sure that we can load without ACME enabled using our currentkeys
	w, err = New(Config{
		Server:       TESTAPI,
		AcmeDisabled: true,

		TLSCertFile: "cert.crt",
		TLSKeyFile:  "key.pem",
	})

	require.NoError(t, err)
	require.False(t, hadfailure)
	require.True(t, bytes.Equal(w.GetCertificate().Certificate[0], c))

	// Lastly: Make sure we can start without ACME enabled, but enable it later.
	// NOTE: This also tests our renewal function
	renewnum := 0
	w, err = New(Config{
		AcmeDisabled: true,
		Server:       TESTAPI,
		TOSCallback:  TOSAgree,
		Domains:      TESTDOMAINS,
		Address:      TLSADDRESS,

		TLSCertFile: "cert.crt",
		TLSKeyFile:  "key.pem",

		RenewCallback: func() {
			renewnum++
		},

		RenewTime:  999999999999999, // A ridiculous value so that renew always happens
		RenewCheck: 3,
		RetryDelay: 3,

		RenewFailedCallback: func(err error) {
			require.NoError(t, err)
		},
	})

	require.NoError(t, err)

	// Now start a server with the config
	listener, err := tls.Listen("tcp", TLSADDRESS, w.TLSConfig())
	require.NoError(t, err)
	go func() {
		http.Serve(listener, nil)
	}()

	fmt.Printf("acmeenable\n")
	require.NoError(t, w.AcmeDisabled(false))

	// Now the certificate should be set
	fmt.Printf("getcert\n")
	crt := w.GetCertificate()
	fmt.Printf("Sleeping for 8 seconds...\n")
	time.Sleep(8 * time.Second)
	fmt.Printf("Done sleeping\n")

	// The certificate should be renewed
	require.NotEqual(t, crt, w.GetCertificate())
	require.True(t, renewnum >= 2)

	// Stop it from being annoying in the background anymore
	w.Config.RenewCheck = 9999999999
	w.Config.RetryDelay = 9999999999
	w.Config.RenewTime = 500
	listener.Close()
}
