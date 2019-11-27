const path = require("path");
const fetch = require("node-fetch");
const argv = require('minimist')(process.argv.slice(2));


// import required objects
const { default:Node, Contract } = require("evm-lite-core");

// account address
const password = "test";
const defaultGas = 10000000;
const defaultGasPrice = 0;
const serverAddress0 = "172.77.5.10"
const serverAddress1 = "172.77.5.11"

const serverPort = "8080"

const node0 = new Node(serverAddress0, serverPort);
const node1 = new Node(serverAddress1, serverPort);

// import keystore and datadirectory objects
const { default:Keystore } = require("evm-lite-keystore");
const { default:DataDirectory } = require("evm-lite-datadir");

const sleep = (milliseconds) => {
    return new Promise(resolve => setTimeout(resolve, milliseconds))
  }

const checkWhitelist= async(contract, account, address) => { 
    const checkTrans = contract.methods.checkAuthorised({
        from: account.address, 
        gas: defaultGas,
        gasPrice: defaultGasPrice,
        value : 0,
    },address);

    const checkReceipt = await node0.callTx(checkTrans, account)

    return checkReceipt;
}

// submitEviction (address _nomineeAddress)
const submitEviction = async(contract, account, evictee) => {
 
    const checkTrans = contract.methods.submitEviction({
        from: account.address, 
        gas: defaultGas,
        gasPrice: defaultGasPrice,
        value : 0,
    }, evictee.address);

    const checkReceipt = await node0.sendTx(checkTrans, account)

    return checkReceipt;
}

const castEvictionVote = async(contract, account, address, vote) => {
 
    const checkTrans = contract.methods.castEvictionVote({
        from: account.address, 
        gas: defaultGas,
        gasPrice: defaultGasPrice,
        value : 0,
    }, address, vote);

    const checkReceipt = await node0.sendTx(checkTrans, account)
 
    return checkReceipt;
}

// get account by address and decrypt with pass
// balance is inaccurate
const getAccount = async (address, password, datadir) => {
    // wait for keyfile to resolve
    const keyfile = await datadir.keystore.get(address);
  
    // return the decrypted account
    return Keystore.decrypt(keyfile, password);
};

// evictNode1 uses the library methods to nominate and vote for the eviction of
// node1 from node0
const evictNode1 = async () => {

    // initialize classes
    const datadirPath = argv.datadir;

    console.log(datadirPath);

    const keystore = new Keystore(path.join(datadirPath, "keystore"));
    const datadir = new DataDirectory(datadirPath, "monetcli", keystore);
   
    console.log(keystore);

    // unlock all of the accounts
    console.log("Decrypting All Accounts")
    const account0 = await getAccount("node0", password, datadir);
    const account1 = await getAccount("node1", password, datadir);

    console.log("Getting POA Contract ABI")
    // Get Contract ABI    
    let url = "http://"+serverAddress0 + ":"+serverPort+"/poa"
    let res = await fetch(url);
    let json = await res.json();
    let abiObj = JSON.parse(json.abi)

    // Create Contract and Initialise it   
    const contract = Contract.load(abiObj, json.address)

    console.log("Checking Whitelist Status");
    // Check the Whitelist Status Status
    let rec0 = await checkWhitelist(contract, account0, account0.address);
    let rec1 = await checkWhitelist(contract, account0, account1.address);

    if ( !rec0 || !rec1 ) {
        console.log("Whitelist should be TT for nodes 0 and 1. Aborting")
        console.log("CheckAuthorised node 0: ", rec0);
        console.log("CheckAuthorised node 1: ", rec1);
        process.exit(181); 
    }

    console.log("Node0 evicts Node1");
    let recnom = await submitEviction(contract, account0, account1);
    let recvote = await castEvictionVote(contract, account0, account1.address, true);
    console.log(recvote);

    console.log("Checking node1 Whitelist");
    rec3 = await checkWhitelist(contract, account0, account1.address);
    if ( rec3 ) { 
        console.log("Expected node1 checkAuthorised to be false. Aborting.")
        process.exit(108);
    }
};

const checkEviction = async () => {
    // get node0 info and check info.peers is down to 1
    let url = "http://"+serverAddress0 + ":"+serverPort+"/info"
    let res = await fetch(url);
    let infoObj = await res.json();
    if (infoObj.num_peers != 1) {
        console.log("Node0.Info.peers should be 1, not ", infoObj.Peers)
        process.exit(108)
    }

    // get node1 info and check that info.status is Suspended
    url = "http://"+ serverAddress1 + ":"+serverPort+"/info"
    res = await fetch(url);
    infoObj = await res.json();
    if (infoObj.state != "Suspended") {
        console.log("Node1.Info.state should be Suspended, not ", infoObj.state)
        process.exit(108)
    }
};

evictNode1()
.then(() => sleep(2000))
.then(() => { return checkEviction();})
.catch(err => console.log(err));


