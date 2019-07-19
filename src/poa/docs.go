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
// TODO.
//
// Keystore
//
// Keys are generated into the keystore subfolder of the configuration folder.
// TODO.
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
