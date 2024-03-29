package server

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"math/big"
	"net"
	"os"
	"time"

	"github.com/rechecked/rcagent/internal/config"
	"github.com/rechecked/rcagent/internal/manager"
)

func GenerateCert(certFn, keyFn string) error {

	// Request a new certificate rather then generate a self signed one
	if config.UsingManager() {
		err := manager.RequestCert(certFn, keyFn)
		if err != nil {
			return err
		}
		return nil
	}

	cert, key, err := selfSignedCert()
	if err != nil {
		return err
	}

	err = writeToFile(certFn, cert)
	if err != nil {
		return err
	}

	err = writeToFile(keyFn, key)
	if err != nil {
		return err
	}

	return nil
}

func writeToFile(file string, b *bytes.Buffer) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()
	n, err := b.WriteTo(f)
	if err != nil {
		return err
	}
	if int(n) < b.Len() {
		return errors.New("could not write certificate to file")
	}
	return nil
}

func selfSignedCert() (certPEM *bytes.Buffer, certPrivKeyPEM *bytes.Buffer, err error) {

	certPEM = new(bytes.Buffer)
	certPrivKeyPEM = new(bytes.Buffer)

	certPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return certPEM, certPrivKeyPEM, err
	}

	// Get SubjectKeyId from priv key
	keyBytes := x509.MarshalPKCS1PublicKey(&certPrivKey.PublicKey)
	keyHash := sha1.Sum(keyBytes)
	ski := keyHash[:]

	// Set up our certificate
	cert := &x509.Certificate{
		SerialNumber: big.NewInt(time.Now().UnixMicro()),
		Subject: pkix.Name{
			Organization: []string{"ReChecked Agent"},
		},
		IPAddresses:  []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(10, 0, 0),
		SubjectKeyId: ski,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, cert, cert, &certPrivKey.PublicKey, certPrivKey)
	if err != nil {
		return certPEM, certPrivKeyPEM, err
	}

	pem.Encode(certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})

	pem.Encode(certPrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(certPrivKey),
	})

	return
}
