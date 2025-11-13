import { Box, Stack, Typography, Chip } from '@mui/material';
import { getApiV1System } from '@/request/Stat';
import Card from '@/components/Card';
import { Icon } from '@ctzhian/ui';
import { useAppSelector } from '@/store';
import { addOpacityToColor } from '@/utils';
import { Ellipsis } from '@ctzhian/ui';
import BlueCard from '@/assets/images/blueCard.png';
import PurpleCard from '@/assets/images/purpleCard.png';
import Nodata from '@/assets/images/nodata.png';
import { useEffect, useState } from 'react';
import { V1SystemResp } from '@/request/Stat';
import SvgIcon from '@/components/SvgIcon';

const System = () => {
  const { kb_id = '' } = useAppSelector(state => state.config);
  const [data, setData] = useState<V1SystemResp | null>(null);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    if (!kb_id) return;
    setLoading(true);
    getApiV1System({ kb_id })
      .then(res => {
        setData(res || null);
      })
      .catch(err => {
        console.error('Failed to get system:', err);
      })
      .finally(() => {
        setLoading(false);
      });
  }, [kb_id]);

  if (loading) {
    return <Box sx={{ p: 2 }}>加载中...</Box>;
  }

  if (!data) {
    return <Box sx={{ p: 2 }}>暂无数据</Box>;
  }

  // 文档统计卡片
  const documentCards = [
    {
      label: '当前文档数',
      value: data.document.current_count,
      color: '#021D70',
      bg: 'linear-gradient( 180deg, #D7EBFD 0%, #BEDDFD 100%)',
      image: BlueCard,
    },
    {
      label: '24h新增文档数',
      value: data.document.new_in_24h,
      color: '#021D70',
      bg: 'linear-gradient( 180deg, #D7EBFD 0%, #BEDDFD 100%)',
      image: BlueCard,
    },
    {
      label: '学习成功数量',
      value: data.document.learning_succeeded,
      color: '#021D70',
      bg: 'linear-gradient( 180deg, #D7EBFD 0%, #BEDDFD 100%)',
      image: BlueCard,
    },
    {
      label: '学习失败数量',
      value: data.document.learning_failed,
      color: '#260A7A',
      bg: 'linear-gradient( 180deg, #F0DDFF 0%, #E6C8FF 100%)',
      image: PurpleCard,
    },
  ];

  const getStatusColor = (status: string) => {
    if (status.includes('running') || status.includes('Up')) {
      return 'success';
    }
    if (status.includes('stopped') || status.includes('Exited')) {
      return 'error';
    }
    return 'warning';
  };

  const getHealthColor = (health: string) => {
    if (health === 'healthy') {
      return 'success';
    }
    if (health === 'unhealthy') {
      return 'error';
    }
    return 'default';
  };

  // 容器图标映射
  const getContainerIcon = (name: string): string => {
    const nameLower = name.toLowerCase();
    if (nameLower.includes('nginx')) return 'nginx';
    if (nameLower.includes('postgres')) return 'postgres';
    if (nameLower.includes('redis')) return 'redis';
    if (nameLower.includes('minio')) return 'minio';
    if (nameLower.includes('qdrant')) return 'qdrant';
    if (nameLower.includes('raglite')) return 'raglite';
    if (nameLower.includes('docker')) return 'docker';
    if (nameLower.includes('nats')) return 'nats';
    if (nameLower.includes('caddy')) return 'caddy';
    if (nameLower.includes('crawler')) return 'crawler';
    if (nameLower.includes('api')) return 'icon-zhinengwenda';
    if (nameLower.includes('app')) return 'icon-zhinengwenda';
    if (nameLower.includes('consumer')) return 'icon-zhinengwenda';
    return 'docker';
  };

  // 状态指示灯颜色
  const getStatusIndicatorColor = (status: string, health?: string) => {
    if (health === 'unhealthy') {
      return '#E53E3E'; // 红色
    }
    if (health === 'degraded') {
      return '#DD6B20'; // 橙色
    }
    if (health === 'downgrade' || health === 'warning' || health === 'error') {
      return '#DD6B20'; // 橙色警告状态
    }
    if (status.includes('stopped') || status.includes('Exited') || status.includes('Failed')) {
      return '#E53E3E'; // 红色
    }
    if (status.includes('running') || status.includes('Up') || health === 'healthy') {
      return '#38A169'; // 绿色
    }
    return '#718096'; // 灰色
  };

  // 简化容器名称显示
  const getDisplayName = (name: string): string => {
    return name.replace('panda-wiki-', '').replace('panda-wiki', 'core');
  };

  // 判断是否需要动态状态灯（仅 RAGLite 和 Qdrant）
  const needsPulseAnimation = (name: string): boolean => {
    const nameLower = name.toLowerCase();
    return nameLower.includes('raglite') || nameLower.includes('qdrant');
  };

  return (
    <Box sx={{ p: 2 }}>
      {/* 容器状态 */}
      <Box sx={{ mb: 3 }}>
        <Typography variant="h6" sx={{ mb: 2, fontWeight: 'bold' }}>
          容器
        </Typography>
        {data.system.components.length === 0 ? (
          <Card sx={{ p: 4, boxShadow: '0px 4px 8px rgba(0, 0, 0, 0.1)' }}>
            <Stack
              alignItems={'center'}
              justifyContent={'center'}
              sx={{ fontSize: 12, color: 'text.disabled' }}
            >
              <img src={Nodata} width={100} />
              未找到容器
            </Stack>
          </Card>
        ) : (
          <Box
            sx={{
              display: 'grid',
              gridTemplateColumns: 'repeat(auto-fill, minmax(260px, 1fr))',
              gap: 2,
            }}
          >
            {data.system.components.map((comp, idx) => {
              const displayName = getDisplayName(comp.name);
              const icon = getContainerIcon(comp.name);
              const statusColor = getStatusIndicatorColor(comp.status, comp.health);
              const hasPulse = needsPulseAnimation(comp.name) && comp.health;

              return (
                <Card
                  key={idx}
                  sx={{
                    p: 2,
                    boxShadow: '0px 4px 8px rgba(0, 0, 0, 0.1)',
                    transition: 'all 0.3s ease',
                    '&:hover': {
                      boxShadow: '0px 6px 16px rgba(0, 0, 0, 0.15)',
                    },
                  }}
                >
                  <Stack spacing={1.5}>
                    {/* 图标和状态指示器 */}
                    <Stack direction="row" alignItems="center" justifyContent="space-between">
                      <Box
                        sx={{
                          width: 40,
                          height: 40,
                          borderRadius: '8px',
                          bgcolor: 'background.paper3',
                          display: 'flex',
                          alignItems: 'center',
                          justifyContent: 'center',
                        }}
                      >
                        {icon.startsWith('icon-') ? (
                          <Icon
                            type={icon}
                            sx={{
                              fontSize: 24,
                              color: 'primary.main',
                            }}
                          />
                        ) : (
                          <SvgIcon name={icon} size={24} color="#1976d2" />
                        )}
                      </Box>
                      <Box
                        sx={{
                          width: 10,
                          height: 10,
                          borderRadius: '50%',
                          bgcolor: statusColor,
                          position: 'relative',
                          ...(hasPulse && {
                            animation: 'pulse 2s ease-in-out infinite',
                            '@keyframes pulse': {
                              '0%, 100%': {
                                boxShadow: `0 0 0 0 ${statusColor}66`,
                              },
                              '50%': {
                                boxShadow: `0 0 0 6px ${statusColor}00`,
                              },
                            },
                          }),
                        }}
                      />
                    </Stack>

                    {/* 容器名称 */}
                    <Box>
                      <Typography
                        sx={{
                          fontSize: 14,
                          fontWeight: 700,
                          color: 'text.primary',
                          mb: 0.5,
                          textTransform: 'capitalize',
                        }}
                      >
                        {displayName}
                      </Typography>
                      <Ellipsis
                        sx={{
                          fontSize: 11,
                          color: 'text.secondary',
                        }}
                      >
                        {comp.image.split(':')[0].split('/').pop()}
                      </Ellipsis>
                    </Box>

                    {/* 状态标签 */}
                    <Stack direction="row" gap={0.5} flexWrap="wrap">
                      {comp.health && (
                        <Chip
                          label={comp.health}
                          size="small"
                          sx={{
                            height: 22,
                            fontSize: 11,
                            fontWeight: 600,
                            bgcolor: `${statusColor}15`,
                            color: statusColor,
                            border: 'none',
                          }}
                        />
                      )}
                      {comp.ports && (
                        <Chip
                          label={comp.ports.split(',')[0]}
                          size="small"
                          sx={{
                            height: 22,
                            fontSize: 11,
                            fontWeight: 600,
                            bgcolor: 'background.paper3',
                            color: 'text.secondary',
                            border: 'none',
                          }}
                        />
                      )}
                    </Stack>

                    {/* 日志状态（如果有） */}
                    {comp.log_status && (
                      <Box
                        sx={{
                          fontSize: 11,
                          color: 'text.secondary',
                          bgcolor: 'background.paper3',
                          p: 1,
                          borderRadius: 1,
                        }}
                      >
                        <Ellipsis>{comp.log_status}</Ellipsis>
                      </Box>
                    )}
                  </Stack>
                </Card>
              );
            })}
          </Box>
        )}
      </Box>

      {/* 学习状态 */}
      <Box sx={{ mb: 3 }}>
        <Typography variant="h6" sx={{ mb: 2, fontWeight: 'bold' }}>
          学习
        </Typography>

        {/* 第一行：队列进度和失败数 */}
        <Stack direction={'row'} gap={2} sx={{ mb: 2 }}>
          {/* 基础处理队列进度 */}
          <Card sx={{ flex: 1, p: 2, boxShadow: '0px 4px 8px rgba(0, 0, 0, 0.1)' }}>
            <Box sx={{ fontSize: 14, fontWeight: 600, mb: 1.5, color: 'text.primary' }}>
              基础处理队列进度
            </Box>
            <Box sx={{ fontSize: 24, fontWeight: 700, mb: 1, color: 'primary.main' }}>
              {data.learning.basic_processing.progress}%
            </Box>
            <Box
              sx={{
                height: 8,
                mb: 1,
                borderRadius: '4px',
                bgcolor: 'background.paper3',
              }}
            >
              <Box
                sx={{
                  height: 8,
                  background: 'linear-gradient( 90deg, #3248F2 0%, #9E68FC 100%)',
                  width: `${data.learning.basic_processing.progress}%`,
                  borderRadius: '4px',
                }}
              ></Box>
            </Box>
            <Stack direction="row" justifyContent="space-between" sx={{ fontSize: 12, color: 'text.secondary' }}>
              <span>等待: {data.learning.basic_processing.pending}</span>
              <span>运行: {data.learning.basic_processing.running}</span>
              <span>总数: {data.learning.basic_processing.total}</span>
            </Stack>
          </Card>

          {/* 基础处理失败数 */}
          <Card sx={{ flex: 1, p: 2, boxShadow: '0px 4px 8px rgba(0, 0, 0, 0.1)' }}>
            <Box sx={{ fontSize: 14, fontWeight: 600, mb: 1.5, color: 'text.primary' }}>
              基础处理失败数
            </Box>
            <Box sx={{ fontSize: 32, fontWeight: 700, color: 'error.main' }}>
              {data.learning.basic_failed}
            </Box>
          </Card>

          {/* 增强处理队列进度 */}
          <Card sx={{ flex: 1, p: 2, boxShadow: '0px 4px 8px rgba(0, 0, 0, 0.1)' }}>
            <Box sx={{ fontSize: 14, fontWeight: 600, mb: 1.5, color: 'text.primary' }}>
              增强处理队列进度
            </Box>
            <Box sx={{ fontSize: 24, fontWeight: 700, mb: 1, color: 'primary.main' }}>
              {data.learning.enhance_processing.progress}%
            </Box>
            <Box
              sx={{
                height: 8,
                mb: 1,
                borderRadius: '4px',
                bgcolor: 'background.paper3',
              }}
            >
              <Box
                sx={{
                  height: 8,
                  background: 'linear-gradient( 90deg, #3248F2 0%, #9E68FC 100%)',
                  width: `${data.learning.enhance_processing.progress}%`,
                  borderRadius: '4px',
                }}
              ></Box>
            </Box>
            <Stack direction="row" justifyContent="space-between" sx={{ fontSize: 12, color: 'text.secondary' }}>
              <span>等待: {data.learning.enhance_processing.pending}</span>
              <span>运行: {data.learning.enhance_processing.running}</span>
              <span>总数: {data.learning.enhance_processing.total}</span>
            </Stack>
          </Card>

          {/* 增强处理失败数 */}
          <Card sx={{ flex: 1, p: 2, boxShadow: '0px 4px 8px rgba(0, 0, 0, 0.1)' }}>
            <Box sx={{ fontSize: 14, fontWeight: 600, mb: 1.5, color: 'text.primary' }}>
              增强处理失败数
            </Box>
            <Box sx={{ fontSize: 32, fontWeight: 700, color: 'error.main' }}>
              {data.learning.enhance_failed}
            </Box>
          </Card>
        </Stack>

        {/* 第二行：失败文档列表 - 参照热门文档样式 */}
        <Stack direction={'row'} gap={2}>
          {/* 基础处理失败文档 */}
          <Card sx={{ flex: 1, p: 2, boxShadow: '0px 4px 8px rgba(0, 0, 0, 0.1)', height: 400 }}>
            <Box sx={{ fontSize: 16, fontWeight: 'bold', mb: 2 }}>
              基础处理失败文档
            </Box>
            {data.learning.basic_failed_docs.length === 0 ? (
              <Stack
                alignItems={'center'}
                justifyContent={'center'}
                sx={{ fontSize: 12, color: 'text.disabled', height: 'calc(100% - 40px)' }}
              >
                <img src={Nodata} width={100} />
                无失败文档
              </Stack>
            ) : (
              <Box sx={{ maxHeight: 'calc(100% - 40px)', overflowY: 'auto' }}>
                <Stack gap={2}>
                  {data.learning.basic_failed_docs.map((doc, index) => (
                    <Box key={index} sx={{ fontSize: 12 }}>
                      <Stack
                        direction={'row'}
                        alignItems={'center'}
                        justifyContent={'space-between'}
                        gap={1}
                        sx={{ mb: 0.5 }}
                      >
                        <Ellipsis sx={{ flex: 1, fontWeight: 600, color: 'text.primary' }}>
                          {doc.node_name || '-'}
                        </Ellipsis>
                      </Stack>
                      <Box
                        sx={{
                          fontSize: 11,
                          color: 'error.main',
                          bgcolor: 'error.lighter',
                          p: 1,
                          borderRadius: 1,
                        }}
                      >
                        <Ellipsis>
                          {doc.reason || '未知原因'}
                        </Ellipsis>
                      </Box>
                    </Box>
                  ))}
                </Stack>
              </Box>
            )}
          </Card>

          {/* 增强处理失败文档 */}
          <Card sx={{ flex: 1, p: 2, boxShadow: '0px 4px 8px rgba(0, 0, 0, 0.1)', height: 400 }}>
            <Box sx={{ fontSize: 16, fontWeight: 'bold', mb: 2 }}>
              增强处理失败文档
            </Box>
            {data.learning.enhance_failed_docs.length === 0 ? (
              <Stack
                alignItems={'center'}
                justifyContent={'center'}
                sx={{ fontSize: 12, color: 'text.disabled', height: 'calc(100% - 40px)' }}
              >
                <img src={Nodata} width={100} />
                无失败文档
              </Stack>
            ) : (
              <Box sx={{ maxHeight: 'calc(100% - 40px)', overflowY: 'auto' }}>
                <Stack gap={2}>
                  {data.learning.enhance_failed_docs.map((doc, index) => (
                    <Box key={index} sx={{ fontSize: 12 }}>
                      <Stack
                        direction={'row'}
                        alignItems={'center'}
                        justifyContent={'space-between'}
                        gap={1}
                        sx={{ mb: 0.5 }}
                      >
                        <Ellipsis sx={{ flex: 1, fontWeight: 600, color: 'text.primary' }}>
                          {doc.node_name || '-'}
                        </Ellipsis>
                      </Stack>
                      <Box
                        sx={{
                          fontSize: 11,
                          color: 'error.main',
                          bgcolor: 'error.lighter',
                          p: 1,
                          borderRadius: 1,
                        }}
                      >
                        <Ellipsis>
                          {doc.reason || '未知原因'}
                        </Ellipsis>
                      </Box>
                    </Box>
                  ))}
                </Stack>
              </Box>
            )}
          </Card>
        </Stack>
      </Box>

    </Box>
  );
};

export default System;
