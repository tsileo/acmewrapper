package acmewrapper

import "time"

// backgroundExpirationChecker is exactly that - it runs in the background
// and ensures that messages regarding certificate expiration as well as
// any renewals if ACME is configured are run on time.
func backgroundExpirationChecker(w *AcmeWrapper) {
	for {
		time.Sleep(time.Duration(w.config.RenewTime) * time.Second)
		if w.CertNeedsUpdate() {
			for {

				if w.config.RenewCallback != nil {
					w.config.RenewCallback()
				}
				if !w.config.AcmeDisabled {
					err := w.Renew()
					if err != nil && w.config.RenewFailedCallback {
						w.config.RenewFailedCallback(err)
					}
				}
				if !w.CertNeedsUpdate() {
					break
				}
				time.Sleep(time.Duration(w.config.RetryDelay) * time.Second)
			}
		}

	}
}
