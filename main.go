package main

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/grpc"
	"grpcCli/helper"
	"grpcCli/services"
	"io"
	"log"
	"sync"
	"time"
)

func main() {
	//creds,err := credentials.NewClientTLSFromFile("keys/ssl.crt","kai.com")
	//if err != nil {
	//	log.Fatalln(err)
	//}

	//向指定端口发起连接
	conn, err := grpc.Dial(":8081", grpc.WithTransportCredentials(helper.GetClientCreds()))
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
	prodInfo, err := prodClient.GetProdInfo(ctx, &services.ProdRequest{ProdId: 1})
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(prodInfo)

	//先创建order service client
	orderClient := services.NewOrderServiceClient(conn)
	//创建订单
	orderReq := &services.OrderRequest{
		OrderMain: &services.OrderMain{
			OrderId:    1,
			OrderNo:    "111",
			OrderPrice: 9.9,
			UserId:     58,
			OrderTime:  &timestamp.Timestamp{Seconds: time.Now().Unix()},
			OrderDetail: []*services.OrderDetail{
				{DetailId: 1, OrderNo: "111", ProdId: 11, ProdPrice: 1.5, ProdNum: 110},
				{DetailId: 2, OrderNo: "111", ProdId: 22, ProdPrice: 18, ProdNum: 8},
				{DetailId: 3, OrderNo: "111", ProdId: 33, ProdPrice: 10.4, ProdNum: 20},
			},
		},
	}

	orderRes, _ := orderClient.NewOrder(ctx, orderReq)
	fmt.Println(orderRes)

	//先创建user service client
	userClient := services.NewUserServiceClient(conn)
	//普通模式，客户端一次性全部发送，服务端一次性全部响应
	users := make([]*services.UserInfo, 0)
	var i int32
	for i = 1; i < 8; i++ {
		user := &services.UserInfo{UserId: i}
		users = append(users, user)
	}
	userReq := &services.UserRequest{Users: users}
	userRes, _ := userClient.GetUserScore(ctx, userReq)
	fmt.Printf("普通模式:%v\n", userRes.Users)

	//服务端开启流模式，客户端一次性全部发送，服务端分批响应
	stream, err := userClient.GetUserScoreByServerStream(ctx, userReq)
	if err != nil {
		log.Fatalln(err)
	}
	//循环从流中读取数据
	for {
		uRes, err := stream.Recv()
		//如果服务端推送结束，会返回io.EOF，那客户端就break
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln(err)
		}
		//这里可以启用goroutines去处理返回的数据
		fmt.Printf("服务端流模式:%v\n", uRes.Users)
	}

	//客户端开启流模式，客户端分批发送，服务端一次性全部响应
	//获取客户端流
	cliStream, err := userClient.GetUserScoreByClientStream(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	//每批次发送2个，一共发送3个批次
	for j := 1; j <= 3; j++ {
		fpUsers := make([]*services.UserInfo, 0)
		for k := 1; k <= 2; k++ {
			id := j*2 + k
			user := &services.UserInfo{UserId: int32(id)}
			fpUsers = append(fpUsers, user)
		}
		time.Sleep(1 * time.Second)
		err := cliStream.Send(&services.UserRequest{Users: fpUsers})
		if err != nil {
			log.Fatalln(err)
		}
	}
	//发送结束后，一次性获取服务端响应
	userRes, err = cliStream.CloseAndRecv()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("客户端流模式:%v\n", userRes.Users)

	//双向流模式，客户端分批发送，服务端分批响应
	var wg sync.WaitGroup
	twStream, err := userClient.GetUserScoreByTwStream(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	//每批次发送2个，一共发送3个批次
	for j := 1; j <= 3; j++ {
		fpUsers := make([]*services.UserInfo, 0)
		for k := 1; k <= 2; k++ {
			id := j*2 + k
			user := &services.UserInfo{UserId: int32(id)}
			fpUsers = append(fpUsers, user)
		}
		err := twStream.Send(&services.UserRequest{Users: fpUsers})
		if err != nil {
			log.Fatalln(err)
		}
		//新建goroutines获取服务端响应并处理
		wg.Add(1)
		go func() {
			defer wg.Done()
			res, err := twStream.Recv()
			if err == io.EOF {
				return
			}
			if err != nil {
				log.Fatalln(err)
			}
			//类似业务处理逻辑
			fmt.Printf("双向流模式:%v\n", res.Users)
		}()
	}
	wg.Wait()
}
