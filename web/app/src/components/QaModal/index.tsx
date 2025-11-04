'use client';
import React, { useState, useEffect, useRef, useMemo } from 'react';
import Logo from '@/assets/images/logo.png';
import { IconZhinengwenda, IconJinsousuo } from '@panda-wiki/icons';
import { useSearchParams } from 'next/navigation';
import {
  Box,
  Button,
  Typography,
  Modal,
  Stack,
  lighten,
  alpha,
} from '@mui/material';
import { CusTabs } from '@ctzhian/ui';
import AiQaContent from './AiQaContent';
import SearchDocContent from './SearchDocContent';
import { useStore } from '@/provider';

interface SearchSuggestion {
  id: string;
  title: string;
  description?: string;
  type?: 'recent' | 'suggestion' | 'trending';
}

interface QaModalProps {
  placeholder?: string;
  initialValue?: string;
  onSearch?: (value?: string, type?: 'search' | 'chat') => void;
  onSearchSuggestions?: (query: string) => Promise<SearchSuggestion[]>;
  defaultSuggestions?: SearchSuggestion[];
}

const QaModal: React.FC<QaModalProps> = () => {
  const { qaModalOpen, setQaModalOpen, kbDetail, mobile } = useStore();
  const [searchMode, setSearchMode] = useState<'chat' | 'search'>('chat');
  const inputRef = useRef<HTMLInputElement>(null);
  const aiQaInputRef = useRef<HTMLInputElement>(null);
  const searchParams = useSearchParams();
  const onClose = () => {
    setQaModalOpen?.(false);
  };

  const placeholder = useMemo(() => {
    return (
      kbDetail?.settings?.web_app_custom_style?.header_search_placeholder ||
      '搜索...'
    );
  }, [kbDetail]);

  const hotSearch = useMemo(() => {
    const bannerConfig = kbDetail?.settings?.web_app_landing_configs?.find(
      item => item.type === 'banner',
    );
    return bannerConfig?.banner_config?.hot_search || [];
  }, [kbDetail]);

  // modal打开时自动聚焦
  useEffect(() => {
    if (qaModalOpen) {
      setTimeout(() => {
        if (searchMode === 'chat') {
          aiQaInputRef.current?.querySelector('textarea')?.focus();
        } else {
          inputRef.current?.querySelector('input')?.focus();
        }
      }, 100);
    }
  }, [qaModalOpen, searchMode]);

  useEffect(() => {
    const cid = searchParams.get('cid');
    const ask = searchParams.get('ask');
    if (cid || ask) {
      setQaModalOpen?.(true);
    }
  }, []);

  return (
    <Modal
      open={qaModalOpen as boolean}
      onClose={onClose}
      sx={{
        display: 'flex',
        justifyContent: 'center',
        alignItems: 'flex-start',
        p: 2,
      }}
    >
      <Box
        sx={theme => ({
          display: 'flex',
          flexDirection: 'column',
          gap: 2,
          flex: 1,
          maxWidth: 800,
          maxHeight: '100%',
          backgroundColor: lighten(theme.palette.background.default, 0.05),
          borderRadius: '10px',
          boxShadow: '0 8px 32px rgba(0, 0, 0, 0.12)',
          overflow: 'hidden',
          outline: 'none',
          '& .Mui-selected': {
            color: `${theme.palette.background.default} !important`,
          },
        })}
        onClick={e => e.stopPropagation()}
      >
        {/* 头部区域 */}
        <Box
          sx={{
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'space-between',
            mx: 2,
            py: 1.5,
            borderBottom: '1px solid',
            borderColor: 'divider',
          }}
        >
          {/* Logo */}
          <Box sx={{ display: 'flex', alignItems: 'center', gap: 1.5 }}>
            <img
              src={kbDetail?.settings?.icon || Logo.src}
              style={{
                width: 24,
                height: 24,
              }}
            />
            <Typography
              variant='h6'
              sx={{ fontSize: 16, color: 'text.primary' }}
            >
              {kbDetail?.settings?.title}
            </Typography>
          </Box>

          <CusTabs
            size='small'
            list={[
              {
                label: (
                  <Stack direction='row' gap={1} alignItems='center'>
                    <IconZhinengwenda sx={{ fontSize: 16 }} />
                    {!mobile && <span>智能问答</span>}
                  </Stack>
                ),
                value: 'chat',
              },
              {
                label: (
                  <Stack direction='row' gap={1} alignItems='center'>
                    <IconJinsousuo sx={{ fontSize: 16 }} />
                    {!mobile && <span>仅搜索文档</span>}
                  </Stack>
                ),
                value: 'search',
              },
            ]}
            value={searchMode}
            onChange={value => setSearchMode(value as 'chat' | 'search')}
          />

          {/* Esc按钮 */}
          {!mobile && (
            <Button
              onClick={onClose}
              size='small'
              sx={theme => ({
                minWidth: 'auto',
                px: 1.5,
                py: 0.5,
                fontSize: 12,
                fontWeight: 500,
                textTransform: 'none',
              })}
            >
              Esc
            </Button>
          )}
        </Box>
        {/* 主内容区域 - 根据模式切换 */}
        <Box sx={{ px: 2, display: searchMode === 'chat' ? 'block' : 'none' }}>
          <AiQaContent
            hotSearch={hotSearch}
            placeholder={placeholder}
            inputRef={aiQaInputRef}
          />
        </Box>
        <Box
          sx={{ px: 2, display: searchMode === 'search' ? 'block' : 'none' }}
        >
          <SearchDocContent inputRef={inputRef} placeholder={placeholder} />
        </Box>
        {/* 底部AI生成提示 */}
        <Box
          sx={{
            px: 3,
            py: kbDetail?.settings?.disclaimer_settings?.content ? 2 : 1,
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
          }}
        >
          <Typography
            variant='caption'
            sx={{
              color: 'text.disabled',
              fontSize: 12,
              display: 'flex',
              alignItems: 'center',
              gap: 1,
            }}
          >
            <Box>{kbDetail?.settings?.disclaimer_settings?.content}</Box>
          </Typography>
        </Box>
      </Box>
    </Modal>
  );
};

export default QaModal;
