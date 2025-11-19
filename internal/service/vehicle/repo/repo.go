package repo

import (
	"context"

	"github.com/laoggu/cabbage-vehicle-backend/internal/domain/entity"
	"github.com/laoggu/cabbage-vehicle-backend/internal/domain/repo"
	"gorm.io/gorm"
)

type mysqlRepo struct {
	db *gorm.DB
}

func NewVehicleRepo(db *gorm.DB) repo.VehicleRepo {
	return &mysqlRepo{db: db}
}

func (r *mysqlRepo) Create(ctx context.Context, v *entity.Vehicle) error {
	return r.db.WithContext(ctx).Create(v).Error
}

func (r *mysqlRepo) Update(ctx context.Context, v *entity.Vehicle) error {
	return r.db.WithContext(ctx).Save(v).Error
}

func (r *mysqlRepo) Get(ctx context.Context, id string) (*entity.Vehicle, error) {
	var v entity.Vehicle
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&v).Error
	return &v, err
}

func (r *mysqlRepo) List(ctx context.Context, tenantID string, offset, limit int) ([]*entity.Vehicle, error) {
	var list []*entity.Vehicle
	err := r.db.WithContext(ctx).
		Where("tenant_id = ?", tenantID).
		Offset(offset).Limit(limit).
		Find(&list).Error
	return list, err
}
