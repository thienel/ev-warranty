import type { ConfigProviderProps } from 'antd/es/config-provider'

export const antdThemeConfig: ConfigProviderProps['theme'] = {
  token: {
    // Color tokens
    colorPrimary: '#697565',
    colorSuccess: '#4CAF50',
    colorWarning: '#FF9800',
    colorError: '#F44336',
    colorInfo: '#2196F3',

    // Background colors
    colorBgBase: '#fafafa',
    colorBgContainer: '#ffffff',
    colorBgElevated: '#ffffff',
    colorBgLayout: '#f5f5f5',
    colorBgSpotlight: '#ffffff',
    colorBgMask: 'rgba(60, 61, 55, 0.45)',

    // Text colors
    colorText: '#2c2d2a',
    colorTextSecondary: '#5c6659',
    colorTextTertiary: '#8b9788',
    colorTextQuaternary: '#bbc5b8',
    colorTextDescription: '#6b7866',
    colorTextHeading: '#181c14',
    colorTextLabel: '#3c3d37',
    colorTextPlaceholder: '#a8b5a5',
    colorTextDisabled: '#d1dace',

    // Border colors
    colorBorder: '#e0e6dd',
    colorBorderSecondary: '#eef2ea',
    colorBorderBg: '#f5f7f3',

    // Fill colors
    colorFill: '#f8faf7',
    colorFillSecondary: '#f0f4ed',
    colorFillTertiary: '#e8ede4',
    colorFillQuaternary: '#dfe6db',

    // Border radius
    borderRadius: 8,
    borderRadiusLG: 12,
    borderRadiusSM: 6,
    borderRadiusXS: 4,

    // Font
    fontFamily:
      "-apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, 'Noto Sans', sans-serif",
    fontSize: 14,
    fontSizeSM: 12,
    fontSizeLG: 16,
    fontSizeXL: 20,
    fontSizeHeading1: 32,
    fontSizeHeading2: 26,
    fontSizeHeading3: 22,
    fontSizeHeading4: 18,

    // Font weight
    fontWeightStrong: 600,

    // Control heights
    controlHeight: 36,
    controlHeightLG: 44,
    controlHeightSM: 28,

    // Layout
    padding: 16,
    paddingLG: 24,
    paddingSM: 12,
    paddingXS: 8,
    paddingXXS: 4,

    // Motion
    motionDurationFast: '0.15s',
    motionDurationMid: '0.25s',
    motionDurationSlow: '0.35s',

    // Box shadow
    boxShadow:
      '0 2px 8px 0 rgba(60, 61, 55, 0.06), 0 1px 4px -1px rgba(60, 61, 55, 0.04), 0 4px 8px 0 rgba(60, 61, 55, 0.04)',
    boxShadowSecondary:
      '0 8px 24px 0 rgba(60, 61, 55, 0.12), 0 4px 8px -4px rgba(60, 61, 55, 0.16), 0 12px 32px 8px rgba(60, 61, 55, 0.08)',
  },

  components: {
    Layout: {
      headerBg: '#ffffff',
      headerHeight: 64,
      bodyBg: '#f5f5f5',
      siderBg: '#ffffff',
    },

    Button: {
      primaryColor: '#ffffff',
      defaultBg: '#ffffff',
      defaultBorderColor: '#e0e6dd',
      defaultColor: '#2c2d2a',
    },

    Input: {
      activeBorderColor: '#697565',
      hoverBorderColor: '#9fb09a',
    },

    Select: {
      optionActiveBg: '#f5f7f3',
      optionSelectedBg: '#eef2ea',
    },

    Menu: {
      itemBg: 'transparent',
      itemSelectedBg: '#f5f7f3',
      itemSelectedColor: '#697565',
      itemHoverBg: '#eef2ea',
      itemActiveBg: '#eef2ea',
    },

    Card: {
      colorBgContainer: '#ffffff',
      boxShadow:
        '0 2px 8px 0 rgba(60, 61, 55, 0.06), 0 1px 4px -1px rgba(60, 61, 55, 0.04), 0 4px 8px 0 rgba(60, 61, 55, 0.04)',
    },

    Table: {
      headerBg: '#f8faf7',
      rowHoverBg: '#f5f7f3',
    },

    Message: {
      contentBg: '#ffffff',
    },

    Notification: {
      colorBgElevated: '#ffffff',
    },
  },
}
