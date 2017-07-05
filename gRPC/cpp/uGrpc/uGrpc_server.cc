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
#include "helper.h"
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

/*
std::string GetFeatureName(const Point& point, const std::vector<Feature>& feature_list) 
{
  for (const Feature& f : feature_list) {
    if (f.location().latitude() == point.latitude() &&
        f.location().longitude() == point.longitude()) {
      return f.name();
    }
  }
  return "";
}
*/

class UGrpcImpl final : public UGrpc::Service 
{
 public:
  explicit UGrpcImpl() 
  {
    //routeguide::ParseDb(db, &feature_list_);
  }
  //Status GetFeature(ServerContext* context, const Point* point, Feature* feature) override 
  Status GetUResponse(::grpc::ServerContext* context, const ::uGrpc::UArgs* request, ::uGrpc::UResponse* response) override
  {
    //feature->set_name(GetFeatureName(*point, feature_list_));
    //feature->mutable_location()->CopyFrom(*point);

// int32 i32 = 1;
    response->set_i32(request->i32() + 1);

    return Status::OK;
  }

  //Status ListFeatures(ServerContext* context, const uGrpc::Rectangle* rectangle, ServerWriter<Feature>* writer) override 
Status ListUResponses(::grpc::ServerContext* context, const ::uGrpc::EmbUArgs* request, ::grpc::ServerWriter< ::uGrpc::UResponse>* writer) override
{
    uGrpc::UResponse resp;

    for(int i = 0; i < 3; i++)
    {
        resp.set_i32(request->lo().i32() + i);
        writer->Write(resp);
    }

    /*
    auto lo = rectangle->lo();
    auto hi = rectangle->hi();
    long left = (std::min)(lo.longitude(), hi.longitude());
    long right = (std::max)(lo.longitude(), hi.longitude());
    long top = (std::max)(lo.latitude(), hi.latitude());
    long bottom = (std::min)(lo.latitude(), hi.latitude());
    for (const Feature& f : feature_list_) {
      if (f.location().longitude() >= left &&
          f.location().longitude() <= right &&
          f.location().latitude() >= bottom &&
          f.location().latitude() <= top) {
        writer->Write(f);
      }
    }
    */
    return Status::OK;
  }

  //Status RecordRoute(ServerContext* context, ServerReader<Point>* reader, RouteSummary* summary) override 
Status RecordRoute(::grpc::ServerContext* context, ::grpc::ServerReader< ::uGrpc::EmbUArgs>* reader, ::uGrpc::EmbUResponse* response) override
{

    EmbUArgs embUArg;
    for(int i = 0; i < 3; i++)
    {
        std::cout<<"RecordRoute(), i:"<< i << "embUArg.lo.i32" << embUArg.lo().i32() <<std::endl;
    }
    ////////////////
    /*
    Point point;
    int point_count = 0;
    int feature_count = 0;
    float distance = 0.0;
    Point previous;

    system_clock::time_point start_time = system_clock::now();
    while (reader->Read(&point)) {
      point_count++;
      if (!GetFeatureName(point, feature_list_).empty()) {
        feature_count++;
      }
      if (point_count != 1) {
        distance += GetDistance(previous, point);
      }
      previous = point;
    }
    system_clock::time_point end_time = system_clock::now();
    summary->set_point_count(point_count);
    summary->set_feature_count(feature_count);
    summary->set_distance(static_cast<long>(distance));
    auto secs = std::chrono::duration_cast<std::chrono::seconds>(
        end_time - start_time);
    summary->set_elapsed_time(secs.count());
    */

    return Status::OK;
}

//  Status RouteChat(ServerContext* context, ServerReaderWriter<RouteNote, RouteNote>* stream) override 
Status RouteChat(::grpc::ServerContext* context, ::grpc::ServerReaderWriter< ::uGrpc::EmbUResponse, ::uGrpc::EmbUArgs>* stream) override
{

    for(int i = 0; i < 3; i++)
    {
        EmbUArgs embUArg;
        EmbUResponse embUResp;

        stream->Read(&embUArg);
        embUResp.mutable_lo()->set_i32(embUArg.lo().i32() + 1);

        stream->Write(embUResp);
    }
    /*
    std::vector<RouteNote> received_notes;
    RouteNote note;
    while (stream->Read(&note)) {
      for (const RouteNote& n : received_notes) {
        if (n.location().latitude() == note.location().latitude() &&
            n.location().longitude() == note.location().longitude()) {
          stream->Write(n);
        }
      }
      received_notes.push_back(note);
    }
    */

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
  // Expect only arg: --db_path=path/to/route_guide_db.json.
  //std::string db = uGrpc::GetDbFileContent(argc, argv);
  RunServer();

  return 0;
}











