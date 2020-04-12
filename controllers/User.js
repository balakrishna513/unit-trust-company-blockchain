'use strict';

var utils = require('../utils/writer.js');
var User = require('../service/UserService');

module.exports.addUser = async function addUser(req, res, next) {
  var body = req.swagger.params['body'].value;
  let user = req.headers['api_key'];

  try {
    const response = await User.addUser(body, user);
    utils.writeJson(res, response);
  } catch (error) {
    console.log("Error:", error)
    utils.writeJson(res, error);
  }
};

module.exports.approveUser = async function approveUser(req, res, next) {
  var body = req.swagger.params['body'].value;
  let user = req.headers['api_key'];

  try {
    const response = await User.approveUser(body, user);
    utils.writeJson(res, response);
  } catch (error) {
    console.log("Error:", error)
    utils.writeJson(res, error);
  }
};

module.exports.getAllUsers = async function getAllUsers(req, res, next) {
  let user = req.headers['api_key'];
  try {
    const response = await User.getAllUsers(user);
    return utils.writeJson(res, response);
  } catch (error) {
    console.log("Error:", error)
    return utils.writeJson(res, []);
  }
};

module.exports.getUserById = async function getUserById(req, res, next) {
  let user = req.headers['api_key'];
  var id = req.swagger.params['id'].value;
  try {
    const response = await User.getUserById(id, user);
    utils.writeJson(res, response);
  } catch (error) {
    console.log("Error:", error)
    utils.writeJson(res, error);
  }
};

module.exports.rejectUser = async function rejectUser(req, res, next) {
  let user = req.headers['api_key'];
  var body = req.swagger.params['body'].value;
  try {
    const response = await User.rejectUser(body, user);
    utils.writeJson(res, response);
  } catch (error) {
    console.log("Error:", error)
    utils.writeJson(res, error);
  }
};
