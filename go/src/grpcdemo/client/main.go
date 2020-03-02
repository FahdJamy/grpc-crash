package main

import (
	"flag"
	"fmt"
	"grpcdemo/src/grpcdemo/pb"
	"io"
	"log"
	"os"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

// PORT defines the port on which the server is running on
const PORT = "9002"

func main() {
	option := flag.Int("o", 1, "Command to run")
	flag.Parse()

	// get credentials from cert.pem and key.pem
	creds, err := credentials.NewClientTLSFromFile("cert.pem", "")
	if err != nil {
		log.Fatal(err)
	}
	opts := []grpc.DialOption{grpc.WithTransportCredentials(creds)}
	conn, err := grpc.Dial("localhost"+PORT, opts...)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewEmployeeServiceClient(conn)

	switch *option {
	case 1:
		SendMetadata(client)
	case 2:
		GetByBadgeNumber(client)
	case 3:
		GetAll(client)
	case 4:
		AddPhoto(client)
	case 5:
		SaveAll(client)
	}
}

// GetByBadgeNumber gets employee by badge number
func GetByBadgeNumber(client pb.EmployeeServiceClient) {
	res, err := client.GetByBadgeNumber(context.Background(), &pb.GetByBadgeNumberRequest{BadgeNumber: 28})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.Employee)
}

// SaveAll saves a number of employees
func SaveAll(client pb.EmployeeServiceClient) {
	employees := []pb.Employee{
		pb.Employee{
			Id:                  10,
			BadgeNumber:         443,
			FirstName:           "John",
			LastName:            443,
			VacationAccrualRate: 2,
			VacationAccrued:     10,
		},
		pb.Employee{
			Id:                  12,
			BadgeNumber:         410,
			FirstName:           "Smith",
			LastName:            410,
			VacationAccrualRate: 4,
			VacationAccrued:     20,
		},
	}

	stream, err := client.SaveAll(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	doneChan := make(chan struct{})
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				doneChan <- struct{}{}
				break
			}
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(res.Employee)
		}
	}()

	for _, emp := range employees {
		err := stream.Send(&pb.EmployeeRequest{Employee: &emp})
		if err != nil {
			log.Fatal(err)
		}
	}
	stream.CloseSend()
	<-doneChan
}

// AddPhoto creates a photo for the client
func AddPhoto(client pb.EmployeeServiceClient) {
	f, err := os.Open("Photo.jpg")
	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}
	md := metadata.New(map[string]string{"badgenumber": "2080"})
	ctx := context.Background()
	ctx = metadata.NewOutgoingContext(ctx, md)

	stream, err := client.AddPhoto(ctx)
	if err != nil {
		log.Fatal(err)
	}
	for {
		chunk := make([]byte, 64*1024)
		n, err := f.Read(chunk)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		if n < len(chunk) {
			chunk = chunk[:n]
		}
		stream.Send(&pb.AddPhotoRequest{Data: chunk})
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.IsOk)
}

// GetAll gets all employees
func GetAll(client pb.EmployeeServiceClient) {
	stream, err := client.GetAll(context.Background(), &pb.GetAllRequest{})
	if err != nil {
		log.Fatal(err)
	}
	for {
		res, err := stream.Recv()
		if err != io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(res.Employee)
	}
}

// SendMetadata sends meta data to the server
func SendMetadata(client pb.EmployeeServiceClient) {
	md := metadata.MD{}
	md["user"] = []string{"invisible"}
	md["password"] = []string{"krs1krs1"}

	ctx := context.Background()
	ctx = metadata.NewOutgoingContext(ctx, md)

	client.GetByBadgeNumber(ctx, &pb.GetByBadgeNumberRequest{})
}
