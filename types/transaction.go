package types

import (
	"crypto/sha256"

	"github.com/leehaowei/blocker/crypto"
	"github.com/leehaowei/blocker/proto"
	pb "google.golang.org/protobuf/proto"
)

func SignTransaction(pk *crypto.PrivateKey, tx *proto.Transaction) *crypto.Signature {
	return pk.Sing(HashTransaction(tx))
}

func HashTransaction(tx *proto.Transaction) []byte {
	b, err := pb.Marshal(tx)
	if err != nil {
		panic(err)
	}
	hash := sha256.Sum256(b)
	return hash[:]
}

func VerifyTransaction(tx *proto.Transaction) bool {
	for _, input := range tx.Input {
		if len(input.Signature) == 0 {
			panic("the ttransaction has no signature")
		}
		var (
			sig    = crypto.SignatureFromBytes(input.Signature)
			pubKey = crypto.PublicKeyFromBytes(input.PublicKey)
		)
		// TODO: to fix: ensure we don't run into problem after verifaction
		// because we've set the signature to nil
		tempSig := input.Signature
		input.Signature = nil
		if !sig.Verify(pubKey, HashTransaction(tx)) {
			return false
		}
		input.Signature = tempSig
	}
	return true
}
