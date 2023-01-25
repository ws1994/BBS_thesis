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
// const client = require('fabric-common')
// const client = require('fabric-client')
const object = require('object-assign')
var log4js = require('log4js')
const cli = require('fabric-common')
// const channel = require('fabric-channel')

// Main program function
async function main() {

    // A wallet stores a collection of identities for use
    const wallet = await Wallets.newFileSystemWallet('../identity/user/isabella/wallet');

    // A gateway defines the peers used to access Fabric networks
    const gateway = new Gateway();

    try {

        // Specify userName for network access
        // const userName = 'isabella.issuer@magnetocorp.com';
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
        // var endorsingPeer0Org3 = channel.getEndorser('peer0.org3.example.com:11051');

        // Get addressability to commercial paper contract
        console.log('Use bigFileTransfer2 smart contract.');
        const contract = await network.getContract('bigFileTransfer2');

        // issue commercial paper
        console.log('Start initializing...');


        console.log('\n### Init config ...')
        var issueResponse = await contract.createTransaction('InitConfig')
            .setEndorsingPeers([endorsingPeer0Org1,endorsingPeer0Org2])
            .submit();
        console.log('Init Config transaction response: '+issueResponse);


        console.log('\n### Set IP of ip.Org1MSP ...')
        var issueResponse = await contract.createTransaction('AddIPInfo')
            .setEndorsingPeers([endorsingPeer0Org1,endorsingPeer0Org2])
            .submit("{\"name\":\"ip.Org1MSP\",\"IpAdr\":\"172.27.0.11\",\"IpPort\":22}");
        console.log('\n ### Set IP of ip.Org1MSP transaction response: '+issueResponse);

        console.log('\n### Get ip for ip.Org1MSP ...')
        var getResponse = await contract.createTransaction('GetIpInfo')
            .setEndorsingPeers([endorsingPeer0Org1])
            .evaluate('ip.Org1MSP');
        console.log('\n### Get ip for ip.Org1MSP query response: '+getResponse);

        console.log('\n### Get key policy of ip.Org1MSP ...')
        var getResponse = await contract.createTransaction('GetKeyPolicy')
            .setEndorsingPeers([endorsingPeer0Org1])
            .evaluate('ip.Org1MSP');
        console.log('\n### Get key policy of ip.Org1MSP query response: '+getResponse);

        console.log('\n### Set key policy for ip.Org1MSP ...')
        var issueResponse = await contract.createTransaction('SetKeyPolicy')
            .setEndorsingPeers([endorsingPeer0Org1,endorsingPeer0Org2])
            .submit('ip.Org1MSP','Org2MSP');
        console.log('\n### Set key policy for ip.Org1MSP transaction response: '+issueResponse);

        console.log('\n### Get key policy of ip.Org1MSP ...')
        var getResponse = await contract.createTransaction('GetKeyPolicy')
            .setEndorsingPeers([endorsingPeer0Org1])
            .evaluate('ip.Org1MSP');
        console.log('\n### Get key policy of ip.Org1MSP query response: '+getResponse);


        console.log('\n### Set IP of ip.Org2MSP ...')
        var issueResponse = await contract.createTransaction('AddIPInfo')
            .setEndorsingPeers([endorsingPeer0Org1,endorsingPeer0Org2])
            .submit("{\"name\":\"ip.Org2MSP\",\"IpAdr\":\"172.27.0.10\",\"IpPort\":22}");
        console.log('\n### Set IP of ip.Org2MSP transaction response: '+issueResponse);

        console.log('\n### Set key policy for ip.Org2MSP ...')
        issueResponse = await contract.createTransaction('SetKeyPolicy')
            .setEndorsingPeers([endorsingPeer0Org1,endorsingPeer0Org2])
            .submit("ip.Org2MSP","Org1MSP");
        console.log('\n### Set key policy for ip.Org2MSP transaction response: '+issueResponse);

        console.log('\n### Get key policy of ip.Org2MSP ...')
        var getResponse = await contract.createTransaction('GetKeyPolicy')
            .setEndorsingPeers([endorsingPeer0Org1])
            .evaluate('ip.Org2MSP');
        console.log('\n### Get key policy of ip.Org2MSP query response: '+getResponse);


        console.log('\n### Get ip for ip.Org2MSP ...')
        getResponse = await contract.createTransaction('GetIpInfo')
            .setEndorsingPeers([endorsingPeer0Org2])
            .evaluate('ip.Org2MSP');
        console.log('\n### GetIpInfo of ip.Org2MSP query response: '+getResponse);


        console.log('\n### Init org1 rsa key ...')
        issueResponse = await contract.createTransaction('Set')
            .setEndorsingPeers([endorsingPeer0Org1,endorsingPeer0Org2])
            .submit("Org1MSPPublicKey","initRSA");
        console.log('\n### Init org1 rsa key transaction response: '+issueResponse);

        console.log('\n### Set Org1MSPPublicKey key policy ...')
        issueResponse = await contract.createTransaction('SetKeyPolicy')
            .setEndorsingPeers([endorsingPeer0Org1,endorsingPeer0Org2])
            .submit("Org1MSPPublicKey","Org1MSP");
        console.log('\n### Set Org1MSPPublicKey key policy transaction response: '+issueResponse);

        console.log('\n### Set Org1MSPPublicKey org1 RSA key ...')
        issueResponse = await contract.createTransaction('InitRSAKeyForCC')
            .setEndorsingPeers([endorsingPeer0Org1])
            .submit();
        console.log('\n### Set Org1MSPPublicKey org1 RSA key transaction response: '+issueResponse);

        console.log('\n### Get Org1MSPPublicKey ...')
        var getResponse = await contract.createTransaction('Get')
            .setEndorsingPeers([endorsingPeer0Org1])
            .evaluate('Org1MSPPublicKey');
        console.log('\n### Get Org1MSPPublicKey query response: '+getResponse);


        console.log('\n### Init org2 rsa key ...')
        issueResponse = await contract.createTransaction('Set')
            .setEndorsingPeers([endorsingPeer0Org1,endorsingPeer0Org2])
            .submit("Org2MSPPublicKey","initRSA");
        console.log('\n### Init org2 rsa key transaction response: '+issueResponse);

        console.log('\n### Set Org2MSPPublicKey key policy ...')
        issueResponse = await contract.createTransaction('SetKeyPolicy')
            .setEndorsingPeers([endorsingPeer0Org1,endorsingPeer0Org2])
            .submit("Org2MSPPublicKey","Org2MSP");
        console.log('\n### Set Org2MSPPublicKey key policy transaction response: '+issueResponse);

        console.log('\n### Set Org2MSPPublicKey org2 RSA key ...')
        issueResponse = await contract.createTransaction('InitRSAKeyForCC')
            .setEndorsingPeers([endorsingPeer0Org2])
            .submit();
        console.log('\n### Set Org2MSPPublicKey org1 RSA key transaction response: '+issueResponse);

        console.log('\n### Get Org2MSPPublicKey ...')
        var getResponse = await contract.createTransaction('Get')
            .setEndorsingPeers([endorsingPeer0Org2])
            .evaluate('Org2MSPPublicKey');
        console.log('\n### Get Org2MSPPublicKey query response: '+getResponse);



        console.log('\n#############################################################');
        console.log('End initializing. Please install the SSH in chaincode dockers');
        console.log('#############################################################');

        console.log('\n#############################################################');
        console.log('Please map the host port to containers ip:22 port');
        console.log('#############################################################');

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
main().then(() => {

    console.log('Issue program complete.');

}).catch((e) => {

    console.log('Issue program exception.');
    console.log(e);
    console.log(e.stack);
    process.exit(-1);

});