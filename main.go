package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"grpcCli/services"
	"io/ioutil"
	"log"
)

func main()  {
	//creds,err := credentials.NewClientTLSFromFile("keys/ssl.crt","kai.com")
	//if err != nil {
	//	log.Fatalln(err)
	//}

	//加载客户端证书
	cert,_ := tls.LoadX509KeyPair("cert/client.pem","cert/client.key")
	//加载ca证书，双向验证，客户端ca证书主要用来验证服务端证书是否合法
	certPool := x509.NewCertPool()
	ca,_ := ioutil.ReadFile("cert/ca.pem")
	certPool.AppendCertsFromPEM(ca)

	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		ServerName: "localhost",
		RootCAs: certPool,
	})

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
