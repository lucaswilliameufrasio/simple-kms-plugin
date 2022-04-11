package main

import (
	"context"
	"log"
	v1beta1 "lucaswilliameufrasio/simple-kms-plugin/proto"
	"lucaswilliameufrasio/simple-kms-plugin/utils"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	ServerAddress = utils.GetEnv("SERVER_ADDRESS", "localhost:9997")
)

func TCPConnect(addr string, t time.Duration) (net.Conn, error) {
	tcp_address, err := net.ResolveTCPAddr("tcp", ServerAddress)
	if err != nil {
		log.Fatalf("Failed to resolve unix address: %v", err)
	}
	conn, err := net.DialTCP("tcp", nil, tcp_address)
	return conn, err
}

func encrypt(client v1beta1.KeyManagementServiceClient) {
	name := "world"

	response, err := client.Encrypt(context.Background(), &v1beta1.EncryptRequest{Plain: []byte(utils.Encode([]byte(name)))})
	if err != nil {
		log.Fatal("could not encrypt: ", err)
	}
	log.Println(response)
	log.Printf("Cipher: %s", response.Cipher)
}

func decrypt(client v1beta1.KeyManagementServiceClient) {
	cipher := "PE9uTqE="

	r, err := client.Decrypt(context.Background(), &v1beta1.DecryptRequest{Cipher: []byte(cipher)})
	if err != nil {
		log.Fatal("could not decrypt: ", err)
	}
	log.Println(r)
	log.Printf("Plain: %s", r.Plain)
}

func GrpcClient() {
	log.Println("start tcp gprc client")

	conn, err := grpc.Dial(ServerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDialer(TCPConnect))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := v1beta1.NewKeyManagementServiceClient(conn)

	encrypt(client)
	decrypt(client)
}

func main() {
	GrpcClient()
}
