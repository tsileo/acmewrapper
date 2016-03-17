package acmewrapper

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

const TESTAPI = "https://acme-staging.api.letsencrypt.org/directory"

func TestUserErrors(t *testing.T) {
	_, err := New(Config{
		Server:      TESTAPI,
		TOSCallback: TOSAgree,
	})
	require.Error(t, err)
	_, err = New(Config{
		Server:  TESTAPI,
		Domains: []string{"localhost"},
	})
	require.Error(t, err)

	_, err = New(Config{
		Server:         TESTAPI,
		TOSCallback:    TOSAgree,
		Domains:        []string{"localhost"},
		PrivateKeyFile: "testinguser.key",
	})
	require.Error(t, err)

	_, err = New(Config{
		Server:           TESTAPI,
		TOSCallback:      TOSAgree,
		Domains:          []string{"localhost"},
		RegistrationFile: "testinguser.reg",
	})
	require.Error(t, err)

	_, err = New(Config{
		Server: TESTAPI,
		TOSCallback: func(tosurl string) bool {
			fmt.Printf("TOS URL: %s\n", tosurl)
			return false
		},
		Domains: []string{"localhost"},
	})

	require.Error(t, err)

}

func TestUser(t *testing.T) {
	// Test that an anonymous user can be successfully created

	w, err := New(Config{
		Server:      TESTAPI,
		TOSCallback: TOSAgree,
		Domains:     []string{"localhost"},
	})

	require.NoError(t, err)
	require.Equal(t, w.GetEmail(), "")
	require.NotNil(t, w.GetRegistration())
	require.NotNil(t, w.GetPrivateKey())

	os.Remove("testinguser.key")
	os.Remove("testinguser.reg")

	w, err = New(Config{
		Server:           TESTAPI,
		TOSCallback:      TOSAgree,
		Domains:          []string{"localhost"},
		PrivateKeyFile:   "testinguser.key",
		RegistrationFile: "testinguser.reg",
	})

	require.NoError(t, err)

	// Now that the files are created, it should load fine without TOS
	w, err = New(Config{
		Server:           TESTAPI,
		TOSCallback:      TOSDecline,
		Domains:          []string{"localhost"},
		PrivateKeyFile:   "testinguser.key",
		RegistrationFile: "testinguser.reg",
	})

	require.NoError(t, err)
	require.Equal(t, w.GetEmail(), "")
	require.NotNil(t, w.GetRegistration())
	require.NotNil(t, w.GetPrivateKey())
}
