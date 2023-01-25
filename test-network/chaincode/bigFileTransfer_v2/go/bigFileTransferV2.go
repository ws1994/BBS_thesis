package main

import (
	"bytes"
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/hyperledger/fabric-chaincode-go/pkg/statebased"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	// "github.com/hyperledger/fabric/core/chaincode/shim/ext/statebased"
)

type File struct {
	Name        string   `json:"name"`
	FileHash    string   `json:"filehash"`
	ObjectType  string   `json:"objecttype"`
	Description string   `json:"description"`
	Owner       string   `json:"owner"`
	Rule        []string `json:"rule"`
	PeerMSP     string   `json:"peermsp"`
}

type Event struct {
	// TxID           string `json:"txid"`
	EventID         string `json:"eventid"`
	Flag            bool   `json:"flag"`
	FileName        string `json:"fileName"`
	SourceOrg       string `json:"sourceorg"`
	SourceUser      string `json:"sourceuser"`
	DestinationOrg  string `json:"destinationorg"`
	DestinationUser string `json:"destinationuser"`
}

type RSAKey struct {
	PrivateKey string `json:"privatekey"`
	PublicKey  string `json:"publickey"`
}

type IP struct {
	Name   string `json:"name"`
	IpAdr  string `json:"ipadr"`
	IpPort int    `json:"ipport"`
}

type EncryptionKey struct {
	KeyID             string `json:"keyid"`
	EncryptedFileHash string `json:"encryptedfilehash"`
	AESKey            []byte `json:"aeskey"`
	Signature         []byte `json:"signature"`
	RSAPublicKey      string `json:"raspublickey"`
}

type SmartContract struct {
	contractapi.Contract
}

// ============================================================
// ID manage
// ============================================================
func (s *SmartContract) InitIP(ctx contractapi.TransactionContextInterface) error {
	s.Set(ctx, "ip.Org1MSP", "initIP")
	s.Set(ctx, "ip.Org2MSP", "initIP")
	s.Set(ctx, "ip.Org3MSP", "initIP")

	return nil
}

func (s *SmartContract) CreateOffState(ctx contractapi.TransactionContextInterface) error {
	submitterByte, _ := ctx.GetStub().GetCreator()
	clientName, err := getSubmitterName(submitterByte)
	if err != nil {
		return fmt.Errorf("fail to get client id")
	}
	log.Print("client name: " + clientName)
	path := "./home/chaincode/"
	path += clientName
	os.Mkdir(path, 0777)
	return nil
}

func (s *SmartContract) InitConfig(ctx contractapi.TransactionContextInterface) (string, error) {

	// os.Mkdir("./home/chaincode/off_state", 0777)

	// cmd1 := "cd opt" + "\n"
	// cmd2 := "mkdir off_state" + "\n"
	cmd2 := "apk add expect" + "\n"
	cmd3 := "apk add openssh" + "\n"
	cmd4 := "sed -i \"s/#PermitRootLogin.*/PermitRootLogin yes/g\" /etc/ssh/sshd_config && \\" + "\n"
	cmd5 := "ssh-keygen -t dsa -P \"\" -f /etc/ssh/ssh_host_dsa_key && \\" + "\n"
	cmd6 := "ssh-keygen -t rsa -P \"\" -f /etc/ssh/ssh_host_rsa_key && \\" + "\n"
	cmd7 := "ssh-keygen -t ecdsa -P \"\" -f /etc/ssh/ssh_host_ecdsa_key && \\" + "\n"
	cmd8 := "ssh-keygen -t ed25519 -P \"\" -f /etc/ssh/ssh_host_ed25519_key && \\" + "\n"
	cmd9 := "echo \"root:admin\" | chpasswd" + "\n"
	cmd10 := "/usr/sbin/sshd -D &"

	command := cmd2 + cmd3 + cmd4 + cmd5 + cmd6 + cmd7 + cmd8 + cmd9 + cmd10

	f, err := os.OpenFile("/home/chaincode/SSHShell.sh", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0777)
	defer f.Close()
	if err != nil {
		return "open: ", fmt.Errorf(err.Error())
	} else {
		_, err4 := f.Write([]byte(command))
		if err4 != nil {
			return "write: ", fmt.Errorf(err4.Error())
		}
	}

	b, err1 := ioutil.ReadFile("/home/chaincode/SSHShell.sh") // just pass the file name
	if err1 != nil {
		return "read: ", fmt.Errorf(err1.Error())
	}

	return string(b), nil
}

func createPkcs8Keys(keyLength int) (privateKey, publicKey string) {
	rsaPrivateKey, err := rsa.GenerateKey(rand.Reader, keyLength)
	if err != nil {
		return
	}

	objPkcs8, _ := x509.MarshalPKCS8PrivateKey(rsaPrivateKey)

	privateKey = string(pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: objPkcs8,
	}))

	objPkix, err := x509.MarshalPKIXPublicKey(&rsaPrivateKey.PublicKey)
	if err != nil {
		return
	}

	publicKey = string(pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: objPkix,
	}))
	return
}

func (s *SmartContract) InitRSAKeyPair(ctx contractapi.TransactionContextInterface, org string) error {
	var rasKey RSAKey

	privateKey, publicKey := createPkcs8Keys(2048)

	if privateKey == "" || privateKey == "" {
		return fmt.Errorf("failed to init RSA key pair")
	}

	rasKey.PrivateKey = privateKey
	rasKey.PublicKey = publicKey

	RSAKeyInfoJSONasBytes, err := json.Marshal(rasKey)
	if err != nil {
		return fmt.Errorf(err.Error())
	}

	err1 := ctx.GetStub().PutPrivateData("_implicit_org_"+org, org+"_RSA", RSAKeyInfoJSONasBytes)
	if err1 != nil {
		return fmt.Errorf("failed to set rsa key: %s", err1.Error())
	}

	return nil
}

func (s *SmartContract) TestBufferSizeEncDec(ctx contractapi.TransactionContextInterface, sizeStr string, filename string) (string, error) {

	size, _ := strconv.Atoi(sizeStr)

	var key []byte
	key32 := make([]byte, 32)
	_, errkey := rand.Read(key)
	if errkey != nil {
		return "", fmt.Errorf("%s", errkey.Error())
	}
	// println(key32)
	key = key32

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("%s", err.Error())
	}

	startEncrypt := time.Now()
	infile, err := os.Open("./home/chaincode/" + filename)
	if err != nil {
		return "", fmt.Errorf("%s", err.Error())
	}
	defer infile.Close()

	// Never use more than 2^32 random nonces with a given key
	// because of the risk of repeat.
	iv := make([]byte, block.BlockSize())
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", fmt.Errorf("%s", err.Error())
	}

	outfile, err := os.OpenFile("./home/chaincode/ciphertext.bin", os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		return "", fmt.Errorf("%s", err.Error())
	}
	defer outfile.Close()

	// The buffer size must be multiple of 16 bytes

	buf := make([]byte, 1024*size)
	stream := cipher.NewCTR(block, iv)
	for {
		n, err := infile.Read(buf)
		if n > 0 {
			stream.XORKeyStream(buf, buf[:n])
			// Write into file
			outfile.Write(buf[:n])
		}

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Printf("Read %d bytes: %v", n, err)
			break
		}
	}
	// Append the IV
	outfile.Write(iv)

	endEncrypt := time.Since(startEncrypt)
	timeEncrypt := strconv.FormatFloat(float64(endEncrypt.Seconds()), 'f', 6, 64)
	fmt.Println("encrypt time：", timeEncrypt)

	startDecrypt := time.Now()
	infile, err = os.Open("./home/chaincode/ciphertext.bin")
	if err != nil {
		return "", fmt.Errorf("%s", err.Error())
	}
	defer infile.Close()

	block, err = aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("%s", err.Error())
	}

	fi, err := infile.Stat()
	if err != nil {
		return "", fmt.Errorf("%s", err.Error())
	}

	iv = make([]byte, block.BlockSize())
	msgLen := fi.Size() - int64(len(iv))
	_, err = infile.ReadAt(iv, msgLen)
	if err != nil {
		return "", fmt.Errorf("%s", err.Error())
	}

	outfile, err = os.OpenFile("./home/chaincode/ciphertext.bin", os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		return "", fmt.Errorf("%s", err.Error())
	}
	defer outfile.Close()

	// The buffer size must be multiple of 16 bytes

	buf = make([]byte, 1024*size)
	stream = cipher.NewCTR(block, iv)
	for {
		n, err := infile.Read(buf)
		if n > 0 {
			// The last bytes are the IV, don't belong the original message
			if n > int(msgLen) {
				n = int(msgLen)
			}
			msgLen -= int64(n)
			stream.XORKeyStream(buf, buf[:n])
			// Write into file
			outfile.Write(buf[:n])
		}

		if err == io.EOF {
			break
		}

		if err != nil {
			return "", fmt.Errorf("%s", err.Error())
			// break
		}
	}
	endDecrypt := time.Since(startDecrypt)
	timeDecrypt := strconv.FormatFloat(float64(endDecrypt.Seconds()), 'f', 6, 64)
	// fmt.Println("decrypt time：", timeDecrypt)
	return "encTime: " + timeEncrypt + " decTime: " + timeDecrypt, nil
}

func (s *SmartContract) InitRSAKeyForCC(ctx contractapi.TransactionContextInterface) error {
	peerMSP, _ := shim.GetMSPID()
	rsaPrivateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return fmt.Errorf("failed to generate RSA key pair")
	}

	objPkcs8, _ := x509.MarshalPKCS8PrivateKey(rsaPrivateKey)

	privateKey := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: objPkcs8,
	}

	file, err := os.Create("/home/chaincode/ccPrivate.pem")
	if err != nil {
		return fmt.Errorf("create pem file failed")
	}
	err = pem.Encode(file, privateKey)

	objPkix, err := x509.MarshalPKIXPublicKey(&rsaPrivateKey.PublicKey)
	if err != nil {
		return fmt.Errorf("failed to init RSA key pair")
	}

	publicKey := string(pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: objPkix,
	}))

	err1 := ctx.GetStub().PutState(peerMSP+"PublicKey", []byte(publicKey))
	if err1 != nil {
		return fmt.Errorf("failed to set rsa public key: %s", err1.Error())
	}

	return nil
}

func (s *SmartContract) SetRsaPKForClient(ctx contractapi.TransactionContextInterface, pk string) error {

	submitterByte, _ := ctx.GetStub().GetCreator()
	clientName, err := getSubmitterName(submitterByte)
	if err != nil {
		return fmt.Errorf("fail to get client id")
	}

	err1 := ctx.GetStub().PutState(clientName+"PublicKey", []byte(pk))
	log.Print(clientName + "PublicKey")
	log.Print(pk)
	if err1 != nil {
		return fmt.Errorf("failed to set rsa public key for client: %s", err1.Error())
	}

	return nil
}

// func (s *SmartContract) TestRsaPKForClientEnc(ctx contractapi.TransactionContextInterface) (string, error) {

// 	submitterByte, _ := ctx.GetStub().GetCreator()
// 	clientName, err := getSubmitterName(submitterByte)
// 	if err != nil {
// 		return "", fmt.Errorf("fail to get client id")
// 	}

// 	pkval, err1 := ctx.GetStub().GetState(clientName + "PublicKey")
// 	if err1 != nil || pkval == nil {
// 		return "", fmt.Errorf("failed to get rsa public key for client: %s", err1.Error())
// 	}
// 	// log.Print(pkval)
// 	log.Print("start encrypt hello rsa")

// 	// ciphertext, _ := useRSAKeyEnc(pkval, []byte("hello rsa"))

// 	block, _ := pem.Decode(pkval)
// 	if block == nil {
// 		return "", fmt.Errorf("failed to parse pk block")
// 	}
// 	pkInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to parse public key interface: %s", err.Error())
// 	}
// 	pk, isKey := pkInterface.(*rsa.PublicKey)
// 	if isKey != true {
// 		return "", fmt.Errorf("is not a public key")
// 	}

// 	ciphertext, _ := rsa.EncryptOAEP(sha1.New(), rand.Reader, pk, []byte("hello rsa"), nil)

// 	// err = ctx.GetStub().PutState("org1mspCipher", ciphertext)

// 	// return ciphertext, nil

// 	log.Print("end encrypt hello rsa")

// 	return string(ciphertext), nil
// }

func (s *SmartContract) Set(ctx contractapi.TransactionContextInterface, key string, value string) error {

	err2 := ctx.GetStub().PutState(key, []byte(value))
	if err2 != nil {
		return fmt.Errorf("failed to set key: %s", err2.Error())
	}

	return nil
}

func (s *SmartContract) Sleep(ctx contractapi.TransactionContextInterface, key string, value string) error {

	err2 := ctx.GetStub().PutState(key, []byte(value))
	if err2 != nil {
		return fmt.Errorf("failed to set key: %s", err2.Error())
	}
	log.Print("start sleep")
	time.Sleep(15 * time.Second)
	return nil
}

func (s *SmartContract) Get(ctx contractapi.TransactionContextInterface, key string) (string, error) {

	val, err2 := ctx.GetStub().GetState(key)
	if err2 != nil {
		return "", fmt.Errorf("failed to get key: %s", err2.Error())
	}
	if val == nil {
		return "", fmt.Errorf("%s does not exist.", key)
	}

	return string(val), nil
}

func (s *SmartContract) Delete(ctx contractapi.TransactionContextInterface, key string) (string, error) {

	err2 := ctx.GetStub().DelState(key)
	if err2 != nil {
		return "", fmt.Errorf("failed to get key: %s", err2.Error())
	}

	return key, nil
}

func (s *SmartContract) GetPvtData(ctx contractapi.TransactionContextInterface, key string, collection string) (string, error) {

	val, err2 := ctx.GetStub().GetPrivateData(collection, key)
	if err2 != nil {
		return "", fmt.Errorf("failed to get key: %s", err2.Error())
	}

	return string(val), nil
}

func (s *SmartContract) SetPvtData(ctx contractapi.TransactionContextInterface, key string, collection string, value string) (string, error) {

	err2 := ctx.GetStub().PutPrivateData(collection, key, []byte(value))
	if err2 != nil {
		return "", fmt.Errorf("failed to put key: %s", err2.Error())
	}
	return string(value), nil
}

func (s *SmartContract) SetKeyPolicy(ctx contractapi.TransactionContextInterface, key string, EP string) (string, error) {
	// key := "ip.peer0.org1"
	// EP := "Org1MSP"

	newEP, err := statebased.NewStateEP(nil)
	if err != nil {
		return "", fmt.Errorf("failed to init newEP: %s", err.Error())
	}
	err = newEP.AddOrgs(statebased.RoleTypeMember, EP)
	if err != nil {
		return "", fmt.Errorf("failed to invoke AddOrgs: %s", err.Error())
	}

	policyByte, err := newEP.Policy()
	if err != nil {
		return "", fmt.Errorf("failed to get policyByte: %s", err.Error())
	}

	epBytes, err := ctx.GetStub().GetStateValidationParameter(key)
	ep, err := statebased.NewStateEP(epBytes)
	orgs := ep.ListOrgs()
	orgsList, err := json.Marshal(orgs)

	err = ctx.GetStub().SetStateValidationParameter(key, policyByte)
	if err != nil {
		return "", fmt.Errorf("failed to invoke SetStateValidationParameter: %s", err.Error())
	}

	return "Current policy for key " + key + " is: " + string(orgsList), nil
}

func (s *SmartContract) GetKeyPolicy(ctx contractapi.TransactionContextInterface, key string) (string, error) {

	epBytes, err := ctx.GetStub().GetStateValidationParameter(key)
	ep, err := statebased.NewStateEP(epBytes)
	orgs := ep.ListOrgs()
	orgsList, err := json.Marshal(orgs)

	if err != nil {
		return "", fmt.Errorf("failed to get policy for key: %s", err.Error())
	}

	return "Policy for key " + key + " is: " + string(orgsList), nil
}

func (s *SmartContract) SetIP(ctx contractapi.TransactionContextInterface, key string) (string, error) {

	cmd := exec.Command("ifconfig") //| grep \"inet addr\"
	output, err := cmd.Output()
	if err != nil {
		return "cmd error", fmt.Errorf("failed to execute cmd: %s", err.Error())
	}

	outputStr := string(output)
	var theIndex int = strings.Index(outputStr, "inet addr:")
	var theIndexEnd int = strings.Index(outputStr, "Bcast")

	ipStr := outputStr[theIndex+10 : theIndexEnd-1]

	err2 := ctx.GetStub().PutState(key, []byte(ipStr))
	if err2 != nil {
		return "", fmt.Errorf("failed to init ip of peer0.org1: %s", err2.Error())
	}

	return ipStr, nil
}

func (s *SmartContract) GetIpInfo(ctx contractapi.TransactionContextInterface, key string) (string, error) {

	ipInfo, err := ctx.GetStub().GetState(key) //get the marble from chaincode state
	if err != nil {
		return "", fmt.Errorf("failed to read ip.peer0.org1. %s", err.Error())
	}
	if ipInfo == nil {
		return "", fmt.Errorf("%s does not exist", key)
	}

	return string(ipInfo), nil
}

func (s *SmartContract) ListFolder(ctx contractapi.TransactionContextInterface, path string) (string, error) {
	cmd := exec.Command("/bin/ls", path) //"/opt/off_state"
	output, _ := cmd.Output()
	// if err != nil {
	// 	return "cmd error", nil
	// }

	return string(output), nil
}

func (s *SmartContract) ListOffStateData(ctx contractapi.TransactionContextInterface) (string, error) {
	submitterByte, _ := ctx.GetStub().GetCreator()
	clientName, err := getSubmitterName(submitterByte)
	if err != nil {
		return "", fmt.Errorf("fail to get client id")
	}
	clientMSPID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return "", fmt.Errorf("fail to get client mspid")
	}
	peerMSPID, err := shim.GetMSPID()
	if err != nil {
		return "", fmt.Errorf("fail to get peer id")
	}
	if clientMSPID != peerMSPID {
		return "", fmt.Errorf("client and the peer belong to different org")
	}

	cmd := exec.Command("/bin/ls", "/home/chaincode/"+clientName)
	output, _ := cmd.Output()
	return string(output), nil
}

func (s *SmartContract) GetSubmitterPK(ctx contractapi.TransactionContextInterface) (string, error) {
	submitterByte, _ := ctx.GetStub().GetCreator()
	certStart := bytes.IndexAny(submitterByte, "-----BEGIN")
	if certStart == -1 {
		return "", fmt.Errorf("no certificate found")
	}
	certText := submitterByte[certStart:]
	bl, _ := pem.Decode(certText)
	if bl == nil {
		return "", fmt.Errorf("could not decode the PEM structure")
	}
	cert, err := x509.ParseCertificate(bl.Bytes)
	if err != nil {
		return "", fmt.Errorf("parse certificate failed")
	}
	rsaPublicKey, _ := x509.MarshalPKIXPublicKey(cert.PublicKey)
	rsaPublicKeyBlock := pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: rsaPublicKey,
	}
	publicKeyPem := string(pem.EncodeToMemory(&rsaPublicKeyBlock))
	return publicKeyPem, nil
}

func (s *SmartContract) ExcShell(ctx contractapi.TransactionContextInterface, shellPath string) (string, error) {

	// cmd := exec.Command("/bin/ls", path) //"/opt/off_state"
	cmd := exec.Command("/bin/sh", "-c", shellPath)

	err := cmd.Run()

	return "end", err
}

// ===============================================
// File
// ===============================================

func (s *SmartContract) QueryFileInfo(ctx contractapi.TransactionContextInterface, filename string) (*File, error) {
	// if !strings.HasPrefix(fileID, "f") {
	// 	return nil, fmt.Errorf("This is not a file ID: " + fileID)
	// }

	fileJSON, err := ctx.GetStub().GetState(filename) //get the marble from chaincode state
	if err != nil {
		return nil, fmt.Errorf("failed to read fileInfo %s", err.Error())
	}
	if fileJSON == nil {
		return nil, fmt.Errorf("%s does not exist", filename)
	}

	file := new(File)
	_ = json.Unmarshal(fileJSON, file)

	return file, nil
}

func checksum(filePath string) (string, error) {
	f, err := os.Open("/home/chaincode/" + filePath)
	if err != nil {
		return "", err
	}

	defer func() {
		_ = f.Close()
	}()

	copyBuf := make([]byte, 1024*1024)

	h := sha256.New()
	if _, err := io.CopyBuffer(h, f, copyBuf); err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

func getSubmitterName(submitterByte []byte) (string, error) {
	// submitterByte, _ := ctx.GetStub().GetCreator()
	certStart := bytes.IndexAny(submitterByte, "-----BEGIN")
	if certStart == -1 {
		return "", fmt.Errorf("no certificate found")
	}
	certText := submitterByte[certStart:]
	bl, _ := pem.Decode(certText)
	if bl == nil {
		return "", fmt.Errorf("could not decode the PEM structure")
	}
	cert, err := x509.ParseCertificate(bl.Bytes)

	if err != nil {
		return "", fmt.Errorf("parse certificate failed")
	}
	submitterName := cert.Subject.CommonName
	return string([]byte(submitterName)), nil
}

func (s *SmartContract) AddIPInfo(ctx contractapi.TransactionContextInterface, ipInfoJson string) error {

	ipInfoArr := []byte(ipInfoJson)
	var ipInfoInput IP
	err := json.Unmarshal(ipInfoArr, &ipInfoInput)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %s", err.Error())
	}

	ipInfoJSONasBytes, err := json.Marshal(ipInfoInput)
	if err != nil {
		return fmt.Errorf(err.Error())
	}

	// === Save marble to state ===
	err2 := ctx.GetStub().PutState(ipInfoInput.Name, ipInfoJSONasBytes)
	if err2 != nil {
		return fmt.Errorf("failed to put ip: %s", err2.Error())
	}
	return nil
}

// ===============================================
// record file to the chain
// ===============================================
func (s *SmartContract) AddFileInfo(ctx contractapi.TransactionContextInterface, fileInfoJson string) error {

	fileInfoArr := []byte(fileInfoJson)
	var fileInfoInput File
	err := json.Unmarshal(fileInfoArr, &fileInfoInput)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %s", err.Error())
	}

	if fileInfoInput.FileHash != "" {
		return fmt.Errorf("FileHash field should be an empty string")
	}
	if len(fileInfoInput.Description) == 0 {
		return fmt.Errorf("Description field must be a non-empty string")
	}
	if len(fileInfoInput.Name) == 0 {
		return fmt.Errorf("Name field must be a non-empty string")
	}
	if len(fileInfoInput.Rule) == 0 {
		return fmt.Errorf("Owner field must be a non-empty string")
	}
	if fileInfoInput.ObjectType != "file" {
		return fmt.Errorf("ObjectType field must be: file")
	}

	//ensure client is from the same org as the endorser peer
	clientMSPID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("fail to get client mspid")
	}
	log.Print("client mspid: " + clientMSPID)
	peerMSPID, err := shim.GetMSPID()
	if err != nil {
		return fmt.Errorf("fail to get peer id")
	}
	log.Print("peer mspid: " + peerMSPID)

	if clientMSPID != peerMSPID {
		return fmt.Errorf("you cannot add file info through this peer: you belong to different org")
	}

	epBytes, err3 := ctx.GetStub().GetStateValidationParameter(fileInfoInput.Name)
	ep, err3 := statebased.NewStateEP(epBytes)
	orgs := ep.ListOrgs()
	// orgsList, err3 := json.Marshal(orgs)
	if err3 != nil {
		return fmt.Errorf(err3.Error())
	}
	for _, org := range orgs {
		if org != peerMSPID {
			return fmt.Errorf("Please first set the key-level policy to AND(%d "+fileInfoInput.Owner+"); EP: %d "+org, len(fileInfoInput.Owner), len(org))
		}
	}

	// ctx.GetClientIdentity().GetX509Certificate()

	submitterByte, _ := ctx.GetStub().GetCreator()
	clientName, err := getSubmitterName(submitterByte)
	if err != nil {
		return fmt.Errorf("fail to get client id")
	}
	log.Print("client name: " + clientName)

	fileInfoInput.Owner = clientName
	fileInfoInput.PeerMSP = peerMSPID

	filePathStr := clientName
	filePathStr += "/"
	filePathStr += fileInfoInput.Name
	filehashstr, err4 := checksum(filePathStr)
	if err4 != nil {
		return fmt.Errorf(err4.Error())
	}
	fileInfoInput.FileHash = filehashstr

	FileInfoJSONasBytes, err := json.Marshal(fileInfoInput)
	if err != nil {
		return fmt.Errorf(err.Error())
	}

	// === Save marble to state ===
	err2 := ctx.GetStub().PutState(fileInfoInput.Name, FileInfoJSONasBytes)
	if err2 != nil {
		return fmt.Errorf("failed to put File: %s", err2.Error())
	}

	return nil
}

// ===============================================
// query all File d
// ===============================================
func (s *SmartContract) ListAllFiles(ctx contractapi.TransactionContextInterface) ([]*File, error) {

	queryString := fmt.Sprintf("{\"selector\":{\"objecttype\":\"file\"}}")

	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, fmt.Errorf("failed to read all files %s", err.Error())
	}

	defer resultsIterator.Close()

	var files []*File
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var file File
		err = json.Unmarshal(queryResponse.Value, &file)
		if err != nil {
			return nil, err
		}
		files = append(files, &file)
	}

	return files, nil
}

func contains(slice []string, s string) int {
	for index, value := range slice {
		if value == s {
			return index
		}
	}
	return -1
}

// ===============================================
// request File
// ===============================================
func (s *SmartContract) RequestSftp(ctx contractapi.TransactionContextInterface, fileName string) (*Event, error) {

	event := new(Event)

	event.FileName = fileName
	// event.SourceOrg = sourceOrg
	// event.DestinationOrg = destinationOrg
	event.Flag = false

	fileJSON, _ := ctx.GetStub().GetState(fileName)
	if fileJSON == nil {
		return event, fmt.Errorf("%s does not exist. Please choose an existing file.", fileName)
	}

	file := new(File)
	_ = json.Unmarshal(fileJSON, file)

	event.SourceOrg = file.PeerMSP
	event.SourceUser = file.Owner

	clientMSPID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return nil, fmt.Errorf("fail to get client mspid")
	}
	event.DestinationOrg = clientMSPID

	submitterByte, _ := ctx.GetStub().GetCreator()
	clientName, err := getSubmitterName(submitterByte)
	if err != nil {
		return nil, fmt.Errorf("fail to get client id")
	}
	event.DestinationUser = clientName
	event.EventID = fileName + file.Owner + clientName

	permit := contains(file.Rule, clientMSPID)
	if permit != -1 {
		event.Flag = true
	}

	EventInfoJSONasBytes, err := json.Marshal(event)
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}

	err2 := ctx.GetStub().PutState(event.EventID, EventInfoJSONasBytes)
	if err2 != nil {
		return nil, fmt.Errorf("failed to put Event: %s", err2.Error())
	}

	err2 = ctx.GetStub().PutState(fileName+"RequestSftp", []byte("RequestSftp"))
	if err2 != nil {
		return nil, fmt.Errorf("failed to put Event: %s", err2.Error())
	}

	return event, nil
}

const (
	username = "root"
	password = "admin"
)

func sftpconnect(user string, password string, host string, port int) (*sftp.Client, error) {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		sshClient    *ssh.Client
		sftpClient   *sftp.Client
		err          error
	)
	// get auth method
	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(password))

	clientConfig = &ssh.ClientConfig{
		User:            user,
		Auth:            auth,
		Timeout:         30 * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// connet to ssh
	addr = fmt.Sprintf("%s:%d", host, port)

	if sshClient, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return nil, err
	}

	// create sftp client
	if sftpClient, err = sftp.NewClient(sshClient); err != nil {
		return nil, err
	}

	return sftpClient, nil
}

func encryptLargeFiles(userPath string, fileName string) ([]byte, error) {
	localFilePath := "/home/chaincode/" + userPath + "/" + fileName // "/opt/off_state/bigFile.txt"
	log.Print("open file ")
	infile, err := os.Open(localFilePath)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer infile.Close()

	// var key []byte
	log.Print("start key ")
	key := make([]byte, 32)
	_, errkey := rand.Read(key)
	if errkey != nil {
		// println(err.Error())
		return nil, errkey
	}
	// println(key32)
	// key = key32

	block, err := aes.NewCipher(key)
	if err != nil {
		log.Panic(err)
		return nil, err
	}

	// Never use more than 2^32 random nonces with a given key
	// because of the risk of repeat.
	iv := make([]byte, block.BlockSize())
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		log.Fatal(err)
		return nil, err
	}
	log.Print("start open outfile ")
	outfile, err := os.OpenFile("/home/chaincode/"+userPath+"/"+fileName+"_ciphertext.bin", os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer outfile.Close()

	// The buffer size must be multiple of 16 bytes bb
	log.Print("start buffer size ")
	buf := make([]byte, 1024)
	stream := cipher.NewCTR(block, iv)
	for {
		n, err := infile.Read(buf)
		if n > 0 {
			stream.XORKeyStream(buf, buf[:n])
			// Write into file
			outfile.Write(buf[:n])
		}

		if err == io.EOF {
			log.Print("EOF")
			break
		}

		if err != nil {
			log.Printf("Read %d bytes: %v", n, err)
			break
			return nil, err
		}
	}
	// Append the IV
	log.Print("start IV")
	outfile.Write(iv)

	return key, nil
}

func checksumForSign(file string) ([]byte, error) {
	f, err := os.Open("/home/chaincode/off_state/" + file)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = f.Close()
	}()

	copyBuf := make([]byte, 1024*1024)

	h := sha256.New()
	if _, err := io.CopyBuffer(h, f, copyBuf); err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}

func sign(data []byte, sHash crypto.Hash, pvtKey string) ([]byte, error) {

	block, _ := pem.Decode([]byte(pvtKey))
	var rsaPrivateKey *rsa.PrivateKey
	if strings.Index(pvtKey, "BEGIN RSA") > 0 {
		rsaPrivateKey, _ = x509.ParsePKCS1PrivateKey(block.Bytes)
	} else {
		privateKey, _ := x509.ParsePKCS8PrivateKey(block.Bytes)
		rsaPrivateKey = privateKey.(*rsa.PrivateKey)
	}

	sign, err := rsa.SignPKCS1v15(rand.Reader, rsaPrivateKey, sHash, data)
	if err != nil {
		return nil, err
	}
	return sign, nil
}

func verify(data []byte, sign []byte, sHash crypto.Hash, pbkKey string) bool {

	block, _ := pem.Decode([]byte(pbkKey))
	publicKey, _ := x509.ParsePKIXPublicKey(block.Bytes)
	rsaPublicKey := publicKey.(*rsa.PublicKey)

	return rsa.VerifyPKCS1v15(rsaPublicKey, sHash, data, sign) == nil
}

func useRSAKeyEnc(keyVal []byte, plaintext []byte) ([]byte, error) {
	block, _ := pem.Decode(keyVal)
	if block == nil {
		return nil, fmt.Errorf("failed to parse pk block")
	}
	// log.Print("get pk interface")
	pkInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key interface: %s", err.Error())
	}
	pk, isKey := pkInterface.(*rsa.PublicKey)
	if isKey != true {
		return nil, fmt.Errorf("failed to get public key interface")
	}
	// log.Print("start encrypt")
	// log.Print(pk)
	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, pk, plaintext)
	if err != nil {
		// log.Print(err.Error())
		return nil, fmt.Errorf("failed to encrypt kr1: %s", err.Error())
	}
	// log.Print("cypher text")
	// log.Print(ciphertext)

	// err = ctx.GetStub().PutState("org1mspCipher", ciphertext)

	return ciphertext, nil
}

func useRSAKeyDec(cipherText []byte) (string, error) {
	privateKeyFile, err := os.Open("/home/chaincode/ccPrivate.pem")
	defer privateKeyFile.Close()
	if err != nil {
		return "", fmt.Errorf("failed to load sk file")
	}
	privateKeyLoad := make([]byte, 4096)
	num, err := privateKeyFile.Read(privateKeyLoad)
	if err != nil {
		return "", fmt.Errorf("failed to read sk file %s", err.Error())
	}
	privateKey := privateKeyLoad[:num]

	block, _ := pem.Decode(privateKey)
	if block == nil {
		return "", fmt.Errorf("failed to parse sk block")
	}

	var priv *rsa.PrivateKey
	if strings.Index(string(privateKey), "BEGIN RSA") > 0 {
		priv, _ = x509.ParsePKCS1PrivateKey(block.Bytes)
	} else {
		privTemp, _ := x509.ParsePKCS8PrivateKey(block.Bytes)
		priv = privTemp.(*rsa.PrivateKey)
	}

	plaintext, err := rsa.DecryptPKCS1v15(rand.Reader, priv, cipherText)
	if err != nil {
		return "", fmt.Errorf("failed to dec key %s", err.Error())
	}
	return string(plaintext), nil
}

func clientRsaPKEnc(pkval []byte, plaintext []byte) ([]byte, error) {

	block, _ := pem.Decode(pkval)
	if block == nil {
		return nil, fmt.Errorf("failed to parse client pk block")
	}
	pkInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key interface: %s", err.Error())
	}
	pk, isKey := pkInterface.(*rsa.PublicKey)
	if isKey != true {
		return nil, fmt.Errorf("is not a public key")
	}
	ciphertext, _ := rsa.EncryptOAEP(sha1.New(), rand.Reader, pk, plaintext, nil)

	return ciphertext, nil
}

func (s *SmartContract) SendFile_SSH_sftp(ctx contractapi.TransactionContextInterface, filename string) (string, error) {
	fileJSON, _ := ctx.GetStub().GetState(filename)
	if fileJSON == nil {
		return "", fmt.Errorf("File %s does not exist. Please choose an existing file.", filename)
	}
	file := new(File)
	_ = json.Unmarshal(fileJSON, file)

	clientMSPID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return "", fmt.Errorf("fail to get client mspid")
	}

	submitterByte, _ := ctx.GetStub().GetCreator()
	clientName, err := getSubmitterName(submitterByte)
	if err != nil {
		return "", fmt.Errorf("fail to get client id")
	}

	eventid := filename + file.Owner + clientName

	eventJSON, _ := ctx.GetStub().GetState(eventid)
	if eventJSON == nil {
		return "", fmt.Errorf("Event %s does not exist. Please choose an existing file.", eventid)
	}
	event := new(Event)
	_ = json.Unmarshal(eventJSON, event)

	if event.Flag == false {
		return "", fmt.Errorf("This event is not permitted")
	}

	destinationOrg := event.DestinationOrg
	sender := file.PeerMSP

	log.Print("start encrypt file " + filename)
	startEncrypt := time.Now()
	key, errEncrypt := encryptLargeFiles(file.Owner, filename)
	if key == nil {
		return "encryption fail.", fmt.Errorf(errEncrypt.Error())
	}
	if errEncrypt != nil {
		return "encryption fail.", fmt.Errorf(errEncrypt.Error())
	}
	endEncrypt := time.Since(startEncrypt)
	fmt.Println("encrypt time：", endEncrypt)
	log.Print("end encrypt file " + filename)

	log.Print("start sftp " + filename)
	var (
		sftpClient *sftp.Client
	)
	ipJSON, err := ctx.GetStub().GetState("ip." + destinationOrg) //ip.Org2MSP
	if err != nil {
		return "", fmt.Errorf("failed to read ip.%s %s", destinationOrg, err.Error())
	}
	if ipJSON == nil {
		return "", fmt.Errorf("ip.%s does not exist", destinationOrg)
	}
	ip := new(IP)
	_ = json.Unmarshal(ipJSON, ip)

	// ipStr2 := strings.Replace(string(ipStr), " ", "", -1)
	// ipstr := "129.63.205.173:14051"
	//

	// ipStr2 := string(ipStr[0 : len(ipStr)-6])
	// portStr := string(ipStr[len(ipStr)-5 : len(ipStr)])
	// port, _ := strconv.Atoi(portStr)

	startSftp := time.Now()
	// log.Print(ip.IpAdr)
	// log.Print(ip.IpPort)
	if clientMSPID == file.PeerMSP {
		sftpClient, err = sftpconnect(username, password, "localhost", 22)
		if err != nil {
			log.Fatal(err)
			return "", fmt.Errorf("failed to sftpconnect %s", err.Error())
		}
	} else {
		sftpClient, err = sftpconnect(username, password, ip.IpAdr, ip.IpPort)
		if err != nil {
			log.Fatal(err)
			return "", fmt.Errorf("failed to sftpconnect %s", err.Error())
		}
	}
	defer sftpClient.Close()

	localFilePath := "/home/chaincode/" + file.Owner + "/" + filename + "_ciphertext.bin"
	var remoteDir = "/home/chaincode/" + event.DestinationUser
	srcFile, err := os.Open(localFilePath)
	if err != nil {
		log.Fatal(err)
		return "", fmt.Errorf("failed to os.Open(localFilePath) %s", err.Error())
	}
	defer srcFile.Close()

	var remoteFileName = path.Base(localFilePath)
	dstFile, err := sftpClient.Create(path.Join(remoteDir, remoteFileName))
	if err != nil {
		log.Fatal(err)
		return "", fmt.Errorf("failed to sftpClient.Create %s", err.Error())
	}
	defer dstFile.Close()

	// buf := make([]byte, 1024)
	buf := make([]byte, 1024*1024)
	for {
		n, _ := srcFile.Read(buf)
		if n == 0 {
			break
		}
		dstFile.Write(buf[0:n])
	}
	elapsed := time.Since(startSftp)
	fmt.Println("sftp time：", elapsed)

	log.Print("end sftp " + filename)

	log.Print("start keys " + filename)
	// rsaKeyJSON, err1 := ctx.GetStub().GetPrivateData("_implicit_org_"+sender, sender+"_RSA")
	// if err1 != nil {
	// 	return "", fmt.Errorf("failed to read rsa key: %s", err1.Error())
	// }
	// if rsaKeyJSON == nil {
	// 	return "", fmt.Errorf("Rsa key for RSA does not exist!")
	// }

	// rsaKey := new(RSAKey)
	// _ = json.Unmarshal(rsaKeyJSON, rsaKey)

	////sign
	// startSign := time.Now()
	// digest, _ := checksumForSign(filename + "_ciphertext.bin")
	// sign, _ := sign(digest, crypto.SHA256, rsaKey.PrivateKey)
	// endSign := time.Since(startSign)
	// fmt.Println("sign time：", endSign)

	// *******************
	// encrypt key using receiver key and get kr1
	// *******************
	pkReceiver, err := ctx.GetStub().GetState(clientName + "PublicKey")
	// log.Print("pk receiver")
	// log.Print(string(pkReceiver))
	if err != nil {
		return "", fmt.Errorf("failed to get rsa public key for client: %s", err.Error())
	}
	if pkReceiver == nil {
		return "", fmt.Errorf("rsa public key for client does not exit")
	}
	kr1, err := clientRsaPKEnc(pkReceiver, key)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt key using pk of client: %s", err.Error())
	}
	// log.Print("k r1")
	// log.Print(string(kr1))
	// *******************
	// then encrypt kr1 using peer key and get kz
	// *******************
	keyVal, err := ctx.GetStub().GetState(clientMSPID + "PublicKey")
	// log.Print("keyVal")
	// log.Print(string(keyVal))
	if err != nil {
		return "", fmt.Errorf("failed to get the public key")
	}
	if keyVal == nil {
		return "", fmt.Errorf("public key does not exist")
	}
	keyZ, err := useRSAKeyEnc(keyVal, kr1)
	if err != nil {
		return "", fmt.Errorf(err.Error())
	}
	// log.Print("k z")
	// log.Print(keyZ)

	var encryptKey EncryptionKey
	encryptKey.KeyID = "key_" + eventid
	encryptKey.AESKey = keyZ
	encryptKey.EncryptedFileHash, _ = checksum(clientName + "/" + filename + "_ciphertext.bin")
	// encryptKey.RSAPublicKey = rsaKey.PublicKey
	// encryptKey.Signature = sign

	EncryptionKeyInfoJSONasBytes, err3 := json.Marshal(encryptKey)
	if err3 != nil {
		return "", fmt.Errorf(err3.Error())
	}

	err2 := ctx.GetStub().PutPrivateData("_implicit_org_"+sender, encryptKey.KeyID, EncryptionKeyInfoJSONasBytes)
	if err2 != nil {
		return "", fmt.Errorf("failed to set key: %s", err2.Error())
	}
	log.Print("end keys " + filename)
	// log.Print("end")

	err2 = ctx.GetStub().PutState(filename+"SendFile_SSH_sftp", []byte("SendFile_SSH_sftp"))
	if err2 != nil {
		return "", fmt.Errorf("failed to put Event: %s", err2.Error())
	}

	timeEncrypt := strconv.FormatFloat(float64(endEncrypt.Seconds()), 'f', 6, 64)
	timeSftp := strconv.FormatFloat(float64(elapsed.Seconds()), 'f', 6, 64)
	// timeSign := strconv.FormatFloat(float64(endSign.Seconds()), 'f', 6, 64)

	return "copy file to remote server finished!\n" + "key_" + eventid + "\n" + "encrypt time:" + timeEncrypt + "\n" + "sftp time:" + timeSftp + " aes key: " + string(key) + "\n", nil
}

func (s *SmartContract) Test_timeout(ctx contractapi.TransactionContextInterface, t string) (string, error) {

	tt, err := strconv.Atoi(t)

	time.Sleep(time.Duration(tt) * time.Second)

	log.Print("end timeout test")

	return t, err
}

func (s *SmartContract) QueryEnctyptionKey(ctx contractapi.TransactionContextInterface, key string) (string, error) {

	clientMSPID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return "", fmt.Errorf("fail to get client mspid")
	}

	keyByte, err := ctx.GetStub().GetPrivateData("_implicit_org_"+clientMSPID, key)
	if err != nil {
		return "", fmt.Errorf("failed to read encryption key. %d", err.Error())
	}
	if keyByte == nil {
		return "", fmt.Errorf("%s does not exist", keyByte)
	}

	return string(keyByte), nil
}

func decryptLargeFiles(owner string, fileName string, key []byte) error {
	infile, err := os.Open("/home/chaincode/" + owner + "/" + fileName + "_ciphertext.bin")
	if err != nil {
		return err
	}
	defer infile.Close()

	block, err := aes.NewCipher(key)
	if err != nil {
		// log.Panic(err)
		return err
	}

	// Never use more than 2^32 random nonces with a given key
	// because of the risk of repeat.
	fi, err := infile.Stat()
	if err != nil {
		// log.Fatal(err)
		return err
	}

	iv := make([]byte, block.BlockSize())
	msgLen := fi.Size() - int64(len(iv))
	_, err = infile.ReadAt(iv, msgLen)
	if err != nil {
		// log.Fatal(err)
		return err
	}

	outfile, err := os.OpenFile("/home/chaincode/"+owner+"/ecn_"+fileName, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0777)
	if err != nil {
		// log.Fatal(err)
		return err
	}
	defer outfile.Close()

	// The buffer size must be multiple of 16 bytes
	buf := make([]byte, 10240)
	stream := cipher.NewCTR(block, iv)
	for {
		n, err := infile.Read(buf)
		if n > 0 {
			// The last bytes are the IV, don't belong the original message
			if n > int(msgLen) {
				n = int(msgLen)
			}
			msgLen -= int64(n)
			stream.XORKeyStream(buf, buf[:n])
			// Write into file
			outfile.Write(buf[:n])
		}

		if err == io.EOF {
			break
		}

		if err != nil {
			// log.Printf("Read %d bytes: %v", n, err)
			break
			return err
		}
	}
	return nil
}

func (s *SmartContract) CheckReceivedFile(ctx contractapi.TransactionContextInterface, filename string) (bool, error) {
	submitterByte, _ := ctx.GetStub().GetCreator()
	clientName, _ := getSubmitterName(submitterByte)

	_, err := os.Stat("/home/chaincode/" + clientName + "/" + filename + "_ciphertext.bin")
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, err
	}

	// filehashstr, err4 := checksum("/home/chaincode/" + clientName + "/" + filename + "_ciphertext.bin")

	return false, err
}

// request kz
func (s *SmartContract) RequestKey(ctx contractapi.TransactionContextInterface, filename string) (string, error) {
	log.Print("start RequestKey" + filename)

	fileJSON, _ := ctx.GetStub().GetState(filename)
	if fileJSON == nil {
		return "", fmt.Errorf("File %s does not exist. Please choose an existing file.", filename)
	}
	file := new(File)
	_ = json.Unmarshal(fileJSON, file)

	submitterByte, _ := ctx.GetStub().GetCreator()
	clientName, err := getSubmitterName(submitterByte)
	if err != nil {
		return "", fmt.Errorf("fail to get client id")
	}
	eventID := filename + file.Owner + clientName

	sender := file.PeerMSP

	keyByte, err := ctx.GetStub().GetPrivateData("_implicit_org_"+sender, "key_"+eventID)
	if err != nil {
		return "false", fmt.Errorf("failed to read encryption key. %d", err.Error())
	}
	if keyByte == nil {
		return "false", fmt.Errorf("%d does not exist", "key_"+eventID)
	}
	// log.Print("keyByte")
	// log.Print(keyByte)

	errSet := ctx.GetStub().PutPrivateData("collectionKey_org1", "keyShare_"+eventID, keyByte)
	// errSet := ctx.GetStub().PutState("keyShare_"+eventID, keyByte)
	if errSet != nil {
		return "false", fmt.Errorf("faild to set keyShare.", errSet.Error())
	}

	err2 := ctx.GetStub().PutState(filename+"RequestKey", []byte("RequestKey"))
	if err2 != nil {
		return "", fmt.Errorf("failed to put Event: %s", err2.Error())
	}

	log.Print("end RequestKey" + filename)
	return string(keyByte), nil
}

//request kr1
func (s *SmartContract) RequestKeyReceiver(ctx contractapi.TransactionContextInterface, filename string) error {
	// log.Print("start get enc key")

	fileJSON, _ := ctx.GetStub().GetState(filename)
	if fileJSON == nil {
		return fmt.Errorf("File %s does not exist. Please choose an existing file.", filename)
	}
	file := new(File)
	_ = json.Unmarshal(fileJSON, file)

	submitterByte, _ := ctx.GetStub().GetCreator()
	clientName, err := getSubmitterName(submitterByte)
	if err != nil {
		return fmt.Errorf("fail to get client id")
	}
	eventID := filename + file.Owner + clientName

	keyByte, err := ctx.GetStub().GetPrivateData("collectionKey_org1", "keyShare_"+eventID)
	// keyByte, err := ctx.GetStub().GetState("keyShare_"+eventID)
	if err != nil {
		return fmt.Errorf("failed to read encryption key. %d", err.Error())
	}
	if keyByte == nil {
		return fmt.Errorf("%d does not exist", "keyShare_"+eventID)
	}

	encryptionKey := new(EncryptionKey)
	_ = json.Unmarshal(keyByte, encryptionKey)

	log.Print("decrypt kz using peer sk and get kr1")
	kr1, err := useRSAKeyDec(encryptionKey.AESKey)
	if err != nil {
		log.Print(err.Error())
		return fmt.Errorf("failed to dec key. %d", err.Error())
	}
	// log.Print(kr1)

	err2 := ctx.GetStub().PutPrivateData("collectionKey_org1", "keyRShare_"+eventID, []byte(kr1))
	if err2 != nil {
		return fmt.Errorf("failed to set keyRShare: %s", err2.Error())
	}

	err2 = ctx.GetStub().PutState(filename+"RequestKeyReceiver", []byte("RequestKeyReceiver"))
	if err2 != nil {
		return fmt.Errorf("failed to put Event: %s", err2.Error())
	}

	return nil
}

func (s *SmartContract) QueryKeyReceiver(ctx contractapi.TransactionContextInterface, filename string) (string, error) {
	log.Print("start get enc key")

	fileJSON, _ := ctx.GetStub().GetState(filename)
	if fileJSON == nil {
		return "", fmt.Errorf("File %s does not exist. Please choose an existing file.", filename)
	}
	file := new(File)
	_ = json.Unmarshal(fileJSON, file)

	submitterByte, _ := ctx.GetStub().GetCreator()
	clientName, err := getSubmitterName(submitterByte)
	if err != nil {
		return "", fmt.Errorf("fail to get client id")
	}
	eventID := filename + file.Owner + clientName

	keyByte, err := ctx.GetStub().GetPrivateData("collectionKey_org1", "keyRShare_"+eventID)
	if err != nil {
		return "", fmt.Errorf("failed to read encryption key. %d", err.Error())
	}
	if keyByte == nil {
		return "", fmt.Errorf("%d does not exist", "keyRShare_"+eventID)
	}

	return string(keyByte), nil
}

func (s *SmartContract) DecryptFile(ctx contractapi.TransactionContextInterface, filename string) (bool, error) {
	// log.Print("start get enc key")

	fileJSON, _ := ctx.GetStub().GetState(filename)
	if fileJSON == nil {
		return false, fmt.Errorf("File %s does not exist. Please choose an existing file.", filename)
	}
	file := new(File)
	_ = json.Unmarshal(fileJSON, file)

	submitterByte, _ := ctx.GetStub().GetCreator()
	clientName, err := getSubmitterName(submitterByte)
	if err != nil {
		return false, fmt.Errorf("fail to get client id")
	}

	keyMap, err := ctx.GetStub().GetTransient()
	if err != nil {
		return false, fmt.Errorf("Fail to get transient. %d", err.Error())
	}
	kr1, ok := keyMap["key"]
	if !ok {
		return false, fmt.Errorf("Fail to get transient key. %d", err.Error())
	}

	err2 := decryptLargeFiles(clientName, filename, []byte(kr1))
	if err2 != nil {
		return false, fmt.Errorf("Fail to decrypt files. %d", err2.Error())
	}
	log.Print("end decrypt file")

	// check hash
	filePathStr := clientName
	filePathStr += "/"
	filePathStr += "ecn_" + filename
	filehashstr, err4 := checksum(filePathStr)
	if err4 != nil {
		return false, fmt.Errorf(err4.Error())
	}
	if file.FileHash != filehashstr {
		return false, fmt.Errorf("hash does not match")
	}

	err2 = ctx.GetStub().PutState(filename+"decrypt", []byte("decrypt"))
	if err2 != nil {
		return false, fmt.Errorf("failed to put Event: %s", err2.Error())
	}

	return true, nil
}

func (s *SmartContract) QueryKeyAndDecryptFile(ctx contractapi.TransactionContextInterface, filename string) (bool, error) {
	// log.Print("start get enc key")

	fileJSON, _ := ctx.GetStub().GetState(filename)
	if fileJSON == nil {
		return false, fmt.Errorf("File %s does not exist. Please choose an existing file.", filename)
	}
	file := new(File)
	_ = json.Unmarshal(fileJSON, file)

	submitterByte, _ := ctx.GetStub().GetCreator()
	clientName, err := getSubmitterName(submitterByte)
	if err != nil {
		return false, fmt.Errorf("fail to get client id")
	}
	eventID := filename + file.Owner + clientName

	keyByte, err := ctx.GetStub().GetPrivateData("collectionKey_org1", "keyRShare_"+eventID)
	// keyByte, err := ctx.GetStub().GetState("keyShare_"+eventID)
	if err != nil {
		return false, fmt.Errorf("failed to read encryption key. %d", err.Error())
	}
	if keyByte == nil {
		return false, fmt.Errorf("%d does not exist", "keyShare_"+eventID)
	}

	err2 := decryptLargeFiles(clientName, filename, keyByte)
	if err2 != nil {
		return false, fmt.Errorf("Fail to decrypt files. %d", err2.Error())
	}
	log.Print("end decrypt file")

	// check hash
	filePathStr := clientName
	filePathStr += "/"
	filePathStr += "ecn_" + filename
	filehashstr, err4 := checksum(filePathStr)
	if err4 != nil {
		return false, fmt.Errorf(err4.Error())
	}
	if file.FileHash != filehashstr {
		return false, fmt.Errorf("hash does not match")
	}

	return true, nil
}

func (s *SmartContract) DeleteFileInDocker(ctx contractapi.TransactionContextInterface, filePath string) (bool, error) {
	err1 := os.RemoveAll("/home/chaincode/" + filePath)
	if err1 != nil {
		return false, err1
	}
	return true, nil
}

// ==================================================
// query history transactions for file (filehash)
//addddd
// ==================================================
func (s *SmartContract) QueryTxHistoryForFile(ctx contractapi.TransactionContextInterface, filehash string) (string, error) {
	resultsIterator, err := ctx.GetStub().GetHistoryForKey(filehash)
	if err != nil {
		return "", fmt.Errorf(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing historic values for the file
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return "", fmt.Errorf(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"TxId\":")
		buffer.WriteString("\"")
		buffer.WriteString(response.TxId)
		buffer.WriteString("\"")
		// buffer.WriteString(", \"Value\":")
		// if response.IsDelete {
		// 	buffer.WriteString("null")
		// } else {
		// 	buffer.WriteString(string(response.Value))
		// }

		buffer.WriteString(", \"Timestamp\":")
		buffer.WriteString("\"")
		buffer.WriteString(time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String())
		buffer.WriteString("\"")

		// buffer.WriteString(", \"IsDelete\":")
		// buffer.WriteString("\"")
		// buffer.WriteString(strconv.FormatBool(response.IsDelete))
		// buffer.WriteString("\"")

		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- getHistoryForMarble returning:\n%s\n", buffer.String())

	// return shim.Success(buffer.Bytes())
	// return buffer.Bytes(), nil
	return buffer.String(), nil
}

func main() {

	chaincode, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		fmt.Printf("Error creating sector chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting sector chaincode: %s", err.Error())
	}
}
