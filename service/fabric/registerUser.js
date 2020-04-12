/*
 * SPDX-License-Identifier: Apache-2.0
 */

'use strict';

const { Wallets } = require('fabric-network');
const FabricCAServices = require('fabric-ca-client');
const fs = require('fs');
const path = require('path');

exports.registerUser = async function(username, org, metadata) {
    console.log("registering new user", username, org, metadata);
    let OrgMSP = capitalizeFirstLetter(org)+"MSP";
    console.log("\n\n OrgMSP::", OrgMSP);
    
    try {
        // load the network configuration
        const ccpPath = path.resolve(__dirname, '..', '..', 'test-network', 'organizations', 'peerOrganizations', ''+org+'.example.com', 'connection-'+org+'.json');
        const ccp = JSON.parse(fs.readFileSync(ccpPath, 'utf8'));

        // Create a new CA client for interacting with the CA.
        const caURL = ccp.certificateAuthorities['ca.'+org+'.example.com'].url;
        const ca = new FabricCAServices(caURL);

        // Create a new file system based wallet for managing identities.
        const walletPath = path.join(process.cwd(), 'wallet');
        const wallet = await Wallets.newFileSystemWallet(walletPath);
        console.log(`Wallet path: ${walletPath}`);

        // Check to see if we've already enrolled the user.
        const userIdentity = await wallet.get(username);
        if (userIdentity) {
            console.log('An identity for the user "appUser" already exists in the wallet');
            return;
        }

        // Check to see if we've already enrolled the admin user.
        const adminIdentity = await wallet.get('admin');
        if (!adminIdentity) {
            console.log('An identity for the admin user "admin" does not exist in the wallet');
            console.log('Run the enrollAdmin.js application before retrying');
            return;
        }

        // build a user object for authenticating with the CA
        const provider = wallet.getProviderRegistry().getProvider(adminIdentity.type);
        const adminUser = await provider.getUserContext(adminIdentity, 'admin');

        // Register the user, enroll the user, and import the new identity into the wallet.
        const secret = await ca.register({
            affiliation: ''+org+'.department1',
            enrollmentID: username,
            role: 'client',
            attrs: [{
                name: "role",
                value: metadata.role,
				ecert: true
			},{
                name: "orgId",
                value: metadata.orgId,
				ecert: true
			},{
                name: "orgName",
                value: metadata.orgName,
				ecert: true
			},{
                name: "orgType",
                value: metadata.orgType,
				ecert: true
			}]
        }, adminUser);
        const enrollment = await ca.enroll({
            enrollmentID: username,
            enrollmentSecret: secret,
            attr_reqs: [{
				name: "role",
				optional: false
			},{
				name: "orgId",
				optional: false
			},,{
				name: "orgName",
				optional: false
			},{
				name: "orgType",
				optional: false
			}]
        });
        const x509Identity = {
            credentials: {
                certificate: enrollment.certificate,
                privateKey: enrollment.key.toBytes(),
            },
            mspId: OrgMSP,
            type: 'X.509',
        };
        await wallet.put(username, x509Identity);
        console.log(`Successfully registered and enrolled admin user "${username}" and imported it into the wallet`);
        return username;
    } catch (error) {
        console.error(`Failed to register user "${username}": ${error}`);
        throw error;
    }
}

function capitalizeFirstLetter(string) {
    return string.charAt(0).toUpperCase() + string.slice(1);
}