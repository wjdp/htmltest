package htmltest

import (
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
)

type CertChainErr struct {
	cert    *x509.Certificate
	chain   *x509.CertPool
	hintErr error
}

func (e CertChainErr) Error() string {
	s := "x509: could not validate certificate chain"
	if e.hintErr != nil {
		s += fmt.Sprintf(" (possibly because of %q)", e.hintErr)
	}
	return s
}

func statusCodeValid(code int) bool {
	return code == http.StatusPartialContent || code == http.StatusOK
}

// fetchCertificate : Retrieve and parse a single certificate. Split out of
// validateCertChain so the response body is closed on every exit path; the
// loop it is called from returns early on each error.
func fetchCertificate(url string) (*x509.Certificate, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("could not fetch certificate at %s (status %d)",
			url, resp.StatusCode)
	}

	certBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return x509.ParseCertificate(certBytes)
}

func validateCertChain(cert *x509.Certificate) (err error) {
	if cert.IssuingCertificateURL == nil {
		return CertChainErr{cert: cert}
	}

	intermediates := x509.NewCertPool()
	//roots, err := x509.SystemCertPool()
	if err != nil {
		return CertChainErr{cert: cert}
	}
	var certsToFetch []string = cert.IssuingCertificateURL

	for i := 0; i < len(certsToFetch); i++ {
		url := certsToFetch[i]

		newCert, err := fetchCertificate(url)
		if err != nil {
			return CertChainErr{cert: cert, chain: intermediates, hintErr: err}
		}

		if newCert.CheckSignatureFrom(cert) == nil {
			// we have out root
			break
		}

		intermediates.AddCert(newCert)

		if newCert.IssuingCertificateURL != nil {
			certsToFetch = append(certsToFetch, newCert.IssuingCertificateURL...)
		}
	}

	_, err = cert.Verify(x509.VerifyOptions{
		Intermediates: intermediates,
	})
	return
}
