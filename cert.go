package acmewrapper

import (
	"time"

	"github.com/xenolf/lego/acme"
)

// writeCert takes an acme CertificateResource (as returned from the acme.RenewCertificate
// and the acme.ObtainCertificate functions), and writes the cert and key files from it.
// If the files already exist, it renames the old versions by adding .bak to them. This makes
// sure that a little accident doesn't cause too much damage.
func writeCert(certfile, keyfile string, crt acme.CertificateResource) error {
	//crt.
}

// Renew generates a new certificate
func (w *AcmeWrapper) Renew() error {
	w.configmutex.Lock()
	defer w.configmutex.Unlock()

}

// CertNeedsUpdate returns whether the current certificate either
// does not exist, or is <X days from expiration, where X is set up in config
func (w *AcmeWrapper) CertNeedsUpdate() bool {
	if !w.cert {
		// The cert doesn't exist - it certainly needs update
		return true
	}
	timeLeft := w.cert.Leaf.NotAfter.Sub(time.Now().UTC())
	return int64(timeLeft.Seconds()) < w.config.RenewTime
}
