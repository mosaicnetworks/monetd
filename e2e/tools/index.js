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
const networkpath = argv.configdir;
const total = argv.total;
const pre = argv.pre;

const defaultGas = 10000000;
const defaultGasPrice = 0;

const keystorepath = path.join(networkpath, "keystore");
const password = "test";  //TODO live read
const keystore = new Keystore(keystorepath);
const datadir = new DataDirectory(networkpath, "monetcli", keystore);
 

var fs = require("fs");



const getAccount = async (address, password, datadir) => {
//  console.log("Decrypting "+address+", "+password)
  const keyfile = await datadir.keystore.get(address);
  return Keystore.decrypt(keyfile, password);

};

function NewNode(name, host, port) {
//  console.log("Creating Node "+name);
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
   
    process.stdout.write(receipt.status);
  };



  const asynchTransferRaw = async (node, from, to, value, nonce) => {
    //	console.group('Locally Signed Transfer');
  //  console.log("AsynchTransferRAW");  
      const tx = new Transaction({from: from.address, to: to, value: value, gas: defaultGas, gasPrice: defaultGasPrice});
      tx.beforeSubmission();
      tx.nonce = nonce;
 //     receipt = await node.api.sendTx(tx, from);
      tx.signed = await from.signTx(tx);
 //     console.log(tx.signed);
      return rawURI = "http://"+nodehost+":"+nodeport+"/rawtx "+tx.signed.rawTransaction
    };
  
  



const processAccount = async (accountName) => {

 // console.log("\n Loading "+transfile+" searching for "+accountName+"\n");
  var content = fs.readFileSync(transfile);

  var arrTrans = {};
  var node;
  var thisAccount;
  var runSynchronously = false;


  
  try {

      var found = false;
      const data = JSONbig.parse(content)
 //     console.log(data);
  
  
      const arraylength = data.length;
      for (var i = 0; i < arraylength; i++ ) {
          if (data[i].Moniker == accountName) {
          //   console.log("Found Account "+accountName+" "+ i)
             if (i==0) { runSynchronously = true;}
             node = new NewNode(nodename, nodehost, nodeport);
             thisAccount = await getAccount(data[i].Moniker, password, datadir);
             

             var baseAccount = await node.api.getAccount(thisAccount.address);
             console.log("Account "+accountName+" has  : "+baseAccount.balance.toFixed()+ "aŦ needs: "+data[i].Debits.toFixed()+ "aŦ");
  
// The types of these variables is variable - so we force them into a BigNumber.
             var has = new BigNumber(baseAccount.balance);
             var needs = new BigNumber(data[i].Debits);
             
             
             if (needs.isGreaterThan(has)) {
                console.log("Account has insufficient credit to guarantee successful execution")
                process.exit(3);
                
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
               //       console.log("Asynchronous execution selected")
                      uri = await asynchTransferRaw(node, thisAccount, data[i].Transactions[j].To, data[i].Transactions[j].Amount, baseAccount.nonce+data[i].Transactions[j].Nonce-1);
                      arrTrans=arrTrans+"\n"+uri;
                  }
              }
             }  


             if (arrTrans!="") {
 //             console.log("Writing "+outfile); 
              fs.writeFile(outfile, arrTrans, function(err) {
                 if (err) {console.log("error writing file : "+err)}
              });
             }
             found=true;
             break; 
          }
        }    

        if (! found) {
        console.log("Account "+accountName+" is not in transaction data. Aborting")
        process.exit(2);
        }
   
    } catch(err) {
      console.error(err);
    }


  //  console.log("Total check")
    if ( total ) {
      try{
          const data = JSONbig.parse(content)
          var totals = [];
    
        //   console.log(data);
          const arraylength = data.length;
          for (var i = 0; i < arraylength; i++ ) {
            baseAccount = await node.api.getAccount(data[i].Address);
    
           // console.log(baseAccount.address + " "+ baseAccount.balance.toNumber().toString());
            totals.push({"address": baseAccount.address, 
               "balance": baseAccount.balance.toFixed()
              });
          }
    
          fs.writeFile (total, JSONbig.stringify(totals), function(err) {
            if (err) throw err;
           // console.log('complete');
            }
        );
    
        } catch(err) {
          console.error(err);
        }
      }
  } 


// Verify that totals match
const checkTotals = async () => {
//  console.log("\n Loading "+transfile+"\n");
  var transcontent = fs.readFileSync(transfile);

//  console.log("\n Loading "+pre+"\n");
  var precontent = fs.readFileSync(pre);

  node = new NewNode(nodename, nodehost, nodeport);
 //  thisAccount = await getAccount(data[i].Moniker, password, datadir);
 //  var baseAccount = await node.api.getAccount(thisAccount.address);



  
  try {
      const predata = JSONbig.parse(precontent)
  //    console.log(predata);
      const transdata = JSONbig.parse(transcontent)
  //    console.log(deltadata);

      const prearraylength = predata.length;
      const transarraylength = transdata.length;



      console.log("Calculating Totals ("+prearraylength+", "+transarraylength+"):");

      var failures = 0; 

      for (var i = 0; i < prearraylength; i++ ) {        
        var addr = predata[i].address.replace(/^(0[xX])/,"").toUpperCase();
//        console.log("Pre "+ i+ " " + addr );
        for (var j = 1; j < transarraylength; j++ ) {  
          var addr2 = transdata[j].Address.replace(/^(0[xX])/,"").toUpperCase();
//          console.log("trans "+ j + " " + addr2) ;       
          if (addr == addr2){
            predata[i]["PreBalance"] = new BigNumber(predata[i].balance);
            var baseAccount = await node.api.getAccount(predata[i].address);
            predata[i]["PostBalance"] = new BigNumber(baseAccount.balance);
            predata[i]["PostNet"] = predata[j]["PostBalance"].minus(predata[i]["PreBalance"]);
            predata[i]["TransNet"] = transdata[j]["Delta"];         
            predata[i]["Diff"] = predata[i]["PostNet"].minus(transdata[j]["Delta"]);
            console.log(transdata[j].Moniker + ":");
            console.log("  Post   : "+baseAccount.balance.toFixed());              
            console.log("  Pre    : "+predata[i]["PreBalance"].toFixed());               
            console.log("  Calc   : "+predata[i]["TransNet"].toFixed());                            
            console.log("  Actual : "+predata[i]["PostNet"].toFixed());                            
            console.log("  Diff   : "+predata[i]["Diff"].toFixed());  
            
            if ( ! predata[i]["Diff"].isZero() ){
              failures++;
            } 
            
            break;
          }
        }
      }
      console.log("Calculated Totals"); 
      if (failures > 0) {
          console.log("Balance Failures: "+failures)
             process.exit(1);
      }
      console.log("Balance Checks Passed.")
//      console.log(deltadata);
  } catch(err) {
    console.error(err);
    process.exit(2);
  }
}





if (pre){
  checkTotals()
  .catch(err => console.log(err));
}
else
{
  processAccount(acct)
  .catch(err => console.log(err));
}