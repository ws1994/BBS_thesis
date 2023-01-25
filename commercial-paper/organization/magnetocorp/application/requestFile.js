/*
SPDX-License-Identifier: Apache-2.0
*/

/*
 * This application has 6 basic steps:
 * 1. Select an identity from a wallet
 * 2. Connect to network gateway
 * 3. Access PaperNet network
 * 4. Construct request to issue commercial paper
 * 5. Submit transaction
 * 6. Process response
 */

'use strict';

// Bring key classes into scope, most importantly Fabric SDK network class
const fs = require('fs');
const yaml = require('js-yaml');
const { Wallets, Gateway } = require('fabric-network');
const CommercialPaper = require('../contract/lib/paper.js');
// const Client = require('fabric-common')
const object = require('object-assign')
var log4js = require('log4js')
const cli = require('fabric-common')
// const channel = require('fabric-channel')
// This is a package implemneted by shan
let Ut = require("./common");
var NodeRSA = require('node-rsa')
const rsa = require('trsa')

var version = '1'

// Main program function
async function main() {


    // A wallet stores a collection of identities for use
    const wallet = await Wallets.newFileSystemWallet('../identity/user/isabella/wallet');

    // A gateway defines the peers used to access Fabric networks
    const gateway = new Gateway();

    version = process.argv.splice(2)

    try {
        console.time('wholeProcess' + version)

        console.time('connect' + version)
        const userName = 'isabella';

        // Load connection profile; will be used to locate a gateway
        let connectionProfile = yaml.safeLoad(fs.readFileSync('../gateway/connection-org2.yaml', 'utf8'));

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


        var endorsingOrg1 = endorsingPeer0Org1
        var endorsingOrg2 = endorsingPeer0Org2
        var des = 'peer0.org2'



        // Get addressability to commercial paper contract
        console.log('Use bigFileTransfer smart contract.');
        const contract = await network.getContract('bigFileTransfer2');
        console.timeEnd('connect' + version)

        // issue commercial paper
        console.log('Start file transfer...');

        const fileName = version

        console.log('\n### Request permission ...' + version)
        console.time('RequestSftp_' + version)
        var issueResponse = await contract.createTransaction('RequestSftp')
            .setEndorsingPeers([endorsingOrg1,endorsingOrg2])
            .submit(fileName);
        console.timeEnd('RequestSftp_' + version)
        console.log(version + 'RequestSftp transaction response: '+issueResponse);


        console.time('sleep')
        await Ut.sleep(3000);
        console.timeEnd('sleep')

        console.log('\n### Transfer file by shell sftp ...' + version)
        console.time('SendFile_SSH_sftp_' + version)
        issueResponse = await contract.createTransaction('SendFile_SSH_sftp')
            .setEndorsingPeers([endorsingOrg1])
            .submit(fileName);
        console.timeEnd('SendFile_SSH_sftp_' + version)
        console.log(version + 'Transfer file by shell sftp transaction response: '+issueResponse); 

        // sleep(3000)
        console.time('sleep' + version)
        await Ut.sleep(3000);
        console.timeEnd('sleep' + version)


        try{
            console.log('\n### Request key_z from peer0.org1 ...' + version)
            console.time('RequestKeyZ_' + version)
            issueResponse = await contract.createTransaction('RequestKey')
                .setEndorsingPeers([endorsingOrg1])//endorsingOrg1
                .submit(fileName);
            console.timeEnd('RequestKeyZ_' + version)
            console.log(version + 'Request key_z from peer0.org3 transaction response: '+issueResponse);
        } catch (error) {
                console.log(`*** Caught the Request key error: \n    ${error}`);
        }

         // sleep(3000)
        console.time('sleep')
        await Ut.sleep(3000);
        console.timeEnd('sleep')


        try{
            console.log('\n### Request key_r from peer0.org1 ...' + version)
            console.time('RequestKeyR_' + version)
            issueResponse = await contract.createTransaction('RequestKeyReceiver')
                .setEndorsingPeers([endorsingOrg2])//endorsingOrg1
                .submit(fileName);
            console.timeEnd('RequestKeyR_' + version)
            console.log(version + 'Request key_r from peer0.org3 transaction response: '+issueResponse);
        } catch (error) {
                console.log(`*** Caught the Request key error: \n    ${error}`);
        }
        

        // sleep(3000)
        console.time('sleep')
        await Ut.sleep(3000);
        console.timeEnd('sleep')


        console.log('\n### Get k_r1 ...')
        var kr = await contract.createTransaction('QueryKeyReceiver')
            .setEndorsingPeers([endorsingPeer0Org2])
            .evaluate(fileName);
        console.log('Get k_r query response: '+ kr);

        console.log('Start decrypt k_r');
        var data = fs.readFileSync('./clientRSAKeys/ccPrivate.pem');
        console.log(data.toString());
        var key = new NodeRSA(data);
        let plaintext = key.decrypt(kr);
        console.log(plaintext);
        console.log('End decrypt');

        console.time('DecryptFile_' + version)
        try{
            // var k_r1 := string(getResponse)
            var key_buff = plaintext
            console.log('\n### start decrypt ...')
            var getResponse = await contract.createTransaction('DecryptFile')
                .setEndorsingPeers([endorsingPeer0Org2])
                .setTransient({"key":key_buff})
                .submit(fileName);
            console.log('Get k_r query response: '+getResponse);
        }catch (error) {
            console.log(version + 'decrypt file transaction response: '+issueResponse);
        }
        console.timeEnd('DecryptFile_' + version)



        console.log('\n### List folder in User1@org2.example.com ...' + version)
        console.time('ListAndDelete' + version)
        getResponse = await contract.createTransaction('ListFolder')
            .setEndorsingPeers([endorsingOrg1])
            .evaluate('/home/chaincode/User1@org2.example.com');
        console.log(version + 'List folder in User1@org2.example.com: '+getResponse);

        // try{
        //     console.log('\n### Remove file ' + fileName + '_ciphertext.bin in peer0.org3...' + version)
        //     getResponse = await contract.createTransaction('DeleteFileInDocker')
        //         .setEndorsingPeers([endorsingOrg3])
        //         .submit('User1@org3.example.com/' + fileName + '_ciphertext.bin');
        //     console.log(version + 'Remove file ' + fileName + '_ciphertext.bin query response: '+getResponse);
        //     console.timeEnd('ListAndDelete' + version)
        // }catch (error) {

        // }


        // try{
        //     console.log('\n### Remove file ' + fileName + '_ciphertext.bin in peer0.org1...' + version)
        //     getResponse = await contract.createTransaction('DeleteFileInDocker')
        //         .setEndorsingPeers([endorsingOrg1])
        //         .submit('Admin@org1.example.com/' + fileName + '_ciphertext.bin');
        //     console.log(version + 'Remove file ' + fileName + '_ciphertext.bin query response: '+getResponse);
        //     console.timeEnd('ListAndDelete' + version)
        // }catch (error) {

        // }


        // try{
        //     console.log('\n### Remove file ' + fileName + 'enc_ in peer0.org1...' + version)
        //     getResponse = await contract.createTransaction('DeleteFileInDocker')
        //         .setEndorsingPeers([endorsingOrg1])
        //         .submit('Admin@org1.example.com/' + 'ecn_' + fileName);
        //     console.log(version + 'Remove file enc_' + fileName + ' query response: '+getResponse);
        //     console.timeEnd('ListAndDelete' + version)
        // }catch (error) {

        // }


        console.log('\n#################');
        console.log('End transfer file.' + version);
        console.log('###################');

   
        // let paper = CommercialPaper.fromBuffer(issueResponse);

        // console.log(`${paper.issuer} commercial paper : ${paper.paperNumber} successfully issued for value ${paper.faceValue}`);
        // console.log('Transaction complete.');
        console.timeEnd('wholeProcess' + version)

    } catch (error) {

        console.log(`Error processing transaction. ${error}`);
        console.log(error.stack);

    } finally {

        // Disconnect from the gateway
        console.log('\nDisconnect from Fabric gateway.' + version);
        // console.log('Process issue transaction response.'+issueResponse);

        gateway.disconnect();

    }
}
main().then(() => {

    console.log('Issue program complete.' + version);

}).catch((e) => {

    console.log('Issue program exception.' + version);
    console.log(e);
    console.log(e.stack);
    process.exit(-1);

});