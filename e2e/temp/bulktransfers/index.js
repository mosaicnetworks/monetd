const path = require("path");
const fetch = require("node-fetch");
const argv = require('minimist')(process.argv.slice(2));

// import required objects
const { default:Node, Contract, Account } = require("evm-lite-core");

const networkpath = "/home/jon/.giverny/networks/funded2/"
const transdir = networkpath + "trans/"
const faucetfile = transdir + "faucet.json"

var fs = require("fs");

console.log("\n Loading Faucet \n");
var content = fs.readFileSync(faucetfile);

try {
    const data = JSON.parse(content)
    console.log(data);
  } catch(err) {
    console.error(err);
  }

// Load faucet JSON to set initial values - and populate the accounts array

