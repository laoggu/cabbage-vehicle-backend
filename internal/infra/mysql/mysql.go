package mysql

import (
	"github.com/laoggu/cabbage-vehicle-backend/internal/pkg/tenant"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func New(dsn string, log *zap.Logger) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}
	// 多租户插件：自动注入 tenant_id
	if err := db.Use(&tenant.Plugin{}); err != nil {
		return nil, err
	}
	return db, nil
}
