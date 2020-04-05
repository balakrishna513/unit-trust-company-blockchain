'use strict';

var utils = require('../utils/writer.js');
var Fund = require('../service/FundService');

module.exports.createFund = async function createFund(req, res, next) {
  var body = req.swagger.params['body'].value;
  let user = req.headers['api_key'];
  try {
    const response = await Fund.createFund(body, user);
    utils.writeJson(res, response);
  } catch (error) {
    console.log("Failed to create fund. Error:", error);
    utils.writeJson(res, error);
  }
};

module.exports.deleteFund = async function deleteFund(req, res, next) {
  var id = req.swagger.params['id'].value;
  let user = req.headers['api_key'];
  try {
    const response = await Fund.deleteFund(id, user);
    utils.writeJson(res, response);
  } catch (error) {
    console.log("Failed to delete fund.", id, error);
    utils.writeJson(res, error);
  }
};

module.exports.getAllFunds = async function getAllFunds(req, res, next) {
  try {
    let user = req.headers['api_key'];
    const response = await Fund.getAllFunds(user);
    utils.writeJson(res, response);
  } catch (error) {
    console.log("Failed to get all funds.", error);
    utils.writeJson(res, []);
  }
};

module.exports.getFundById = async function getFundById(req, res, next) {
  var id = req.swagger.params['id'].value;
  let user = req.headers['api_key'];
  try {
    const response = await Fund.getFundById(id, user);
    utils.writeJson(res, response);
  } catch (error) {
    console.log("Failed to get fund by id.", id, error);
    utils.writeJson(res, error);
  }
};

module.exports.sellFund = async function sellFund(req, res, next) {
  var body = req.swagger.params['body'].value;
  let user = req.headers['api_key'];
  try {
    const response = await Fund.sellFund(body, user)
    utils.writeJson(res, response);
  } catch (error) {
    console.log("Failed to sell fund. Error:", error);
    utils.writeJson(res, error);
  }
};

module.exports.buyFund = async function buyFund(req, res, next) {
  var body = req.swagger.params['body'].value;
  let user = req.headers['api_key'];
  try {
    const response = await Fund.buyFund(body, user)
    utils.writeJson(res, response);
  } catch (error) {
    console.log("Failed to buy fund. Error:", error);
    utils.writeJson(res, error);
  }
};