package iotago_test

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/iotaledger/hive.go/serializer/v2"
	"github.com/iotaledger/iota.go/v3"
	"github.com/iotaledger/iota.go/v3/tpkg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSignatureSelector(t *testing.T) {
	_, err := iotago.SignatureSelector(100)
	assert.True(t, errors.Is(err, iotago.ErrUnknownSignatureType))
}

func TestSignatures_DeSerialize(t *testing.T) {
	tests := []deSerializeTest{
		{
			name:   "ok - Ed25519",
			source: tpkg.RandEd25519Signature(),
			target: &iotago.Ed25519Signature{},
		},
		{
			name:   "ok - BLS",
			source: tpkg.RandBLSSignature(),
			target: &iotago.BLSSignature{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.deSerialize)
	}
}

func TestEd25519Signature_Valid(t *testing.T) {
	type test struct {
		Address   tpkg.HexBytes `json:"address"`
		Message   tpkg.HexBytes `json:"message"`
		PublicKey tpkg.HexBytes `json:"pub_key"`
		Signature tpkg.HexBytes `json:"signature"`
		Valid     bool          `json:"valid"`
	}
	var tests []test
	// load the tests from file
	b, err := ioutil.ReadFile(filepath.Join("testdata", t.Name()+".json"))
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(b, &tests))

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			// deserialize the address from the test
			addr := &iotago.Ed25519Address{}
			_, err = addr.Deserialize(tt.Address, serializer.DeSeriModePerformValidation, nil)
			require.NoError(t, err)
			// create the signature type
			sig := &iotago.Ed25519Signature{}
			copy(sig.PublicKey[:], tt.PublicKey)
			copy(sig.Signature[:], tt.Signature)

			sigError := sig.Valid(tt.Message, addr)
			switch tt.Valid {
			case true:
				assert.NoError(t, sigError)
			case false:
				assert.Error(t, sigError)
			}
		})
	}
}
