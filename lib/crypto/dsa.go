package crypto

import (
	"crypto/dsa"
	"crypto/rand"
	"crypto/sha1"
	"io"
	"math/big"
)

var dsap = new(big.Int).SetBytes([]byte{
	0x9c, 0x05, 0xb2, 0xaa, 0x96, 0x0d, 0x9b, 0x97, 0xb8, 0x93, 0x19, 0x63, 0xc9, 0xcc, 0x9e, 0x8c,
	0x30, 0x26, 0xe9, 0xb8, 0xed, 0x92, 0xfa, 0xd0, 0xa6, 0x9c, 0xc8, 0x86, 0xd5, 0xbf, 0x80, 0x15,
	0xfc, 0xad, 0xae, 0x31, 0xa0, 0xad, 0x18, 0xfa, 0xb3, 0xf0, 0x1b, 0x00, 0xa3, 0x58, 0xde, 0x23,
	0x76, 0x55, 0xc4, 0x96, 0x4a, 0xfa, 0xa2, 0xb3, 0x37, 0xe9, 0x6a, 0xd3, 0x16, 0xb9, 0xfb, 0x1c,
	0xc5, 0x64, 0xb5, 0xae, 0xc5, 0xb6, 0x9a, 0x9f, 0xf6, 0xc3, 0xe4, 0x54, 0x87, 0x07, 0xfe, 0xf8,
	0x50, 0x3d, 0x91, 0xdd, 0x86, 0x02, 0xe8, 0x67, 0xe6, 0xd3, 0x5d, 0x22, 0x35, 0xc1, 0x86, 0x9c,
	0xe2, 0x47, 0x9c, 0x3b, 0x9d, 0x54, 0x01, 0xde, 0x04, 0xe0, 0x72, 0x7f, 0xb3, 0x3d, 0x65, 0x11,
	0x28, 0x5d, 0x4c, 0xf2, 0x95, 0x38, 0xd9, 0xe3, 0xb6, 0x05, 0x1f, 0x5b, 0x22, 0xcc, 0x1c, 0x93,
})

var dsaq = new(big.Int).SetBytes([]byte{
	0xa5, 0xdf, 0xc2, 0x8f, 0xef, 0x4c, 0xa1, 0xe2, 0x86, 0x74, 0x4c, 0xd8, 0xee, 0xd9, 0xd2, 0x9d,
	0x68, 0x40, 0x46, 0xb7,
})

var dsag = new(big.Int).SetBytes([]byte{
	0x0c, 0x1f, 0x4d, 0x27, 0xd4, 0x00, 0x93, 0xb4, 0x29, 0xe9, 0x62, 0xd7, 0x22, 0x38, 0x24, 0xe0,
	0xbb, 0xc4, 0x7e, 0x7c, 0x83, 0x2a, 0x39, 0x23, 0x6f, 0xc6, 0x83, 0xaf, 0x84, 0x88, 0x95, 0x81,
	0x07, 0x5f, 0xf9, 0x08, 0x2e, 0xd3, 0x23, 0x53, 0xd4, 0x37, 0x4d, 0x73, 0x01, 0xcd, 0xa1, 0xd2,
	0x3c, 0x43, 0x1f, 0x46, 0x98, 0x59, 0x9d, 0xda, 0x02, 0x45, 0x18, 0x24, 0xff, 0x36, 0x97, 0x52,
	0x59, 0x36, 0x47, 0xcc, 0x3d, 0xdc, 0x19, 0x7d, 0xe9, 0x85, 0xe4, 0x3d, 0x13, 0x6c, 0xdc, 0xfc,
	0x6b, 0xd5, 0x40, 0x9c, 0xd2, 0xf4, 0x50, 0x82, 0x11, 0x42, 0xa5, 0xe6, 0xf8, 0xeb, 0x1c, 0x3a,
	0xb5, 0xd0, 0x48, 0x4b, 0x81, 0x29, 0xfc, 0xf1, 0x7b, 0xce, 0x4f, 0x7f, 0x33, 0x32, 0x1c, 0x3c,
	0xb3, 0xdb, 0xb1, 0x4a, 0x90, 0x5e, 0x7b, 0x2b, 0x3e, 0x93, 0xbe, 0x47, 0x08, 0xcb, 0xcc, 0x82,
})

var param = dsa.Parameters{
	P: dsap,
	Q: dsaq,
	G: dsag,
}

// generate a dsa keypair
func DSAGenerate(priv *dsa.PrivateKey, rand io.Reader) error {
	// put our paramters in
	priv.P = param.P
	priv.Q = param.Q
	priv.G = param.G
	// generate the keypair
	return dsa.GenerateKey(priv, rand)
}

// create i2p dsa public key given its public component
func createDSAPublicKey(Y *big.Int) *dsa.PublicKey {
	return &dsa.PublicKey{
		Parameters: param,
		Y:          Y,
	}
}

// createa i2p dsa private key given its public component
func createDSAPrivkey(X *big.Int) *dsa.PrivateKey {
	Y := new(big.Int)
	Y.Exp(dsag, X, dsap)
	return &dsa.PrivateKey{
		PublicKey: dsa.PublicKey{
			Parameters: param,
			Y:          Y,
		},
		X: X,
	}
}

type DSAVerifier struct {
	k *dsa.PublicKey
}

type DSAPublicKey [128]byte

// create a new dsa verifier
func (k DSAPublicKey) NewVerifier() (v Verifier, err error) {
	v = &DSAVerifier{
		k: createDSAPublicKey(new(big.Int).SetBytes(k[:])),
	}
	return
}

// verify data with a dsa public key
func (v *DSAVerifier) Verify(data, sig []byte) (err error) {
	h := sha1.Sum(data)
	err = v.VerifyHash(h[:], sig)
	return
}

// verify hash of data with a dsa public key
func (v *DSAVerifier) VerifyHash(h, sig []byte) (err error) {
	if len(sig) == 40 {
		r := new(big.Int).SetBytes(sig[:20])
		s := new(big.Int).SetBytes(sig[20:])
		if dsa.Verify(v.k, h, r, s) {
			// valid signature
		} else {
			// invalid signature
			err = ErrInvalidSignature
		}
	} else {
		err = ErrBadSignatureSize
	}
	return
}

func (k DSAPublicKey) Len() int {
	return len(k)
}

type DSASigner struct {
	k *dsa.PrivateKey
}

type DSAPrivateKey [20]byte

// create a new dsa signer
func (k DSAPrivateKey) NewSigner() (s Signer, err error) {
	s = &DSASigner{
		k: createDSAPrivkey(new(big.Int).SetBytes(k[:])),
	}
	return
}

func (ds *DSASigner) Sign(data []byte) (sig []byte, err error) {
	h := sha1.Sum(data)
	sig, err = ds.SignHash(h[:])
	return
}

func (ds *DSASigner) SignHash(h []byte) (sig []byte, err error) {
	var r, s *big.Int
	r, s, err = dsa.Sign(rand.Reader, ds.k, h)
	if err == nil {
		sig = make([]byte, 40)
		copy(sig, r.Bytes())
		copy(sig[20:], s.Bytes())
	}
	return
}

func (k DSAPrivateKey) Len() int {
	return len(k)
}
