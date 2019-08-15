package service

import (
	"fmt"
	"net"
	"log"
	"google.golang.org/grpc"
	"monitor/pb"
)

/*
 * rpc service interface
 */

//rpc server info
type RpcServer struct {
	address string `rpc service address`
	service *grpc.Server `rpc service`
}

//construct
func NewRpcServer(port int) *RpcServer {
	address := fmt.Sprintf(":%d", port)
	this := &RpcServer{
		address:address,
		service:nil,
	}
	//create rpc service
	this.createService()
	return this
}

//stop
func (r *RpcServer) Stop() {
	if r.service != nil {
		r.service.Stop()
		log.Println("rpc service stopped.")
	}
}

//create rpc service
func (r *RpcServer) createService() {
	var tips string
	var err error

	//try listen tcp port
	listen, err := net.Listen("tcp", r.address)
	if err != nil {
		tips = "Create rpc service failed, error:" + err.Error()
		log.Println(tips)
		panic(tips)
	}

	//create rpc server with rpc stat support
	r.service = grpc.NewServer(grpc.StatsHandler(NewRpcStat()))

	//register call back
	monitor.RegisterMonitorServiceServer(r.service, &RpcNode{})

	//begin rpc service
	go r.beginService(listen)
}

//begin rpc service
func (r *RpcServer) beginService(listen net.Listener) {
	//service listen
	err := r.service.Serve(listen)
	if err != nil {
		tips := "Failed for rpc service, error:" + err.Error()
		log.Println(tips)
		panic(tips)
	}
}