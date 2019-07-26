// Package contract supports the compilation of solidity contracts.
//
// If you wish to implement building a smart contract you would need to call the
// following functions:
//
//  func GetFinalSoliditySource(peers poatypes.PeerRecordList) (string, error)
// GetFinalSoliditySource contains the current POA contract embedded within it.
// It then applies the peerlist to populate the initial whitelist on the
// contract.
//  func CompileSolidityContract(soliditySource string)(map[string]*compiler.Contract, error)
// CompileSolidityContract takes the solidity source as a string, and returns
// a struct containing contract into, including bytecode and ABI.
package contract
