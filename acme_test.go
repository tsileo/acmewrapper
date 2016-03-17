package acmewrapper

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUserErrors(t *testing.T) {
	_, err := New(Config{
		Server:      TESTAPI,
		TOSCallback: TOSAgree,
		Address:     TLSADDRESS,
	})
	require.Error(t, err)
	_, err = New(Config{
		Server:  TESTAPI,
		Domains: TESTDOMAINS,
		Address: TLSADDRESS,
	})
	require.Error(t, err)

	_, err = New(Config{
		Server:         TESTAPI,
		TOSCallback:    TOSAgree,
		Domains:        TESTDOMAINS,
		Address:        TLSADDRESS,
		PrivateKeyFile: "testinguser.key",
	})
	require.Error(t, err)

	_, err = New(Config{
		Server:           TESTAPI,
		TOSCallback:      TOSAgree,
		Domains:          TESTDOMAINS,
		Address:          TLSADDRESS,
		RegistrationFile: "testinguser.reg",
	})
	require.Error(t, err)

	_, err = New(Config{
		Server: TESTAPI,
		TOSCallback: func(tosurl string) bool {
			fmt.Printf("TOS URL: %s\n", tosurl)
			return false
		},
		Domains: TESTDOMAINS,
	})

	require.Error(t, err)

}

func TestUser(t *testing.T) {
	// Test that an anonymous user can be successfully created

	w, err := New(Config{
		Server:      TESTAPI,
		TOSCallback: TOSAgree,
		Domains:     TESTDOMAINS,
		Address:     TLSADDRESS,
	})

	require.NoError(t, err)
	require.Equal(t, w.GetEmail(), "")
	require.NotNil(t, w.GetRegistration())
	require.NotNil(t, w.GetPrivateKey())
	require.NotNil(t, w.GetCertificate())

	os.Remove("testinguser.key")
	os.Remove("testinguser.reg")

	w, err = New(Config{
		Server:           TESTAPI,
		TOSCallback:      TOSAgree,
		Domains:          TESTDOMAINS,
		PrivateKeyFile:   "testinguser.key",
		RegistrationFile: "testinguser.reg",
		Address:          TLSADDRESS,
	})

	require.NoError(t, err)

	// Now that the files are created, it should load fine without TOS
	w, err = New(Config{
		Server:           TESTAPI,
		TOSCallback:      TOSDecline,
		Domains:          TESTDOMAINS,
		PrivateKeyFile:   "testinguser.key",
		RegistrationFile: "testinguser.reg",
		Address:          TLSADDRESS,
	})

	require.NoError(t, err)
	require.Equal(t, w.GetEmail(), "")
	require.NotNil(t, w.GetRegistration())
	require.NotNil(t, w.GetPrivateKey())
}
