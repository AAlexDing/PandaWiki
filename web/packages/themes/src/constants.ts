import {
  bluePalette,
  greenPalette,
  orangePalette,
  blackPalette,
  deepTealPalette,
  redPalette,
  whitePalette,
  electricBluePalette,
  darkDeepForestPalette,
} from './index';

export const THEME_LIST = [
  {
    label: '魅力蓝',
    value: 'blue',
    palette: bluePalette,
  },
  {
    label: '清新绿',
    value: 'green',
    palette: greenPalette,
  },
  {
    label: '活力橙',
    value: 'orange',
    palette: orangePalette,
  },
  {
    label: '纯净白',
    value: 'black',
    palette: blackPalette,
  },
  {
    label: '幽静深蓝',
    value: 'deepTeal',
    palette: deepTealPalette,
  },
  {
    label: '热情红',
    value: 'red',
    palette: redPalette,
  },
  {
    label: '深墨绿',
    value: 'darkDeepForest',
    palette: darkDeepForestPalette,
  },
  // {
  //   label: '深邃黑',
  //   value: 'white',
  //   palette: whitePalette,
  // },
  {
    label: '电光蓝',
    value: 'electricBlue',
    palette: electricBluePalette,
  },
];

export const THEME_TO_PALETTE = THEME_LIST.reduce(
  (acc, item) => {
    acc[item.value] = {
      value: item.value,
      label: item.label,
      palette: item.palette,
    };
    return acc;
  },
  {} as Record<string, { value: string; label: string; palette: any }>,
);
