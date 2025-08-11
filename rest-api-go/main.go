package main

import (
	"crypto/x509"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-gateway/pkg/identity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Config holds the configuration for the application
type Config struct {
	MspID         string `json:"mspID"`
	CryptoPath    string `json:"cryptoPath"`
	CertPath      string `json:"certPath"`
	KeyPath       string `json:"keyPath"`
	TlsCertPath   string `json:"tlsCertPath"`
	PeerEndpoint  string `json:"peerEndpoint"`
	GatewayPeer   string `json:"gatewayPeer"`
	ChannelName   string `json:"channelName"`
	ChaincodeName string `json:"chaincodeName"`
	Port          string `json:"port"`
}

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

var contract *client.Contract
var appConfig Config

func main() {
	log.Println("--- Starting Asset Tracker API ---")

	// Load configuration
	if err := loadConfig(); err != nil {
		log.Fatalf("FATAL: could not load config: %v", err)
	}

	// Set up Fabric Gateway connection
	log.Println("Initializing connection to Fabric Gateway...")
	gw, err := connectToGateway()
	if err != nil {
		log.Fatalf("FATAL: Failed to connect to gateway: %v", err)
	}
	defer gw.Close()
	log.Println("Successfully connected to Fabric Gateway.")

	network := gw.GetNetwork(appConfig.ChannelName)
	contract = network.GetContract(appConfig.ChaincodeName)

	// Set up Gin router
	router := gin.Default()
	setupRoutes(router)

	log.Printf("Starting API server on port %s", appConfig.Port)
	if err := router.Run(":" + appConfig.Port); err != nil {
		log.Fatalf("FATAL: could not start API server: %v", err)
	}
}

func loadConfig() error {
	file, err := os.ReadFile("config.json")
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}
	if err := json.Unmarshal(file, &appConfig); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}
	// Expand the crypto path
	appConfig.CryptoPath = os.ExpandEnv(appConfig.CryptoPath)
	return nil
}

func setupRoutes(router *gin.Engine) {
	router.POST("/assets", createAsset)
	router.GET("/assets/:id", readAsset)
	router.PUT("/assets/:id", updateAsset)
	router.DELETE("/assets/:id", deleteAsset)
	router.GET("/assets", getAllAssets)
}

func connectToGateway() (*client.Gateway, error) {
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
		return nil, fmt.Errorf("failed to create sign implementation: %w", err)
	}

	gw, err := client.Connect(
		id,
		client.WithSign(sign),
		client.WithClientConnection(clientConnection),
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
	// In a real-world scenario, you would want to use TLS
	log.Printf("Connecting to peer at %s without TLS", appConfig.PeerEndpoint)
	connection, err := grpc.Dial(appConfig.PeerEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to dial peer: %w", err)
	}
	return connection, nil
}

func newIdentity() (*identity.X509Identity, error) {
	certPath := filepath.Join(appConfig.CryptoPath, appConfig.CertPath)
	certificate, err := loadCertificate(certPath)
	if err != nil {
		return nil, err
	}
	return identity.NewX509Identity(appConfig.MspID, certificate)
}

func newSign() (identity.Sign, error) {
	keyDir := filepath.Join(appConfig.CryptoPath, appConfig.KeyPath)
	files, err := os.ReadDir(keyDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read key directory %s: %w", keyDir, err)
	}
	if len(files) == 0 {
		return nil, fmt.Errorf("no key files found in directory %s", keyDir)
	}
	privateKeyPEM, err := os.ReadFile(filepath.Join(keyDir, files[0].Name()))
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
		return nil, fmt.Errorf("failed to read certificate file %s: %w", path, err)
	}
	return identity.CertificateFromPEM(certificatePEM)
}

// --- API Handlers ---

func createAsset(c *gin.Context) {
	var newAsset Asset
	if err := c.ShouldBindJSON(&newAsset); err != nil {
		log.Printf("ERROR: Invalid request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	log.Printf("Received request to create asset: %s", newAsset.DEALERID)
	_, err := contract.SubmitTransaction("CreateAsset", newAsset.DEALERID, newAsset.MSISDN, newAsset.MPIN, newAsset.BALANCE, newAsset.STATUS, newAsset.TRANSAMOUNT, newAsset.TRANSTYPE, newAsset.REMARKS)
	if err != nil {
		log.Printf("ERROR: Failed to submit CreateAsset transaction for asset %s: %v", newAsset.DEALERID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to submit transaction: " + err.Error()})
		return
	}

	log.Printf("Successfully created asset: %s", newAsset.DEALERID)
	c.JSON(http.StatusCreated, newAsset)
}

func readAsset(c *gin.Context) {
	id := c.Param("id")
	log.Printf("Received request to read asset: %s", id)

	evaluateResult, err := contract.EvaluateTransaction("ReadAsset", id)
	if err != nil {
		log.Printf("ERROR: Failed to evaluate ReadAsset transaction for asset %s: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read asset: " + err.Error()})
		return
	}

	var asset Asset
	if err := json.Unmarshal(evaluateResult, &asset); err != nil {
		log.Printf("ERROR: Failed to unmarshal asset data for %s: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse asset data: " + err.Error()})
		return
	}

	log.Printf("Successfully read asset: %s", id)
	c.JSON(http.StatusOK, asset)
}

func updateAsset(c *gin.Context) {
	id := c.Param("id")
	var updatedAsset Asset
	if err := c.ShouldBindJSON(&updatedAsset); err != nil {
		log.Printf("ERROR: Invalid request body for update: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	if id != updatedAsset.DEALERID {
		log.Printf("ERROR: Asset ID mismatch in path and body: path=%s, body=%s", id, updatedAsset.DEALERID)
		c.JSON(http.StatusBadRequest, gin.H{"error": "asset ID in path does not match ID in body"})
		return
	}

	log.Printf("Received request to update asset: %s", id)
	_, err := contract.SubmitTransaction("UpdateAsset", updatedAsset.DEALERID, updatedAsset.MSISDN, updatedAsset.MPIN, updatedAsset.BALANCE, updatedAsset.STATUS, updatedAsset.TRANSAMOUNT, updatedAsset.TRANSTYPE, updatedAsset.REMARKS)
	if err != nil {
		log.Printf("ERROR: Failed to submit UpdateAsset transaction for asset %s: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update asset: " + err.Error()})
		return
	}

	log.Printf("Successfully updated asset: %s", id)
	c.JSON(http.StatusOK, updatedAsset)
}

func deleteAsset(c *gin.Context) {
	id := c.Param("id")
	log.Printf("Received request to delete asset: %s", id)

	_, err := contract.SubmitTransaction("DeleteAsset", id)
	if err != nil {
		log.Printf("ERROR: Failed to submit DeleteAsset transaction for asset %s: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete asset: " + err.Error()})
		return
	}

	log.Printf("Successfully deleted asset: %s", id)
	c.Status(http.StatusNoContent)
}

func getAllAssets(c *gin.Context) {
	log.Println("Received request to get all assets")

	evaluateResult, err := contract.EvaluateTransaction("GetAllAssets")
	if err != nil {
		log.Printf("ERROR: Failed to evaluate GetAllAssets transaction: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get all assets: " + err.Error()})
		return
	}

	var assets []*Asset
	if err := json.Unmarshal(evaluateResult, &assets); err != nil {
		log.Printf("ERROR: Failed to unmarshal asset list: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse asset list: " + err.Error()})
		return
	}

	log.Printf("Successfully retrieved %d assets", len(assets))
	c.JSON(http.StatusOK, assets)
}
