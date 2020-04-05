'use strict';

const invokeService = require('./fabric/invoke');
const queryService = require('./fabric/query');

/**
 * Add a new user
 * 
 *
 * body User User object that needs to store
 * no response value expected for this operation
 **/
exports.addUser = function (body, user) {
  return new Promise(function (resolve, reject) {
    resolve();
  });
}


/**
 * Approve user
 * 
 *
 * body  User that needs to approved (optional)
 * no response value expected for this operation
 **/
exports.approveUser = function (body, user) {
  return new Promise(function (resolve, reject) {
    resolve();
  });
}


/**
 * Get all users
 * 
 *
 * returns User
 **/
exports.getAllUsers = function (user) {
  return new Promise(function (resolve, reject) {
    var examples = {};
    examples['application/json'] = {
      "orgType": "orgType",
      "orgName": "orgName",
      "name": "name",
      "id": "id",
      "userType": "userType",
      "orgId": "orgId"
    };
    if (Object.keys(examples).length > 0) {
      resolve(examples[Object.keys(examples)[0]]);
    } else {
      resolve();
    }
  });
}


/**
 * Get user by user id
 * 
 *
 * id String The user id that needs to be fetched.
 * returns User
 **/
exports.getUserById = function (id) {
  return new Promise(function (resolve, reject) {
    var examples = {};
    examples['application/json'] = {
      "orgType": "orgType",
      "orgName": "orgName",
      "name": "name",
      "id": "id",
      "userType": "userType",
      "orgId": "orgId"
    };
    if (Object.keys(examples).length > 0) {
      resolve(examples[Object.keys(examples)[0]]);
    } else {
      resolve();
    }
  });
}


/**
 * Reject user
 * 
 *
 * body  User that needs to be rejected (optional)
 * no response value expected for this operation
 **/
exports.rejectUser = function (body, user) {
  return new Promise(function (resolve, reject) {
    resolve();
  });
}