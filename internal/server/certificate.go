package server

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"github.com/rechecked/rcagent/internal/config"
	"math/big"
	"net"
	"os"
	"time"
)

func GenerateCert() error {

	cert, key, err := selfSignedCert()
	if err != nil {
		return err
	}

	err = writeToFile(config.ConfigDir+"/rcagent.pem", cert)
	if err != nil {
		return err
	}

	err = writeToFile(config.ConfigDir+"/rcagent.key", key)
	if err != nil {
		return err
	}

	return nil
}

func writeToFile(file string, bytes *bytes.Buffer) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = bytes.WriteTo(f)
	if err != nil {
		return err
	}
	return nil
}

func selfSignedCert() (certPEM *bytes.Buffer, certPrivKeyPEM *bytes.Buffer, err error) {

	certPEM = new(bytes.Buffer)
	certPrivKeyPEM = new(bytes.Buffer)

	// Set up our certificate
	cert := &x509.Certificate{
		SerialNumber: big.NewInt(2019),
		Subject: pkix.Name{
			Organization: []string{"ReChecked"},
			Country:      []string{"US"},
		},
		IPAddresses:  []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(10, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}

	certPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return certPEM, certPrivKeyPEM, err
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
