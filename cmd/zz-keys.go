/*
Copyright © 2020 Alessandro Segala (@ItalyPaleAle)

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/

package cmd

import (
	"crypto/subtle"
	"errors"
	"strings"

	"github.com/ItalyPaleAle/prvt/crypto"
	"github.com/ItalyPaleAle/prvt/infofile"
	"github.com/ItalyPaleAle/prvt/keys"

	"github.com/manifoldco/promptui"
)

// PromptPassphrase prompts the user for a passphrase
func PromptPassphrase() (string, error) {
	prompt := promptui.Prompt{
		Validate: func(input string) error {
			if len(input) < 1 {
				return errors.New("Passphrase must not be empty")
			}
			return nil
		},
		Label: "Passphrase",
		Mask:  '*',
	}

	key, err := prompt.Run()
	if err != nil {
		return "", err
	}

	return key, err
}

// NewInfoFile generates a new info file with a brand-new master key, wrapped either with a passphrase-derived key, or with GPG
func NewInfoFile(gpgKey string) (info *infofile.InfoFile, errMessage string, err error) {
	// First, create the info file
	info, err = infofile.New()
	if err != nil {
		return nil, "Error creating info file", err
	}

	// Generate the master key
	masterKey, err := crypto.NewKey()
	if err != nil {
		return nil, "Error generating the master key", err
	}

	// Add the key
	errMessage, err = AddKey(info, masterKey, gpgKey)
	if err != nil {
		info = nil
	}

	return info, "", nil
}

// UpgradeInfoFile upgrades an info file to the latest version
func UpgradeInfoFile(info *infofile.InfoFile) (errMessage string, err error) {
	// Can only upgrade info files version 1 and 2
	if info.Version != 1 && info.Version != 2 {
		return "Unsupported repository version", errors.New("This repository has already been upgraded or is using an unsupported version")
	}

	// Upgrade 1 -> 2
	if info.Version < 2 {
		errMessage, err = upgradeInfoFileV1(info)
		if err != nil {
			return errMessage, err
		}
	}

	// Upgrade 2 -> 3
	// Nothing to do here, as the change is just in the index file
	// However, we still want to update the info file so older versions of prvt won't try to open a protobuf-encoded index file
	/*if info.Version < 3 {
	}*/

	// Update the version
	info.Version = 3

	return "", nil
}

// Upgrade an info file from version 1 to 2
func upgradeInfoFileV1(info *infofile.InfoFile) (errMessage string, err error) {
	// GPG keys are already migrated into the Keys slice
	// But passphrases need to be migrated
	if len(info.Salt) > 0 && len(info.ConfirmationHash) > 0 {
		// Prompt for the passphrase to get the current master key
		passphrase, err := PromptPassphrase()
		if err != nil {
			return "Error getting passphrase", err
		}

		// Get the current master key from the passphrase
		masterKey, confirmationHash, err := crypto.KeyFromPassphrase(passphrase, info.Salt)
		if err != nil || subtle.ConstantTimeCompare(info.ConfirmationHash, confirmationHash) == 0 {
			return "Cannot unlock the repository", errors.New("Invalid passphrase")
		}

		// Create a new salt
		newSalt, err := crypto.NewSalt()
		if err != nil {
			return "Error generating a new salt", err
		}

		// Create a new wrapping key
		wrappingKey, newConfirmationHash, err := crypto.KeyFromPassphrase(passphrase, newSalt)
		if err != nil {
			return "Error deriving the wrapping key", err
		}

		// Wrap the key
		wrappedKey, err := crypto.WrapKey(wrappingKey, masterKey)
		if err != nil {
			return "Error wrapping the master key", err
		}

		// Add the key
		err = info.AddPassphrase(newSalt, newConfirmationHash, wrappedKey)
		if err != nil {
			return "Error adding the key", err
		}

		// Remove the old key
		info.Salt = nil
		info.ConfirmationHash = nil
	}

	return "", nil
}

// AddKey adds a key to an info file
// If the GPG Key is empty, will prompt for a passphrase
func AddKey(info *infofile.InfoFile, masterKey []byte, gpgKey string) (errMessage string, err error) {
	if gpgKey == "" {
		// Add the passphrase
		return addKeyPassphrase(info, masterKey)
	} else {
		// Before adding the key, check if it's already there
		// Lowercase the key ID for normalization
		keyId := strings.ToLower(gpgKey)
		for _, k := range info.Keys {
			if strings.ToLower(k.GPGKey) == keyId {
				return "Key already added", errors.New("This GPG key has already been added to the repository")
			}
		}

		// Add the GPG key
		return addKeyGPG(info, masterKey, gpgKey)
	}
}

// Used by AddKey to add a new passphrase
func addKeyPassphrase(info *infofile.InfoFile, masterKey []byte) (errMessage string, err error) {
	var salt, confirmationHash, wrappedKey []byte

	// No GPG key specified, so we need to prompt for a passphrase first
	passphrase, err := PromptPassphrase()
	if err != nil {
		return "Error getting passphrase", err
	}

	// Before adding the key, check if it's already there
	_, _, _, testErr := keys.GetMasterKeyWithPassphrase(info, passphrase)
	if testErr == nil {
		return "Key already added", errors.New("This passphrase has already been added to the repository")
	}

	// Derive the wrapping key, after generating a new salt
	salt, err = crypto.NewSalt()
	if err != nil {
		return "Error generating a new salt", err
	}
	var wrappingKey []byte
	wrappingKey, confirmationHash, err = crypto.KeyFromPassphrase(passphrase, salt)
	if err != nil {
		return "Error deriving the wrapping key", err
	}

	// Wrap the key
	wrappedKey, err = crypto.WrapKey(wrappingKey, masterKey)
	if err != nil {
		return "Error wrapping the master key", err
	}

	// Add the key
	err = info.AddPassphrase(salt, confirmationHash, wrappedKey)
	if err != nil {
		return "Error adding the key", err
	}

	return "", nil
}

// Used by AddKey to add a new GPG key
func addKeyGPG(info *infofile.InfoFile, masterKey []byte, gpgKey string) (errMessage string, err error) {
	var wrappedKey []byte

	// Use GPG to wrap the master key
	wrappedKey, err = keys.GPGEncrypt(masterKey, gpgKey)
	if err != nil {
		return "Error encrypting the master key with GPG", err
	}

	// Add the key
	err = info.AddGPGWrappedKey(gpgKey, wrappedKey)
	if err != nil {
		return "Error adding the key", err
	}

	return "", nil
}

// GetMasterKey gets the master key, either unwrapping it with a passphrase or with GPG
func GetMasterKey(info *infofile.InfoFile) (masterKey []byte, keyId string, errMessage string, err error) {
	// First, try unwrapping the key using GPG
	masterKey, keyId, errMessage, err = keys.GetMasterKeyWithGPG(info)
	if err == nil {
		return
	}

	// No GPG key specified or unlocking with a GPG key was not successful
	// We'll try with passphrases; first, prompt for it
	passphrase, err := PromptPassphrase()
	if err != nil {
		return nil, "", "Error getting passphrase", err
	}

	// Try unwrapping using the passphrase
	return keys.GetMasterKeyWithPassphrase(info, passphrase)
}
