package main

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"net"
	"sync"

	pb "example.com/AOS_PRJ2/tokenmanager"
	"google.golang.org/grpc"
)

type Domain struct {
	Low  uint64
	Mid  uint64
	High uint64
}

type State struct {
	Partial uint64
	Final   uint64
}

type Token struct {
	Id     string
	mut    sync.Mutex
	Name   string
	domain *Domain
	state  *State
}

//Tokens gets stored in this map
var Tokens = make(map[string]*Token)

type tokenService struct {
	pb.UnimplementedTokenManagerServer
}

// Hash concatentates a message and a nonce and generates a hash value.
func Hash(name string, nonce uint64) uint64 {
	hasher := sha256.New()
	hasher.Write([]byte(fmt.Sprintf("%s %d", name, nonce)))
	return binary.BigEndian.Uint64(hasher.Sum(nil))
}

//Returns minimum two passed values
func minOf(val1 uint64, val2 uint64) uint64 {
	if val1 < val2 {
		return val1
	} else {
		return val2
	}
}

//Return minimum hash value
func minHash(name string, min uint64, max uint64) uint64 {
	var minVal uint64 = Hash(name, min)

	for i := min + 1; i < max; i++ {
		minVal = minOf(minVal, Hash(name, i))
	}
	return minVal
}

//Prints all tokens
func printAllTokes() {
	var i = 1
	fmt.Println("\n----------------------------------------------------------------")
	for key, val := range Tokens {
		fmt.Print("\nToken ", i, ":- Id: ", key)
		if val != nil && &val.Name != nil && val.Name != "" {
			fmt.Print(" Name: ", val.Name)
		}
		if val != nil && val.state != nil && &val.state.Partial != nil {
			fmt.Print(" Partial: ", val.state.Partial)
		}
		if val != nil && val.state != nil && &val.state.Final != nil {
			fmt.Print(" Final: ", val.state.Final)
		}
		i += 1
	}
	fmt.Println("\n----------------------------------------------------------------")
}

//Prints token with given id
func printToken(id string) {
	log.Print("\n")
	fmt.Println("----------------------------------------------------------------")
	fmt.Println("Current Token:-")
	//fmt.Println("Name: ", Tokens[id].Name, " Partial Value: ", Tokens[id].state.Partial, " Final Value: ", Tokens[id].state.Final)
	fmt.Printf("Id: %s", id)
	if Tokens[id] != nil && &Tokens[id].Name != nil && Tokens[id].Name != "" {
		fmt.Printf(" Name: %s", Tokens[id].Name)
	}
	if Tokens[id] != nil && Tokens[id].state != nil && &Tokens[id].state.Partial != nil {
		fmt.Printf(" Partial: %d", Tokens[id].state.Partial)
	}
	if Tokens[id] != nil && Tokens[id].state != nil && &Tokens[id].state.Final != nil {
		fmt.Printf(" Final: %d", Tokens[id].state.Final)
	}
	fmt.Println("\n----------------------------------------------------------------")

}

//Create rpc implementation
func (t *tokenService) Create(_ context.Context, req *pb.NormalRequest) (*pb.ServerResponse, error) {

	var createToken = &Token{
		Id: req.GetId(),
	}
	Tokens[req.GetId()] = createToken
	printToken(req.GetId())
	printAllTokes()
	return &pb.ServerResponse{
		Message: 1,
	}, nil
}

//Read rpc implementation
func (t *tokenService) Read(_ context.Context, req *pb.NormalRequest) (*pb.ServerResponse, error) {

	var tempFinalVal uint64 = minOf(minHash(Tokens[req.GetId()].Name, Tokens[req.GetId()].domain.Mid, Tokens[req.GetId()].domain.High), Tokens[req.GetId()].state.Partial)

	Tokens[req.GetId()].state.Final = tempFinalVal
	printToken(req.GetId())
	printAllTokes()
	return &pb.ServerResponse{
		Message: Tokens[req.GetId()].state.Final,
	}, nil
}

//Go routine that writes into token
func writeIntoToken(mut *sync.Mutex, req *pb.WriteRequest, wg *sync.WaitGroup) {

	mut.Lock()
	Tokens[req.GetId()] = &Token{
		Name: req.GetName(),
		domain: &Domain{
			Low:  req.GetLow(),
			Mid:  req.GetMid(),
			High: req.GetHigh(),
		},
		state: &State{
			Partial: minHash(req.GetName(), req.GetLow(), req.GetMid()),
		},
	}
	mut.Unlock()

	wg.Done()

}

//Write rpc implementation
func (t *tokenService) Write(_ context.Context, req *pb.WriteRequest) (*pb.ServerResponse, error) {
	var wg sync.WaitGroup
	var mutex sync.Mutex

	wg.Add(1)
	go writeIntoToken(&mutex, req, &wg)

	wg.Wait()
	printToken(req.GetId())
	printAllTokes()

	return &pb.ServerResponse{
		Message: Tokens[req.GetId()].state.Partial,
	}, nil

}

//Drop rpc implmentation
func (t *tokenService) Drop(_ context.Context, req *pb.NormalRequest) (*pb.ServerResponse, error) {

	delete(Tokens, req.GetId())
	log.Print("\n")
	printAllTokes()
	return &pb.ServerResponse{
		Message: 1,
	}, nil
}

func main() {
	//Flag that specifies from which port to listen for incomming requests
	var port = flag.Int("port", 4200, "Defines from which port to listen")
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	rpcServer := grpc.NewServer()
	tokenServer := &tokenService{}

	pb.RegisterTokenManagerServer(rpcServer, tokenServer)
	log.Printf("Server is listening at %v", lis.Addr())

	if err := rpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// type TokenManagerServer interface {
// 	Create(context.Context, *NormalRequest) (*ServerResponse, error)
// 	Write(context.Context, *WriteRequest) (*ServerResponse, error)
// 	Read(context.Context, *NormalRequest) (*ServerResponse, error)
// 	Drop(context.Context, *NormalRequest) (*ServerResponse, error)
// 	mustEmbedUnimplementedTokenManagerServer()
// }
