package rpc

import (
	"google.golang.org/grpc"
)

func Dial(host string, insecure bool) (conn *grpc.ClientConn, err error) {
	var option grpc.DialOption
	if insecure {
		option = grpc.WithInsecure()
	}

	conn, err = grpc.Dial(host, option)
	if err != nil {
		return &grpc.ClientConn{}, err
	}
	return
}
