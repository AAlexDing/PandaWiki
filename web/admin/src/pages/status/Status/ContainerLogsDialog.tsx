import { useState, useCallback, useRef, useEffect } from 'react';
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  Typography,
  Box,
  CircularProgress,
  IconButton,
  Paper,
  Divider,
  Chip,
  Stack,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  OutlinedInput,
  InputAdornment,
  Tooltip,
} from '@mui/material';
import {
  Close as CloseIcon,
  Refresh as RefreshIcon,
  ArrowUpward as ArrowUpwardIcon,
  ArrowDownward as ArrowDownwardIcon,
  FilterList as FilterListIcon,
  Clear as ClearIcon,
} from '@mui/icons-material';
import { getApiV1SystemLogsContainerName, V1ContainerLogsResp } from '@/request/Stat';

interface ContainerLogsDialogProps {
  open: boolean;
  onClose: () => void;
  container: {
    name: string;
    status: string;
    image: string;
  } | null;
}

interface LogEntry {
  id: string;
  timestamp: string;
  message: string;
  level?: 'info' | 'warn' | 'error' | 'debug';
}

const ContainerLogsDialog: React.FC<ContainerLogsDialogProps> = ({
  open,
  onClose,
  container,
}) => {
  const [logs, setLogs] = useState<LogEntry[]>([]);
  const [loading, setLoading] = useState(false);
  const [hasMore, setHasMore] = useState(true);
  const [autoScroll, setAutoScroll] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [selectedLevels, setSelectedLevels] = useState<Set<LogEntry['level']>>(new Set(['info', 'warn', 'error', 'debug', undefined]));
  const [currentPage, setCurrentPage] = useState(1);
  const [allLogsCount, setAllLogsCount] = useState(0);
  const [showLoadMoreIndicator, setShowLoadMoreIndicator] = useState(false);

  const logsContainerRef = useRef<HTMLDivElement>(null);
  const pageRef = useRef(1);
  const loadingRef = useRef(false);

  // 解析日志行
  const parseLogLine = (line: string): LogEntry => {
    try {
      // 尝试解析 JSON 格式的日志
      const parsed = JSON.parse(line);
      const timestamp = parsed.timestamp || parsed.time || new Date().toISOString();
      const message = parsed.message || parsed.msg || line;
      return {
        id: `${timestamp}_${message.slice(0, 50)}`, // 使用时间戳和消息前50字符作为唯一ID
        timestamp: timestamp,
        message: message,
        level: parsed.level || parsed.severity || 'info',
      };
    } catch {
      // 如果不是 JSON 格式，尝试提取时间戳
      const timestampRegex = /^(\d{4}-\d{2}-\d{2}[T ]\d{2}:\d{2}:\d{2})/;
      const match = line.match(timestampRegex);
      const timestamp = match?.[1] || new Date().toISOString();
      return {
        id: `${timestamp}_${line.slice(0, 50)}`, // 使用时间戳和消息前50字符作为唯一ID
        timestamp: timestamp,
        message: line,
        level: 'info',
      };
    }
  };

  // 获取容器日志
  const fetchLogs = useCallback(async (page: number = 1, isLoadMore: boolean = false, resetPage: boolean = true): Promise<void> => {
    if (!container || loadingRef.current) return Promise.resolve();

    setLoading(true);
    setError(null);
    loadingRef.current = true;

    try {
      const response: V1ContainerLogsResp = await getApiV1SystemLogsContainerName({
        containerName: container.name,
        page,
        limit: 100,
      });

      // 转换 API 响应的日志条目格式
      const newLogs = (response.logs || []).map(log => ({
        id: `${log.timestamp}_${log.message.slice(0, 50)}`, // 使用时间戳和消息前50字符作为唯一ID
        timestamp: log.timestamp,
        message: log.message,
        level: log.level as 'info' | 'warn' | 'error' | 'debug' | undefined,
      }));

      // 如果是加载更多，将新日志追加到前面（历史日志）
      if (isLoadMore) {
        setLogs(prev => [...newLogs, ...prev]);
        setCurrentPage(prev => prev + 1);
      } else {
        setLogs(newLogs);
        setCurrentPage(page);
        setAllLogsCount(response.total);
      }

      setHasMore(response.has_more);
    } catch (err) {
      setError(err instanceof Error ? err.message : '获取日志失败');
    } finally {
      setLoading(false);
      loadingRef.current = false;
    }
    return Promise.resolve();
  }, [container]);

  
  // 处理滚动事件
  const handleScroll = useCallback(() => {
    if (!logsContainerRef.current) return;

    const { scrollTop, scrollHeight, clientHeight } = logsContainerRef.current;

    // 控制加载提示的显示
    if (scrollTop === 0 && hasMore && !loading) {
      setShowLoadMoreIndicator(true);
    } else {
      setShowLoadMoreIndicator(false);
    }

    // 当滚动到顶部且还有更多日志时，加载更多历史日志
    if (scrollTop === 0 && hasMore && !loading && pageRef.current < 50) { // 限制最大页数避免无限加载
      // 找到第一个完全可见的日志项
      const containerRect = logsContainerRef.current.getBoundingClientRect();
      const logElements = Array.from(logsContainerRef.current.querySelectorAll('[data-log-id]:not([style*="display: none"])'));

      let firstVisibleLogId: string | null = null;
      for (const element of logElements) {
        const elementRect = element.getBoundingClientRect();
        const relativeTop = elementRect.top - containerRect.top;

        // 找到第一个完全可见的日志项
        if (relativeTop >= 0) {
          firstVisibleLogId = element.getAttribute('data-log-id');
          break;
        }
      }

      fetchLogs(currentPage + 1, true, false).then(() => {
        // 加载完成后，定位到之前查看的日志条目
        if (logsContainerRef.current && firstVisibleLogId) {
          requestAnimationFrame(() => {
            const targetElement = logsContainerRef.current?.querySelector(`[data-log-id="${firstVisibleLogId}"]`);
            if (targetElement) {
              const elementRect = targetElement.getBoundingClientRect();
              const containerRect = logsContainerRef.current!.getBoundingClientRect();
              const relativeTop = elementRect.top - containerRect.top;

              logsContainerRef.current!.scrollTop = relativeTop;
            }
          });
        }
      });
    }
  }, [hasMore, loading, fetchLogs, currentPage, pageRef]);

  // 筛选日志
  const filteredLogs = logs.filter(log =>
    log.level ? selectedLevels.has(log.level) : true
  );

  // 加载更多日志
  const loadMoreLogs = useCallback(() => {
    if (!loading && hasMore) {
      fetchLogs(pageRef.current + 1, true);
    }
  }, [loading, hasMore, fetchLogs]);

  // 刷新日志
  const refreshLogs = useCallback(() => {
    setLogs([]);
    pageRef.current = 1;
    setCurrentPage(1);
    setHasMore(true);
    setShowLoadMoreIndicator(false);
    fetchLogs(1, false);
  }, [fetchLogs]);

  // 监听滚动事件
  // 当对话框打开或容器变化时获取日志
  useEffect(() => {
    if (open && container) {
      refreshLogs();
    } else {
      setLogs([]);
      setError(null);
      setSelectedLevels(new Set(['info', 'warn', 'error', 'debug', undefined]));
      setCurrentPage(1);
      setAllLogsCount(0);
      setShowLoadMoreIndicator(false);
    }
  }, [open, container, refreshLogs]);

  // 滚动到顶部
  const scrollToTop = () => {
    if (logsContainerRef.current) {
      logsContainerRef.current.scrollTop = 0;
    }
  };

  // 滚动到底部
  const scrollToBottom = () => {
    if (logsContainerRef.current) {
      logsContainerRef.current.scrollTop = logsContainerRef.current.scrollHeight;
    }
  };

  // 获取日志颜色
  const getLogColor = (level: string) => {
    switch (level?.toLowerCase()) {
      case 'error':
        return '#E53E3E';
      case 'warn':
        return '#DD6B20';
      case 'info':
        return '#38A169';
      case 'debug':
        return '#718096';
      case 'undefined':
      case null:
      case undefined:
        return '#9CA3AF';
      default:
        return '#2D3748';
    }
  };

  // 日志等级选项配置
  const logLevelOptions: { value: LogEntry['level']; label: string; color: string }[] = [
    { value: 'error', label: 'ERROR', color: '#E53E3E' },
    { value: 'warn', label: 'WARN', color: '#DD6B20' },
    { value: 'info', label: 'INFO', color: '#38A169' },
    { value: 'debug', label: 'DEBUG', color: '#718096' },
    { value: undefined, label: 'UNDEF', color: '#9CA3AF' },
  ];

  // 切换日志等级选中状态
  const toggleLogLevel = (level: LogEntry['level']) => {
    setSelectedLevels(prev => {
      const newSet = new Set(prev);
      if (newSet.has(level)) {
        newSet.delete(level);
      } else {
        newSet.add(level);
      }
      return newSet;
    });
  };

  return (
    <Dialog open={open} onClose={onClose} maxWidth="lg" fullWidth>
      <DialogTitle>
        <Box display="flex" alignItems="center" justifyContent="space-between">
          <Typography variant="h6">
            容器日志 - {container?.name}
          </Typography>
          <IconButton onClick={onClose} size="small">
            <CloseIcon />
          </IconButton>
        </Box>
      </DialogTitle>

      <DialogContent sx={{ p: 0 }}>
        {/* 工具栏 */}
        <Box sx={{ p: 2, borderBottom: 1, borderColor: 'divider', display: 'flex', gap: 1 }}>
          <Button
            startIcon={<RefreshIcon />}
            onClick={refreshLogs}
            disabled={loading}
            size="small"
          >
            刷新
          </Button>
          <Button
            startIcon={<ArrowUpwardIcon />}
            onClick={scrollToTop}
            size="small"
          >
            顶部
          </Button>
          <Button
            startIcon={<ArrowDownwardIcon />}
            onClick={scrollToBottom}
            size="small"
          >
            底部
          </Button>
          <Box flex={1} />

          {/* 日志等级筛选标签 */}
          <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
            <Typography variant="caption" color="text.secondary" sx={{ mr: 1 }}>
              筛选:
            </Typography>
            {logLevelOptions.map((option) => (
              <Chip
                key={option.label}
                label={option.label}
                size="small"
                variant={selectedLevels.has(option.value) ? "filled" : "outlined"}
                onClick={() => toggleLogLevel(option.value)}
                sx={{
                  height: 24,
                  fontSize: 10,
                  minWidth: 45,
                  backgroundColor: selectedLevels.has(option.value) ? option.color : 'transparent',
                  color: selectedLevels.has(option.value) ? 'white' : option.color,
                  borderColor: option.color,
                  fontFamily: 'Arial, sans-serif',
                  fontWeight: 'bold',
                  cursor: 'pointer',
                  '& .MuiChip-label': {
                    color: selectedLevels.has(option.value) ? 'white' : option.color,
                  },
                  '&:hover': {
                    backgroundColor: selectedLevels.has(option.value)
                      ? option.color + 'DD'
                      : option.color + '20',
                    transform: 'scale(1.05)',
                  },
                  transition: 'all 0.2s ease-in-out',
                }}
              />
            ))}

            <Tooltip title="重置筛选">
              <IconButton
                onClick={() => setSelectedLevels(new Set(['info', 'warn', 'error', 'debug', undefined]))}
                size="small"
                sx={{ ml: 1 }}
              >
                <ClearIcon />
              </IconButton>
            </Tooltip>
          </Box>

          <Typography variant="caption" color="text.secondary">
            显示 {filteredLogs.length}/{logs.length} 条日志
          </Typography>
        </Box>

        {/* 日志内容 */}
        <Box
          ref={logsContainerRef}
          sx={{
            height: 500,
            overflow: 'auto',
            backgroundColor: '#1e1e1e',
            fontFamily: 'Monaco, Consolas, "Courier New", monospace',
            fontSize: 12,
            p: 2,
          }}
          onScroll={handleScroll}
        >
          {error ? (
            <Box
              display="flex"
              flexDirection="column"
              alignItems="center"
              justifyContent="center"
              height="100%"
              color="error.main"
            >
              <Typography color="error">{error}</Typography>
              <Button onClick={refreshLogs} sx={{ mt: 2 }}>
                重试
              </Button>
            </Box>
          ) : loading && logs.length === 0 ? (
            <Box
              display="flex"
              alignItems="center"
              justifyContent="center"
              height="100%"
            >
              <CircularProgress />
              <Typography sx={{ ml: 2 }}>加载日志中...</Typography>
            </Box>
          ) : (
            <>
              {logs.length === 0 && !loading && (
                <Box
                  display="flex"
                  alignItems="center"
                  justifyContent="center"
                  height="100%"
                  color="text.secondary"
                >
                  <Typography>暂无日志</Typography>
                </Box>
              )}

              {filteredLogs.length === 0 && logs.length > 0 && (
                <Box
                  display="flex"
                  flexDirection="column"
                  alignItems="center"
                  justifyContent="center"
                  height="100%"
                  color="text.secondary"
                  sx={{ pt: 8 }}
                >
                  <Typography>没有符合筛选条件的日志</Typography>
                  <Button
                    onClick={() => setSelectedLevels(new Set(['info', 'warn', 'error', 'debug', undefined]))}
                    sx={{ mt: 2 }}
                    size="small"
                  >
                    显示全部日志
                  </Button>
                </Box>
              )}

              {/* 加载更多指示器 - 只在滑动到顶部时显示 */}
              {showLoadMoreIndicator && hasMore && (
                <Box
                  display="flex"
                  alignItems="center"
                  justifyContent="center"
                  py={2}
                  color="text.secondary"
                >
                  <CircularProgress size={16} sx={{ mr: 1 }} />
                  <Typography variant="caption">向上滚动加载更多日志...</Typography>
                </Box>
              )}

              {filteredLogs.map((log, index) => (
                <Box
                  key={log.id || index}
                  data-log-id={log.id}
                  sx={{
                    mb: 1.5,
                    '&:last-child': { mb: 0 },
                    backgroundColor: 'rgba(255, 255, 255, 0.02)',
                    borderRadius: 1,
                    p: 1,
                    border: '1px solid rgba(255, 255, 255, 0.08)',
                  }}
                >
                  {/* 标签行：时间戳和日志等级 */}
                  <Stack
                    direction="row"
                    spacing={1}
                    alignItems="center"
                    sx={{ mb: 0.5 }}
                  >
                    {/* 时间戳标签 */}
                    <Box
                      sx={{
                        backgroundColor: '#374151',
                        color: '#f3f4f6',
                        px: 1,
                        py: 0.25,
                        borderRadius: 0.5,
                        fontSize: 10,
                        fontFamily: 'Monaco, Consolas, monospace',
                        minWidth: 140,
                        textAlign: 'center',
                        fontWeight: 500,
                      }}
                    >
                      {new Date(log.timestamp).toLocaleString()}
                    </Box>

                    {/* 日志等级标签 */}
                    <Chip
                      label={log.level?.toUpperCase() || 'UNDEF'}
                      size="small"
                      sx={{
                        height: 22,
                        fontSize: 10,
                        minWidth: 55,
                        backgroundColor: getLogColor(log.level || 'undefined'),
                        color: 'white',
                        fontFamily: 'Arial, sans-serif',
                        fontWeight: 'bold',
                        boxShadow: '0 1px 3px rgba(0,0,0,0.2)',
                        '& .MuiChip-label': {
                          color: 'white',
                          padding: '0 8px',
                        },
                      }}
                    />
                  </Stack>

                  {/* 日志消息 */}
                  <Typography
                    component="pre"
                    sx={{
                      color: '#e5e7eb',
                      fontSize: 12,
                      whiteSpace: 'pre-wrap',
                      wordBreak: 'break-word',
                      userSelect: 'text',
                      m: 0,
                      mt: 0.5,
                      pl: 1,
                      borderLeft: `3px solid ${getLogColor(log.level || 'undefined')}`,
                      backgroundColor: 'rgba(0, 0, 0, 0.2)',
                      borderRadius: 0.5,
                      py: 0.5,
                    }}
                  >
                    {log.message}
                  </Typography>
                </Box>
              ))}
            </>
          )}
        </Box>
      </DialogContent>

      <DialogActions>
        <Button onClick={onClose}>关闭</Button>
      </DialogActions>
    </Dialog>
  );
};

export default ContainerLogsDialog;