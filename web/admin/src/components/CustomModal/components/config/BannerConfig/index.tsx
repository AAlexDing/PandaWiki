import React, { useState, useEffect } from 'react';
import { TextField, Chip, Autocomplete, Box } from '@mui/material';
import { CommonItem, StyledCommonWrapper } from '../../components/StyledCommon';
import type { ConfigProps } from '../type';
import { useForm, Controller } from 'react-hook-form';
import { useAppSelector } from '@/store';
import DragList from './DragList';
import UploadFile from '@/components/UploadFile';
import { DEFAULT_DATA } from '../../../constants';
import useDebounceAppPreviewData from '@/hooks/useDebounceAppPreviewData';
import { handleLandingConfigs, findConfigById } from '../../../utils';

const Config: React.FC<ConfigProps> = ({ setIsEdit, id }) => {
  const { appPreviewData } = useAppSelector(state => state.config);
  const [inputValue, setInputValue] = useState('');
  const { control, watch, setValue, subscribe } = useForm<
    typeof DEFAULT_DATA.banner
  >({
    defaultValues: findConfigById(
      appPreviewData?.settings?.web_app_landing_configs || [],
      id,
    ),
  });

  const debouncedDispatch = useDebounceAppPreviewData();
  const btns = watch('btns') || [];

  const handleAddButton = () => {
    const nextId = `${Date.now()}`;
    setValue('btns', [
      ...(btns || []),
      { id: nextId, text: '', type: 'contained', href: '' },
    ]);
  };

  useEffect(() => {
    const callback = subscribe({
      formState: {
        values: true,
      },
      callback: ({ values }) => {
        const previewData = {
          ...appPreviewData,
          settings: {
            ...appPreviewData?.settings,
            web_app_landing_configs: handleLandingConfigs({
              id,
              config: appPreviewData?.settings?.web_app_landing_configs || [],
              values: {
                ...values,
              },
            }),
          },
        };
        setIsEdit(true);
        debouncedDispatch(previewData);
      },
    });
    return () => {
      callback();
    };
  }, [subscribe]);

  return (
    <StyledCommonWrapper>
      <CommonItem title='主标题'>
        <Controller
          control={control}
          name='title'
          render={({ field }) => (
            <TextField
              label='文字'
              slotProps={{
                inputLabel: {
                  shrink: true,
                },
              }}
              fullWidth
              {...field}
              placeholder='请输入'
            />
          )}
        />
      </CommonItem>
      <CommonItem title='副标题'>
        <Controller
          control={control}
          name='subtitle'
          render={({ field }) => (
            <TextField
              label='文字'
              slotProps={{
                inputLabel: {
                  shrink: true,
                },
              }}
              multiline
              minRows={2}
              fullWidth
              {...field}
              placeholder='请输入'
            />
          )}
        />
      </CommonItem>
      <CommonItem title='背景图' desc='(推荐 1920*720)'>
        <Controller
          control={control}
          name='bg_url'
          render={({ field }) => (
            <UploadFile
              name='bg_url'
              id='bannerconfig_bgUrl'
              type='url'
              disabled={false}
              accept='image/*'
              width={354}
              height={129}
              value={field.value || ''}
              onChange={(url: string) => {
                field.onChange(url);
                setIsEdit(true);
              }}
            />
          )}
        />
      </CommonItem>
      <CommonItem title='搜索框'>
        <Controller
          control={control}
          name='placeholder'
          render={({ field }) => <TextField {...field} placeholder='请输入' />}
        />
      </CommonItem>
      <CommonItem title='热门搜索'>
        <Controller
          control={control}
          name='hot_search'
          render={({ field }) => (
            <Autocomplete
              {...field}
              value={field.value || []}
              multiple
              freeSolo
              fullWidth
              options={[]}
              inputValue={inputValue}
              onInputChange={(_, newInputValue) => setInputValue(newInputValue)}
              onChange={(_, newValue) => {
                setIsEdit(true);
                const newValues = [...new Set(newValue as string[])];
                field.onChange(newValues);
              }}
              onBlur={() => {
                setIsEdit(true);
                const trimmedValue = inputValue.trim();
                if (trimmedValue && !field.value?.includes(trimmedValue)) {
                  field.onChange([...(field.value || []), trimmedValue]);
                }
                setInputValue('');
              }}
              renderValue={(value, getTagProps) => {
                return value.map((option, index: number) => {
                  return (
                    <Chip
                      variant='outlined'
                      size='small'
                      label={
                        <Box sx={{ fontSize: '12px' }}>
                          {option as React.ReactNode}
                        </Box>
                      }
                      {...getTagProps({ index })}
                      key={index}
                    />
                  );
                });
              }}
              renderInput={params => (
                <TextField
                  {...params}
                  placeholder='回车确认，填写下一个热门搜索'
                  variant='outlined'
                />
              )}
            />
          )}
        />
      </CommonItem>
      <CommonItem title='主按钮' onAdd={handleAddButton}>
        <DragList
          data={btns as Required<(typeof btns)[0]>[]}
          onChange={btns => {
            setValue('btns', btns);
            setIsEdit(true);
          }}
          setIsEdit={setIsEdit}
        />
      </CommonItem>
    </StyledCommonWrapper>
  );
};

export default Config;
