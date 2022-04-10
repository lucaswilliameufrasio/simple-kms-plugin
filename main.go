package main

import (
	"context"
	"fmt"
	"log"
	v1beta1 "lucaswilliameufrasio/simple-kms-plugin/proto"
	"lucaswilliameufrasio/simple-kms-plugin/utils"
	"net"
	"os"
	"strings"

	"google.golang.org/grpc"
)

type KeyManagementServer struct {
	v1beta1.UnimplementedKeyManagementServiceServer
}

var (
	// This is required, so the fallback is empty
	EncryptionSecretKey = utils.GetEnv("ENCRYPTION_SECRET", "")
	ServerFile          = utils.GetEnv("SOCKET_SERVER_FILE", "/tmp/test_kms.sock")
	Protocol            = utils.GetEnv("PROTOCOL", "tcp")
	VersionSupported    = "v1beta1"
)

func (kms *KeyManagementServer) Encrypt(ctx context.Context, req *v1beta1.EncryptRequest) (*v1beta1.EncryptResponse, error) {
	plainString := string(req.Plain)

	encryptedPlain, err := utils.Encrypt(plainString, EncryptionSecretKey)

	if err != nil {
		return nil, err
	}

	response := v1beta1.EncryptResponse{
		Cipher: []byte(encryptedPlain),
	}

	fmt.Println(string(response.Cipher))

	return &response, nil
}

func (kms *KeyManagementServer) Decrypt(ctx context.Context, req *v1beta1.DecryptRequest) (*v1beta1.DecryptResponse, error) {
	cipherString := string(req.Cipher)

	decrypted, err := utils.Decrypt(cipherString, EncryptionSecretKey)

	if err != nil {
		fmt.Println(err.Error())

		return nil, err
	}

	fmt.Printf("decrypted: %s\n", decrypted)

	response := v1beta1.DecryptResponse{
		Plain: []byte(decrypted),
	}
	return &response, nil
}

func (kms *KeyManagementServer) Version(ctx context.Context, req *v1beta1.VersionRequest) (*v1beta1.VersionResponse, error) {
	version := req.Version

	if version != VersionSupported {
		versionNotSupportedError := fmt.Errorf("VersionNotSupportedError")
		return nil, versionNotSupportedError
	}

	response := &v1beta1.VersionResponse{
		Version:        VersionSupported,
		RuntimeName:    "simple-kms-plugin",
		RuntimeVersion: "0.0.1",
	}

	return response, nil
}

func getUnixListenerConnection() net.Listener {
	os.Remove(ServerFile)
	server_addr, err := net.ResolveUnixAddr("unix", ServerFile)

	if err != nil {
		log.Fatal("failed to resolve unix addr")
	}
	conn, err := net.ListenUnix("unix", server_addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	return conn
}

func getTCPListenerConnection() net.Listener {
	address := "localhost:9997"
	conn, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("tcp connection err: ", err.Error())
	}

	fmt.Println("Listening on address: ", address)

	return conn
}

func startGRPCServer() {
	fmt.Println(len(EncryptionSecretKey))

	var conn net.Listener

	if strings.ToLower(Protocol) == "tcp" {
		conn = getTCPListenerConnection()
	} else {
		conn = getUnixListenerConnection()
	}

	grpcServer := grpc.NewServer()

	kmsServer := KeyManagementServer{}

	v1beta1.RegisterKeyManagementServiceServer(grpcServer, &kmsServer)

	if err := grpcServer.Serve(conn); err != nil {
		log.Fatal(err)
	}
}

func main() {
	startGRPCServer()
}
