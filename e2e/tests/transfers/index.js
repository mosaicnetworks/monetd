const path = require("path");
const fetch = require("node-fetch");
const argv = require('minimist')(process.argv.slice(2));


// import required objects
const { EVMLC, Contract, Account } = require("evm-lite-core");

// account address
const password = "test";
const defaultGas = 10000000;
const defaultGasPrice = 0;
const defaultTimeout = 10000;
const serverAddress = "172.77.5.10"
const serverAddress1 = "172.77.5.11"
const serverAddress2 = "172.77.5.12"
const serverAddress3 = "172.77.5.13"

const serverPort = "8080"

const sleep = (milliseconds) => {
    return new Promise(resolve => setTimeout(resolve, milliseconds))
  }


function Node(name, host, port) {
	this.name = name;

	this.api = new EVMLC(host, port);

	this.account = {};
}


const node0 = new Node("node0", serverAddress, serverPort);
const node1 = new Node("node1", serverAddress1, serverPort);
const node2 = new Node("node2", serverAddress2, serverPort);
const node3 = new Node("node3", serverAddress3, serverPort);


// import keystore and datadirectory objects
const { Keystore } = require("evm-lite-keystore");
const { DataDirectory } = require("evm-lite-datadir");


const trans0=[[2,1620],[1,1890],[3,2160],[3,2430],[3,2700],[1,1080],[1,1350],[3,1620],[1,1890],[2,2160],[2,2430],[2,2700],[2,1080],[3,1350],[1,1620],[3,1890],[3,2160],[2,2430],[1,2700],[1,1080],[1,1350],[3,1620],[2,1890],[2,2160],[2,2430],[2,2700],[2,1080],[2,1350],[1,1620],[3,1890],[2,2160],[2,2430],[2,2700],[2,1080],[3,1350],[1,1620],[3,1890],[3,2160],[3,2430],[3,2700],[1,1080],[2,1350],[3,1620],[3,1890],[1,2160],[2,2430],[2,2700],[1,1080],[2,1350],[1,1620],[3,1890],[2,2160],[1,2430],[1,2700],[2,1080],[2,1350],[2,1620],[2,1890],[1,2160],[1,2430],[3,2700],[3,1080],[3,1350],[3,1620],[1,1890],[2,2160],[2,2430],[3,2700],[2,1080],[3,1350],[2,1620],[2,1890],[1,2160],[3,2430],[2,2700],[2,1080],[1,1350],[1,1620],[1,1890],[2,2160],[1,2430],[2,2700],[2,1080],[1,1350],[2,1620],[2,1890],[3,2160],[1,2430],[1,2700],[3,1080],[1,1350],[2,1620],[3,1890],[2,2160],[1,2430],[3,2700],[2,1080],[2,1350],[2,1620],[1,1890]]

const trans1=[[2,1620],[2,1890],[0,2160],[3,2430],[0,2700],[0,1080],[3,1350],[3,1620],[0,1890],[0,2160],[2,2430],[3,2700],[3,1080],[2,1350],[3,1620],[2,1890],[0,2160],[2,2430],[0,2700],[0,1080],[0,1350],[2,1620],[3,1890],[2,2160],[3,2430],[0,2700],[2,1080],[2,1350],[3,1620],[2,1890],[3,2160],[2,2430],[3,2700],[0,1080],[2,1350],[2,1620],[3,1890],[3,2160],[3,2430],[2,2700],[2,1080],[3,1350],[0,1620],[0,1890],[3,2160],[0,2430],[2,2700],[0,1080],[2,1350],[0,1620],[0,1890],[3,2160],[2,2430],[3,2700],[2,1080],[0,1350],[3,1620],[3,1890],[0,2160],[0,2430],[3,2700],[0,1080],[3,1350],[0,1620],[3,1890],[2,2160],[3,2430],[0,2700],[2,1080],[0,1350],[2,1620],[2,1890],[0,2160],[0,2430],[3,2700],[0,1080],[0,1350],[2,1620],[0,1890],[2,2160],[2,2430],[0,2700],[0,1080],[0,1350],[2,1620],[2,1890],[0,2160],[0,2430],[0,2700],[0,1080],[2,1350],[0,1620],[3,1890],[2,2160],[0,2430],[2,2700],[2,1080],[0,1350],[3,1620],[0,1890]]

const trans2=[[0,1620],[3,1890],[1,2160],[1,2430],[1,2700],[1,1080],[0,1350],[0,1620],[1,1890],[0,2160],[3,2430],[1,2700],[0,1080],[0,1350],[1,1620],[0,1890],[3,2160],[0,2430],[1,2700],[3,1080],[0,1350],[3,1620],[0,1890],[3,2160],[1,2430],[1,2700],[1,1080],[0,1350],[0,1620],[0,1890],[0,2160],[3,2430],[1,2700],[1,1080],[1,1350],[1,1620],[3,1890],[1,2160],[3,2430],[0,2700],[0,1080],[3,1350],[0,1620],[0,1890],[1,2160],[1,2430],[0,2700],[1,1080],[0,1350],[3,1620],[0,1890],[1,2160],[3,2430],[0,2700],[0,1080],[0,1350],[1,1620],[3,1890],[1,2160],[3,2430],[0,2700],[3,1080],[0,1350],[3,1620],[3,1890],[1,2160],[0,2430],[3,2700],[1,1080],[1,1350],[1,1620],[0,1890],[0,2160],[0,2430],[1,2700],[0,1080],[0,1350],[3,1620],[1,1890],[0,2160],[3,2430],[1,2700],[3,1080],[0,1350],[0,1620],[3,1890],[3,2160],[3,2430],[1,2700],[3,1080],[3,1350],[1,1620],[3,1890],[1,2160],[3,2430],[3,2700],[3,1080],[0,1350],[1,1620],[1,1890]]

const trans3=[[2,1620],[1,1890],[0,2160],[2,2430],[1,2700],[1,1080],[1,1350],[1,1620],[1,1890],[1,2160],[1,2430],[0,2700],[2,1080],[2,1350],[0,1620],[0,1890],[1,2160],[2,2430],[2,2700],[2,1080],[0,1350],[1,1620],[1,1890],[0,2160],[0,2430],[2,2700],[1,1080],[2,1350],[2,1620],[2,1890],[1,2160],[2,2430],[2,2700],[0,1080],[1,1350],[0,1620],[0,1890],[0,2160],[2,2430],[2,2700],[1,1080],[0,1350],[2,1620],[2,1890],[0,2160],[1,2430],[2,2700],[1,1080],[1,1350],[1,1620],[0,1890],[2,2160],[0,2430],[2,2700],[0,1080],[2,1350],[0,1620],[2,1890],[1,2160],[0,2430],[2,2700],[1,1080],[0,1350],[0,1620],[2,1890],[1,2160],[0,2430],[1,2700],[2,1080],[2,1350],[0,1620],[1,1890],[0,2160],[1,2430],[2,2700],[2,1080],[0,1350],[1,1620],[0,1890],[1,2160],[0,2430],[2,2700],[1,1080],[0,1350],[0,1620],[1,1890],[0,2160],[2,2430],[1,2700],[1,1080],[0,1350],[1,1620],[1,1890],[2,2160],[0,2430],[0,2700],[2,1080],[1,1350],[2,1620],[0,1890]]

// set the keystore object as the keystore for datadir object



// get account by address and decrypt with pass
// balance is inaccurate
const getAccount = async (address, password, datadir) => {
    // wait for keyfile to resolve
    const keyfile = await datadir.keystore.get(address);
  
    // return the decrypted account
    return Keystore.decrypt(keyfile, password);
  };



const displayBalance = async (node, address, name) => {
  //      console.log(node)


		baseAccount = await node.api.getAccount(address);
        console.log(`${name}: `, '\n', baseAccount, '\n');
        return "";
};


const getNonce = async (node, address) => {
    baseAccount = await node.api.getAccount(address);
    console.log(`Node: `, '\n', baseAccount, '\n');
    
    return baseAccount.nonce;
};






const transferRaw = async (node, from, to, value, nonce) => {
//	console.group('Locally Signed Transfer');

	const transaction = Account.prepareTransfer(
		from.address,
		to.address,
		value,
		defaultGas,
		defaultGasPrice
	);

    transaction.nonce = nonce;
//	console.log('Transaction: ', transaction, '\n');

	const receipt = await node.api.sendTransaction(transaction, 
		from,
	 	defaultTimeout);
	
//	console.log('Receipt: ', receipt);

//	console.groupEnd();
};




const transfers = async (nodeno) => {

    // initialize classes
    const datadirPath = argv.datadir;

//    console.log(datadirPath);

    const datadir = new DataDirectory(datadirPath);
    const keystore = new Keystore(path.join(datadirPath, "keystore"));

//    console.log(keystore);


    datadir.setKeystore(keystore);


// unlock all of the accounts
    console.log("Decrypting All Accounts")
    const account0 = await getAccount("node0", password, datadir);
    const account1 = await getAccount("node1", password, datadir);
    const account2 = await getAccount("node2", password, datadir);
    const account3 = await getAccount("node3", password, datadir);


    var account;
    var trans;
    var node;


    switch (argv.nodeno) {
        case 0: 
            account=account0;
            trans=trans0;
            node=node0;
            break;
        case 1: 
            account=account1;
            trans=trans1;
            node=node1;
            break;
        case 2: 
            account=account2;
            trans=trans2;
            node=node2;
            break;
            
        case 3: 
            account=account3;
            trans=trans3;
            node=node3;
            break;
    
        default:
            console.log("Unknown node number")
            return
    }
    
//    console.log(account)

//    let ret =  await displayBalance(node0, account.address, "node0");

    let nonce = await getNonce(node, account.address);

  //  nonce+=1;
    console.log("Nonce is : "+nonce)


    for(let i = 0; i < trans.length; i++){

        let payee=trans[i][0];
        let amt=trans[i][1];
        let payeenode;
        let payeeaccount;

        switch (payee) {
            case 0:   payeenode = node0; payeeaccount = account0;  break;
            case 1:   payeenode = node1; payeeaccount = account1;  break;
            case 2:   payeenode = node2; payeeaccount = account2;  break;
            case 3:   payeenode = node3; payeeaccount = account3;  break;
            default:  
                return
        }


        let ret = await transferRaw(node,  account, payeeaccount, amt, nonce)

        nonce++;
        process.stdout.write(".");

        let slp = await sleep(100);
   
        if (i>100) {break;} 
   
     }
     
     console.log("complete")
  return "done"
}; 

console.log("Node number "+argv.nodeno)

transfers(argv.nodeno)
.then(console.log)
.catch(console.log);

