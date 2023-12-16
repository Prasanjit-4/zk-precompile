// SPDX-License-Identifier: MIT
pragma solidity >=0.8.0;

interface IEdDSA {

  function verify(
    uint x,
    uint y
  ) external view returns (bool result);
}