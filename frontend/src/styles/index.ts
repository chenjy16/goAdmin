// 样式常量
export * from './constants';

// 样式混合函数
export * from './mixins';

// 样式工具函数
export * from './utils';

// 重新导出常用的样式工具，提供更简洁的API
export {
  COLORS,
  SPACING,
  FONT_SIZES,
  FONT_WEIGHTS,
  BORDER_RADIUS,
  SHADOWS,
  Z_INDEX,
  BREAKPOINTS,
  ANIMATION_DURATION,
  ANIMATION_EASING,
  LAYOUT,
  COMPONENT_SIZES,
} from './constants';

export {
  spacing,
  colors,
  shadows,
  borderRadius,
  layout,
  text,
  size,
  animation,
  combine,
  conditional,
  theme,
} from './utils';

export {
  flexCenter,
  textEllipsis,
  cardStyle,
  primaryButton,
  secondaryButton,
  inputBase,
  codeBlock,
  successState,
  warningState,
  errorState,
  loadingOverlay,
  emptyState,
  errorContainer,
  highlightContainer,
  fadeIn,
  slideIn,
  absoluteCenter,
  getResponsiveStyles,
  getScrollbarStyles,
} from './mixins';