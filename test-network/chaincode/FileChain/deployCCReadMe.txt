./network.sh down
./network.sh up createChannel -s couchdb -i 2.0.0


cd ./chaincode/FileChain/go

GO111MODULE=on go mod vendor

cd ../../../


export PATH=${PWD}/../bin:${PWD}:$PATH
export FABRIC_CFG_PATH=$PWD/../config/

export PEER0_ORG1_CA=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export PEER0_ORG2_CA=${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt

export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="Org1MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_ADDRESS=localhost:7051
export ORDERER_CA=${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem


peer lifecycle chaincode package filechain.tar.gz --lang golang --path ./chaincode/FileChain/go --label filechian_1

peer lifecycle chaincode install filechain.tar.gz


export version=1
export packageID=filechian_1:33659edfba278a976e6f5488ea7a1fd445028f81cc1a67cb60223d75a4e4764e


peer lifecycle chaincode approveformyorg --orderer localhost:7050 --ordererTLSHostnameOverride orderer.example.com --channelID mychannel --name filechain -v $version --package-id $packageID --sequence $version --tls --cafile $ORDERER_CA --peerAddresses localhost:7051 --tlsRootCertFiles ${PEER0_ORG1_CA}



export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="Org2MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp
export CORE_PEER_ADDRESS=localhost:9051
export ORDERER_CA=${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem


peer lifecycle chaincode install filechain.tar.gz

peer lifecycle chaincode approveformyorg --orderer localhost:7050 --ordererTLSHostnameOverride orderer.example.com --channelID mychannel --name filechain -v $version --package-id $packageID --sequence $version --tls --cafile $ORDERER_CA --peerAddresses localhost:9051 --tlsRootCertFiles ${PEER0_ORG2_CA}


export PEER0_ORG1_CA=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export PEER0_ORG2_CA=${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt

peer lifecycle chaincode commit -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --peerAddresses localhost:7051 --tlsRootCertFiles ${PEER0_ORG1_CA} --peerAddresses localhost:9051 --tlsRootCertFiles ${PEER0_ORG2_CA} --channelID mychannel --name filechain -v $version  --sequence $version --tls --cafile $ORDERER_CA --waitForEvent



%********************
% Init
%********************

peer chaincode invoke -o localhost:7050 --tls true --cafile $ORDERER_CA -C mychannel -n filechain --peerAddresses localhost:7051 --tlsRootCertFiles ${PEER0_ORG1_CA} --peerAddresses localhost:9051 --tlsRootCertFiles ${PEER0_ORG2_CA} -c '{"Args":["Init","0", "0","0"]}' --waitForEvent

peer chaincode query -C mychannel -n filechain -c '{"Args":["Get","ip.peer0.org2"]}'


peer chaincode invoke -o localhost:7050 --tls true --cafile $ORDERER_CA -C mychannel -n filechain --peerAddresses localhost:7051 --tlsRootCertFiles ${PEER0_ORG1_CA} --peerAddresses localhost:10051 --tlsRootCertFiles ${PEER1_ORG1_CA} --peerAddresses localhost:11051 --tlsRootCertFiles ${PEER2_ORG1_CA} --peerAddresses localhost:12051 --tlsRootCertFiles ${PEER3_ORG1_CA} --peerAddresses localhost:9051 --tlsRootCertFiles ${PEER0_ORG2_CA} -c '{"Args":["InitConfig"]}' --waitForEvent


peer chaincode invoke -o localhost:7050 --tls true --cafile $ORDERER_CA -C mychannel -n filechain --peerAddresses localhost:7051 --tlsRootCertFiles ${PEER0_ORG1_CA} --peerAddresses localhost:9051 --tlsRootCertFiles ${PEER0_ORG2_CA} -c '{"Args":["InitConfig"]}' --waitForEvent

docker ps

export org1CCID=54ea586818ea  
docker exec -u root -it $org1CCID /bin/sh
cd /home/chaincode/
chmod 777 SSHShell.sh
./SSHShell.sh
exit

export org2CCID=880bfc8fb9c8
docker exec -u root -it $org2CCID /bin/sh
cd /home/chaincode/
chmod 777 SSHShell.sh
./SSHShell.sh
exit


sftp root@192.168.240.8:/home/chaincode/off_state <<< $'put /home/cgao/Hyperledger/fabric-samples/test-network/2_6g.mp4'



peer chaincode invoke -o localhost:7050 --tls true --cafile $ORDERER_CA -C mychannel -n bigFileTransfer --peerAddresses localhost:7051 --tlsRootCertFiles ${PEER0_ORG1_CA} -c '{"Args":["InitRSAKeyPair","Org1MSP"]}' --waitForEvent


peer chaincode invoke -o localhost:7050 --tls true --cafile $ORDERER_CA -C mychannel -n bigFileTransfer --peerAddresses localhost:9051 --tlsRootCertFiles ${PEER0_ORG2_CA} -c '{"Args":["InitRSAKeyPair","Org2MSP"]}' --waitForEvent


peer chaincode query -C mychannel -n bigFileTransfer --peerAddresses localhost:7051 --tlsRootCertFiles ${PEER0_ORG1_CA} -c '{"Args":["GetPvtData","Org1MSP_RSA","_implicit_org_Org1MSP"]}'


peer chaincode query -C mychannel -n bigFileTransfer --peerAddresses localhost:9051 --tlsRootCertFiles ${PEER0_ORG2_CA} -c '{"Args":["GetPvtData","Org2MSP_RSA","_implicit_org_Org2MSP"]}'



export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="Org1MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_ADDRESS=localhost:7051

peer chaincode invoke -o localhost:7050 --tls true --cafile $ORDERER_CA -C mychannel -n bigFileTransfer --peerAddresses localhost:7051 --tlsRootCertFiles ${PEER0_ORG1_CA} --peerAddresses localhost:9051 --tlsRootCertFiles ${PEER0_ORG2_CA} -c '{"Args":["SetKeyPolicy","ip.peer0.org1","Org1MSP"]}' --waitForEvent

peer chaincode invoke -o localhost:7050 --tls true --cafile $ORDERER_CA -C mychannel -n bigFileTransfer --peerAddresses localhost:10051 --tlsRootCertFiles ${PEER1_ORG1_CA} --peerAddresses localhost:9051 --tlsRootCertFiles ${PEER0_ORG2_CA} -c '{"Args":["SetKeyPolicy","ip.peer1.org1","Org1MSP"]}' --waitForEvent

peer chaincode invoke -o localhost:7050 --tls true --cafile $ORDERER_CA -C mychannel -n bigFileTransfer --peerAddresses localhost:11051 --tlsRootCertFiles ${PEER2_ORG1_CA} --peerAddresses localhost:9051 --tlsRootCertFiles ${PEER0_ORG2_CA} -c '{"Args":["SetKeyPolicy","ip.peer2.org1","Org1MSP"]}' --waitForEvent

peer chaincode invoke -o localhost:7050 --tls true --cafile $ORDERER_CA -C mychannel -n bigFileTransfer --peerAddresses localhost:12051 --tlsRootCertFiles ${PEER3_ORG1_CA} --peerAddresses localhost:9051 --tlsRootCertFiles ${PEER0_ORG2_CA} -c '{"Args":["SetKeyPolicy","ip.peer3.org1","Org1MSP"]}' --waitForEvent

peer chaincode query -C mychannel -n bigFileTransfer -c '{"Args":["GetKeyPolicy","ip.peer0.org1"]}'

peer chaincode query -C mychannel -n bigFileTransfer -c '{"Args":["GetKeyPolicy","ip.peer1.org1"]}'

peer chaincode query -C mychannel -n bigFileTransfer -c '{"Args":["GetKeyPolicy","ip.peer2.org1"]}'

peer chaincode query -C mychannel -n bigFileTransfer -c '{"Args":["GetKeyPolicy","ip.peer3.org1"]}'

peer chaincode invoke -o localhost:7050 --tls true --cafile $ORDERER_CA -C mychannel -n bigFileTransfer --peerAddresses localhost:7051 --tlsRootCertFiles ${PEER0_ORG1_CA} -c '{"Args":["SetIP","ip.peer0.org1"]}' --waitForEvent

peer chaincode invoke -o localhost:7050 --tls true --cafile $ORDERER_CA -C mychannel -n bigFileTransfer --peerAddresses localhost:10051 --tlsRootCertFiles ${PEER1_ORG1_CA} -c '{"Args":["SetIP","ip.peer1.org1"]}' --waitForEvent

peer chaincode invoke -o localhost:7050 --tls true --cafile $ORDERER_CA -C mychannel -n bigFileTransfer --peerAddresses localhost:11051 --tlsRootCertFiles ${PEER2_ORG1_CA} -c '{"Args":["SetIP","ip.peer2.org1"]}' --waitForEvent

peer chaincode invoke -o localhost:7050 --tls true --cafile $ORDERER_CA -C mychannel -n bigFileTransfer --peerAddresses localhost:12051 --tlsRootCertFiles ${PEER3_ORG1_CA} -c '{"Args":["SetIP","ip.peer3.org1"]}' --waitForEvent

peer chaincode query -C mychannel -n bigFileTransfer -c '{"Args":["GetIpInfo","ip.peer0.org1"]}'

peer chaincode query -C mychannel -n bigFileTransfer -c '{"Args":["GetIpInfo","ip.peer1.org1"]}'

peer chaincode query -C mychannel -n bigFileTransfer -c '{"Args":["GetIpInfo","ip.peer2.org1"]}'

peer chaincode query -C mychannel -n bigFileTransfer -c '{"Args":["GetIpInfo","ip.peer3.org1"]}'




export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="Org2MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp
export CORE_PEER_ADDRESS=localhost:9051

peer chaincode invoke -o localhost:7050 --tls true --cafile $ORDERER_CA -C mychannel -n bigFileTransfer --peerAddresses localhost:7051 --tlsRootCertFiles ${PEER0_ORG1_CA} --peerAddresses localhost:9051 --tlsRootCertFiles ${PEER0_ORG2_CA} -c '{"Args":["SetKeyPolicy","ip.peer0.org2","Org2MSP"]}' --waitForEvent

peer chaincode query -C mychannel -n bigFileTransfer -c '{"Args":["GetKeyPolicy","ip.peer0.org2"]}'

peer chaincode invoke -o localhost:7050 --tls true --cafile $ORDERER_CA -C mychannel -n bigFileTransfer --peerAddresses localhost:9051 --tlsRootCertFiles ${PEER0_ORG2_CA} -c '{"Args":["SetIP","ip.peer0.org2"]}' --waitForEvent

peer chaincode query -C mychannel -n bigFileTransfer -c '{"Args":["GetIpInfo","ip.peer0.org2"]}'




%********************
% Application small txt file
%********************


peer chaincode invoke -o localhost:7050 --tls true --cafile $ORDERER_CA -C mychannel -n bigFileTransfer2 --peerAddresses localhost:7051 --tlsRootCertFiles ${PEER0_ORG1_CA} --peerAddresses localhost:9051 --tlsRootCertFiles ${PEER0_ORG2_CA} -c '{"Args":["Set","bigFile.txt","initFile"]}'

peer chaincode invoke -o localhost:7050 --tls true --cafile $ORDERER_CA -C mychannel -n bigFileTransfer2 --peerAddresses localhost:7051 --tlsRootCertFiles ${PEER0_ORG1_CA} --peerAddresses localhost:9051 --tlsRootCertFiles ${PEER0_ORG2_CA} -c '{"Args":["SetKeyPolicy","bigFile.txt","Org1MSP"]}'

peer chaincode invoke -o localhost:7050 --tls true --cafile $ORDERER_CA -C mychannel -n bigFileTransfer2 --peerAddresses localhost:7051 --tlsRootCertFiles ${PEER0_ORG1_CA} -c '{"Args":["GetKeyPolicy","bigFile.txt"]}'

export org1CCID=8573f6cc46f1 
docker cp /home/cgao/Hyperledger/fabric-samples/test-network/bigFile.txt $org1CCID:/home/chaincode/User1@org1.example.com/

peer chaincode query --peerAddresses localhost:7051 --tlsRootCertFiles ${PEER0_ORG1_CA} -C mychannel -n bigFileTransfer2 -c '{"Args":["ListFolder","/home/chaincode/User1@org1.example.com"]}'




peer chaincode invoke -o localhost:7050 --tls true --cafile $ORDERER_CA -C mychannel -n bigFileTransfer --peerAddresses localhost:7051 --tlsRootCertFiles ${PEER0_ORG1_CA} --peerAddresses localhost:9051 --tlsRootCertFiles ${PEER0_ORG2_CA} -c '{"Args":["Set","bigFilev1.txt","initFile"]}'

peer chaincode invoke -o localhost:7050 --tls true --cafile $ORDERER_CA -C mychannel -n bigFileTransfer --peerAddresses localhost:7051 --tlsRootCertFiles ${PEER0_ORG1_CA} --peerAddresses localhost:9051 --tlsRootCertFiles ${PEER0_ORG2_CA} -c '{"Args":["SetKeyPolicy","bigFilev1.txt","Org1MSP"]}'

peer chaincode invoke -o localhost:7050 --tls true --cafile $ORDERER_CA -C mychannel -n bigFileTransfer --peerAddresses localhost:7051 --tlsRootCertFiles ${PEER0_ORG1_CA} -c '{"Args":["GetKeyPolicy","bigFilev1.txt"]}'



peer chaincode invoke -o localhost:7050 --tls true --cafile $ORDERER_CA -C mychannel -n bigFileTransfer --peerAddresses localhost:7051 --tlsRootCertFiles ${PEER0_ORG1_CA} --peerAddresses localhost:9051 --tlsRootCertFiles ${PEER0_ORG2_CA} -c '{"Args":["Set","bigFilev2.txt","initFile"]}'

peer chaincode invoke -o localhost:7050 --tls true --cafile $ORDERER_CA -C mychannel -n bigFileTransfer --peerAddresses localhost:7051 --tlsRootCertFiles ${PEER0_ORG1_CA} --peerAddresses localhost:9051 --tlsRootCertFiles ${PEER0_ORG2_CA} -c '{"Args":["SetKeyPolicy","bigFilev2.txt","Org1MSP"]}'

peer chaincode invoke -o localhost:7050 --tls true --cafile $ORDERER_CA -C mychannel -n bigFileTransfer --peerAddresses localhost:7051 --tlsRootCertFiles ${PEER0_ORG1_CA} -c '{"Args":["GetKeyPolicy","bigFilev2.txt"]}'



peer chaincode invoke -o localhost:7050 --tls true --cafile $ORDERER_CA -C mychannel -n bigFileTransfer --peerAddresses localhost:7051 --tlsRootCertFiles ${PEER0_ORG1_CA} --peerAddresses localhost:9051 --tlsRootCertFiles ${PEER0_ORG2_CA} -c '{"Args":["Set","bigFilev3.txt","initFile"]}'

peer chaincode invoke -o localhost:7050 --tls true --cafile $ORDERER_CA -C mychannel -n bigFileTransfer --peerAddresses localhost:7051 --tlsRootCertFiles ${PEER0_ORG1_CA} --peerAddresses localhost:9051 --tlsRootCertFiles ${PEER0_ORG2_CA} -c '{"Args":["SetKeyPolicy","bigFilev3.txt","Org1MSP"]}'

peer chaincode invoke -o localhost:7050 --tls true --cafile $ORDERER_CA -C mychannel -n bigFileTransfer --peerAddresses localhost:7051 --tlsRootCertFiles ${PEER0_ORG1_CA} -c '{"Args":["GetKeyPolicy","bigFilev3.txt"]}'



peer chaincode invoke -o localhost:7050 --tls true --cafile $ORDERER_CA -C mychannel -n bigFileTransfer --peerAddresses localhost:7051 --tlsRootCertFiles ${PEER0_ORG1_CA} --peerAddresses localhost:9051 --tlsRootCertFiles ${PEER0_ORG2_CA} -c '{"Args":["Set","bigFilev4.txt","initFile"]}'

peer chaincode invoke -o localhost:7050 --tls true --cafile $ORDERER_CA -C mychannel -n bigFileTransfer --peerAddresses localhost:7051 --tlsRootCertFiles ${PEER0_ORG1_CA} --peerAddresses localhost:9051 --tlsRootCertFiles ${PEER0_ORG2_CA} -c '{"Args":["SetKeyPolicy","bigFilev4.txt","Org1MSP"]}'

peer chaincode invoke -o localhost:7050 --tls true --cafile $ORDERER_CA -C mychannel -n bigFileTransfer --peerAddresses localhost:7051 --tlsRootCertFiles ${PEER0_ORG1_CA} -c '{"Args":["GetKeyPolicy","bigFilev4.txt"]}'




peer chaincode invoke -o localhost:7050 --tls true --cafile $ORDERER_CA -C mychannel -n bigFileTransfer2 --peerAddresses localhost:7051 --tlsRootCertFiles ${PEER0_ORG1_CA} -c '{"Args":["AddFileInfo","{\"name\":\"bigFile.txt\",\"filehash\":\"\",\"objecttype\":\"file\",\"description\":\"File1 belong to org1\",\"owner\":\"User1@org1.example.com\",\"rule\":[\"Org2MSP\"]}"]}'

peer chaincode invoke -o localhost:7050 --tls true --cafile $ORDERER_CA -C mychannel -n bigFileTransfer --peerAddresses localhost:10051 --tlsRootCertFiles ${PEER1_ORG1_CA} -c '{"Args":["AddFileInfo","{\"name\":\"bigFilev2.txt\",\"filehash\":\"\",\"objecttype\":\"file\",\"description\":\"File1 belong to org1\",\"owner\":\"Org1MSP\"}"]}'

peer chaincode invoke -o localhost:7050 --tls true --cafile $ORDERER_CA -C mychannel -n bigFileTransfer --peerAddresses localhost:11051 --tlsRootCertFiles ${PEER2_ORG1_CA} -c '{"Args":["AddFileInfo","{\"name\":\"bigFilev3.txt\",\"filehash\":\"\",\"objecttype\":\"file\",\"description\":\"File1 belong to org1\",\"owner\":\"Org1MSP\"}"]}'

peer chaincode invoke -o localhost:7050 --tls true --cafile $ORDERER_CA -C mychannel -n bigFileTransfer --peerAddresses localhost:12051 --tlsRootCertFiles ${PEER3_ORG1_CA} -c '{"Args":["AddFileInfo","{\"name\":\"bigFilev4.txt\",\"filehash\":\"\",\"objecttype\":\"file\",\"description\":\"File1 belong to org1\",\"owner\":\"Org1MSP\"}"]}'


peer chaincode invoke -o localhost:7050 --tls true --cafile $ORDERER_CA -C mychannel -n bigFileTransfer --peerAddresses localhost:7051 --tlsRootCertFiles ${PEER0_ORG1_CA} -c '{"Args":["QueryFileInfo","bigFile.txt"]}'


peer chaincode invoke -o localhost:7050 --tls true --cafile $ORDERER_CA -C mychannel -n bigFileTransfer --peerAddresses localhost:7051 --tlsRootCertFiles ${PEER0_ORG1_CA} -c '{"Args":["ListAllFiles"]}'




peer chaincode invoke -o localhost:7050 --tls true --cafile $ORDERER_CA -C mychannel -n bigFileTransfer --peerAddresses localhost:7051 --tlsRootCertFiles ${PEER0_ORG1_CA} --peerAddresses localhost:9051 --tlsRootCertFiles ${PEER0_ORG2_CA} -c '{"Args":["RequestSftp","bigFile.txt","Org1MSP","Org2MSP"]}'

peer chaincode query -C mychannel -n bigFileTransfer -c '{"Args":["Get","bigFile.txtOrg1MSPOrg2MSP"]}'



peer chaincode invoke -o localhost:7050 --tls true --cafile $ORDERER_CA -C mychannel -n bigFileTransfer --peerAddresses localhost:7051 --tlsRootCertFiles ${PEER0_ORG1_CA} -c '{"Args":["SendFile_SSH_sftp","bigFilev1.txtOrg1MSPOrg2MSP","bigFilev1.txt","org2"]}'

peer chaincode invoke -o localhost:7050 --tls true --cafile $ORDERER_CA -C mychannel -n bigFileTransfer --peerAddresses localhost:10051 --tlsRootCertFiles ${PEER1_ORG1_CA} -c '{"Args":["SendFile_SSH_sftp","bigFilev2.txtOrg1MSPOrg2MSP","bigFilev2.txt","org2"]}'

peer chaincode invoke -o localhost:7050 --tls true --cafile $ORDERER_CA -C mychannel -n bigFileTransfer --peerAddresses localhost:11051 --tlsRootCertFiles ${PEER2_ORG1_CA} -c '{"Args":["SendFile_SSH_sftp","bigFilev3.txtOrg1MSPOrg2MSP","bigFilev3.txt","org2"]}'

peer chaincode invoke -o localhost:7050 --tls true --cafile $ORDERER_CA -C mychannel -n bigFileTransfer --peerAddresses localhost:12051 --tlsRootCertFiles ${PEER3_ORG1_CA} -c '{"Args":["SendFile_SSH_sftp","bigFilev4.txtOrg1MSPOrg2MSP","bigFilev4.txt","org2"]}'



peer chaincode query --peerAddresses localhost:9051 --tlsRootCertFiles ${PEER0_ORG2_CA} -C mychannel -n bigFileTransfer -c '{"Args":["CheckReceivedFile","bigFile.txt"]}'


peer chaincode invoke -o localhost:7050 --tls true --cafile $ORDERER_CA -C mychannel -n bigFileTransfer --peerAddresses localhost:7051 --tlsRootCertFiles ${PEER0_ORG1_CA} -c '{"Args":["RequestKey","bigFile.txtOrg1MSPOrg2MSP"]}'


peer chaincode query --peerAddresses localhost:9051 --tlsRootCertFiles ${PEER0_ORG2_CA} -C mychannel -n bigFileTransfer -c '{"Args":["QueryKeyAndDecryptFile","bigFile.txt","bigFile.txtOrg1MSPOrg2MSP"]}'





%*****************************
% Application image 77.9M
%*****************************

peer chaincode invoke -o localhost:7050 --tls true --cafile $ORDERER_CA -C mychannel -n bigFileTransfer --peerAddresses localhost:7051 --tlsRootCertFiles ${PEER0_ORG1_CA} --peerAddresses localhost:9051 --tlsRootCertFiles ${PEER0_ORG2_CA} -c '{"Args":["Set","Rheinfall.jpg","initFile"]}' --waitForEvent

peer chaincode invoke -o localhost:7050 --tls true --cafile $ORDERER_CA -C mychannel -n bigFileTransfer --peerAddresses localhost:7051 --tlsRootCertFiles ${PEER0_ORG1_CA} --peerAddresses localhost:9051 --tlsRootCertFiles ${PEER0_ORG2_CA} -c '{"Args":["SetKeyPolicy","Rheinfall.jpg","Org1MSP"]}' --waitForEvent

peer chaincode invoke -o localhost:7050 --tls true --cafile $ORDERER_CA -C mychannel -n bigFileTransfer --peerAddresses localhost:7051 --tlsRootCertFiles ${PEER0_ORG1_CA} -c '{"Args":["GetKeyPolicy","Rheinfall.jpg"]}'


peer chaincode query -C mychannel -n bigFileTransfer -c '{"Args":["GetIpInfo","ip.peer0.org1"]}'
sftp root@192.168.112.7:/home/chaincode/off_state <<< $'put /home/cgao/Hyperledger/fabric-samples/test-network/Rheinfall.jpg'


peer chaincode query --peerAddresses localhost:7051 --tlsRootCertFiles ${PEER0_ORG1_CA} -C mychannel -n bigFileTransfer -c '{"Args":["ListFolder","/home/chaincode/off_state"]}'


peer chaincode invoke -o localhost:7050 --tls true --cafile $ORDERER_CA -C mychannel -n bigFileTransfer --peerAddresses localhost:7051 --tlsRootCertFiles ${PEER0_ORG1_CA} -c '{"Args":["AddFileInfo","{\"name\":\"Rheinfall.jpg\",\"filehash\":\"\",\"objecttype\":\"file\",\"description\":\"image 77.9M belong to org1\",\"owner\":\"Org1MSP\"}"]}' --waitForEvent

peer chaincode invoke -o localhost:7050 --tls true --cafile $ORDERER_CA -C mychannel -n bigFileTransfer --peerAddresses localhost:7051 --tlsRootCertFiles ${PEER0_ORG1_CA} -c '{"Args":["QueryFileInfo","Rheinfall.jpg"]}'




peer chaincode invoke -o localhost:7050 --tls true --cafile $ORDERER_CA -C mychannel -n bigFileTransfer --peerAddresses localhost:7051 --tlsRootCertFiles ${PEER0_ORG1_CA} -c '{"Args":["ListAllFiles"]}'




peer chaincode invoke -o localhost:7050 --tls true --cafile $ORDERER_CA -C mychannel -n bigFileTransfer --peerAddresses localhost:7051 --tlsRootCertFiles ${PEER0_ORG1_CA} --peerAddresses localhost:9051 --tlsRootCertFiles ${PEER0_ORG2_CA} -c '{"Args":["RequestSftp","Rheinfall.jpg","Org1MSP","Org2MSP"]}' --waitForEvent

peer chaincode query -C mychannel -n bigFileTransfer -c '{"Args":["Get","Rheinfall.jpgOrg1MSPOrg2MSP"]}'


peer chaincode invoke -o localhost:7050 --tls true --cafile $ORDERER_CA -C mychannel -n bigFileTransfer --peerAddresses localhost:7051 --tlsRootCertFiles ${PEER0_ORG1_CA} -c '{"Args":["SendFile_SSH_sftp","Rheinfall.jpgOrg1MSPOrg2MSP","Rheinfall.jpg","org2"]}' --waitForEvent


peer chaincode query --peerAddresses localhost:9051 --tlsRootCertFiles ${PEER0_ORG2_CA} -C mychannel -n bigFileTransfer -c '{"Args":["CheckReceivedFile","Rheinfall.jpg"]}'


peer chaincode invoke -o localhost:7050 --tls true --cafile $ORDERER_CA -C mychannel -n bigFileTransfer --peerAddresses localhost:7051 --tlsRootCertFiles ${PEER0_ORG1_CA} -c '{"Args":["RequestKey","Rheinfall.jpgOrg1MSPOrg2MSP"]}' --waitForEvent


peer chaincode invoke -o localhost:7050 --tls true --cafile $ORDERER_CA -C mychannel -n bigFileTransfer --peerAddresses localhost:9051 --tlsRootCertFiles ${PEER0_ORG2_CA} -C mychannel -n bigFileTransfer -c '{"Args":["QueryKeyAndDecryptFile","Rheinfall.jpg","Rheinfall.jpgOrg1MSPOrg2MSP"]}' --waitForEvent


peer chaincode query --peerAddresses localhost:9051 --tlsRootCertFiles ${PEER0_ORG2_CA} -C mychannel -n bigFileTransfer -c '{"Args":["ListFolder","/home/chaincode/off_state"]}'


export org2CCID=8c1134a5daf3
docker exec -it -u root $org2CCID /bin/sh

export org2CCID=7fd6e934ed6b
docker cp $org2CCID:/home/chaincode/sftp_Shell.sh /home/cgao/Hyperledger/fabric-samples/test-network/






%*****************************
% Application image 576M  wac_nearside.tif
%*****************************


peer chaincode invoke -o localhost:7050 --tls true --cafile $ORDERER_CA -C mychannel -n bigFileTransfer --waitForEvent --peerAddresses localhost:7051 --tlsRootCertFiles ${PEER0_ORG1_CA} --peerAddresses localhost:9051 --tlsRootCertFiles ${PEER0_ORG2_CA} -c '{"Args":["Set","wac_nearside.tif","initFile"]}' --waitForEvent

peer chaincode invoke -o localhost:7050 --tls true --cafile $ORDERER_CA -C mychannel -n bigFileTransfer --waitForEvent --peerAddresses localhost:7051 --tlsRootCertFiles ${PEER0_ORG1_CA} --peerAddresses localhost:9051 --tlsRootCertFiles ${PEER0_ORG2_CA} -c '{"Args":["SetKeyPolicy","wac_nearside.tif","Org1MSP"]}' --waitForEvent

peer chaincode invoke -o localhost:7050 --tls true --cafile $ORDERER_CA -C mychannel -n bigFileTransfer --peerAddresses localhost:7051 --tlsRootCertFiles ${PEER0_ORG1_CA} -c '{"Args":["GetKeyPolicy","wac_nearside.tif"]}' --waitForEvent


peer chaincode query -C mychannel -n bigFileTransfer -c '{"Args":["GetIpInfo","ip.peer0.org1"]}'
sftp root@192.168.112.7:/home/chaincode/off_state <<< $'put /home/cgao/Hyperledger/fabric-samples/test-network/wac_nearside.tif'


peer chaincode query --peerAddresses localhost:7051 --tlsRootCertFiles ${PEER0_ORG1_CA} -C mychannel -n bigFileTransfer -c '{"Args":["ListFolder","/home/chaincode/off_state"]}'


peer chaincode invoke -o localhost:7050 --tls true --cafile $ORDERER_CA -C mychannel -n bigFileTransfer --waitForEvent --peerAddresses localhost:7051 --tlsRootCertFiles ${PEER0_ORG1_CA} -c '{"Args":["AddFileInfo","{\"name\":\"wac_nearside.tif\",\"filehash\":\"\",\"objecttype\":\"file\",\"description\":\"image wac_nearside.tif 576M belong to org1\",\"owner\":\"Org1MSP\"}"]}'

peer chaincode query -C mychannel -n bigFileTransfer -c '{"Args":["QueryFileInfo","wac_nearside.tif"]}'

peer chaincode query -C mychannel -n bigFileTransfer -c '{"Args":["ListAllFiles"]}'



peer chaincode invoke -o localhost:7050 --tls true --cafile $ORDERER_CA -C mychannel -n bigFileTransfer --waitForEvent --peerAddresses localhost:7051 --tlsRootCertFiles ${PEER0_ORG1_CA} --peerAddresses localhost:9051 --tlsRootCertFiles ${PEER0_ORG2_CA} -c '{"Args":["RequestSftp","wac_nearside.tif","Org1MSP","Org2MSP"]}'

peer chaincode query -C mychannel -n bigFileTransfer -c '{"Args":["Get","wac_nearside.tifOrg1MSPOrg2MSP"]}'


peer chaincode invoke -o localhost:7050 --tls true --cafile $ORDERER_CA -C mychannel -n bigFileTransfer --waitForEvent --waitForEventTimeout 3000s --peerAddresses localhost:7051 --tlsRootCertFiles ${PEER0_ORG1_CA} -c '{"Args":["SendFile_SSH_sftp","wac_nearside.tifOrg1MSPOrg2MSP","wac_nearside.tif","org2"]}'


peer chaincode query --peerAddresses localhost:9051 --tlsRootCertFiles ${PEER0_ORG2_CA} -C mychannel -n bigFileTransfer -c '{"Args":["CheckReceivedFile","wac_nearside.tif"]}'


peer chaincode invoke -o localhost:7050 --tls true --cafile $ORDERER_CA -C mychannel -n bigFileTransfer --peerAddresses localhost:7051 --tlsRootCertFiles ${PEER0_ORG1_CA} -c '{"Args":["RequestKey","wac_nearside.tifOrg1MSPOrg2MSP"]}'


peer chaincode invoke -o localhost:7050 --tls true --cafile $ORDERER_CA -C mychannel -n bigFileTransfer --peerAddresses localhost:9051 --tlsRootCertFiles ${PEER0_ORG2_CA} -C mychannel -n bigFileTransfer -c '{"Args":["QueryKeyAndDecryptFile","wac_nearside.tif","wac_nearside.tifOrg1MSPOrg2MSP"]}'


peer chaincode query --peerAddresses localhost:9051 --tlsRootCertFiles ${PEER0_ORG2_CA} -C mychannel -n bigFileTransfer -c '{"Args":["ListFolder","/home/chaincode/off_state"]}'


export org2CCID=dafcd37db64b
docker exec -it -u root $org2CCID /bin/sh

export org2CCID=47f4ad13b448
docker cp $org2CCID:/home/chaincode/off_state/ecn_wac_nearside.tif /home/cgao/Hyperledger/fabric-samples/test-network/




%*****************************
% Application 67.pdf
%*****************************


peer chaincode invoke -o localhost:7050 --tls true --cafile $ORDERER_CA -C mychannel -n bigFileTransfer --waitForEvent --peerAddresses localhost:7051 --tlsRootCertFiles ${PEER0_ORG1_CA} --peerAddresses localhost:9051 --tlsRootCertFiles ${PEER0_ORG2_CA} -c '{"Args":["Set","67.pdf","initFile"]}' --waitForEvent

peer chaincode invoke -o localhost:7050 --tls true --cafile $ORDERER_CA -C mychannel -n bigFileTransfer --waitForEvent --peerAddresses localhost:7051 --tlsRootCertFiles ${PEER0_ORG1_CA} --peerAddresses localhost:9051 --tlsRootCertFiles ${PEER0_ORG2_CA} -c '{"Args":["SetKeyPolicy","67.pdf","Org1MSP"]}' --waitForEvent

peer chaincode invoke -o localhost:7050 --tls true --cafile $ORDERER_CA -C mychannel -n bigFileTransfer --peerAddresses localhost:7051 --tlsRootCertFiles ${PEER0_ORG1_CA} -c '{"Args":["GetKeyPolicy","67.pdf"]}'


peer chaincode query -C mychannel -n bigFileTransfer -c '{"Args":["GetIpInfo","ip.peer0.org1"]}'
sftp root@192.168.208.7:/home/chaincode/off_state <<< $'put /home/cgao/Hyperledger/fabric-samples/test-network/67.pdf'




peer chaincode query --peerAddresses localhost:7051 --tlsRootCertFiles ${PEER0_ORG1_CA} -C mychannel -n bigFileTransfer -c '{"Args":["ListFolder","/home/chaincode/off_state"]}'


peer chaincode invoke -o localhost:7050 --tls true --cafile $ORDERER_CA -C mychannel -n bigFileTransfer --waitForEvent --peerAddresses localhost:7051 --tlsRootCertFiles ${PEER0_ORG1_CA} -c '{"Args":["AddFileInfo","{\"name\":\"67.pdf\",\"filehash\":\"\",\"objecttype\":\"file\",\"description\":\"67.pdf belong to org1\",\"owner\":\"Org1MSP\"}"]}'

peer chaincode invoke -o localhost:7050 --tls true --cafile $ORDERER_CA -C mychannel -n bigFileTransfer --peerAddresses localhost:7051 --tlsRootCertFiles ${PEER0_ORG1_CA} -c '{"Args":["QueryFileInfo","67.pdf"]}'




peer chaincode invoke -o localhost:7050 --tls true --cafile $ORDERER_CA -C mychannel -n bigFileTransfer --peerAddresses localhost:7051 --tlsRootCertFiles ${PEER0_ORG1_CA} -c '{"Args":["ListAllFiles"]}'




peer chaincode invoke -o localhost:7050 --tls true --cafile $ORDERER_CA -C mychannel -n bigFileTransfer --waitForEvent --peerAddresses localhost:7051 --tlsRootCertFiles ${PEER0_ORG1_CA} --peerAddresses localhost:9051 --tlsRootCertFiles ${PEER0_ORG2_CA} -c '{"Args":["RequestSftp","67.pdf","Org1MSP","Org2MSP"]}'

peer chaincode query -C mychannel -n bigFileTransfer -c '{"Args":["Get","67.pdfOrg1MSPOrg2MSP"]}'


peer chaincode invoke -o localhost:7050 --tls true --cafile $ORDERER_CA -C mychannel -n bigFileTransfer --waitForEvent --waitForEventTimeout 3000s --peerAddresses localhost:7051 --tlsRootCertFiles ${PEER0_ORG1_CA} -c '{"Args":["SendFile_SSH_sftp","67.pdfOrg1MSPOrg2MSP","67.pdf","org2"]}'


peer chaincode query --peerAddresses localhost:9051 --tlsRootCertFiles ${PEER0_ORG2_CA} -C mychannel -n bigFileTransfer -c '{"Args":["CheckReceivedFile","67.pdf"]}'


peer chaincode invoke -o localhost:7050 --tls true --cafile $ORDERER_CA -C mychannel -n bigFileTransfer --peerAddresses localhost:7051 --tlsRootCertFiles ${PEER0_ORG1_CA} -c '{"Args":["RequestKey","67.pdfOrg1MSPOrg2MSP"]}'


peer chaincode invoke -o localhost:7050 --tls true --cafile $ORDERER_CA -C mychannel -n bigFileTransfer --peerAddresses localhost:9051 --tlsRootCertFiles ${PEER0_ORG2_CA} -C mychannel -n bigFileTransfer -c '{"Args":["QueryKeyAndDecryptFile","67.pdf","67.pdfOrg1MSPOrg2MSP"]}'


peer chaincode query --peerAddresses localhost:9051 --tlsRootCertFiles ${PEER0_ORG2_CA} -C mychannel -n bigFileTransfer -c '{"Args":["ListFolder","/home/chaincode/off_state"]}'


export org2CCID=dafcd37db64b
docker exec -it -u root $org2CCID /bin/sh

export org2CCID=dafcd37db64b
docker cp $org2CCID:/home/chaincode/off_state/ecn_67.pdf /home/cgao/Hyperledger/fabric-samples/test-network/




%********************
% fetch blocks
%********************

peer channel fetch 19 block_19.pb -o localhost:7050 -c mychannel --tls --cafile $ORDERER_CA

configtxlator proto_decode --input block_19.pb --type common.Block --output block_19.json

chmod 777 block_19.json

sudo docker cp cli:/opt/gopath/src/github.com/hyperledger/fabric/peer/block_12.json /home/cgao/Hyperledger/fabric-samples/3org-privatedata-200/blocks/

chmod 777 -R blocks












peer chaincode query --peerAddresses localhost:7051 --tlsRootCertFiles ${PEER0_ORG1_CA} -C mychannel -n bigFileTransfer -c '{"Args":["ExcShell","/home/chaincode/off_state/testShell.sh"]}'



sftp root@192.168.128.8:/home/chaincode/off_state <<< $'put /home/cgao/Hyperledger/fabric-samples/test-network/bigFileTransfer.tar.gz'



export org1CCID=b6130b2c2b51
docker cp /home/cgao/Hyperledger/fabric-samples/test-network/testShell.sh $org1CCID:/home/chaincode/off_state/

export org1CCID=b6130b2c2b51
docker cp /home/cgao/Hyperledger/fabric-samples/test-network/bigFile.txt $org1CCID:/home/chaincode/off_state/



export org1CCID=a2bd7458e004
docker exec -it -u root $org1CCID /bin/sh

export org1CCID=43a5d2279411
docker cp /home/cgao/Hyperledger/fabric-samples/test-network/bigFile.txt $org1CCID:/opt/off_state/




export org2CCID=ac4d67a763ba
docker cp /home/cgao/Hyperledger/fabric-samples/test-network/testShell.sh $org2CCID:/opt/

export org2CCID=6bebb0ef10b6 
docker exec -it -u root $org2CCID /bin/sh



export org2CCID=157ddc1b4041 
docker exec -it $org2CCID /bin/sh



adduser  chaincode -u 20001 -D -S -s /bin/bash -G groupA



cd opt
mkdir off_state
cd off_state

apk add openssh

    sed -i "s/#PermitRootLogin.*/PermitRootLogin yes/g" /etc/ssh/sshd_config && \

    ssh-keygen -t dsa -P "" -f /etc/ssh/ssh_host_dsa_key && \

    ssh-keygen -t rsa -P "" -f /etc/ssh/ssh_host_rsa_key && \

    ssh-keygen -t ecdsa -P "" -f /etc/ssh/ssh_host_ecdsa_key && \

    ssh-keygen -t ed25519 -P "" -f /etc/ssh/ssh_host_ed25519_key && \

    echo "root:admin" | chpasswd

/usr/sbin/sshd -D &




ssh-keygen

cd /root/.ssh

ls

cat id_rsa.pub

touch authorized_keys

vi authorized_keys



scp -v /opt/off_state/bigFile.txt root@172.24.0.7:/opt/off_state



export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="Org1MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_ADDRESS=localhost:7051

peer chaincode query --peerAddresses localhost:7051 --tlsRootCertFiles ${PEER0_ORG1_CA} -C mychannel -n bigFileTransfer -c '{"Args":["SendFiles","bigFile.txt"]}'

peer chaincode query --peerAddresses localhost:7051 --tlsRootCertFiles ${PEER0_ORG1_CA} -C mychannel -n bigFileTransfer -c '{"Args":["ListFolder","/opt/"]}'

peer chaincode query --peerAddresses localhost:9051 --tlsRootCertFiles ${PEER0_ORG2_CA} -C mychannel -n bigFileTransfer -c '{"Args":["ListFolder","/opt/"]}'

peer chaincode query --peerAddresses localhost:7051 --tlsRootCertFiles ${PEER0_ORG1_CA} -C mychannel -n bigFileTransfer -c '{"Args":["ExcShell","/opt/off_state/scp.sh"]}'

peer chaincode query --peerAddresses localhost:7051 --tlsRootCertFiles ${PEER0_ORG1_CA} -C mychannel -n bigFileTransfer -c '{"Args":["ExcShell","/opt/testShell.sh"]}'

peer chaincode query --peerAddresses localhost:7051 --tlsRootCertFiles ${PEER0_ORG1_CA} -C mychannel -n bigFileTransfer -c '{"Args":["Test_SSH_sftp"]}'



export org1CCID=b3bdcf524b09
docker exec -it -u root $org1CCID /bin/sh
cd /home/chaincode/
./SSHShell.sh
cd off_state


export org1CCID=f6fad548780b
docker exec -it $org1CCID /bin/sh


sudo passwd chaincode


peer chaincode query -C mychannel -n bigFileTransfer --peerAddresses localhost:9051 --tlsRootCertFiles ${PEER0_ORG2_CA} -c '{"Args":["ExcShell","/opt/shell.sh"]}'

peer chaincode query -C mychannel -n bigFileTransfer --peerAddresses localhost:7051 --tlsRootCertFiles ${PEER0_ORG1_CA} -c '{"Args":["Init"]}'










%*****************************
% Application 100.mp4
%*****************************


peer chaincode invoke -o localhost:7050 --tls true --cafile $ORDERER_CA -C mychannel -n bigFileTransfer --waitForEvent --peerAddresses localhost:7051 --tlsRootCertFiles ${PEER0_ORG1_CA} --peerAddresses localhost:9051 --tlsRootCertFiles ${PEER0_ORG2_CA} -c '{"Args":["Set","100.mp4","initFile"]}' --waitForEvent

peer chaincode invoke -o localhost:7050 --tls true --cafile $ORDERER_CA -C mychannel -n bigFileTransfer --waitForEvent --peerAddresses localhost:7051 --tlsRootCertFiles ${PEER0_ORG1_CA} --peerAddresses localhost:9051 --tlsRootCertFiles ${PEER0_ORG2_CA} -c '{"Args":["SetKeyPolicy","100.mp4","Org1MSP"]}' --waitForEvent

peer chaincode invoke -o localhost:7050 --tls true --cafile $ORDERER_CA -C mychannel -n bigFileTransfer --peerAddresses localhost:7051 --tlsRootCertFiles ${PEER0_ORG1_CA} -c '{"Args":["GetKeyPolicy","100.mp4"]}'


peer chaincode query -C mychannel -n bigFileTransfer -c '{"Args":["GetIpInfo","ip.peer0.org1"]}'
sftp root@172.22.0.8:/home/chaincode/off_state <<< $'put /home/cgao/Hyperledger/fabric-samples/test-network/100.mp4'




peer chaincode query --peerAddresses localhost:7051 --tlsRootCertFiles ${PEER0_ORG1_CA} -C mychannel -n bigFileTransfer -c '{"Args":["ListFolder","/home/chaincode/off_state"]}'


peer chaincode invoke -o localhost:7050 --tls true --cafile $ORDERER_CA -C mychannel -n bigFileTransfer --waitForEvent --peerAddresses localhost:7051 --tlsRootCertFiles ${PEER0_ORG1_CA} -c '{"Args":["AddFileInfo","{\"name\":\"100.mp4\",\"filehash\":\"\",\"objecttype\":\"file\",\"description\":\"100.mp4 belong to org1\",\"owner\":\"Org1MSP\"}"]}'

peer chaincode invoke -o localhost:7050 --tls true --cafile $ORDERER_CA -C mychannel -n bigFileTransfer --peerAddresses localhost:7051 --tlsRootCertFiles ${PEER0_ORG1_CA} -c '{"Args":["QueryFileInfo","100.mp4"]}'




peer chaincode invoke -o localhost:7050 --tls true --cafile $ORDERER_CA -C mychannel -n bigFileTransfer --peerAddresses localhost:7051 --tlsRootCertFiles ${PEER0_ORG1_CA} -c '{"Args":["ListAllFiles"]}'




peer chaincode invoke -o localhost:7050 --tls true --cafile $ORDERER_CA -C mychannel -n bigFileTransfer --waitForEvent --peerAddresses localhost:7051 --tlsRootCertFiles ${PEER0_ORG1_CA} --peerAddresses localhost:9051 --tlsRootCertFiles ${PEER0_ORG2_CA} -c '{"Args":["RequestSftp","100.mp4","Org1MSP","Org2MSP"]}'

peer chaincode query -C mychannel -n bigFileTransfer -c '{"Args":["Get","100.mp4Org1MSPOrg2MSP"]}'


peer chaincode invoke -o localhost:7050 --tls true --cafile $ORDERER_CA -C mychannel -n bigFileTransfer --waitForEvent --waitForEventTimeout 3000s --peerAddresses localhost:7051 --tlsRootCertFiles ${PEER0_ORG1_CA} -c '{"Args":["SendFile_SSH_sftp","100.mp4Org1MSPOrg2MSP","100.mp4","org2"]}'


peer chaincode query --peerAddresses localhost:9051 --tlsRootCertFiles ${PEER0_ORG2_CA} -C mychannel -n bigFileTransfer -c '{"Args":["CheckReceivedFile","100.mp4"]}'


peer chaincode invoke -o localhost:7050 --tls true --cafile $ORDERER_CA -C mychannel -n bigFileTransfer --peerAddresses localhost:7051 --tlsRootCertFiles ${PEER0_ORG1_CA} -c '{"Args":["RequestKey","100.mp4Org1MSPOrg2MSP"]}'


peer chaincode invoke -o localhost:7050 --tls true --cafile $ORDERER_CA -C mychannel -n bigFileTransfer --peerAddresses localhost:9051 --tlsRootCertFiles ${PEER0_ORG2_CA} -C mychannel -n bigFileTransfer -c '{"Args":["QueryKeyAndDecryptFile","100.mp4","100.mp4Org1MSPOrg2MSP"]}'


peer chaincode query --peerAddresses localhost:9051 --tlsRootCertFiles ${PEER0_ORG2_CA} -C mychannel -n bigFileTransfer -c '{"Args":["ListFolder","/home/chaincode/off_state"]}'


export org2CCID=1444f4cfc3ff
docker exec -it -u root $org2CCID /bin/sh

export org2CCID=dafcd37db64b
docker cp $org2CCID:/home/chaincode/off_state/ecn_100.mp4 /home/cgao/Hyperledger/fabric-samples/test-network/




locate 1_2g.zip_ciphertext.bin

/var/snap/docker/common/var-lib-docker/overlay2/78e842828addcfeee8a7b19fc343d50c93593e132671d9ca565d9cd2e8dca603/diff/home/chaincode/off_state/1_2g.zip_ciphertext.bin


./evlTransferLatency.sh 2>&1 | tee evlTransferLatency_1_2g.log


peer chaincode query -C mychannel -n bigFileTransfer --peerAddresses localhost:7051 --tlsRootCertFiles ${PEER0_ORG1_CA} -c '{"Args":["GetPvtData","key_1_2gv8.zipOrg1MSPOrg2MSP","_implicit_org_Org1MSP"]}'



peer chaincode invoke -o localhost:7050 --tls true --cafile $ORDERER_CA -C mychannel -n bigFileTransfer --waitForEvent --peerAddresses localhost:7051 --tlsRootCertFiles ${PEER0_ORG1_CA} --peerAddresses localhost:9051 --tlsRootCertFiles ${PEER0_ORG2_CA} -c '{"Args":["SetPvtData","pvtkey1","val1","collectionKey_org1"]}'


peer chaincode invoke -o localhost:7050 --tls true --cafile $ORDERER_CA -C mychannel -n bigFileTransfer --waitForEvent --peerAddresses localhost:7051 --tlsRootCertFiles ${PEER0_ORG1_CA} --peerAddresses localhost:9051 --tlsRootCertFiles ${PEER0_ORG2_CA} -c '{"Args":["GetPvtData","pvtkey1","collectionKey_org1"]}'


peer chaincode invoke -o localhost:7050 --tls true --cafile $ORDERER_CA -C mychannel -n bigFileTransfer --waitForEvent --peerAddresses localhost:7051 --tlsRootCertFiles ${PEER0_ORG1_CA} -c '{"Args":["GetPvtData","pvtkey1","collectionKey_org1"]}'


peer chaincode invoke -o localhost:7050 --tls true --cafile $ORDERER_CA -C mychannel -n bigFileTransfer --waitForEvent --peerAddresses localhost:7051 --tlsRootCertFiles ${PEER0_ORG1_CA} --peerAddresses localhost:9051 --tlsRootCertFiles ${PEER0_ORG2_CA} -c '{"Args":["SetPDCKeyPolicy","pvtkey1","collectionKey_org1","Org1MSP"]}'