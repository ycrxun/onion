package server

import (
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"
	"github.com/ycrxun/onion/services/account/storage"
	pb "github.com/ycrxun/onion/services/account/proto"
	"google.golang.org/grpc"
	"fmt"
	"net"
	"log"
	"golang.org/x/net/context"
	"github.com/ycrxun/onion/metrics"
	"time"
)

type Server struct {
	tracer  opentracing.Tracer
	storage storage.Storage
	m       metrics.Metrics
}

// NewServer returns a new server
func NewServer(tr opentracing.Tracer, metrics metrics.Metrics,
	storage storage.Storage) *Server {
	return &Server{
		tracer:  tr,
		storage: storage,
		m:       metrics,
	}
}

// Run starts the server
func (s *Server) Run(port int) error {
	srv := grpc.NewServer(
		grpc.UnaryInterceptor(
			otgrpc.OpenTracingServerInterceptor(s.tracer),
		),
	)
	pb.RegisterAccountServiceServer(srv, s)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	return srv.Serve(lis)
}

func (s *Server) List(ctx context.Context, req *pb.ListAccountsRequest) (*pb.ListAccountsResponse, error) {
	defer func(begin time.Time) {
		s.m.Counter.With("method", "list").Add(1)
		s.m.Histogram.With("method", "list").Observe(time.Since(begin).Seconds())
	}(time.Now())

	list, next, err := s.storage.List(req.PageSize, req.PageToken)
	if err != nil {
		return nil, err
	}
	var accounts []*pb.Account
	for _, v := range list {
		account := pb.Account{
			Id:                 v.ID,
			Name:               v.Name,
			Email:              v.Email,
			ConfirmToken:       v.ConfirmationToken,
			PasswordResetToken: v.PasswordResetToken,
			Metadata:           v.Metadata,
		}
		accounts = append(accounts, &account)
	}
	response := &pb.ListAccountsResponse{
		Accounts:      accounts,
		NextPageToken: next,
	}
	return response, nil
}
func (s *Server) GetById(ctx context.Context,
	req *pb.GetByIdRequest) (*pb.Account,
	error) {
	a, err := s.storage.ReadByID(req.Id)

	return &pb.Account{
		Id:    a.ID,
		Email: a.Email,
		Name:  a.Name,
	}, err
}
func (s *Server) GetByEmail(context.Context, *pb.GetByEmailRequest) (*pb.Account, error) {
	return nil, nil
}
func (s *Server) AuthenticateByEmail(context.Context, *pb.AuthenticateByEmailRequest) (*pb.Account, error) {
	return nil, nil
}
func (s *Server) GeneratePasswordToken(context.Context, *pb.GeneratePasswordTokenRequest) (*pb.GeneratePasswordTokenResponse, error) {
	return nil, nil
}
func (s *Server) ResetPassword(context.Context, *pb.ResetPasswordRequest) (*pb.Account, error) {
	return nil, nil
}
func (s *Server) ConfirmAccount(context.Context, *pb.ConfirmAccountRequest) (*pb.Account, error) {
	return nil, nil
}
func (s *Server) Create(context.Context, *pb.CreateAccountRequest) (*pb.Account, error) {
	return nil, nil
}
func (s *Server) Update(context.Context, *pb.UpdateAccountRequest) (*pb.Account, error) {
	return nil, nil
}
func (s *Server) Delete(context.Context, *pb.DeleteAccountRequest) (*pb.DeleteAccountResponse, error) {
	return nil, nil
}
