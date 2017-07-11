
package main

import (
	"flag"
	"io"
	//"math/rand"
	//"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	pb "github.com/falconray0704/u4go/gRPC/go/uGrpc"

	"google.golang.org/grpc/grpclog"
	"fmt"
	"bytes"
	//"golang.org/x/tools/go/gcimporter15/testdata"
)

var (
	serverAddr         = flag.String("server_addr", "127.0.0.1:50051", "The server address in the format of host:port")
)

// printFeature gets the feature for the given point.
func printUResponse(client pb.UGrpcClient) {

	var uArg = pb.UArgs{BS: []byte("Client request GetUResponse()")}

	resp, err := client.GetUResponse(context.Background(), &uArg)
	if err != nil {
		grpclog.Fatalf("--- %v.GetUResponse(ctx, %v) = _, %v ---", client, uArg, err)
	} else {
		if resp.UArgs == nil {
			fmt.Printf("--- Server returns incomplete feature. ---\n")
		} else {
			if resp.UArgs.BS == nil {
				fmt.Printf("--- resp.UArgs.BS is nil ---\n")
			} else {
				grpclog.Println(resp)
			}
		}
		fmt.Printf("GetUResponse() success.")
	}
}

// printFeatures lists all the features within the given bounding Rectangle.
func printUResponses(client pb.UGrpcClient) {

	var embArgs = &pb.EmbUArgs{Lo: &pb.UArgs{BS: []byte("Client request ListUResponses()")}}
	//grpclog.Printf("Looking for features within %v", rect)
	stream, err := client.ListUResponses(context.Background(),embArgs)
	if err != nil {
		grpclog.Fatalf("--- %v.ListResponses(ctx, %%v) = _, %v ---", client, embArgs, err)
	} else {
		var i int32
		for i = 0; i < 3; i++{
			uResp, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				grpclog.Fatalf("--- %v.ListFeatures(ctx, %v) = _, %v", client, embArgs, err)
			} else {
				if uResp != nil && uResp.I32 == i && uResp.BS != nil &&
					bytes.Compare(uResp.BS, []byte("Server response ListUResponses()")) == 0 {
					fmt.Printf("ListUResponses() success.\n")
				} else {
					fmt.Printf("--- ListUResponses() fail uResp:%v ---\n", uResp)
				}
			}
		}
		if i != 3 {
			fmt.Printf("--- ListUResponses() fail, i: %d ---", i)
		}
	}
}

// runRecordRoute sends a sequence of points to server and expects to get a RouteSummary from server.
func runRecordRoute(client pb.UGrpcClient) {
	var i int32
	stream, err := client.RecordRoute(context.Background())
	if err != nil {
		grpclog.Fatalf("--- %v.RecordRoute(_) = _, %v ---", client, err)
	} else {
		for i = 0; i < 3; i++ {
			var embUArg = pb.EmbUArgs{Lo: &pb.UArgs{I32: i, BS: []byte("Client RecordRoute() request")}}
			if err := stream.Send(&embUArg); err != nil {
				grpclog.Fatalf("--- %v.Send(%v) = %v ---", stream, embUArg, err)
			}
		}

		reply, err := stream.CloseAndRecv()
		if err != nil {
			grpclog.Fatalf("--- %v.CloseAndRecv() got error %v, ---", stream, err)
		} else {
			if reply.Lo.I32 != 3 || bytes.Compare(reply.Lo.BS, []byte("Server RecordRoute() get success")) != 0 {
				grpclog.Printf("--- resp: %v ---", reply)
			} else {
				fmt.Printf("Client RecordRoute() success \n")
			}
		}

	}
}

// runRouteChat receives a sequence of route notes, while sending notes for various locations.
func runRouteChat(client pb.UGrpcClient) {

	stream, err := client.RouteChat(context.Background())
	if err != nil {
		grpclog.Fatalf("--- %v.RouteChat(_) = _, %v", client, err)
	} else {
		waitc := make(chan struct{})
		go func() {
			var i int32
			for i = 0; i < 3; i++ {
				resp, err := stream.Recv()
				if err != nil {
					grpclog.Fatalf("--- Failed to receive a uResponse : %v", err)
				} else if resp.Lo.I32 != i || bytes.Compare(resp.Lo.BS, []byte("Server RouteChat() response")) != 0 {
					fmt.Printf("--- Server RouteChat() response incorrect, resp:%v ---", resp)
				}
			}
			fmt.Printf("runRouteChat() success \n")
			close(waitc)
		}()

		var i int32
		for i = 0; i < 3; i++ {
			var embUArg = pb.EmbUArgs{Lo: &pb.UArgs{I32:i, BS: []byte("Client RouteChat() request")}}

			if err := stream.Send(&embUArg); err != nil {
				grpclog.Fatalf("Failed to send a note: %v", err)
			}
		}
		stream.CloseSend()
		<-waitc
	}
}

func main() {
	flag.Parse()
	var opts []grpc.DialOption

	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial(*serverAddr, opts...)

	if err != nil {
		grpclog.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewUGrpcClient(conn)

	// Looking for a valid feature
	fmt.Printf("+++ printUResponse() +++\n")
	printUResponse(client)

	// Feature missing.
	fmt.Printf("+++ printUResponse() +++\n")
	printUResponses(client)

	// RecordRoute
	fmt.Printf("+++ runRecordRoute() +++\n")
	runRecordRoute(client)

	// RouteChat
	fmt.Printf("+++ runRouteChat() +++\n")
	runRouteChat(client)
}
