package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"grpcCli/services"
	"log"
)

func main()  {
	creds,err := credentials.NewClientTLSFromFile("keys/ssl.crt","kai.com")
	if err != nil {
		log.Fatalln(err)
	}
	//向指定端口发起连接
	conn,err := grpc.Dial(":8081",grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	//创建service client
	prodClient := services.NewProdServiceClient(conn)
	//调用方法
	prodRes,err := prodClient.GetProdStock(context.Background(),&services.ProdRequest{ProdId: 1})
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("stock:%d\n",prodRes.ProdStock)
}
