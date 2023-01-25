async function buildChainCustody(keyName) {

'use strict';

// Bring key classes into scope, most importantly Fabric SDK network class
const fs = require('fs');
const yaml = require('js-yaml');
const { Wallets, Gateway } = require('fabric-network');
const CommercialPaper = require('../contract/lib/paper.js');
const BlockDecoder = require('./node_modules/fabric-common/lib/BlockDecoder.js');
// const client = require('fabric-common')
// const client = require('fabric-client')
const object = require('object-assign')
var log4js = require('log4js')
const cli = require('fabric-common')
// const channel = require('fabric-channel')


// Main program function

    // A wallet stores a collection of identities for use
    const wallet = await Wallets.newFileSystemWallet('../../digibank/identity/user/balaji/wallet');

    // A gateway defines the peers used to access Fabric networks
    const gateway = new Gateway();

    // var fileName = '1'
    // fileName = process.argv.splice(2)

    // Main try/catch block
    // var issueResponse
    try {

        // Specify userName for network access
        // const userName = 'isabella.issuer@magnetocorp.com';
        const userName = 'balaji';

        // Load connection profile; will be used to locate a gateway
        let connectionProfile = yaml.safeLoad(fs.readFileSync('../../digibank/gateway/connection-org1.yaml', 'utf8'));

        // Set connection options; identity and wallet
        let connectionOptions = {
            identity: userName,
            wallet: wallet,
            discovery: { enabled:true, asLocalhost: false }
        };

        // Connect to gateway using application specified parameters
        console.log('Connect to Fabric gateway.');
        await gateway.connect(connectionProfile, connectionOptions);

        // Access PaperNet network
        console.log('Use network channel: mychannel.');
        const network = await gateway.getNetwork('mychannel');
        const channel = network.getChannel();

        var endorsingPeer0Org1 = channel.getEndorser('peer0.org1.example.com:7051');
        var endorsingPeer0Org2 = channel.getEndorser('peer0.org2.example.com:9051');
        // var endorsingPeer0Org3 = channel.getEndorser('peer0.org3.example.com:11051');

        // Get addressability to commercial paper contract
        console.log('Use bigFileTransfer smart contract.');
        const contract = await network.getContract('bigFileTransfer2');

        // var keyName = "218v1.mp4";//ip.Org1MSP  218v1.mp4 67v1.pdf



        var resJsonArr = [];

        console.log('\n### Get info of Tx history ...')
        var issueResponse = await contract.createTransaction('QueryTxHistoryForFile')
            .setEndorsingPeers([endorsingPeer0Org1])
            .evaluate(keyName);// 
        // console.log(issueResponse.toString());
        var resJson0 = JSON.parse(issueResponse);
        console.log(resJson0);
        resJsonArr.push(resJson0);

        console.log('\n### Get info of Tx RequestSftp ...')
        issueResponse = await contract.createTransaction('QueryTxHistoryForFile')
            .setEndorsingPeers([endorsingPeer0Org1])
            .evaluate(keyName+"RequestSftp");//
        // console.log(issueResponse.toString());
        var resJson1 = JSON.parse(issueResponse)
        console.log(resJson1);
        resJsonArr.push(resJson1);

        console.log('\n### Get info of Tx SendFile_SSH_sftp ...')
        issueResponse = await contract.createTransaction('QueryTxHistoryForFile')
            .setEndorsingPeers([endorsingPeer0Org1])
            .evaluate(keyName+"SendFile_SSH_sftp");//
        // console.log(issueResponse.toString());
        var resJson2 = JSON.parse(issueResponse)
        console.log(resJson2);
        resJsonArr.push(resJson2);

        console.log('\n### Get info of Tx RequestKey ...')
        issueResponse = await contract.createTransaction('QueryTxHistoryForFile')
            .setEndorsingPeers([endorsingPeer0Org1])
            .evaluate(keyName+"RequestKey");//
        // console.log(issueResponse.toString());
        var resJson3 = JSON.parse(issueResponse)
        console.log(resJson3);
        resJsonArr.push(resJson3);

        console.log('\n### Get info of Tx RequestKeyReceiver ...')
        issueResponse = await contract.createTransaction('QueryTxHistoryForFile')
            .setEndorsingPeers([endorsingPeer0Org1])
            .evaluate(keyName+"RequestKeyReceiver");//
        // console.log(issueResponse.toString());
        var resJson4 = JSON.parse(issueResponse)
        console.log(resJson4);
        resJsonArr.push(resJson4);

        console.log('\n### Get info of Tx decrypt ...')
        issueResponse = await contract.createTransaction('QueryTxHistoryForFile')
            .setEndorsingPeers([endorsingPeer0Org1])
            .evaluate(keyName+"decrypt");//
        // console.log(issueResponse.toString());
        var resJson5 = JSON.parse(issueResponse)
        console.log(resJson5);
        resJsonArr.push(resJson5);

        var txChainArr = [];

        for(let n in resJsonArr){
            var resJson = resJsonArr[n];
            for (let i in resJson){
                var TxChain = {
                    "txID": "",
                    "operation": "",
                    "user": "",
                    "time":"",
                    "blockHash":"",
                    "blockPreviousHash":""
                };

                console.log(i);
                // i++;
                console.log(resJson[i]);
                TxChain.txID = resJson[i].TxId;
                TxChain.time = resJson[i].Timestamp;

                const contractSys = await network.getContract('qscc');
                // console.log('\n### Get info of Tx ...')
                issueResponse = await contractSys.createTransaction('GetTransactionByID')
                    .setEndorsingPeers([endorsingPeer0Org1])
                    .evaluate("mychannel",resJson[i].TxId);

                var blockResponse = await contractSys.createTransaction('GetBlockByTxID')
                    .setEndorsingPeers([endorsingPeer0Org1])
                    .evaluate("mychannel",resJson[i].TxId);
                let blockJson = BlockDecoder.decode(blockResponse);
                var blockID = blockJson.header.data_hash.toString('base64');
                var blockPreID = blockJson.header.previous_hash.toString('base64');
                console.log("blockHash is  " + blockID);
                console.log("blockPreviousHash is  " + blockPreID);

                TxChain.blockHash = blockID;
                TxChain.blockPreviousHash = blockPreID;
                    
                let processedTransaction = BlockDecoder.decodeTransaction(issueResponse);

                // Iterate over actions
                const actions = processedTransaction.transactionEnvelope.payload.data.actions;
                for (const action of actions) {
                  TxChain.user = action.header.creator.mspid
                  var args = action.payload.chaincode_proposal_payload.input.chaincode_spec.input.args;
                  for (var j = 0;j<args.length;j++){
                    console.log(args[j].toString());
                  }
                  TxChain.operation = args[0].toString();
                  // console.log(args[0].toString);
                }

                if(TxChain.operation != "Set"){
                    txChainArr.push(TxChain);
                }  
            }
        }

        console.log(txChainArr);


        console.log('\n##############');
        console.log('End list off-state file.')
        console.log('################');

        return txChainArr;

    } catch (error) {

        console.log(`Error processing transaction. ${error}`);
        console.log(error.stack);

    } finally {

        // Disconnect from the gateway
        console.log('\nDisconnect from Fabric gateway.');
        // console.log('Process issue transaction response.'+issueResponse);

        gateway.disconnect();

    }
}

module.exports = { buildChainCustody };