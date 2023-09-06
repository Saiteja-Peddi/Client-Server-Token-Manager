package main

import (
	"context"
	"flag"
	"log"
	"time"

	pb "example.com/AOS_PRJ2/tokenmanager"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

//Flags that can be used while running.
var (
	create = flag.Bool("create", false, "Calls Create method if true")
	write  = flag.Bool("write", false, "Calls Write method if true")
	read   = flag.Bool("read", false, "Calls Read method if true")
	drop   = flag.Bool("drop", false, "Calls Drop method if true")
	id     = flag.String("id", "", "Defines what id to use")
	name   = flag.String("name", "", "Defines what name to use")
	host   = flag.String("host", "localhost", "Defines hosting domain")
	port   = flag.String("port", "4200", "Defines port")
	low    = flag.Uint64("low", 0, "Defines low")
	mid    = flag.Uint64("mid", 10, "Defines mid")
	high   = flag.Uint64("high", 20, "Defines high")
)

func main() {
	flag.Parse()
	var address = *host + ":" + *port
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewTokenManagerClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	//Checks for what remote procedure to call from flags.
	if *create {
		res, err := c.Create(ctx, &pb.NormalRequest{Id: *id})
		if err != nil {
			log.Fatalf("Unable to create: %v", err)
		}
		if res.GetMessage() == 1 {
			log.Printf("Successfully created specified token")
		}
	} else if *read {
		res, err := c.Read(ctx, &pb.NormalRequest{Id: *id})
		if err != nil {
			log.Fatalf("Unable to read: %v", err)
		}
		log.Println("Final Value: ", res.GetMessage())

	} else if *write {
		res, err := c.Write(ctx, &pb.WriteRequest{Id: *id, Name: *name, Low: *low, Mid: *mid, High: *high})
		if err != nil {
			log.Fatalf("Unable to write: %v", err)
		}
		log.Println("Partial Value: ", res.GetMessage())

	} else if *drop {
		res, err := c.Drop(ctx, &pb.NormalRequest{Id: *id})
		if err != nil {
			log.Fatalf("Unable to drop: %v", err)
		}
		if res.GetMessage() == 1 {
			log.Printf("Successfully droped specified token")
		}
	}

}
