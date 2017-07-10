package main

import (
	//"encoding/json"
	"flag"
	"fmt"
	//"io"
	//"io/ioutil"
	//"math"
	"net"
	//"time"
	"bytes"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	//"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"

	//"github.com/golang/protobuf/proto"

	pb "github.com/falconray0704/u4go/gRPC/go/uGrpc"
)

var (
	port       = flag.Int("port", 50051, "The server port")
)

type uGrpcServer struct {
}

// GetFeature returns the feature at the given point.
//func (s *uGrpcServer) GetFeature(ctx context.Context, point *pb.Point) (*pb.Feature, error) {
func (s *uGrpcServer) GetUResponse(ctx context.Context, uArgs *pb.UArgs) (*pb.UResponse, error) {

	if bytes.Compare(uArgs.BS, []byte("Client request GetUResponse()")) != 0 {
		fmt.Printf("--- GetUResponse() server get client arg incorrect,request->bs():%s \n",uArgs.BS)
		return &pb.UResponse{BS:[]byte("Server response GetUResponse() fail")}, nil
	}

	return &pb.UResponse{BS:[]byte("Server response GetUResponse()")}, nil
}

// ListFeatures lists all features contained within the given bounding Rectangle.
//func (s *uGrpcServer) ListFeatures(rect *pb.Rectangle, stream pb.RouteGuide_ListFeaturesServer) error {
func (s *uGrpcServer) ListUResponses(embUArgs *pb.EmbUArgs, stream pb.UGrpc_ListUResponsesServer) error {

	var i int32
	var resp pb.UResponse

	//resp.BS = []byte{}

	if bytes.Compare(embUArgs.Lo.BS, []byte("Client request ListUResponses()")) == 0 {
		resp.BS = []byte("Server response ListUResponses()")
	} else {
		resp.BS = []byte("Server get args incorrect.")
	}

	for i = 0; i < 3 ; i++ {
		resp.I32 = embUArgs.Lo.I32 + i
		if err := stream.Send(&resp); err != nil {
			return err
		}
	}
	return nil
}

// RecordRoute records a route composited of a sequence of points.
//
// It gets a stream of points, and responds with statistics about the "trip":
// number of points,  number of known features visited, total distance traveled, and
// total time spent.
//func (s *uGrpcServer) RecordRoute(stream pb.RouteGuide_RecordRouteServer) error {
func (s *uGrpcServer) RecordRoute(stream pb.UGrpc_RecordRouteServer) error {

	var i int32
	var err error
	//var embUArg = &pb.EmbUArgs{&pb.UArgs{}, &pb.UArgs{}}
	var embUArg *pb.EmbUArgs

	for i = 0; i < 3; i++ {
		embUArg, err = stream.Recv()
		if err != nil {
			return err
		}
		if embUArg.Lo.I32 != i || bytes.Compare(embUArg.Lo.BS, []byte("Client RecordRoute() request")) != 0 {
			fmt.Printf("--- RecordRoute() err --- , i:%d embUArg.lo().i32():%d embUArg.lo().bs():%s \n",
			i, embUArg.Lo.I32, embUArg.Lo.BS)
			break;
		}
		fmt.Printf("RecordRoute(), i:%d embUArg.lo().i32():%d \n", i, embUArg.Lo.I32)
	}
	var resp = pb.EmbUResponse{Lo: &pb.UArgs{}}
	if i == 3 {
		resp.Lo.I32 = i
		resp.Lo.BS = []byte("Server RecordRoute() get success")
	} else {
		resp.Lo.BS = []byte("Server RecordRoute() get fail")
	}

	return stream.SendAndClose(&resp)

}

// RouteChat receives a stream of message/location pairs, and responds with a stream of all
// previous messages at each of those locations.
//func (s *uGrpcServer) RouteChat(stream pb.RouteGuide_RouteChatServer) error {
func (s *uGrpcServer) RouteChat(stream pb.UGrpc_RouteChatServer) error {

	var i int32
	for i = 0; i < 3; i++ {
		var resp = pb.EmbUResponse{Lo: &pb.UArgs{}}
		embUArg, err := stream.Recv()
		if err != nil {
			return err
		}

		if embUArg.Lo.I32 != i || bytes.Compare(embUArg.Lo.BS, []byte("Client RouteChat() request")) != 0 {
			fmt.Printf("--- RouteChat() get args incorrect, i32:%d lo().bs():%s \n",
			embUArg.Lo.I32, embUArg.Lo.BS)

			resp.Lo.BS = []byte("Server RouteChat() response fail")
		} else {
			resp.Lo.BS = []byte("Server RouteChat() response")
		}

		resp.Lo.I32 = i

		if err := stream.Send(&resp); err != nil {
			return err
		}
	}
	return nil
}

func newServer() *uGrpcServer {
	s := new(uGrpcServer)
	return s
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		grpclog.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterUGrpcServer(grpcServer, newServer())
	grpcServer.Serve(lis)
}
