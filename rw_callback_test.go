package acmewrapper

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSaveLoadCallback(t *testing.T) {
	memFS := map[string][]byte{}
	_, err := New(Config{
		Server:           TESTAPI,
		TOSCallback:      TOSAgree,
		Domains:          TESTDOMAINS,
		PrivateKeyFile:   "testinguser.key",
		RegistrationFile: "testinguser.reg",
		Address:          TLSADDRESS,
		TLSCertFile:      "cert.crt",
		TLSKeyFile:       "key.pem",
		SaveFileCallback: func(path string, contents []byte) error {
			memFS[path] = contents
			return nil
		},
		LoadFileCallback: func(path string) ([]byte, error) {
			contents, ok := memFS[path]
			if !ok {
				return nil, os.ErrNotExist
			}
			return contents, nil
		},
	})
	require.NoError(t, err)
	require.Len(t, memFS, 4)
	require.NotNil(t, memFS["testinguser.key"])
	require.NotNil(t, memFS["testinguser.reg"])
	require.NotNil(t, memFS["cert.crt"])
	require.NotNil(t, memFS["key.pem"])
}
