package auth

import (
	"context"

	"github.com/laoggu/cabbage-vehicle-backend/api"
	"github.com/laoggu/cabbage-vehicle-backend/internal/pkg/jwt"
	"github.com/laoggu/cabbage-vehicle-backend/internal/service/auth/app"
	"go.uber.org/zap"
)

type Server struct {
	api.UnimplementedAuthServiceServer
	log *zap.Logger
}

func NewServer(log *zap.Logger) *Server {
	return &Server{log: log}
}

func (s *Server) WxLogin(ctx context.Context, req *api.WxLoginReq) (*api.WxLoginResp, error) {
	wx, err := app.Code2Session(ctx, req.Code)
	if err != nil {
		return nil, err
	}
	acc, ref, err := app.GenTokens(wx.OpenID)
	if err != nil {
		return nil, err
	}
	return &api.WxLoginResp{
		AccessToken:  acc,
		RefreshToken: ref,
		ExpiresIn:    1800,
	}, nil
}

func (s *Server) Refresh(ctx context.Context, req *api.RefreshReq) (*api.RefreshResp, error) {
	cl, err := jwt.Parse(req.RefreshToken)
	if err != nil {
		return nil, err
	}
	acc, _, err := app.GenTokens(cl.Sub)
	if err != nil {
		return nil, err
	}
	return &api.RefreshResp{AccessToken: acc, ExpiresIn: 1800}, nil
}
