package wallet

import (
	"errors"
	"fmt"
)

var (
	// ErrNullNetwork ...
	ErrNullNetwork = errors.New("network must not be null")
	// ErrNullInputWitnessUtxo ...
	ErrNullInputWitnessUtxo = errors.New("input witness utxo must not be null")
	// ErrNullSigningMnemonic ...
	ErrNullSigningMnemonic = errors.New("signing mnemonic is null")
	// ErrNullBlindingMnemonic ...
	ErrNullBlindingMnemonic = errors.New("blinding mnemonic is null")
	// ErrNullSigningMasterKey ...
	ErrNullSigningMasterKey = errors.New("signing master key is null")
	// ErrNullBlindingMasterKey ...
	ErrNullBlindingMasterKey = errors.New("blinding master key is null")
	// ErrNullPassphrase ...
	ErrNullPassphrase = errors.New("passphrase must not be null")
	// ErrNullPlainText ...
	ErrNullPlainText = errors.New("text to encrypt must not be null")
	// ErrNullCypherText ...
	ErrNullCypherText = errors.New("cypher to decrypt must not be null")
	// ErrNullDerivationPath ...
	ErrNullDerivationPath = errors.New("derivation path must not be null")
	// ErrNullOutputDerivationPath ...
	ErrNullOutputDerivationPath = fmt.Errorf("output %v", ErrNullDerivationPath)
	// ErrNullChangeDerivationPath ...
	ErrNullChangeDerivationPath = fmt.Errorf("change %v", ErrNullDerivationPath)
	// ErrNullOutputScript ...
	ErrNullOutputScript = errors.New("output script must not be null")
	// ErrNullPset ...
	ErrNullPset = errors.New("pset base64 must not be null")
	// ErrNullChangePathsByAsset ...
	ErrNullChangePathsByAsset = errors.New(
		"derivation paths for eventual change(s) must not be null",
	)

	// ErrInvalidSigningMnemonic ...
	ErrInvalidSigningMnemonic = errors.New("signing mnemonic is invalid")
	// ErrInvalidEntropySize ...
	ErrInvalidEntropySize = errors.New(
		"entropy size must be a multiple of 32 in the range [128,256]",
	)
	// ErrInvalidBlindingMnemonic ...
	ErrInvalidBlindingMnemonic = errors.New("blinding mnemonic is invalid")
	// ErrInvalidCypherText ...
	ErrInvalidCypherText = errors.New("cypher must be in base64 format")
	// ErrInvalidDerivationPath ...
	ErrInvalidDerivationPath = errors.New("invalid derivation path")
	// ErrInvalidDerivationPathLength ...
	ErrInvalidDerivationPathLength = errors.New(
		"derivation path must be a relative path in the form \"account'/branch/index\"",
	)
	// ErrInvalidDerivationPathAccount ...
	ErrInvalidDerivationPathAccount = errors.New(
		"derivation path's account (first elem) must be hardened (suffix \"'\")",
	)
	// ErrInvalidInputAsset ...
	ErrInvalidInputAsset = errors.New("input asset must be a 32 byte array in hex format")
	// ErrInvalidOutputAsset ...
	ErrInvalidOutputAsset = errors.New("output asset must be a 32 byte array in hex format")
	// ErrInvalidOutputAddress ...
	ErrInvalidOutputAddress = errors.New("output address must be a valid address")
	// ErrInvalidChangeAddress ...
	ErrInvalidChangeAddress = errors.New("change address must be a valid address")
	// ErrInvalidMilliSatsPerBytes ...
	ErrInvalidMilliSatsPerBytes = errors.New("unit of mSats/byte must be at least 100 (0.1 sats/byte)")
	// ErrInvalidOutputBlindingKeysLen ...
	ErrInvalidOutputBlindingKeysLen = errors.New(
		"number of output blinding keys must match number of outputs",
	)
	// ErrInvalidPassphrase ...
	ErrInvalidPassphrase = errors.New("passphrase provided is not correct")
	// ErrInvalidSignatures ...
	ErrInvalidSignatures = errors.New("transaction contains invalid signature(s)")
	// ErrInvalidAttempts ...
	ErrInvalidAttempts = fmt.Errorf(
		"attempts must be a number in range [0, %d]",
		MaxBlindingAttempts,
	)

	// ErrEmptyDerivationPaths ...
	ErrEmptyDerivationPaths = errors.New("derivation path list must not be empty")
	// ErrEmptyUnspents ...
	ErrEmptyUnspents = errors.New("unspents list must not be empty")

	// ErrMalformedDerivationPath ...
	ErrMalformedDerivationPath = errors.New(
		"path must not start or end with a '/' and " +
			"can optionally start with 'm/' for absolute paths",
	)
	// ErrOutOfRangeDerivationPathAccount ...
	ErrOutOfRangeDerivationPathAccount = fmt.Errorf(
		"account index must be in hardened range [0, %d]",
		MaxHardenedValue,
	)
	// ErrZeroInputAmount ...
	ErrZeroInputAmount = errors.New("input amount must not be zero")
	// ErrZeroOutputAmount ...
	ErrZeroOutputAmount = errors.New("output amount must not be zero")

	// ErrReachedMaxBlindingAttempts ...
	ErrReachedMaxBlindingAttempts = errors.New(
		"max number of attempts reached for blinding the transaction",
	)
)
