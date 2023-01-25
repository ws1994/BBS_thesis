cd chaincode/bigFileTransfer/go

GO111MODULE=on go mod vendor

docker exec cli peer lifecycle chaincode package fileTransfer.tar.gz --path github.com/hyperledger/fabric/peer/chaincode/bigFileTransfer/go/ --label fileTransfer_1

docker exec cli peer lifecycle chaincode install fileTransfer.tar.gz

docker exec -e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp -e CORE_PEER_ADDRESS=peer0.org2.example.com:9051 -e CORE_PEER_LOCALMSPID="Org2MSP" -e CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt cli peer lifecycle chaincode install fileTransfer.tar.gz

docker exec -e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org3.example.com/users/Admin@org3.example.com/msp -e CORE_PEER_ADDRESS=peer0.org3.example.com:11051 -e CORE_PEER_LOCALMSPID="Org3MSP" -e CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org3.example.com/peers/peer0.org3.example.com/tls/ca.crt cli peer lifecycle chaincode install fileTransfer.tar.gz

docker exec cli peer lifecycle chaincode approveformyorg --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem --channelID mychannel --name fileTransfer --version 1 --sequence 1 --waitForEvent --package-id fileTransfer_1:515c2b6a08725fe1cc6b1bf91e5412115001920c7f30421a6e60f71dbffe1ae4

docker exec -e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp -e CORE_PEER_ADDRESS=peer0.org2.example.com:9051 -e CORE_PEER_LOCALMSPID="Org2MSP" -e CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt cli peer lifecycle chaincode approveformyorg --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem --channelID mychannel --name fileTransfer --version 1 --sequence 1 --waitForEvent --package-id fileTransfer_1:515c2b6a08725fe1cc6b1bf91e5412115001920c7f30421a6e60f71dbffe1ae4

docker exec -e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org3.example.com/users/Admin@org3.example.com/msp -e CORE_PEER_ADDRESS=peer0.org3.example.com:11051 -e CORE_PEER_LOCALMSPID="Org3MSP" -e CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org3.example.com/peers/peer0.org3.example.com/tls/ca.crt cli peer lifecycle chaincode approveformyorg --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem --channelID mychannel --name fileTransfer --version 1 --sequence 1 --waitForEvent --package-id fileTransfer_1:515c2b6a08725fe1cc6b1bf91e5412115001920c7f30421a6e60f71dbffe1ae4

docker exec cli peer lifecycle chaincode commit -o orderer.example.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem --peerAddresses peer0.org1.example.com:7051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt --peerAddresses peer0.org2.example.com:9051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt --peerAddresses peer0.org3.example.com:11051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org3.example.com/peers/peer0.org3.example.com/tls/ca.crt --channelID mychannel --name fileTransfer --version 1 --sequence 1 --waitForEvent




%********************
% update
%********************

docker exec cli peer lifecycle chaincode package fileTransfer.tar.gz --path github.com/hyperledger/fabric/peer/chaincode/bigFileTransfer/go/ --label fileTransfer_1

docker exec cli peer lifecycle chaincode install fileTransfer.tar.gz

docker exec -e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp -e CORE_PEER_ADDRESS=peer0.org2.example.com:9051 -e CORE_PEER_LOCALMSPID="Org2MSP" -e CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt cli peer lifecycle chaincode install fileTransfer.tar.gz

docker exec -e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org3.example.com/users/Admin@org3.example.com/msp -e CORE_PEER_ADDRESS=peer0.org3.example.com:11051 -e CORE_PEER_LOCALMSPID="Org3MSP" -e CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org3.example.com/peers/peer0.org3.example.com/tls/ca.crt cli peer lifecycle chaincode install fileTransfer.tar.gz

docker exec cli peer lifecycle chaincode approveformyorg --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem --channelID mychannel --name fileTransfer --version 2 --sequence 2 --waitForEvent --package-id fileTransfer_1:515c2b6a08725fe1cc6b1bf91e5412115001920c7f30421a6e60f71dbffe1ae4

docker exec -e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp -e CORE_PEER_ADDRESS=peer0.org2.example.com:9051 -e CORE_PEER_LOCALMSPID="Org2MSP" -e CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt cli peer lifecycle chaincode approveformyorg --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem --channelID mychannel --name fileTransfer --version 2 --sequence 2 --waitForEvent --package-id fileTransfer_1:515c2b6a08725fe1cc6b1bf91e5412115001920c7f30421a6e60f71dbffe1ae4

docker exec -e CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org3.example.com/users/Admin@org3.example.com/msp -e CORE_PEER_ADDRESS=peer0.org3.example.com:11051 -e CORE_PEER_LOCALMSPID="Org3MSP" -e CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org3.example.com/peers/peer0.org3.example.com/tls/ca.crt cli peer lifecycle chaincode approveformyorg --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem --channelID mychannel --name fileTransfer --version 2 --sequence 2 --waitForEvent --package-id fileTransfer_1:515c2b6a08725fe1cc6b1bf91e5412115001920c7f30421a6e60f71dbffe1ae4

docker exec cli peer lifecycle chaincode commit -o orderer.example.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem --peerAddresses peer0.org1.example.com:7051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt --peerAddresses peer0.org2.example.com:9051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt --peerAddresses peer0.org3.example.com:11051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org3.example.com/peers/peer0.org3.example.com/tls/ca.crt --channelID mychannel --name fileTransfer --version 2 --sequence 2 --waitForEvent




%********************
% test 
%********************

docker exec cli peer chaincode invoke -o orderer.example.com:7050 --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C mychannel -n fileTransfer --peerAddresses peer0.org1.example.com:7051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt --peerAddresses peer0.org2.example.com:9051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt -c '{"Args":["callScp"]}'

docker exec cli peer chaincode query -C mychannel -n fileTransfer -c '{"Args":["callScp"]}'


docker exec cli peer chaincode query -C mychannel -n fileTransfer -c '{"Args":["Init"]}'

docker exec cli peer chaincode query -C mychannel -n fileTransfer -c '{"Args":["QueryFileInfo","filehash"]}'




docker exec -it peer0.org1.example.com /bin/ls /

docker exec -it peer0.org1.example.com /bin/sh


docker cp /home/cgao/Hyperledger/fabric-samples/3org-privatedata-200/blocks/block_6.json dev-peer0.org1.example.com-fileTransfer_1-515c2b6a08725fe1cc6b1bf91e5412115001920c7f30421a6e60f71dbffe1ae4:/opt/blocks/ 


docker exec -it dev-peer0.org1.example.com-fileTransfer_1-515c2b6a08725fe1cc6b1bf91e5412115001920c7f30421a6e60f71dbffe1ae4 /bin/sh

docker exec -it -u root dev-peer0.org1.example.com-fileTransfer_1-515c2b6a08725fe1cc6b1bf91e5412115001920c7f30421a6e60f71dbffe1ae4 /bin/sh

docker exec -it -u root dev-peer0.org2.example.com-fileTransfer_1-515c2b6a08725fe1cc6b1bf91e5412115001920c7f30421a6e60f71dbffe1ae4 /bin/sh





apk add openssh

apk del ***


sed -i "s/#PermitRootLogin.*/PermitRootLogin yes/g" /etc/ssh/sshd_config && \
ssh-keygen -t dsa -P "" -f /etc/ssh/ssh_host_dsa_key && \
ssh-keygen -t rsa -P "" -f /etc/ssh/ssh_host_rsa_key && \
ssh-keygen -t ecdsa -P "" -f /etc/ssh/ssh_host_ecdsa_key && \
ssh-keygen -t ed25519 -P "" -f /etc/ssh/ssh_host_ed25519_key && \
echo "root:admin" | chpasswd


/usr/sbin/sshd -D &





ssh-keygen -t dsa -P "" -f /etc/ssh/ssh_host_dsa_key && \
ssh-keygen -t rsa -P "" -f /etc/ssh/ssh_host_rsa_key && \
ssh-keygen -t ecdsa -P "" -f /etc/ssh/ssh_host_ecdsa_key && \
ssh-keygen -t ed25519 -P "" -f /etc/ssh/ssh_host_ed25519_key && \
echo "root:shan" | chpasswd