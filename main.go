package main

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/grpc"
	"grpcCli/helper"
	"grpcCli/services"
	"log"
	"time"
)

func main()  {
	//creds,err := credentials.NewClientTLSFromFile("keys/ssl.crt","kai.com")
	//if err != nil {
	//	log.Fatalln(err)
	//}

	//向指定端口发起连接
	conn,err := grpc.Dial(":8081",grpc.WithTransportCredentials(helper.GetClientCreds()))
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	ctx := context.Background()

	//创建prod service client
	prodClient := services.NewProdServiceClient(conn)
	////获取单个prod
	////根据区域获取，如果不传ProdArea，默认是A
	//prodRes,err := prodClient.GetProdStock(context.Background(),&services.ProdRequest{ProdId: 1,ProdArea: services.ProdAreas_B})
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//fmt.Printf("stock:%d\n",prodRes.ProdStock)

	////获取多个prod
	//prodListRes,err := prodClient.GetProdStocks(context.Background(),&services.QuerySize{Size: 1})
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//
	//fmt.Println(prodListRes.Prods)

	//获取商品信息
	prodInfo,err := prodClient.GetProdInfo(ctx,&services.ProdRequest{ProdId: 1})
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(prodInfo)

	//先创建order service client
	orderClient := services.NewOrderServiceClient(conn)
	//创建订单
	orderReq := &services.OrderRequest{
		OrderMain: &services.OrderMain{
			OrderId:1,
			OrderNo: "111",
			OrderPrice: 9.9,
			UserId: 58,
			OrderTime: &timestamp.Timestamp{Seconds: time.Now().Unix()},
			OrderDetail: []*services.OrderDetail{
				{DetailId: 1,OrderNo: "111",ProdId: 11,ProdPrice: 1.5,ProdNum: 110},
				{DetailId: 2,OrderNo: "111",ProdId: 22,ProdPrice: 18,ProdNum: 8},
				{DetailId: 3,OrderNo: "111",ProdId: 33,ProdPrice: 10.4,ProdNum: 20},
			},
		},
	}

	orderRes,_ := orderClient.NewOrder(ctx,orderReq)
	fmt.Println(orderRes)


	//先创建user service client
	userClient := services.NewUserServiceClient(conn)
	//获取用户积分
	users := make([]*services.UserInfo,0)
	var i int32
	for i=1;i<8;i++ {
		user := &services.UserInfo{UserId: i}
		users = append(users, user)
	}
	userReq := &services.UserRequest{Users: users}
	userRes,_ := userClient.GetUserScore(ctx,userReq)
	fmt.Println(userRes.Users)
}
