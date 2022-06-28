const Web3 = require('web3');
const {numberToHex, padLeft} = require("web3-utils"),
  {rlp, toBuffer} = require("ethereumjs-util");
const fs = require('fs');

const DEFAULT_CHAIN_PARAMS = {
  '0x38': {
    blockConfirmations: 12,
    epochLength: 200,
    nativeTokenSymbol: 'BNB',
    nativeTokenName: 'BNB',
  },
  '0x61': {
    blockConfirmations: 12,
    epochLength: 200,
    nativeTokenSymbol: 'BNB',
    nativeTokenName: 'BNB',
  },
};

const DEFAULT_BAS_PARAMS = {
  blockConfirmations: 12,
  epochLength: 1200,
  nativeTokenSymbol: 'BAS',
  nativeTokenName: 'BAS',
};

const DEFAULT_RPC_URL = 'https://data-seed-prebsc-1-s1.binance.org:8545/'

const main = async () => {
  const [, , rpcUrl, targetFile] = process.argv;
  const web3 = new Web3(rpcUrl || DEFAULT_RPC_URL);
  const chainId = await web3.eth.getChainId(),
    genesisBlock = await web3.eth.getBlock('0')
  const rawBlock = rlp.encode([
    toBuffer(genesisBlock.parentHash),
    toBuffer(genesisBlock.sha3Uncles),
    toBuffer(genesisBlock.miner),
    toBuffer(genesisBlock.stateRoot),
    toBuffer(genesisBlock.transactionsRoot),
    toBuffer(genesisBlock.receiptsRoot),
    toBuffer(genesisBlock.logsBloom),
    Number(genesisBlock.difficulty),
    Number(genesisBlock.number),
    Number(genesisBlock.gasLimit),
    Number(genesisBlock.gasUsed),
    Number(genesisBlock.timestamp),
    toBuffer(genesisBlock.extraData),
    toBuffer(genesisBlock.mixHash),
    padLeft(genesisBlock.nonce, 8),
  ])
  console.log(`Hex chain id: ${numberToHex(chainId)}`);
  console.log(`Raw genesis block: 0x${rawBlock.toString('hex')}`)
  if (targetFile) {
    console.log(`Dumping chain params to file: ${targetFile}`);
    fs.writeFileSync(targetFile, JSON.stringify({
      chainId: numberToHex(chainId),
      genesisBlock: `0x${rawBlock.toString('hex')}`
    }, null, 2));
  }
};

main().catch(console.error)