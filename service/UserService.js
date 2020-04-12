'use strict';

const invokeService = require('./fabric/invoke');
const queryService = require('./fabric/query');
const fabricUserService = require('./fabric/registerUser')

/**
 * Add a new user
 * 
 *
 * body User User object that needs to store
 * no response value expected for this operation
 **/
exports.addUser = async function (body, user) {
  console.log("addUser service started..", body, user);

  let invokeReq = {
    funcName: "RegisterUser",
    args: JSON.stringify({
      docType: "user",
      ...body
    })
  }

  await fabricUserService.registerUser(body.id, body.orgId, body);

  return await invokeService.invoke(invokeReq, user);
}


/**
 * Approve user
 * 
 *
 * body  User that needs to approved (optional)
 * no response value expected for this operation
 **/
exports.approveUser = async function (body, user) {
  console.log("approveUser service started..", body, user);

  let invokeReq = {
    funcName: "ApproveUser",
    args: JSON.stringify({
      ...body
    })
  }
  return await invokeService.invoke(invokeReq, user);
}


/**
 * Get all users
 * 
 *
 * returns User
 **/
exports.getAllUsers = async function (user) {
  console.log("getAllUsers started...");

  try {
    let queryReq = {
      funcName: "QueryAllUsers",
      args: {}
    }

    const assets = await queryService.query(queryReq, user);
    console.log("users::", assets);

    return assets;
  } catch (error) {
    console.log("failed to get all users. Error:", error);
    throw error;
  }
}


/**
 * Get user by user id
 * 
 *
 * id String The user id that needs to be fetched.
 * returns User
 **/
exports.getUserById = async function (id, user) {
  let queryReq = {
    funcName: "GetUser",
    args: {
      id: id
    }
  }

  return await queryService.query(queryReq, user);
}


/**
 * Reject user
 * 
 *
 * body  User that needs to be rejected (optional)
 * no response value expected for this operation
 **/
exports.rejectUser = async function (body, user) {
  console.log("rejectUser service started..", body, user);

  let invokeReq = {
    funcName: "RejectUser",
    args: JSON.stringify({
      ...body
    })
  }
  return await invokeService.invoke(invokeReq, user);
}