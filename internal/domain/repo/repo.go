package repo

import (
	"context"

	"github.com/laoggu/cabbage-vehicle-backend/internal/domain/entity"
)

type VehicleRepo interface {
	Create(ctx context.Context, v *entity.Vehicle) error
	Update(ctx context.Context, v *entity.Vehicle) error
	Get(ctx context.Context, id string) (*entity.Vehicle, error)
	List(ctx context.Context, tenantID string, offset, limit int) ([]*entity.Vehicle, error)
}
