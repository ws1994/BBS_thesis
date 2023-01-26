# Hyperledger Fabric Source Code 

Modify the source code and re-compile this project will generate modified docker images, which can be used to bring up a Fabric system and further develop the big-data sharing application as shown in this link https://github.com/ws1994/BBS_thesis .

## Integrate the CORE Scheme

This project only shows the code example to integrate CORE encryption scheme to the private data collection mechanism in Fabric. There is no code example for other modifications.

The modified code file list for CORE is as follows:

https://github.com/ws1994/BBS_thesis/blob/fabric/images/peer/Dockerfile

https://github.com/ws1994/BBS_thesis/blob/fabric/gossip/service/gossip_service.go

https://github.com/ws1994/BBS_thesis/blob/fabric/gossip/privdata/pvtdataprovider.go

Note: please search " ws add " to find out the modified and newly added codes in the code file
