package main

import (
	"bytes"
	"context"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-gateway/pkg/identity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	mspID         = "Org1MSP"
	cryptoPath    = "/path/to/crypto/material" // This should be mounted into the container
	certPath      = cryptoPath + "/users/User1@org1.example.com/msp/signcerts/cert.pem"
	keyPath       = cryptoPath + "/users/User1@org1.example.com/msp/keystore/"
	tlsCertPath   = cryptoPath + "/peers/peer0.org1.example.com/tls/ca.crt"
	peerEndpoint  = "localhost:7051"
	gatewayPeer   = "peer0.org1.example.com"
	channelName   = "mychannel"
	chaincodeName = "asset-tracker"
)

var contract *client.Contract

// Asset describes basic details of what makes up a simple asset
type Asset struct {
	DEALERID    string `json:"DEALERID"`
	MSISDN      string `json:"MSISDN"`
	MPIN        string `json:"MPIN"`
	BALANCE     string `json:"BALANCE"`
	STATUS      string `json:"STATUS"`
	TRANSAMOUNT string `json:"TRANSAMOUNT"`
	TRANSTYPE   string `json:"TRANSTYPE"`
	REMARKS     string `json:"REMARKS"`
}

func main() {
	// Set up the Fabric Gateway connection
	gw, err := connectToGateway()
	if err != nil {
		log.Fatalf("Failed to connect to gateway: %v", err)
	}
	defer gw.Close()

	network := gw.GetNetwork(channelName)
	contract = network.GetContract(chaincodeName)

	// Set up the Gin router
	router := gin.Default()

	router.POST("/assets", createAsset)
	router.GET("/assets/:id", readAsset)
	router.PUT("/assets/:id", updateAsset)
	router.DELETE("/assets/:id", deleteAsset)
	router.GET("/assets", getAllAssets)

	router.Run(":8080")
}

func connectToGateway() (*client.Gateway, error) {
	// The gRPC client connection should be shared by all Gateway connections to this endpoint
	clientConnection, err := newGrpcConnection()
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection: %w", err)
	}

	id, err := newIdentity()
	if err != nil {
		return nil, fmt.Errorf("failed to create identity: %w", err)
	}

	sign, err := newSign()
	if err != nil {
		return nil, fmt.Errorf("failed to create sign: %w", err)
	}

	// Create a Gateway connection for a specific client identity
	gw, err := client.Connect(
		id,
		client.WithSign(sign),
		client.WithClientConnection(clientConnection),
		// Default timeouts for different gRPC calls
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to gateway: %w", err)
	}

	return gw, nil
}

func newGrpcConnection() (*grpc.ClientConn, error) {
	certificate, err := loadCertificate(tlsCertPath)
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	certPool.AddCert(certificate)
	// transportCredentials := credentials.NewClientTLSFromCert(certPool, gatewayPeer)

	// connection, err := grpc.Dial(peerEndpoint, grpc.WithTransportCredentials(transportCredentials))
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to create gRPC connection: %w", err)
	// }
    // Bypassing TLS for now
	connection, err := grpc.Dial(peerEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        return nil, fmt.Errorf("failed to create gRPC connection: %w", err)
    }


	return connection, nil
}

func newIdentity() (*identity.X509Identity, error) {
	certificate, err := loadCertificate(certPath)
	if err != nil {
		return nil, err
	}

	return identity.NewX509Identity(mspID, certificate)
}

func newSign() (identity.Sign, error) {
	files, err := os.ReadDir(keyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read key directory: %w", err)
	}
	privateKeyPEM, err := os.ReadFile(keyPath + files[0].Name())

	if err != nil {
		return nil, fmt.Errorf("failed to read private key file: %w", err)
	}

	privateKey, err := identity.PrivateKeyFromPEM(privateKeyPEM)
	if err != nil {
		return nil, err
	}

	return identity.NewPrivateKeySign(privateKey)
}

func loadCertificate(path string) (*x509.Certificate, error) {
	certificatePEM, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read certificate file: %w", err)
	}
	return identity.CertificateFromPEM(certificatePEM)
}

func createAsset(c *gin.Context) {
	var newAsset Asset
	if err := c.ShouldBindJSON(&newAsset); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := contract.SubmitTransaction("CreateAsset", newAsset.DEALERID, newAsset.MSISDN, newAsset.MPIN, newAsset.BALANCE, newAsset.STATUS, newAsset.TRANSAMOUNT, newAsset.TRANSTYPE, newAsset.REMARKS)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to submit transaction: %v", err)})
		return
	}

	c.JSON(http.StatusCreated, newAsset)
}

func readAsset(c *gin.Context) {
	id := c.Param("id")

	evaluateResult, err := contract.EvaluateTransaction("ReadAsset", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to evaluate transaction: %v", err)})
		return
	}

	var asset Asset
	err = json.Unmarshal(evaluateResult, &asset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to unmarshal asset: %v", err)})
		return
	}

	c.JSON(http.StatusOK, asset)
}

func updateAsset(c *gin.Context) {
	id := c.Param("id")
	var updatedAsset Asset
	if err := c.ShouldBindJSON(&updatedAsset); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if id != updatedAsset.DEALERID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "dealer ID in path and body do not match"})
		return
	}

	_, err := contract.SubmitTransaction("UpdateAsset", updatedAsset.DEALERID, updatedAsset.MSISDN, updatedAsset.MPIN, updatedAsset.BALANCE, updatedAsset.STATUS, updatedAsset.TRANSAMOUNT, updatedAsset.TRANSTYPE, updatedAsset.REMARKS)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to submit transaction: %v", err)})
		return
	}

	c.JSON(http.StatusOK, updatedAsset)
}

func deleteAsset(c *gin.Context) {
	id := c.Param("id")

	_, err := contract.SubmitTransaction("DeleteAsset", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to submit transaction: %v", err)})
		return
	}

	c.Status(http.StatusNoContent)
}

func getAllAssets(c *gin.Context) {
	evaluateResult, err := contract.EvaluateTransaction("GetAllAssets")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to evaluate transaction: %v", err)})
		return
	}

	var assets []*Asset
	err = json.Unmarshal(evaluateResult, &assets)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to unmarshal assets: %v", err)})
		return
	}

	c.JSON(http.StatusOK, assets)
}
