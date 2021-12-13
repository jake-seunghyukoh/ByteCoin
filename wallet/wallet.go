package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"github.com/ohshyuk5/ByteCoin/utils"
	"math/big"
	"os"
)

const walletName = "bytecoin.wallet"

var w *wallet

type wallet struct {
	privateKey *ecdsa.PrivateKey
	Address    string
}

func hasWalletFile() bool {
	_, err := os.Stat(walletName)
	return !os.IsNotExist(err)
}

func createPrivateKey() (privateKey *ecdsa.PrivateKey) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	utils.HandleErr(err)
	return
}

func persistPrivateKey(key *ecdsa.PrivateKey) {
	bytes, err := x509.MarshalECPrivateKey(key)
	utils.HandleErr(err)

	err = os.WriteFile(walletName, bytes, os.FileMode(0644))
	utils.HandleErr(err)
}

func restorePivateKey() (key *ecdsa.PrivateKey) {
	bytes, err := os.ReadFile(walletName)
	utils.HandleErr(err)

	key, err = x509.ParseECPrivateKey(bytes)
	utils.HandleErr(err)

	return
}

func encodeBigInts(a, b []byte) string {
	c := append(a, b...)
	return fmt.Sprintf("%x", c)
}

func decodeBigInts(payload string) (*big.Int, *big.Int, error) {
	bytes, err := hex.DecodeString(payload)
	if err != nil {
		return nil, nil, err
	}

	firstHalfBytes := bytes[:len(bytes)/2]
	secondHalfBytes := bytes[len(bytes)/2:]

	bigA, bigB := big.Int{}, big.Int{}
	bigA.SetBytes(firstHalfBytes)
	bigB.SetBytes(secondHalfBytes)
	return &bigA, &bigB, nil
}

func addressFromKey(key *ecdsa.PrivateKey) string {
	x := key.X.Bytes()
	y := key.Y.Bytes()
	return encodeBigInts(x, y)
}

func Sign(w *wallet, payload string) string {
	bytes, err := hex.DecodeString(payload)
	utils.HandleErr(err)

	r, s, err := ecdsa.Sign(rand.Reader, w.privateKey, bytes)
	utils.HandleErr(err)

	return encodeBigInts(r.Bytes(), s.Bytes())
}

func Verify(signature, payload, address string) bool {
	r, s, err := decodeBigInts(signature)
	utils.HandleErr(err)

	x, y, err := decodeBigInts(address)
	utils.HandleErr(err)

	publickKey := ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     x,
		Y:     y,
	}

	payloadBytes, err := hex.DecodeString(payload)
	utils.HandleErr(err)

	valid := ecdsa.Verify(&publickKey, payloadBytes, r, s)
	return valid
}

func Wallet() *wallet {
	if w == nil {
		w = &wallet{}
		var key *ecdsa.PrivateKey

		if hasWalletFile() {
			key = restorePivateKey()
		} else {
			key = createPrivateKey()
			persistPrivateKey(key)
		}

		w.privateKey = key
		w.Address = addressFromKey(key)
	}
	return w
}
