'use strict';

var utils = require('../utils/writer.js');
var User = require('../service/UserService');

module.exports.addUser = function addUser (req, res, next) {
  var body = req.swagger.params['body'].value;
  let user = req.headers['api_key'];
  User.addUser(body)
    .then(function (response) {
      utils.writeJson(res, response);
    })
    .catch(function (response) {
      utils.writeJson(res, response);
    });
};

module.exports.approveUser = function approveUser (req, res, next) {
  var body = req.swagger.params['body'].value;
  let user = req.headers['api_key'];
  User.approveUser(body)
    .then(function (response) {
      utils.writeJson(res, response);
    })
    .catch(function (response) {
      utils.writeJson(res, response);
    });
};

module.exports.getAllUsers = function getAllUsers (req, res, next) {
  let user = req.headers['api_key'];
  User.getAllUsers()
    .then(function (response) {
      utils.writeJson(res, response);
    })
    .catch(function (response) {
      utils.writeJson(res, response);
    });
};

module.exports.getUserById = function getUserById (req, res, next) {
  let user = req.headers['api_key'];
  var id = req.swagger.params['id'].value;
  User.getUserById(id)
    .then(function (response) {
      utils.writeJson(res, response);
    })
    .catch(function (response) {
      utils.writeJson(res, response);
    });
};

module.exports.rejectUser = function rejectUser (req, res, next) {
  let user = req.headers['api_key'];
  var body = req.swagger.params['body'].value;
  User.rejectUser(body)
    .then(function (response) {
      utils.writeJson(res, response);
    })
    .catch(function (response) {
      utils.writeJson(res, response);
    });
};
