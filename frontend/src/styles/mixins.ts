import type { CSSProperties } from 'react';
import { COLORS, SPACING, FONT_SIZES, FONT_WEIGHTS, BORDER_RADIUS, SHADOWS, ANIMATION_DURATION, ANIMATION_EASING } from './constants';

// 基础样式混合
export const flexCenter = (): CSSProperties => ({
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
});

export const flexBetween = (): CSSProperties => ({
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'space-between',
});

export const flexStart = (): CSSProperties => ({
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'flex-start',
});

export const flexEnd = (): CSSProperties => ({
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'flex-end',
});

export const flexColumn = (): CSSProperties => ({
  display: 'flex',
  flexDirection: 'column',
});

export const flexColumnCenter = (): CSSProperties => ({
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'center',
  justifyContent: 'center',
});

// 文本样式混合
export const textEllipsis = (): CSSProperties => ({
  whiteSpace: 'nowrap',
  overflow: 'hidden',
  textOverflow: 'ellipsis',
});

export const textCenter = (): CSSProperties => ({
  textAlign: 'center',
});

export const textBold = (): CSSProperties => ({
  fontWeight: FONT_WEIGHTS.bold,
});

export const textMedium = (): CSSProperties => ({
  fontWeight: FONT_WEIGHTS.medium,
});

export const textSecondary = (): CSSProperties => ({
  color: COLORS.text.secondary,
});

export const textPrimary = (): CSSProperties => ({
  color: COLORS.text.primary,
});

// 卡片样式混合
export const cardStyle = (padding?: string): CSSProperties => ({
  backgroundColor: COLORS.background.card,
  borderRadius: BORDER_RADIUS.lg,
  border: `1px solid ${COLORS.border.light}`,
  padding: padding || SPACING.card,
  boxShadow: SHADOWS.card,
});

export const innerCardStyle = (): CSSProperties => ({
  backgroundColor: COLORS.background.default,
  borderRadius: BORDER_RADIUS.md,
  border: `1px solid ${COLORS.border.light}`,
  padding: SPACING.lg,
});

// 按钮样式混合
export const buttonBase = (): CSSProperties => ({
  border: 'none',
  borderRadius: BORDER_RADIUS.sm,
  cursor: 'pointer',
  transition: `all ${ANIMATION_DURATION.normal} ${ANIMATION_EASING.easeInOut}`,
  fontWeight: FONT_WEIGHTS.medium,
});

export const primaryButton = (): CSSProperties => ({
  ...buttonBase(),
  backgroundColor: COLORS.primary,
  color: COLORS.text.inverse,
});

export const secondaryButton = (): CSSProperties => ({
  ...buttonBase(),
  backgroundColor: COLORS.text.secondary,
  color: COLORS.text.inverse,
});

export const dangerButton = (): CSSProperties => ({
  ...buttonBase(),
  backgroundColor: COLORS.error,
  color: COLORS.text.inverse,
});

// 输入框样式混合
export const inputBase = (): CSSProperties => ({
  width: '100%',
  padding: SPACING.sm,
  border: `1px solid ${COLORS.border.default}`,
  borderRadius: BORDER_RADIUS.sm,
  fontSize: FONT_SIZES.md,
  transition: `border-color ${ANIMATION_DURATION.normal} ${ANIMATION_EASING.easeInOut}`,
});

export const inputFocus = (): CSSProperties => ({
  borderColor: COLORS.primary,
  outline: 'none',
  boxShadow: `0 0 0 2px ${COLORS.primary}20`,
});

// 代码块样式混合
export const codeBlock = (): CSSProperties => ({
  fontSize: FONT_SIZES.xs,
  backgroundColor: COLORS.special.code,
  padding: SPACING.sm,
  borderRadius: BORDER_RADIUS.sm,
  fontFamily: 'monospace',
  whiteSpace: 'pre-wrap',
  overflow: 'auto',
});

export const inlineCode = (): CSSProperties => ({
  fontSize: FONT_SIZES.xs,
  backgroundColor: COLORS.special.code,
  padding: '2px 4px',
  borderRadius: BORDER_RADIUS.xs,
  fontFamily: 'monospace',
});

// 状态样式混合
export const successState = (): CSSProperties => ({
  color: COLORS.success,
});

export const errorState = (): CSSProperties => ({
  color: COLORS.error,
});

export const warningState = (): CSSProperties => ({
  color: COLORS.warning,
});

export const infoState = (): CSSProperties => ({
  color: COLORS.info,
});

// 加载状态样式混合
export const loadingOverlay = (): CSSProperties => ({
  position: 'absolute',
  top: 0,
  left: 0,
  right: 0,
  bottom: 0,
  backgroundColor: COLORS.background.overlay,
  ...flexCenter(),
  zIndex: 1000,
});

export const loadingSpinner = (size: number = 24, color: string = COLORS.primary): CSSProperties => ({
  width: `${size}px`,
  height: `${size}px`,
  border: `2px solid ${color}20`,
  borderTop: `2px solid ${color}`,
  borderRadius: BORDER_RADIUS.round,
  animation: 'spin 1s linear infinite',
});

// 空状态样式混合
export const emptyState = (): CSSProperties => ({
  ...flexColumnCenter(),
  padding: SPACING.page,
  color: COLORS.text.secondary,
  backgroundColor: COLORS.background.light,
  borderRadius: BORDER_RADIUS.md,
  border: `1px dashed ${COLORS.border.default}`,
});

// 错误状态样式混合
export const errorContainer = (): CSSProperties => ({
  padding: SPACING.lg,
  backgroundColor: COLORS.special.danger,
  border: `1px solid ${COLORS.special.dangerBorder}`,
  borderRadius: BORDER_RADIUS.md,
  color: COLORS.error,
});

// 高亮样式混合
export const highlightContainer = (): CSSProperties => ({
  backgroundColor: COLORS.special.highlight,
  border: `1px solid ${COLORS.warning}40`,
  borderRadius: BORDER_RADIUS.md,
  padding: SPACING.lg,
});

// 响应式样式辅助函数（需要在CSS-in-JS库中使用）
export const getResponsiveStyles = () => ({
  mobile: '@media (max-width: 768px)',
  tablet: '@media (min-width: 769px) and (max-width: 1024px)',
  desktop: '@media (min-width: 1025px)',
});

// 滚动条样式（需要在全局CSS中定义）
export const getScrollbarStyles = (width: string = '6px') => ({
  width,
  trackBackground: COLORS.background.light,
  thumbBackground: COLORS.border.dark,
  thumbHoverBackground: COLORS.text.secondary,
  borderRadius: BORDER_RADIUS.xs,
});

// 动画样式混合
export const fadeIn = (duration: string = ANIMATION_DURATION.normal): CSSProperties => ({
  animation: `fadeIn ${duration} ${ANIMATION_EASING.easeInOut}`,
});

export const slideIn = (direction: 'left' | 'right' | 'up' | 'down' = 'up', duration: string = ANIMATION_DURATION.normal): CSSProperties => ({
  animation: `slideIn${direction.charAt(0).toUpperCase() + direction.slice(1)} ${duration} ${ANIMATION_EASING.easeInOut}`,
});

// 工具样式混合
export const absoluteCenter = (): CSSProperties => ({
  position: 'absolute',
  top: '50%',
  left: '50%',
  transform: 'translate(-50%, -50%)',
});

export const fullSize = (): CSSProperties => ({
  width: '100%',
  height: '100%',
});

export const visuallyHidden = (): CSSProperties => ({
  position: 'absolute',
  width: '1px',
  height: '1px',
  padding: 0,
  margin: '-1px',
  overflow: 'hidden',
  clip: 'rect(0, 0, 0, 0)',
  whiteSpace: 'nowrap',
  border: 0,
});

// 组合样式混合
export const pageContainer = (): CSSProperties => ({
  padding: SPACING.section,
  maxWidth: '1200px',
  margin: '0 auto',
});

export const sectionHeader = (): CSSProperties => ({
  ...flexBetween(),
  marginBottom: SPACING.section,
});

export const formSection = (): CSSProperties => ({
  marginBottom: SPACING.xl,
});

export const tableHeader = (): CSSProperties => ({
  ...flexBetween(),
  marginBottom: SPACING.lg,
});

export const modalContent = (): CSSProperties => ({
  padding: SPACING.xl,
});

export const cardHeader = (): CSSProperties => ({
  ...flexBetween(),
  marginBottom: SPACING.lg,
  paddingBottom: SPACING.md,
  borderBottom: `1px solid ${COLORS.border.light}`,
});

export const cardBody = (): CSSProperties => ({
  padding: SPACING.lg,
});

export const cardFooter = (): CSSProperties => ({
  ...flexEnd(),
  paddingTop: SPACING.md,
  borderTop: `1px solid ${COLORS.border.light}`,
  gap: SPACING.sm,
});