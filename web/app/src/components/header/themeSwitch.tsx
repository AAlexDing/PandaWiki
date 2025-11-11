import { IconButton, alpha } from '@mui/material';
import { IconDark, IconLight } from '../icons';
import { useThemeStore } from '@/provider/themeStore';

const ThemeSwitch = () => {
  const { themeMode, setThemeMode } = useThemeStore();
  return (
    <IconButton
      size='small'
      onClick={() => setThemeMode(themeMode === 'dark' ? 'light' : 'dark')}
    >
      {themeMode === 'dark' ? (
        <IconDark
          sx={theme => ({ color: alpha(theme.palette.text.primary, 0.65) })}
        />
      ) : (
        <IconLight sx={{ fontSize: 20 }} />
      )}
    </IconButton>
  );
};

export default ThemeSwitch;
