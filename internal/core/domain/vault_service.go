package domain

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"reflect"
	"sort"

	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/tdex-network/tdex-daemon/pkg/wallet"
	"github.com/vulpemventures/go-elements/address"
	"github.com/vulpemventures/go-elements/network"
)

// IsZero returns whether the Vault is initialized without holding any data
func (v *Vault) IsZero() bool {
	return reflect.DeepEqual(*v, Vault{})
}

// IsInitialized returnes whether the Vault has been inizitialized
func (v *Vault) IsInitialized() bool {
	return len(v.EncryptedMnemonic) > 0
}

// IsLocked returns whether the Vault is initialized and locked
func (v *Vault) IsLocked() bool {
	return !v.IsInitialized() || !MnemonicStoreManager.IsSet()
}

// GetMnemonicSafe is getter for Vault's mnemonic in plain text
func (v *Vault) GetMnemonicSafe() ([]string, error) {
	if v.IsLocked() {
		return nil, ErrVaultMustBeUnlocked
	}

	return MnemonicStoreManager.Get(), nil
}

// Lock locks the Vault by wiping its mnemonic field
func (v *Vault) Lock() {
	if v.IsLocked() {
		return
	}
	// flush mnemonic in plain text
	MnemonicStoreManager.Unset()
}

// Unlock attempts to decrypt the mnemonic with the provided passphrase
func (v *Vault) Unlock(passphrase string) error {
	if !v.IsLocked() {
		return nil
	}

	mnemonic, err := EncrypterManager.Decrypt(v.EncryptedMnemonic, passphrase)
	if err != nil {
		return err
	}

	MnemonicStoreManager.Set(mnemonic)
	return nil
}

// ChangePassphrase attempts to unlock the
func (v *Vault) ChangePassphrase(currentPassphrase, newPassphrase string) error {
	if !v.IsLocked() {
		return ErrVaultMustBeLocked
	}
	if !v.isValidPassphrase(currentPassphrase) {
		return ErrVaultInvalidPassphrase
	}

	mnemonic, err := EncrypterManager.Decrypt(v.EncryptedMnemonic, currentPassphrase)
	if err != nil {
		return err
	}

	encryptedMnemonic, err := EncrypterManager.Encrypt(mnemonic, newPassphrase)
	if err != nil {
		return err
	}

	v.EncryptedMnemonic = encryptedMnemonic
	v.PassphraseHash = btcutil.Hash160([]byte(newPassphrase))
	return nil
}

// InitAccount creates a new account in the current Vault if not existing
func (v *Vault) InitAccount(accountIndex int) {
	if _, ok := v.Accounts[accountIndex]; !ok {
		account, _ := NewAccount(accountIndex)
		v.Accounts[accountIndex] = account
	}
}

// AccountByIndex returns the account with the given index
func (v *Vault) AccountByIndex(accountIndex int) (*Account, error) {
	account, ok := v.Accounts[accountIndex]
	if !ok {
		return nil, ErrVaultAccountNotFound
	}
	return account, nil
}

// AccountByAddress returns the account to which the provided address belongs
func (v *Vault) AccountByAddress(addr string) (*Account, int, error) {
	info, ok := v.AccountAndKeyByAddress[addr]
	if !ok {
		return nil, -1, ErrVaultAccountNotFound
	}
	account, err := v.AccountByIndex(info.AccountIndex)
	if err != nil {
		return nil, -1, err
	}
	return account, info.AccountIndex, nil
}

// DeriveNextExternalAddressForAccount returns the next unused address, the corresponding output script, the blinding key.
func (v *Vault) DeriveNextExternalAddressForAccount(accountIndex int) (*AddressInfo, error) {
	if v.IsLocked() {
		return nil, ErrVaultMustBeUnlocked
	}

	return v.deriveNextAddressForAccount(accountIndex, ExternalChain)
}

// DeriveNextInternalAddressForAccount returns the next unused change address for the
// provided account and the corresponding output script
func (v *Vault) DeriveNextInternalAddressForAccount(accountIndex int) (*AddressInfo, error) {
	if v.IsLocked() {
		return nil, ErrVaultMustBeUnlocked
	}

	return v.deriveNextAddressForAccount(accountIndex, InternalChain)
}

// AllDerivedAddressesInfo returns the info of all the external and internal
// addresses derived by the daemon. This method does not require the Vault to
// be unlocked since it does not make use of the mnemonic in plain text.
// The info returned for each address are the account index, the derivation
// path, and the private blinding key.
func (v *Vault) AllDerivedAddressesInfo() AddressesInfo {
	return v.allDerivedAddressesInfo()
}

// AllDerivedAddressesInfoForAccount returns info about all the external and
// internal addresses derived for the provided account.
func (v *Vault) AllDerivedAddressesInfoForAccount(
	accountIndex int,
) (AddressesInfo, error) {
	if v.IsLocked() {
		return nil, ErrVaultMustBeUnlocked
	}

	includeInternal := true
	return v.allDerivedAddressesInfoForAccount(accountIndex, includeInternal)
}

// AllDerivedExternalAddressesInfoForAccount returns info about all external
// addresses derived for the provided account.
func (v *Vault) AllDerivedExternalAddressesInfoForAccount(
	accountIndex int,
) (AddressesInfo, error) {
	if v.IsLocked() {
		return nil, ErrVaultMustBeUnlocked
	}

	includeInternal := true
	return v.allDerivedAddressesInfoForAccount(accountIndex, !includeInternal)
}

func (v *Vault) isValidPassphrase(passphrase string) bool {
	return bytes.Equal(v.PassphraseHash, btcutil.Hash160([]byte(passphrase)))
}

func (v *Vault) isPassphraseSet() bool {
	return len(v.PassphraseHash) > 0
}

func (v *Vault) deriveNextAddressForAccount(accountIndex, chainIndex int) (*AddressInfo, error) {
	w, err := wallet.NewWalletFromMnemonic(wallet.NewWalletFromMnemonicOpts{
		SigningMnemonic: MnemonicStoreManager.Get(),
	})
	if err != nil {
		return nil, err
	}

	account, ok := v.Accounts[accountIndex]
	if !ok {
		account, err = NewAccount(accountIndex)
		if err != nil {
			return nil, err
		}
		v.Accounts[accountIndex] = account
	}

	addressIndex := account.LastExternalIndex
	if chainIndex == InternalChain {
		addressIndex = account.LastInternalIndex
	}
	derivationPath := fmt.Sprintf(
		"%d'/%d/%d",
		account.AccountIndex, chainIndex, addressIndex,
	)
	addr, script, err := w.DeriveConfidentialAddress(wallet.DeriveConfidentialAddressOpts{
		DerivationPath: derivationPath,
		Network:        v.Network,
	})
	if err != nil {
		return nil, err
	}

	blindingKey, _, _ := w.DeriveBlindingKeyPair(wallet.DeriveBlindingKeyPairOpts{
		Script: script,
	})

	account.addDerivationPath(hex.EncodeToString(script), derivationPath)
	if chainIndex == InternalChain {
		account.nextInternalIndex()
	} else {
		account.nextExternalIndex()
	}
	v.AccountAndKeyByAddress[addr] = AccountAndKey{
		AccountIndex: account.AccountIndex,
		BlindingKey:  blindingKey.Serialize(),
	}

	return &AddressInfo{
		AccountIndex:   accountIndex,
		Address:        addr,
		Script:         hex.EncodeToString(script),
		BlindingKey:    blindingKey.Serialize(),
		DerivationPath: derivationPath,
	}, nil
}

func (v *Vault) allDerivedAddressesInfo() AddressesInfo {
	listLen := len(v.AccountAndKeyByAddress)
	list := make(AddressesInfo, listLen, listLen)
	i := 0

	for addr, info := range v.AccountAndKeyByAddress {
		account, _ := v.AccountByIndex(info.AccountIndex)
		script, _ := address.ToOutputScript(addr)
		path, _ := account.DerivationPathByScript[hex.EncodeToString(script)]

		list[i] = AddressInfo{
			AccountIndex:   info.AccountIndex,
			Address:        addr,
			BlindingKey:    info.BlindingKey,
			DerivationPath: path,
			Script:         hex.EncodeToString(script),
		}
		i++
	}

	// sorting list by derivation path also groups addresses by accountIndex
	sort.SliceStable(list, func(i, j int) bool {
		return list[i].DerivationPath < list[j].DerivationPath
	})
	return list
}

func (v *Vault) allDerivedAddressesInfoForAccount(
	accountIndex int,
	includeInternal bool,
) (AddressesInfo, error) {
	account, err := v.AccountByIndex(accountIndex)
	if err != nil {
		return nil, err
	}

	w, err := wallet.NewWalletFromMnemonic(wallet.NewWalletFromMnemonicOpts{
		SigningMnemonic: MnemonicStoreManager.Get(),
	})
	if err != nil {
		return nil, err
	}
	net := v.Network

	infoLen := account.LastExternalIndex + account.LastInternalIndex
	info := make(AddressesInfo, 0, infoLen)

	extAddressesInfo := deriveAddressesInRange(
		w,
		net,
		accountIndex,
		ExternalChain,
		0,
		account.LastExternalIndex-1,
	)
	info = append(info, extAddressesInfo...)

	if includeInternal {
		inAddressesInfo := deriveAddressesInRange(
			w,
			net,
			accountIndex,
			InternalChain,
			0,
			account.LastInternalIndex-1,
		)
		info = append(info, inAddressesInfo...)
	}

	return info, nil
}

func validateAccountIndex(accIndex int) error {
	if accIndex < 0 {
		return ErrInvalidAccount
	}

	return nil
}

// NextExternalIndex increments the last external index by one and returns the new last
func (a *Account) nextExternalIndex() (next int) {
	// restart from 0 if index has reached the its max value
	next = 0
	if a.LastExternalIndex != hdkeychain.HardenedKeyStart-1 {
		next = a.LastExternalIndex + 1
	}
	a.LastExternalIndex = next
	return
}

// NextInternalIndex increments the last internal index by one and returns the new last
func (a *Account) nextInternalIndex() (next int) {
	next = 0
	if a.LastInternalIndex != hdkeychain.HardenedKeyStart-1 {
		next = a.LastInternalIndex + 1
	}
	a.LastInternalIndex = next
	return
}

// AddDerivationPath adds an entry outputScript-derivationPath to the inner to
// the inner derivationPathByScript map
func (a *Account) addDerivationPath(outputScript, derivationPath string) {
	if _, ok := a.DerivationPathByScript[outputScript]; !ok {
		a.DerivationPathByScript[outputScript] = derivationPath
	}
}

func deriveAddressesInRange(
	w *wallet.Wallet,
	net *network.Network,
	accountIndex,
	chainIndex,
	firstAddressIndex,
	lastAddressIndex int,
) AddressesInfo {
	infoLen := lastAddressIndex - firstAddressIndex + 1
	info := make(AddressesInfo, infoLen, infoLen)

	for i := firstAddressIndex; i <= lastAddressIndex; i++ {
		derivationPath := fmt.Sprintf("%d'/%d/%d", accountIndex, chainIndex, i)
		addr, script, _ := w.DeriveConfidentialAddress(wallet.DeriveConfidentialAddressOpts{
			DerivationPath: derivationPath,
			Network:        net,
		})
		key, _, _ := w.DeriveBlindingKeyPair(wallet.DeriveBlindingKeyPairOpts{
			Script: script,
		})
		info[i] = AddressInfo{
			AccountIndex:   accountIndex,
			Address:        addr,
			BlindingKey:    key.Serialize(),
			DerivationPath: derivationPath,
			Script:         hex.EncodeToString(script),
		}
	}
	return info
}
