// evm-lite-js imports
const { EVMLC, Account, Contract } = require('evm-lite-core');
const { Keystore } = require('evm-lite-keystore');

const util = require('util');
const path = require('path');
const fs = require('fs');
const argv = require('minimist')(process.argv.slice(2));
const solc = require('solc');
const prompt = require('prompt');

const FgRed = '\x1b[31m';
const FgGreen = '\x1b[32m';
const FgYellow = '\x1b[33m';
const FgBlue = '\x1b[34m';
const FgMagenta = '\x1b[35m';
const FgCyan = '\x1b[36m';
const FgWhite = '\x1b[37m';

const log = (color, text) => {
	console.log(color + text + '\x1b[0m');
};


const finalDelta = [-500,+500,0,0];

var online = true;

const defaultTimeout = 15000;

const schema = {
	properties: {
		enters: {
			description: 'PRESS ENTER TO CONTINUE',
			ask: function() {
				if (!online) {
					console.log('Skipping prompt');
					sleep(2);
				}
				return online;
			}
		}
	}
};

const step = message => {
	log(FgWhite, '\n' + message);
	return new Promise(resolve => {
		prompt.get(schema, function(err, res) {
			resolve();
		});
	});
};

const hardstep = message => {
	log(FgWhite, '\n' + message);
	return new Promise(resolve => {
		prompt.get('PRESS ENTER TO CONTINUE', function(err, res) {
			resolve();
		});
	});
};

const explain = message => {
	log(FgCyan, util.format('\nEXPLANATION:\n%s', message));
};

const space = () => {
	console.log('\n');
};

const sleep = function(time) {
	return new Promise(resolve => setTimeout(resolve, time));
};

/**
 * Demo starts here.
 */

/** DEFAULT CONFIG VALUES */
const DEFAULT_GAS = 100000000;
const DEFAULT_GASPRICE = 0;

function Node(name, host, port) {
	this.name = name;

	this.api = new EVMLC(host, port);

	this.account = {};
}

// Util functions
const readPasswordFile = path => {
	return fs.readFileSync(path, { encoding: 'utf8' });
};

// State functions
const allAccounts = [];
const allNodes = [];
const allMonikers = [];
const initBalance = [];
const lastBalance = [];

var crowdFunding = {};
var contractPath = '';

const init = async () => {
	console.group('Initialize Nodes: ');

	if (argv.offline) { online = false}

	const ips = argv.ips
		.replace(/\s/g, '')
		.split(',')
		.sort();
	const port = argv.port;

	const keystore = new Keystore(argv.keystore);

	const passwordPath = argv.pwd;
	const password = readPasswordFile(passwordPath);

	contractPath = argv.contract;

	console.log('Sorted IPs: ', ips);
	console.log('Port: ', port);
	console.log('Keystore Path: ', keystore.path);
	console.log('Password Path: ', passwordPath);
	console.log('Contract Path: ', contractPath);

	for (i = 0; i < ips.length; i++) {
		node = new Node(util.format('node%d', i + 1), ips[i], port);
		allNodes.push(node);
	}

	console.groupEnd();

	return {
		keystore,
		password
	};
};

const decryptAccounts = async ({ keystore, password }) => {
	console.group('Decrypt Accounts');
	console.log('Password: ', password);

	const keyfiles = await keystore.list();

    for (var moniker of Object.keys(keyfiles)) {

		keyfile = keyfiles[moniker];

		let account;

		try {
			account = await Keystore.decrypt(keyfile, password);
		} catch (e) {
			console.error(
				`Decryption Failed: ${keyfile.address} (${password})`
			);
		}

		try {
			if (account) {
				const base = await allNodes[0].api.getAccount(account.address);

				account.balance = base.balance;
				account.nonce = base.nonce;
				account.moniker = moniker;
			}
		} catch (e) {
			// pass
		}

		if (account) {
			let balance = 0;

			if (typeof account.balance === 'object') {
				balance = account.balance.toFormat(0);
			} else {
				balance = account.balance;
			}

			console.log('Decrypted: ', `${account.address} (${balance || 0})`);
			allAccounts.push(account);
			allMonikers.push(moniker);
			initBalance.push(balance)
		}
	}

	for (i = 0; i < allNodes.length; i++) {
		allNodes[i].api.defaultFrom = allAccounts[i].address;
		allNodes[i].account = allAccounts[i];
		allNodes[i].name =allMonikers[i];

	}

	console.groupEnd();
};


const checkBalances = async () => {
	console.group('Check Balances: ');
	i = 0; 
	let failed = false;
	for (const node of allNodes) {		
		let expected = initBalance[i] + finalDelta[i]
		if (expected == lastBalance[i]) { console.log(node.name + " balance as expected.");}
		else {
			console.log("ERROR: "+ node.name + " expected " + expected + " got " + lastBalance[i]);
			failed = true;
		}
		i++;
	}	
	console.groupEnd();
	if (failed) {process.exit(1)} 
};


const displayAllBalances = async () => {
	console.group('Current Account Balances');

	let i = 0;
	for (const node of allNodes) {
		const baseAccount = await node.api.getAccount(node.account.address);
		const account = {
			...baseAccount
		};

		let balance = 0;

		if (typeof account.balance === 'object') {
			balance = account.balance.toFormat(0);
		} else {
			balance = account.balance;
		}

		account.balance = balance;
		lastBalance[i] = balance;
		i++;
		console.log(`${node.name}: `, '\n', account, '\n');
	}
	console.groupEnd();
};

const transferRaw = async (from, to, value) => {
	console.group('Locally Signed Transfer');

	const transaction = Account.prepareTransfer(
		from.account.address,
		to.account.address,
		value,
		DEFAULT_GAS,
		DEFAULT_GASPRICE
	);

	console.log('Transaction: ', transaction, '\n');

	const receipt = await from.api.sendTransaction(transaction, 
		from.account,
	 	defaultTimeout);
	
	console.log('Receipt: ', receipt);

	console.groupEnd();
};

const compileContract = async () => {
	const input = fs.readFileSync(contractPath, {
		encoding: 'utf8'
	});
	const output = solc.compile(input.toString(), 1);
	const bytecode = output.contracts[`:CrowdFunding`].bytecode;
	const abi = output.contracts[`:CrowdFunding`].interface;

	console.log('ABI: ', abi, '\n');
	return Contract.create(JSON.parse(abi), bytecode);
};

class CrowdFunding {
	constructor(contract, node) {
		this.contract = contract;
		this.node = node;
	}

	async deploy(value) {
		const tx = this.contract.deployTransaction(
			[value],
			this.node.account.address,
			DEFAULT_GAS,
			DEFAULT_GASPRICE
		);

		const receipt = await this.node.api.sendTransaction(
			tx,
			this.node.account,
			defaultTimeout,
		);
		console.log('Receipt:', receipt);

		this.contract.setAddressAndAddFunctions(receipt.contractAddress);

		return this;
	}

	async contribute(value) {
		const tx = this.contract.methods.contribute({
			from: this.node.account.address,
			value,
			gas: DEFAULT_GAS,
			gasPrice: DEFAULT_GASPRICE
		});

		console.log('Transaction: ', tx, '\n');

		const receipt = await this.node.api.sendTransaction(
			tx,
			this.node.account,
			defaultTimeout,
		);

		for (const log of receipt.logs) {
			console.log(
				log.event || 'No Event Name',
				JSON.stringify(log.args, null, 2)
			);
		}

		return tx;
	}

	async checkGoalReached() {
		const tx = this.contract.methods.checkGoalReached({
			gas: DEFAULT_GAS,
			gasPrice: DEFAULT_GASPRICE
		});

		const response = await this.node.api.callTransaction(tx);

		const parsedResponse = {
			goalReached: response[0],
			beneficiary: response[1],
			fundingTarget: response[2].toFormat(0),
			current: response[3].toFormat(0)
		};

		log(FgBlue, JSON.stringify(parsedResponse, null, 2));

		return response;
	}

	async settle() {
		const tx = this.contract.methods.settle({
			from: this.node.account.address,
			gas: DEFAULT_GAS,
			gasPrice: DEFAULT_GASPRICE
		});

		console.log('Transaction: ', tx, '\n');

		const receipt = await this.node.api.sendTransaction(
			tx,
			this.node.account,
			defaultTimeout,
		);

		for (const log of receipt.logs) {
			console.log(
				log.event || 'No Event Name',
				JSON.stringify(log.args, null, 2)
			);
		}

		return tx;
	}
}

init()
	.then(object => decryptAccounts(object))
	.then(() => step('STEP 1) Get Accounts'))
	.then(() => {
		space();
		return displayAllBalances();
	})
	.then(() =>
		explain(
			'Each node controls one account which allows it to send and receive Ether. \n' +
				'The private keys reside locally and directly on the evm-lite nodes. In a \n' +
				'production setting, access to the nodes would be restricted to the people  \n' +
				'allowed to sign messages with the private key. We also keep a local copy \n' +
				'of all the private keys to demonstrate client-side signing.'
		)
	)
	.then(() => step('STEP 2) Send 500 tokens from Amelia to Becky'))
	.then(() => {
		space();
		return transferRaw(allNodes[0], allNodes[1], 500);
	})
	.then(() =>
		explain(
			'We created an EVM transaction to send 500 tokens from Amelia to Becky. The \n' +
				"transaction was signed localy with Amelia's private key and sent through Amelia's node. \n" +
				'The client-facing service running in EVM-Lite relayed the transaction to Babble \n' +
				'for consensus ordering. Babble gossiped the raw transaction to the other Babble \n' +
				'nodes which ran it through the consensus algorithm before committing it back to \n' +
				'EVM-Lite as part of Block. So each node received and processed the transaction. \n' +
				'They each applied the same changes to their local copy of the ledger.\n'
		)
	)
	.then(() => step('STEP 3) Check balances again'))
	.then(() => {
		space();
		return displayAllBalances();
	})
	.then(() =>
		explain('Notice how the balances of Amelia and Becky have changed.')
	)
	.then(() =>
		step(
			'STEP 4) Deploy a CrowdFunding SmartContract for 1000 tokens from Amelia'
		)
	)
	.then(() => {
		space();
		return compileContract();
	})
	.then(async contract => {
		crowdFunding = new CrowdFunding(contract, allNodes[0]);
		await crowdFunding.deploy(1000);
	})
	.then(() =>
		explain(
			'Here we compiled and deployed the CrowdFunding SmartContract. \n' +
				'The contract was written in the high-level Solidity language which compiles \n' +
				'down to EVM bytecode. To deploy the SmartContract we created an EVM transaction \n' +
				"with a 'data' field containing the bytecode. After going through consensus, the \n" +
				'transaction is applied on every node, so every participant will run a copy of \n' +
				'the same code with the same data.'
		)
	)
	.then(() => step('STEP 5) Contribute 499 tokens from Amelia'))
	.then(() => {
		space();
		return crowdFunding.contribute(499);
	})
	.then(() =>
		explain(
			"We created an EVM transaction to call the 'contribute' method of the SmartContract. \n" +
				"The 'value' field of the transaction is the amount that the caller is actually \n" +
				'going to contribute. The operation would fail if the account did not have enough Ether. \n' +
				'As an exercise you can check that the transaction was run through every Babble \n' +
				"node and that Becky's balance has changed."
		)
	)
	.then(() => step('STEP 6) Check goal reached'))
	.then(() => {
		space();
		return crowdFunding.checkGoalReached();
	})
	.then(() =>
		explain(
			'Here we called another method of the SmartContract to check if the funding goal \n' +
				'was met. Since only 499 of 1000 were received, the answer is no.'
		)
	)
	.then(() => step('STEP 7) Contribute 501 wei from Amelia again'))
	.then(() => {
		space();
		return crowdFunding.contribute(501);
	})
	.then(() => step('STEP 8) Check goal reached'))
	.then(() => {
		space();
		return crowdFunding.checkGoalReached();
	})
	.then(() =>
		step(
			'STEP 9) Before we `settle` lets check balances again to show that Amelia\'s balance decreased by a total of 1000.'
		)
	)
	.then(() => {
		space();
		return displayAllBalances();
	})
	.then(() =>
		explain('Since the funding goal was reached we can now settle.')
	)
	.then(() => step('STEP 10) Settle'))
	.then(() => {
		space();
		return crowdFunding.settle();
	})
	.then(() =>
		explain(
			'The funds were transferred from the SmartContract back to Amelia.'
		)
	)
	.then(() => step('STEP 11) Check balances again'))
	.then(() => {
		space();
		return displayAllBalances();
	})	
	.then(() => {
		space();
		return checkBalances();
	})
	.catch(err => console.log(err));
