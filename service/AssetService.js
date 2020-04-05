'use strict';

const invokeService = require('./fabric/invoke');
const queryService = require('./fabric/query');
/**
 * Create new asset
 * 
 *
 * body Asset asset object that needs to be added to the store
 * no response value expected for this operation
 **/
exports.createAsset = async function (body, user) {
  console.log("createAsset service started..", body, user);

  let invokeReq = {
    funcName: "CreateAsset",
    args: JSON.stringify({
      docType: "asset",
      ...body
    })
  }
  return await invokeService.invoke(invokeReq, user);
}


/**
 * Get all assets
 * 
 *
 * returns Asset
 **/
exports.getAllAssets = async function (user) {
  console.log("getAllAssets started...");

  try {
    let queryReq = {
      funcName: "QueryAllAssets",
      args: {}
    }

    const assets = await queryService.query(queryReq, user);
    console.log("assets::", assets);

    return assets;
  } catch (error) {
    console.log("failed to get all assets. Error:", error);
    throw error;
  }
}


/**
 * Get asset by id
 * 
 *
 * id String The asset id that needs to be fetched.
 * returns Asset
 **/
exports.getAssetById = async function (id, user) {
  let queryReq = {
    funcName: "GetAsset",
    args: {
      id: id
    }
  }

  return await queryService.query(queryReq, user);
}

