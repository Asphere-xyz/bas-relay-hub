/** @var artifacts {Array} */
/** @var web3 {Web3} */
/** @function contract */
/** @function it */
/** @function before */
/** @var assert */

const ParliaBlockVerifier = artifacts.require("ParliaBlockVerifier");

contract("ParliaBlockVerifier", async (accounts) => {
  const [owner] = accounts
  it("system fee is well calculated", async () => {
    const verifier = await ParliaBlockVerifier.new()
    const result = await debug(verifier.extractSigningData('0xf903c341a0e06d1e696e78bef7671cc0936b4428d9f6b3aa5ef1dd1b21620fca6742d6faf1a01dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d493479470f657164e5b75689b64b7fd1fa275f334f28e18a04f954c11b1721ed1816b4eef4a57ba99c7288f1a5dbaf09f375efbb888bc6a8ea088d550352ebf720520dc48b2a16a82da57524145862bbedd46f2d22166bf7069a0ff6fc0e069e4e43b42d84b218d3a8bbead88cc5f06044c8f6cb64fa446520c26b90100fffff77feefbffffffbffb7dfffffffffffffffebfffff7ffef5fffffefb7ffbfeffffff75fdfffb7fffbfbfdf5fffeffff7ffffffcf7f7fffffb7bdffffffffffbfffffbfeffffffff7ffddb7ffffefffff9ffffffffedfffffffffffbffffffff5fffff3fff7fbeeef3f7ffdefffffbbbfffebefbffffffcfb3fffbfdffdffbffdffeedffffffff7ffffddffffffbff7fdfc5fb7f6dffdfbffef7fff6dffffffffffbefffffff96efff7bdff7fefdfbfbfde7f7fcfdffdffffbdffffeffbfff7b7ffffbffdd5de7ffbfffdffffd7bffffdf7ffef7fffffffeffffffffffff3fdffffedfdfdfddded7fd3fef7fedfdffffbffff5fffffffff3bdebfffffffbf0283c79d908405f5e1008405ee8b078461a65ec5b901c4d883010105846765746888676f312e31372e32856c696e7578000000c3167bdf2465176c461afb316ebc773c61faee85a6515daa295e26495cef6f69dfa69911d9d8e4f3bbadb89b29a97c6effb8a411dabc6adeefaa84f5067c8bbe2a7cdd959bfe8d9487b2a43b33565295a698f7e22d4c407bbe49438ed859fe965b140dcf1aab71a93f349bbafec1551819b8be1efea2fc46ca749aa14430b3230294d12c6ab2aac5c2cd68e80b16b581685b1ded8013785d6623cc18d214320b6bb6475970f657164e5b75689b64b7fd1fa275f334f28e187ae2f5b9e386cd1b50a4550696d957cb4900f03a8b6c8fd93d6f4cea42bbb345dbc6f0dfdb5bec739bb832254baf4e8b4cc26bd2b52b31389b56e98b9f8ccdafcc39f3c7d6ebf637c9151673cbc36b88a6f79b60359f141df90a0c745125b131caaffd12b8f7166496996a7da21cf1f1b04d9b3e26a3d077be807dddb074639cd9fa61b47676c064fc50d62cce2fd7544e0b2cc94692d4a704debef7bcb61328e2d3a739effcd3a99387d015e260eefac72ebea1e9ae3261a475a27bb1028f140bc2a7c843318afdea0a6e3c511bbd10f4519ece37dc24887e11b55dee226379db83cffc681495730c11fdde79ba4c0ca00000000000000000000000000000000000000000000000000000000000000000880000000000000000'));
    console.log(result)
    console.log(JSON.stringify(result.logs.map(l => l), null, 2))
  })
});
