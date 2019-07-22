// Package poa provides standard function for create and managing POA networks.
//
// There is no functionality in package poa itself, it all resides in subpackages.
//
// A typical monetd configuration looks like this:
//
//  /home/user/.monet
//  ├── babble
//  │   ├── peers.genesis.json
//  │   ├── peers.json
//  │   └── priv_key
//  ├── eth
//  │   ├── genesis.json
//  │   ├── keystore
//  │   │   └── node0.json
//  │   └── pwd.txt
//  ├── monetd.toml
//  └── poa
//      ├── compile.toml
//      ├── contract0.abi
//      └── contract0.sol
//
// The aim of this package is to provide tools that allow you to create all of
// these files.
//
//
// Location
//
// The first stage is determining the location of the configuration directory.
// There is a generic configuration location function where you can pass the
// name of the configuration directory:
//  func DefaultConfigDir(configDir string) (string, error)
//
// But for most typical cases we will be using:
//  func DefaultMonetConfigDir() (string, error)
// DefaultMonetConfigDir wraps DefaultConfigDir and handles the different folder
// names on different OS.
//
// Directories
//
// These function ensure that all the expected directories are in place. This
// allows a sledgehammer approach to directories in other commands - i.e. run
// these commands and be assured that all the relevant items exist.
//
//  func CreateMonetConfigFolders(configDir string) error
//
// CreateMonetConfigFolders creates the default folder structure in configDir. It is a simple wrapper for
// CreateDirsIfNotExists.
//
// Keystore
//
// Keys are generated into the keystore subfolder of the configuration folder.
// The core function is GenerateKeyPair, whcih  generates an Ethereum key pair. keyfilepath is the path to write the new keyfile to.
// passwordFile is a plain text file containing the passphrase to use for the keyfile. privateKeyfile is the
// path to a private key. If specified, this function does not generate a new keyfile, it instead
// generates a keyfile from the private key. outputJSON controls whether the output to stdio is in
// JSON format or not.  The function returns a key object which can be used to retrive public or private
// keys or the address.
//
//  func GenerateKeyPair(keyfilepath, passwordFile, privateKeyfile string, outputJSON bool) (*keystore.Key, error)
//
// It is not envisaged that this function would be directly used in the Monet architecture. One of the
// wrappers detailed below would be more appropriate as they also handle placing the key files
// in the right configs folders.
//
//  func NewKeyPair(configDir, moniker, passwordFile string) (*keystore.Key, error)
//
// NewKeyPair is a wrapper to GenerateKeyPair. It does not support setting a private key.
// Additionally it does not support outputting to JSON format - if required, that can be
// achieved calling GenerateKeyPair directly.
//
// To inspect keys we have a pair of functions. InspectKey expects a full filepath to the actual keyfile. Normally
// you would use InspectKeyMoniker which looks for the key for a given moniker.
//
//  func InspectKeyMoniker(configDir string, moniker string,  PasswordFile string, showPrivate bool, outputJSON bool) error
//  func InspectKey(keyfilepath string, PasswordFile string, showPrivate bool, outputJSON bool) error
//
// We have a similar pair of functions for updating keys too:
//
//  func UpdateKeysMoniker(configDir string, moniker string, PasswordFile string, newPasswordFile string) error
//  func UpdateKeys(keyfilepath string, PasswordFile string, newPasswordFile string) error
//
// Monetd Toml File
//
// TODO.
//
// Peers
//
// TODO.
//
// Contract
//
// TODO.
//
// Genesis File
//
// TODO.
//
// Config Build
//
// TODO.
//
// Config Join
//
// TODO.
package poa
