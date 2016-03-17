package acmewrapper

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"

	"github.com/xenolf/lego/acme"
)

// initACME initailizes the acme client - it does everything from reading/writing the
// user private key and registration files, to ensuring that the user is registered
// on the ACME server and has accepted the TOS.
// It expects w.config to be set up.
// It sets up:
//	- w.privatekey
//	- w.registration
//	- w.client (init the user + agree to TOS)
//
// Its input is whether there is a server running already. If the server is running,
// then the SNI query will succeed. If it isn't (ie, we are just setting up), then
// initACME must set up its own temporary server to get any initial certificates.
func (w *AcmeWrapper) initACME(serverRunning bool) (err error) {
	// We are modifying and using some of the config properties, so lock them
	w.configmutex.Lock()
	defer w.configmutex.Unlock()

	// Just in case initACME is being run on an existing AcmeWrapper
	w.registration = nil
	w.privatekey = nil

	if len(w.config.Domains) == 0 {
		return errors.New("No domains set - can't initialize ACME client")
	}

	if w.config.PrivateKeyFile != "" {
		if w.config.RegistrationFile == "" {
			return errors.New("A filename was set for the private key but not the registration file")
		}

		// We are to use file-backed registration. See if the files exist already. We first load
		// the key file, then we load the registration file

		w.privatekey, err = LoadPrivateKey(w.config.PrivateKeyFile)
		if err != nil {
			if !os.IsNotExist(err) {
				return err
			}
			// The private key file doesn't exist - w.privatekey is left at nil
		}

		f, err := os.Open(w.config.RegistrationFile)
		if err == nil {
			defer f.Close()

			if err = json.NewDecoder(f).Decode(&w.registration); err != nil {
				return err
			}

		} else if !os.IsNotExist(err) {
			return err
		}

		// If only one exists, but not the other, return an error. Reemember that these are nil if the file didn't exist
		if (w.privatekey != nil || w.registration != nil) && (w.privatekey == nil || w.registration == nil) {
			return errors.New("One of the files (registration or key) exists, but the other is missing")
		}

	} else if w.config.RegistrationFile != "" {
		return errors.New("A filename was set for the registration file but not the private key")
	}

	if w.privatekey == nil {
		// If privatekey is nil, it means that either there are no files, or we are running in memory only
		// Whatever the case, we generate our acme user

		// Generate the key
		if w.config.PrivateKeyType == acme.RSA2048 {
			w.privatekey, err = rsa.GenerateKey(rand.Reader, 2048)
		} else if w.config.PrivateKeyType == acme.RSA4096 {
			w.privatekey, err = rsa.GenerateKey(rand.Reader, 4096)
		} else if w.config.PrivateKeyType == acme.RSA8192 {
			w.privatekey, err = rsa.GenerateKey(rand.Reader, 8192)
		} else if w.config.PrivateKeyType == acme.EC256 {
			w.privatekey, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		} else if w.config.PrivateKeyType == acme.EC384 {
			w.privatekey, err = ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
		} else {
			return errors.New("Unrecognized key type")
		}
		if err != nil {
			return err
		}

		if w.config.PrivateKeyFile != "" {
			// If we are to use a file, write it now
			if err = SavePrivateKey(w.config.PrivateKeyFile, w.privatekey); err != nil {
				return err
			}
		}

	}

	// Now that we have the key and necessary setup info, we prepare the ACME client.
	w.client, err = acme.NewClient(w.config.Server, w, w.config.PrivateKeyType)
	if err != nil {
		return err
	}

	if w.registration == nil {
		// There is no registration - register with the ACME server
		w.registration, err = w.client.Register()
		if err != nil {
			return err
		}

		if !w.config.TOSCallback(w.registration.TosURL) {
			return errors.New("Terms of service were not accepted")
		}

		if err = w.client.AgreeToTOS(); err != nil {
			return err
		}

		// If we are to use a registration file, write the file now
		if w.config.RegistrationFile != "" {
			jsonBytes, err := json.MarshalIndent(w.registration, "", "\t")
			if err != nil {
				return err
			}
			if err = ioutil.WriteFile(w.config.RegistrationFile, jsonBytes, 0600); err != nil {
				return err
			}
		}
	}

	// Now that the user and client basics are intialized, we set up the client
	// so that it uses our custom SNI provider and disables the rest. We don't want
	// to start custom servers, but rather plug into our certificate updater
	w.client.SetChallengeProvider(acme.TLSSNI01, wrapperChallengeProvider{
		w: w,
	})
	// All of the challenges are disabled EXCEPT SNI
	w.client.ExcludeChallenges([]acme.Challenge{acme.HTTP01, acme.DNS01})

	return nil

}
