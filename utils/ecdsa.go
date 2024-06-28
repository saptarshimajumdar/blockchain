package utils

import (
	"fmt"
	"math/big"
)

type Signature struct {
	//represents the x co-ord of a point in the elliptic curve, computed by a random nuber during creating public and private keys
	R *big.Int
	// s is another part of the dig ital signature that also represents the result of operations involving the private key, the message hash, and the random number k.
	S *big.Int
}

func (s Signature) String() string {
	return fmt.Sprintf("%x%x", s.R, s.S)
}
