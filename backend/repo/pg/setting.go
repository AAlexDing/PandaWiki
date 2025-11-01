package pg

import (
	"context"
	"errors"

	"github.com/chaitin/panda-wiki/domain"
	"github.com/chaitin/panda-wiki/log"
	"github.com/chaitin/panda-wiki/store/pg"
	"gorm.io/gorm"
)

type SettingRepository struct {
	db     *pg.DB
	logger *log.Logger
}

func NewSettingRepository(db *pg.DB, logger *log.Logger) *SettingRepository {
	return &SettingRepository{
		db:     db,
		logger: logger,
	}
}

func (r *SettingRepository) CreateSetting(ctx context.Context, setting *domain.Setting) error {
	return r.db.WithContext(ctx).Table("settings").Create(setting).Error
}

func (r *SettingRepository) GetSetting(ctx context.Context, kbID, key string) (*domain.Setting, error) {
	var setting domain.Setting
	err := r.db.WithContext(ctx).Table("settings").
		Where("kb_id = ? AND key = ?", kbID, key).
		First(&setting).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &setting, nil
}

func (r *SettingRepository) UpdateSetting(ctx context.Context, kbID, key string, value []byte) error {
	// 先尝试获取现有设置
	setting, err := r.GetSetting(ctx, kbID, key)
	if err != nil {
		return err
	}

	// 如果存在则更新，否则创建
	if setting != nil {
		return r.db.WithContext(ctx).Table("settings").
			Where("kb_id = ? AND key = ?", kbID, key).
			Update("value", value).Error
	}

	// 创建新记录
	newSetting := &domain.Setting{
		KBID:  kbID,
		Key:   key,
		Value: value,
	}
	return r.CreateSetting(ctx, newSetting)
}
