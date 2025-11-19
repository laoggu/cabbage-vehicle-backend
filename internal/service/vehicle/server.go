package vehicle

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/laoggu/cabbage-vehicle-backend/api"
	"github.com/laoggu/cabbage-vehicle-backend/internal/domain/entity"
	"github.com/laoggu/cabbage-vehicle-backend/internal/domain/repo"
	"github.com/laoggu/cabbage-vehicle-backend/internal/infra/oss"
	"github.com/laoggu/cabbage-vehicle-backend/internal/pkg/tenant"
)

type Server struct {
	api.UnimplementedVehicleServiceServer
	log  *zap.Logger
	repo repo.VehicleRepo
	oss  *oss.Client
}

func NewServer(log *zap.Logger, repo repo.VehicleRepo, oss *oss.Client) *Server {
	return &Server{log: log, repo: repo, oss: oss}
}

func (s *Server) CreateVehicle(ctx context.Context, req *api.CreateVehicleReq) (*api.Vehicle, error) {
	v := &entity.Vehicle{
		ID:        uuid.NewString(),
		Plate:     req.Plate,
		Model:     req.Model,
		ColdChain: req.ColdChain,
		IceBoxNo:  req.IceBoxNo,
	}
	if err := s.repo.Create(ctx, v); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return toPb(v), nil
}

func (s *Server) UpdateVehicle(ctx context.Context, req *api.UpdateVehicleReq) (*api.Vehicle, error) {
	v, err := s.repo.Get(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, "vehicle not found")
	}
	v.Plate = req.Plate
	v.Model = req.Model
	v.ColdChain = req.ColdChain
	v.IceBoxNo = req.IceBoxNo
	if err := s.repo.Update(ctx, v); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return toPb(v), nil
}

func (s *Server) GetVehicle(ctx context.Context, req *api.GetVehicleReq) (*api.Vehicle, error) {
	v, err := s.repo.Get(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, "vehicle not found")
	}
	return toPb(v), nil
}

func (s *Server) ListVehicles(ctx context.Context, req *api.ListVehiclesReq) (*api.ListVehiclesResp, error) {
	items, err := s.repo.List(ctx, tenant.ID(ctx), int(req.Offset), int(req.Limit))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	resp := &api.ListVehiclesResp{Items: make([]*api.Vehicle, len(items))}
	for i, v := range items {
		resp.Items[i] = toPb(v)
	}
	return resp, nil
}

func (s *Server) GetUploadURL(ctx context.Context, req *api.GetUploadURLReq) (*api.GetUploadURLResp, error) {
	url, obj, err := s.oss.PresignPut(req.Suffix)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &api.GetUploadURLResp{Url: url, Object: obj}, nil
}

func toPb(v *entity.Vehicle) *api.Vehicle {
	return &api.Vehicle{
		Id:        v.ID,
		Plate:     v.Plate,
		Model:     v.Model,
		ColdChain: v.ColdChain,
		IceBoxNo:  v.IceBoxNo,
	}
}
