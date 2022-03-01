package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"time"

	"greet/greetpb"

	"google.golang.org/grpc"
)

type server struct{}

func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	log.Printf("Greet function was invoked with %v", req)
	first_name := req.GetGreeting().GetFirstName()
	result := "Hello " + first_name
	res := &greetpb.GreetResponse{
		Result: result,
	}

	return res, nil
}

func (*server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
	result := ""
	for {
		req, err := stream.Recv()

		if err == io.EOF {
			return stream.SendAndClose(&greetpb.LongGreetResponse{
				Result: result,
			})
		}

		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
		}

		result = result + "Hello " + req.GetGreeting().GetFirstName() + "!\n"
	}
}

func (*server) GreetEveryone(stream greetpb.GreetService_GreetEveryoneServer) error {
  
  for {
    req, err := stream.Recv()
    
    if err == io.EOF {
      return nil
    }
    if err != nil {
      log.Fatalf("Error receiving client stream: %v",err)
      return err
    }

    log.Printf("Received: %v",req.GetGreeting().GetFirstName())
    err1 := stream.Send(&greetpb.GreetEveryoneResponse{
      Result: "Hello " + req.GetGreeting().GetFirstName() + "!",
    })
    
    if err1 != nil {

      log.Fatalf("Error sending response: %v",err1)
      return err1
    }
  }
}

func (*server) Sum(ctx context.Context, req *greetpb.SumRequest) (*greetpb.SumResponse, error) {
	log.Printf("Sum function was invoked with %v", req)

	op_1 := req.GetOperand_1()
	op_2 := req.GetOperand_2()

	res := &greetpb.SumResponse{
		Result: op_1 + op_2,
	}

	return res, nil
}

func (*server) FindMax(stream greetpb.GreetService_FindMaxServer) error {
  max := int64(0)

  for {
    req, err := stream.Recv()

    if err == io.EOF {
      return nil
    }

    if err != nil {
      log.Fatalf("Error reading stream: %v",err)
      return err
    }

    num := req.GetNumber()

    if num > max {
      max = num

      err := stream.Send(&greetpb.FindMaxResponse{
        Maximum: max,
      })

      if err != nil {
        log.Fatalf("Error sending response: %v", err)
        return err
      }
    }
  }
}

func (*server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	firstName := req.GetGreeting().GetFirstName()
	for i := 0; i < 10; i++ {
		res := &greetpb.GreetManyTimesResponse{
			Result: "Hello " + firstName + " number " + strconv.Itoa(i),
		}
		stream.Send(res)
		time.Sleep(1000 * time.Millisecond)
	}
	return nil
}

func (*server) PrimeNumberDecomposition(req *greetpb.PrimeNumberDecompositionRequest, stream greetpb.GreetService_PrimeNumberDecompositionServer) error {
	N := req.GetNumber()

	k := int64(2)

	for N > 1 {
		if N%k == 0 {
			stream.Send(&greetpb.PrimeNumberDecompositionResponse{
				Result: k,
			})
			N = N / k
		} else {
			k++
		}
	}
	//decompose(stream, N, 2)
	return nil
}

func (*server) Average(stream greetpb.GreetService_AverageServer) error {
  sum := int64(0)
  count := int64(0)

  for {
    req, err := stream.Recv()
    if err == io.EOF {
      if count == 0 {
        return stream.SendAndClose(&greetpb.AverageResponse{
          Average: 0,
        })
      }
      return stream.SendAndClose(&greetpb.AverageResponse{
        Average: float32(sum) / float32(count),
      })
    }
    if err != nil {
      log.Fatalf("Error reading request from stream: %v",err)
    }

    sum += req.GetNumber()
    count ++;
  }

}

func decompose(stream greetpb.GreetService_PrimeNumberDecompositionServer, N int64, k int64) {
	if N > 1 {
		if N%k == 0 {
			stream.Send(&greetpb.PrimeNumberDecompositionResponse{
				Result: k,
			})
			decompose(stream, N/k, k)
		} else {
			decompose(stream, N, k+1)
		}
	}
}

func main() {
	fmt.Println("Hello world")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
