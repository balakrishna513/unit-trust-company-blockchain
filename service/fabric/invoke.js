/*
 * SPDX-License-Identifier: Apache-2.0
 */

'use strict';

const { Gateway, Wallets } = require('fabric-network');
const fs = require('fs');
const path = require('path');
const commonUtil = require('./common');

exports.invoke = async function (requestObject, user) {
    try {
        //queryPeer, channelName,request.chaincodeName, request.args, request.fcn, userObj.userId, userObj.org
        // load the network configuration
        const ccpPath = path.resolve(__dirname, '..', '..', 'test-network', 'organizations', 'peerOrganizations', 'org1.example.com', 'connection-org1.json');
        let ccp = JSON.parse(fs.readFileSync(ccpPath, 'utf8'));

        // Create a new file system based wallet for managing identities.
        const walletPath = path.join(process.cwd(), 'wallet');
        const wallet = await Wallets.newFileSystemWallet(walletPath);
        console.log(`Wallet path: ${walletPath}`);

        // Check to see if we've already enrolled the user.
        let identity = await commonUtil.getIdentity(wallet, user);
        if (!identity) {
            console.log(`An identity for the user ${user} does not exist in the wallet`);
            throw new Error("IDENTITY_NOF_FOUND");
        }

        // Create a new gateway for connecting to our peer node.
        const gateway = new Gateway();
        await gateway.connect(ccp, { wallet, identity: user, discovery: { enabled: true, asLocalhost: true } });

        // Get the network (channel) our contract is deployed to.
        const network = await gateway.getNetwork('mychannel');

        // Get the contract from the network.
        const contract = network.getContract('mtct_cc');

        // Submit the specified transaction.
        // createCar transaction - requires 5 argument, ex: ('createCar', 'CAR12', 'Honda', 'Accord', 'Black', 'Tom')
        // changeCarOwner transaction - requires 2 args , ex: ('changeCarOwner', 'CAR12', 'Dave')
        //await contract.submitTransaction('createCar', 'CAR12', 'Honda', 'Accord', 'Black', 'Tom');
        console.log("requestObject.args::", requestObject.args)
        let response = await contract.submitTransaction(requestObject.funcName, requestObject.args);
        console.log('Transaction has been submitted');

        // Disconnect from the gateway.
        await gateway.disconnect();
        return response;
    } catch (error) {
        console.error(`\n\n Failed to submit transaction: ${error}`);
        if(error.responses && error.responses[0] && error.responses[0].response.message) {
            error = error.responses[0].response.message;
        }else{
            error = error.message;
        }
        //error = (error.responses[0] && error.responses[0].response.message) ? error.responses[0].response.message : error.message;
        throw error;
    }
}
