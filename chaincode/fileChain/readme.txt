./network.sh deployCC -ccn abstore -ccp ../chaincode/abstore/go -ccl go -ccep "OR('Org1MSP.peer','Org2MSP.peer')"


./network.sh deployCC -ccn filechain -ccp ../chaincode/fileChain/go -ccl go -ccep "OR('Org1MSP.peer','Org2MSP.peer')"


./network.sh deployCC -ccn filechain -ccp ../chaincode/fileChain/go -ccv 4 -ccs 4 -ccl go -ccep "OR('Org1MSP.peer','Org2MSP.peer')"



 ./network.sh deployCC -ccn private -ccp ../asset-transfer-private-data/chaincode-go -ccl go -ccep "OR('Org1MSP.peer','Org2MSP.peer')" -cccg '../asset-transfer-private-data/chaincode-go/collections_config.json' -ccep "OR('Org1MSP.peer','Org2MSP.peer')"


export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="Org2MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt
erOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp
export CORE_PEER_ADDRESS=localhost:9051
export ORDERER_CA=${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem


export PEER0_ORG1_CA=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export PEER0_ORG2_CA=${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt


peer chaincode invoke -o localhost:7050 --tls true --cafile $ORDERER_CA -C mychannel -n filechain --peerAddresses localhost:7051 --tlsRootCertFiles ${PEER0_ORG1_CA} --peerAddresses localhost:9051 --tlsRootCertFiles ${PEER0_ORG2_CA} -c '{"Args":["Init","0", "0","0"]}' --waitForEvent


peer chaincode invoke -o localhost:7050 --tls true --cafile $ORDERER_CA -C mychannel -n filechain --peerAddresses localhost:7051 --tlsRootCertFiles ${PEER0_ORG1_CA} --peerAddresses localhost:9051 --tlsRootCertFiles ${PEER0_ORG2_CA} -c '{"Args":["uploadFile","file1", "color","type","size"]}' --waitForEvent


peer chaincode invoke -o localhost:7050 --tls true --cafile $ORDERER_CA -C mychannel -n filechain --peerAddresses localhost:7051 --tlsRootCertFiles ${PEER0_ORG1_CA} --peerAddresses localhost:9051 --tlsRootCertFiles ${PEER0_ORG2_CA} -c '{"Args":["uploadFile","file2", "color2","type2","size2"]}' --waitForEvent


peer chaincode query -C mychannel -n filechain -c '{"Args":["query","statistical"]}'

peer chaincode query -C mychannel -n filechain -c '{"Args":["query","uploadFile_file1"]}'


peer chaincode invoke -o localhost:7050 --tls true --cafile $ORDERER_CA -C mychannel -n filechain --peerAddresses localhost:7051 --tlsRootCertFiles ${PEER0_ORG1_CA} --peerAddresses localhost:9051 --tlsRootCertFiles ${PEER0_ORG2_CA} -c '{"Args":["readFile","file1", "color","type","size","num1","num2"]}' --waitForEvent


peer chaincode invoke -o localhost:7050 --tls true --cafile $ORDERER_CA -C mychannel -n filechain --peerAddresses localhost:7051 --tlsRootCertFiles ${PEER0_ORG1_CA} --peerAddresses localhost:9051 --tlsRootCertFiles ${PEER0_ORG2_CA} -c '{"Args":["downloadFile","file1", "color","type","size","num1","num2","num3","num4"]}' --waitForEvent


peer chaincode invoke -o localhost:7050 --tls true --cafile $ORDERER_CA -C mychannel -n filechain --peerAddresses localhost:7051 --tlsRootCertFiles ${PEER0_ORG1_CA} --peerAddresses localhost:9051 --tlsRootCertFiles ${PEER0_ORG2_CA} -c '{"Args":["deleteReadFile","file1&file2"]}' --waitForEvent

