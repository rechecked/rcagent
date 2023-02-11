
package server

import (
    "bytes"
    "crypto/rand"
    "crypto/rsa"
    //"crypto/tls"
    "crypto/x509"
    "crypto/x509/pkix"
    "encoding/pem"
    "os"
    //"fmt"
    //"io/ioutil"
    "math/big"
    "net"
    //"net/http"
    //"net/http/httptest"
    //"strings"
    "time"
)

func GenerateCert() error {

    cert, key, err := selfSignedCert()
    if err != nil {
        return err
    }

    err = writeToFile("rcagent.pem", cert)
    if err != nil {
        return err
    }

    err = writeToFile("rcagent.key", key)
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
    
    // set up our CA certificate
    /*
    ca := &x509.Certificate{
        SerialNumber: big.NewInt(2019),
        Subject: pkix.Name{
            Organization:  []string{"ReChecked"},
            Country:       []string{"US"},
        },
        NotBefore:             time.Now(),
        NotAfter:              time.Now().AddDate(10, 0, 0),
        IsCA:                  true,
        ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
        KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
        BasicConstraintsValid: true,
    }

    // create our private and public key
    caPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
    if err != nil {
        return nil, nil, err
    }

    // create the CA
    caBytes, err := x509.CreateCertificate(rand.Reader, ca, ca, &caPrivKey.PublicKey, caPrivKey)
    if err != nil {
        return nil, nil, err
    }

    // pem encode
    caPEM := new(bytes.Buffer)
    pem.Encode(caPEM, &pem.Block{
        Type:  "CERTIFICATE",
        Bytes: caBytes,
    })

    caPrivKeyPEM := new(bytes.Buffer)
    pem.Encode(caPrivKeyPEM, &pem.Block{
        Type:  "RSA PRIVATE KEY",
        Bytes: x509.MarshalPKCS1PrivateKey(caPrivKey),
    })
    */

    certPEM = new(bytes.Buffer)
    certPrivKeyPEM = new(bytes.Buffer)

    // set up our server certificate
    cert := &x509.Certificate{
        SerialNumber: big.NewInt(2019),
        Subject: pkix.Name{
            Organization:  []string{"ReChecked"},
            Country:       []string{"US"},
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

    /*
    serverCert, err := tls.X509KeyPair(certPEM.Bytes(), certPrivKeyPEM.Bytes())
    if err != nil {
        return certPEM, certPrivKey, err
    }

    serverTLSConf = &tls.Config{
        Certificates: []tls.Certificate{serverCert},
    }

    certpool := x509.NewCertPool()
    certpool.AppendCertsFromPEM(caPEM.Bytes())
    clientTLSConf = &tls.Config{
        RootCAs: certpool,
    }
    */

    return
}
