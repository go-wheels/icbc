package icbc

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"strings"
)

const (
	pemDelim              = "-----"
	pemTagBegin           = pemDelim + "BEGIN"
	pemTagEnd             = pemDelim + "END"
	pemLabelRSAPrivateKey = "RSA PRIVATE KEY"
	pemLabelPublicKey     = "PUBLIC KEY"
)

func parseRSAPublicKey(data string) (key *rsa.PublicKey, err error) {
	block, err := decodePEMData(data, pemLabelPublicKey)
	if err != nil {
		return
	}
	itf, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return
	}
	key, ok := itf.(*rsa.PublicKey)
	if !ok {
		err = errors.New("icbc: unknown type of public key")
	}
	return
}

func parseRSAPrivateKey(data string) (key *rsa.PrivateKey, err error) {
	block, err := decodePEMData(data, pemLabelRSAPrivateKey)
	if err != nil {
		return
	}
	itf, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return
	}
	key, ok := itf.(*rsa.PrivateKey)
	if !ok {
		err = errors.New("icbc: unknown type of private key")
	}
	return
}

func decodePEMData(data string, label string) (block *pem.Block, err error) {
	standardizedData := standardizePEMData(data, label)
	block, _ = pem.Decode([]byte(standardizedData))
	if block == nil {
		err = errors.New("icbc: failed to decode PEM block")
	}
	return
}

func standardizePEMData(data string, label string) string {
	standardizedData := strings.TrimSpace(data)
	if !strings.Contains(standardizedData, pemTagBegin) {
		standardizedData = pemTagBegin + " " + label + " " + pemDelim + "\n" + standardizedData
	}
	if !strings.Contains(standardizedData, pemTagEnd) {
		standardizedData += "\n" + pemTagEnd + " " + label + " " + pemDelim
	}
	return standardizedData
}
