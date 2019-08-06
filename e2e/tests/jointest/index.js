const path = require("path");
const fetch = require("node-fetch");
const argv = require('minimist')(process.argv.slice(2));


// import required objects
const { EVMLC, Contract } = require("evm-lite-core");

// account address
const password = "test";
const defaultGas = 10000000;
const defaultGasPrice = 0;
const serverAddress = "172.77.5.10"
const serverAddress1 = "172.77.5.11"
const serverAddress2 = "172.77.5.12"
const serverAddress3 = "172.77.5.13"

const serverPort = "8080"

const node = new EVMLC(serverAddress, serverPort);
const node2 = new EVMLC(serverAddress1, serverPort);
const node3 = new EVMLC(serverAddress2, serverPort);
const node4 = new EVMLC(serverAddress3, serverPort);


// import keystore and datadirectory objects
const { Keystore } = require("evm-lite-keystore");
const { DataDirectory } = require("evm-lite-datadir");






// set the keystore object as the keystore for datadir object



const checkWhitelist= async(contract, account, address) => {
 
    const checkTrans = contract.methods.checkAuthorised({
        from: account.address, 
        gas: defaultGas,
        gasPrice: defaultGasPrice,
        value : 0,
    },address);

    const checkReceipt = await node.callTransaction(checkTrans, account)
 //   console.log(checkReceipt)

    return checkReceipt;
}


const getWhiteListCount = async(contract, account) => {
 
    const checkTrans = contract.methods.getWhiteListCount({
        from: account.address, 
        gas: defaultGas,
        gasPrice: defaultGasPrice,
        value : 0,
    });

    const checkReceipt = await node.callTransaction(checkTrans, account)
 //   console.log(checkReceipt)

    return checkReceipt.toNumber();
}


const getNomineeCount = async(contract, account) => {
 
    const checkTrans = contract.methods.getNomineeCount({
        from: account.address, 
        gas: defaultGas,
        gasPrice: defaultGasPrice,
        value : 0,
    });

    const checkReceipt = await node.callTransaction(checkTrans, account)
 //   console.log(checkReceipt)

    return checkReceipt.toNumber();
}




const selfnominate = async(contract, account) => {
 
    const checkTrans = contract.methods.submitNominee({
        from: account.address, 
        gas: defaultGas,
        gasPrice: defaultGasPrice,
        value : 0,
    }, account.address, account.name);

    const checkReceipt = await node.sendTransaction(checkTrans, account)
 //   console.log(checkReceipt)

    return checkReceipt;
}



const castvote = async(contract, account, address, vote) => {
 
    //castNomineeVote(address _nomineeAddress, bool _accepted)

    const checkTrans = contract.methods.castNomineeVote({
        from: account.address, 
        gas: defaultGas,
        gasPrice: defaultGasPrice,
        value : 0,
    }, address, vote);

    const checkReceipt = await node.sendTransaction(checkTrans, account)
 //   console.log(checkReceipt)

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



const run = async () => {

    // initialize classes
    const datadirPath = argv.datadir;

    console.log(datadirPath);

    const datadir = new DataDirectory(datadirPath);
    const keystore = new Keystore(path.join(datadirPath, "keystore"));

    console.log(keystore);


    datadir.setKeystore(keystore);


// unlock all of the accounts
    console.log("Decrypting All Accounts")
    const account0 = await getAccount("node0", password, datadir);
    const account1 = await getAccount("node1", password, datadir);
    const account2 = await getAccount("node2", password, datadir);
    const account3 = await getAccount("node3", password, datadir);

    console.log("Getting POA Contract ABI")
// Get Contract ABI    
    let url = "http://"+serverAddress + ":"+serverPort+"/poa"
    let res = await fetch(url);
    let json = await res.json();
    let abiObj = JSON.parse(json.abi)

// Create Contract and Initialise it   

    console.log("Running POA init")
    const contract = Contract.load(abiObj, json.address)
    const initTrans = contract.methods.init({
        from: account0.address, 
        gas: defaultGas,
        gasPrice: defaultGasPrice,
        value : 0,
    })

    const initReceipt = await node.sendTransaction(initTrans, account0);
  
 //   TODO uncomment this line   
 //   console.log(initReceipt);

     console.log("Checking Whitelist Status");
 // Check the Whitelist Status Status
    let rec0 = await checkWhitelist(contract, account0, account0.address);
    let rec1 = await checkWhitelist(contract, account0, account1.address);
    let rec2 = await checkWhitelist(contract, account0, account2.address);
    let rec3 = await checkWhitelist(contract, account0, account3.address);

    if ( ( ! rec0 ) || ( ! rec1 ) || ( ! rec2 ) || ( rec3 )) {
        console.log("Whitelist should be TTTF for nodes 0 to 3. Aborting")
        console.log("CheckAuthorised node 0: ", rec0);
        console.log("CheckAuthorised node 1: ", rec1);
        console.log("CheckAuthorised node 2: ", rec2);
        console.log("CheckAuthorised node 3: ", rec3);
        process.exit(101); 
    }

    console.log("Checking Whitelist count");
    let reccnt = await getWhiteListCount(contract, account0);

    if ( reccnt != 3 ) {
        console.log("Expected Whitelist count of 3, got "+ reccnt+". Aborting");
        process.exit(102);
    }

    console.log("Checking Nominee count");
    reccnt = await getNomineeCount(contract, account0);

    if ( reccnt != 0 ) {
        console.log("Expected Nominee count of 0, got "+ reccnt+". Aborting");
        process.exit(103);
    }


    console.log("Node3 self nominates");

    let recnom = await selfnominate(contract, account3);

 //   console.log(recnom);



    console.log("Checking Nominee count");
    reccnt = await getNomineeCount(contract, account0);

    if ( reccnt != 1 ) {
        console.log("Expected Nominee count of 1, got "+ reccnt+". Aborting");
        process.exit(105);
    }

    console.log("Checking node3 Whitelist");
    rec3 = await checkWhitelist(contract, account0, account3.address);
    if ( rec3 ) { 
        console.log("Expected node3 checkAuthorised to be false. Got true. Aborting.")
        process.exit(106);
    }


    console.log("Node 0 votes for node 3");
    let recvote = await castvote(contract, account0, account3.address, true);
    console.log(recvote);


    console.log("Checking Nominee count");
    reccnt = await getNomineeCount(contract, account0);

    if ( reccnt != 1 ) {
        console.log("Expected Nominee count of 1, got "+ reccnt+". Aborting");
        process.exit(107);
    }

    console.log("Checking node3 Whitelist");
    rec3 = await checkWhitelist(contract, account0, account3.address);
    if ( rec3 ) { 
        console.log("Expected node3 checkAuthorised to be false. Got true. Aborting.")
        process.exit(108);
    }



    console.log("Node 1 votes for node 3");
    recvote = await castvote(contract, account1, account3.address, true);
    console.log(recvote);

    console.log("Checking Nominee count");
    reccnt = await getNomineeCount(contract, account0);

    if ( reccnt != 1 ) {
        console.log("Expected Nominee count of 1, got "+ reccnt+". Aborting");
        process.exit(108);
    }

    console.log("Checking node3 Whitelist");
    rec3 = await checkWhitelist(contract, account0, account3.address);
    if ( rec3 ) { 
        console.log("Expected node3 checkAuthorised to be false. Got true. Aborting.")
        process.exit(109);
    }


    console.log("Node 2 votes for node 3");
    recvote = await castvote(contract, account2, account3.address, true);
    console.log(recvote);



    console.log("Checking Whitelist Status");
    // Check the Whitelist Status Status
       rec0 = await checkWhitelist(contract, account0, account0.address);
       rec1 = await checkWhitelist(contract, account0, account1.address);
       rec2 = await checkWhitelist(contract, account0, account2.address);
       rec3 = await checkWhitelist(contract, account0, account3.address);
   
       if ( ( ! rec0 ) || ( ! rec1 ) || ( ! rec2 ) || ( ! rec3 )) {
           console.log("Whitelist should be TTTT for nodes 0 to 3. Aborting")
           console.log("CheckAuthorised node 0: ", rec0);
           console.log("CheckAuthorised node 1: ", rec1);
           console.log("CheckAuthorised node 2: ", rec2);
           console.log("CheckAuthorised node 3: ", rec3);
           process.exit(111); 
       }
   
       console.log("Checking Whitelist count");
       reccnt = await getWhiteListCount(contract, account0);
   
       if ( reccnt != 4 ) {
           console.log("Expected Whitelist count of 4, got "+ reccnt+". Aborting");
           process.exit(112);
       }
   

       console.log("Checking Nominee count");
       reccnt = await getNomineeCount(contract, account0);
   
       if ( reccnt != 0 ) {
           console.log("Expected Nominee count of 0, got "+ reccnt+". Aborting");
           process.exit(103);
       }
   


    return "done"
};


run()
  .then(console.log)
  .catch(console.log);
