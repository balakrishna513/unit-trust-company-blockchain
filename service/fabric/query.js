/*
 * SPDX-License-Identifier: Apache-2.0
 */

'use strict';

const { Gateway, Wallets } = require('fabric-network');
const path = require('path');
const fs = require('fs');
const commonUtil = require('./common');

exports.query = async function(requestObject, user) {
    try {
        // load the network configuration
        const ccpPath = path.resolve(__dirname, '..', '..', 'test-network', 'organizations', 'peerOrganizations', 'org1.example.com', 'connection-org1.json');
        const ccp = JSON.parse(fs.readFileSync(ccpPath, 'utf8'));

        // Create a new file system based wallet for managing identities.
        const walletPath = path.join(process.cwd(), 'wallet');
        const wallet = await Wallets.newFileSystemWallet(walletPath);
        console.log(`Wallet path: ${walletPath}`);

        // Check to see if we've already enrolled the user.
        let identity = await commonUtil.getIdentity(wallet, user);
        if (!identity) {
            console.log(`An identity for the user ${user} does not exist in the wallet`);
            return;
        }

        // Create a new gateway for connecting to our peer node.
        const gateway = new Gateway();
        await gateway.connect(ccp, { wallet, identity: user, discovery: { enabled: true, asLocalhost: true } });

        // Get the network (channel) our contract is deployed to.
        const network = await gateway.getNetwork('mychannel');

        // Get the contract from the network.
        const contract = network.getContract('mtct_cc');

        // Evaluate the specified transaction.
        // queryCar transaction - requires 1 argument, ex: ('queryCar', 'CAR4')
        // queryAllCars transaction - requires no arguments, ex: ('queryAllCars')
        console.log("args::", Object.values(requestObject.args));
        
        const result = await contract.evaluateTransaction(requestObject.funcName, Object.values(requestObject.args));
        console.log(`Transaction has been evaluated, result is: ${result.toString()}`);
        return result ? JSON.parse(result): result;
    } catch (error) {
        console.error(`Failed to evaluate transaction: ${error}`);
        throw error;
    }
}
