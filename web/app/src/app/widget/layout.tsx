import StoreProvider from '@/provider';
import { getShareV1AppWidgetInfo } from '@/request/ShareApp';
import React from 'react';

const Layout = async ({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) => {
  const widgetDetail: any = await getShareV1AppWidgetInfo();

  return <StoreProvider widget={widgetDetail}>{children}</StoreProvider>;
};

export default Layout;
