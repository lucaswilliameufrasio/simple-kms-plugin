package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	v1beta1 "lucaswilliameufrasio/simple-kms-plugin/proto"
	"lucaswilliameufrasio/simple-kms-plugin/utils"
	"net"
	"os"

	"google.golang.org/grpc"
)

type KeyManagementServer struct {
	v1beta1.UnimplementedKeyManagementServiceServer
}

var (
	EncryptionSecretKey = os.Getenv("ENCRYPTION_SECRET")
)

func (kms *KeyManagementServer) Encrypt(ctx context.Context, req *v1beta1.EncryptRequest) (*v1beta1.EncryptResponse, error) {
	plainDecoded, err := base64.StdEncoding.DecodeString(string(req.Plain))

	if err != nil {
		fmt.Println(err.Error())

		return nil, err
	}
	plainString := string(plainDecoded)
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
	versionSupported := "v1beta1"
	version := req.Version

	if version != "v1beta1" {
		versionNotSupportedError := fmt.Errorf("VersionNotSupportedError")
		return nil, versionNotSupportedError
	}

	response := &v1beta1.VersionResponse{
		Version:        versionSupported,
		RuntimeName:    "simple-kms-plugin",
		RuntimeVersion: "0.0.1",
	}

	return response, nil
}

func main() {
	fmt.Println(len(EncryptionSecretKey))
	address := "localhost:9997"
	conn, err := net.Listen("tcp", address)

	if err != nil {
		log.Fatal("tcp connection err: ", err.Error())
	}

	grpcServer := grpc.NewServer()

	kmsServer := KeyManagementServer{}

	v1beta1.RegisterKeyManagementServiceServer(grpcServer, &kmsServer)

	fmt.Println("Starting gRPC server at: ", address)

	if err := grpcServer.Serve(conn); err != nil {
		log.Fatal(err)
	}
}
