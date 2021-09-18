package helper

import (
	"crypto/tls"
	"crypto/x509"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
)

//GetClientCreds 获取客户端相关证书配置
func GetClientCreds() credentials.TransportCredentials {
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
	return creds
}
