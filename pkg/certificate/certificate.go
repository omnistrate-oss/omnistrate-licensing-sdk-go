package certificate

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	_ "embed"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"time"
)

// Root cert for RSA from https://letsencrypt.org/certificates/
//
//go:embed isrgrootx1.pem
var rootPEM []byte

// Intermediate cert for RSA from https://letsencrypt.org/certificates/
//
//go:embed r10.pem
var intermediatePEM10 []byte

// Intermediate cert for RSA from https://letsencrypt.org/certificates/
//
//go:embed r11.pem
var intermediatePEM11 []byte

func LoadCertificate(certPath string) (*x509.Certificate, error) {
	data, err := os.ReadFile(certPath)
	if err != nil {
		return nil, err
	}
	return LoadCertificateFromBytes(data)
}

func LoadCertificateFromBytes(data []byte) (*x509.Certificate, error) {
	cert, _ := Decode(data)
	if len(cert) == 0 {
		return nil, fmt.Errorf("failed to parse certificate PEM")
	}
	return x509.ParseCertificate(cert)
}

func LoadPrivateKey(keyPath string) (*rsa.PrivateKey, error) {
	data, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}
	return LoadPrivateKeyFromBytes(data)
}

func LoadPrivateKeyFromBytes(data []byte) (*rsa.PrivateKey, error) {
	cert, _ := Decode(data)
	if len(cert) == 0 {
		return nil, fmt.Errorf("failed to Sig certificate PEM")
	}

	// Parse the RSA private key
	return x509.ParsePKCS1PrivateKey(cert)
}

// Decode certificate from PEM
func Decode(pemCert []byte) (cert []byte, err error) {
	// Decode the PEM certificate
	var pemBlock *pem.Block
	if pemBlock, _ = pem.Decode(pemCert); pemBlock == nil {
		err = errors.New("failed to decode certificate from PEM")
		return
	}

	// Return the decoded certificate
	cert = pemBlock.Bytes
	return
}

func Sign(key *rsa.PrivateKey, data []byte) ([]byte, error) {
	return key.Sign(rand.Reader, hash(data), crypto.SHA256)
}

func VerifySignature(cert *x509.Certificate, signature, data []byte) error {
	return rsa.VerifyPKCS1v15(cert.PublicKey.(*rsa.PublicKey), crypto.SHA256, hash(data), signature)
}

func hash(data []byte) []byte {
	hash := sha256.New()
	hash.Write(data)
	return hash.Sum(nil)
}

func VerifyCertificate(cert *x509.Certificate, dnsName string, currentTime time.Time) error {
	// Create a certificate pool with CA and intermediate certificates
	rootsPool := x509.NewCertPool()
	rootsPool.AppendCertsFromPEM(rootPEM)

	intermediatesPool := x509.NewCertPool()
	intermediatesPool.AppendCertsFromPEM(intermediatePEM10)
	intermediatesPool.AppendCertsFromPEM(intermediatePEM11)

	_, err := cert.Verify(x509.VerifyOptions{
		DNSName:       dnsName,
		CurrentTime:   currentTime,
		Roots:         rootsPool,
		Intermediates: intermediatesPool,
	})
	return err
}
