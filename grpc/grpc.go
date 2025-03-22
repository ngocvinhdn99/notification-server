package grpc

import (
	"context"
	"log"
	"net"
	"time"

	pb "draft-notification/proto"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Server
type server struct {
	pb.UnimplementedMessengerServer
}

func (s *server) SendMessage(ctx context.Context, req *pb.MessageRequest) (*pb.MessageResponse, error) {
	log.Printf("Received message: %v", req.Content)
	return &pb.MessageResponse{Result: "Vinh received Message: " + req.Content}, nil
}

func runServer() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterMessengerServer(s, &server{})

	log.Println("gRPC server is running on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func runClient() {
	conn, err :=
		grpc.NewClient(
			"localhost:50051",
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
		)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewMessengerClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.SendMessage(ctx, &pb.MessageRequest{Content: "Hello from client!"})
	if err != nil {
		log.Fatalf("could not send message: %v", err)
	}
	log.Printf("Server response: %s", r.Result)
}

func RunGrpc() {
	// Start server in a goroutine
	go runServer()

	// Wait for server to start
	time.Sleep(2 * time.Second)

	// Run client
	runClient()

	// Keep the server running
	select {}
}
