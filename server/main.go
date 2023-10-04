package main

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Philipp15b/go-steam/v3"
	pb "github.com/twonull/grpc-inspect/proto"
	"google.golang.org/grpc"
)

type Job struct {
	S            uint64
	A            uint64
	D            uint64
	M            uint64
	ResponseChan chan *pb.ItemInfo
}

type Worker struct {
	Queue chan *Job
	Res   chan *pb.ItemInfo
	Ready bool
}

type InspectorServer struct {
	pb.UnimplementedInspectServer
	Workers []*Worker
}

var file *string = flag.String("file", "", "Specifies the location of accounts.txt containing Steam accounts in user:pass format. (Required)")

func (s *InspectorServer) SendInspect(ctx context.Context, request *pb.InspectRequest) (*pb.InspectResponse, error) {
	if len(request.Fields) != 4 {
		return nil, errors.New("malformed request")
	}

	worker := s.findWorkerWithShortestQueue()

	if worker == nil {
		return nil, errors.New("no available worker")
	}

	resChan := make(chan *pb.ItemInfo)

	worker.Queue <- &Job{
		S:            request.Fields[0],
		A:            request.Fields[1],
		D:            request.Fields[2],
		M:            request.Fields[3],
		ResponseChan: resChan,
	}

	answer, ok := <-resChan

	if !ok {
		for unfinished := range worker.Queue {
			newWorker := s.findWorkerWithShortestQueue()
			newWorker.Queue <- unfinished
		}
		return nil, errors.New("worker encountered fatal error")
	}

	return &pb.InspectResponse{
		ItemInfo: answer,
	}, nil
}

func (s *InspectorServer) findWorkerWithShortestQueue() *Worker {
	var minQueueLen = -1
	var selectedWorker *Worker

	for _, worker := range s.Workers {
		if worker.Ready {
			queueLen := len(worker.Queue)
			if minQueueLen == -1 || queueLen < minQueueLen {
				minQueueLen = queueLen
				selectedWorker = worker
			}
		}
	}
	return selectedWorker
}

func main() {
	flag.Parse()
	conv := filepath.ToSlash(*file)
	_, err := os.Stat(conv)
	if err != nil {
		log.Fatal("Error finding accounts.txt. Please create this file according to the Readme or correct the path.")
	}
	words, err := readLines(conv)
	if err != nil {
		log.Fatal("accounts.txt formatted incorrectly.")
	}
	service := &InspectorServer{}
	service.Workers = make([]*Worker, len(words))

	for i := 0; i < len(words); i++ {
		service.Workers[i] = &Worker{
			Queue: make(chan *Job),
			Res:   make(chan *pb.ItemInfo),
			Ready: false,
		}
		go service.startWorker(service.Workers[i], words[i][0], words[i][1])
	}

	listenAddr := "localhost:50051"
	lis, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterInspectServer(s, service)

	log.Printf("Server listening on %s\n", listenAddr)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func (s *InspectorServer) startWorker(worker *Worker, user string, pass string) {
	loginInfo := new(steam.LogOnDetails)
	loginInfo.Username = user
	loginInfo.Password = pass

	client := steam.NewClient()
	cs := InitClient(client)
	client.Connect()

	go func() {
		for event := range client.Events() {
			switch e := event.(type) {
			case *steam.ConnectedEvent:
				client.Auth.LogOn(loginInfo)
			case *steam.LoggedOnEvent:
				log.Println(user, "Login Success")
				cs.SendHello()
				cs.SetPlaying(true)
				worker.Ready = true
			case *pb.ItemInfo:
				worker.Res <- e
			case steam.FatalErrorEvent:
				log.Println(user, "Encountered Fatal Error", event)
				close(worker.Queue)
				close(worker.Res)
				return
			}
		}
	}()

	for {
		if worker.Ready {
			job, ok := <-worker.Queue
			if !ok {
				close(job.ResponseChan)
				worker.Ready = false
				return
			}
			cs.InspectItem(job.S, job.A, job.D, job.M)
			res, ok := <-worker.Res
			if !ok {
				close(job.ResponseChan)
				worker.Ready = false
				return
			}
			job.ResponseChan <- res
		}
		time.Sleep(1000 * time.Millisecond)
	}
}

func readLines(path string) ([][]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var accts [][]string
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		word := scanner.Text()
		combo := strings.Split(word, ":")
		if len(combo) == 2 {
			accts = append(accts, combo)
		}
	}
	return accts, scanner.Err()
}
