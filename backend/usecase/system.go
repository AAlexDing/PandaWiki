package usecase

import (
	"context"
	"fmt"

	v1 "github.com/chaitin/panda-wiki/api/system/v1"
	"github.com/chaitin/panda-wiki/log"
	"github.com/chaitin/panda-wiki/repo/pg"
	utilsDocker "github.com/chaitin/panda-wiki/utils"
)

type SystemUseCase struct {
	nodeRepo *pg.NodeRepository
	logger   *log.Logger
}

func NewSystemUseCase(nodeRepo *pg.NodeRepository, logger *log.Logger) *SystemUseCase {
	return &SystemUseCase{
		nodeRepo: nodeRepo,
		logger:   logger.WithModule("usecase.system"),
	}
}

// GetSystem 获取系统状态信息
func (u *SystemUseCase) GetSystem(ctx context.Context, kbID string) (*v1.SystemResp, error) {
	// 获取文档统计
	currentCount, newIn24h, learningSucceeded, learningFailed, err := u.nodeRepo.GetStatusDocumentStats(ctx, kbID)
	if err != nil {
		return nil, fmt.Errorf("failed to get document stats: %w", err)
	}

	// 获取学习状态统计
	basicPending, basicRunning, basicFailed, basicSucceeded, enhancePending, enhanceRunning, enhanceFailed, enhanceSucceeded, basicFailedDocs, enhanceFailedDocs, err := u.nodeRepo.GetStatusLearningStats(ctx, kbID)
	if err != nil {
		return nil, fmt.Errorf("failed to get learning stats: %w", err)
	}

	// 计算增强处理队列进度（不使用时间限制）
	enhanceTotal := enhancePending + enhanceRunning + enhanceFailed + enhanceSucceeded
	enhanceProgress := 0
	if enhanceTotal > 0 {
		enhanceProgress = int((enhanceTotal - enhancePending - enhanceRunning - enhanceFailed) * 100 / enhanceTotal)
	}

	// 计算基础处理队列进度（不使用时间限制）
	basicTotal := basicPending + basicRunning + basicFailed + basicSucceeded + enhanceTotal
	basicProgress := 0
	if basicTotal > 0 {
		basicProgress = int((basicTotal - basicPending - basicRunning - basicFailed) * 100 / basicTotal)
	}

	// 转换失败文档格式
	basicFailedDocsResp := make([]v1.FailedDoc, len(basicFailedDocs))
	for i, doc := range basicFailedDocs {
		basicFailedDocsResp[i] = v1.FailedDoc{
			NodeID:   doc.NodeID,
			NodeName: doc.NodeName,
			Reason:   doc.Reason,
		}
	}

	enhanceFailedDocsResp := make([]v1.FailedDoc, len(enhanceFailedDocs))
	for i, doc := range enhanceFailedDocs {
		enhanceFailedDocsResp[i] = v1.FailedDoc{
			NodeID:   doc.NodeID,
			NodeName: doc.NodeName,
			Reason:   doc.Reason,
		}
	}

	// 获取Docker组件状态
	dockerComponents, err := utilsDocker.GetDockerComposeStatusSimple(ctx)
	if err != nil {
		u.logger.Warn("failed to get docker status", log.Error(err))
		// 不返回错误，只记录警告
		dockerComponents = []utilsDocker.ComponentStatus{}
	}

	systemComponents := make([]v1.ComponentStatus, len(dockerComponents))
	for i, comp := range dockerComponents {
		systemComponents[i] = v1.ComponentStatus{
			Name:      comp.Name,
			Status:    comp.Status,
			Image:     comp.Image,
			Ports:     comp.Ports,
			Health:    comp.Health,
			LogStatus: comp.LogStatus,
		}
	}

	return &v1.SystemResp{
		Document: v1.DocumentInfo{
			CurrentCount:      currentCount,
			NewIn24h:          newIn24h,
			LearningSucceeded: learningSucceeded,
			LearningFailed:    learningFailed,
		},
		Learning: v1.LearningInfo{
			BasicProcessing: v1.QueueProgress{
				Pending:  basicPending,
				Running:  basicRunning,
				Total:    basicTotal,
				Progress: basicProgress,
			},
			BasicFailed: basicFailed,
			EnhanceProcessing: v1.QueueProgress{
				Pending:  enhancePending,
				Running:  enhanceRunning,
				Total:    enhanceTotal,
				Progress: enhanceProgress,
			},
			EnhanceFailed:     enhanceFailed,
			BasicFailedDocs:   basicFailedDocsResp,
			EnhanceFailedDocs: enhanceFailedDocsResp,
		},
		System: v1.SystemInfo{
			Components: systemComponents,
		},
	}, nil
}

// GetContainerLogs 获取容器日志
func (u *SystemUseCase) GetContainerLogs(ctx context.Context, containerName string, page, limit int) (*v1.ContainerLogsResp, error) {
	// 设置默认值
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 100
	}

	// 获取分页日志
	logs, hasMore, err := utilsDocker.GetContainerLogsPaginated(ctx, containerName, page, limit)
	if err != nil {
		u.logger.Error("failed to get container logs",
			log.String("container", containerName),
			log.Error(err))
		return nil, fmt.Errorf("failed to get container logs: %w", err)
	}

	// 转换日志格式
	logEntries := make([]v1.LogEntry, len(logs))
	for i, log := range logs {
		logEntries[i] = v1.LogEntry{
			Timestamp: log.Timestamp,
			Message:   log.Message,
			Level:     log.Level,
		}
	}

	return &v1.ContainerLogsResp{
		Logs:    logEntries,
		HasMore: hasMore,
		Total:   int64(len(logEntries)),
	}, nil
}
