const antdTheme = {
  token: {
    // === MÀUS CHÍNH ===
    colorPrimary: '#697565', // Sage Green - màu chính dễ nhìn hơn
    colorPrimaryHover: '#5a6358', // Hover đậm hơn một chút
    colorPrimaryActive: '#4d5549', // Active đậm hơn
    colorPrimaryBg: '#f5f7f3', // Background rất nhạt của primary
    colorPrimaryBgHover: '#eef2ea', // Hover background nhạt
    colorPrimaryBorder: '#b8c4b3', // Border nhạt hơn
    colorPrimaryBorderHover: '#9fb09a', // Hover border
    colorPrimaryText: '#697565', // Text primary
    colorPrimaryTextHover: '#5a6358', // Hover text primary
    colorPrimaryTextActive: '#4d5549', // Active text primary

    // === MÀUS THÀNH CÔNG ===
    colorSuccess: '#4CAF50', // Xanh lá tự nhiên
    colorSuccessBg: '#f3f8f3',
    colorSuccessBorder: '#c8e6c9',

    // === MÀUS CẢNH BÁO ===
    colorWarning: '#FF9800', // Cam ấm
    colorWarningBg: '#fff8f0',
    colorWarningBorder: '#ffcc80',

    // === MÀUS LỖI ===
    colorError: '#F44336', // Đỏ mềm mại hơn
    colorErrorBg: '#fdf2f2',
    colorErrorBorder: '#ffcdd2',

    // === MÀUS THÔNG TIN ===
    colorInfo: '#2196F3', // Xanh dương nhẹ nhàng
    colorInfoBg: '#f3f8ff',
    colorInfoBorder: '#bbdefb',

    // === MÀUS NỀN ===
    colorBgBase: '#fafafa', // Background trắng xám nhẹ thay vì cream
    colorBgContainer: '#ffffff', // Background container trắng tinh
    colorBgElevated: '#ffffff', // Background elevated
    colorBgLayout: '#f5f5f5', // Background layout xám rất nhạt
    colorBgSpotlight: '#ffffff', // Background spotlight
    colorBgMask: 'rgba(60, 61, 55, 0.45)', // Mask overlay từ charcoal

    // === MÀUS TEXT ===
    colorText: '#2c2d2a', // Text chính đậm hơn một chút từ charcoal
    colorTextSecondary: '#5c6659', // Text phụ từ sage green đậm hơn
    colorTextTertiary: '#8b9788', // Text tertiary
    colorTextQuaternary: '#bbc5b8', // Text quaternary
    colorTextDescription: '#6b7866', // Text mô tả
    colorTextHeading: '#181c14', // Text heading giữ dark forest
    colorTextLabel: '#3c3d37', // Text label từ charcoal
    colorTextPlaceholder: '#a8b5a5', // Placeholder text
    colorTextDisabled: '#d1dace', // Text disabled

    // === MÀUS BORDER ===
    colorBorder: '#e0e6dd', // Border chính nhạt hơn
    colorBorderSecondary: '#eef2ea', // Border phụ rất nhạt
    colorBorderBg: '#f5f7f3', // Border background

    // === MÀUS FILL ===
    colorFill: '#f8faf7', // Fill nhạt nhất
    colorFillSecondary: '#f0f4ed', // Fill nhạt
    colorFillTertiary: '#e8ede4', // Fill trung bình
    colorFillQuaternary: '#dfe6db', // Fill đậm nhất

    // === KÍCH THƯỚC ===
    borderRadius: 8, // Border radius tăng lên cho modern hơn
    borderRadiusLG: 12, // Border radius lớn
    borderRadiusSM: 6, // Border radius nhỏ
    borderRadiusXS: 4, // Border radius rất nhỏ

    // === KHOẢNG CÁCH ===
    padding: 16, // Padding cơ bản
    paddingLG: 24, // Padding lớn
    paddingSM: 12, // Padding nhỏ
    paddingXS: 8, // Padding rất nhỏ
    paddingXXS: 4, // Padding cực nhỏ
    margin: 16, // Margin cơ bản
    marginLG: 24, // Margin lớn
    marginSM: 12, // Margin nhỏ
    marginXS: 8, // Margin rất nhỏ
    marginXXS: 4, // Margin cực nhỏ

    // === FONT ===
    fontFamily:
      '-apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, "Noto Sans", sans-serif',
    fontSize: 14, // Font size cơ bản
    fontSizeLG: 16, // Font size lớn
    fontSizeSM: 12, // Font size nhỏ
    fontSizeXL: 20, // Font size rất lớn
    fontSizeHeading1: 32, // H1 giảm xuống cho cân bằng hơn
    fontSizeHeading2: 26, // H2
    fontSizeHeading3: 22, // H3
    fontSizeHeading4: 18, // H4
    fontSizeHeading5: 16, // H5
    fontWeightStrong: 600, // Font weight đậm

    // === CHIỀU CAO ===
    controlHeight: 36, // Tăng chiều cao cho dễ tương tác hơn
    controlHeightLG: 44, // Chiều cao control lớn
    controlHeightSM: 28, // Chiều cao control nhỏ

    // === LINE HEIGHT ===
    lineHeight: 1.5714285714285714, // Line height cơ bản
    lineHeightLG: 1.5, // Line height lớn
    lineHeightSM: 1.6666666666666667, // Line height nhỏ

    // === BOX SHADOW ===
    boxShadow:
      '0 2px 8px 0 rgba(60, 61, 55, 0.06), 0 1px 4px -1px rgba(60, 61, 55, 0.04), 0 4px 8px 0 rgba(60, 61, 55, 0.04)',
    boxShadowSecondary:
      '0 8px 24px 0 rgba(60, 61, 55, 0.12), 0 4px 8px -4px rgba(60, 61, 55, 0.16), 0 12px 32px 8px rgba(60, 61, 55, 0.08)',

    // === WIRE FRAME (Tùy chọn) ===
    wireframe: false, // Bật/tắt wireframe mode

    // === MOTION ===
    motionDurationFast: '0.15s', // Animation tăng lên một chút
    motionDurationMid: '0.25s', // Animation trung bình
    motionDurationSlow: '0.35s', // Animation chậm
  },

  // === ALGORITHM (Tùy chọn) ===
  algorithm: [], // Có thể sử dụng theme.darkAlgorithm, theme.compactAlgorithm

  // === COMPONENT CUSTOMIZATION ===
  components: {
    // Button customization
    Button: {
      colorPrimary: '#697565',
      colorPrimaryHover: '#5a6358',
      colorPrimaryActive: '#4d5549',
      borderRadius: 8,
      algorithm: true,
    },

    // Menu customization
    Menu: {
      itemBg: 'transparent',
      itemSelectedBg: '#f5f7f3',
      itemHoverBg: '#f8faf7',
      itemSelectedColor: '#697565',
      itemActiveBg: '#eef2ea',
      borderRadius: 6,
    },

    // Layout customization
    Layout: {
      headerBg: '#ffffff',
      bodyBg: '#fafafa', // Thay đổi từ cream sang xám nhạt
      siderBg: '#181C14',
      footerBg: '#f5f5f5',
    },

    // Table customization
    Table: {
      headerBg: '#f8faf7',
      headerColor: '#2c2d2a',
      rowHoverBg: '#f5f7f3',
      borderColor: '#e0e6dd',
      colorBgContainer: '#ffffff',
    },

    // Input customization
    Input: {
      colorBorder: '#e0e6dd',
      colorBorderHover: '#b8c4b3',
      colorBgContainer: '#ffffff',
      borderRadius: 8,
      controlHeight: 36,
    },

    // Card customization
    Card: {
      colorBorderSecondary: '#eef2ea',
      colorBgContainer: '#ffffff',
      borderRadiusLG: 12,
      boxShadowTertiary: '0 2px 8px 0 rgba(60, 61, 55, 0.06)',
    },

    // Typography customization
    Typography: {
      colorText: '#2c2d2a',
      colorTextHeading: '#181c14',
      colorTextDescription: '#6b7866',
      colorTextSecondary: '#5c6659',
    },

    // Notification customization
    Notification: {
      colorBgElevated: '#ffffff',
      borderRadiusLG: 12,
    },

    // Modal customization
    Modal: {
      colorBgElevated: '#ffffff',
      borderRadiusLG: 16,
    },

    // Dropdown customization
    Dropdown: {
      colorBgElevated: '#ffffff',
      borderRadiusOuter: 8,
      boxShadowSecondary:
        '0 8px 24px 0 rgba(60, 61, 55, 0.12), 0 4px 8px -4px rgba(60, 61, 55, 0.16)',
    },

    // Switch customization
    Switch: {
      colorPrimary: '#697565',
      colorPrimaryHover: '#5a6358',
    },

    // Progress customization
    Progress: {
      colorInfo: '#697565',
      remainingColor: '#f0f4ed',
    },

    // Tag customization
    Tag: {
      colorFillSecondary: '#f8faf7',
      colorBorderSecondary: '#eef2ea',
      borderRadiusSM: 6,
    },

    // Divider customization
    Divider: {
      colorSplit: '#e0e6dd',
    },
  },
}

export default antdTheme
