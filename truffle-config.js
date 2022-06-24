const HDWalletProvider = require('@truffle/hdwallet-provider');

const createTruffleConfig = require('./create-config');

module.exports = createTruffleConfig(HDWalletProvider, {
  compilerVersion: '0.8.14'
})
