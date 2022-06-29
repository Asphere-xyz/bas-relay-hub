const Web3 = require("web3");

const RelayHub = artifacts.require("RelayHub");
const CrossChainBridge = artifacts.require("CrossChainBridge");

const BINDINGS = [
  {rootChain: 'chapel', childChain: 'bas-devnet-1'},
];

const findChainParams = (chainName) => {
  return require(`../params/${chainName}.json`);
}

const getContractForNetwork = async (deployer, contractType, networkName) => {
  if (deployer.network === networkName) {
    return contractType.deployed();
  }
  const network = deployer.networks[networkName];
  if (!network) throw new Error(`There is no provider for network: ${networkName}`)
  const newType = contractType.clone(contractType.toJSON());
  // feels dirty, but what we can do?
  const providers = network.provider().engine['_providers'],
    {rpcUrl} = providers[providers.length - 1];
  const web3 = new Web3(rpcUrl)
  newType.configureNetwork({provider: web3.currentProvider});
  newType.setNetwork(networkName);
  return newType.deployed();
}

module.exports = async (deployer) => {
  if (['test', 'soliditycoverage', 'ganache'].includes(deployer.network)) {
    return
  }

  const bindings = BINDINGS.filter(config => deployer.network === config.rootChain || deployer.network === config.childChain);
  if (!bindings.length) throw new Error(`There is no bindings for network: ${deployer.network}`);

  for (const {rootChain, childChain} of bindings) {
    const rootRelayHub = await getContractForNetwork(deployer, RelayHub, rootChain),
      childRelayHub = await getContractForNetwork(deployer, RelayHub, childChain);
    const rootBridge = await getContractForNetwork(deployer, CrossChainBridge, rootChain),
      childBridge = await getContractForNetwork(deployer, CrossChainBridge, childChain);
    const rootParams = findChainParams(rootChain),
      childParams = findChainParams(childChain)
    if (rootChain === deployer.network) {
      await rootRelayHub.registerCertifiedBAS(childParams.chainId, childParams.genesisBlock, rootBridge.address, childParams.epochLength);
    } else if (childChain === deployer.network) {
      await childRelayHub.registerCertifiedBAS(rootParams.chainId, rootParams.genesisBlock, childBridge.address, rootParams.epochLength);
    }
  }
};
