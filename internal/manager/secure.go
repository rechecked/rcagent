package manager

import (
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/rechecked/rcagent/internal/config"
)

const (
	INTERNAL_CERT_SERIAL_NUMBER = 0
	DAYS_TO_EXPIRATION          = -30
)

type Cert struct {
	Certificate string `json:"cert"`
	PrivateKey  string `json:"key"`
}

func DecodePEMCert(p []byte) *x509.Certificate {
	block, _ := pem.Decode(p)
	if block == nil {
		return nil
	}

	if block.Type == "CERTIFICATE" {
		cert, _ := x509.ParseCertificate(block.Bytes)
		return cert
	}

	return nil
}

// We need to check cetificate validation every so often and request a new
// certificate if ours is going to expire
func validateCert() {

	certFn := config.GetConfigFilePath("rcagent.pem")
	keyFn := config.GetConfigFilePath("rcagent.key")

	// Check if we need to update the cert
	if !isCertRequestNeeded(certFn) {
		return
	}

	// Request a new cert
	i := getHostInfo()
	data := map[string]string{
		"machineId": i.MachineId,
		"hostname":  i.Hostname,
		"address":   getOutboundIP(),
	}

	resp, err := sendPost("agents/certificate", data)
	if err != nil {
		fmt.Println(err)
		return
	}

	var cert Cert
	err = json.Unmarshal(resp, &cert)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Save certs to location
	err = ioutil.WriteFile(certFn, []byte(cert.Certificate), 0600)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = ioutil.WriteFile(keyFn, []byte(cert.PrivateKey), 0600)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func isCertRequestNeeded(fn string) bool {

	if _, err := os.Stat(fn); err != nil {
		return true
	}

	bytes, err := ioutil.ReadFile(fn)
	if err != nil {
		return false
	}

	cert := DecodePEMCert(bytes)
	if cert == nil {
		return true
	}

	// Internally generated certificates will be overwritten
	if cert.SerialNumber.Int64() == INTERNAL_CERT_SERIAL_NUMBER {
		return true
	}

	// Check if we are close to expiration or not
	if time.Now().After(cert.NotAfter.AddDate(0, 0, DAYS_TO_EXPIRATION)) {
		return true
	}

	return false
}
