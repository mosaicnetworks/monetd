const path = require("path");
const BigNumber = require('bignumber.js');

const argv = require('minimist')(process.argv.slice(2));
const JSONbig = require('json-bigint');

// import required objects
const { default:Node, Contract, Account, Transaction } = require("evm-lite-core");
const { default:Keystore } = require("evm-lite-keystore");
const { default:DataDirectory } = require("evm-lite-datadir");

const acct = argv.account;
const nodename = argv.nodename;
const nodehost = argv.nodehost;
const nodeport = argv.nodeport;
const transfile = argv.transfile;
const outfile = argv.outfile;
const networkpath = argv.configdir


const defaultGas = 10000000;
const defaultGasPrice = 0;

const keystorepath = path.join(networkpath, "keystore");
const password = "test";  //TODO live read
const keystore = new Keystore(keystorepath);
const datadir = new DataDirectory(networkpath, "monetcli", keystore);
 

var fs = require("fs");



const getAccount = async (address, password, datadir) => {
  console.log("Decrypting "+address+", "+password)
  const keyfile = await datadir.keystore.get(address);
  return Keystore.decrypt(keyfile, password);

};

function NewNode(name, host, port) {
  console.log("Creating Node "+name);
	this.name = name;
	this.api = new Node(host, port);
	this.account = {};
}


const transferRaw = async (node, from, to, value) => {
    const receipt = await node.api.transfer(
      from,
      to,
      value,
      defaultGas,
      defaultGasPrice
    );
   
    console.log('Status: ', receipt.status);
  };



  const asynchTransferRaw = async (node, from, to, value, nonce) => {
    //	console.group('Locally Signed Transfer');
    console.log("AsynchTransferRAW");
  
      const tx = new Transaction({from: from.address, to: to, value: value, gas: defaultGas, gasPrice: defaultGasPrice});
      tx.beforeSubmission();

      tx.nonce = nonce;

 //     receipt = await node.api.sendTx(tx, from);
      tx.signed = await from.signTx(tx);

      console.log(tx.signed);
     
      return rawURI = "http://"+nodehost+":"+nodeport+"/rawtx "+tx.signed.rawTransaction
    };
  
  



const processAccount = async (accountName) => {

  console.log("\n Loading "+transfile+" searching for "+accountName+"\n");
  var content = fs.readFileSync(transfile);

  var arrTrans = {};
  var node;
  var thisAccount;
  var runSynchronously = false;


  
  try {
      const data = JSONbig.parse(content)
 //     console.log(data);
  
  
      const arraylength = data.length;
      for (var i = 0; i < arraylength; i++ ) {
          if (data[i].Moniker == accountName) {
             console.log("Found Account "+accountName+" "+ i)
             if (i==0) { runSynchronously = true;}
             node = new NewNode(nodename, nodehost, nodeport);
             thisAccount = await getAccount(data[i].Moniker, password, datadir);
             

             var baseAccount = await node.api.getAccount(thisAccount.address);
             console.log("Account has  : "+baseAccount.balance.toFixed()+ "aŦ");
             console.log("Account needs: "+data[i].Debits.toFixed()+ "aŦ");
             
             if (baseAccount.balance<data[i].Debits ) {
                console.log("Account has insufficient credit to guarantee successful execution")
//TODO                process.exit(3);
                
             }

             arrTrans = "";

             if (data[i].Transactions)
             {
              for (var j = 0; j<data[i].Transactions.length; j++) {

                  if (data[i].Transactions[j].Amount == 0) {
                      continue;
                  }

                  if (runSynchronously) 
                  {
                      await transferRaw(node, thisAccount, data[i].Transactions[j].To, data[i].Transactions[j].Amount);
                  } else {
                      console.log("Asynchronous execution selected")
                      uri = await asynchTransferRaw(node, thisAccount, data[i].Transactions[j].To, data[i].Transactions[j].Amount, baseAccount.nonce+data[i].Transactions[j].Nonce-1);
                      arrTrans=arrTrans+"\n"+uri;
                  }
              }
             }  


             if (arrTrans!="") {
              console.log("Writing "+outfile); 
              fs.writeFile(outfile, arrTrans, function(err) {
                 if (err) {console.log("error writing file : "+err)}
              });
             }
             return; 
          }
        }    

        console.log("Account "+accountName+" is not in transaction data. Aborting")
        process.exit(2);

        return ;
   
    } catch(err) {
      console.error(err);
    }
  
  } 



  processAccount(acct)
  .catch(err => console.log(err));
