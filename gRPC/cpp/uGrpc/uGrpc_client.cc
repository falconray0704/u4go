#include <chrono>
#include <iostream>
#include <memory>
#include <random>
#include <string>
#include <thread>

#include <grpc/grpc.h>
#include <grpc++/channel.h>
#include <grpc++/client_context.h>
#include <grpc++/create_channel.h>
#include <grpc++/security/credentials.h>
#include "uGrpc.grpc.pb.h"

using grpc::Channel;
using grpc::ClientContext;
using grpc::ClientReader;
using grpc::ClientReaderWriter;
using grpc::ClientWriter;
using grpc::Status;

using uGrpc::EmbUArgs;
using uGrpc::EmbUResponse;
using uGrpc::UArgs;
using uGrpc::UResponse;
using uGrpc::UGrpc;

class UGrpcClient
{
    public:
        UGrpcClient(std::shared_ptr<Channel> channel) : stub_(UGrpc::NewStub(channel))
        {
        }

        void GetUResponse()
        {
            UArgs uArg;
            UResponse uResp;

            uArg.set_bs(std::string("Client request GetUResponse()"));

            GetOneResponse(uArg, &uResp);

            if( uResp.bs() == std::string("Server response GetUResponse()"))
            {
                std::cout << "GetUResponse() success." << std::endl;
            }
            else
            {
                std::cout << "--- GetUResponse() fail:" << std::endl;
                std::cout << "--- Server response uResp.bs():" << uResp.bs() << std::endl;
            }

        }

        void ListUResponses()
        {
            ClientContext context;
            EmbUArgs embUArgsReq;
            UResponse uResponseResp;

            embUArgsReq.mutable_lo()->set_bs(std::string("Client request ListUResponses()"));
            std::unique_ptr<ClientReader<UResponse> > reader( stub_->ListUResponses(&context, embUArgsReq));
            for(int i = 0; i < 3 && reader->Read(&uResponseResp); i++)
            {
                std::cout << "ListUResponses() found UResponse called :" << std::endl;
                if (uResponseResp.i32() == i && uResponseResp.bs() == std::string("Server response ListUResponses()"))
                {
                    std::cout << "ListUResponses() success." << std::endl;
                }
                else
                {
                    std::cout << "--- ListUResponses() fail ---" << std::endl;
                    std::cout << "--- uResponseResp.i32():" << uResponseResp.i32() << "uResponseResp.bs():" << uResponseResp.bs() << std::endl;
                }
            }
            Status status = reader->Finish();
            if (status.ok())
            {
                std::cout << "ListFeatures rpc succeeded." << std::endl;
            }
            else
            {
                std::cout << "ListFeatures rpc failed." << std::endl;
            }
        }

        void RecordRoute()
        {
            int i = 0;
            ClientContext context;
            EmbUResponse embResp;

            std::unique_ptr<ClientWriter<EmbUArgs> > writer( stub_->RecordRoute(&context, &embResp));

            for (i = 0; i < 3; i++)
            {
                EmbUArgs embUArg;

                embUArg.mutable_lo()->set_i32(i);
                embUArg.mutable_lo()->set_bs(std::string("Client RecordRoute() request"));

                if (!writer->Write(embUArg))
                {
                    // Broken stream.
                    break;
                }
            }

            writer->WritesDone();
            Status status = writer->Finish();

            if (status.ok())
            {
                if(embResp.lo().i32() != 3 || embResp.lo().bs() != std::string("Server RecordRoute() get success"))
                {
                    std::cout << "--- Client RecordRoute() err ---, embResp.lo().i32():" << embResp.lo().i32()
                                << "embResp.lo().bs():" << embResp.lo().bs() << std::endl;
                }
                else
                {
                    std::cout << "Client RecordRoute() success" << std::endl;
                }
            }
            else
            {
              std::cout << "Client RecordRoute() rpc failed." << std::endl;
            }
        }

        void RouteChat()
        {
            ClientContext context;

            std::shared_ptr<ClientReaderWriter<EmbUArgs, EmbUResponse> > stream( stub_->RouteChat(&context));

            std::thread writer([stream]() {
                int i = 0;
                for (i = 0; i < 3; i++)
                {
                    EmbUArgs embUArg;
                    embUArg.mutable_lo()->set_i32(i);
                    embUArg.mutable_lo()->set_bs(std::string("Client RouteChat() request"));
                    stream->Write(embUArg);
                }
                stream->WritesDone();
            });

            EmbUResponse embUResp;
            int i = 0;
            while (stream->Read(&embUResp)) 
            {
                if(embUResp.lo().i32() != i ||
                        embUResp.lo().bs() != std::string("Server RouteChat() response"))
                {
                   std::cout << "--- Server RouteChat() response incorrect ---" << std::endl;
                }
                i++;
            }
            writer.join();
            Status status = stream->Finish();
            if (!status.ok()) 
            {
                std::cout << "RouteChat rpc failed." << std::endl;
            }
            else
            {
                    std::cout << "Client RouteChat() success" << std::endl;
            }
        }


    private:
        bool GetOneResponse(const UArgs& uArgs, UResponse* resp)
        {
            ClientContext context;
            Status status = stub_->GetUResponse(&context, uArgs, resp);
            if (!status.ok())
            {
                std::cout << "GetUResponse rpc failed." << std::endl;
                return false;
            }
            if (!resp->has_uargs())
            {
                std::cout << "Server returns incomplete feature." << std::endl;
                return false;
            }
            if (resp->uargs().bs().empty())
            {
                std::cout << "resp->uargs().bs().empty() == true" << std::endl;
            }
            else
            {
                std::cout << "resp->uargs().bs():" << resp->uargs().bs() << std::endl;
            }
            return true;
        }

        const float kCoordFactor_ = 10000000.0;
        std::unique_ptr<UGrpc::Stub> stub_;
};

int main(int argc, char** argv)
{
    UGrpcClient uGrpcClient( grpc::CreateChannel("localhost:50051", grpc::InsecureChannelCredentials()));

    std::cout << "-------------- GetFeature --------------" << std::endl;
    uGrpcClient.GetUResponse();
    std::cout << "-------------- ListFeatures --------------" << std::endl;
    uGrpcClient.ListUResponses();
    std::cout << "-------------- RecordRoute --------------" << std::endl;
    uGrpcClient.RecordRoute();
    std::cout << "-------------- RouteChat --------------" << std::endl;
    uGrpcClient.RouteChat();

    return 0;
}

