package main

import (
	"context"
	"fmt"
	"greet/greetpb"
	"io"
	"log"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Hello I'm a client")

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}

	defer conn.Close()

	c := greetpb.NewGreetServiceClient(conn)

	//doUnary(c)

	//doSum(c)

	//doStream(c)

  //doPrimeStream(c)

  //doLongGreet(c)

  //doAverage(c)

  //doGreetEveryone(c)

  doFindMax(c)
}

func doGreetEveryone(c greetpb.GreetServiceClient) {
  stream, err := c.GreetEveryone(context.Background())

  if err != nil {
    log.Fatalf("Error creating stream: %v",err)
    return
  }

  requests := []*greetpb.GreetEveryoneRequest{
    {
      Greeting: &greetpb.Greeting {
        FirstName: "Karim",
      },
    },
    {
      Greeting: &greetpb.Greeting {
        FirstName: "Semina",
      },
    },
    {
      Greeting: &greetpb.Greeting {
        FirstName: "Rimmy",
      },
    },
    {
      Greeting: &greetpb.Greeting {
        FirstName: "Insha",
      },
    },
    {
      Greeting: &greetpb.Greeting {
        FirstName: "Ahan",
      },
    },
  }
  waitc := make(chan struct{})
  go func() {
    for _, req := range requests {
      stream.Send(req)
      //time.Sleep(1000 * time.Millisecond)
    }
    stream.CloseSend()
  }()

  go func(){
    for {
      res, err := stream.Recv()
      if err == io.EOF {
        break
      }
      if err != nil {
        log.Fatalf("Error receiving response: %v",err)
        break
      }

      log.Printf("Received: %v",res.GetResult())
    }
    close(waitc)
  }()

  <-waitc
}

func doFindMax(c greetpb.GreetServiceClient) {
  stream, err := c.FindMax(context.Background())

  if err != nil {
    log.Fatalf("Error getting stream: %v",err)
    return
  }

  waitc := make (chan struct{})
  numbers := []int64 {1,5,3,6,2,20}
  go func() {
    for _, num := range numbers {
      err := stream.Send(&greetpb.FindMaxRequest{
        Number: num,
      })
      if err != nil {
        log.Fatalf("Error sending request: %v",err)
        break
      }
    }
    stream.CloseSend()
  }()

  go func(){
    for {
      req, err := stream.Recv()

      if err == io.EOF {
        break
      }
      if err != nil {
        log.Fatalf("Error receiving response: %v",err)
        break
      }
      log.Printf("Received: %v",req.GetMaximum())
    }

    close(waitc)
  }()

  <-waitc
}

func doAverage(c greetpb.GreetServiceClient) {
  stream, err := c.Average(context.Background())

  if err != nil {
    log.Fatalf("Error sending request: %v",err)
  }

  for i := 1; i < 5; i++ {
    stream.Send(&greetpb.AverageRequest{
      Number: int64(i),
    })
  }

  res, err := stream.CloseAndRecv()

  if err != nil {
    log.Fatalf("Error receiving response: %v",err)
  }

  log.Printf("Average is: %v", res.GetAverage())
}

func doLongGreet(c greetpb.GreetServiceClient) {
  stream, err := c.LongGreet(context.Background())
  
  if err != nil {
    log.Fatalf("Error sending request: %v",err)
  }

  stream.Send(&greetpb.LongGreetRequest{
    Greeting: &greetpb.Greeting {
      FirstName: "Karim",
      LastName: "Lalani",
    },
  })
  stream.Send(&greetpb.LongGreetRequest{
    Greeting: &greetpb.Greeting {
      FirstName: "Semina",
      LastName: "Lalani",
    },
  })
  stream.Send(&greetpb.LongGreetRequest{
    Greeting: &greetpb.Greeting {
      FirstName: "Rimmy",
      LastName: "Lalani",
    },
  })
  res, err := stream.CloseAndRecv()

  if err != nil {
    log.Fatalf("Error receiving response: %v",err)
  } 

  log.Printf("Received response: %v",res.GetResult())
}

func doPrimeStream(c greetpb.GreetServiceClient) {
  resStream, err := c.PrimeNumberDecomposition(context.Background(), &greetpb.PrimeNumberDecompositionRequest {
    Number: 120,
  })

  if err != nil {
    log.Fatalf("Error sending request: %v", err)
  }

  for {
    msg, err := resStream.Recv()

    if err == io.EOF {
      break
    }

    if err != nil {
      log.Fatalf("error reading stream: %v",err)
    }

    log.Printf("Received Prime Decomposition: %v", msg.GetResult())
  }
}

func doStream(c greetpb.GreetServiceClient) {

	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Karim",
			LastName:  "Lalani",
		},
	}

	resStream, err := c.GreetManyTimes(context.Background(), req)

	if err != nil {
		log.Fatalf("error sending request: %v", err)
	}

	for {
		msg, err := resStream.Recv()

		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error reading stream: %v", err)
		}

    log.Printf("Response from GreetManyTymes: %v", msg.GetResult())
	}
}

func doUnary(c greetpb.GreetServiceClient) {

	//fmt.Printf("created client: %f", c)

	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Karim",
			LastName:  "Lalani",
		},
	}
	res, err := c.Greet(context.Background(), req)

	if err != nil {
		log.Fatalf("Error while calling Greet RPC: %v", err)
	}

	log.Printf("Response from Greet: %v", res)
}

func doSum(c greetpb.GreetServiceClient) {
	res, err := c.Sum(context.Background(), &greetpb.SumRequest{
		Operand_1: 101,
		Operand_2: 99,
	})

	if err != nil {
		log.Fatalf("Error while calling Sum RPC: %v", err)
	}

	log.Printf("Response from Sum: %v", res)
}
