package main

import (
	
	"fmt"
	"log"
	"net"
    "sync"
	"server_stream/proto"
    "time"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)
type server struct{
     proto.UnimplementedStreamServiceServer
}
func (*server) FetchResponse(req *proto.Request,stream proto.StreamService_FetchResponseServer) error {
	log.Printf("Fetch response for id: %d",req.Id)
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func (count int64)  {
			defer wg.Done()
			time.Sleep(time.Duration(count) * time.Second)
			result:=fmt.Sprintf("Request #%d For Id:%d", count, req.Id)
		    resp:=proto.Response{
			    Result:result,
		    }
		if err := stream.Send(&resp);err!=nil {
		     	log.Printf("error in sending : %v",err)
		}
		} (int64(i))
	}
     wg.Wait()

	return nil
}
func main()  {
	lis,err:=net.Listen("tcp" ,":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v ",err)
	}
	s:=grpc.NewServer()

   proto.RegisterStreamServiceServer(s,&server{})
   reflection.Register(s)

   log.Println("server started......")
   if err:=s.Serve(lis);err!=nil {
	   log.Fatalf("failed to serve: %v",err)
   }

}