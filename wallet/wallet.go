package wallet

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"github.com/ohshyuk5/ByteCoin/utils"
	"math/big"
)

const hashedMessage = "185f8db32271fe25f561a6fc938b2e264306ec304eda518007d1764826381969"
const privateKey = "30770201010420efda9719cf88346df34dc48c7b88befcf9399cd2b08e6cf398390dc74bec0297a00a06082a8648ce3d030107a14403420004dba86a5a6472e12b5a39ce409a552d4c5858363256f2dd5bb834698a03f9c06db997c1c202ff028ce48965f45985594d0d2098a73bcefefde04cb70ad2c54ba3"
const signature = "81fedf5207ca7fe8d74f944557a1bbd69866c77c18c138335fa1e1ab53ce0209521f083d9ee78321e80a92faf1d9d4082fd3d1fc83a71be4ab1c356a28749190"

func Start() {
	privateKeyAsBytes, err := hex.DecodeString(privateKey)
	utils.HandleErr(err)

	restoredKey, err := x509.ParseECPrivateKey(privateKeyAsBytes)
	utils.HandleErr(err)

	signatureAsBytes, err := hex.DecodeString(signature)
	utils.HandleErr(err)

	rBytes := signatureAsBytes[:len(signatureAsBytes)/2]
	sBytes := signatureAsBytes[len(signatureAsBytes)/2:]

	var bigR, bigS = big.Int{}, big.Int{}

	bigR.SetBytes(rBytes)
	bigS.SetBytes(sBytes)

	hashAsBytes, err := hex.DecodeString(hashedMessage)
	utils.HandleErr(err)

	ok := ecdsa.Verify(&restoredKey.PublicKey, hashAsBytes, &bigR, &bigS)
	fmt.Println(ok)
}
