'use strict';

var utils = require('../utils/writer.js');
var assetService = require('../service/AssetService');

module.exports.createAsset = async function createAsset(req, res, next) {
  console.log("createAsset controller started...");
  var body = req.swagger.params['body'].value;
  let user = req.headers['api_key'];

  try {
    const response = await assetService.createAsset(body, user);
    utils.writeJson(res, response);
  } catch (error) {
    console.log("Error:", error)
    utils.writeJson(res, error);
  }
};

module.exports.getAllAssets = async function getAllAssets(req, res, next) {
  try {
    let user = req.headers['api_key'];
    const response = await assetService.getAllAssets(user);
    return utils.writeJson(res, response);
  } catch (error) {
    console.log("Error:", error)
    return utils.writeJson(res, []);
  }
};

module.exports.getAssetById = async function getAssetById(req, res, next) {
  var id = req.swagger.params['id'].value;
  try {
    let user = req.headers['api_key'];
    const response = await assetService.getAssetById(id, user);
    utils.writeJson(res, response);
  } catch (error) {
    console.log("Error:", error)
    utils.writeJson(res, error);
  }
};
