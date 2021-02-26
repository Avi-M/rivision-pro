package main

import (
	"grpc_chat/proto"

	"log"
	"net"
	"os"
	"sync"
	"google.golang.org/grpc/reflection" 
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
	glog "google.golang.org/grpc/grpclog"
)

var grpcLog glog.LoggerV2

func init(){
	grpcLog=glog.NewLoggerV2(os.Stdout,os.Stdout,os.Stdout)
}

type Connection struct{
	stream proto.Broadcast_CreateStreamServer
	id string
	active bool
	error chan error
}

type Server struct{
	proto.UnimplementedBroadcastServer
	Connection []*Connection
} 

func (s *Server) CreateStream(pconn *proto.Connect,stream proto.Broadcast_CreateStreamServer)error{
	conn:=&Connection{
		stream:stream,
		id: pconn.User.Id,
		active:true,
		error:make(chan error),
	}
	s.Connection=append(s.Connection,conn)
	return <-conn.error
}


func (s *Server) BroadcastMessage(ctx context.Context,msg *proto.Message)(*proto.Close,error){
	wait:=sync.WaitGroup{}
	done:=make(chan int)

	for _,conn:=range s.Connection{
		wait.Add(1)

		go func(msg *proto.Message,conn *Connection){
			defer wait.Done()

			if conn.active{
				err:=conn.stream.Send(msg)
				grpcLog.Info("Sending message to:",conn.stream)

				if err!=nil{
					grpcLog.Errorf("Error with Stream :%v - Error:%v",conn.stream,err)
					conn.active=false
					conn.error<-err
				}
			}
		}(msg,conn)
	}

	go func(){
		wait.Wait()
		close(done)
	}()
	<-done
	return &proto.Close{},nil

}



func main(){
	var connection []*Connection

	server:=&Server{Connection:connection}

	grpcServer :=grpc.NewServer()

	listener, err:=net.Listen("tcp",":8080")

	if err!=nil{
		log.Fatalf("Error Creating the server %v",err)
	}

	grpcLog.Info("Starting server at port :8080")

	proto.RegisterBroadcastServer(grpcServer,server)

		
// Register reflection service on gRPC server.
   reflection.Register(grpcServer)
	grpcServer.Serve(listener)
}