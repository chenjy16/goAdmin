import type { CSSProperties } from 'react';
import { SPACING, COLORS, BORDER_RADIUS, SHADOWS } from './constants';

// 间距工具函数
export const spacing = {
  // 获取间距值
  get: (size: keyof typeof SPACING) => SPACING[size],
  
  // 生成padding样式
  padding: (size: keyof typeof SPACING) => ({ padding: SPACING[size] }),
  paddingX: (size: keyof typeof SPACING) => ({ 
    paddingLeft: SPACING[size], 
    paddingRight: SPACING[size] 
  }),
  paddingY: (size: keyof typeof SPACING) => ({ 
    paddingTop: SPACING[size], 
    paddingBottom: SPACING[size] 
  }),
  paddingTop: (size: keyof typeof SPACING) => ({ paddingTop: SPACING[size] }),
  paddingBottom: (size: keyof typeof SPACING) => ({ paddingBottom: SPACING[size] }),
  paddingLeft: (size: keyof typeof SPACING) => ({ paddingLeft: SPACING[size] }),
  paddingRight: (size: keyof typeof SPACING) => ({ paddingRight: SPACING[size] }),
  
  // 生成margin样式
  margin: (size: keyof typeof SPACING) => ({ margin: SPACING[size] }),
  marginX: (size: keyof typeof SPACING) => ({ 
    marginLeft: SPACING[size], 
    marginRight: SPACING[size] 
  }),
  marginY: (size: keyof typeof SPACING) => ({ 
    marginTop: SPACING[size], 
    marginBottom: SPACING[size] 
  }),
  marginTop: (size: keyof typeof SPACING) => ({ marginTop: SPACING[size] }),
  marginBottom: (size: keyof typeof SPACING) => ({ marginBottom: SPACING[size] }),
  marginLeft: (size: keyof typeof SPACING) => ({ marginLeft: SPACING[size] }),
  marginRight: (size: keyof typeof SPACING) => ({ marginRight: SPACING[size] }),
};

// 颜色工具函数
export const colors = {
  // 获取颜色值
  primary: () => COLORS.primary,
  success: () => COLORS.success,
  warning: () => COLORS.warning,
  error: () => COLORS.error,
  info: () => COLORS.info,
  healthy: () => COLORS.healthy,
  unhealthy: () => COLORS.unhealthy,
  
  // 文本颜色
  text: (type: 'primary' | 'secondary' | 'disabled' | 'inverse') => COLORS.text[type],
  
  // 背景颜色
  background: (type: 'default' | 'light' | 'dark' | 'card' | 'overlay') => COLORS.background[type],
  
  // 边框颜色
  border: (type: 'default' | 'light' | 'dark') => COLORS.border[type],
  
  // 特殊颜色
  special: (type: 'code' | 'highlight' | 'danger' | 'dangerBorder') => COLORS.special[type],
};

// 阴影工具函数
export const shadows = {
  get: (level: keyof typeof SHADOWS) => ({ boxShadow: SHADOWS[level] }),
  none: () => ({ boxShadow: 'none' }),
  hover: () => ({ boxShadow: SHADOWS.md }),
  focus: () => ({ boxShadow: `0 0 0 2px ${COLORS.primary}` }),
};

// 边框半径工具函数
export const borderRadius = {
  get: (size: keyof typeof BORDER_RADIUS) => ({ borderRadius: BORDER_RADIUS[size] }),
  top: (size: keyof typeof BORDER_RADIUS) => ({ 
    borderTopLeftRadius: BORDER_RADIUS[size],
    borderTopRightRadius: BORDER_RADIUS[size]
  }),
  bottom: (size: keyof typeof BORDER_RADIUS) => ({ 
    borderBottomLeftRadius: BORDER_RADIUS[size],
    borderBottomRightRadius: BORDER_RADIUS[size]
  }),
  left: (size: keyof typeof BORDER_RADIUS) => ({ 
    borderTopLeftRadius: BORDER_RADIUS[size],
    borderBottomLeftRadius: BORDER_RADIUS[size]
  }),
  right: (size: keyof typeof BORDER_RADIUS) => ({ 
    borderTopRightRadius: BORDER_RADIUS[size],
    borderBottomRightRadius: BORDER_RADIUS[size]
  }),
};

// 布局工具函数
export const layout = {
  // Flexbox 工具
  flex: (direction?: 'row' | 'column', align?: string, justify?: string): CSSProperties => ({
    display: 'flex',
    flexDirection: direction || 'row',
    alignItems: align || 'center',
    justifyContent: justify || 'flex-start',
  }),
  
  flexCenter: (): CSSProperties => ({
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
  }),
  
  flexBetween: (): CSSProperties => ({
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'space-between',
  }),
  
  flexColumn: (): CSSProperties => ({
    display: 'flex',
    flexDirection: 'column',
  }),
  
  // Grid 工具
  grid: (columns?: number, gap?: keyof typeof SPACING): CSSProperties => ({
    display: 'grid',
    gridTemplateColumns: columns ? `repeat(${columns}, 1fr)` : 'auto',
    gap: gap ? SPACING[gap] : SPACING.md,
  }),
  
  // 定位工具
  absolute: (top?: string, right?: string, bottom?: string, left?: string): CSSProperties => ({
    position: 'absolute',
    ...(top && { top }),
    ...(right && { right }),
    ...(bottom && { bottom }),
    ...(left && { left }),
  }),
  
  relative: (): CSSProperties => ({
    position: 'relative',
  }),
  
  fixed: (top?: string, right?: string, bottom?: string, left?: string): CSSProperties => ({
    position: 'fixed',
    ...(top && { top }),
    ...(right && { right }),
    ...(bottom && { bottom }),
    ...(left && { left }),
  }),
};

// 文本工具函数
export const text = {
  // 文本截断
  ellipsis: (lines?: number): CSSProperties => {
    if (lines && lines > 1) {
      return {
        display: '-webkit-box',
        WebkitLineClamp: lines,
        WebkitBoxOrient: 'vertical',
        overflow: 'hidden',
        textOverflow: 'ellipsis',
      };
    }
    return {
      overflow: 'hidden',
      textOverflow: 'ellipsis',
      whiteSpace: 'nowrap',
    };
  },
  
  // 文本对齐
  align: (alignment: 'left' | 'center' | 'right' | 'justify'): CSSProperties => ({
    textAlign: alignment,
  }),
  
  // 文本大小写
  transform: (transform: 'uppercase' | 'lowercase' | 'capitalize' | 'none'): CSSProperties => ({
    textTransform: transform,
  }),
  
  // 行高
  lineHeight: (height: number): CSSProperties => ({
    lineHeight: height,
  }),
};

// 尺寸工具函数
export const size = {
  // 宽度
  width: (value: string | number): CSSProperties => ({
    width: typeof value === 'number' ? `${value}px` : value,
  }),
  
  // 高度
  height: (value: string | number): CSSProperties => ({
    height: typeof value === 'number' ? `${value}px` : value,
  }),
  
  // 最小宽度
  minWidth: (value: string | number): CSSProperties => ({
    minWidth: typeof value === 'number' ? `${value}px` : value,
  }),
  
  // 最小高度
  minHeight: (value: string | number): CSSProperties => ({
    minHeight: typeof value === 'number' ? `${value}px` : value,
  }),
  
  // 最大宽度
  maxWidth: (value: string | number): CSSProperties => ({
    maxWidth: typeof value === 'number' ? `${value}px` : value,
  }),
  
  // 最大高度
  maxHeight: (value: string | number): CSSProperties => ({
    maxHeight: typeof value === 'number' ? `${value}px` : value,
  }),
  
  // 正方形
  square: (value: string | number): CSSProperties => {
    const sizeValue = typeof value === 'number' ? `${value}px` : value;
    return {
      width: sizeValue,
      height: sizeValue,
    };
  },
  
  // 圆形
  circle: (value: string | number): CSSProperties => {
    const sizeValue = typeof value === 'number' ? `${value}px` : value;
    return {
      width: sizeValue,
      height: sizeValue,
      borderRadius: '50%',
    };
  },
};

// 动画工具函数
export const animation = {
  // 过渡效果
  transition: (property?: string, duration?: string, easing?: string): CSSProperties => ({
    transition: `${property || 'all'} ${duration || '0.3s'} ${easing || 'ease'}`,
  }),
  
  // 变换
  transform: (value: string): CSSProperties => ({
    transform: value,
  }),
  
  // 旋转
  rotate: (degrees: number): CSSProperties => ({
    transform: `rotate(${degrees}deg)`,
  }),
  
  // 缩放
  scale: (value: number): CSSProperties => ({
    transform: `scale(${value})`,
  }),
  
  // 平移
  translate: (x: string | number, y?: string | number): CSSProperties => {
    const xValue = typeof x === 'number' ? `${x}px` : x;
    const yValue = y ? (typeof y === 'number' ? `${y}px` : y) : '0';
    return {
      transform: `translate(${xValue}, ${yValue})`,
    };
  },
};

// 组合样式工具函数
export const combine = (...styles: (CSSProperties | undefined)[]): CSSProperties => {
  return styles.filter(Boolean).reduce((acc, style) => ({ ...acc, ...style }), {});
};

// 条件样式工具函数
export const conditional = (condition: boolean, trueStyle: CSSProperties, falseStyle?: CSSProperties): CSSProperties => {
  return condition ? trueStyle : (falseStyle || {});
};

// 主题工具函数
export const theme = {
  // 获取主题相关的样式
  getButtonStyle: (variant: 'primary' | 'success' | 'warning' | 'error' = 'primary'): CSSProperties => ({
    backgroundColor: colors[variant](),
    color: '#ffffff',
    border: 'none',
    ...borderRadius.get('sm'),
    ...spacing.padding('sm'),
    cursor: 'pointer',
    ...animation.transition('all', '0.2s'),
  }),
  
  getCardStyle: (): CSSProperties => ({
    backgroundColor: colors.background('card'),
    ...borderRadius.get('md'),
    ...shadows.get('sm'),
    ...spacing.padding('lg'),
  }),
  
  getInputStyle: (): CSSProperties => ({
    border: `1px solid ${colors.border('default')}`,
    ...borderRadius.get('sm'),
    ...spacing.padding('sm'),
    fontSize: '14px',
    ...animation.transition('border-color', '0.2s'),
  }),
};