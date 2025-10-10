// 颜色常量
export const COLORS = {
  // 主色调
  primary: '#1890ff',
  success: '#52c41a',
  warning: '#faad14',
  error: '#ff4d4f',
  info: '#1890ff',
  
  // 状态颜色
  healthy: '#52c41a',
  unhealthy: '#ff4d4f',
  
  // 文本颜色
  text: {
    primary: '#262626',
    secondary: '#8c8c8c',
    disabled: '#bfbfbf',
    inverse: '#ffffff',
  },
  
  // 背景颜色
  background: {
    default: '#ffffff',
    light: '#fafafa',
    dark: '#f5f5f5',
    card: '#ffffff',
    overlay: 'rgba(255, 255, 255, 0.8)',
  },
  
  // 边框颜色
  border: {
    default: '#d9d9d9',
    light: '#f0f0f0',
    dark: '#bfbfbf',
  },
  
  // 特殊颜色
  special: {
    code: '#f5f5f5',
    highlight: '#fff7e6',
    danger: '#fff2f0',
    dangerBorder: '#ffccc7',
  },
} as const;

// 间距常量
export const SPACING = {
  // 基础间距
  xs: '4px',
  sm: '8px',
  md: '12px',
  lg: '16px',
  xl: '20px',
  xxl: '24px',
  xxxl: '32px',
  
  // 特殊间距
  section: '24px',
  page: '40px',
  card: '20px',
  form: '16px',
} as const;

// 字体大小常量
export const FONT_SIZES = {
  xs: '12px',
  sm: '13px',
  md: '14px',
  lg: '15px',
  xl: '16px',
  xxl: '18px',
  xxxl: '20px',
  
  // 标题字体
  h1: '24px',
  h2: '20px',
  h3: '18px',
  h4: '16px',
  h5: '14px',
  h6: '12px',
} as const;

// 字体权重常量
export const FONT_WEIGHTS = {
  normal: 400,
  medium: 500,
  semibold: 600,
  bold: 700,
} as const;

// 边框半径常量
export const BORDER_RADIUS = {
  xs: '2px',
  sm: '4px',
  md: '6px',
  lg: '8px',
  xl: '12px',
  round: '50%',
} as const;

// 阴影常量
export const SHADOWS = {
  sm: '0 1px 2px rgba(0, 0, 0, 0.05)',
  md: '0 1px 3px rgba(0, 0, 0, 0.1), 0 1px 2px rgba(0, 0, 0, 0.06)',
  lg: '0 4px 6px rgba(0, 0, 0, 0.07), 0 2px 4px rgba(0, 0, 0, 0.06)',
  xl: '0 10px 15px rgba(0, 0, 0, 0.1), 0 4px 6px rgba(0, 0, 0, 0.05)',
  card: '0 1px 3px rgba(0, 0, 0, 0.12), 0 1px 2px rgba(0, 0, 0, 0.24)',
} as const;

// Z-index 常量
export const Z_INDEX = {
  dropdown: 1000,
  sticky: 1020,
  fixed: 1030,
  modal: 1040,
  popover: 1050,
  tooltip: 1060,
  toast: 1070,
} as const;

// 断点常量
export const BREAKPOINTS = {
  xs: '480px',
  sm: '576px',
  md: '768px',
  lg: '992px',
  xl: '1200px',
  xxl: '1600px',
} as const;

// 动画时长常量
export const ANIMATION_DURATION = {
  fast: '0.15s',
  normal: '0.3s',
  slow: '0.5s',
} as const;

// 动画缓动函数常量
export const ANIMATION_EASING = {
  ease: 'ease',
  easeIn: 'ease-in',
  easeOut: 'ease-out',
  easeInOut: 'ease-in-out',
  linear: 'linear',
} as const;

// 布局常量
export const LAYOUT = {
  header: {
    height: '64px',
    padding: '0 16px',
  },
  sidebar: {
    width: '200px',
    collapsedWidth: '80px',
  },
  content: {
    padding: '24px',
    minHeight: 'calc(100vh - 64px)',
  },
  footer: {
    height: '48px',
    padding: '12px 16px',
  },
} as const;

// 组件尺寸常量
export const COMPONENT_SIZES = {
  button: {
    small: { height: '24px', padding: '4px 8px', fontSize: '12px' },
    medium: { height: '32px', padding: '8px 16px', fontSize: '14px' },
    large: { height: '40px', padding: '12px 20px', fontSize: '16px' },
  },
  input: {
    small: { height: '24px', fontSize: '12px' },
    medium: { height: '32px', fontSize: '14px' },
    large: { height: '40px', fontSize: '16px' },
  },
  card: {
    padding: '20px',
    borderRadius: '8px',
    border: '1px solid #f0f0f0',
  },
} as const;