const path = require("path");
const fetch = require("node-fetch");
const argv = require('minimist')(process.argv.slice(2));

// import required objects
const { default:Node, Contract, Account, Transaction } = require("evm-lite-core");
const { default:Keystore } = require("evm-lite-keystore");
const { default:DataDirectory } = require("evm-lite-datadir");



const network = argv.network
const acct = argv.account
const total = argv.totals
const givdir = argv.givdir

const networkpath = givdir+"/networks/"+network+"/"

const transdir = networkpath + "trans/"
const faucetfile = transdir + acct + ".json"
const defaultGas = 10000000;
const defaultGasPrice = 0;


var fs = require("fs");

var NodeCollection = {};
var AccountCollection = {};



const keystorepath = path.join(networkpath, "keystore");

// console.log(keystorepath);

const password = "test";  //TODO live read
const keystore = new Keystore(keystorepath);
const datadir = new DataDirectory(networkpath, "monetcli", keystore);
 


const getAccount = async (address, password, datadir) => {
  // wait for keyfile to resolve

  console.log("Decrypting "+address+", "+password)
//  console.log(datadir.config)

  const keyfile = await datadir.keystore.get(address);

  // return the decrypted account
  return Keystore.decrypt(keyfile, password);

};



const InitAccount = async (name, password, datadir, nodename, address) => {
  if (! AccountCollection.hasOwnProperty(name)){
      AccountCollection[name] = await getAccount(name, password, datadir);

      baseAccount = await NodeCollection[nodename].api.getAccount(address);
      console.log(`GET ACCOUNT: `, '\n', baseAccount, '\n');

      AccountCollection[name].nonce = baseAccount.nonce+10;

  }
}


const InitNode = async (name, host, port) => {
   if (! NodeCollection.hasOwnProperty(name)){
       NodeCollection[name] = new DemoNode(name,host,port)
   }
}

function DemoNode(name, host, port) {
  console.log("Creating Node "+name);
	this.name = name;
	this.api = new Node(host, port);
	this.account = {};
}


const transferRaw = async (node, from, to, value) => {
  //	console.group('Locally Signed Transfer');
  console.log("TransferRAW");
/*
    const tx = new Transaction({from: from.address, to: to, value: value, gas: defaultGas, gasPrice: defaultGasPrice});
    tx.beforeSubmission();
console.log(tx);
   const receipt =  await  node.api.sendTx(tx, from);
  */  
    const receipt = await node.api.transfer(
      from,
      to,
      value,
      defaultGas,
      defaultGasPrice
    );
   
    console.log('Receipt: ', receipt);


  };



// Faucet parse...


const processJSON = async (input) => {

console.log("\n Loading "+input+"\n");
var content = fs.readFileSync(input);
var node;

try {
    const data = JSON.parse(content)
 //   console.log(data);

    const arraylength = data.Transactions.length;
    for (var i = 0; i < arraylength; i++ ) {
        arrSplit = data.Transactions[i].Node.split(":");

        if (i==0) {
          InitNode(data.Transactions[i].NodeName, arrSplit[0], arrSplit[1]);
          node = NodeCollection[data.Transactions[i].NodeName];
        }  

        await InitAccount(data.Transactions[i].FromName, password, datadir, data.Transactions[i].NodeName, data.Transactions[i].From);

        console.log(AccountCollection[data.Transactions[i].FromName].to + " "+ data.Transactions[i].Amount);

         let payer = data.Transactions[i].From;
         let payee = data.Transactions[i].To;
         let value = data.Transactions[i].Amount;

 //        console.log(AccountCollection[data.Transactions[i].FromName]);

 //         console.log(NodeCollection[data.Transactions[i].NodeName]);

         await transferRaw(node, 
              AccountCollection[data.Transactions[i].FromName], payee, value);
    }
 
  } catch(err) {
    console.error(err);
  }



  if (  total ) {

  try{
      const data = JSON.parse(content)
      var totals = [];

    //   console.log(data);
      const arraylength = data.Transactions.length;
      for (var i = 0; i < arraylength; i++ ) {
        baseAccount = await node.api.getAccount(data.Transactions[i].To);
       // console.log(baseAccount.address + " "+ baseAccount.balance.toNumber().toString());
        totals.push({"address": baseAccount.address, "balance": baseAccount.balance.toNumber()});
      }

      fs.writeFile (total, JSON.stringify(totals), function(err) {
        if (err) throw err;
        console.log('complete');
        }
    );

    } catch(err) {
      console.error(err);
    }
  }

};



processJSON(faucetfile)
.then(console.log)
.catch(console.log);

