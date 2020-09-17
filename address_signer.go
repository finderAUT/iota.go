package iota

import (
	"crypto/ed25519"
	"errors"
	"fmt"
)

var (
	// Returned if the needed keys to sign a message are absent.
	ErrAddressKeysMissing = errors.New("keys for address missing")
)

// AddressSigner produces signatures for messages which get verified against a given address.
type AddressSigner interface {
	// Sign produces the signature for the given message.
	Sign(addr Serializable, msg []byte) (signature Serializable, err error)
}

// AddressSignerFunc implements the AddressSigner interface.
type AddressSignerFunc func(addr Serializable, msg []byte) (signature Serializable, err error)

func (s AddressSignerFunc) Sign(addr Serializable, msg []byte) (signature Serializable, err error) {
	return s(addr, msg)
}

// AddressKeys pairs an address and its source key(s).
type AddressKeys struct {
	// The target address.
	Address Serializable `json:"address"`
	// The signing keys.
	Keys interface{} `json:"keys"`
}

// NewInMemoryAddressSigner creates a new InMemoryAddressSigner holding the given AddressKeys.
func NewInMemoryAddressSigner(addrKeys ...AddressKeys) AddressSigner {
	ss := &InMemoryAddressSigner{
		addrKeys: map[Serializable]interface{}{},
	}
	for _, c := range addrKeys {
		ss.addrKeys[c.Address] = c.Keys
	}
	return ss
}

// InMemoryAddressSigner implements AddressSigner by holding keys simply in-memory.
type InMemoryAddressSigner struct {
	addrKeys map[Serializable]interface{}
}

func (s *InMemoryAddressSigner) Sign(addr Serializable, msg []byte) (signature Serializable, err error) {
	switch addr.(type) {
	case *WOTSAddress:
		// TODO: implement
		return nil, ErrWOTSNotImplemented
	case *Ed25519Address:
		prvKey, ok := s.addrKeys[addr].(ed25519.PrivateKey)
		if !ok {
			return nil, fmt.Errorf("%w: can't sign message for Ed25519 address", ErrAddressKeysMissing)
		}

		ed25519Sig := &Ed25519Signature{}
		copy(ed25519Sig.Signature[:], ed25519.Sign(prvKey, msg))
		copy(ed25519Sig.PublicKey[:], prvKey.Public().(ed25519.PublicKey))

		return ed25519Sig, nil
	default:
		return nil, fmt.Errorf("%w: unknown address type %T", ErrUnknownAddrType, addr)
	}
}
