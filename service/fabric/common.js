
const adminEnrollService = require('./enrollAdmin')

exports.getIdentity = async function (wallet, user) {
    console.log("getIdentity called..");

    // Check to see if we've already enrolled the user.
    let identity = await wallet.get(user);
    if (!identity && user == "admin") {
        console.log("enrolling admin user....");
        await adminEnrollService.enrollAdmin();
        identity = await wallet.get(user);
    }

    if (!identity) {
        console.log(`111An identity for the user ${user} does not exist in the wallet`);
        return;
    }

    return identity;
}
