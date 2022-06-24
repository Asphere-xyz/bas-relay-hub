const ParliaBlockVerifier = artifacts.require("ParliaBlockVerifier");
const RelayHub = artifacts.require("RelayHub");
const SimpleTokenFactory = artifacts.require("SimpleTokenFactory");
const BridgeRouter = artifacts.require("BridgeRouter");
const CrossChainBridge = artifacts.require("CrossChainBridge");

module.exports = async (deployer) => {
  // deploy parlia block verifier as default verifier
  await deployer.deploy(ParliaBlockVerifier, 12, 200)
  const parliaBlockVerifier = await ParliaBlockVerifier.deployed();
  // deploy relay hub
  await deployer.deploy(RelayHub, parliaBlockVerifier.address);
  const relayHub = await RelayHub.deployed();
  // // deploy simple token factory
  await deployer.deploy(SimpleTokenFactory)
  const simpleTokenFactory = await SimpleTokenFactory.deployed();
  // deploy bridge router
  await deployer.deploy(BridgeRouter)
  const bridgeRouter = await BridgeRouter.deployed();
  // deploy cross chain bridge
  await deployer.deploy(CrossChainBridge);
  const crossChainBridge = await CrossChainBridge.deployed();
  // inti cross bridge
  await crossChainBridge.initialize(relayHub.address, relayHub.address, simpleTokenFactory.address, bridgeRouter.address, "BNB", "BNB");
};
