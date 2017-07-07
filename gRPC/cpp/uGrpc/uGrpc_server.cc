#include <algorithm>
#include <chrono>
#include <cmath>
#include <iostream>
#include <memory>
#include <string>

#include <grpc/grpc.h>
#include <grpc++/server.h>
#include <grpc++/server_builder.h>
#include <grpc++/server_context.h>
#include <grpc++/security/server_credentials.h>
#include "uGrpc.grpc.pb.h"

using grpc::Server;
using grpc::ServerBuilder;
using grpc::ServerContext;
using grpc::ServerReader;
using grpc::ServerReaderWriter;
using grpc::ServerWriter;
using grpc::Status;
using uGrpc::EmbUArgs;
using uGrpc::EmbUResponse;
using uGrpc::UArgs;
using uGrpc::UResponse;
using uGrpc::UGrpc;
using std::chrono::system_clock;

class UGrpcImpl final : public UGrpc::Service 
{

    public:
        explicit UGrpcImpl() 
        {
        }

        Status GetUResponse(::grpc::ServerContext* context, const ::uGrpc::UArgs* request, ::uGrpc::UResponse* response) override
        {
            if(request->bs() != std::string("Client request GetUResponse()"))
            {
                std::cout << "--- GetUResponse() server get client arg incorrect,request->bs():" << request->bs() << std::endl;
                response->set_bs(std::string("Server response GetUResponse() fail"));
            }
            else
            {
                response->set_bs(std::string("Server response GetUResponse()"));
            }
            return Status::OK;
        }

        Status ListUResponses(::grpc::ServerContext* context, const ::uGrpc::EmbUArgs* request, ::grpc::ServerWriter< ::uGrpc::UResponse>* writer) override
        {
            uGrpc::UResponse resp;

            if (request->lo().bs() == std::string("Client request ListUResponses()"))
            {
                resp.set_bs(std::string("Server response ListUResponses()"));
            }
            else
            {
                resp.set_bs(std::string("Server get args incorrect."));
            }

            for(int i = 0; i < 3; i++)
            {
                resp.set_i32(request->lo().i32() + i);
                writer->Write(resp);
            }

            return Status::OK;
        }

        Status RecordRoute(::grpc::ServerContext* context, ::grpc::ServerReader< ::uGrpc::EmbUArgs>* reader, ::uGrpc::EmbUResponse* response) override
        {
            int i = 0;
            EmbUArgs embUArg;
            for(i = 0; i < 3 && reader->Read(&embUArg); i++)
            {
                if (embUArg.lo().i32() != i || embUArg.lo().bs() != std::string("Client RecordRoute() request"))
                {
                    std::cout<<"--- RecordRoute() err --- , i:"<< i << " embUArg.lo().i32():" << embUArg.lo().i32() << "embUArg.lo().bs():" << embUArg.lo().bs() <<std::endl;
                    break;
                }
                std::cout<<"RecordRoute(), i:"<< i << " embUArg.lo().i32():" << embUArg.lo().i32() <<std::endl;
            }

            if(i == 3)
            {
                response->mutable_lo()->set_i32(i);
                response->mutable_lo()->set_bs(std::string("Server RecordRoute() get success"));
            }
            else
            {
                response->mutable_lo()->set_bs(std::string("Server RecordRoute() get fail"));
            }
            return Status::OK;
        }

        Status RouteChat(::grpc::ServerContext* context, ::grpc::ServerReaderWriter< ::uGrpc::EmbUResponse, ::uGrpc::EmbUArgs>* stream) override
        {

            for(int i = 0; i < 3; i++)
            {
                EmbUArgs embUArg;
                EmbUResponse embUResp;

                stream->Read(&embUArg);

                if(embUArg.lo().i32() != i || embUArg.lo().bs() != std::string("Client RouteChat() request"))
                {
                    std::cout << "--- RouteChat() get args incorrect, i32:" << embUArg.lo().i32() << "lo().bs():" << embUArg.lo().i32() << std::endl;
                    embUResp.mutable_lo()->set_bs(std::string("Server RouteChat() response fail"));
                }
                else
                {
                    embUResp.mutable_lo()->set_bs(std::string("Server RouteChat() response"));
                }
                embUResp.mutable_lo()->set_i32(i);

                stream->Write(embUResp);
            }

            return Status::OK;
        }

    private:

};

void RunServer()
{
    std::string server_address("0.0.0.0:50051");
    UGrpcImpl service;

    ServerBuilder builder;
    builder.AddListeningPort(server_address, grpc::InsecureServerCredentials());
    builder.RegisterService(&service);
    std::unique_ptr<Server> server(builder.BuildAndStart());
    std::cout << "Server listening on " << server_address << std::endl;
    server->Wait();
}

int main(int argc, char** argv)
{
    RunServer();

    return 0;
}

