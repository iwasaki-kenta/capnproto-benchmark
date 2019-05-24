package main

import (
	"context"
	"flag"
	"fmt"
	"math/rand"
	"net"
	"sync/atomic"
	"time"
	"zombiezen.com/go/capnproto2/rpc"
)

const PayloadSize = 600

var sendCount, recvCount uint64

type benchmarkServer struct{}

func (benchmarkServer) Send(call Benchmark_send) error {
	req, err := call.Params.Req()
	if err != nil {
		return err
	}

	bytes, err := req.Bytes()
	if err != nil {
		return err
	}

	if len(bytes) != PayloadSize {
		panic("Got an unexpected payload size.")
	}

	atomic.AddUint64(&recvCount, 1)

	var buf [PayloadSize]byte

	if _, err := rand.Read(buf[:]); err != nil {
		panic(err)
	}

	res, err := call.Results.NewRes()
	if err != nil {
		return err
	}

	if err = res.SetBytes(buf[:]); err != nil {
		return err
	}

	atomic.AddUint64(&sendCount, 1)

	return nil
}

func singleFlags() {
	panic("go run single.go [client/server] [benchmark endpoint address]")
}

func main() {
	flag.Parse()

	if len(flag.Args()) < 1 {
		singleFlags()
	}

	go func() {
		for range time.Tick(1 * time.Second) {
			fmt.Printf("Sent %d messages, and received %d messages.\n", atomic.SwapUint64(&sendCount, 0), atomic.SwapUint64(&recvCount, 0))
		}
	}()

	switch flag.Arg(0) {
	case "client":
		if len(flag.Args()) != 2 {
			singleFlags()
		}

		sock, err := net.Dial("tcp", flag.Arg(1))
		if err != nil {
			panic(err)
		}

		conn := rpc.NewConn(rpc.StreamTransport(sock))
		ctx := context.Background()

		benchmark := Benchmark{Client: conn.Bootstrap(ctx)}

		for {
			res, err := benchmark.Send(ctx, func(params Benchmark_send_Params) error {
				req, err := params.NewReq()
				if err != nil {
					return err
				}

				var buf [PayloadSize]byte

				if _, err := rand.Read(buf[:]); err != nil {
					panic(err)
				}

				if err = req.SetBytes(buf[:]); err != nil {
					return err
				}

				atomic.AddUint64(&sendCount, 1)

				return nil
			}).Res().Struct()

			if err != nil {
				panic(err)
			}

			bytes, err := res.Bytes()
			if err != nil {
				panic(err)
			}

			if len(bytes) != PayloadSize {
				panic("Got an unexpected payload size.")
			}

			atomic.AddUint64(&recvCount, 1)
		}
	case "server":
		if len(flag.Args()) != 1 {
			singleFlags()
		}

		listener, err := net.Listen("tcp", ":0")
		if err != nil {
			panic(err)
		}

		fmt.Printf("Listening on port %d.\n", listener.Addr().(*net.TCPAddr).Port)

		server := Benchmark_ServerToClient(benchmarkServer{})

		for {
			conn, err := listener.Accept()
			if err != nil {
				panic(err)
			}

			c := rpc.NewConn(rpc.StreamTransport(conn), rpc.MainInterface(server.Client))
			c.Wait()
		}
	default:
		singleFlags()
	}

}
