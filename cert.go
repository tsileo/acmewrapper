package acmewrapper

import (
	"crypto/tls"
	"crypto/x509"
	"time"

	"github.com/xenolf/lego/acme"
)

// writeCert takes an acme CertificateResource (as returned from the acme.RenewCertificate
// and the acme.ObtainCertificate functions), and writes the cert and key files from it.
// If the files already exist, it renames the old versions by adding .bak to them. This makes
// sure that a little accident doesn't cause too much damage.
func (w *AcmeWrapper) writeCert(certfile, keyfile string, crt acme.CertificateResource) (err error) {
	//only create a backup version if we're definitely saving to disk
	if w.Config.SaveFileCallback == nil {
		if err := w.backupAndSaveFile(certfile, crt.Certificate); err != nil {
			return err
		}
		if err := w.backupAndSaveFile(keyfile, crt.PrivateKey); err != nil {
			return err
		}
	} else {
		if err := w.saveFile(certfile, crt.Certificate); err != nil {
			return err
		}
		if err := w.saveFile(keyfile, crt.PrivateKey); err != nil {
			return err
		}
	}
	return nil
}

func tlsCert(crt acme.CertificateResource) (*tls.Certificate, error) {
	cert, err := tls.X509KeyPair(crt.Certificate, crt.PrivateKey)
	return &cert, err
}

// CertNeedsUpdate returns whether the current certificate either
// does not exist, or is <X days from expiration, where X is set up in config
func (w *AcmeWrapper) CertNeedsUpdate() bool {
	if w.cert == nil {
		// The cert doesn't exist - it certainly needs update
		return true
	}

	// w.cert.Leaf is not set, so we have to manually parse the certs
	// and make sure that they don't expire soon
	for _, c := range w.cert.Certificate {
		crt, err := x509.ParseCertificate(c)
		if err != nil {
			// If there's an error, we assume the cert is broken, and needs update
			return true
		}
		timeLeft := crt.NotAfter.Sub(time.Now().UTC())
		if timeLeft < w.Config.RenewTime {
			return true
		}
	}

	return false
}
