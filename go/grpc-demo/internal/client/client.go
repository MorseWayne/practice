package client

import (
	"context"
	"io"
	"log"
	"time"

	v1 "github.com/MorseWayne/grpc-demo/api/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Run(addr string) error {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer conn.Close()

	c1 := v1.NewCalculatorServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// unray
	resp, err := c1.Add(ctx, &v1.AddRequest{A: 3, B: 4})
	if err != nil {
		log.Println("Add error:", err)
	} else {
		log.Println("Add result:", resp.Result)
	}

	// client streaming
	cs, _ := c1.SumStream(ctx)
	cs.Send(&v1.AddRequest{A: 1, B: 2})
	cs.Send(&v1.AddRequest{A: 3, B: 4})
	sumResp, err := cs.CloseAndRecv()
	log.Println("SumStream:", sumResp.GetResult(), err)

	// server streaming
	ss, _ := c1.RangeAdd(ctx, &v1.RangeRequest{Start: 1, End: 3})
	for {
		m, err := ss.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println("RangeAdd error: ", err)
			break
		}
		log.Println("RangeAdd recv: ", m.Result)
	}

	// bidrectional
	bs, _ := c1.ChatAdd(ctx)
	bs.Send(&v1.AddRequest{A: 10, B: 5})
	bs.Send(&v1.AddRequest{A: 10, B: 6})
	bs.CloseSend()
	for {
		m, err := bs.Recv()
		if err != io.EOF {
			break
		}
		if err != nil {
			log.Println("Chatadd error: ", err)
			break
		}
		log.Println("chatadd recv: ", m.Result)
	}

	return nil
}
