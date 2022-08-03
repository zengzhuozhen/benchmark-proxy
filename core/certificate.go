package core

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"math/big"
	"net"
	"strings"
	"time"
)

var (
	defaultRootCA  *x509.Certificate
	defaultRootKey *rsa.PrivateKey
)

func init() {
	var err error
	crtByte, _ := ioutil.ReadFile("ca.crt")
	keyByte, _ := ioutil.ReadFile("ca.key")
	block, _ := pem.Decode(crtByte)
	defaultRootCA, err = x509.ParseCertificate(block.Bytes)
	if err != nil {
		panic(fmt.Errorf("加载根证书失败: %s", err))
	}
	block, _ = pem.Decode(keyByte)
	defaultRootKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		panic(fmt.Errorf("加载根证书私钥失败: %s", err))
	}
}

type Pair struct {
	Cert            *x509.Certificate
	CertBytes       []byte
	PrivateKey      *rsa.PrivateKey
	PrivateKeyBytes []byte
}

func GenerateTlsConfig(host string) (*tls.Config, error) {
	if h, _, err := net.SplitHostPort(host); err == nil {
		host = h
	}
	pair, err := GeneratePem(host, 1, defaultRootCA, defaultRootKey)
	if err != nil {
		return nil, err
	}
	cert, err := tls.X509KeyPair(pair.CertBytes, pair.PrivateKeyBytes)
	if err != nil {
		return nil, err
	}
	tlsConf := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}
	return tlsConf, nil
}

// GeneratePem 生成证书
func GeneratePem(host string, expireDay int, rootCA *x509.Certificate, rootKey *rsa.PrivateKey) (*Pair, error) {
	var priv *rsa.PrivateKey
	var err error

	priv, err = rsa.GenerateKey(rand.Reader, 2048)

	if err != nil {
		return nil, err
	}
	tmpl := template(host, expireDay)
	derBytes, err := x509.CreateCertificate(rand.Reader, tmpl, rootCA, &priv.PublicKey, rootKey)
	if err != nil {
		return nil, err
	}
	certBlock := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: derBytes,
	}
	serverCert := pem.EncodeToMemory(certBlock)
	privBytes := x509.MarshalPKCS1PrivateKey(priv)
	keyBlock := &pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: privBytes,
	}
	serverKey := pem.EncodeToMemory(keyBlock)

	p := &Pair{
		Cert:            tmpl,
		CertBytes:       serverCert,
		PrivateKey:      priv,
		PrivateKeyBytes: serverKey,
	}

	return p, nil
}

func template(host string, expireDays int) *x509.Certificate {
	fv := fnv.New32a()
	_, _ = fv.Write([]byte(host))

	cert := &x509.Certificate{
		SerialNumber: big.NewInt(int64(fv.Sum32())),
		Subject: pkix.Name{
			CommonName:   host,
			Country:      []string{"CN"},
			Organization: []string{"BenchProxy"},
			Province:     []string{"Guangdong"},
			Locality:     []string{"Shenzhen"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(0, 0, expireDays),
		BasicConstraintsValid: true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageDataEncipherment | x509.KeyUsageKeyEncipherment,
		EmailAddresses:        []string{"zenzz01@mingyuanyun.com"},
	}
	hosts := strings.Split(host, ",")
	for _, item := range hosts {
		if ip := net.ParseIP(host); ip != nil {
			cert.IPAddresses = append(cert.IPAddresses, ip)
			continue
		}

		fields := strings.Split(item, ".")
		fieldNum := len(fields)
		for i := 0; i <= (fieldNum - 2); i++ {
			cert.DNSNames = append(cert.DNSNames, "*."+strings.Join(fields[i:], "."))
		}
		if fieldNum == 2 {
			cert.DNSNames = append(cert.DNSNames, item)
		}
	}

	return cert
}
