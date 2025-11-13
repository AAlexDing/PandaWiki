package utils

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

// ComponentStatus Docker组件状态
type ComponentStatus struct {
	Name      string
	Status    string
	Image     string
	Ports     string
	Health    string
	LogStatus string
}

// GetDockerComposeStatusSimple 使用 Docker SDK 获取容器状态
func GetDockerComposeStatusSimple(ctx context.Context) ([]ComponentStatus, error) {
	// 创建 Docker 客户端
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("failed to create docker client: %w", err)
	}
	defer cli.Close()

	// 列出所有容器（包括停止的）
	containers, err := cli.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		return nil, fmt.Errorf("failed to list containers: %w", err)
	}

	var components []ComponentStatus

	for _, c := range containers {
		// 获取容器名称（去掉前导斜杠）
		name := strings.TrimPrefix(c.Names[0], "/")

		// 只处理 panda-wiki 相关的容器
		nameLower := strings.ToLower(name)
		if !strings.Contains(nameLower, "panda-wiki") &&
			!strings.Contains(nameLower, "raglite") &&
			!strings.Contains(nameLower, "qdrant") &&
			!strings.Contains(nameLower, "anydoc") {
			continue
		}

		// 构建端口字符串
		ports := buildPortsString(c.Ports)

		comp := ComponentStatus{
			Name:   name,
			Status: c.Status,
			Image:  c.Image,
			Ports:  ports,
		}

		// 对 RAGLite 和 Qdrant 进行特殊日志解析
		if strings.Contains(nameLower, "raglite") {
			comp.Health, comp.LogStatus = parseRAGLiteLogsSDK(ctx, cli, c.ID)
		} else if strings.Contains(nameLower, "qdrant") {
			comp.Health, comp.LogStatus = parseQdrantLogsSDK(ctx, cli, c.ID)
		}

		components = append(components, comp)
	}

	return components, nil
}

// buildPortsString 构建端口字符串
func buildPortsString(ports []types.Port) string {
	if len(ports) == 0 {
		return ""
	}

	var portStrs []string
	for _, port := range ports {
		if port.PublicPort > 0 {
			portStrs = append(portStrs, fmt.Sprintf("%d->%d/%s", port.PublicPort, port.PrivatePort, port.Type))
		} else {
			portStrs = append(portStrs, fmt.Sprintf("%d/%s", port.PrivatePort, port.Type))
		}
	}

	return strings.Join(portStrs, ", ")
}

// parseRAGLiteLogsSDK 使用 SDK 解析 RAGLite 日志
func parseRAGLiteLogsSDK(ctx context.Context, cli *client.Client, containerID string) (health, logStatus string) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	logs, err := getContainerLogs(ctx, cli, containerID, 200)
	if err != nil {
		return "unknown", "failed to read logs"
	}

	lines := strings.Split(logs, "\n")

	// 默认状态
	health = "unknown"
	logStatus = "No recent logs"

	// 查找关键信息
	hasError := false
	hasFatal := false
	isListening := false
	isStarted := false
	lastMeaningfulLine := ""

	for i := len(lines) - 1; i >= 0; i-- {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}

		lineLower := strings.ToLower(line)

		// 记录最后一条有意义的日志
		if lastMeaningfulLine == "" && len(line) > 10 {
			lastMeaningfulLine = line
			// 限制长度
			if len(lastMeaningfulLine) > 200 {
				lastMeaningfulLine = lastMeaningfulLine[:200] + "..."
			}
		}

		// 检查启动和监听状态
		if strings.Contains(lineLower, "listening on") ||
			strings.Contains(lineLower, "server started") ||
			strings.Contains(lineLower, "started on port") {
			isListening = true
		}

		if strings.Contains(lineLower, "started") || strings.Contains(lineLower, "ready") {
			isStarted = true
		}

		// 检查错误
		if strings.Contains(lineLower, "fatal") {
			hasFatal = true
		}
		if strings.Contains(lineLower, "error") && !strings.Contains(lineLower, "level=error") {
			hasError = true
		}
	}

	// 判断健康状态
	if hasFatal {
		health = "unhealthy"
		logStatus = "Fatal error detected: " + lastMeaningfulLine
	} else if hasError {
		health = "degraded"
		logStatus = "Error detected: " + lastMeaningfulLine
	} else if isListening || isStarted {
		health = "healthy"
		logStatus = "Running normally"
	} else if lastMeaningfulLine != "" {
		logStatus = lastMeaningfulLine
	}

	return health, logStatus
}

// parseQdrantLogsSDK 使用 SDK 解析 Qdrant 日志
func parseQdrantLogsSDK(ctx context.Context, cli *client.Client, containerID string) (health, logStatus string) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	logs, err := getContainerLogs(ctx, cli, containerID, 200)
	if err != nil {
		return "unknown", "failed to read logs"
	}

	lines := strings.Split(logs, "\n")

	// 默认状态
	health = "unknown"
	logStatus = "No recent logs"

	// 查找关键信息
	hasPanic := false
	hasError := false
	hasCollectionLoaded := false
	isListening := false
	lastMeaningfulLine := ""
	lastErrorLine := ""

	for i := len(lines) - 1; i >= 0; i-- {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}

		lineLower := strings.ToLower(line)

		// 记录最后一条有意义的日志
		if lastMeaningfulLine == "" && len(line) > 10 {
			lastMeaningfulLine = line
			// 限制长度
			if len(lastMeaningfulLine) > 200 {
				lastMeaningfulLine = lastMeaningfulLine[:200] + "..."
			}
		}

		// 检查 Qdrant 特有的启动标志
		if strings.Contains(lineLower, "qdrant is ready") ||
			strings.Contains(lineLower, "listening") ||
			strings.Contains(lineLower, "access web ui at") {
			isListening = true
		}

		// 检查集合加载
		if strings.Contains(lineLower, "loading collection") ||
			strings.Contains(lineLower, "collection loaded") {
			hasCollectionLoaded = true
		}

		// 检查 panic
		if strings.Contains(lineLower, "panic") {
			hasPanic = true
			if lastErrorLine == "" {
				lastErrorLine = line
			}
		}

		// 检查错误 (Qdrant 的 ERROR 级别日志)
		if strings.Contains(line, "ERROR") ||
			(strings.Contains(lineLower, "error") && strings.Contains(line, "qdrant::startup")) {
			hasError = true
			if lastErrorLine == "" {
				lastErrorLine = line
			}
		}
	}

	// 判断健康状态
	if hasPanic {
		health = "unhealthy"
		if lastErrorLine != "" {
			logStatus = "Panic detected: " + lastErrorLine
		} else {
			logStatus = "Panic detected in logs"
		}
	} else if hasError {
		health = "unhealthy"
		if lastErrorLine != "" {
			logStatus = "Error: " + lastErrorLine
		} else {
			logStatus = "Error detected in logs"
		}
	} else if isListening {
		health = "healthy"
		if hasCollectionLoaded {
			logStatus = "Running - Collections loaded"
		} else {
			logStatus = "Running normally"
		}
	} else if hasCollectionLoaded {
		health = "degraded"
		logStatus = "Collections loaded but not fully ready"
	} else if lastMeaningfulLine != "" {
		logStatus = lastMeaningfulLine
	}

	return health, logStatus
}

// getContainerLogs 获取容器日志
func getContainerLogs(ctx context.Context, cli *client.Client, containerID string, tailLines int) (string, error) {
	options := container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Tail:       fmt.Sprintf("%d", tailLines),
	}

	reader, err := cli.ContainerLogs(ctx, containerID, options)
	if err != nil {
		return "", err
	}
	defer reader.Close()

	// 读取日志内容
	logBytes, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}

	// Docker logs 包含 8 字节的 header，需要去除
	logs := stripDockerLogHeaders(string(logBytes))

	return logs, nil
}

// stripDockerLogHeaders 去除 Docker 日志的 header
func stripDockerLogHeaders(logs string) string {
	lines := strings.Split(logs, "\n")
	var cleanLines []string

	for _, line := range lines {
		// Docker 日志每行前 8 字节是 header: [stream_type, 0, 0, 0, size1, size2, size3, size4]
		// 我们需要跳过这 8 字节
		if len(line) > 8 {
			cleanLines = append(cleanLines, line[8:])
		} else if len(line) > 0 {
			cleanLines = append(cleanLines, line)
		}
	}

	return strings.Join(cleanLines, "\n")
}
