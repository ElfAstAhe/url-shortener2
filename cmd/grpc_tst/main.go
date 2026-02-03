package main

import (
	"context"
	"fmt"

	pb "github.com/ElfAstAhe/url-shortener2/api/proto"
	"github.com/ElfAstAhe/url-shortener2/internal/bll/service/auth"
	appgrpc "github.com/ElfAstAhe/url-shortener2/internal/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func main() {
	userInfo := auth.BuildRandomUserInfo()
	tokenString, err := auth.NewJWTStringFromUserInfo(userInfo)
	if err != nil {
		panic(err)
	}

	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client := pb.NewShortenerServiceClient(conn)

	// key
	key := "d12e2b1b2197d3634eae3d911427c1f3"

	// context
	ctx := metadata.NewOutgoingContext(context.Background(), metadata.Pairs(appgrpc.MetaDataAuthorization, tokenString))

	// get url
	req := pb.URLExpandRequest_builder{
		Id: key,
	}.Build()

	fmt.Println(req)

	data, err := client.ExpandURL(ctx, req)
	if err != nil {
		panic(err)
	}

	fmt.Printf("get url result:\nstring \n[%s]\nvariable \n[%v]\n", data, data)
}
