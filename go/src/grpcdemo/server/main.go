package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"

	"grpcdemo/src/grpcdemo/pb"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const port = ":9002"

var employees = []pb.Employee{
	pb.Employee{
		Id:                  1,
		BadgeNumber:         2080,
		FirstName:           "Grace",
		LastName:            2080,
		VacationAccrualRate: 2,
		VacationAccrued:     30,
	},
	pb.Employee{
		Id:                  2,
		BadgeNumber:         7890,
		FirstName:           "Rola",
		LastName:            2080,
		VacationAccrualRate: 2,
		VacationAccrued:     30,
	},
	pb.Employee{
		Id:                  3,
		BadgeNumber:         7320,
		FirstName:           "Mike",
		LastName:            2080,
		VacationAccrualRate: 2,
		VacationAccrued:     30,
	},
	pb.Employee{
		Id:                  4,
		BadgeNumber:         2080,
		FirstName:           "Baery",
		LastName:            2080,
		VacationAccrualRate: 2,
		VacationAccrued:     30,
	},
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal(err)
	}
	// get credentials from cert.pem and key.pem
	creds, err := credentials.NewServerTLSFromFile("cert.pem", "key.pem")
	if err != nil {
		log.Fatal(err)
	}
	// pass the credentials to the server that is gonna be created
	opts := []grpc.ServerOption{grpc.Creds(creds)}
	// create the server
	s := grpc.NewServer(opts...)
	// register the service using the pb package
	pb.RegisterEmployeeServiceServer(s, &employeeService{})
	log.Println("Running server on port " + port)
	// start the server
	s.Serve(lis)
}

// employeeService can be embedded to have forward compatible implementations.
type employeeService struct{}

func (s *employeeService) GetByBadgeNumber(ctx context.Context, req *pb.GetByBadgeNumberRequest) (*pb.EmployeeResponse, error) {

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		fmt.Printf("Metadata received is : ===> %v\n", md)
	}

	for _, emp := range employees {
		if req.BadgeNumber == emp.BadgeNumber {
			return &pb.EmployeeResponse{Employee: &emp}, nil
		}
	}

	return nil, errors.New("Employee not found")
}

func (s *employeeService) GetAll(req *pb.GetAllRequest, stream pb.EmployeeService_GetAllServer) error {
	for _, emp := range employees {
		stream.Send(&pb.EmployeeResponse{Employee: &emp})
	}
	return nil
}

func (s *employeeService) Save(ctx context.Context, req *pb.EmployeeRequest) (*pb.EmployeeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Save not implemented")
}

func (s *employeeService) SaveAll(stream pb.EmployeeService_SaveAllServer) error {
	for {
		emp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		employees = append(employees, *emp.Employee)
		stream.Send(&pb.EmployeeResponse{Employee: emp.Employee})
	}
	for _, e := range employees {
		fmt.Println(e)
	}

	return nil
}

func (s *employeeService) AddPhoto(stream pb.EmployeeService_AddPhotoServer) error {

	if md, ok := metadata.FromIncomingContext(stream.Context()); ok {
		fmt.Println("Received badge number is ===> %v", md["badgenumber"][0])
	}
	imgData := []byte{} // hold the entire IMG data

	for {
		data, err := stream.Recv()
		// if done reading the file
		if err == io.EOF {
			fmt.Println("All image is received and it is of size ==> %v", len(imgData))
			// Use sendAndClose to let the client know that you've received every thing and done
			return stream.SendAndClose(&pb.AddPhotoResponse{IsOk: true})
		}
		if err != nil {
			return err
		}
		fmt.Println("Received bytes size ===> %v", len(data.Data))
		imgData = append(imgData, data.Data...)
	}

	return nil
}
