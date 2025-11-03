package pg

import (
	"context"

	"gorm.io/gorm"

	"github.com/chaitin/panda-wiki/consts"
	"github.com/chaitin/panda-wiki/domain"
	"github.com/chaitin/panda-wiki/log"
	"github.com/chaitin/panda-wiki/store/pg"
)

type ContributeRepo struct {
	db     *pg.DB
	logger *log.Logger
}

func NewContributeRepo(db *pg.DB, logger *log.Logger) *ContributeRepo {
	return &ContributeRepo{
		db:     db,
		logger: logger,
	}
}

func (r *ContributeRepo) Create(ctx context.Context, contribute *domain.Contribute) error {
	return r.db.WithContext(ctx).Create(contribute).Error
}

func (r *ContributeRepo) GetListByKBID(ctx context.Context, kbID string, page, perPage int) ([]*domain.Contribute, int64, error) {
	var contributes []*domain.Contribute
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.Contribute{}).Where("contributes.kb_id = ?", kbID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage
	// 关联查询 nodes 表获取文档标题，关联 auths 表获取用户信息
	if err := query.
		Joins("left join nodes on contributes.node_id = nodes.id").
		Joins("left join auths on contributes.auth_id = auths.id").
		Select("contributes.*, nodes.name as node_name, auths.user_info").
		Order("contributes.created_at DESC").
		Offset(offset).
		Limit(perPage).
		Scan(&contributes).Error; err != nil {
		return nil, 0, err
	}

	return contributes, total, nil
}

func (r *ContributeRepo) UpdateStatus(ctx context.Context, id string, status consts.ContributeStatus, auditUserID, reason string) error {
	updates := map[string]interface{}{
		"status":        status,
		"audit_user_id": auditUserID,
		"reason":        reason,
	}
	return r.db.WithContext(ctx).Model(&domain.Contribute{}).Where("id = ?", id).Updates(updates).Error
}

func (r *ContributeRepo) GetByID(ctx context.Context, id string) (*domain.Contribute, error) {
	var contribute domain.Contribute
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&contribute).Error; err != nil {
		return nil, err
	}
	return &contribute, nil
}

func (r *ContributeRepo) Delete(ctx context.Context, id, kbID string) error {
	result := r.db.WithContext(ctx).Where("id = ? AND kb_id = ?", id, kbID).Delete(&domain.Contribute{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
