package pg

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/chaitin/panda-wiki/domain"
	"github.com/chaitin/panda-wiki/log"
	"github.com/chaitin/panda-wiki/store/pg"
	"gorm.io/gorm"
)

type PromptRepo struct {
	db     *pg.DB
	logger *log.Logger
}

type promptJson struct {
	Content string
}

func NewPromptRepo(db *pg.DB, logger *log.Logger) *PromptRepo {
	return &PromptRepo{
		db:     db,
		logger: logger,
	}
}

func (r *PromptRepo) GetPrompt(ctx context.Context, kbID string) (string, error) {
	var setting domain.Setting
	var prompt promptJson
	err := r.db.WithContext(ctx).Table("settings").
		Where("kb_id = ? AND key = ?", kbID, domain.SettingKeySystemPrompt).
		First(&setting).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil
		}
		return "", err
	}
	if err := json.Unmarshal(setting.Value, &prompt); err != nil {
		return "", err
	}
	return prompt.Content, nil
}

func (r *PromptRepo) UpdatePrompt(ctx context.Context, kbID, content string) error {
	prompt := promptJson{Content: content}
	value, err := json.Marshal(prompt)
	if err != nil {
		return err
	}

	// 先尝试获取现有设置
	var setting domain.Setting
	err = r.db.WithContext(ctx).Table("settings").
		Where("kb_id = ? AND key = ?", kbID, domain.SettingKeySystemPrompt).
		First(&setting).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 不存在则创建
			newSetting := &domain.Setting{
				KBID:  kbID,
				Key:   domain.SettingKeySystemPrompt,
				Value: value,
			}
			return r.db.WithContext(ctx).Table("settings").Create(newSetting).Error
		}
		return err
	}

	// 存在则更新
	return r.db.WithContext(ctx).Table("settings").
		Where("kb_id = ? AND key = ?", kbID, domain.SettingKeySystemPrompt).
		Update("value", value).Error
}
