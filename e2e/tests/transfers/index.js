const path = require("path");
const fetch = require("node-fetch");
const argv = require('minimist')(process.argv.slice(2));


// import required objects
const { default:Node, Contract, Account } = require("evm-lite-core");

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


function DemoNode(name, host, port) {
	this.name = name;

	this.api = new Node(host, port);

	this.account = {};
}


const node0 = new DemoNode("node0", serverAddress, serverPort);
const node1 = new DemoNode("node1", serverAddress1, serverPort);
const node2 = new DemoNode("node2", serverAddress2, serverPort);
const node3 = new DemoNode("node3", serverAddress3, serverPort);


// import keystore and datadirectory objects
const { default:Keystore } = require("evm-lite-keystore");
const { default:DataDirectory } = require("evm-lite-datadir");



// Prebaked transaction lists. 
// The first parameter is the node the payment is being made to. The second is the amount.

// Payments from Node 0
const trans0=[[2,1620],[1,1890],[3,2160],[3,2430],[3,2700],[1,1080],[1,1350],[3,1620],[1,1890],[2,2160],[2,2430],[2,2700],[2,1080],[3,1350],[1,1620],[3,1890],[3,2160],[2,2430],[1,2700],[1,1080],[1,1350],[3,1620],[2,1890],[2,2160],[2,2430],[2,2700],[2,1080],[2,1350],[1,1620],[3,1890],[2,2160],[2,2430],[2,2700],[2,1080],[3,1350],[1,1620],[3,1890],[3,2160],[3,2430],[3,2700],[1,1080],[2,1350],[3,1620],[3,1890],[1,2160],[2,2430],[2,2700],[1,1080],[2,1350],[1,1620],[3,1890],[2,2160],[1,2430],[1,2700],[2,1080],[2,1350],[2,1620],[2,1890],[1,2160],[1,2430],[3,2700],[3,1080],[3,1350],[3,1620],[1,1890],[2,2160],[2,2430],[3,2700],[2,1080],[3,1350],[2,1620],[2,1890],[1,2160],[3,2430],[2,2700],[2,1080],[1,1350],[1,1620],[1,1890],[2,2160],[1,2430],[2,2700],[2,1080],[1,1350],[2,1620],[2,1890],[3,2160],[1,2430],[1,2700],[3,1080],[1,1350],[2,1620],[3,1890],[2,2160],[1,2430],[3,2700],[2,1080],[2,1350],[2,1620],[1,1890]]

// Payments from Node 1
const trans1=[[2,1620],[2,1890],[0,2160],[3,2430],[0,2700],[0,1080],[3,1350],[3,1620],[0,1890],[0,2160],[2,2430],[3,2700],[3,1080],[2,1350],[3,1620],[2,1890],[0,2160],[2,2430],[0,2700],[0,1080],[0,1350],[2,1620],[3,1890],[2,2160],[3,2430],[0,2700],[2,1080],[2,1350],[3,1620],[2,1890],[3,2160],[2,2430],[3,2700],[0,1080],[2,1350],[2,1620],[3,1890],[3,2160],[3,2430],[2,2700],[2,1080],[3,1350],[0,1620],[0,1890],[3,2160],[0,2430],[2,2700],[0,1080],[2,1350],[0,1620],[0,1890],[3,2160],[2,2430],[3,2700],[2,1080],[0,1350],[3,1620],[3,1890],[0,2160],[0,2430],[3,2700],[0,1080],[3,1350],[0,1620],[3,1890],[2,2160],[3,2430],[0,2700],[2,1080],[0,1350],[2,1620],[2,1890],[0,2160],[0,2430],[3,2700],[0,1080],[0,1350],[2,1620],[0,1890],[2,2160],[2,2430],[0,2700],[0,1080],[0,1350],[2,1620],[2,1890],[0,2160],[0,2430],[0,2700],[0,1080],[2,1350],[0,1620],[3,1890],[2,2160],[0,2430],[2,2700],[2,1080],[0,1350],[3,1620],[0,1890]]

// Payments from Node 2
const trans2=[[0,1620],[3,1890],[1,2160],[1,2430],[1,2700],[1,1080],[0,1350],[0,1620],[1,1890],[0,2160],[3,2430],[1,2700],[0,1080],[0,1350],[1,1620],[0,1890],[3,2160],[0,2430],[1,2700],[3,1080],[0,1350],[3,1620],[0,1890],[3,2160],[1,2430],[1,2700],[1,1080],[0,1350],[0,1620],[0,1890],[0,2160],[3,2430],[1,2700],[1,1080],[1,1350],[1,1620],[3,1890],[1,2160],[3,2430],[0,2700],[0,1080],[3,1350],[0,1620],[0,1890],[1,2160],[1,2430],[0,2700],[1,1080],[0,1350],[3,1620],[0,1890],[1,2160],[3,2430],[0,2700],[0,1080],[0,1350],[1,1620],[3,1890],[1,2160],[3,2430],[0,2700],[3,1080],[0,1350],[3,1620],[3,1890],[1,2160],[0,2430],[3,2700],[1,1080],[1,1350],[1,1620],[0,1890],[0,2160],[0,2430],[1,2700],[0,1080],[0,1350],[3,1620],[1,1890],[0,2160],[3,2430],[1,2700],[3,1080],[0,1350],[0,1620],[3,1890],[3,2160],[3,2430],[1,2700],[3,1080],[3,1350],[1,1620],[3,1890],[1,2160],[3,2430],[3,2700],[3,1080],[0,1350],[1,1620],[1,1890]]

// Payments from Node 3
const trans3=[[2,1620],[1,1890],[0,2160],[2,2430],[1,2700],[1,1080],[1,1350],[1,1620],[1,1890],[1,2160],[1,2430],[0,2700],[2,1080],[2,1350],[0,1620],[0,1890],[1,2160],[2,2430],[2,2700],[2,1080],[0,1350],[1,1620],[1,1890],[0,2160],[0,2430],[2,2700],[1,1080],[2,1350],[2,1620],[2,1890],[1,2160],[2,2430],[2,2700],[0,1080],[1,1350],[0,1620],[0,1890],[0,2160],[2,2430],[2,2700],[1,1080],[0,1350],[2,1620],[2,1890],[0,2160],[1,2430],[2,2700],[1,1080],[1,1350],[1,1620],[0,1890],[2,2160],[0,2430],[2,2700],[0,1080],[2,1350],[0,1620],[2,1890],[1,2160],[0,2430],[2,2700],[1,1080],[0,1350],[0,1620],[2,1890],[1,2160],[0,2430],[1,2700],[2,1080],[2,1350],[0,1620],[1,1890],[0,2160],[1,2430],[2,2700],[2,1080],[0,1350],[1,1620],[0,1890],[1,2160],[0,2430],[2,2700],[1,1080],[0,1350],[0,1620],[1,1890],[0,2160],[2,2430],[1,2700],[1,1080],[0,1350],[1,1620],[1,1890],[2,2160],[0,2430],[0,2700],[2,1080],[1,1350],[2,1620],[0,1890]]



// Offsets. The total change from the initial balance for nodes 0 to 3. 
// This is of course predicated on each account having the same number of transactions from the list. 
const offsets=[[0,-1620,3240,-1620],[-1890,270,3240,-1620],[270,270,1080,-1620],[-2160,270,1080,810],[-2160,2970,-1620,810],[-2160,5130,-2700,-270],[-2160,6480,-4050,-270],[-2160,6480,-5670,1350],[-2160,10260,-7560,-540],[0,10260,-7560,-2700],[-2430,10260,-5130,-2700],[-2430,10260,-5130,-2700],[-2430,9180,-4050,-2700],[-2430,7830,-2700,-2700],[-2430,9450,-4320,-2700],[-540,7560,-4320,-2700],[-540,7560,-6480,-540],[-540,5130,-1620,-2970],[-540,7830,-1620,-5670],[-540,7830,-1620,-5670],[2160,7830,-2970,-7020],[540,7830,-2970,-5400],[540,7830,-2970,-5400],[540,5670,-810,-5400],[540,5670,-810,-5400],[540,5670,1890,-8100],[-540,6750,2970,-9180],[-540,5400,5670,-10530],[-540,5400,5670,-10530],[-540,3510,7560,-10530],[-540,3510,7560,-10530],[-2970,1080,12420,-10530],[-5670,1080,15120,-10530],[-4590,1080,15120,-11610],[-5940,2430,15120,-11610],[-5940,4050,15120,-13230],[-5940,2160,13230,-9450],[-5940,2160,11070,-7290],[-8370,-270,11070,-2430],[-8370,-2970,13770,-2430],[-8370,-1890,13770,-3510],[-8370,-3240,13770,-2160],[-6750,-4860,13770,-2160],[-4860,-6750,13770,-2160],[-4860,-4590,11610,-2160],[-4860,-2160,11610,-4590],[-4860,-4860,17010,-7290],[-4860,-2700,15930,-8370],[-4860,-2700,17280,-9720],[-4860,-1080,15660,-9720],[-1080,-2970,13770,-9720],[-3240,-2970,15930,-9720],[-3240,-2970,15930,-9720],[-3240,-2970,15930,-9720],[-2160,-4050,17010,-10800],[-810,-5400,18360,-12150],[-810,-5400,18360,-12150],[-2700,-7290,20250,-10260],[-2700,-2970,18090,-12420],[-270,-2970,15660,-12420],[-270,-5670,15660,-9720],[-270,-5670,14580,-8640],[1080,-7020,13230,-7290],[2700,-8640,11610,-5670],[810,-8640,11610,-3780],[-1350,-6480,13770,-5940],[1080,-8910,13770,-5940],[1080,-8910,11070,-3240],[0,-8910,13230,-4320],[0,-8910,13230,-4320],[0,-8910,14850,-5940],[0,-8910,16740,-7830],[4320,-8910,14580,-9990],[6750,-8910,12150,-9990],[4050,-8910,14850,-9990],[5130,-9990,15930,-11070],[7830,-9990,14580,-12420],[6210,-8370,14580,-12420],[8100,-6480,12690,-14310],[8100,-6480,14850,-16470],[8100,-6480,14850,-16470],[8100,-6480,17550,-19170],[8100,-6480,17550,-19170],[10800,-6480,16200,-20520],[12420,-8100,17820,-22140],[10530,-8100,19710,-22140],[12690,-10260,17550,-19980],[12690,-10260,17550,-19980],[12690,-4860,14850,-22680],[12690,-4860,13770,-21600],[12690,-4860,13770,-21600],[12690,-3240,13770,-23220],[10800,-3240,11880,-19440],[8640,-3240,16200,-21600],[11070,-3240,13770,-21600],[11070,-5940,13770,-18900],[9990,-7020,15930,-18900],[11340,-7020,15930,-20250],[9720,-7020,17550,-20250],[11610,-5130,15660,-22140]]

// set the keystore object as the keystore for datadir object



// get account by address and decrypt with pass
// balance is inaccurate
const getAccount = async (address, password, datadir) => {
    // wait for keyfile to resolve
    const keyfile = await datadir.keystore.get(address);
  
    // return the decrypted account
    return Keystore.decrypt(keyfile, password);
  };



const checkFinalBalances = async () => {
    const datadirPath = argv.datadir;
    const numTrans = argv.transcount;
    const keystore = new Keystore(path.join(datadirPath, "keystore"));
    const datadir = new DataDirectory(datadirPath, "monetcli", keystore);
    console.log("Decrypting All Accounts")
    const account0 = await getAccount("node0", password, datadir);
    const account1 = await getAccount("node1", password, datadir);
    const account2 = await getAccount("node2", password, datadir);
    const account3 = await getAccount("node3", password, datadir);


    let ret0=await displayBalance(node0, account0.address, "node0");
    let ret1=await displayBalance(node1, account1.address, "node1");
    let ret2=await displayBalance(node2, account2.address, "node2");
    let ret3=await displayBalance(node3, account3.address, "node3");


    let diff0=ret0-1000000000
    let diff1=ret1-1000000000
    let diff2=ret2-1000000000
    let diff3=ret3-1000000000

    let exitcode = 0

    if (diff0 != offsets[numTrans-1][0]) {
        console.log("[FAIL] Node 0 expected " +offsets[numTrans-1][0]+ " got diff0");
        exitcode=600
    }

    if (diff1 != offsets[numTrans-1][1]) {
        console.log("[FAIL] Node 1 expected " +offsets[numTrans-1][1]+ " got diff1");
        exitcode=601
    }

    if (diff2 != offsets[numTrans-1][2]) {
        console.log("[FAIL] Node 2 expected " +offsets[numTrans-1][2]+ " got diff2");
        exitcode=602
    }

    if (diff3 != offsets[numTrans-1][3]) {
        console.log("[FAIL] Node 3 expected " +offsets[numTrans-1][3]+ " got diff3");
        exitcode=603
    }

    if (exitcode != 0) { process.exit(exitcode);}

    return "done"
}



const displayBalance = async (node, address, name) => {
  //      console.log(node)


		baseAccount = await node.api.getAccount(address);
        console.log(`${name}: `, '\n', baseAccount, '\n');
        return baseAccount.balance;
};


const getNonce = async (node, address) => {
    baseAccount = await node.api.getAccount(address);
    console.log(`Node: `, '\n', baseAccount, '\n');
    
    return baseAccount.nonce;
};






const transferRaw = async (node, from, to, value, nonce) => {
//	console.group('Locally Signed Transfer');

	const receipt = await node.api.transfer(
		from,
		to.address,
		value,
		defaultGas,
		defaultGasPrice
	);

	console.log('Receipt: ', receipt);
};




const transfers = async (nodeno) => {

    // initialize classes
    const datadirPath = argv.datadir;


    const numTrans = argv.transcount;

//    console.log(datadirPath);

    const keystore = new Keystore(path.join(datadirPath, "keystore"));    
    const datadir = new DataDirectory(datadirPath, "monetcli", keystore);
    

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


        if (i>=numTrans) {break;} 

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
   
     }
     
     console.log("complete")
  return "done"
}; 

console.log("Node number "+argv.nodeno)

if (argv.nodeno < 0) {
// This is the final check to verify changes are as expected.

     checkFinalBalances()
     .then(console.log)
     .catch(console.log)

} else {
// This is actual transfers. This script is run four times. Once for each account.

    transfers(argv.nodeno)
    .then(console.log)
    .catch(console.log);
}


