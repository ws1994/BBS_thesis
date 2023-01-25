// class SearchFile{

async function searchAllFile() {
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
    const wallet = await Wallets.newFileSystemWallet('../../magnetocorp/identity/user/isabella/wallet');

    // A gateway defines the peers used to access Fabric networks
    const gateway = new Gateway();

    // var version = '1'
    // version = process.argv.splice(2)

    try {

        // Specify userName for network access
        // const userName = 'isabella.issuer@magnetocorp.com';
        const userName = 'isabella';

        // Load connection profile; will be used to locate a gateway
        let connectionProfile = yaml.safeLoad(fs.readFileSync('../../magnetocorp/gateway/connection-org2.yaml', 'utf8'));

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

        // Get addressability to commercial paper contract
        console.log('Use bigFileTransfer smart contract.');
        const contract = await network.getContract('bigFileTransfer2');
     

        console.log('\n### List All Files ...')
        var getResponse = await contract.createTransaction('ListAllFiles')
            .setEndorsingPeers([endorsingPeer0Org1])
            .evaluate();
        console.log('List all files query response: ');
        // var result = getResponse.toString();
        // result = result.replace(/(^\[)|(\]$)/g,'');
        // var resultArr = result.split("},{");
        // for (let i = 0; i < resultArr.length; ++i) {
        //     resultArr[i] = resultArr[i].replace(/(^{)|(}$)/g,'')
        //     console.log(resultArr[i]);
        //     console.log();
        // }

        var resJson = JSON.parse(getResponse);
        console.log(resJson);

        console.log('\n##############');
        console.log('End search file.')
        console.log('################');

        return resJson;

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
// }
// module.exports = SearchFile;
module.exports = { searchAllFile };