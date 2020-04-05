'use strict';

const invokeService = require('./fabric/invoke');
const queryService = require('./fabric/query');

/**
 * Create new fund
 * 
 *
 * body Fund fund object that needs to be added to the store
 * no response value expected for this operation
 **/
exports.createFund = async function(body, user) {
  console.log("createFund service started..", body, user);

  let invokeReq = {
    funcName: "CreateFund",
    args: JSON.stringify({
      docType: "fund",
      ...body
    })
  }
  return await invokeService.invoke(invokeReq, user);
}


/**
 * Delete fund by ID
 * Delete fund by ID
 *
 * id String ID of the fund that needs to be deleted
 * no response value expected for this operation
 **/
exports.deleteFund = async function(id, user) {
  console.log("deleteFund service started..", id, user);

  let invokeReq = {
    funcName: "DeleteFund",
    args: JSON.stringify({
      id: id
    })
  }
  return await invokeService.invoke(invokeReq, user);
}


/**
 * Get all funds
 * 
 *
 * id String The fund id that needs to be sold.
 * returns Fund
 **/
exports.getAllFunds = async function(user) {
  console.log("getAllFunds service started...");
  
  try {
    let queryReq = {
      funcName: "QueryAllFunds",
      args: {}
    }

    const funds = await queryService.query(queryReq, user);
    console.log("funds::", funds);
    return funds;
  } catch (error) {
    console.log("failed to get all funds. Error:", error);
    throw error;
  }
}


/**
 * Get fund by id
 * 
 *
 * id String The fund id that needs to be fetched.
 * returns Fund
 **/
exports.getFundById = async function(id, user) {
  try {
    let queryReq = {
      funcName: "GetFund",
      args: {
        id: id
      }
    }
  
    return await queryService.query(queryReq, user);
  } catch (error) {
    console.log("Failed to get fund by id:", id, error);
    throw error;
  }
}


/**
 * sell fund
 * 
 *
 * id String The fund id that needs to be sold.
 * no response value expected for this operation
 **/
exports.sellFund = async function(body, user) {
  console.log("sellFund service started..", body, user);

  let invokeReq = {
    funcName: "SellFund",
    args: JSON.stringify({
      ...body
    })
  }
  return await invokeService.invoke(invokeReq, user);
}

/**
 * buy fund
 * 
 *
 * id String The fund id that needs to buy.
 * no response value expected for this operation
 **/
exports.buyFund = async function(body, user) {
  console.log("buyFund service started..", body, user);

  let invokeReq = {
    funcName: "BuyFund",
    args: JSON.stringify({
      ...body
    })
  }
  return await invokeService.invoke(invokeReq, user);
}
