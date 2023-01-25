
async function addFile(fileName,fileDes, fileRule) {
'use strict';

// Bring key classes into scope, most importantly Fabric SDK network class
const fs = require('fs');
const yaml = require('js-yaml');
const { Wallets, Gateway } = require('fabric-network');
const CommercialPaper = require('../contract/lib/paper.js');
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

        // const fileName = '67v'+ version + '.pdf'
        // const fileName = '218v'+ version + '.mp4'
        // const fileName = '576Mv'+ version + '.tif'
        // const fileName = '1_2gv'+ version + '.zip'
        // const fileName = '2_6gv'+ version + '.rar'
        // const fileName = '5_3gv'+ version + '.zip'
        // const fileName = '10_2gv'+ version + '.zip'

        console.log('\n### Set key for file ' + fileName + ' ...')
        var issueResponse = await contract.createTransaction('Set')
            .setEndorsingPeers([endorsingPeer0Org1,endorsingPeer0Org2])
            .submit(fileName,"initFile");
        console.log('Set key for file transaction response: '+issueResponse);


        console.log('\n### Get info of ' + fileName + ' ...')
        issueResponse = await contract.createTransaction('Get')
            .setEndorsingPeers([endorsingPeer0Org1])
            .evaluate(fileName);
        console.log('Get key for file transaction response: '+issueResponse);


        console.log('\n### Set key policy for ' + fileName + ' ...')
        var issueResponse = await contract.createTransaction('SetKeyPolicy')
            .setEndorsingPeers([endorsingPeer0Org1,endorsingPeer0Org2])
            .submit(fileName,'Org1MSP');
        console.log('Set key policy for '+ fileName + ' transaction response: '+issueResponse);

        console.log('\n### Get key policy for ' + fileName + ' ...')
        var getResponse = await contract.createTransaction('GetKeyPolicy')
            .setEndorsingPeers([endorsingPeer0Org1])
            .evaluate(fileName);
        console.log('Get key policy for '+ fileName + ' query response: '+getResponse);

        // console.log('\n### Get ip for ip.peer0.org1 ...')
        //     getResponse = await contract.createTransaction('GetIpInfo')
        //         .setEndorsingPeers([endorsingPeer0Org1])
        //         .evaluate('ip.peer0.org1');
        //     console.log('GetIpInfo of ip.peer0.org1 query response: '+getResponse);
            // sleep 3s

        console.log('\n### List folder in User1@org1.example.com...')
        var getResponse = await contract.createTransaction('ListFolder')
            .setEndorsingPeers([endorsingPeer0Org1])
            .evaluate('/home/chaincode/User1@org1.example.com');
        console.log('List folder in User1@org1.example.com query response: '+getResponse);


        // var endorsingOrg1 = endorsingPeer0Org1

        // result = result.replace(/(^\[)|(\]$)/g,'');
        var ruleArr = fileRule.split(";");
        var rule1 = ""
        var rule2 = ""
        if(ruleArr.length == 2){
            rule1 = ruleArr[0];
            rule2 = ruleArr[1];
            console.log('\n### Add file info for ' + fileName + ' ...')
            var issueResponse = await contract.createTransaction('AddFileInfo')
            .setEndorsingPeers([endorsingPeer0Org1])
            .submit("{\"name\":\"" + fileName + "\",\"filehash\":\"\",\"objecttype\":\"file\",\"description\":\" " + fileDes + " \",\"rule\":[\"" +rule1+ "\",\"" +rule2+ "\"]}");
            console.log('Add file info for ' + fileName + ' transaction response: '+issueResponse);
        }else{
            rule1 = fileRule;
            console.log('\n### Add file info for ' + fileName + ' ...')
            var issueResponse = await contract.createTransaction('AddFileInfo')
            .setEndorsingPeers([endorsingPeer0Org1])
            .submit("{\"name\":\"" + fileName + "\",\"filehash\":\"\",\"objecttype\":\"file\",\"description\":\" " + fileDes + " \",\"rule\":[\""+rule1+"\"]}");
            console.log('Add file info for ' + fileName + ' transaction response: '+issueResponse);
        }

        // console.log('\n### Add file info for ' + fileName + ' ...')
        // var issueResponse = await contract.createTransaction('AddFileInfo')
        //     .setEndorsingPeers([endorsingPeer0Org1])
        //     .submit("{\"name\":\"" + fileName + "\",\"filehash\":\"\",\"objecttype\":\"file\",\"description\":\" " + fileDes + " \",\"rule\":[\"Org2MSP\",\"Org3MSP\"]}");
        // console.log('Add file info for ' + fileName + ' transaction response: '+issueResponse);

        

        // console.log('\n### List All Files ...')
        // var getResponse = await contract.createTransaction('ListAllFiles')
        //     .setEndorsingPeers([endorsingPeer0Org1])
        //     .evaluate(fileName);
        // console.log('List all files query response: '+getResponse);


        console.log('\n##############');
        console.log('End add file.')
        console.log('################');

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
module.exports = { addFile };