package utils

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"hash"
	"sync"
)

var _ Signer = (*Signature)(nil)

type (
	//Signature model
	Signature struct {
		id       string
		hashFunc hash.Hash
		alg      string
		lock     sync.Mutex
	}
	//Signer interface
	Signer interface {
		ID() string
		Algorithm() string
		Create(b []byte) []byte
		CreateString(b []byte) string
		Validate(b []byte, ex string) bool
	}
)

//NewSignature create sign sha256
func NewSHA256(id, secret string) *Signature {
	return NewCustomSignature(id, secret, "hmac-sha256", sha256.New)
}

//NewMD5 create sign md5
func NewMD5(id, secret string) *Signature {
	return NewCustomSignature(id, secret, "hmac-md5", md5.New)
}

//NewSHA512 create sign sha512
func NewSHA512(id, secret string) *Signature {
	return NewCustomSignature(id, secret, "hmac-sha512", sha512.New)
}

//NewCustomSignature create sign with custom hash function
func NewCustomSignature(id, secret, alg string, h func() hash.Hash) *Signature {
	return &Signature{
		id:       id,
		alg:      alg,
		hashFunc: hmac.New(h, []byte(secret)),
	}
}

//ID signature
func (s *Signature) ID() string {
	return s.id
}

//Algorithm getter
func (s *Signature) Algorithm() string {
	return s.alg
}

//Create getting hash as bytes
func (s *Signature) Create(b []byte) []byte {
	s.lock.Lock()
	defer func() {
		s.hashFunc.Reset()
		s.lock.Unlock()
	}()
	s.hashFunc.Write(b)
	return s.hashFunc.Sum(nil)
}

//CreateString getting hash as string
func (s *Signature) CreateString(b []byte) string {
	return hex.EncodeToString(s.Create(b))
}

//Validate signature
func (s *Signature) Validate(b []byte, ex string) bool {
	v, err := hex.DecodeString(ex)
	if err != nil {
		return false
	}
	return hmac.Equal(s.Create(b), v)
}

//SignatureStorage ------------------------------------------------------------------

//SignatureStorage storage
type SignatureStorage struct {
	list map[string]Signer
	lock sync.RWMutex
}

//NewSignatureStorage init storage
func NewSignatureStorage() *SignatureStorage {
	return &SignatureStorage{
		list: make(map[string]Signer),
	}
}

//Add adding to storage
func (ss *SignatureStorage) Add(s Signer) {
	ss.lock.Lock()
	defer ss.lock.Unlock()

	ss.list[s.ID()] = s
}

//Get getting to storage
func (ss *SignatureStorage) Get(id string) Signer {
	ss.lock.RLock()
	defer ss.lock.RUnlock()

	if s, ok := ss.list[id]; ok {
		return s
	}
	return nil
}

//Count count sign in storage
func (ss *SignatureStorage) Count() int {
	ss.lock.RLock()
	defer ss.lock.RUnlock()
	l := len(ss.list)
	return l
}

//Del deleting from storage
func (ss *SignatureStorage) Del(id string) {
	ss.lock.Lock()
	defer ss.lock.Unlock()

	delete(ss.list, id)
}

//CleanAll removing all from storage
func (ss *SignatureStorage) CleanAll() {
	ss.lock.Lock()
	defer ss.lock.Unlock()

	for k := range ss.list {
		delete(ss.list, k)
	}
}
