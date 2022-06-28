#!/usr/bin/env bash
NETWORK_NAMES=$1
set -e
for n in $(echo "$NETWORK_NAMES" | tr "," " "); do
  echo "Running first migration for network: ${n}"
  npx truffle migrate --network=$n --to=1
done
for n in $(echo "$NETWORK_NAMES" | tr "," " "); do
  echo "Running second migration for network: ${n}"
  npx truffle migrate --network=$n --to=2
done