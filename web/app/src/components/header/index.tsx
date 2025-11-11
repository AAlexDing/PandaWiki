'use client';

import Logo from '@/assets/images/logo.png';
import { Box } from '@mui/material';
import { useStore } from '@/provider';
import { usePathname } from 'next/navigation';
import { useMemo } from 'react';
import {
  Header as CustomHeader,
  WelcomeHeader as WelcomeHeaderComponent,
} from '@panda-wiki/ui';
import QaModal from '../QaModal';
import ThemeSwitch from './themeSwitch';
interface HeaderProps {
  isDocPage?: boolean;
  isWelcomePage?: boolean;
}

const Header = ({ isDocPage = false, isWelcomePage = false }: HeaderProps) => {
  const { mobile = false, kbDetail, catalogWidth, setQaModalOpen } = useStore();
  const pathname = usePathname();
  const docWidth = useMemo(() => {
    if (isWelcomePage) return 'full';
    return kbDetail?.settings?.theme_and_style?.doc_width || 'full';
  }, [kbDetail, isWelcomePage]);

  const handleSearch = (value?: string, type: 'chat' | 'search' = 'chat') => {
    if (value?.trim()) {
      if (type === 'chat') {
        sessionStorage.setItem('chat_search_query', value.trim());
        setQaModalOpen?.(true);
      } else {
        sessionStorage.setItem('chat_search_query', value.trim());
      }
    }
  };

  const showSearch = useMemo(() => {
    return pathname !== '/welcome' && !pathname.startsWith('/chat');
  }, [pathname]);

  return (
    <CustomHeader
      isDocPage={isDocPage}
      mobile={mobile}
      docWidth={docWidth}
      catalogWidth={catalogWidth}
      logo={kbDetail?.settings?.icon || Logo.src}
      title={kbDetail?.settings?.title}
      placeholder={
        kbDetail?.settings?.web_app_custom_style?.header_search_placeholder
      }
      showSearch
      btns={kbDetail?.settings?.btns}
      onSearch={handleSearch}
      onQaClick={() => setQaModalOpen?.(true)}
    >
      <Box sx={{ ml: 2 }}>
        <ThemeSwitch />
      </Box>
      <QaModal />
    </CustomHeader>
  );
};

export const WelcomeHeader = () => {
  const { mobile = false, kbDetail, catalogWidth, setQaModalOpen } = useStore();
  const handleSearch = (value?: string, type: 'chat' | 'search' = 'chat') => {
    if (value?.trim()) {
      if (type === 'chat') {
        sessionStorage.setItem('chat_search_query', value.trim());
        setQaModalOpen?.(true);
      } else {
        sessionStorage.setItem('chat_search_query', value.trim());
      }
    }
  };
  return (
    <WelcomeHeaderComponent
      isDocPage={false}
      mobile={mobile}
      docWidth='full'
      catalogWidth={catalogWidth}
      logo={kbDetail?.settings?.icon || Logo.src}
      title={kbDetail?.settings?.title}
      placeholder={
        kbDetail?.settings?.web_app_custom_style?.header_search_placeholder
      }
      showSearch
      btns={kbDetail?.settings?.btns}
      onSearch={handleSearch}
      onQaClick={() => setQaModalOpen?.(true)}
    >
      <QaModal />
    </WelcomeHeaderComponent>
  );
};

export default Header;
