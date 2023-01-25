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
var NodeRSA = require('node-rsa')
const rsa = require('trsa')
// const channel = require('fabric-channel')

// Main program function
async function main() {

    // var logger = log4js.getLogger()//'Helper'
    // logger.level = 'DEBUG'
    // client.setLogger(logger)

    // A wallet stores a collection of identities for use
    const wallet = await Wallets.newFileSystemWallet('../identity/user/isabella/wallet');

    // A gateway defines the peers used to access Fabric networks
    const gateway = new Gateway();

    // const client = Client.require('fabric-common')
    // const default_option = client.getConfigSetting('connection-options');
    // console.log(default_option)
    // const new_option = {
    //     'grpc.keepalive_timeout_ms': 400000,
    //     'request-timeout': 400000
    // };
    // const new_defaults = Object.assign(default_option,new_option);
    // client.setConfigSetting('connection-options',new_defaults);
    // // client.setConfigSetting('request-timeout',60000);
    // console.log(default_option)

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
        console.log('Use bigFileTransfer smart contract.');
        const contract = await network.getContract('bigFileTransfer2');

        var pk = '';
        var data = fs.readFileSync('./clientRSAKeys/ccPublic.pem','utf8');
        pk = data;
        console.log(pk);

        console.log('\n### SetRsaPKForClient ...');
        var getResponse = await contract.createTransaction('SetRsaPKForClient')
            .setEndorsingPeers([endorsingPeer0Org1,endorsingPeer0Org2])
            .submit(pk);
        console.log('SetRsaPKForClient response: '+ getResponse);

        console.log('\n### get Rsa PK For Client ...');
        var getResponse = await contract.createTransaction('Get')
            .setEndorsingPeers([endorsingPeer0Org1])
            .evaluate("User1@org2.example.comPublicKey");
        console.log('get Rsa PK For Client response: '+ getResponse);

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