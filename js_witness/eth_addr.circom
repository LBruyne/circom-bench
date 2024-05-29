pragma circom 2.0.2;

include "./ecdsa/eth_addr.circom";

component main {public [privkey]} = PrivKeyToAddr(64, 4);