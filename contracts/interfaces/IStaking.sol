// SPDX-License-Identifier: GPL-3.0-only
pragma solidity ^0.8.0;

interface IStaking {

    function currentEpoch() external view returns (uint64);

    function nextEpoch() external view returns (uint64);

    function isValidatorActive(address validator) external view returns (bool);

    function isValidator(address validator) external view returns (bool);

    function getValidators() external view returns (address[] memory);
}