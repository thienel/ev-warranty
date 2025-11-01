-- =========================================================
-- * EV Warranty Seed Data Script
-- * This script is idempotent - safe to run multiple times
-- =========================================================

-- Set required SQL Server options for indexed views and computed columns
SET QUOTED_IDENTIFIER ON;
SET ANSI_NULLS ON;
SET ANSI_PADDING ON;
SET ANSI_WARNINGS ON;
SET CONCAT_NULL_YIELDS_NULL ON;
SET ARITHABORT ON;

-- Check if seed data already exists
IF EXISTS (SELECT 1 FROM vehicle_models WHERE id = '650e8400-e29b-41d4-a716-446655440101')
BEGIN
    PRINT 'Seed data already exists. Skipping insertion.';
    RETURN;
END

BEGIN TRANSACTION;

INSERT INTO warranty_policies (id, policy_name, warranty_duration_months, kilometer_limit, terms_and_conditions, created_at)
VALUES
    -- VinFast Models
    ('050e8400-e29b-41d4-a716-446655440001', 'VF8_2023_WARRANTY', 120, 200000, N'Bảo hành 10 năm hoặc 200,000 km (tùy điều kiện nào đến trước). Pin và động cơ điện được bảo hành trọn đời theo điều kiện sử dụng. Bảo dưỡng định kỳ theo khuyến cáo nhà sản xuất. Không áp dụng cho xe sử dụng vào mục đích thương mại.', GETDATE()),
    ('050e8400-e29b-41d4-a716-446655440002', 'VF8_2024_WARRANTY', 120, 200000, N'Bảo hành 10 năm hoặc 200,000 km (tùy điều kiện nào đến trước). Pin và động cơ điện được bảo hành trọn đời theo điều kiện sử dụng. Bảo dưỡng định kỳ theo khuyến cáo nhà sản xuất. Không áp dụng cho xe sử dụng vào mục đích thương mại.', GETDATE()),
    ('050e8400-e29b-41d4-a716-446655440003', 'VF8_2025_WARRANTY', 120, 200000, N'Bảo hành 10 năm hoặc 200,000 km (tùy điều kiện nào đến trước). Pin và động cơ điện được bảo hành trọn đời theo điều kiện sử dụng. Bảo dưỡng định kỳ theo khuyến cáo nhà sản xuất. Không áp dụng cho xe sử dụng vào mục đích thương mại.', GETDATE()),
    ('050e8400-e29b-41d4-a716-446655440004', 'VF9_2023_WARRANTY', 120, 200000, N'Bảo hành 10 năm hoặc 200,000 km (tùy điều kiện nào đến trước). Pin và động cơ điện được bảo hành trọn đời theo điều kiện sử dụng. Bảo dưỡng định kỳ theo khuyến cáo nhà sản xuất. Không áp dụng cho xe sử dụng vào mục đích thương mại.', GETDATE()),
    ('050e8400-e29b-41d4-a716-446655440005', 'VF9_2024_WARRANTY', 120, 200000, N'Bảo hành 10 năm hoặc 200,000 km (tùy điều kiện nào đến trước). Pin và động cơ điện được bảo hành trọn đời theo điều kiện sử dụng. Bảo dưỡng định kỳ theo khuyến cáo nhà sản xuất. Không áp dụng cho xe sử dụng vào mục đích thương mại.', GETDATE()),
    ('050e8400-e29b-41d4-a716-446655440006', 'VF9_2025_WARRANTY', 120, 200000, N'Bảo hành 10 năm hoặc 200,000 km (tùy điều kiện nào đến trước). Pin và động cơ điện được bảo hành trọn đời theo điều kiện sử dụng. Bảo dưỡng định kỳ theo khuyến cáo nhà sản xuất. Không áp dụng cho xe sử dụng vào mục đích thương mại.', GETDATE()),
    ('050e8400-e29b-41d4-a716-446655440007', 'VFE34_2022_WARRANTY', 84, 150000, N'Bảo hành 7 năm hoặc 150,000 km (tùy điều kiện nào đến trước). Pin được bảo hành 8 năm hoặc 160,000 km. Bảo dưỡng định kỳ miễn phí 3 năm đầu. Áp dụng cho cả mục đích cá nhân và thương mại.', GETDATE()),
    ('050e8400-e29b-41d4-a716-446655440008', 'VFE34_2023_WARRANTY', 84, 150000, N'Bảo hành 7 năm hoặc 150,000 km (tùy điều kiện nào đến trước). Pin được bảo hành 8 năm hoặc 160,000 km. Bảo dưỡng định kỳ miễn phí 3 năm đầu. Áp dụng cho cả mục đích cá nhân và thương mại.', GETDATE()),
    ('050e8400-e29b-41d4-a716-446655440009', 'VFE34_2024_WARRANTY', 84, 150000, N'Bảo hành 7 năm hoặc 150,000 km (tùy điều kiện nào đến trước). Pin được bảo hành 8 năm hoặc 160,000 km. Bảo dưỡng định kỳ miễn phí 3 năm đầu. Áp dụng cho cả mục đích cá nhân và thương mại.', GETDATE()),
    ('050e8400-e29b-41d4-a716-446655440010', 'VFE34_2025_WARRANTY', 84, 150000, N'Bảo hành 7 năm hoặc 150,000 km (tùy điều kiện nào đến trước). Pin được bảo hành 8 năm hoặc 160,000 km. Bảo dưỡng định kỳ miễn phí 3 năm đầu. Áp dụng cho cả mục đích cá nhân và thương mại.', GETDATE()),
    ('050e8400-e29b-41d4-a716-446655440011', 'VF5_2023_WARRANTY', 120, 200000, N'Bảo hành 10 năm hoặc 200,000 km (tùy điều kiện nào đến trước). Pin được bảo hành trọn đời theo điều kiện sử dụng. Hỗ trợ cứu hộ 24/7 miễn phí trong thời gian bảo hành. Không áp dụng cho xe sử dụng vào mục đích thương mại.', GETDATE()),
    ('050e8400-e29b-41d4-a716-446655440012', 'VF5_2024_WARRANTY', 120, 200000, N'Bảo hành 10 năm hoặc 200,000 km (tùy điều kiện nào đến trước). Pin được bảo hành trọn đời theo điều kiện sử dụng. Hỗ trợ cứu hộ 24/7 miễn phí trong thời gian bảo hành. Không áp dụng cho xe sử dụng vào mục đích thương mại.', GETDATE()),
    ('050e8400-e29b-41d4-a716-446655440013', 'VF5_2025_WARRANTY', 120, 200000, N'Bảo hành 10 năm hoặc 200,000 km (tùy điều kiện nào đến trước). Pin được bảo hành trọn đời theo điều kiện sử dụng. Hỗ trợ cứu hộ 24/7 miễn phí trong thời gian bảo hành. Không áp dụng cho xe sử dụng vào mục đích thương mại.', GETDATE()),
    ('050e8400-e29b-41d4-a716-446655440014', 'VF6_2024_WARRANTY', 120, 200000, N'Bảo hành 10 năm hoặc 200,000 km (tùy điều kiện nào đến trước). Pin được bảo hành trọn đời theo điều kiện sử dụng. Hỗ trợ cứu hộ 24/7 miễn phí trong thời gian bảo hành. Không áp dụng cho xe sử dụng vào mục đích thương mại.', GETDATE()),
    ('050e8400-e29b-41d4-a716-446655440015', 'VF6_2025_WARRANTY', 120, 200000, N'Bảo hành 10 năm hoặc 200,000 km (tùy điều kiện nào đến trước). Pin được bảo hành trọn đời theo điều kiện sử dụng. Hỗ trợ cứu hộ 24/7 miễn phí trong thời gian bảo hành. Không áp dụng cho xe sử dụng vào mục đích thương mại.', GETDATE()),
    ('050e8400-e29b-41d4-a716-446655440016', 'VF7_2024_WARRANTY', 120, 200000, N'Bảo hành 10 năm hoặc 200,000 km (tùy điều kiện nào đến trước). Pin được bảo hành trọn đời theo điều kiện sử dụng. Hỗ trợ cứu hộ 24/7 miễn phí trong thời gian bảo hành. Không áp dụng cho xe sử dụng vào mục đích thương mại.', GETDATE()),
    ('050e8400-e29b-41d4-a716-446655440017', 'VF7_2025_WARRANTY', 120, 200000, N'Bảo hành 10 năm hoặc 200,000 km (tùy điều kiện nào đến trước). Pin được bảo hành trọn đời theo điều kiện sử dụng. Hỗ trợ cứu hộ 24/7 miễn phí trong thời gian bảo hành. Không áp dụng cho xe sử dụng vào mục đích thương mại.', GETDATE()),
    ('050e8400-e29b-41d4-a716-446655440018', 'VF3_2024_WARRANTY', 96, 160000, N'Bảo hành 8 năm hoặc 160,000 km (tùy điều kiện nào đến trước). Pin được bảo hành 10 năm. Bảo dưỡng định kỳ miễn phí 2 năm đầu. Không áp dụng cho xe sử dụng vào mục đích thương mại.', GETDATE()),
    ('050e8400-e29b-41d4-a716-446655440019', 'VF3_2025_WARRANTY', 96, 160000, N'Bảo hành 8 năm hoặc 160,000 km (tùy điều kiện nào đến trước). Pin được bảo hành 10 năm. Bảo dưỡng định kỳ miễn phí 2 năm đầu. Không áp dụng cho xe sử dụng vào mục đích thương mại.', GETDATE()),

    -- BYD Atto 3
    ('050e8400-e29b-41d4-a716-446655440020', 'BYD_ATTO3_2022_WARRANTY', 72, 150000, N'Bảo hành 6 năm hoặc 150,000 km (tùy điều kiện nào đến trước). Pin và hệ thống điện được bảo hành 8 năm hoặc 150,000 km. Bảo dưỡng định kỳ theo khuyến cáo. Áp dụng cho xe sử dụng cá nhân.', GETDATE()),
    ('050e8400-e29b-41d4-a716-446655440021', 'BYD_ATTO3_2023_WARRANTY', 72, 150000, N'Bảo hành 6 năm hoặc 150,000 km (tùy điều kiện nào đến trước). Pin và hệ thống điện được bảo hành 8 năm hoặc 150,000 km. Bảo dưỡng định kỳ theo khuyến cáo. Áp dụng cho xe sử dụng cá nhân.', GETDATE()),
    ('050e8400-e29b-41d4-a716-446655440022', 'BYD_ATTO3_2024_WARRANTY', 72, 150000, N'Bảo hành 6 năm hoặc 150,000 km (tùy điều kiện nào đến trước). Pin và hệ thống điện được bảo hành 8 năm hoặc 150,000 km. Bảo dưỡng định kỳ theo khuyến cáo. Áp dụng cho xe sử dụng cá nhân.', GETDATE()),
    ('050e8400-e29b-41d4-a716-446655440023', 'BYD_ATTO3_2025_WARRANTY', 72, 150000, N'Bảo hành 6 năm hoặc 150,000 km (tùy điều kiện nào đến trước). Pin và hệ thống điện được bảo hành 8 năm hoặc 150,000 km. Bảo dưỡng định kỳ theo khuyến cáo. Áp dụng cho xe sử dụng cá nhân.', GETDATE()),

    -- BYD Dolphin
    ('050e8400-e29b-41d4-a716-446655440024', 'BYD_DOLPHIN_2023_WARRANTY', 72, 150000, N'Bảo hành 6 năm hoặc 150,000 km (tùy điều kiện nào đến trước). Pin và hệ thống điện được bảo hành 8 năm hoặc 150,000 km. Bảo dưỡng định kỳ theo khuyến cáo. Áp dụng cho xe sử dụng cá nhân.', GETDATE()),
    ('050e8400-e29b-41d4-a716-446655440025', 'BYD_DOLPHIN_2024_WARRANTY', 72, 150000, N'Bảo hành 6 năm hoặc 150,000 km (tùy điều kiện nào đến trước). Pin và hệ thống điện được bảo hành 8 năm hoặc 150,000 km. Bảo dưỡng định kỳ theo khuyến cáo. Áp dụng cho xe sử dụng cá nhân.', GETDATE()),
    ('050e8400-e29b-41d4-a716-446655440026', 'BYD_DOLPHIN_2025_WARRANTY', 72, 150000, N'Bảo hành 6 năm hoặc 150,000 km (tùy điều kiện nào đến trước). Pin và hệ thống điện được bảo hành 8 năm hoặc 150,000 km. Bảo dưỡng định kỳ theo khuyến cáo. Áp dụng cho xe sử dụng cá nhân.', GETDATE()),

    -- BYD Seal
    ('050e8400-e29b-41d4-a716-446655440027', 'BYD_SEAL_2023_WARRANTY', 72, 150000, N'Bảo hành 6 năm hoặc 150,000 km (tùy điều kiện nào đến trước). Pin và hệ thống điện được bảo hành 8 năm hoặc 150,000 km. Hỗ trợ cứu hộ 24/7 trong thời gian bảo hành. Áp dụng cho xe sử dụng cá nhân.', GETDATE()),
    ('050e8400-e29b-41d4-a716-446655440028', 'BYD_SEAL_2024_WARRANTY', 72, 150000, N'Bảo hành 6 năm hoặc 150,000 km (tùy điều kiện nào đến trước). Pin và hệ thống điện được bảo hành 8 năm hoặc 150,000 km. Hỗ trợ cứu hộ 24/7 trong thời gian bảo hành. Áp dụng cho xe sử dụng cá nhân.', GETDATE()),
    ('050e8400-e29b-41d4-a716-446655440029', 'BYD_SEAL_2025_WARRANTY', 72, 150000, N'Bảo hành 6 năm hoặc 150,000 km (tùy điều kiện nào đến trước). Pin và hệ thống điện được bảo hành 8 năm hoặc 150,000 km. Hỗ trợ cứu hộ 24/7 trong thời gian bảo hành. Áp dụng cho xe sử dụng cá nhân.', GETDATE()),

    -- BYD Tang
    ('050e8400-e29b-41d4-a716-446655440030', 'BYD_TANG_2024_WARRANTY', 72, 150000, N'Bảo hành 6 năm hoặc 150,000 km (tùy điều kiện nào đến trước). Pin và hệ thống điện được bảo hành 8 năm hoặc 150,000 km. Hỗ trợ cứu hộ 24/7 trong thời gian bảo hành. Áp dụng cho xe sử dụng cá nhân.', GETDATE()),
    ('050e8400-e29b-41d4-a716-446655440031', 'BYD_TANG_2025_WARRANTY', 72, 150000, N'Bảo hành 6 năm hoặc 150,000 km (tùy điều kiện nào đến trước). Pin và hệ thống điện được bảo hành 8 năm hoặc 150,000 km. Hỗ trợ cứu hộ 24/7 trong thời gian bảo hành. Áp dụng cho xe sử dụng cá nhân.', GETDATE()),

    -- BYD Seal U
    ('050e8400-e29b-41d4-a716-446655440032', 'BYD_SEALU_2024_WARRANTY', 72, 150000, N'Bảo hành 6 năm hoặc 150,000 km (tùy điều kiện nào đến trước). Pin và hệ thống điện được bảo hành 8 năm hoặc 150,000 km. Hỗ trợ cứu hộ 24/7 trong thời gian bảo hành. Áp dụng cho xe sử dụng cá nhân.', GETDATE()),
    ('050e8400-e29b-41d4-a716-446655440033', 'BYD_SEALU_2025_WARRANTY', 72, 150000, N'Bảo hành 6 năm hoặc 150,000 km (tùy điều kiện nào đến trước). Pin và hệ thống điện được bảo hành 8 năm hoặc 150,000 km. Hỗ trợ cứu hộ 24/7 trong thời gian bảo hành. Áp dụng cho xe sử dụng cá nhân.', GETDATE()),

    -- BYD M6
    ('050e8400-e29b-41d4-a716-446655440034', 'BYD_M6_2024_WARRANTY', 72, 150000, N'Bảo hành 6 năm hoặc 150,000 km (tùy điều kiện nào đến trước). Pin và hệ thống điện được bảo hành 8 năm hoặc 150,000 km. Hỗ trợ cứu hộ 24/7 trong thời gian bảo hành. Áp dụng cho xe sử dụng cá nhân.', GETDATE()),
    ('050e8400-e29b-41d4-a716-446655440035', 'BYD_M6_2025_WARRANTY', 72, 150000, N'Bảo hành 6 năm hoặc 150,000 km (tùy điều kiện nào đến trước). Pin và hệ thống điện được bảo hành 8 năm hoặc 150,000 km. Hỗ trợ cứu hộ 24/7 trong thời gian bảo hành. Áp dụng cho xe sử dụng cá nhân.', GETDATE()),

    -- Mercedes-Benz EQA
    ('050e8400-e29b-41d4-a716-446655440036', 'MB_EQA_2023_WARRANTY', 48, 100000, N'Bảo hành 4 năm hoặc 100,000 km (tùy điều kiện nào đến trước). Pin được bảo hành 8 năm hoặc 160,000 km. Bảo dưỡng định kỳ tại hệ thống đại lý ủy quyền Mercedes-Benz. Áp dụng cho xe sử dụng cá nhân.', GETDATE()),
    ('050e8400-e29b-41d4-a716-446655440037', 'MB_EQA_2024_WARRANTY', 48, 100000, N'Bảo hành 4 năm hoặc 100,000 km (tùy điều kiện nào đến trước). Pin được bảo hành 8 năm hoặc 160,000 km. Bảo dưỡng định kỳ tại hệ thống đại lý ủy quyền Mercedes-Benz. Áp dụng cho xe sử dụng cá nhân.', GETDATE()),
    ('050e8400-e29b-41d4-a716-446655440038', 'MB_EQA_2025_WARRANTY', 48, 100000, N'Bảo hành 4 năm hoặc 100,000 km (tùy điều kiện nào đến trước). Pin được bảo hành 8 năm hoặc 160,000 km. Bảo dưỡng định kỳ tại hệ thống đại lý ủy quyền Mercedes-Benz. Áp dụng cho xe sử dụng cá nhân.', GETDATE()),

    -- Mercedes-Benz EQB
    ('050e8400-e29b-41d4-a716-446655440039', 'MB_EQB_2023_WARRANTY', 48, 100000, N'Bảo hành 4 năm hoặc 100,000 km (tùy điều kiện nào đến trước). Pin được bảo hành 8 năm hoặc 160,000 km. Bảo dưỡng định kỳ tại hệ thống đại lý ủy quyền Mercedes-Benz. Áp dụng cho xe sử dụng cá nhân.', GETDATE()),
    ('050e8400-e29b-41d4-a716-446655440040', 'MB_EQB_2024_WARRANTY', 48, 100000, N'Bảo hành 4 năm hoặc 100,000 km (tùy điều kiện nào đến trước). Pin được bảo hành 8 năm hoặc 160,000 km. Bảo dưỡng định kỳ tại hệ thống đại lý ủy quyền Mercedes-Benz. Áp dụng cho xe sử dụng cá nhân.', GETDATE()),
    ('050e8400-e29b-41d4-a716-446655440041', 'MB_EQB_2025_WARRANTY', 48, 100000, N'Bảo hành 4 năm hoặc 100,000 km (tùy điều kiện nào đến trước). Pin được bảo hành 8 năm hoặc 160,000 km. Bảo dưỡng định kỳ tại hệ thống đại lý ủy quyền Mercedes-Benz. Áp dụng cho xe sử dụng cá nhân.', GETDATE()),

    -- Mercedes-Benz EQC
    ('050e8400-e29b-41d4-a716-446655440042', 'MB_EQC_2022_WARRANTY', 48, 100000, N'Bảo hành 4 năm hoặc 100,000 km (tùy điều kiện nào đến trước). Pin được bảo hành 8 năm hoặc 160,000 km. Bảo dưỡng định kỳ tại hệ thống đại lý ủy quyền Mercedes-Benz. Áp dụng cho xe sử dụng cá nhân.', GETDATE()),
    ('050e8400-e29b-41d4-a716-446655440043', 'MB_EQC_2023_WARRANTY', 48, 100000, N'Bảo hành 4 năm hoặc 100,000 km (tùy điều kiện nào đến trước). Pin được bảo hành 8 năm hoặc 160,000 km. Bảo dưỡng định kỳ tại hệ thống đại lý ủy quyền Mercedes-Benz. Áp dụng cho xe sử dụng cá nhân.', GETDATE()),
    ('050e8400-e29b-41d4-a716-446655440044', 'MB_EQC_2024_WARRANTY', 48, 100000, N'Bảo hành 4 năm hoặc 100,000 km (tùy điều kiện nào đến trước). Pin được bảo hành 8 năm hoặc 160,000 km. Bảo dưỡng định kỳ tại hệ thống đại lý ủy quyền Mercedes-Benz. Áp dụng cho xe sử dụng cá nhân.', GETDATE()),

    -- Mercedes-Benz EQE
    ('050e8400-e29b-41d4-a716-446655440045', 'MB_EQE_2023_WARRANTY', 48, 100000, N'Bảo hành 4 năm hoặc 100,000 km (tùy điều kiện nào đến trước). Pin được bảo hành 8 năm hoặc 160,000 km. Bảo dưỡng định kỳ tại hệ thống đại lý ủy quyền Mercedes-Benz. Hỗ trợ cứu hộ 24/7. Áp dụng cho xe sử dụng cá nhân.', GETDATE()),
    ('050e8400-e29b-41d4-a716-446655440046', 'MB_EQE_2024_WARRANTY', 48, 100000, N'Bảo hành 4 năm hoặc 100,000 km (tùy điều kiện nào đến trước). Pin được bảo hành 8 năm hoặc 160,000 km. Bảo dưỡng định kỳ tại hệ thống đại lý ủy quyền Mercedes-Benz. Hỗ trợ cứu hộ 24/7. Áp dụng cho xe sử dụng cá nhân.', GETDATE()),
    ('050e8400-e29b-41d4-a716-446655440047', 'MB_EQE_2025_WARRANTY', 48, 100000, N'Bảo hành 4 năm hoặc 100,000 km (tùy điều kiện nào đến trước). Pin được bảo hành 8 năm hoặc 160,000 km. Bảo dưỡng định kỳ tại hệ thống đại lý ủy quyền Mercedes-Benz. Hỗ trợ cứu hộ 24/7. Áp dụng cho xe sử dụng cá nhân.', GETDATE()),

    -- Mercedes-Benz EQE SUV
    ('050e8400-e29b-41d4-a716-446655440048', 'MB_EQESUV_2023_WARRANTY', 48, 100000, N'Bảo hành 4 năm hoặc 100,000 km (tùy điều kiện nào đến trước). Pin được bảo hành 8 năm hoặc 160,000 km. Bảo dưỡng định kỳ tại hệ thống đại lý ủy quyền Mercedes-Benz. Hỗ trợ cứu hộ 24/7. Áp dụng cho xe sử dụng cá nhân.', GETDATE()),
    ('050e8400-e29b-41d4-a716-446655440049', 'MB_EQESUV_2024_WARRANTY', 48, 100000, N'Bảo hành 4 năm hoặc 100,000 km (tùy điều kiện nào đến trước). Pin được bảo hành 8 năm hoặc 160,000 km. Bảo dưỡng định kỳ tại hệ thống đại lý ủy quyền Mercedes-Benz. Hỗ trợ cứu hộ 24/7. Áp dụng cho xe sử dụng cá nhân.', GETDATE()),
    ('050e8400-e29b-41d4-a716-446655440050', 'MB_EQESUV_2025_WARRANTY', 48, 100000, N'Bảo hành 4 năm hoặc 100,000 km (tùy điều kiện nào đến trước). Pin được bảo hành 8 năm hoặc 160,000 km. Bảo dưỡng định kỳ tại hệ thống đại lý ủy quyền Mercedes-Benz. Hỗ trợ cứu hộ 24/7. Áp dụng cho xe sử dụng cá nhân.', GETDATE()),

    ('050e8400-e29b-41d4-a716-446655440051', 'MB_EQS_2023_WARRANTY', 48, 100000, N'Bảo hành 4 năm hoặc 100,000 km (tùy điều kiện nào đến trước). Pin được bảo hành 8 năm hoặc 160,000 km. Bảo dưỡng định kỳ tại hệ thống đại lý ủy quyền Mercedes-Benz. Hỗ trợ cứu hộ 24/7. Áp dụng cho xe sử dụng cá nhân.', GETDATE()),
    ('050e8400-e29b-41d4-a716-446655440052', 'MB_EQS_2024_WARRANTY', 48, 100000, N'Bảo hành 4 năm hoặc 100,000 km (tùy điều kiện nào đến trước). Pin được bảo hành 8 năm hoặc 160,000 km. Bảo dưỡng định kỳ tại hệ thống đại lý ủy quyền Mercedes-Benz. Hỗ trợ cứu hộ 24/7. Áp dụng cho xe sử dụng cá nhân.', GETDATE()),
    ('050e8400-e29b-41d4-a716-446655440053', 'MB_EQS_2025_WARRANTY', 48, 100000, N'Bảo hành 4 năm hoặc 100,000 km (tùy điều kiện nào đến trước). Pin được bảo hành 8 năm hoặc 160,000 km. Bảo dưỡng định kỳ tại hệ thống đại lý ủy quyền Mercedes-Benz. Hỗ trợ cứu hộ 24/7. Áp dụng cho xe sử dụng cá nhân.', GETDATE()),

    ('050e8400-e29b-41d4-a716-446655440054', 'MB_EQSSUV_2023_WARRANTY', 48, 100000, N'Bảo hành 4 năm hoặc 100,000 km (tùy điều kiện nào đến trước). Pin được bảo hành 8 năm hoặc 160,000 km. Bảo dưỡng định kỳ tại hệ thống đại lý ủy quyền Mercedes-Benz. Hỗ trợ cứu hộ 24/7. Áp dụng cho xe sử dụng cá nhân.', GETDATE()),
    ('050e8400-e29b-41d4-a716-446655440055', 'MB_EQSSUV_2024_WARRANTY', 48, 100000, N'Bảo hành 4 năm hoặc 100,000 km (tùy điều kiện nào đến trước). Pin được bảo hành 8 năm hoặc 160,000 km. Bảo dưỡng định kỳ tại hệ thống đại lý ủy quyền Mercedes-Benz. Hỗ trợ cứu hộ 24/7. Áp dụng cho xe sử dụng cá nhân.', GETDATE()),
    ('050e8400-e29b-41d4-a716-446655440056', 'MB_EQSSUV_2025_WARRANTY', 48, 100000, N'Bảo hành 4 năm hoặc 100,000 km (tùy điều kiện nào đến trước). Pin được bảo hành 8 năm hoặc 160,000 km. Bảo dưỡng định kỳ tại hệ thống đại lý ủy quyền Mercedes-Benz. Hỗ trợ cứu hộ 24/7. Áp dụng cho xe sử dụng cá nhân.', GETDATE());


-- Insert Vehicle Models (EV models from VinFast, BYD, and Mercedes-Benz available in Vietnam)

-- =========================================================
-- * DECLARE GUIDs for Vehicle Models
-- =========================================================
DECLARE
    -- VinFast 4401xx
    @vf8_2023 UNIQUEIDENTIFIER = '650e8400-e29b-41d4-a716-446655440101',
    @vf8_2024 UNIQUEIDENTIFIER = '650e8400-e29b-41d4-a716-446655440102',
    @vf8_2025 UNIQUEIDENTIFIER = '650e8400-e29b-41d4-a716-446655440103',

    @vf9_2023 UNIQUEIDENTIFIER = '650e8400-e29b-41d4-a716-446655440104',
    @vf9_2024 UNIQUEIDENTIFIER = '650e8400-e29b-41d4-a716-446655440105',
    @vf9_2025 UNIQUEIDENTIFIER = '650e8400-e29b-41d4-a716-446655440106',

    @vfe34_2022 UNIQUEIDENTIFIER = '650e8400-e29b-41d4-a716-446655440107',
    @vfe34_2023 UNIQUEIDENTIFIER = '650e8400-e29b-41d4-a716-446655440108',
    @vfe34_2024 UNIQUEIDENTIFIER = '650e8400-e29b-41d4-a716-446655440109',
    @vfe34_2025 UNIQUEIDENTIFIER = '650e8400-e29b-41d4-a716-446655440110',

    @vf5_2023 UNIQUEIDENTIFIER = '650e8400-e29b-41d4-a716-446655440111',
    @vf5_2024 UNIQUEIDENTIFIER = '650e8400-e29b-41d4-a716-446655440112',
    @vf5_2025 UNIQUEIDENTIFIER = '650e8400-e29b-41d4-a716-446655440113',

    @vf6_2024 UNIQUEIDENTIFIER = '650e8400-e29b-41d4-a716-446655440114',
    @vf6_2025 UNIQUEIDENTIFIER = '650e8400-e29b-41d4-a716-446655440115',

    @vf7_2024 UNIQUEIDENTIFIER = '650e8400-e29b-41d4-a716-446655440116',
    @vf7_2025 UNIQUEIDENTIFIER = '650e8400-e29b-41d4-a716-446655440117',

    @vf3_2024 UNIQUEIDENTIFIER = '650e8400-e29b-41d4-a716-446655440118',
    @vf3_2025 UNIQUEIDENTIFIER = '650e8400-e29b-41d4-a716-446655440119',

    -- BYD 4402xx
    @byd_atto3_2022 UNIQUEIDENTIFIER = '650e8400-e29b-41d4-a716-446655440201',
    @byd_atto3_2023 UNIQUEIDENTIFIER = '650e8400-e29b-41d4-a716-446655440202',
    @byd_atto3_2024 UNIQUEIDENTIFIER = '650e8400-e29b-41d4-a716-446655440203',
    @byd_atto3_2025 UNIQUEIDENTIFIER = '650e8400-e29b-41d4-a716-446655440204',

    @byd_dolphin_2023 UNIQUEIDENTIFIER = '650e8400-e29b-41d4-a716-446655440205',
    @byd_dolphin_2024 UNIQUEIDENTIFIER = '650e8400-e29b-41d4-a716-446655440206',
    @byd_dolphin_2025 UNIQUEIDENTIFIER = '650e8400-e29b-41d4-a716-446655440207',

    @byd_seal_2023 UNIQUEIDENTIFIER = '650e8400-e29b-41d4-a716-446655440208',
    @byd_seal_2024 UNIQUEIDENTIFIER = '650e8400-e29b-41d4-a716-446655440209',
    @byd_seal_2025 UNIQUEIDENTIFIER = '650e8400-e29b-41d4-a716-446655440210',

    @byd_tang_2024 UNIQUEIDENTIFIER = '650e8400-e29b-41d4-a716-446655440211',
    @byd_tang_2025 UNIQUEIDENTIFIER = '650e8400-e29b-41d4-a716-446655440212',

    @byd_sealu_2024 UNIQUEIDENTIFIER = '650e8400-e29b-41d4-a716-446655440213',
    @byd_sealu_2025 UNIQUEIDENTIFIER = '650e8400-e29b-41d4-a716-446655440214',

    @byd_m6_2024 UNIQUEIDENTIFIER = '650e8400-e29b-41d4-a716-446655440215',
    @byd_m6_2025 UNIQUEIDENTIFIER = '650e8400-e29b-41d4-a716-446655440216',

    -- Mercedes-Benz 4403xx
    @mb_eqa_2023 UNIQUEIDENTIFIER = '650e8400-e29b-41d4-a716-446655440301',
    @mb_eqa_2024 UNIQUEIDENTIFIER = '650e8400-e29b-41d4-a716-446655440302',
    @mb_eqa_2025 UNIQUEIDENTIFIER = '650e8400-e29b-41d4-a716-446655440303',

    @mb_eqb_2023 UNIQUEIDENTIFIER = '650e8400-e29b-41d4-a716-446655440304',
    @mb_eqb_2024 UNIQUEIDENTIFIER = '650e8400-e29b-41d4-a716-446655440305',
    @mb_eqb_2025 UNIQUEIDENTIFIER = '650e8400-e29b-41d4-a716-446655440306',

    @mb_eqc_2022 UNIQUEIDENTIFIER = '650e8400-e29b-41d4-a716-446655440307',
    @mb_eqc_2023 UNIQUEIDENTIFIER = '650e8400-e29b-41d4-a716-446655440308',
    @mb_eqc_2024 UNIQUEIDENTIFIER = '650e8400-e29b-41d4-a716-446655440309',

    @mb_eqe_2023 UNIQUEIDENTIFIER = '650e8400-e29b-41d4-a716-446655440310',
    @mb_eqe_2024 UNIQUEIDENTIFIER = '650e8400-e29b-41d4-a716-446655440311',
    @mb_eqe_2025 UNIQUEIDENTIFIER = '650e8400-e29b-41d4-a716-446655440312',

    @mb_eqesuv_2023 UNIQUEIDENTIFIER = '650e8400-e29b-41d4-a716-446655440313',
    @mb_eqesuv_2024 UNIQUEIDENTIFIER = '650e8400-e29b-41d4-a716-446655440314',
    @mb_eqesuv_2025 UNIQUEIDENTIFIER = '650e8400-e29b-41d4-a716-446655440315',

    @mb_eqs_2023 UNIQUEIDENTIFIER = '650e8400-e29b-41d4-a716-446655440316',
    @mb_eqs_2024 UNIQUEIDENTIFIER = '650e8400-e29b-41d4-a716-446655440317',
    @mb_eqs_2025 UNIQUEIDENTIFIER = '650e8400-e29b-41d4-a716-446655440318',

    @mb_eqssuv_2023 UNIQUEIDENTIFIER = '650e8400-e29b-41d4-a716-446655440319',
    @mb_eqssuv_2024 UNIQUEIDENTIFIER = '650e8400-e29b-41d4-a716-446655440320',
    @mb_eqssuv_2025 UNIQUEIDENTIFIER = '650e8400-e29b-41d4-a716-446655440321';

-- =========================================================
-- * INSERT DATA (VinFast / BYD / Mercedes-Benz)
-- =========================================================

INSERT INTO vehicle_models (id, brand, model_name, year, created_at, policy_id)
VALUES
    -- VinFast
    (@vf8_2023, 'VinFast', 'VF8', 2023, GETDATE(), '050e8400-e29b-41d4-a716-446655440001'),
    (@vf8_2024, 'VinFast', 'VF8', 2024, GETDATE(), '050e8400-e29b-41d4-a716-446655440002'),
    (@vf8_2025, 'VinFast', 'VF8', 2025, GETDATE(), '050e8400-e29b-41d4-a716-446655440003'),
    (@vf9_2023, 'VinFast', 'VF9', 2023, GETDATE(), '050e8400-e29b-41d4-a716-446655440004'),
    (@vf9_2024, 'VinFast', 'VF9', 2024, GETDATE(), '050e8400-e29b-41d4-a716-446655440005'),
    (@vf9_2025, 'VinFast', 'VF9', 2025, GETDATE(), '050e8400-e29b-41d4-a716-446655440006'),
    (@vfe34_2022, 'VinFast', 'VF e34', 2022, GETDATE(), '050e8400-e29b-41d4-a716-446655440007'),
    (@vfe34_2023, 'VinFast', 'VF e34', 2023, GETDATE(), '050e8400-e29b-41d4-a716-446655440008'),
    (@vfe34_2024, 'VinFast', 'VF e34', 2024, GETDATE(), '050e8400-e29b-41d4-a716-446655440009'),
    (@vfe34_2025, 'VinFast', 'VF e34', 2025, GETDATE(), '050e8400-e29b-41d4-a716-446655440010'),
    (@vf5_2023, 'VinFast', 'VF5 Plus', 2023, GETDATE(), '050e8400-e29b-41d4-a716-446655440011'),
    (@vf5_2024, 'VinFast', 'VF5 Plus', 2024, GETDATE(), '050e8400-e29b-41d4-a716-446655440012'),
    (@vf5_2025, 'VinFast', 'VF5 Plus', 2025, GETDATE(), '050e8400-e29b-41d4-a716-446655440013'),
    (@vf6_2024, 'VinFast', 'VF6', 2024, GETDATE(), '050e8400-e29b-41d4-a716-446655440014'),
    (@vf6_2025, 'VinFast', 'VF6', 2025, GETDATE(), '050e8400-e29b-41d4-a716-446655440015'),
    (@vf7_2024, 'VinFast', 'VF7', 2024, GETDATE(), '050e8400-e29b-41d4-a716-446655440016'),
    (@vf7_2025, 'VinFast', 'VF7', 2025, GETDATE(), '050e8400-e29b-41d4-a716-446655440017'),
    (@vf3_2024, 'VinFast', 'VF3', 2024, GETDATE(), '050e8400-e29b-41d4-a716-446655440018'),
    (@vf3_2025, 'VinFast', 'VF3', 2025, GETDATE(), '050e8400-e29b-41d4-a716-446655440019'),

    -- BYD
    (@byd_atto3_2022, 'BYD', 'Atto 3', 2022, GETDATE(), '050e8400-e29b-41d4-a716-446655440020'),
    (@byd_atto3_2023, 'BYD', 'Atto 3', 2023, GETDATE(), '050e8400-e29b-41d4-a716-446655440021'),
    (@byd_atto3_2024, 'BYD', 'Atto 3', 2024, GETDATE(), '050e8400-e29b-41d4-a716-446655440022'),
    (@byd_atto3_2025, 'BYD', 'Atto 3', 2025, GETDATE(), '050e8400-e29b-41d4-a716-446655440023'),
    (@byd_dolphin_2023, 'BYD', 'Dolphin', 2023, GETDATE(), '050e8400-e29b-41d4-a716-446655440024'),
    (@byd_dolphin_2024, 'BYD', 'Dolphin', 2024, GETDATE(), '050e8400-e29b-41d4-a716-446655440025'),
    (@byd_dolphin_2025, 'BYD', 'Dolphin', 2025, GETDATE(), '050e8400-e29b-41d4-a716-446655440026'),
    (@byd_seal_2023, 'BYD', 'Seal', 2023, GETDATE(), '050e8400-e29b-41d4-a716-446655440027'),
    (@byd_seal_2024, 'BYD', 'Seal', 2024, GETDATE(), '050e8400-e29b-41d4-a716-446655440028'),
    (@byd_seal_2025, 'BYD', 'Seal', 2025, GETDATE(), '050e8400-e29b-41d4-a716-446655440029'),
    (@byd_tang_2024, 'BYD', 'Tang EV', 2024, GETDATE(), '050e8400-e29b-41d4-a716-446655440030'),
    (@byd_tang_2025, 'BYD', 'Tang EV', 2025, GETDATE(), '050e8400-e29b-41d4-a716-446655440031'),
    (@byd_sealu_2024, 'BYD', 'Seal U', 2024, GETDATE(), '050e8400-e29b-41d4-a716-446655440032'),
    (@byd_sealu_2025, 'BYD', 'Seal U', 2025, GETDATE(), '050e8400-e29b-41d4-a716-446655440033'),
    (@byd_m6_2024, 'BYD', 'M6', 2024, GETDATE(), '050e8400-e29b-41d4-a716-446655440034'),
    (@byd_m6_2025, 'BYD', 'M6', 2025, GETDATE(), '050e8400-e29b-41d4-a716-446655440035'),

    -- Mercedes-Benz
    (@mb_eqa_2023, 'Mercedes-Benz', 'EQA', 2023, GETDATE(), '050e8400-e29b-41d4-a716-446655440036'),
    (@mb_eqa_2024, 'Mercedes-Benz', 'EQA', 2024, GETDATE(), '050e8400-e29b-41d4-a716-446655440037'),
    (@mb_eqa_2025, 'Mercedes-Benz', 'EQA', 2025, GETDATE(), '050e8400-e29b-41d4-a716-446655440038'),
    (@mb_eqb_2023, 'Mercedes-Benz', 'EQB', 2023, GETDATE(), '050e8400-e29b-41d4-a716-446655440039'),
    (@mb_eqb_2024, 'Mercedes-Benz', 'EQB', 2024, GETDATE(), '050e8400-e29b-41d4-a716-446655440040'),
    (@mb_eqb_2025, 'Mercedes-Benz', 'EQB', 2025, GETDATE(), '050e8400-e29b-41d4-a716-446655440041'),
    (@mb_eqc_2022, 'Mercedes-Benz', 'EQC', 2022, GETDATE(), '050e8400-e29b-41d4-a716-446655440042'),
    (@mb_eqc_2023, 'Mercedes-Benz', 'EQC', 2023, GETDATE(), '050e8400-e29b-41d4-a716-446655440043'),
    (@mb_eqc_2024, 'Mercedes-Benz', 'EQC', 2024, GETDATE(), '050e8400-e29b-41d4-a716-446655440044'),
    (@mb_eqe_2023, 'Mercedes-Benz', 'EQE', 2023, GETDATE(), '050e8400-e29b-41d4-a716-446655440045'),
    (@mb_eqe_2024, 'Mercedes-Benz', 'EQE', 2024, GETDATE(), '050e8400-e29b-41d4-a716-446655440046'),
    (@mb_eqe_2025, 'Mercedes-Benz', 'EQE', 2025, GETDATE(), '050e8400-e29b-41d4-a716-446655440047'),
    (@mb_eqesuv_2023, 'Mercedes-Benz', 'EQE SUV', 2023, GETDATE(), '050e8400-e29b-41d4-a716-446655440048'),
    (@mb_eqesuv_2024, 'Mercedes-Benz', 'EQE SUV', 2024, GETDATE(), '050e8400-e29b-41d4-a716-446655440049'),
    (@mb_eqesuv_2025, 'Mercedes-Benz', 'EQE SUV', 2025, GETDATE(), '050e8400-e29b-41d4-a716-446655440050'),
    (@mb_eqs_2023, 'Mercedes-Benz', 'EQS', 2023, GETDATE(), '050e8400-e29b-41d4-a716-446655440051'),
    (@mb_eqs_2024, 'Mercedes-Benz', 'EQS', 2024, GETDATE(), '050e8400-e29b-41d4-a716-446655440052'),
    (@mb_eqs_2025, 'Mercedes-Benz', 'EQS', 2025, GETDATE(), '050e8400-e29b-41d4-a716-446655440053'),
    (@mb_eqssuv_2023, 'Mercedes-Benz', 'EQS SUV', 2023, GETDATE(), '050e8400-e29b-41d4-a716-446655440054'),
    (@mb_eqssuv_2024, 'Mercedes-Benz', 'EQS SUV', 2024, GETDATE(), '050e8400-e29b-41d4-a716-446655440055'),
    (@mb_eqssuv_2025, 'Mercedes-Benz', 'EQS SUV', 2025, GETDATE(), '050e8400-e29b-41d4-a716-446655440056');

INSERT INTO customers (id, first_name, last_name, phone_number, email, address, created_at) VALUES
    ('850e8400-e29b-41d4-a716-446655440001', N'Minh', N'Nguyễn Văn', '0901234567', 'minh.nv@gmail.com', N'123 Nguyễn Huệ, Quận 1, TP.HCM', GETDATE()),
    ('850e8400-e29b-41d4-a716-446655440002', N'Lan', N'Trần Thị', '0902345678', 'lan.tt@gmail.com', N'456 Lê Lợi, Quận 1, TP.HCM', GETDATE()),
    ('850e8400-e29b-41d4-a716-446655440003', N'Hùng', N'Lê Văn', '0903456789', 'hung.lv@gmail.com', N'789 Hai Bà Trưng, Quận 3, TP.HCM', GETDATE()),
    ('850e8400-e29b-41d4-a716-446655440004', N'Hoa', N'Phạm Thị', '0904567890', 'hoa.pt@gmail.com', N'321 Võ Văn Tần, Quận 3, TP.HCM', GETDATE()),
    ('850e8400-e29b-41d4-a716-446655440005', N'Tuấn', N'Hoàng Văn', '0905678901', 'tuan.hv@gmail.com', N'654 Nguyễn Thị Minh Khai, Quận 1, TP.HCM', GETDATE()),
    ('850e8400-e29b-41d4-a716-446655440006', N'Linh', N'Vũ Thị', '0906789012', 'linh.vt@gmail.com', N'987 Trần Hưng Đạo, Quận 5, TP.HCM', GETDATE()),
    ('850e8400-e29b-41d4-a716-446655440007', N'Đức', N'Đặng Văn', '0907890123', 'duc.dv@gmail.com', N'159 Lý Thường Kiệt, Quận 10, TP.HCM', GETDATE()),
    ('850e8400-e29b-41d4-a716-446655440008', N'Mai', N'Bùi Thị', '0908901234', 'mai.bt@gmail.com', N'753 Cách Mạng Tháng 8, Quận Tân Bình, TP.HCM', GETDATE()),
    ('850e8400-e29b-41d4-a716-446655440009', N'Nam', N'Võ Văn', '0909012345', 'nam.vv@gmail.com', N'852 Lạc Long Quân, Quận 11, TP.HCM', GETDATE()),
    ('850e8400-e29b-41d4-a716-44665544000a', N'Thảo', N'Đinh Thị', '0900123456', 'thao.dt@gmail.com', N'246 Âu Cơ, Quận Tân Phú, TP.HCM', GETDATE()),

    ('850e8400-e29b-41d4-a716-44665544000b', N'Phong', N'Ngô Văn', '0911234567', 'phong.nv@gmail.com', N'135 Nguyễn Văn Cừ, Quận 5, TP.HCM', GETDATE()),
    ('850e8400-e29b-41d4-a716-44665544000c', N'Hương', N'Lý Thị', '0912345678', 'huong.lt@gmail.com', N'468 Hoàng Văn Thụ, Quận Phú Nhuận, TP.HCM', GETDATE()),
    ('850e8400-e29b-41d4-a716-44665544000d', N'Bình', N'Trương Văn', '0913456789', 'binh.tv@gmail.com', N'791 Điện Biên Phủ, Quận 3, TP.HCM', GETDATE()),
    ('850e8400-e29b-41d4-a716-44665544000e', N'Nga', N'Dương Thị', '0914567890', 'nga.dt@gmail.com', N'357 Phan Đăng Lưu, Quận Phú Nhuận, TP.HCM', GETDATE()),
    ('850e8400-e29b-41d4-a716-44665544000f', N'Kiên', N'Phan Văn', '0915678901', 'kien.pv@gmail.com', N'680 Xô Viết Nghệ Tĩnh, Quận Bình Thạnh, TP.HCM', GETDATE()),
    ('850e8400-e29b-41d4-a716-446655440010', N'Yến', N'Tô Thị', '0916789012', 'yen.tt@gmail.com', N'913 Hồng Bàng, Quận 6, TP.HCM', GETDATE()),
    ('850e8400-e29b-41d4-a716-446655440011', N'Long', N'Cao Văn', '0917890123', 'long.cv@gmail.com', N'147 Phạm Văn Đồng, Quận Bình Thạnh, TP.HCM', GETDATE()),
    ('850e8400-e29b-41d4-a716-446655440012', N'Trang', N'Lương Thị', '0918901234', 'trang.lt@gmail.com', N'582 Lê Văn Sỹ, Quận 3, TP.HCM', GETDATE()),
    ('850e8400-e29b-41d4-a716-446655440013', N'Khánh', N'Đỗ Văn', '0919012345', 'khanh.dv@gmail.com', N'926 Ba Tháng Hai, Quận 10, TP.HCM', GETDATE()),
    ('850e8400-e29b-41d4-a716-446655440014', N'Ánh', N'Hà Thị', '0910123456', 'anh.ht@gmail.com', N'369 Tân Sơn Nhì, Quận Tân Phú, TP.HCM', GETDATE()),

    ('850e8400-e29b-41d4-a716-446655440015', N'Quân', N'Trịnh Văn', '0921234567', 'quan.tv@gmail.com', N'78 Trần Quang Khải, Quận 1, TP.HCM', GETDATE()),
    ('850e8400-e29b-41d4-a716-446655440016', N'Chi', N'Võ Thị', '0922345678', 'chi.vt@gmail.com', N'142 Yersin, Quận 1, TP.HCM', GETDATE()),
    ('850e8400-e29b-41d4-a716-446655440017', N'Thắng', N'Lâm Văn', '0923456789', 'thang.lv@gmail.com', N'286 Nguyễn Trãi, Quận 1, TP.HCM', GETDATE()),
    ('850e8400-e29b-41d4-a716-446655440018', N'Nhung', N'Đào Thị', '0924567890', 'nhung.dt@gmail.com', N'514 Võ Thị Sáu, Quận 3, TP.HCM', GETDATE()),
    ('850e8400-e29b-41d4-a716-446655440019', N'Sơn', N'Mai Văn', '0925678901', 'son.mv@gmail.com', N'628 Nguyễn Đình Chiểu, Quận 3, TP.HCM', GETDATE()),
    ('850e8400-e29b-41d4-a716-44665544001a', N'Thúy', N'Trần Thị Thanh', '0926789012', 'thuy.ttt@gmail.com', N'771 Pasteur, Quận 3, TP.HCM', GETDATE()),
    ('850e8400-e29b-41d4-a716-44665544001b', N'Dũng', N'Nguyễn Văn', '0927890123', 'dung.nv@gmail.com', N'853 Nam Kỳ Khởi Nghĩa, Quận 3, TP.HCM', GETDATE()),
    ('850e8400-e29b-41d4-a716-44665544001c', N'Hằng', N'Phạm Thị', '0928901234', 'hang.pt@gmail.com', N'965 Cộng Hòa, Quận Tân Bình, TP.HCM', GETDATE()),
    ('850e8400-e29b-41d4-a716-44665544001d', N'Tiến', N'Lê Văn', '0929012345', 'tien.lv@gmail.com', N'177 Hoàng Hoa Thám, Quận Tân Bình, TP.HCM', GETDATE()),
    ('850e8400-e29b-41d4-a716-44665544001e', N'Vy', N'Hoàng Thị', '0920123456', 'vy.ht@gmail.com', N'389 Trường Chinh, Quận Tân Bình, TP.HCM', GETDATE()),

    ('850e8400-e29b-41d4-a716-44665544001f', N'Đạt', N'Vũ Văn', '0931234567', 'dat.vv@gmail.com', N'512 Nguyễn Kiệm, Quận Phú Nhuận, TP.HCM', GETDATE()),
    ('850e8400-e29b-41d4-a716-446655440020', N'Giang', N'Đặng Thị', '0932345678', 'giang.dt@gmail.com', N'643 Nguyễn Văn Trỗi, Quận Phú Nhuận, TP.HCM', GETDATE()),
    ('850e8400-e29b-41d4-a716-446655440021', N'Hải', N'Bùi Văn', '0933456789', 'hai.bv@gmail.com', N'775 Hoa Lan, Quận Phú Nhuận, TP.HCM', GETDATE()),
    ('850e8400-e29b-41d4-a716-446655440022', N'Ngọc', N'Võ Thị', '0934567890', 'ngoc.vt@gmail.com', N'886 Ung Văn Khiêm, Quận Bình Thạnh, TP.HCM', GETDATE()),
    ('850e8400-e29b-41d4-a716-446655440023', N'Trung', N'Ngô Văn', '0935678901', 'trung.nv@gmail.com', N'918 Đinh Tiên Hoàng, Quận Bình Thạnh, TP.HCM', GETDATE()),
    ('850e8400-e29b-41d4-a716-446655440024', N'Thanh', N'Lý Thị', '0936789012', 'thanh.lt@gmail.com', N'241 Nơ Trang Long, Quận Bình Thạnh, TP.HCM', GETDATE()),
    ('850e8400-e29b-41d4-a716-446655440025', N'Cường', N'Trương Văn', '0937890123', 'cuong.tv@gmail.com', N'367 Bạch Đằng, Quận Bình Thạnh, TP.HCM', GETDATE()),
    ('850e8400-e29b-41d4-a716-446655440026', N'Phương', N'Dương Thị', '0938901234', 'phuong.dt@gmail.com', N'494 Lê Quang Định, Quận Bình Thạnh, TP.HCM', GETDATE()),
    ('850e8400-e29b-41d4-a716-446655440027', N'Huy', N'Phan Văn', '0939012345', 'huy.pv@gmail.com', N'628 D2, Quận Bình Thạnh, TP.HCM', GETDATE()),
    ('850e8400-e29b-41d4-a716-446655440028', N'Thu', N'Tô Thị', '0930123456', 'thu.tt@gmail.com', N'751 D1, Quận Bình Thạnh, TP.HCM', GETDATE()),

    ('850e8400-e29b-41d4-a716-446655440029', N'Khoa', N'Cao Văn', '0941234567', 'khoa.cv@gmail.com', N'892 Quang Trung, Quận Gò Vấp, TP.HCM', GETDATE()),
    ('850e8400-e29b-41d4-a716-44665544002a', N'Duyên', N'Lương Thị', '0942345678', 'duyen.lt@gmail.com', N'134 Nguyễn Oanh, Quận Gò Vấp, TP.HCM', GETDATE()),
    ('850e8400-e29b-41d4-a716-44665544002b', N'Tài', N'Đỗ Văn', '0943456789', 'tai.dv@gmail.com', N'267 Phan Văn Trị, Quận Gò Vấp, TP.HCM', GETDATE()),
    ('850e8400-e29b-41d4-a716-44665544002c', N'Xuân', N'Hà Thị', '0944567890', 'xuan.ht@gmail.com', N'398 Lê Đức Thọ, Quận Gò Vấp, TP.HCM', GETDATE()),
    ('850e8400-e29b-41d4-a716-44665544002d', N'Vinh', N'Trịnh Văn', '0945678901', 'vinh.tv@gmail.com', N'521 Phan Huy Ích, Quận Gò Vấp, TP.HCM', GETDATE()),
    ('850e8400-e29b-41d4-a716-44665544002e', N'Diệu', N'Võ Thị', '0946789012', 'dieu.vt@gmail.com', N'654 Nguyễn Văn Lượng, Quận Gò Vấp, TP.HCM', GETDATE()),
    ('850e8400-e29b-41d4-a716-44665544002f', N'Toàn', N'Lâm Văn', '0947890123', 'toan.lv@gmail.com', N'787 Thống Nhất, Quận Gò Vấp, TP.HCM', GETDATE()),
    ('850e8400-e29b-41d4-a716-446655440030', N'My', N'Đào Thị', '0948901234', 'my.dt@gmail.com', N'819 Võ Văn Ngân, Quận Thủ Đức, TP.HCM', GETDATE()),
    ('850e8400-e29b-41d4-a716-446655440031', N'Hiếu', N'Mai Văn', '0949012345', 'hieu.mv@gmail.com', N'952 Tô Ngọc Vân, Quận Thủ Đức, TP.HCM', GETDATE()),
    ('850e8400-e29b-41d4-a716-446655440032', N'Loan', N'Trần Thị', '0940123456', 'loan.tt@gmail.com', N'186 Kha Vạn Cân, Quận Thủ Đức, TP.HCM', GETDATE()),

    ('850e8400-e29b-41d4-a716-446655440033', N'An', N'Nguyễn Văn', '0951234567', 'an.nv@gmail.com', N'314 Hoàng Diệu 2, Quận Thủ Đức, TP.HCM', GETDATE()),
    ('850e8400-e29b-41d4-a716-446655440034', N'Nhi', N'Phạm Thị', '0952345678', 'nhi.pt@gmail.com', N'447 Man Thiện, Quận 9, TP.HCM', GETDATE()),
    ('850e8400-e29b-41d4-a716-446655440035', N'Bảo', N'Lê Văn', '0953456789', 'bao.lv@gmail.com', N'579 Đỗ Xuân Hợp, Quận 9, TP.HCM', GETDATE()),
    ('850e8400-e29b-41d4-a716-446655440036', N'Tâm', N'Hoàng Thị', '0954567890', 'tam.ht@gmail.com', N'612 Lê Văn Việt, Quận 9, TP.HCM', GETDATE()),
    ('850e8400-e29b-41d4-a716-446655440037', N'Hưng', N'Vũ Văn', '0955678901', 'hungvu.vv@gmail.com', N'745 Xa Lộ Hà Nội, Quận 9, TP.HCM', GETDATE()),
    ('850e8400-e29b-41d4-a716-446655440038', N'Thư', N'Đặng Thị', '0956789012', 'thu.dt@gmail.com', N'878 Đỗ Bí, Quận 12, TP.HCM', GETDATE()),
    ('850e8400-e29b-41d4-a716-446655440039', N'Quang', N'Bùi Văn', '0957890123', 'quang.bv@gmail.com', N'921 Tô Ký, Quận 12, TP.HCM', GETDATE()),
    ('850e8400-e29b-41d4-a716-44665544003a', N'Vân', N'Võ Thị', '0958901234', 'van.vt@gmail.com', N'154 Lê Thị Riêng, Quận 12, TP.HCM', GETDATE()),
    ('850e8400-e29b-41d4-a716-44665544003b', N'Dương', N'Ngô Văn', '0959012345', 'duong.nv@gmail.com', N'287 Hà Huy Giáp, Quận 12, TP.HCM', GETDATE()),
    ('850e8400-e29b-41d4-a716-44665544003c', N'Trinh', N'Lý Thị', '0950123456', 'trinh.lt@gmail.com', N'419 Trường Chinh, Quận 12, TP.HCM', GETDATE()),

    ('850e8400-e29b-41d4-a716-44665544003d', N'Hoàng', N'Trương Văn', '0961234567', 'hoang.tv@gmail.com', N'542 Hậu Giang, Quận 6, TP.HCM', GETDATE()),
    ('850e8400-e29b-41d4-a716-44665544003e', N'Nhã', N'Dương Thị', '0962345678', 'nha.dt@gmail.com', N'675 Phú Thọ Hòa, Quận Tân Phú, TP.HCM', GETDATE()),
    ('850e8400-e29b-41d4-a716-44665544003f', N'Tâm', N'Phan Văn', '0963456789', 'tamphan.pv@gmail.com', N'708 Tân Kỳ Tân Quý, Quận Tân Phú, TP.HCM', GETDATE()),
    ('850e8400-e29b-41d4-a716-446655440040', N'Hạnh', N'Tô Thị', '0964567890', 'hanh.tt@gmail.com', N'831 Trường Sa, Quận Phú Nhuận, TP.HCM', GETDATE()),
    ('850e8400-e29b-41d4-a716-446655440041', N'Khang', N'Cao Văn', '0965678901', 'khang.cv@gmail.com', N'964 Nguyễn Văn Đậu, Quận Bình Thạnh, TP.HCM', GETDATE()),
    ('850e8400-e29b-41d4-a716-446655440042', N'Như', N'Lương Thị', '0966789012', 'nhu.lt@gmail.com', N'197 Ngô Tất Tố, Quận Bình Thạnh, TP.HCM', GETDATE()),
    ('850e8400-e29b-41d4-a716-446655440043', N'Đăng', N'Đỗ Văn', '0967890123', 'dang.dv@gmail.com', N'329 Bùi Đình Túy, Quận Bình Thạnh, TP.HCM', GETDATE()),
    ('850e8400-e29b-41d4-a716-446655440044', N'Oanh', N'Hà Thị', '0968901234', 'oanh.ht@gmail.com', N'452 Nguyễn Xí, Quận Bình Thạnh, TP.HCM', GETDATE()),
    ('850e8400-e29b-41d4-a716-446655440045', N'Phúc', N'Trịnh Văn', '0969012345', 'phuc.tv@gmail.com', N'585 Đinh Bộ Lĩnh, Quận Bình Thạnh, TP.HCM', GETDATE()),
    ('850e8400-e29b-41d4-a716-446655440046', N'Liên', N'Võ Thị', '0960123456', 'lien.vt@gmail.com', N'718 Nguyễn Hữu Cảnh, Quận Bình Thạnh, TP.HCM', GETDATE()),

    ('850e8400-e29b-41d4-a716-446655440047', N'Tuấn', N'Lâm Văn', '0971234567', 'tuanlam.lv@gmail.com', N'841 Mai Chí Thọ, Quận 2, TP.HCM', GETDATE()),
    ('850e8400-e29b-41d4-a716-446655440048', N'Quyên', N'Đào Thị', '0972345678', 'quyen.dt@gmail.com', N'174 Lương Định Của, Quận 2, TP.HCM', GETDATE()),
    ('850e8400-e29b-41d4-a716-446655440049', N'Thịnh', N'Mai Văn', '0973456789', 'thinh.mv@gmail.com', N'306 Nguyễn Duy Trinh, Quận 2, TP.HCM', GETDATE()),
    ('850e8400-e29b-41d4-a716-44665544004a', N'Kim', N'Trần Thị', '0974567890', 'kim.tt@gmail.com', N'439 Trần Não, Quận 2, TP.HCM', GETDATE()),
    ('850e8400-e29b-41d4-a716-44665544004b', N'Thiện', N'Nguyễn Văn', '0975678901', 'thien.nv@gmail.com', N'572 Xa Lộ Hà Nội, Quận 2, TP.HCM', GETDATE()),
    ('850e8400-e29b-41d4-a716-44665544004c', N'Hiền', N'Phạm Thị', '0976789012', 'hien.pt@gmail.com', N'605 Quốc Hương, Quận 2, TP.HCM', GETDATE()),
    ('850e8400-e29b-41d4-a716-44665544004d', N'Tùng', N'Lê Văn', '0977890123', 'tung.lv@gmail.com', N'738 Nguyễn Ư Dĩ, Quận 2, TP.HCM', GETDATE()),
    ('850e8400-e29b-41d4-a716-44665544004e', N'Hồng', N'Hoàng Thị', '0978901234', 'hong.ht@gmail.com', N'861 Nguyễn Thị Định, Quận 2, TP.HCM', GETDATE()),
    ('850e8400-e29b-41d4-a716-44665544004f', N'Đức', N'Vũ Văn', '0979012345', 'ducvu.vv@gmail.com', N'194 Song Hành, Quận 2, TP.HCM', GETDATE()),
    ('850e8400-e29b-41d4-a716-446655440050', N'Tuyết', N'Đặng Thị', '0970123456', 'tuyet.dt@gmail.com', N'327 Đồng Văn Cống, Quận 2, TP.HCM', GETDATE()),

    ('850e8400-e29b-41d4-a716-446655440051', N'Hải', N'Bùi Văn', '0981234567', 'haibui.bv@gmail.com', N'459 Phạm Văn Nghị, Quận Tân Phú, TP.HCM', GETDATE()),
    ('850e8400-e29b-41d4-a716-446655440052', N'Ngân', N'Võ Thị', '0982345678', 'ngan.vt@gmail.com', N'592 Tân Hương, Quận Tân Phú, TP.HCM', GETDATE()),
    ('850e8400-e29b-41d4-a716-446655440053', N'Lâm', N'Ngô Văn', '0983456789', 'lam.nv@gmail.com', N'625 Tây Thạnh, Quận Tân Phú, TP.HCM', GETDATE()),
    ('850e8400-e29b-41d4-a716-446655440054', N'Dung', N'Lý Thị', '0984567890', 'dung.lt@gmail.com', N'758 Lũy Bán Bích, Quận Tân Phú, TP.HCM', GETDATE()),
    ('850e8400-e29b-41d4-a716-446655440055', N'Việt', N'Trương Văn', '0985678901', 'viet.tv@gmail.com', N'881 Bình Long, Quận Tân Phú, TP.HCM', GETDATE()),
    ('850e8400-e29b-41d4-a716-446655440056', N'Huyền', N'Dương Thị', '0986789012', 'huyen.dt@gmail.com', N'214 Đào Duy Anh, Quận Phú Nhuận, TP.HCM', GETDATE()),
    ('850e8400-e29b-41d4-a716-446655440057', N'Tân', N'Phan Văn', '0987890123', 'tan.pv@gmail.com', N'347 Phan Đình Phùng, Quận Phú Nhuận, TP.HCM', GETDATE()),
    ('850e8400-e29b-41d4-a716-446655440058', N'Châu', N'Tô Thị', '0988901234', 'chau.tt@gmail.com', N'479 Trần Huy Liệu, Quận Phú Nhuận, TP.HCM', GETDATE()),
    ('850e8400-e29b-41d4-a716-446655440059', N'Phát', N'Cao Văn', '0989012345', 'phat.cv@gmail.com', N'512 Huỳnh Văn Bánh, Quận Phú Nhuận, TP.HCM', GETDATE()),
    ('850e8400-e29b-41d4-a716-44665544005a', N'Trâm', N'Lương Thị', '0980123456', 'tram.lt@gmail.com', N'645 Phạm Ngọc Thạch, Quận 3, TP.HCM', GETDATE()),

    ('850e8400-e29b-41d4-a716-44665544005b', N'Sỹ', N'Đỗ Văn', '0991234567', 'sy.dv@gmail.com', N'778 Lý Chính Thắng, Quận 3, TP.HCM', GETDATE()),
    ('850e8400-e29b-41d4-a716-44665544005c', N'Hà', N'Hà Thị', '0992345678', 'ha.ht@gmail.com', N'811 Trần Quốc Thảo, Quận 3, TP.HCM', GETDATE()),
    ('850e8400-e29b-41d4-a716-44665544005d', N'Hưng', N'Trịnh Văn', '0993456789', 'hungtrinh.tv@gmail.com', N'944 Võ Văn Tần, Quận 3, TP.HCM', GETDATE()),
    ('850e8400-e29b-41d4-a716-44665544005e', N'Uyên', N'Võ Thị', '0994567890', 'uyen.vt@gmail.com', N'177 Cô Giang, Quận 1, TP.HCM', GETDATE()),
    ('850e8400-e29b-41d4-a716-44665544005f', N'Kiệt', N'Lâm Văn', '0995678901', 'kiet.lv@gmail.com', N'209 Ký Con, Quận 1, TP.HCM', GETDATE()),
    ('850e8400-e29b-41d4-a716-446655440060', N'Hạ', N'Đào Thị', '0996789012', 'ha.dt@gmail.com', N'342 Tôn Thất Đạm, Quận 1, TP.HCM', GETDATE()),
    ('850e8400-e29b-41d4-a716-446655440061', N'Minh', N'Mai Văn', '0997890123', 'minhmai.mv@gmail.com', N'475 Lý Tự Trọng, Quận 1, TP.HCM', GETDATE()),
    ('850e8400-e29b-41d4-a716-446655440062', N'Thùy', N'Trần Thị', '0998901234', 'thuy.tt@gmail.com', N'508 Mạc Thị Bưởi, Quận 1, TP.HCM', GETDATE()),
    ('850e8400-e29b-41d4-a716-446655440063', N'Công', N'Nguyễn Văn', '0999012345', 'cong.nv@gmail.com', N'631 Đồng Khởi, Quận 1, TP.HCM', GETDATE()),
    ('850e8400-e29b-41d4-a716-446655440064', N'Thảo', N'Phạm Thị', '0990123456', 'thaopham.pt@gmail.com', N'764 Lê Thánh Tôn, Quận 1, TP.HCM', GETDATE());

INSERT INTO vehicles (id, vin, license_plate, customer_id, model_id, purchase_date, created_at)
VALUES
    ('750e8400-e29b-41d4-a716-446655440001', 'VF8ABC123XYZ45678', '1HGBH41JXMN109186', '850e8400-e29b-41d4-a716-446655440001', @vf8_2023, '2023-03-15', GETDATE()),
    ('750e8400-e29b-41d4-a716-446655440002', 'VF9DEF456ABC78901', '2T3ZFREV8HW345678', '850e8400-e29b-41d4-a716-446655440002', @vf9_2024, '2024-01-20', GETDATE()),
    ('750e8400-e29b-41d4-a716-446655440003', 'VFE34GHJ789DEF234', '5YJSA1E26HF234567', '850e8400-e29b-41d4-a716-446655440003', @vfe34_2023, '2023-07-10', GETDATE()),
    ('750e8400-e29b-41d4-a716-446655440004', 'VF5KLM012GHJ567', '3VW2B7AJ9KM456789', '850e8400-e29b-41d4-a716-446655440004', @vf5_2024, '2024-05-22', GETDATE()),
    ('750e8400-e29b-41d4-a716-446655440005', 'VF6NOP345KLM890', 'WBAJF1C59JA567890', '850e8400-e29b-41d4-a716-446655440005', @vf6_2024, '2024-08-14', GETDATE()),
    ('750e8400-e29b-41d4-a716-446655440006', 'VF7PQR678NOP123', 'KNDCC3L74M7345678', '850e8400-e29b-41d4-a716-446655440006', @vf7_2025, '2025-02-18', GETDATE()),
    ('750e8400-e29b-41d4-a716-446655440007', 'VF3STU901PQR456', 'JM1DKFC75M0234567', '850e8400-e29b-41d4-a716-446655440007', @vf3_2024, '2024-06-30', GETDATE()),
    ('750e8400-e29b-41d4-a716-446655440008', 'BYD8VWX234STU789', '4T1BF1FK7GU456789', '850e8400-e29b-41d4-a716-446655440008', @byd_atto3_2023, '2023-09-12', GETDATE()),
    ('750e8400-e29b-41d4-a716-446655440009', 'BYDDOL567VWX123', 'JTDKARFP8K3234567', '850e8400-e29b-41d4-a716-446655440009', @byd_dolphin_2024, '2024-03-25', GETDATE()),
    ('750e8400-e29b-41d4-a716-44665544000a', 'BYDSEA890YZA567', '1C4RJFBG2JC567890', '850e8400-e29b-41d4-a716-44665544000a', @byd_seal_2023, '2023-11-08', GETDATE()),
    ('750e8400-e29b-41d4-a716-44665544000b', 'BYDTAN123BCD789', '5UXCR6C03L9345678', '850e8400-e29b-41d4-a716-44665544000b', @byd_tang_2024, '2024-04-17', GETDATE()),
    ('750e8400-e29b-41d4-a716-44665544000c', 'BYDSLU456EFG901', 'WDDWJ8GB3KF234567', '850e8400-e29b-41d4-a716-44665544000c', @byd_sealu_2025, '2025-01-09', GETDATE()),
    ('750e8400-e29b-41d4-a716-44665544000d', 'BYDM6J789HJK234', '3GNAXUEV9ML456789', '850e8400-e29b-41d4-a716-44665544000d', @byd_m6_2024, '2024-07-21', GETDATE()),
    ('750e8400-e29b-41d4-a716-44665544000e', 'MBEQA012KLM567', '2HGFC2F59MH567890', '850e8400-e29b-41d4-a716-44665544000e', @mb_eqa_2023, '2023-05-13', GETDATE()),
    ('750e8400-e29b-41d4-a716-44665544000f', 'MBEQB345NOP890', 'WAUFFAFL2HN234567', '850e8400-e29b-41d4-a716-44665544000f', @mb_eqb_2024, '2024-09-05', GETDATE()),
    ('750e8400-e29b-41d4-a716-446655440010', 'MBEQC678PQR123', '1FMCU9GD7KUA45678', '850e8400-e29b-41d4-a716-446655440010', @mb_eqc_2023, '2023-12-19', GETDATE()),
    ('750e8400-e29b-41d4-a716-446655440011', 'MBEQE901STU456', 'KM8J33A49MU567890', '850e8400-e29b-41d4-a716-446655440011', @mb_eqe_2024, '2024-02-28', GETDATE()),
    ('750e8400-e29b-41d4-a716-446655440012', 'MBEQESUV234VWX789', '5NPD84LF9KH234567', '850e8400-e29b-41d4-a716-446655440012', @mb_eqesuv_2025, '2025-06-11', GETDATE()),
    ('750e8400-e29b-41d4-a716-446655440013', 'MBEQS567YZA012', 'JN8AZ2NE2M9345678', '850e8400-e29b-41d4-a716-446655440013', @mb_eqs_2023, '2023-08-07', GETDATE()),
    ('750e8400-e29b-41d4-a716-446655440014', 'MBEQSSUV890BCD345', '1G1ZE5ST8MF456789', '850e8400-e29b-41d4-a716-446655440014', @mb_eqssuv_2024, '2024-10-15', GETDATE()),
    ('750e8400-e29b-41d4-a716-446655440015', 'VF8EFG123HJK678', '3MW39FS05R8567890', '850e8400-e29b-41d4-a716-446655440015', @vf8_2024, '2024-12-03', GETDATE()),
    ('750e8400-e29b-41d4-a716-446655440016', 'VF9KLM456NOP901', 'JTHBK1GG9M2234567', '850e8400-e29b-41d4-a716-446655440016', @vf9_2023, '2023-04-26', GETDATE()),
    ('750e8400-e29b-41d4-a716-446655440017', 'VFE34PQR789STU234', 'WBA3B3C55EK345678', '850e8400-e29b-41d4-a716-446655440017', @vfe34_2024, '2024-11-12', GETDATE()),
    ('750e8400-e29b-41d4-a716-446655440018', 'VF5VWX012YZA567', 'SALWR2TF9KA456789', '850e8400-e29b-41d4-a716-446655440018', @vf5_2025, '2025-03-08', GETDATE()),
    ('750e8400-e29b-41d4-a716-446655440019', 'VF6BCD345EFG890', 'YV4H60CF8M1567890', '850e8400-e29b-41d4-a716-446655440019', @vf6_2025, '2025-07-19', GETDATE()),
    ('750e8400-e29b-41d4-a716-44665544001a', 'VF7HJK678KLM123', 'ZFF79ALA6M0234567', '850e8400-e29b-41d4-a716-44665544001a', @vf7_2024, '2024-01-31', GETDATE()),
    ('750e8400-e29b-41d4-a716-44665544001b', 'VF3NOP901PQR456', '19VDE1F38ME345678', '850e8400-e29b-41d4-a716-44665544001b', @vf3_2025, '2025-05-14', GETDATE()),
    ('750e8400-e29b-41d4-a716-44665544001c', 'BYDATT234STU789', '5YFBURHE7MP456789', '850e8400-e29b-41d4-a716-44665544001c', @byd_atto3_2024, '2024-08-27', GETDATE()),
    ('750e8400-e29b-41d4-a716-44665544001d', 'BYDDOL567VWX012', 'WMWXP7C51KWR67890', '850e8400-e29b-41d4-a716-44665544001d', @byd_dolphin_2025, '2025-09-16', GETDATE()),
    ('750e8400-e29b-41d4-a716-44665544001e', 'BYDSEA890YZA345', 'KNDJT2A20M7234567', '850e8400-e29b-41d4-a716-44665544001e', @byd_seal_2024, '2024-06-05', GETDATE()),
    ('750e8400-e29b-41d4-a716-44665544001f', 'BYDTAN123BCD678', 'ML32A3HJ3MH345678', '850e8400-e29b-41d4-a716-44665544001f', @byd_tang_2025, '2025-10-22', GETDATE()),
    ('750e8400-e29b-41d4-a716-446655440020', 'BYDSLU456EFG900', 'JTMBFREV0MJ456789', '850e8400-e29b-41d4-a716-446655440020', @byd_sealu_2024, '2024-03-11', GETDATE()),
    ('750e8400-e29b-41d4-a716-446655440021', 'BYDM6J789HJL234', '1N4BL4EV8MN567890', '850e8400-e29b-41d4-a716-446655440021', @byd_m6_2025, '2025-04-29', GETDATE()),
    ('750e8400-e29b-41d4-a716-446655440022', 'MBEQA012MMM567', 'WVWZZZ3CZMW234567', '850e8400-e29b-41d4-a716-446655440022', @mb_eqa_2024, '2024-07-07', GETDATE()),
    ('750e8400-e29b-41d4-a716-446655440023', 'MBEQB345NOP891', '2C4RDGCG9MR345678', '850e8400-e29b-41d4-a716-446655440023', @mb_eqb_2025, '2025-08-18', GETDATE()),
    ('750e8400-e29b-41d4-a716-446655440024', 'MBEQC678PQR120', 'SAJDA1CH1MCP56789', '850e8400-e29b-41d4-a716-446655440024', @mb_eqc_2024, '2024-02-14', GETDATE()),
    ('750e8400-e29b-41d4-a716-446655440025', 'MBEQE901STU458', 'WBA8E9C50KAK67890', '850e8400-e29b-41d4-a716-446655440025', @mb_eqe_2025, '2025-01-06', GETDATE()),
    ('750e8400-e29b-41d4-a716-446655440026', 'MBEQESUV234VWX780', '5YMGZ0C21M0234567', '850e8400-e29b-41d4-a716-446655440026', @mb_eqesuv_2024, '2024-05-23', GETDATE()),
    ('750e8400-e29b-41d4-a716-446655440027', 'MBEQS567YZA010', '1FTEW1E50MKE45678', '850e8400-e29b-41d4-a716-446655440027', @mb_eqs_2024, '2024-09-30', GETDATE()),
    ('750e8400-e29b-41d4-a716-446655440028', 'MBEQSSUV890BCD045', 'KNDMC3LD6M6456789', '850e8400-e29b-41d4-a716-446655440028', @mb_eqssuv_2025, '2025-02-25', GETDATE()),
    ('750e8400-e29b-41d4-a716-446655440029', 'VF8GHJ345KLM078', 'JN1AZ4EH4MM567890', '850e8400-e29b-41d4-a716-446655440029', @vf8_2025, '2025-06-13', GETDATE()),
    ('750e8400-e29b-41d4-a716-44665544002a', 'VF9NOP678PQR001', '5YJSA1E14HF234567', '850e8400-e29b-41d4-a716-44665544002a', @vf9_2025, '2025-03-19', GETDATE()),
    ('750e8400-e29b-41d4-a716-44665544002b', 'VFE34STU901VWX034', 'WDDUX8EB7MA345678', '850e8400-e29b-41d4-a716-44665544002b', @vfe34_2025, '2025-07-04', GETDATE()),
    ('750e8400-e29b-41d4-a716-44665544002c', 'VF5YZA234BCD067', '1G1BE5SM2M7456789', '850e8400-e29b-41d4-a716-44665544002c', @vf5_2023, '2023-10-28', GETDATE()),
    ('750e8400-e29b-41d4-a716-44665544002d', 'VF6EFG567HJK090', 'ZACCJBBN0MPZ67890', '850e8400-e29b-41d4-a716-44665544002d', @vf6_2024, '2024-12-09', GETDATE()),
    ('750e8400-e29b-41d4-a716-44665544002e', 'VF7KLM890NOP023', 'JM3TCBCY9M0234567', '850e8400-e29b-41d4-a716-44665544002e', @vf7_2025, '2025-04-16', GETDATE()),
    ('750e8400-e29b-41d4-a716-44665544002f', 'VF3PQR123STU056', '3CZRU6H78MG345678', '850e8400-e29b-41d4-a716-44665544002f', @vf3_2024, '2024-11-21', GETDATE()),
    ('750e8400-e29b-41d4-a716-446655440030', 'BYDATT789VWX112', 'YV1A22MK4M2456789', '850e8400-e29b-41d4-a716-446655440030', @byd_atto3_2025, '2025-08-03', GETDATE()),
    ('750e8400-e29b-41d4-a716-446655440031', 'BYDDOL012YZA045', 'KL1TD6DE3MB567890', '850e8400-e29b-41d4-a716-446655440031', @byd_dolphin_2023, '2023-06-17', GETDATE()),
    ('750e8400-e29b-41d4-a716-446655440032', 'BYDSEA345BCD078', 'JTDZN3EU7MJ234567', '850e8400-e29b-41d4-a716-446655440032', @byd_seal_2025, '2025-09-28', GETDATE()),
    ('750e8400-e29b-41d4-a716-446655440033', 'VF8VWX678EFG001', 'WP0AA2A98MS345678', '850e8400-e29b-41d4-a716-446655440033', @vf8_2023, '2023-02-11', GETDATE()),
    ('750e8400-e29b-41d4-a716-446655440034', 'VF9YZA901HJK034', '1FMCU9HD3MUA56789', '850e8400-e29b-41d4-a716-446655440034', @vf9_2024, '2024-05-08', GETDATE()),
    ('750e8400-e29b-41d4-a716-446655440035', 'VFE34BCD234KLM567', 'KNDJ23AU2M7567890', '850e8400-e29b-41d4-a716-446655440035', @vfe34_2022, '2023-01-24', GETDATE()),
    ('750e8400-e29b-41d4-a716-446655440036', 'VF5EFG567NOP890', 'LRBFXBSA4MD234567', '850e8400-e29b-41d4-a716-446655440036', @vf5_2024, '2024-08-20', GETDATE()),
    ('750e8400-e29b-41d4-a716-446655440037', 'VF6HJK890PQR123', '2GNAXUEV2M6345678', '850e8400-e29b-41d4-a716-446655440037', @vf6_2025, '2025-10-05', GETDATE()),
    ('750e8400-e29b-41d4-a716-446655440038', 'VF7KLM123STU456', 'VF1RFB00567456789', '850e8400-e29b-41d4-a716-446655440038', @vf7_2024, '2024-03-14', GETDATE()),
    ('750e8400-e29b-41d4-a716-446655440039', 'VF3NOP456VWX789', 'ZHWGJ6AT1MLA67890', '850e8400-e29b-41d4-a716-446655440039', @vf3_2025, '2025-07-27', GETDATE()),
    ('750e8400-e29b-41d4-a716-44665544003a', 'BYDATT789YZA012', '5YJSA1E40KF234567', '850e8400-e29b-41d4-a716-44665544003a', @byd_atto3_2022, '2023-03-02', GETDATE()),
    ('750e8400-e29b-41d4-a716-44665544003b', 'BYDDOL012BCD345', 'WBA5A5C53MD345678', '850e8400-e29b-41d4-a716-44665544003b', @byd_dolphin_2024, '2024-06-18', GETDATE()),
    ('750e8400-e29b-41d4-a716-44665544003c', 'BYDSEA345EFG678', 'WAUENAF44MN456789', '850e8400-e29b-41d4-a716-44665544003c', @byd_seal_2023, '2023-11-30', GETDATE()),
    ('750e8400-e29b-41d4-a716-44665544003d', 'BYDTAN678HJK901', '1HGCV1F36MA567890', '850e8400-e29b-41d4-a716-44665544003d', @byd_tang_2024, '2024-09-11', GETDATE()),
    ('750e8400-e29b-41d4-a716-44665544003e', 'BYDSLU901KLM234', 'JF2GTALC0MH234567', '850e8400-e29b-41d4-a716-44665544003e', @byd_sealu_2025, '2025-01-23', GETDATE()),
    ('750e8400-e29b-41d4-a716-44665544003f', 'BYDM6J234NOP567', '3VW2B7AJ8MM345678', '850e8400-e29b-41d4-a716-44665544003f', @byd_m6_2024, '2024-04-07', GETDATE()),
    ('750e8400-e29b-41d4-a716-446655440040', 'MBEQA567PQR890', 'WDDGF8BB9MR456789', '850e8400-e29b-41d4-a716-446655440040', @mb_eqa_2025, '2025-05-19', GETDATE()),
    ('750e8400-e29b-41d4-a716-446655440041', 'MBEQB890STU123', 'SALGA2EV7MA567890', '850e8400-e29b-41d4-a716-446655440041', @mb_eqb_2023, '2023-08-25', GETDATE()),
    ('750e8400-e29b-41d4-a716-446655440042', 'MBEQC123VWX456', 'JN1BJ1CR8MW234567', '850e8400-e29b-41d4-a716-446655440042', @mb_eqc_2022, '2023-12-13', GETDATE()),
    ('750e8400-e29b-41d4-a716-446655440043', 'MBEQE456YZA789', 'WP1AE2A28MLA45678', '850e8400-e29b-41d4-a716-446655440043', @mb_eqe_2023, '2023-04-08', GETDATE()),
    ('750e8400-e29b-41d4-a716-446655440044', 'MBEQESUV789BCD012', '1C4RJFAG9MC456789', '850e8400-e29b-41d4-a716-446655440044', @mb_eqesuv_2024, '2024-10-01', GETDATE()),
    ('750e8400-e29b-41d4-a716-446655440045', 'MBEQS012EFG345', '5UXCR4C04M9567890', '850e8400-e29b-41d4-a716-446655440045', @mb_eqs_2025, '2025-02-16', GETDATE()),
    ('750e8400-e29b-41d4-a716-446655440046', 'MBEQSSUV345HJK678', 'WAUZZZ4G4MN234567', '850e8400-e29b-41d4-a716-446655440046', @mb_eqssuv_2023, '2023-07-22', GETDATE()),
    ('750e8400-e29b-41d4-a716-446655440047', 'VF8KLM678NOP901', '2T3P1RFV8MW345678', '850e8400-e29b-41d4-a716-446655440047', @vf8_2024, '2024-11-06', GETDATE()),
    ('750e8400-e29b-41d4-a716-446655440048', 'VF9PQR901STU234', '5YJSA1E22MF456789', '850e8400-e29b-41d4-a716-446655440048', @vf9_2025, '2025-03-29', GETDATE()),
    ('750e8400-e29b-41d4-a716-446655440049', 'VFE34VWX234YZA567', 'KNDJT2A60M7567890', '850e8400-e29b-41d4-a716-446655440049', @vfe34_2023, '2023-09-14', GETDATE()),
    ('750e8400-e29b-41d4-a716-44665544004a', 'VF5BCD567EFG890', 'ML32A3HJ8MH234567', '850e8400-e29b-41d4-a716-44665544004a', @vf5_2025, '2025-06-21', GETDATE()),
    ('750e8400-e29b-41d4-a716-44665544004b', 'VF6GHJ890KLM123', 'JTMBFREV4MJ345678', '850e8400-e29b-41d4-a716-44665544004b', @vf6_2024, '2024-01-18', GETDATE()),
    ('750e8400-e29b-41d4-a716-44665544004c', 'VF7NOP123PQR456', '1N4BL4EV2MN456789', '850e8400-e29b-41d4-a716-44665544004c', @vf7_2025, '2025-08-09', GETDATE()),
    ('750e8400-e29b-41d4-a716-44665544004d', 'VF3STU456VWX789', 'WVWZZZ3CZMW567890', '850e8400-e29b-41d4-a716-44665544004d', @vf3_2024, '2024-04-25', GETDATE()),
    ('750e8400-e29b-41d4-a716-44665544004e', 'BYDATT789YZA112', '2C4RDGCG3MR234567', '850e8400-e29b-41d4-a716-44665544004e', @byd_atto3_2023, '2023-05-31', GETDATE()),
    ('750e8400-e29b-41d4-a716-44665544004f', 'BYDDOL012BCD045', 'SAJDA1CH5MCP45678', '850e8400-e29b-41d4-a716-44665544004f', @byd_dolphin_2025, '2025-10-17', GETDATE()),
    ('750e8400-e29b-41d4-a716-446655440050', 'BYDSEA345EFG078', 'WBA8E9C54KAK56789', '850e8400-e29b-41d4-a716-446655440050', @byd_seal_2024, '2024-02-03', GETDATE()),
    ('750e8400-e29b-41d4-a716-446655440051', 'BYDTAN678HJK001', '5YMGZ0C25M0567890', '850e8400-e29b-41d4-a716-446655440051', @byd_tang_2025, '2025-07-15', GETDATE()),
    ('750e8400-e29b-41d4-a716-446655440052', 'BYDSLU901KLM034', '1FTEW1E54MKE34567', '850e8400-e29b-41d4-a716-446655440052', @byd_sealu_2024, '2024-12-28', GETDATE()),
    ('750e8400-e29b-41d4-a716-446655440053', 'BYDM6J234NOP067', 'KNDMC3LD0M6345678', '850e8400-e29b-41d4-a716-446655440053', @byd_m6_2025, '2025-04-11', GETDATE()),
    ('750e8400-e29b-41d4-a716-446655440054', 'MBEQA567PQR800', 'JN1AZ4EH8MM456789', '850e8400-e29b-41d4-a716-446655440054', @mb_eqa_2023, '2023-06-06', GETDATE()),
    ('750e8400-e29b-41d4-a716-446655440055', 'MBEQB890STU103', '5YJSA1E18HF567890', '850e8400-e29b-41d4-a716-446655440055', @mb_eqb_2024, '2024-09-19', GETDATE()),
    ('750e8400-e29b-41d4-a716-446655440056', 'MBEQC123VWX406', 'WDDUX8EB1MA234567', '850e8400-e29b-41d4-a716-446655440056', @mb_eqc_2023, '2023-10-12', GETDATE()),
    ('750e8400-e29b-41d4-a716-446655440057', 'MBEQE456YZA709', '1G1BE5SM6M7345678', '850e8400-e29b-41d4-a716-446655440057', @mb_eqe_2024, '2024-05-27', GETDATE()),
    ('750e8400-e29b-41d4-a716-446655440058', 'MBEQESUV789BCD002', 'ZACCJBBN4MPZ56789', '850e8400-e29b-41d4-a716-446655440058', @mb_eqesuv_2025, '2025-01-14', GETDATE()),
    ('750e8400-e29b-41d4-a716-446655440059', 'MBEQS012EFG305', 'JM3TCBCY3M0567890', '850e8400-e29b-41d4-a716-446655440059', @mb_eqs_2024, '2024-07-30', GETDATE()),
    ('750e8400-e29b-41d4-a716-44665544005a', 'MBEQSSUV345HJK608', '3CZRU6H72MG234567', '850e8400-e29b-41d4-a716-44665544005a', @mb_eqssuv_2024, '2024-11-24', GETDATE()),
    ('750e8400-e29b-41d4-a716-44665544005b', 'VF8NOP678PQR001', 'YV1A22MK8M2345678', '850e8400-e29b-41d4-a716-44665544005b', @vf8_2025, '2025-02-07', GETDATE()),
    ('750e8400-e29b-41d4-a716-44665544005c', 'VF9STU901VWX034', 'KL1TD6DE7MB456789', '850e8400-e29b-41d4-a716-44665544005c', @vf9_2023, '2023-03-21', GETDATE()),
    ('750e8400-e29b-41d4-a716-44665544005d', 'VFE34YZA234BCD067', 'JTDZN3EU1MJ567890', '850e8400-e29b-41d4-a716-44665544005d', @vfe34_2024, '2024-06-14', GETDATE()),
    ('750e8400-e29b-41d4-a716-44665544005e', 'VF5EFG567HJK899', 'WP0AA2A92MS234567', '850e8400-e29b-41d4-a716-44665544005e', @vf5_2023, '2023-08-18', GETDATE()),
    ('750e8400-e29b-41d4-a716-44665544005f', 'VF6KLM890NOP111', '1FMCU9HD7MUA45678', '850e8400-e29b-41d4-a716-44665544005f', @vf6_2025, '2025-09-02', GETDATE()),
    ('750e8400-e29b-41d4-a716-446655440060', 'VF7PQR123STU450', 'KNDJ23AU6M7456789', '850e8400-e29b-41d4-a716-446655440060', @vf7_2024, '2024-10-26', GETDATE()),
    ('750e8400-e29b-41d4-a716-446655440061', 'VF3VWX456YZA780', 'LRBFXBSA8MD567890', '850e8400-e29b-41d4-a716-446655440061', @vf3_2025, '2025-03-12', GETDATE()),
    ('750e8400-e29b-41d4-a716-446655440062', 'BYDATT789BCD010', '2GNAXUEV6M6234567', '850e8400-e29b-41d4-a716-446655440062', @byd_atto3_2024, '2024-05-04', GETDATE()),
    ('750e8400-e29b-41d4-a716-446655440063', 'BYDDOL012EFG340', 'VF1RFB00567345678', '850e8400-e29b-41d4-a716-446655440063', @byd_dolphin_2023, '2023-12-17', GETDATE()),
    ('750e8400-e29b-41d4-a716-446655440064', 'BYDSEA345HJK670', 'ZHWGJ6AT5MLA56789', '850e8400-e29b-41d4-a716-446655440064', @byd_seal_2025, '2025-06-28', GETDATE())

INSERT INTO part_categories (id, category_name, description, parent_category_id, created_at)
VALUES
    ('150e8400-e29b-41d4-a716-446655440001', 'VF8_2025_CATEGORIES', N'Nhóm phụ tùng chính cho dòng xe VinFast VF8 2025', NULL, GETDATE()),
    ('150e8400-e29b-41d4-a716-446655440002', 'VFE34_2025_CATEGORIES', N'Nhóm phụ tùng chính cho dòng xe VinFast VF e34 2025', NULL, GETDATE());
-- VF8_2025_CATEGORIES
INSERT INTO part_categories (id, category_name, description, parent_category_id, created_at)
VALUES
    ('150e8400-e29b-41d4-a716-446655441111', 'VF8_2025_Powertrain & Battery', N'Nhóm động cơ, pin và hệ thống truyền động cho VF8 2025', '150e8400-e29b-41d4-a716-446655440001', GETDATE()),
    ('150e8400-e29b-41d4-a716-446655441222', 'VF8_2025_Exterior Body & Trim', N'Nhóm thân xe, ngoại thất, ốp – trim cho VF8 2025', '150e8400-e29b-41d4-a716-446655440001', GETDATE()),
    ('150e8400-e29b-41d4-a716-446655441333', 'VF8_2025_Interior & Comfort', N'Nhóm nội thất, ghế, vô-lăng, tiện nghi cho VF8 2025', '150e8400-e29b-41d4-a716-446655440001', GETDATE()),
    ('150e8400-e29b-41d4-a716-446655441444', 'VF8_2025_Electronics & Accessories', N'Nhóm điện tử, camera, cảm biến, phụ kiện cho VF8 2025', '150e8400-e29b-41d4-a716-446655440001', GETDATE()),
    ('150e8400-e29b-41d4-a716-446655441555', 'VF8_2025_Charging & Infrastructure', N'Nhóm trạm sạc, dây sạc, phụ kiện sạc cho VF8 2025', '150e8400-e29b-41d4-a716-446655440001', GETDATE());
-- VFE34_2025_CATEGORIES
INSERT INTO part_categories (id, category_name, description, parent_category_id, created_at)
VALUES
    ('150e8400-e29b-41d4-a716-446655442111', 'VFe34_2025_Powertrain & Battery', N'Nhóm động cơ, pin và hệ thống truyền động cho VF e34 2025', '150e8400-e29b-41d4-a716-446655440002', GETDATE()),
    ('150e8400-e29b-41d4-a716-446655442222', 'VFe34_2025_Exterior Body & Trim', N'Nhóm thân xe, ngoại thất, ốp – trim cho VF e34 2025', '150e8400-e29b-41d4-a716-446655440002', GETDATE()),
    ('150e8400-e29b-41d4-a716-446655442333', 'VFe34_2025_Interior & Comfort', N'Nhóm nội thất, ghế, vô-lăng, tiện nghi cho VF e34 2025', '150e8400-e29b-41d4-a716-446655440002', GETDATE()),
    ('150e8400-e29b-41d4-a716-446655442444', 'VFe34_2025_Electronics & Accessories', N'Nhóm điện tử, camera, cảm biến, phụ kiện cho VF e34 2025', '150e8400-e29b-41d4-a716-446655440002', GETDATE()),
    ('150e8400-e29b-41d4-a716-446655442555', 'VFe34_2025_Charging & Infrastructure', N'Nhóm trạm sạc, dây sạc, phụ kiện sạc cho VF e34 2025', '150e8400-e29b-41d4-a716-446655440002', GETDATE());

-- Cấp 1: Các danh mục con trực tiếp của POWERTRAIN & BATTERY
INSERT INTO part_categories (id, category_name, description, parent_category_id, created_at)
VALUES
    ('150e8400-e29b-41d4-a716-446655440012', N'Electric Motor', N'Động cơ điện', '150e8400-e29b-41d4-a716-446655441111', GETDATE()),
    ('150e8400-e29b-41d4-a716-446655440013', N'Battery Pack', N'Bộ pin', '150e8400-e29b-41d4-a716-446655441111', GETDATE()),
    ('150e8400-e29b-41d4-a716-446655440014', N'Transmission', N'Hệ thống truyền động', '150e8400-e29b-41d4-a716-446655441111', GETDATE()),
    ('150e8400-e29b-41d4-a716-446655440015', N'Cooling System', N'Hệ thống làm mát', '150e8400-e29b-41d4-a716-446655441111', GETDATE());

-- Cấp 2: Các danh mục con của Electric Motor
INSERT INTO part_categories (id, category_name, description, parent_category_id, created_at)
VALUES
    ('150e8400-e29b-41d4-a716-446655440016', N'Motor Assembly', N'Cụm động cơ hoàn chỉnh', '150e8400-e29b-41d4-a716-446655440012', GETDATE()),
    ('150e8400-e29b-41d4-a716-446655440017', N'Motor Components', N'Linh kiện cấu thành động cơ', '150e8400-e29b-41d4-a716-446655440012', GETDATE());

-- Cấp 3: Các danh mục con của Motor Assembly
INSERT INTO part_categories (id, category_name, description, parent_category_id, created_at)
VALUES
    ('150e8400-e29b-41d4-a716-446655440018', N'Front Motor Unit', N'Khối động cơ trước', '150e8400-e29b-41d4-a716-446655440016', GETDATE()),
    ('150e8400-e29b-41d4-a716-446655440019', N'Rear Motor Unit', N'Khối động cơ sau', '150e8400-e29b-41d4-a716-446655440016', GETDATE());

-- Cấp 4: Các mục con của Front Motor Unit
INSERT INTO part_categories (id, category_name, description, parent_category_id, created_at)
VALUES
    ('150e8400-e29b-41d4-a716-44665544001A', 'VF8.EM.FM.150KW.ECO', N'Động cơ trước 150kW cho phiên bản Eco.', '150e8400-e29b-41d4-a716-446655440018', GETDATE()),
    ('150e8400-e29b-41d4-a716-44665544001B', 'VF8.EM.FM.150KW.PLUS', N'Động cơ trước 150kW cho phiên bản Plus.', '150e8400-e29b-41d4-a716-446655440018', GETDATE()),
    ('150e8400-e29b-41d4-a716-44665544001C', 'VF8.EM.FM.150KW.SVC.NEW', N'Phụ tùng thay thế mới.', '150e8400-e29b-41d4-a716-446655440018', GETDATE()),
    ('150e8400-e29b-41d4-a716-44665544001D', 'VF8.EM.FM.150KW.SVC.RMD', N'Phụ tùng thay thế tái sản xuất.', '150e8400-e29b-41d4-a716-446655440018', GETDATE());

-- Cấp 4: Các mục con của Rear Motor Unit
INSERT INTO part_categories (id, category_name, description, parent_category_id, created_at)
VALUES
    ('150e8400-e29b-41d4-a716-44665544001E', 'VF8.EM.RM.150KW.ECO', N'Động cơ sau 150kW cho phiên bản Eco.', '150e8400-e29b-41d4-a716-446655440019', GETDATE()),
    ('150e8400-e29b-41d4-a716-44665544001F', 'VF8.EM.RM.150KW.PLUS', N'Động cơ sau 150kW cho phiên bản Plus.', '150e8400-e29b-41d4-a716-446655440019', GETDATE()),
    ('150e8400-e29b-41d4-a716-446655440020', 'VF8.EM.RM.150KW.SVC.NEW', N'Phụ tùng thay thế mới.', '150e8400-e29b-41d4-a716-446655440019', GETDATE()),
    ('150e8400-e29b-41d4-a716-446655440021', 'VF8.EM.RM.150KW.SVC.RMD', N'Phụ tùng thay thế tái sản xuất.', '150e8400-e29b-41d4-a716-446655440019', GETDATE());

-- Cấp 3: Các danh mục con của Motor Components
INSERT INTO part_categories (id, category_name, description, parent_category_id, created_at)
VALUES
    ('150e8400-e29b-41d4-a716-446655440022', N'Inverter/Power Electronics Module', N'Bộ biến tần/Mô-đun điện tử công suất', '150e8400-e29b-41d4-a716-446655440017', GETDATE()),
    ('150e8400-e29b-41d4-a716-446655440023', N'Rotor Assembly', N'Cụm Rô-to', '150e8400-e29b-41d4-a716-446655440017', GETDATE()),
    ('150e8400-e29b-41d4-a716-446655440024', N'Stator Assembly', N'Cụm Stato', '150e8400-e29b-41d4-a716-446655440017', GETDATE()),
    ('150e8400-e29b-41d4-a716-446655440025', N'Motor Position Sensor', N'Cảm biến vị trí động cơ', '150e8400-e29b-41d4-a716-446655440017', GETDATE());

-- Cấp 4: Các mục con của Inverter/Power Electronics Module
INSERT INTO part_categories (id, category_name, description, parent_category_id, created_at)
VALUES
    ('150e8400-e29b-41d4-a716-446655440026', 'PEM.400V.V1.BOSCH', N'Biến tần 400V, thế hệ 1, nhà cung cấp Bosch.', '150e8400-e29b-41d4-a716-446655440022', GETDATE()),
    ('150e8400-e29b-41d4-a716-446655440027', 'PEM.400V.V2.VISTEON', N'Biến tần 400V, thế hệ 2, nhà cung cấp Visteon.', '150e8400-e29b-41d4-a716-446655440022', GETDATE()),
    ('150e8400-e29b-41d4-a716-446655440028', 'PEM.400V.V2.1.SVC', N'Biến tần 400V, thế hệ 2.1, phụ tùng thay thế.', '150e8400-e29b-41d4-a716-446655440022', GETDATE());

-- Cấp 4: Các mục con của Rotor Assembly
INSERT INTO part_categories (id, category_name, description, parent_category_id, created_at)
VALUES
    ('150e8400-e29b-41d4-a716-446655440029', 'ROT.ASM.150KW.G1', N'Cụm Rô-to cho động cơ 150kW thế hệ 1.', '150e8400-e29b-41d4-a716-446655440023', GETDATE()),
    ('150e8400-e29b-41d4-a716-44665544002A', 'ROT.ASM.150KW.G2', N'Cụm Rô-to cho động cơ 150kW thế hệ 2 (cải tiến).', '150e8400-e29b-41d4-a716-446655440023', GETDATE()),
    ('150e8400-e29b-41d4-a716-44665544002B', 'ROT.ASM.150KW.G1.SVC', N'Phụ tùng thay thế cho thế hệ 1.', '150e8400-e29b-41d4-a716-446655440023', GETDATE());

-- Cấp 4: Các mục con của Stator Assembly
INSERT INTO part_categories (id, category_name, description, parent_category_id, created_at)
VALUES
    ('150e8400-e29b-41d4-a716-44665544002C', 'STA.ASM.150KW.G1', N'Cụm Stato cho động cơ 150kW thế hệ 1.', '150e8400-e29b-41d4-a716-446655440024', GETDATE()),
    ('150e8400-e29b-41d4-a716-44665544002D', 'STA.ASM.150KW.G2', N'Cụm Stato cho động cơ 150kW thế hệ 2 (cải tiến).', '150e8400-e29b-41d4-a716-446655440024', GETDATE()),
    ('150e8400-e29b-41d4-a716-44665544002E', 'STA.ASM.150KW.G1.SVC', N'Phụ tùng thay thế cho thế hệ 1.', '150e8400-e29b-41d4-a716-446655440024', GETDATE());

-- Cấp 4: Các mục con của Motor Position Sensor
INSERT INTO part_categories (id, category_name, description, parent_category_id, created_at)
VALUES
    ('150e8400-e29b-41d4-a716-44665544002F', 'SENS.POS.RESOLVER.V1', N'Cảm biến vị trí loại Resolver, phiên bản 1.', '150e8400-e29b-41d4-a716-446655440025', GETDATE()),
    ('150e8400-e29b-41d4-a716-446655440030', 'SENS.POS.RESOLVER.V2.HELLA', N'Cảm biến vị trí loại Resolver, phiên bản 2, nhà cung cấp Hella.', '150e8400-e29b-41d4-a716-446655440025', GETDATE()),
    ('150e8400-e29b-41d4-a716-446655440031', 'SENS.POS.HALL.V3.BOSCH', N'Cảm biến vị trí loại Hall Effect, phiên bản 3, nhà cung cấp Bosch.', '150e8400-e29b-41d4-a716-446655440025', GETDATE());

-- Cấp 2: Các danh mục con của Battery Pack
INSERT INTO part_categories (id, category_name, description, parent_category_id, created_at)
VALUES
    ('150e8400-e29b-41d4-a716-446655440032', N'High-Voltage Battery Module', N'Mô-đun pin cao áp', '150e8400-e29b-41d4-a716-446655440013', GETDATE()),
    ('150e8400-e29b-41d4-a716-446655440033', N'Battery Management System', N'Hệ thống quản lý pin - BMS', '150e8400-e29b-41d4-a716-446655440013', GETDATE()),
    ('150e8400-e29b-41d4-a716-446655440034', N'Battery Housing & Structure', N'Vỏ và kết cấu pin', '150e8400-e29b-41d4-a716-446655440013', GETDATE()),
    ('150e8400-e29b-41d4-a716-446655440035', N'High-Voltage Junction Box', N'Hộp cầu nối cao áp - HVJB', '150e8400-e29b-41d4-a716-446655440013', GETDATE());

-- Cấp 3: Các danh mục con của High-Voltage Battery Module
INSERT INTO part_categories (id, category_name, description, parent_category_id, created_at)
VALUES
    ('150e8400-e29b-41d4-a716-446655440036', N'Battery Cell', N'Cell pin', '150e8400-e29b-41d4-a716-446655440032', GETDATE()),
    ('150e8400-e29b-41d4-a716-446655440037', N'Module Interconnect Board', N'Bảng mạch kết nối mô-đun', '150e8400-e29b-41d4-a716-446655440032', GETDATE());

-- Cấp 4: Các mục con của Battery Cell
INSERT INTO part_categories (id, category_name, description, parent_category_id, created_at)
VALUES
    ('150e8400-e29b-41d4-a716-446655440038', 'CELL.NMC.82KWH.CATL', N'Cell pin NMC cho bộ pin 82kWh, nhà cung cấp CATL.', '150e8400-e29b-41d4-a716-446655440036', GETDATE()),
    ('150e8400-e29b-41d4-a716-446655440039', 'CELL.NMC.87KWH.SDI', N'Cell pin NMC cho bộ pin 87.7kWh, nhà cung cấp Samsung SDI.', '150e8400-e29b-41d4-a716-446655440036', GETDATE()),
    ('150e8400-e29b-41d4-a716-44665544003A', 'CELL.NMC.87KWH.CATL.G2', N'Cell pin NMC thế hệ 2 cho bộ pin 87.7kWh, nhà cung cấp CATL.', '150e8400-e29b-41d4-a716-446655440036', GETDATE()),
    ('150e8400-e29b-41d4-a716-44665544003B', 'CELL.LFP.75KWH.GOTION', N'Cell pin LFP cho bộ pin 75kWh (phiên bản tương lai), nhà cung cấp Gotion.', '150e8400-e29b-41d4-a716-446655440036', GETDATE());

-- Cấp 4: Các mục con của Module Interconnect Board
INSERT INTO part_categories (id, category_name, description, parent_category_id, created_at)
VALUES
    ('150e8400-e29b-41d4-a716-44665544003C', 'MIB.82KWH.V1', N'Bảng mạch cho bộ pin 82kWh, phiên bản 1.', '150e8400-e29b-41d4-a716-446655440037', GETDATE()),
    ('150e8400-e29b-41d4-a716-44665544003D', 'MIB.87KWH.V1', N'Bảng mạch cho bộ pin 87.7kWh, phiên bản 1.', '150e8400-e29b-41d4-a716-446655440037', GETDATE()),
    ('150e8400-e29b-41d4-a716-44665544003E', 'MIB.87KWH.V2.ENH', N'Bảng mạch cho bộ pin 87.7kWh, phiên bản 2, cải tiến.', '150e8400-e29b-41d4-a716-446655440037', GETDATE());

-- Cấp 3: Các danh mục con của Battery Management System
INSERT INTO part_categories (id, category_name, description, parent_category_id, created_at)
VALUES
    ('150e8400-e29b-41d4-a716-44665544003F', N'Master Control Unit', N'Bộ điều khiển chính', '150e8400-e29b-41d4-a716-446655440033', GETDATE()),
    ('150e8400-e29b-41d4-a716-446655440040', N'Slave/Cell Monitoring Unit', N'Bộ giám sát cell pin', '150e8400-e29b-41d4-a716-446655440033', GETDATE());

-- Cấp 4: Các mục con của Master Control Unit
INSERT INTO part_categories (id, category_name, description, parent_category_id, created_at)
VALUES
    ('150e8400-e29b-41d4-a716-446655440041', 'BMS.MCU.HW1.FW25', N'Phần cứng V1, Firmware 2.5.', '150e8400-e29b-41d4-a716-44665544003F', GETDATE()),
    ('150e8400-e29b-41d4-a716-446655440042', 'BMS.MCU.HW2.FW31', N'Phần cứng V2, Firmware 3.1.', '150e8400-e29b-41d4-a716-44665544003F', GETDATE()),
    ('150e8400-e29b-41d4-a716-446655440043', 'BMS.MCU.HW2.FW33.OTA', N'Phần cứng V2, Firmware 3.3 hỗ trợ cập nhật từ xa.', '150e8400-e29b-41d4-a716-44665544003F', GETDATE());

-- Cấp 4: Các mục con của Slave/Cell Monitoring Unit
INSERT INTO part_categories (id, category_name, description, parent_category_id, created_at)
VALUES
    ('150e8400-e29b-41d4-a716-446655440044', 'BMS.CMU.HW1.82KWH', N'Bộ giám sát cho pin 82kWh.', '150e8400-e29b-41d4-a716-446655440040', GETDATE()),
    ('150e8400-e29b-41d4-a716-446655440045', 'BMS.CMU.HW2.87KWH', N'Bộ giám sát cho pin 87.7kWh.', '150e8400-e29b-41d4-a716-446655440040', GETDATE()),
    ('150e8400-e29b-41d4-a716-446655440046', 'BMS.CMU.HW2.SVC', N'Phụ tùng thay thế cho phần cứng V2.', '150e8400-e29b-41d4-a716-446655440040', GETDATE());

-- Cấp 3: Các danh mục con của Battery Housing & Structure
INSERT INTO part_categories (id, category_name, description, parent_category_id, created_at)
VALUES
    ('150e8400-e29b-41d4-a716-446655440047', N'Upper Battery Casing', N'Vỏ pin trên', '150e8400-e29b-41d4-a716-446655440034', GETDATE()),
    ('150e8400-e29b-41d4-a716-446655440048', N'Lower Battery Tray', N'Khay pin dưới', '150e8400-e29b-41d4-a716-446655440034', GETDATE()),
    ('150e8400-e29b-41d4-a716-446655440049', N'Sealing Gaskets', N'Gioăng làm kín', '150e8400-e29b-41d4-a716-446655440034', GETDATE());

-- Cấp 4: Các mục con của Upper Battery Casing
INSERT INTO part_categories (id, category_name, description, parent_category_id, created_at)
VALUES
    ('150e8400-e29b-41d4-a716-44665544004A', 'BAT.CAS.UPR.82KWH', N'Vỏ trên cho pin 82kWh.', '150e8400-e29b-41d4-a716-446655440047', GETDATE()),
    ('150e8400-e29b-41d4-a716-44665544004B', 'BAT.CAS.UPR.87KWH', N'Vỏ trên cho pin 87.7kWh.', '150e8400-e29b-41d4-a716-446655440047', GETDATE()),
    ('150e8400-e29b-41d4-a716-44665544004C', 'BAT.CAS.UPR.87KWH.AL', N'Vỏ trên bằng nhôm cho pin 87.7kWh.', '150e8400-e29b-41d4-a716-446655440047', GETDATE());

-- Cấp 4: Các mục con của Lower Battery Tray
INSERT INTO part_categories (id, category_name, description, parent_category_id, created_at)
VALUES
    ('150e8400-e29b-41d4-a716-44665544004D', 'BAT.TRAY.LWR.82KWH', N'Khay dưới cho pin 82kWh.', '150e8400-e29b-41d4-a716-446655440048', GETDATE()),
    ('150e8400-e29b-41d4-a716-44665544004E', 'BAT.TRAY.LWR.87KWH', N'Khay dưới cho pin 87.7kWh.', '150e8400-e29b-41d4-a716-446655440048', GETDATE()),
    ('150e8400-e29b-41d4-a716-44665544004F', 'BAT.TRAY.LWR.87KWH.REINF', N'Khay dưới cho pin 87.7kWh, có gia cố.', '150e8400-e29b-41d4-a716-446655440048', GETDATE());

-- Cấp 4: Các mục con của Sealing Gaskets
INSERT INTO part_categories (id, category_name, description, parent_category_id, created_at)
VALUES
    ('150e8400-e29b-41d4-a716-446655440050', 'GASKET.BAT.82KWH', N'Bộ gioăng cho pin 82kWh.', '150e8400-e29b-41d4-a716-446655440049', GETDATE()),
    ('150e8400-e29b-41d4-a716-446655440051', 'GASKET.BAT.87KWH', N'Bộ gioăng cho pin 87.7kWh.', '150e8400-e29b-41d4-a716-446655440049', GETDATE()),
    ('150e8400-e29b-41d4-a716-446655440052', 'GASKET.BAT.SVC.KIT', N'Bộ gioăng thay thế dịch vụ.', '150e8400-e29b-41d4-a716-446655440049', GETDATE());

-- Cấp 3: Các danh mục con của High-Voltage Junction Box
INSERT INTO part_categories (id, category_name, description, parent_category_id, created_at)
VALUES
    ('150e8400-e29b-41d4-a716-446655440053', N'Contactors/Relays', N'Công tắc tơ/Rơ-le', '150e8400-e29b-41d4-a716-446655440035', GETDATE()),
    ('150e8400-e29b-41d4-a716-446655440054', N'High-Voltage Fuses', N'Cầu chì cao áp', '150e8400-e29b-41d4-a716-446655440035', GETDATE());

-- Cấp 4: Các mục con của Contactors/Relays
INSERT INTO part_categories (id, category_name, description, parent_category_id, created_at)
VALUES
    ('150e8400-e29b-41d4-a716-446655440055', 'HVJB.CON.400V.MAIN', N'Công tắc tơ chính.', '150e8400-e29b-41d4-a716-446655440053', GETDATE()),
    ('150e8400-e29b-41d4-a716-446655440056', 'HVJB.CON.400V.PRECHG', N'Công tắc tơ sạc trước.', '150e8400-e29b-41d4-a716-446655440053', GETDATE()),
    ('150e8400-e29b-41d4-a716-446655440057', 'HVJB.CON.400V.FASTCHG', N'Công tắc tơ sạc nhanh.', '150e8400-e29b-41d4-a716-446655440053', GETDATE());

-- Cấp 4: Các mục con của High-Voltage Fuses
INSERT INTO part_categories (id, category_name, description, parent_category_id, created_at)
VALUES
    ('150e8400-e29b-41d4-a716-446655440058', 'HVJB.FUSE.400A', N'Cầu chì 400A.', '150e8400-e29b-41d4-a716-446655440054', GETDATE()),
    ('150e8400-e29b-41d4-a716-446655440059', 'HVJB.FUSE.500A', N'Cầu chì 500A.', '150e8400-e29b-41d4-a716-446655440054', GETDATE()),
    ('150e8400-e29b-41d4-a716-44665544005A', 'HVJB.FUSE.SVC.KIT', N'Bộ cầu chì thay thế dịch vụ.', '150e8400-e29b-41d4-a716-446655440054', GETDATE());

-- Cấp 2: Các danh mục con của Transmission
INSERT INTO part_categories (id, category_name, description, parent_category_id, created_at)
VALUES
    ('150e8400-e29b-41d4-a716-44665544005B', N'Gear Reducer Assembly', N'Cụm hộp giảm tốc', '150e8400-e29b-41d4-a716-446655440014', GETDATE()),
    ('150e8400-e29b-41d4-a716-44665544005C', N'Differential', N'Bộ vi sai', '150e8400-e29b-41d4-a716-446655440014', GETDATE()),
    ('150e8400-e29b-41d4-a716-44665544005D', N'Driveshafts & Axles', N'Trục láp và trục truyền động', '150e8400-e29b-41d4-a716-446655440014', GETDATE());

-- Cấp 3: Các danh mục con của Gear Reducer Assembly
INSERT INTO part_categories (id, category_name, description, parent_category_id, created_at)
VALUES
    ('150e8400-e29b-41d4-a716-44665544005E', N'Front Single-Speed Gearbox', N'Hộp số giảm tốc một cấp trước', '150e8400-e29b-41d4-a716-44665544005B', GETDATE()),
    ('150e8400-e29b-41d4-a716-44665544005F', N'Rear Single-Speed Gearbox', N'Hộp số giảm tốc một cấp sau', '150e8400-e29b-41d4-a716-44665544005B', GETDATE());

-- Cấp 4: Các mục con của Front Single-Speed Gearbox
INSERT INTO part_categories (id, category_name, description, parent_category_id, created_at)
VALUES
    ('150e8400-e29b-41d4-a716-446655440060', 'GR.ASM.FR.R9-1', N'Hộp số trước, tỉ số truyền 9.1.', '150e8400-e29b-41d4-a716-44665544005E', GETDATE()),
    ('150e8400-e29b-41d4-a716-446655440061', 'GR.ASM.FR.R9-3.PLUS', N'Hộp số trước, tỉ số truyền 9.3 (bản Plus).', '150e8400-e29b-41d4-a716-44665544005E', GETDATE()),
    ('150e8400-e29b-41d4-a716-446655440062', 'GR.ASM.FR.R9-1.SVC', N'Phụ tùng thay thế, tỉ số truyền 9.1.', '150e8400-e29b-41d4-a716-44665544005E', GETDATE());

-- Cấp 4: Các mục con của Rear Single-Speed Gearbox
INSERT INTO part_categories (id, category_name, description, parent_category_id, created_at)
VALUES
    ('150e8400-e29b-41d4-a716-446655440063', 'GR.ASM.RR.R9-1', N'Hộp số sau, tỉ số truyền 9.1.', '150e8400-e29b-41d4-a716-44665544005F', GETDATE()),
    ('150e8400-e29b-41d4-a716-446655440064', 'GR.ASM.RR.R9-3.PLUS', N'Hộp số sau, tỉ số truyền 9.3 (bản Plus).', '150e8400-e29b-41d4-a716-44665544005F', GETDATE()),
    ('150e8400-e29b-41d4-a716-446655440065', 'GR.ASM.RR.R9-1.SVC', N'Phụ tùng thay thế, tỉ số truyền 9.1.', '150e8400-e29b-41d4-a716-44665544005F', GETDATE());

-- Cấp 3: Các danh mục con của Differential
INSERT INTO part_categories (id, category_name, description, parent_category_id, created_at)
VALUES
    ('150e8400-e29b-41d4-a716-446655440066', N'Front Differential Unit', N'Cụm vi sai trước', '150e8400-e29b-41d4-a716-44665544005C', GETDATE()),
    ('150e8400-e29b-41d4-a716-446655440067', N'Rear Differential Unit', N'Cụm vi sai sau', '150e8400-e29b-41d4-a716-44665544005C', GETDATE());

-- Cấp 4: Các mục con của Front Differential Unit
INSERT INTO part_categories (id, category_name, description, parent_category_id, created_at)
VALUES
    ('150e8400-e29b-41d4-a716-446655440068', 'DIFF.FR.OPEN.V1', N'Vi sai mở phía trước, phiên bản 1.', '150e8400-e29b-41d4-a716-446655440066', GETDATE()),
    ('150e8400-e29b-41d4-a716-446655440069', 'DIFF.FR.OPEN.V2', N'Vi sai mở phía trước, phiên bản 2 (cải tiến).', '150e8400-e29b-41d4-a716-446655440066', GETDATE()),
    ('150e8400-e29b-41d4-a716-44665544006A', 'DIFF.FR.OPEN.SVC.KIT', N'Bộ phụ tùng sửa chữa vi sai.', '150e8400-e29b-41d4-a716-446655440066', GETDATE());

-- Cấp 4: Các mục con của Rear Differential Unit
INSERT INTO part_categories (id, category_name, description, parent_category_id, created_at)
VALUES
    ('150e8400-e29b-41d4-a716-44665544006B', 'DIFF.RR.OPEN.V1', N'Vi sai mở phía sau, phiên bản 1.', '150e8400-e29b-41d4-a716-446655440067', GETDATE()),
    ('150e8400-e29b-41d4-a716-44665544006C', 'DIFF.RR.OPEN.V2', N'Vi sai mở phía sau, phiên bản 2 (cải tiến).', '150e8400-e29b-41d4-a716-446655440067', GETDATE()),
    ('150e8400-e29b-41d4-a716-44665544006D', 'DIFF.RR.OPEN.SVC.KIT', N'Bộ phụ tùng sửa chữa vi sai.', '150e8400-e29b-41d4-a716-446655440067', GETDATE());

-- Cấp 3: Các danh mục con của Driveshafts & Axles
INSERT INTO part_categories (id, category_name, description, parent_category_id, created_at)
VALUES
    ('150e8400-e29b-41d4-a716-44665544006E', N'Front Axle Shaft', N'Trục láp trước', '150e8400-e29b-41d4-a716-44665544005D', GETDATE()),
    ('150e8400-e29b-41d4-a716-44665544006F', N'Rear Axle Shaft', N'Trục láp sau', '150e8400-e29b-41d4-a716-44665544005D', GETDATE());

-- Cấp 4: Các mục con của Front Axle Shaft
INSERT INTO part_categories (id, category_name, description, parent_category_id, created_at)
VALUES
    ('150e8400-e29b-41d4-a716-446655440070', 'AXLE.FR.LH', N'Trục láp trước bên trái.', '150e8400-e29b-41d4-a716-44665544006E', GETDATE()),
    ('150e8400-e29b-41d4-a716-446655440071', 'AXLE.FR.RH', N'Trục láp trước bên phải.', '150e8400-e29b-41d4-a716-44665544006E', GETDATE()),
    ('150e8400-e29b-41d4-a716-446655440072', 'AXLE.FR.ASM.SVC', N'Bộ trục láp trước thay thế.', '150e8400-e29b-41d4-a716-44665544006E', GETDATE());

-- Cấp 4: Các mục con của Rear Axle Shaft
INSERT INTO part_categories (id, category_name, description, parent_category_id, created_at)
VALUES
    ('150e8400-e29b-41d4-a716-446655440073', 'AXLE.RR.LH', N'Trục láp sau bên trái.', '150e8400-e29b-41d4-a716-44665544006F', GETDATE()),
    ('150e8400-e29b-41d4-a716-446655440074', 'AXLE.RR.RH', N'Trục láp sau bên phải.', '150e8400-e29b-41d4-a716-44665544006F', GETDATE()),
    ('150e8400-e29b-41d4-a716-446655440075', 'AXLE.RR.ASM.SVC', N'Bộ trục láp sau thay thế.', '150e8400-e29b-41d4-a716-44665544006F', GETDATE());

-- Cấp 2: Các danh mục con của Cooling System
INSERT INTO part_categories (id, category_name, description, parent_category_id, created_at)
VALUES
    ('150e8400-e29b-41d4-a716-446655440076', N'Battery Thermal Management', N'Quản lý nhiệt cho pin', '150e8400-e29b-41d4-a716-446655440015', GETDATE()),
    ('150e8400-e29b-41d4-a716-446655440077', N'Drivetrain Thermal Management', N'Quản lý nhiệt cho hệ truyền động', '150e8400-e29b-41d4-a716-446655440015', GETDATE());

-- Cấp 3: Các danh mục con của Battery Thermal Management
INSERT INTO part_categories (id, category_name, description, parent_category_id, created_at)
VALUES
    ('150e8400-e29b-41d4-a716-446655440078', N'Electric Coolant Pump', N'Bơm dung dịch làm mát bằng điện', '150e8400-e29b-41d4-a716-446655440076', GETDATE()),
    ('150e8400-e29b-41d4-a716-446655440079', N'Battery Chiller/Heat Exchanger', N'Bộ làm lạnh/trao đổi nhiệt cho pin', '150e8400-e29b-41d4-a716-446655440076', GETDATE()),
    ('150e8400-e29b-41d4-a716-44665544007A', N'Cooling Plates', N'Tấm làm mát pin', '150e8400-e29b-41d4-a716-446655440076', GETDATE());

-- Cấp 4: Các mục con của Electric Coolant Pump
INSERT INTO part_categories (id, category_name, description, parent_category_id, created_at)
VALUES
    ('150e8400-e29b-41d4-a716-44665544007B', 'PUMP.BAT.12V.BOSCH', N'Bơm làm mát pin, 12V, nhà cung cấp Bosch.', '150e8400-e29b-41d4-a716-446655440078', GETDATE()),
    ('150e8400-e29b-41d4-a716-44665544007C', 'PUMP.BAT.12V.VALEO', N'Bơm làm mát pin, 12V, nhà cung cấp Valeo.', '150e8400-e29b-41d4-a716-446655440078', GETDATE()),
    ('150e8400-e29b-41d4-a716-44665544007D', 'PUMP.BAT.12V.HP.GATES', N'Bơm làm mát pin, 12V, hiệu suất cao, nhà cung cấp Gates.', '150e8400-e29b-41d4-a716-446655440078', GETDATE());

-- Cấp 4: Các mục con của Battery Chiller/Heat Exchanger
INSERT INTO part_categories (id, category_name, description, parent_category_id, created_at)
VALUES
    ('150e8400-e29b-41d4-a716-44665544007E', 'CHILLER.BAT.V1', N'Bộ làm lạnh pin, phiên bản 1.', '150e8400-e29b-41d4-a716-446655440079', GETDATE()),
    ('150e8400-e29b-41d4-a716-44665544007F', 'CHILLER.BAT.V2.ENH', N'Bộ làm lạnh pin, phiên bản 2, hiệu suất tăng cường.', '150e8400-e29b-41d4-a716-446655440079', GETDATE()),
    ('150e8400-e29b-41d4-a716-446655440080', 'CHILLER.BAT.V1.SVC', N'Phụ tùng thay thế cho phiên bản 1.', '150e8400-e29b-41d4-a716-446655440079', GETDATE());

-- Cấp 4: Các mục con của Cooling Plates
INSERT INTO part_categories (id, category_name, description, parent_category_id, created_at)
VALUES
    ('150e8400-e29b-41d4-a716-446655440081', 'PLATE.COOL.82KWH', N'Bộ tấm làm mát cho pin 82kWh.', '150e8400-e29b-41d4-a716-44665544007A', GETDATE()),
    ('150e8400-e29b-41d4-a716-446655440082', 'PLATE.COOL.87KWH', N'Bộ tấm làm mát cho pin 87.7kWh.', '150e8400-e29b-41d4-a716-44665544007A', GETDATE()),
    ('150e8400-e29b-41d4-a716-446655440083', 'PLATE.COOL.SVC.KIT', N'Bộ phụ tùng tấm làm mát.', '150e8400-e29b-41d4-a716-44665544007A', GETDATE());

-- Cấp 3: Các danh mục con của Drivetrain Thermal Management
INSERT INTO part_categories (id, category_name, description, parent_category_id, created_at)
VALUES
    ('150e8400-e29b-41d4-a716-446655440084', N'Drivetrain Radiator', N'Két làm mát hệ truyền động', '150e8400-e29b-41d4-a716-446655440077', GETDATE()),
    ('150e8400-e29b-41d4-a716-446655440085', N'Radiator Fan Assembly', N'Cụm quạt két làm mát', '150e8400-e29b-41d4-a716-446655440077', GETDATE()),
    ('150e8400-e29b-41d4-a716-446655440086', N'Three-Way Coolant Valve', N'Van dung dịch làm mát 3 ngả', '150e8400-e29b-41d4-a716-446655440077', GETDATE());

-- Cấp 4: Các mục con của Drivetrain Radiator
INSERT INTO part_categories (id, category_name, description, parent_category_id, created_at)
VALUES
    ('150e8400-e29b-41d4-a716-446655440087', 'RAD.DRV.STD', N'Két làm mát tiêu chuẩn.', '150e8400-e29b-41d4-a716-446655440084', GETDATE()),
    ('150e8400-e29b-41d4-a716-446655440088', 'RAD.DRV.TROPICAL', N'Két làm mát cho thị trường nhiệt đới.', '150e8400-e29b-41d4-a716-446655440084', GETDATE()),
    ('150e8400-e29b-41d4-a716-446655440089', 'RAD.DRV.SVC', N'Két làm mát thay thế.', '150e8400-e29b-41d4-a716-446655440084', GETDATE());

-- Cấp 4: Các mục con của Radiator Fan Assembly
INSERT INTO part_categories (id, category_name, description, parent_category_id, created_at)
VALUES
    ('150e8400-e29b-41d4-a716-44665544008A', 'FAN.RAD.500W.V1', N'Cụm quạt 500W, phiên bản 1.', '150e8400-e29b-41d4-a716-446655440085', GETDATE()),
    ('150e8400-e29b-41d4-a716-44665544008B', 'FAN.RAD.600W.HP', N'Cụm quạt 600W, hiệu suất cao.', '150e8400-e29b-41d4-a716-446655440085', GETDATE()),
    ('150e8400-e29b-41d4-a716-44665544008C', 'FAN.RAD.500W.SVC', N'Phụ tùng thay thế 500W.', '150e8400-e29b-41d4-a716-446655440085', GETDATE());

-- Cấp 4: Các mục con của Three-Way Coolant Valve
INSERT INTO part_categories (id, category_name, description, parent_category_id, created_at)
VALUES
    ('150e8400-e29b-41d4-a716-44665544008D', 'VALVE.3WAY.V1', N'Van 3 ngả, phiên bản 1.', '150e8400-e29b-41d4-a716-446655440086', GETDATE()),
    ('150e8400-e29b-41d4-a716-44665544008E', 'VALVE.3WAY.V2.ENH', N'Van 3 ngả, phiên bản 2, cải tiến.', '150e8400-e29b-41d4-a716-446655440086', GETDATE()),
    ('150e8400-e29b-41d4-a716-44665544008F', 'VALVE.3WAY.V1.SVC', N'Phụ tùng thay thế cho phiên bản 1.', '150e8400-e29b-41d4-a716-446655440086', GETDATE());

--INSERT INTO parts (id, serial_number, part_name, unit_price, category_id, office_location_id, created_at)
--VALUES
-- Parts for Front Motor Unit
-- Category: VF8.EM.FM.150KW.ECO
INSERT INTO parts (id, serial_number, part_name, unit_price, category_id, office_location_id, created_at)
VALUES
    (NEWID(), 'FM150ECO000000001', N'Động cơ trước 150kW cho phiên bản Eco - Kho HN', 75000000, '150e8400-e29b-41d4-a716-44665544001A', '550e8400-e29b-41d4-a716-446655440001', GETDATE()),
    (NEWID(), 'FM150ECO000000002', N'Động cơ trước 150kW cho phiên bản Eco - Kho HP', 75000000, '150e8400-e29b-41d4-a716-44665544001A', '550e8400-e29b-41d4-a716-446655440002', GETDATE()),
    (NEWID(), 'FM150ECO000000003', N'Động cơ trước 150kW cho phiên bản Eco - Kho DN', 75000000, '150e8400-e29b-41d4-a716-44665544001A', '550e8400-e29b-41d4-a716-446655440004', GETDATE());

-- Category: VF8.EM.FM.150KW.PLUS
INSERT INTO parts (id, serial_number, part_name, unit_price, category_id, office_location_id, created_at)
VALUES
    (NEWID(), 'FM150PLUS00000001', N'Động cơ trước 150kW cho phiên bản Plus - Kho HN', 85000000, '150e8400-e29b-41d4-a716-44665544001B', '550e8400-e29b-41d4-a716-446655440001', GETDATE()),
    (NEWID(), 'FM150PLUS00000002', N'Động cơ trước 150kW cho phiên bản Plus - Kho HP', 85000000, '150e8400-e29b-41d4-a716-44665544001B', '550e8400-e29b-41d4-a716-446655440002', GETDATE()),
    (NEWID(), 'FM150PLUS00000003', N'Động cơ trước 150kW cho phiên bản Plus - Kho DN', 85000000, '150e8400-e29b-41d4-a716-44665544001B', '550e8400-e29b-41d4-a716-446655440004', GETDATE());

-- Category: VF8.EM.FM.150KW.SVC.NEW
INSERT INTO parts (id, serial_number, part_name, unit_price, category_id, office_location_id, created_at)
VALUES
    (NEWID(), 'FM150SVCNEW000001', N'Phụ tùng thay thế mới (Động cơ trước) - Kho HN', 90000000, '150e8400-e29b-41d4-a716-44665544001C', '550e8400-e29b-41d4-a716-446655440001', GETDATE()),
    (NEWID(), 'FM150SVCNEW000002', N'Phụ tùng thay thế mới (Động cơ trước) - Kho HP', 90000000, '150e8400-e29b-41d4-a716-44665544001C', '550e8400-e29b-41d4-a716-446655440002', GETDATE()),
    (NEWID(), 'FM150SVCNEW000003', N'Phụ tùng thay thế mới (Động cơ trước) - Kho DN', 90000000, '150e8400-e29b-41d4-a716-44665544001C', '550e8400-e29b-41d4-a716-446655440004', GETDATE());

-- Category: VF8.EM.FM.150KW.SVC.RMD
INSERT INTO parts (id, serial_number, part_name, unit_price, category_id, office_location_id, created_at)
VALUES
    (NEWID(), 'FM150SVCRMD00001', N'Phụ tùng thay thế tái sản xuất (Động cơ trước) - Kho HN', 55000000, '150e8400-e29b-41d4-a716-44665544001D', '550e8400-e29b-41d4-a716-446655440001', GETDATE()),
    (NEWID(), 'FM150SVCRMD00002', N'Phụ tùng thay thế tái sản xuất (Động cơ trước) - Kho HP', 55000000, '150e8400-e29b-41d4-a716-44665544001D', '550e8400-e29b-41d4-a716-446655440002', GETDATE()),
    (NEWID(), 'FM150SVCRMD00003', N'Phụ tùng thay thế tái sản xuất (Động cơ trước) - Kho DN', 55000000, '150e8400-e29b-41d4-a716-44665544001D', '550e8400-e29b-41d4-a716-446655440004', GETDATE());

-- Parts for Rear Motor Unit
-- Category: VF8.EM.RM.150KW.ECO
INSERT INTO parts (id, serial_number, part_name, unit_price, category_id, office_location_id, created_at)
VALUES
    (NEWID(), 'RM150ECO000000001', N'Động cơ sau 150kW cho phiên bản Eco - Kho HN', 75000000, '150e8400-e29b-41d4-a716-44665544001E', '550e8400-e29b-41d4-a716-446655440001', GETDATE()),
    (NEWID(), 'RM150ECO000000002', N'Động cơ sau 150kW cho phiên bản Eco - Kho HP', 75000000, '150e8400-e29b-41d4-a716-44665544001E', '550e8400-e29b-41d4-a716-446655440002', GETDATE()),
    (NEWID(), 'RM150ECO000000003', N'Động cơ sau 150kW cho phiên bản Eco - Kho DN', 75000000, '150e8400-e29b-41d4-a716-44665544001E', '550e8400-e29b-41d4-a716-446655440004', GETDATE());

-- Category: VF8.EM.RM.150KW.PLUS
INSERT INTO parts (id, serial_number, part_name, unit_price, category_id, office_location_id, created_at)
VALUES
    (NEWID(), 'RM150PLUS00000001', N'Động cơ sau 150kW cho phiên bản Plus - Kho HN', 85000000, '150e8400-e29b-41d4-a716-44665544001F', '550e8400-e29b-41d4-a716-446655440001', GETDATE()),
    (NEWID(), 'RM150PLUS00000002', N'Động cơ sau 150kW cho phiên bản Plus - Kho HP', 85000000, '150e8400-e29b-41d4-a716-44665544001F', '550e8400-e29b-41d4-a716-446655440002', GETDATE()),
    (NEWID(), 'RM150PLUS00000003', N'Động cơ sau 150kW cho phiên bản Plus - Kho DN', 85000000, '150e8400-e29b-41d4-a716-44665544001F', '550e8400-e29b-41d4-a716-446655440004', GETDATE());

-- Category: VF8.EM.RM.150KW.SVC.NEW
INSERT INTO parts (id, serial_number, part_name, unit_price, category_id, office_location_id, created_at)
VALUES
    (NEWID(), 'RM150SVCNEW000001', N'Phụ tùng thay thế mới (Động cơ sau) - Kho HN', 90000000, '150e8400-e29b-41d4-a716-446655440020', '550e8400-e29b-41d4-a716-446655440001', GETDATE()),
    (NEWID(), 'RM150SVCNEW000002', N'Phụ tùng thay thế mới (Động cơ sau) - Kho HP', 90000000, '150e8400-e29b-41d4-a716-446655440020', '550e8400-e29b-41d4-a716-446655440002', GETDATE()),
    (NEWID(), 'RM150SVCNEW000003', N'Phụ tùng thay thế mới (Động cơ sau) - Kho DN', 90000000, '150e8400-e29b-41d4-a716-446655440020', '550e8400-e29b-41d4-a716-446655440004', GETDATE());

-- Category: VF8.EM.RM.150KW.SVC.RMD
INSERT INTO parts (id, serial_number, part_name, unit_price, category_id, office_location_id, created_at)
VALUES
    (NEWID(), 'RM150SVCRMD00001', N'Phụ tùng thay thế tái sản xuất (Động cơ sau) - Kho HN', 55000000, '150e8400-e29b-41d4-a716-446655440021', '550e8400-e29b-41d4-a716-446655440001', GETDATE()),
    (NEWID(), 'RM150SVCRMD00002', N'Phụ tùng thay thế tái sản xuất (Động cơ sau) - Kho HP', 55000000, '150e8400-e29b-41d4-a716-446655440021', '550e8400-e29b-41d4-a716-446655440002', GETDATE()),
    (NEWID(), 'RM150SVCRMD00003', N'Phụ tùng thay thế tái sản xuất (Động cơ sau) - Kho DN', 55000000, '150e8400-e29b-41d4-a716-446655440021', '550e8400-e29b-41d4-a716-446655440004', GETDATE());

-- Parts for Inverter/Power Electronics Module
-- Category: PEM.400V.V1.BOSCH
INSERT INTO parts (id, serial_number, part_name, unit_price, category_id, office_location_id, created_at)
VALUES
    (NEWID(), 'PEM400V1BOSCH0001', N'Biến tần 400V, V1, Bosch - Kho HN', 35000000, '150e8400-e29b-41d4-a716-446655440026', '550e8400-e29b-41d4-a716-446655440001', GETDATE()),
    (NEWID(), 'PEM400V1BOSCH0002', N'Biến tần 400V, V1, Bosch - Kho HP', 35000000, '150e8400-e29b-41d4-a716-446655440026', '550e8400-e29b-41d4-a716-446655440002', GETDATE()),
    (NEWID(), 'PEM400V1BOSCH0003', N'Biến tần 400V, V1, Bosch - Kho DN', 35000000, '150e8400-e29b-41d4-a716-446655440026', '550e8400-e29b-41d4-a716-446655440004', GETDATE());

-- Category: PEM.400V.V2.VISTEON
INSERT INTO parts (id, serial_number, part_name, unit_price, category_id, office_location_id, created_at)
VALUES
    (NEWID(), 'PEM400V2VISTEON01', N'Biến tần 400V, V2, Visteon - Kho HN', 38000000, '150e8400-e29b-41d4-a716-446655440027', '550e8400-e29b-41d4-a716-446655440001', GETDATE()),
    (NEWID(), 'PEM400V2VISTEON02', N'Biến tần 400V, V2, Visteon - Kho HP', 38000000, '150e8400-e29b-41d4-a716-446655440027', '550e8400-e29b-41d4-a716-446655440002', GETDATE()),
    (NEWID(), 'PEM400V2VISTEON03', N'Biến tần 400V, V2, Visteon - Kho DN', 38000000, '150e8400-e29b-41d4-a716-446655440027', '550e8400-e29b-41d4-a716-446655440004', GETDATE());

-- Category: PEM.400V.V2.1.SVC
INSERT INTO parts (id, serial_number, part_name, unit_price, category_id, office_location_id, created_at)
VALUES
    (NEWID(), 'PEM400V21SVC00001', N'Biến tần 400V, V2.1, phụ tùng thay thế - Kho HN', 40000000, '150e8400-e29b-41d4-a716-446655440028', '550e8400-e29b-41d4-a716-446655440001', GETDATE()),
    (NEWID(), 'PEM400V21SVC00002', N'Biến tần 400V, V2.1, phụ tùng thay thế - Kho HP', 40000000, '150e8400-e29b-41d4-a716-446655440028', '550e8400-e29b-41d4-a716-446655440002', GETDATE()),
    (NEWID(), 'PEM400V21SVC00003', N'Biến tần 400V, V2.1, phụ tùng thay thế - Kho DN', 40000000, '150e8400-e29b-41d4-a716-446655440028', '550e8400-e29b-41d4-a716-446655440004', GETDATE());

-- Parts for Rotor Assembly
-- Category: ROT.ASM.150KW.G1
INSERT INTO parts (id, serial_number, part_name, unit_price, category_id, office_location_id, created_at)
VALUES
    (NEWID(), 'ROTASM150G1000001', N'Cụm Rô-to 150kW G1 - Kho HN', 22000000, '150e8400-e29b-41d4-a716-446655440029', '550e8400-e29b-41d4-a716-446655440001', GETDATE()),
    (NEWID(), 'ROTASM150G1000002', N'Cụm Rô-to 150kW G1 - Kho HP', 22000000, '150e8400-e29b-41d4-a716-446655440029', '550e8400-e29b-41d4-a716-446655440002', GETDATE()),
    (NEWID(), 'ROTASM150G1000003', N'Cụm Rô-to 150kW G1 - Kho DN', 22000000, '150e8400-e29b-41d4-a716-446655440029', '550e8400-e29b-41d4-a716-446655440004', GETDATE());

-- Category: ROT.ASM.150KW.G2
INSERT INTO parts (id, serial_number, part_name, unit_price, category_id, office_location_id, created_at)
VALUES
    (NEWID(), 'ROTASM150G2000001', N'Cụm Rô-to 150kW G2 (cải tiến) - Kho HN', 25000000, '150e8400-e29b-41d4-a716-44665544002A', '550e8400-e29b-41d4-a716-446655440001', GETDATE()),
    (NEWID(), 'ROTASM150G2000002', N'Cụm Rô-to 150kW G2 (cải tiến) - Kho HP', 25000000, '150e8400-e29b-41d4-a716-44665544002A', '550e8400-e29b-41d4-a716-446655440002', GETDATE()),
    (NEWID(), 'ROTASM150G2000003', N'Cụm Rô-to 150kW G2 (cải tiến) - Kho DN', 25000000, '150e8400-e29b-41d4-a716-44665544002A', '550e8400-e29b-41d4-a716-446655440004', GETDATE());

-- Category: ROT.ASM.150KW.G1.SVC
INSERT INTO parts (id, serial_number, part_name, unit_price, category_id, office_location_id, created_at)
VALUES
    (NEWID(), 'ROTASM150G1SVC001', N'Phụ tùng thay thế Rô-to G1 - Kho HN', 23000000, '150e8400-e29b-41d4-a716-44665544002B', '550e8400-e29b-41d4-a716-446655440001', GETDATE()),
    (NEWID(), 'ROTASM150G1SVC002', N'Phụ tùng thay thế Rô-to G1 - Kho HP', 23000000, '150e8400-e29b-41d4-a716-44665544002B', '550e8400-e29b-41d4-a716-446655440002', GETDATE()),
    (NEWID(), 'ROTASM150G1SVC003', N'Phụ tùng thay thế Rô-to G1 - Kho DN', 23000000, '150e8400-e29b-41d4-a716-44665544002B', '550e8400-e29b-41d4-a716-446655440004', GETDATE());

-- Parts for Stator Assembly
-- Category: STA.ASM.150KW.G1
INSERT INTO parts (id, serial_number, part_name, unit_price, category_id, office_location_id, created_at)
VALUES
    (NEWID(), 'STAASM150G1000001', N'Cụm Stato 150kW G1 - Kho HN', 28000000, '150e8400-e29b-41d4-a716-44665544002C', '550e8400-e29b-41d4-a716-446655440001', GETDATE()),
    (NEWID(), 'STAASM150G1000002', N'Cụm Stato 150kW G1 - Kho HP', 28000000, '150e8400-e29b-41d4-a716-44665544002C', '550e8400-e29b-41d4-a716-446655440002', GETDATE()),
    (NEWID(), 'STAASM150G1000003', N'Cụm Stato 150kW G1 - Kho DN', 28000000, '150e8400-e29b-41d4-a716-44665544002C', '550e8400-e29b-41d4-a716-446655440004', GETDATE());

-- Category: STA.ASM.150KW.G2
INSERT INTO parts (id, serial_number, part_name, unit_price, category_id, office_location_id, created_at)
VALUES
    (NEWID(), 'STAASM150G2000001', N'Cụm Stato 150kW G2 (cải tiến) - Kho HN', 31000000, '150e8400-e29b-41d4-a716-44665544002D', '550e8400-e29b-41d4-a716-446655440001', GETDATE()),
    (NEWID(), 'STAASM150G2000002', N'Cụm Stato 150kW G2 (cải tiến) - Kho HP', 31000000, '150e8400-e29b-41d4-a716-44665544002D', '550e8400-e29b-41d4-a716-446655440002', GETDATE()),
    (NEWID(), 'STAASM150G2000003', N'Cụm Stato 150kW G2 (cải tiến) - Kho DN', 31000000, '150e8400-e29b-41d4-a716-44665544002D', '550e8400-e29b-41d4-a716-446655440004', GETDATE());

-- Category: STA.ASM.150KW.G1.SVC
INSERT INTO parts (id, serial_number, part_name, unit_price, category_id, office_location_id, created_at)
VALUES
    (NEWID(), 'STAASM150G1SVC001', N'Phụ tùng thay thế Stato G1 - Kho HN', 29000000, '150e8400-e29b-41d4-a716-44665544002E', '550e8400-e29b-41d4-a716-446655440001', GETDATE()),
    (NEWID(), 'STAASM150G1SVC002', N'Phụ tùng thay thế Stato G1 - Kho HP', 29000000, '150e8400-e29b-41d4-a716-44665544002E', '550e8400-e29b-41d4-a716-446655440002', GETDATE()),
    (NEWID(), 'STAASM150G1SVC003', N'Phụ tùng thay thế Stato G1 - Kho DN', 29000000, '150e8400-e29b-41d4-a716-44665544002E', '550e8400-e29b-41d4-a716-446655440004', GETDATE());

-- Parts for Motor Position Sensor
-- Category: SENS.POS.RESOLVER.V1
INSERT INTO parts (id, serial_number, part_name, unit_price, category_id, office_location_id, created_at)
VALUES
    (NEWID(), 'SENSPOSRESV100001', N'Cảm biến vị trí Resolver V1 - Kho HN', 1200000, '150e8400-e29b-41d4-a716-44665544002F', '550e8400-e29b-41d4-a716-446655440001', GETDATE()),
    (NEWID(), 'SENSPOSRESV100002', N'Cảm biến vị trí Resolver V1 - Kho HP', 1200000, '150e8400-e29b-41d4-a716-44665544002F', '550e8400-e29b-41d4-a716-446655440002', GETDATE()),
    (NEWID(), 'SENSPOSRESV100003', N'Cảm biến vị trí Resolver V1 - Kho DN', 1200000, '150e8400-e29b-41d4-a716-44665544002F', '550e8400-e29b-41d4-a716-446655440004', GETDATE());

-- Category: SENS.POS.RESOLVER.V2.HELLA
INSERT INTO parts (id, serial_number, part_name, unit_price, category_id, office_location_id, created_at)
VALUES
    (NEWID(), 'SENSPOSRESV2H0001', N'Cảm biến vị trí Resolver V2, Hella - Kho HN', 1500000, '150e8400-e29b-41d4-a716-446655440030', '550e8400-e29b-41d4-a716-446655440001', GETDATE()),
    (NEWID(), 'SENSPOSRESV2H0002', N'Cảm biến vị trí Resolver V2, Hella - Kho HP', 1500000, '150e8400-e29b-41d4-a716-446655440030', '550e8400-e29b-41d4-a716-446655440002', GETDATE()),
    (NEWID(), 'SENSPOSRESV2H0003', N'Cảm biến vị trí Resolver V2, Hella - Kho DN', 1500000, '150e8400-e29b-41d4-a716-446655440030', '550e8400-e29b-41d4-a716-446655440004', GETDATE());

-- Category: SENS.POS.HALL.V3.BOSCH
INSERT INTO parts (id, serial_number, part_name, unit_price, category_id, office_location_id, created_at)
VALUES
    (NEWID(), 'SENSPOSHALLV3B001', N'Cảm biến vị trí Hall V3, Bosch - Kho HN', 1350000, '150e8400-e29b-41d4-a716-446655440031', '550e8400-e29b-41d4-a716-446655440001', GETDATE()),
    (NEWID(), 'SENSPOSHALLV3B002', N'Cảm biến vị trí Hall V3, Bosch - Kho HP', 1350000, '150e8400-e29b-41d4-a716-446655440031', '550e8400-e29b-41d4-a716-446655440002', GETDATE()),
    (NEWID(), 'SENSPOSHALLV3B003', N'Cảm biến vị trí Hall V3, Bosch - Kho DN', 1350000, '150e8400-e29b-41d4-a716-446655440031', '550e8400-e29b-41d4-a716-446655440004', GETDATE());

-- Parts for Battery Cell
-- Category: CELL.NMC.82KWH.CATL
INSERT INTO parts (id, serial_number, part_name, unit_price, category_id, office_location_id, created_at)
VALUES
    (NEWID(), 'CELLNMC82CATL0001', N'Cell pin NMC 82kWh, CATL - Kho HN', 2500000, '150e8400-e29b-41d4-a716-446655440038', '550e8400-e29b-41d4-a716-446655440001', GETDATE()),
    (NEWID(), 'CELLNMC82CATL0002', N'Cell pin NMC 82kWh, CATL - Kho HP', 2500000, '150e8400-e29b-41d4-a716-446655440038', '550e8400-e29b-41d4-a716-446655440002', GETDATE()),
    (NEWID(), 'CELLNMC82CATL0003', N'Cell pin NMC 82kWh, CATL - Kho DN', 2500000, '150e8400-e29b-41d4-a716-446655440038', '550e8400-e29b-41d4-a716-446655440004', GETDATE());

-- Category: CELL.NMC.87KWH.SDI
INSERT INTO parts (id, serial_number, part_name, unit_price, category_id, office_location_id, created_at)
VALUES
    (NEWID(), 'CELLNMC87SDI00001', N'Cell pin NMC 87.7kWh, Samsung SDI - Kho HN', 2800000, '150e8400-e29b-41d4-a716-446655440039', '550e8400-e29b-41d4-a716-446655440001', GETDATE()),
    (NEWID(), 'CELLNMC87SDI00002', N'Cell pin NMC 87.7kWh, Samsung SDI - Kho HP', 2800000, '150e8400-e29b-41d4-a716-446655440039', '550e8400-e29b-41d4-a716-446655440002', GETDATE()),
    (NEWID(), 'CELLNMC87SDI00003', N'Cell pin NMC 87.7kWh, Samsung SDI - Kho DN', 2800000, '150e8400-e29b-41d4-a716-446655440039', '550e8400-e29b-41d4-a716-446655440004', GETDATE());

-- Category: CELL.NMC.87KWH.CATL.G2
INSERT INTO parts (id, serial_number, part_name, unit_price, category_id, office_location_id, created_at)
VALUES
    (NEWID(), 'CELLNMC87CATLG201', N'Cell pin NMC G2 87.7kWh, CATL - Kho HN', 2900000, '150e8400-e29b-41d4-a716-44665544003A', '550e8400-e29b-41d4-a716-446655440001', GETDATE()),
    (NEWID(), 'CELLNMC87CATLG202', N'Cell pin NMC G2 87.7kWh, CATL - Kho HP', 2900000, '150e8400-e29b-41d4-a716-44665544003A', '550e8400-e29b-41d4-a716-446655440002', GETDATE()),
    (NEWID(), 'CELLNMC87CATLG203', N'Cell pin NMC G2 87.7kWh, CATL - Kho DN', 2900000, '150e8400-e29b-41d4-a716-44665544003A', '550e8400-e29b-41d4-a716-446655440004', GETDATE());

-- Category: CELL.LFP.75KWH.GOTION
INSERT INTO parts (id, serial_number, part_name, unit_price, category_id, office_location_id, created_at)
VALUES
    (NEWID(), 'CELLLFP75GOTION01', N'Cell pin LFP 75kWh, Gotion - Kho HN', 2200000, '150e8400-e29b-41d4-a716-44665544003B', '550e8400-e29b-41d4-a716-446655440001', GETDATE()),
    (NEWID(), 'CELLLFP75GOTION02', N'Cell pin LFP 75kWh, Gotion - Kho HP', 2200000, '150e8400-e29b-41d4-a716-44665544003B', '550e8400-e29b-41d4-a716-446655440002', GETDATE()),
    (NEWID(), 'CELLLFP75GOTION03', N'Cell pin LFP 75kWh, Gotion - Kho DN', 2200000, '150e8400-e29b-41d4-a716-44665544003B', '550e8400-e29b-41d4-a716-446655440004', GETDATE());

-- Parts for Module Interconnect Board
-- Category: MIB.82KWH.V1
INSERT INTO parts (id, serial_number, part_name, unit_price, category_id, office_location_id, created_at)
VALUES
    (NEWID(), 'MIB82KWHV10000001', N'Bảng mạch MIB 82kWh, V1 - Kho HN', 4500000, '150e8400-e29b-41d4-a716-44665544003C', '550e8400-e29b-41d4-a716-446655440001', GETDATE()),
    (NEWID(), 'MIB82KWHV10000002', N'Bảng mạch MIB 82kWh, V1 - Kho HP', 4500000, '150e8400-e29b-41d4-a716-44665544003C', '550e8400-e29b-41d4-a716-446655440002', GETDATE()),
    (NEWID(), 'MIB82KWHV10000003', N'Bảng mạch MIB 82kWh, V1 - Kho DN', 4500000, '150e8400-e29b-41d4-a716-44665544003C', '550e8400-e29b-41d4-a716-446655440004', GETDATE());

-- Category: MIB.87KWH.V1
INSERT INTO parts (id, serial_number, part_name, unit_price, category_id, office_location_id, created_at)
VALUES
    (NEWID(), 'MIB87KWHV10000001', N'Bảng mạch MIB 87.7kWh, V1 - Kho HN', 4800000, '150e8400-e29b-41d4-a716-44665544003D', '550e8400-e29b-41d4-a716-446655440001', GETDATE()),
    (NEWID(), 'MIB87KWHV10000002', N'Bảng mạch MIB 87.7kWh, V1 - Kho HP', 4800000, '150e8400-e29b-41d4-a716-44665544003D', '550e8400-e29b-41d4-a716-446655440002', GETDATE()),
    (NEWID(), 'MIB87KWHV10000003', N'Bảng mạch MIB 87.7kWh, V1 - Kho DN', 4800000, '150e8400-e29b-41d4-a716-44665544003D', '550e8400-e29b-41d4-a716-446655440004', GETDATE());

-- Category: MIB.87KWH.V2.ENH
INSERT INTO parts (id, serial_number, part_name, unit_price, category_id, office_location_id, created_at)
VALUES
    (NEWID(), 'MIB87KWHV2ENH0001', N'Bảng mạch MIB 87.7kWh, V2, cải tiến - Kho HN', 5200000, '150e8400-e29b-41d4-a716-44665544003E', '550e8400-e29b-41d4-a716-446655440001', GETDATE()),
    (NEWID(), 'MIB87KWHV2ENH0002', N'Bảng mạch MIB 87.7kWh, V2, cải tiến - Kho HP', 5200000, '150e8400-e29b-41d4-a716-44665544003E', '550e8400-e29b-41d4-a716-446655440002', GETDATE()),
    (NEWID(), 'MIB87KWHV2ENH0003', N'Bảng mạch MIB 87.7kWh, V2, cải tiến - Kho DN', 5200000, '150e8400-e29b-41d4-a716-44665544003E', '550e8400-e29b-41d4-a716-446655440004', GETDATE());

-- Parts for Master Control Unit
-- Category: BMS.MCU.HW1.FW25
INSERT INTO parts (id, serial_number, part_name, unit_price, category_id, office_location_id, created_at)
VALUES
    (NEWID(), 'BMSMCUHW1FW250001', N'BMS MCU HW1, FW2.5 - Kho HN', 6500000, '150e8400-e29b-41d4-a716-446655440041', '550e8400-e29b-41d4-a716-446655440001', GETDATE()),
    (NEWID(), 'BMSMCUHW1FW250002', N'BMS MCU HW1, FW2.5 - Kho HP', 6500000, '150e8400-e29b-41d4-a716-446655440041', '550e8400-e29b-41d4-a716-446655440002', GETDATE()),
    (NEWID(), 'BMSMCUHW1FW250003', N'BMS MCU HW1, FW2.5 - Kho DN', 6500000, '150e8400-e29b-41d4-a716-446655440041', '550e8400-e29b-41d4-a716-446655440004', GETDATE());

-- Category: BMS.MCU.HW2.FW31
INSERT INTO parts (id, serial_number, part_name, unit_price, category_id, office_location_id, created_at)
VALUES
    (NEWID(), 'BMSMCUHW2FW310001', N'BMS MCU HW2, FW3.1 - Kho HN', 7200000, '150e8400-e29b-41d4-a716-446655440042', '550e8400-e29b-41d4-a716-446655440001', GETDATE()),
    (NEWID(), 'BMSMCUHW2FW310002', N'BMS MCU HW2, FW3.1 - Kho HP', 7200000, '150e8400-e29b-41d4-a716-446655440042', '550e8400-e29b-41d4-a716-446655440002', GETDATE()),
    (NEWID(), 'BMSMCUHW2FW310003', N'BMS MCU HW2, FW3.1 - Kho DN', 7200000, '150e8400-e29b-41d4-a716-446655440042', '550e8400-e29b-41d4-a716-446655440004', GETDATE());

-- Category: BMS.MCU.HW2.FW33.OTA
INSERT INTO parts (id, serial_number, part_name, unit_price, category_id, office_location_id, created_at)
VALUES
    (NEWID(), 'BMSMCUHW2FW33OTA1', N'BMS MCU HW2, FW3.3 OTA - Kho HN', 7800000, '150e8400-e29b-41d4-a716-446655440043', '550e8400-e29b-41d4-a716-446655440001', GETDATE()),
    (NEWID(), 'BMSMCUHW2FW33OTA2', N'BMS MCU HW2, FW3.3 OTA - Kho HP', 7800000, '150e8400-e29b-41d4-a716-446655440043', '550e8400-e29b-41d4-a716-446655440002', GETDATE()),
    (NEWID(), 'BMSMCUHW2FW33OTA3', N'BMS MCU HW2, FW3.3 OTA - Kho DN', 7800000, '150e8400-e29b-41d4-a716-446655440043', '550e8400-e29b-41d4-a716-446655440004', GETDATE());

-- Parts for Slave/Cell Monitoring Unit
-- Category: BMS.CMU.HW1.82KWH
INSERT INTO parts (id, serial_number, part_name, unit_price, category_id, office_location_id, created_at)
VALUES
    (NEWID(), 'BMSCMUHW182KWH001', N'BMS CMU HW1 cho pin 82kWh - Kho HN', 1800000, '150e8400-e29b-41d4-a716-446655440044', '550e8400-e29b-41d4-a716-446655440001', GETDATE()),
    (NEWID(), 'BMSCMUHW182KWH002', N'BMS CMU HW1 cho pin 82kWh - Kho HP', 1800000, '150e8400-e29b-41d4-a716-446655440044', '550e8400-e29b-41d4-a716-446655440002', GETDATE()),
    (NEWID(), 'BMSCMUHW182KWH003', N'BMS CMU HW1 cho pin 82kWh - Kho DN', 1800000, '150e8400-e29b-41d4-a716-446655440044', '550e8400-e29b-41d4-a716-446655440004', GETDATE());

-- Category: BMS.CMU.HW2.87KWH
INSERT INTO parts (id, serial_number, part_name, unit_price, category_id, office_location_id, created_at)
VALUES
    (NEWID(), 'BMSCMUHW287KWH001', N'BMS CMU HW2 cho pin 87.7kWh - Kho HN', 2100000, '150e8400-e29b-41d4-a716-446655440045', '550e8400-e29b-41d4-a716-446655440001', GETDATE()),
    (NEWID(), 'BMSCMUHW287KWH002', N'BMS CMU HW2 cho pin 87.7kWh - Kho HP', 2100000, '150e8400-e29b-41d4-a716-446655440045', '550e8400-e29b-41d4-a716-446655440002', GETDATE()),
    (NEWID(), 'BMSCMUHW287KWH003', N'BMS CMU HW2 cho pin 87.7kWh - Kho DN', 2100000, '150e8400-e29b-41d4-a716-446655440045', '550e8400-e29b-41d4-a716-446655440004', GETDATE());

-- Category: BMS.CMU.HW2.SVC
INSERT INTO parts (id, serial_number, part_name, unit_price, category_id, office_location_id, created_at)
VALUES
    (NEWID(), 'BMSCMUHW2SVC00001', N'Phụ tùng thay thế CMU HW2 - Kho HN', 2200000, '150e8400-e29b-41d4-a716-446655440046', '550e8400-e29b-41d4-a716-446655440001', GETDATE()),
    (NEWID(), 'BMSCMUHW2SVC00002', N'Phụ tùng thay thế CMU HW2 - Kho HP', 2200000, '150e8400-e29b-41d4-a716-446655440046', '550e8400-e29b-41d4-a716-446655440002', GETDATE()),
    (NEWID(), 'BMSCMUHW2SVC00003', N'Phụ tùng thay thế CMU HW2 - Kho DN', 2200000, '150e8400-e29b-41d4-a716-446655440046', '550e8400-e29b-41d4-a716-446655440004', GETDATE());

-- Parts for Upper Battery Casing
-- Category: BAT.CAS.UPR.82KWH
INSERT INTO parts (id, serial_number, part_name, unit_price, category_id, office_location_id, created_at)
VALUES
    (NEWID(), 'BATCASUPR82KWH001', N'Vỏ trên pin 82kWh - Kho HN', 5500000, '150e8400-e29b-41d4-a716-44665544004A', '550e8400-e29b-41d4-a716-446655440001', GETDATE()),
    (NEWID(), 'BATCASUPR82KWH002', N'Vỏ trên pin 82kWh - Kho HP', 5500000, '150e8400-e29b-41d4-a716-44665544004A', '550e8400-e29b-41d4-a716-446655440002', GETDATE()),
    (NEWID(), 'BATCASUPR82KWH003', N'Vỏ trên pin 82kWh - Kho DN', 5500000, '150e8400-e29b-41d4-a716-44665544004A', '550e8400-e29b-41d4-a716-446655440004', GETDATE());

-- Category: BAT.CAS.UPR.87KWH
INSERT INTO parts (id, serial_number, part_name, unit_price, category_id, office_location_id, created_at)
VALUES
    (NEWID(), 'BATCASUPR87KWH001', N'Vỏ trên pin 87.7kWh - Kho HN', 5800000, '150e8400-e29b-41d4-a716-44665544004B', '550e8400-e29b-41d4-a716-446655440001', GETDATE()),
    (NEWID(), 'BATCASUPR87KWH002', N'Vỏ trên pin 87.7kWh - Kho HP', 5800000, '150e8400-e29b-41d4-a716-44665544004B', '550e8400-e29b-41d4-a716-446655440002', GETDATE()),
    (NEWID(), 'BATCASUPR87KWH003', N'Vỏ trên pin 87.7kWh - Kho DN', 5800000, '150e8400-e29b-41d4-a716-44665544004B', '550e8400-e29b-41d4-a716-446655440004', GETDATE());

-- Category: BAT.CAS.UPR.87KWH.AL
INSERT INTO parts (id, serial_number, part_name, unit_price, category_id, office_location_id, created_at)
VALUES
    (NEWID(), 'BATCASUPR87KWHAL1', N'Vỏ trên nhôm pin 87.7kWh - Kho HN', 6500000, '150e8400-e29b-41d4-a716-44665544004C', '550e8400-e29b-41d4-a716-446655440001', GETDATE()),
    (NEWID(), 'BATCASUPR87KWHAL2', N'Vỏ trên nhôm pin 87.7kWh - Kho HP', 6500000, '150e8400-e29b-41d4-a716-44665544004C', '550e8400-e29b-41d4-a716-446655440002', GETDATE()),
    (NEWID(), 'BATCASUPR87KWHAL3', N'Vỏ trên nhôm pin 87.7kWh - Kho DN', 6500000, '150e8400-e29b-41d4-a716-44665544004C', '550e8400-e29b-41d4-a716-446655440004', GETDATE());

-- Parts for Lower Battery Tray
-- Category: BAT.TRAY.LWR.82KWH
INSERT INTO parts (id, serial_number, part_name, unit_price, category_id, office_location_id, created_at)
VALUES
    (NEWID(), 'BATTRAYLWR82KWH01', N'Khay dưới pin 82kWh - Kho HN', 8000000, '150e8400-e29b-41d4-a716-44665544004D', '550e8400-e29b-41d4-a716-446655440001', GETDATE()),
    (NEWID(), 'BATTRAYLWR82KWH02', N'Khay dưới pin 82kWh - Kho HP', 8000000, '150e8400-e29b-41d4-a716-44665544004D', '550e8400-e29b-41d4-a716-446655440002', GETDATE()),
    (NEWID(), 'BATTRAYLWR82KWH03', N'Khay dưới pin 82kWh - Kho DN', 8000000, '150e8400-e29b-41d4-a716-44665544004D', '550e8400-e29b-41d4-a716-446655440004', GETDATE());

-- Category: BAT.TRAY.LWR.87KWH
INSERT INTO parts (id, serial_number, part_name, unit_price, category_id, office_location_id, created_at)
VALUES
    (NEWID(), 'BATTRAYLWR87KWH01', N'Khay dưới pin 87.7kWh - Kho HN', 8500000, '150e8400-e29b-41d4-a716-44665544004E', '550e8400-e29b-41d4-a716-446655440001', GETDATE()),
    (NEWID(), 'BATTRAYLWR87KWH02', N'Khay dưới pin 87.7kWh - Kho HP', 8500000, '150e8400-e29b-41d4-a716-44665544004E', '550e8400-e29b-41d4-a716-446655440002', GETDATE()),
    (NEWID(), 'BATTRAYLWR87KWH03', N'Khay dưới pin 87.7kWh - Kho DN', 8500000, '150e8400-e29b-41d4-a716-44665544004E', '550e8400-e29b-41d4-a716-446655440004', GETDATE());

-- Category: BAT.TRAY.LWR.87KWH.REINF
INSERT INTO parts (id, serial_number, part_name, unit_price, category_id, office_location_id, created_at)
VALUES
    (NEWID(), 'BATTRAYLWR87REINF1', N'Khay dưới pin 87.7kWh, gia cố - Kho HN', 9500000, '150e8400-e29b-41d4-a716-44665544004F', '550e8400-e29b-41d4-a716-446655440001', GETDATE()),
    (NEWID(), 'BATTRAYLWR87REINF2', N'Khay dưới pin 87.7kWh, gia cố - Kho HP', 9500000, '150e8400-e29b-41d4-a716-44665544004F', '550e8400-e29b-41d4-a716-446655440002', GETDATE()),
    (NEWID(), 'BATTRAYLWR87REINF3', N'Khay dưới pin 87.7kWh, gia cố - Kho DN', 9500000, '150e8400-e29b-41d4-a716-44665544004F', '550e8400-e29b-41d4-a716-446655440004', GETDATE());

-- Parts for Sealing Gaskets
-- Category: GASKET.BAT.82KWH
INSERT INTO parts (id, serial_number, part_name, unit_price, category_id, office_location_id, created_at)
VALUES
    (NEWID(), 'GASKETBAT82KWH001', N'Bộ gioăng pin 82kWh - Kho HN', 750000, '150e8400-e29b-41d4-a716-446655440050', '550e8400-e29b-41d4-a716-446655440001', GETDATE()),
    (NEWID(), 'GASKETBAT82KWH002', N'Bộ gioăng pin 82kWh - Kho HP', 750000, '150e8400-e29b-41d4-a716-446655440050', '550e8400-e29b-41d4-a716-446655440002', GETDATE()),
    (NEWID(), 'GASKETBAT82KWH003', N'Bộ gioăng pin 82kWh - Kho DN', 750000, '150e8400-e29b-41d4-a716-446655440050', '550e8400-e29b-41d4-a716-446655440004', GETDATE());

-- Category: GASKET.BAT.87KWH
INSERT INTO parts (id, serial_number, part_name, unit_price, category_id, office_location_id, created_at)
VALUES
    (NEWID(), 'GASKETBAT87KWH001', N'Bộ gioăng pin 87.7kWh - Kho HN', 800000, '150e8400-e29b-41d4-a716-446655440051', '550e8400-e29b-41d4-a716-446655440001', GETDATE()),
    (NEWID(), 'GASKETBAT87KWH002', N'Bộ gioăng pin 87.7kWh - Kho HP', 800000, '150e8400-e29b-41d4-a716-446655440051', '550e8400-e29b-41d4-a716-446655440002', GETDATE()),
    (NEWID(), 'GASKETBAT87KWH003', N'Bộ gioăng pin 87.7kWh - Kho DN', 800000, '150e8400-e29b-41d4-a716-446655440051', '550e8400-e29b-41d4-a716-446655440004', GETDATE());

-- Category: GASKET.BAT.SVC.KIT
INSERT INTO parts (id, serial_number, part_name, unit_price, category_id, office_location_id, created_at)
VALUES
    (NEWID(), 'GASKETBATSVCKIT01', N'Bộ gioăng thay thế dịch vụ - Kho HN', 900000, '150e8400-e29b-41d4-a716-446655440052', '550e8400-e29b-41d4-a716-446655440001', GETDATE()),
    (NEWID(), 'GASKETBATSVCKIT02', N'Bộ gioăng thay thế dịch vụ - Kho HP', 900000, '150e8400-e29b-41d4-a716-446655440052', '550e8400-e29b-41d4-a716-446655440002', GETDATE()),
    (NEWID(), 'GASKETBATSVCKIT03', N'Bộ gioăng thay thế dịch vụ - Kho DN', 900000, '150e8400-e29b-41d4-a716-446655440052', '550e8400-e29b-41d4-a716-446655440004', GETDATE());

-- Parts for Contactors/Relays
-- Category: HVJB.CON.400V.MAIN
INSERT INTO parts (id, serial_number, part_name, unit_price, category_id, office_location_id, created_at)
VALUES
    (NEWID(), 'HVJBCON400VMAIN01', N'Công tắc tơ chính - Kho HN', 2500000, '150e8400-e29b-41d4-a716-446655440055', '550e8400-e29b-41d4-a716-446655440001', GETDATE()),
    (NEWID(), 'HVJBCON400VMAIN02', N'Công tắc tơ chính - Kho HP', 2500000, '150e8400-e29b-41d4-a716-446655440055', '550e8400-e29b-41d4-a716-446655440002', GETDATE()),
    (NEWID(), 'HVJBCON400VMAIN03', N'Công tắc tơ chính - Kho DN', 2500000, '150e8400-e29b-41d4-a716-446655440055', '550e8400-e29b-41d4-a716-446655440004', GETDATE());

-- Category: HVJB.CON.400V.PRECHG
INSERT INTO parts (id, serial_number, part_name, unit_price, category_id, office_location_id, created_at)
VALUES
    (NEWID(), 'HVJBCON400VPRECH1', N'Công tắc tơ sạc trước - Kho HN', 1800000, '150e8400-e29b-41d4-a716-446655440056', '550e8400-e29b-41d4-a716-446655440001', GETDATE()),
    (NEWID(), 'HVJBCON400VPRECH2', N'Công tắc tơ sạc trước - Kho HP', 1800000, '150e8400-e29b-41d4-a716-446655440056', '550e8400-e29b-41d4-a716-446655440002', GETDATE()),
    (NEWID(), 'HVJBCON400VPRECH3', N'Công tắc tơ sạc trước - Kho DN', 1800000, '150e8400-e29b-41d4-a716-446655440056', '550e8400-e29b-41d4-a716-446655440004', GETDATE());

-- Category: HVJB.CON.400V.FASTCHG
INSERT INTO parts (id, serial_number, part_name, unit_price, category_id, office_location_id, created_at)
VALUES
    (NEWID(), 'HVJBCON400VFASTC1', N'Công tắc tơ sạc nhanh - Kho HN', 2800000, '150e8400-e29b-41d4-a716-446655440057', '550e8400-e29b-41d4-a716-446655440001', GETDATE()),
    (NEWID(), 'HVJBCON400VFASTC2', N'Công tắc tơ sạc nhanh - Kho HP', 2800000, '150e8400-e29b-41d4-a716-446655440057', '550e8400-e29b-41d4-a716-446655440002', GETDATE()),
    (NEWID(), 'HVJBCON400VFASTC3', N'Công tắc tơ sạc nhanh - Kho DN', 2800000, '150e8400-e29b-41d4-a716-446655440057', '550e8400-e29b-41d4-a716-446655440004', GETDATE());

-- Parts for High-Voltage Fuses
-- Category: HVJB.FUSE.400A
INSERT INTO parts (id, serial_number, part_name, unit_price, category_id, office_location_id, created_at)
VALUES
    (NEWID(), 'HVJBFUSE400A00001', N'Cầu chì 400A - Kho HN', 500000, '150e8400-e29b-41d4-a716-446655440058', '550e8400-e29b-41d4-a716-446655440001', GETDATE()),
    (NEWID(), 'HVJBFUSE400A00002', N'Cầu chì 400A - Kho HP', 500000, '150e8400-e29b-41d4-a716-446655440058', '550e8400-e29b-41d4-a716-446655440002', GETDATE()),
    (NEWID(), 'HVJBFUSE400A00003', N'Cầu chì 400A - Kho DN', 500000, '150e8400-e29b-41d4-a716-446655440058', '550e8400-e29b-41d4-a716-446655440004', GETDATE());

-- Category: HVJB.FUSE.500A
INSERT INTO parts (id, serial_number, part_name, unit_price, category_id, office_location_id, created_at)
VALUES
    (NEWID(), 'HVJBFUSE500A00001', N'Cầu chì 500A - Kho HN', 650000, '150e8400-e29b-41d4-a716-446655440059', '550e8400-e29b-41d4-a716-446655440001', GETDATE()),
    (NEWID(), 'HVJBFUSE500A00002', N'Cầu chì 500A - Kho HP', 650000, '150e8400-e29b-41d4-a716-446655440059', '550e8400-e29b-41d4-a716-446655440002', GETDATE()),
    (NEWID(), 'HVJBFUSE500A00003', N'Cầu chì 500A - Kho DN', 650000, '150e8400-e29b-41d4-a716-446655440059', '550e8400-e29b-41d4-a716-446655440004', GETDATE());

-- Category: HVJB.FUSE.SVC.KIT
INSERT INTO parts (id, serial_number, part_name, unit_price, category_id, office_location_id, created_at)
VALUES
    (NEWID(), 'HVJBFUSESVCKIT001', N'Bộ cầu chì thay thế dịch vụ - Kho HN', 1000000, '150e8400-e29b-41d4-a716-44665544005A', '550e8400-e29b-41d4-a716-446655440001', GETDATE()),
    (NEWID(), 'HVJBFUSESVCKIT002', N'Bộ cầu chì thay thế dịch vụ - Kho HP', 1000000, '150e8400-e29b-41d4-a716-44665544005A', '550e8400-e29b-41d4-a716-446655440002', GETDATE()),
    (NEWID(), 'HVJBFUSESVCKIT003', N'Bộ cầu chì thay thế dịch vụ - Kho DN', 1000000, '150e8400-e29b-41d4-a716-44665544005A', '550e8400-e29b-41d4-a716-446655440004', GETDATE());

-- Parts for Front Single-Speed Gearbox
-- Category: GR.ASM.FR.R9-1
INSERT INTO parts (id, serial_number, part_name, unit_price, category_id, office_location_id, created_at)
VALUES
    (NEWID(), 'GRASMFRR910000001', N'Hộp số trước, tỉ số truyền 9.1 - Kho HN', 45000000, '150e8400-e29b-41d4-a716-446655440060', '550e8400-e29b-41d4-a716-446655440001', GETDATE()),
    (NEWID(), 'GRASMFRR910000002', N'Hộp số trước, tỉ số truyền 9.1 - Kho HP', 45000000, '150e8400-e29b-41d4-a716-446655440060', '550e8400-e29b-41d4-a716-446655440002', GETDATE()),
    (NEWID(), 'GRASMFRR910000003', N'Hộp số trước, tỉ số truyền 9.1 - Kho DN', 45000000, '150e8400-e29b-41d4-a716-446655440060', '550e8400-e29b-41d4-a716-446655440004', GETDATE());

-- Category: GR.ASM.FR.R9-3.PLUS
INSERT INTO parts (id, serial_number, part_name, unit_price, category_id, office_location_id, created_at)
VALUES
    (NEWID(), 'GRASMFRR93PLUS001', N'Hộp số trước, tỉ số truyền 9.3 (Plus) - Kho HN', 48000000, '150e8400-e29b-41d4-a716-446655440061', '550e8400-e29b-41d4-a716-446655440001', GETDATE()),
    (NEWID(), 'GRASMFRR93PLUS002', N'Hộp số trước, tỉ số truyền 9.3 (Plus) - Kho HP', 48000000, '150e8400-e29b-41d4-a716-446655440061', '550e8400-e29b-41d4-a716-446655440002', GETDATE()),
    (NEWID(), 'GRASMFRR93PLUS003', N'Hộp số trước, tỉ số truyền 9.3 (Plus) - Kho DN', 48000000, '150e8400-e29b-41d4-a716-446655440061', '550e8400-e29b-41d4-a716-446655440004', GETDATE());

-- Category: GR.ASM.FR.R9-1.SVC
INSERT INTO parts (id, serial_number, part_name, unit_price, category_id, office_location_id, created_at)
VALUES
    (NEWID(), 'GRASMFRR91SVC0001', N'Phụ tùng hộp số trước, tỉ số truyền 9.1 - Kho HN', 46000000, '150e8400-e29b-41d4-a716-446655440062', '550e8400-e29b-41d4-a716-446655440001', GETDATE()),
    (NEWID(), 'GRASMFRR91SVC0002', N'Phụ tùng hộp số trước, tỉ số truyền 9.1 - Kho HP', 46000000, '150e8400-e29b-41d4-a716-446655440062', '550e8400-e29b-41d4-a716-446655440002', GETDATE()),
    (NEWID(), 'GRASMFRR91SVC0003', N'Phụ tùng hộp số trước, tỉ số truyền 9.1 - Kho DN', 46000000, '150e8400-e29b-41d4-a716-446655440062', '550e8400-e29b-41d4-a716-446655440004', GETDATE());

-- Parts for Rear Single-Speed Gearbox
-- Category: GR.ASM.RR.R9-1
INSERT INTO parts (id, serial_number, part_name, unit_price, category_id, office_location_id, created_at)
VALUES
    (NEWID(), 'GRASMRRR910000001', N'Hộp số sau, tỉ số truyền 9.1 - Kho HN', 45000000, '150e8400-e29b-41d4-a716-446655440063', '550e8400-e29b-41d4-a716-446655440001', GETDATE()),
    (NEWID(), 'GRASMRRR910000002', N'Hộp số sau, tỉ số truyền 9.1 - Kho HP', 45000000, '150e8400-e29b-41d4-a716-446655440063', '550e8400-e29b-41d4-a716-446655440002', GETDATE()),
    (NEWID(), 'GRASMRRR910000003', N'Hộp số sau, tỉ số truyền 9.1 - Kho DN', 45000000, '150e8400-e29b-41d4-a716-446655440063', '550e8400-e29b-41d4-a716-446655440004', GETDATE());

-- Category: GR.ASM.RR.R9-3.PLUS
INSERT INTO parts (id, serial_number, part_name, unit_price, category_id, office_location_id, created_at)
VALUES
    (NEWID(), 'GRASMRRR93PLUS001', N'Hộp số sau, tỉ số truyền 9.3 (Plus) - Kho HN', 48000000, '150e8400-e29b-41d4-a716-446655440064', '550e8400-e29b-41d4-a716-446655440001', GETDATE()),
    (NEWID(), 'GRASMRRR93PLUS002', N'Hộp số sau, tỉ số truyền 9.3 (Plus) - Kho HP', 48000000, '150e8400-e29b-41d4-a716-446655440064', '550e8400-e29b-41d4-a716-446655440002', GETDATE()),
    (NEWID(), 'GRASMRRR93PLUS003', N'Hộp số sau, tỉ số truyền 9.3 (Plus) - Kho DN', 48000000, '150e8400-e29b-41d4-a716-446655440064', '550e8400-e29b-41d4-a716-446655440004', GETDATE());

-- Category: GR.ASM.RR.R9-1.SVC
INSERT INTO parts (id, serial_number, part_name, unit_price, category_id, office_location_id, created_at)
VALUES
    (NEWID(), 'GRASMRRR91SVC0001', N'Phụ tùng hộp số sau, tỉ số truyền 9.1 - Kho HN', 46000000, '150e8400-e29b-41d4-a716-446655440065', '550e8400-e29b-41d4-a716-446655440001', GETDATE()),
    (NEWID(), 'GRASMRRR91SVC0002', N'Phụ tùng hộp số sau, tỉ số truyền 9.1 - Kho HP', 46000000, '150e8400-e29b-41d4-a716-446655440065', '550e8400-e29b-41d4-a716-446655440002', GETDATE()),
    (NEWID(), 'GRASMRRR91SVC0003', N'Phụ tùng hộp số sau, tỉ số truyền 9.1 - Kho DN', 46000000, '150e8400-e29b-41d4-a716-446655440065', '550e8400-e29b-41d4-a716-446655440004', GETDATE());

-- Parts for Front Differential Unit
-- Category: DIFF.FR.OPEN.V1
INSERT INTO parts (id, serial_number, part_name, unit_price, category_id, office_location_id, created_at)
VALUES
    (NEWID(), 'DIFFROPENV1000001', N'Vi sai mở phía trước, V1 - Kho HN', 12000000, '150e8400-e29b-41d4-a716-446655440068', '550e8400-e29b-41d4-a716-446655440001', GETDATE()),
    (NEWID(), 'DIFFROPENV1000002', N'Vi sai mở phía trước, V1 - Kho HP', 12000000, '150e8400-e29b-41d4-a716-446655440068', '550e8400-e29b-41d4-a716-446655440002', GETDATE()),
    (NEWID(), 'DIFFROPENV1000003', N'Vi sai mở phía trước, V1 - Kho DN', 12000000, '150e8400-e29b-41d4-a716-446655440068', '550e8400-e29b-41d4-a716-446655440004', GETDATE());

-- Category: DIFF.FR.OPEN.V2
INSERT INTO parts (id, serial_number, part_name, unit_price, category_id, office_location_id, created_at)
VALUES
    (NEWID(), 'DIFFROPENV2000001', N'Vi sai mở phía trước, V2 (cải tiến) - Kho HN', 13500000, '150e8400-e29b-41d4-a716-446655440069', '550e8400-e29b-41d4-a716-446655440001', GETDATE()),
    (NEWID(), 'DIFFROPENV2000002', N'Vi sai mở phía trước, V2 (cải tiến) - Kho HP', 13500000, '150e8400-e29b-41d4-a716-446655440069', '550e8400-e29b-41d4-a716-446655440002', GETDATE()),
    (NEWID(), 'DIFFROPENV2000003', N'Vi sai mở phía trước, V2 (cải tiến) - Kho DN', 13500000, '150e8400-e29b-41d4-a716-446655440069', '550e8400-e29b-41d4-a716-446655440004', GETDATE());

-- Category: DIFF.FR.OPEN.SVC.KIT
INSERT INTO parts (id, serial_number, part_name, unit_price, category_id, office_location_id, created_at)
VALUES
    (NEWID(), 'DIFFROPENSVCKIT01', N'Bộ phụ tùng sửa chữa vi sai trước - Kho HN', 4500000, '150e8400-e29b-41d4-a716-44665544006A', '550e8400-e29b-41d4-a716-446655440001', GETDATE()),
    (NEWID(), 'DIFFROPENSVCKIT02', N'Bộ phụ tùng sửa chữa vi sai trước - Kho HP', 4500000, '150e8400-e29b-41d4-a716-44665544006A', '550e8400-e29b-41d4-a716-446655440002', GETDATE()),
    (NEWID(), 'DIFFROPENSVCKIT03', N'Bộ phụ tùng sửa chữa vi sai trước - Kho DN', 4500000, '150e8400-e29b-41d4-a716-44665544006A', '550e8400-e29b-41d4-a716-446655440004', GETDATE());

-- Parts for Rear Differential Unit
-- Category: DIFF.RR.OPEN.V1
INSERT INTO parts (id, serial_number, part_name, unit_price, category_id, office_location_id, created_at)
VALUES
    (NEWID(), 'DIFFRROPENV1000001', N'Vi sai mở phía sau, V1 - Kho HN', 12000000, '150e8400-e29b-41d4-a716-44665544006B', '550e8400-e29b-41d4-a716-446655440001', GETDATE()),
    (NEWID(), 'DIFFRROPENV1000002', N'Vi sai mở phía sau, V1 - Kho HP', 12000000, '150e8400-e29b-41d4-a716-44665544006B', '550e8400-e29b-41d4-a716-446655440002', GETDATE()),
    (NEWID(), 'DIFFRROPENV1000003', N'Vi sai mở phía sau, V1 - Kho DN', 12000000, '150e8400-e29b-41d4-a716-44665544006B', '550e8400-e29b-41d4-a716-446655440004', GETDATE());

-- Category: DIFF.RR.OPEN.V2
INSERT INTO parts (id, serial_number, part_name, unit_price, category_id, office_location_id, created_at)
VALUES
    (NEWID(), 'DIFFRROPENV2000001', N'Vi sai mở phía sau, V2 (cải tiến) - Kho HN', 13500000, '150e8400-e29b-41d4-a716-44665544006C', '550e8400-e29b-41d4-a716-446655440001', GETDATE()),
    (NEWID(), 'DIFFRROPENV2000002', N'Vi sai mở phía sau, V2 (cải tiến) - Kho HP', 13500000, '150e8400-e29b-41d4-a716-44665544006C', '550e8400-e29b-41d4-a716-446655440002', GETDATE()),
    (NEWID(), 'DIFFRROPENV2000003', N'Vi sai mở phía sau, V2 (cải tiến) - Kho DN', 13500000, '150e8400-e29b-41d4-a716-44665544006C', '550e8400-e29b-41d4-a716-446655440004', GETDATE());

-- Category: DIFF.RR.OPEN.SVC.KIT
INSERT INTO parts (id, serial_number, part_name, unit_price, category_id, office_location_id, created_at)
VALUES
    (NEWID(), 'DIFFRROPENSVCKIT01', N'Bộ phụ tùng sửa chữa vi sai sau - Kho HN', 4500000, '150e8400-e29b-41d4-a716-44665544006D', '550e8400-e29b-41d4-a716-446655440001', GETDATE()),
    (NEWID(), 'DIFFRROPENSVCKIT02', N'Bộ phụ tùng sửa chữa vi sai sau - Kho HP', 4500000, '150e8400-e29b-41d4-a716-44665544006D', '550e8400-e29b-41d4-a716-446655440002', GETDATE()),
    (NEWID(), 'DIFFRROPENSVCKIT03', N'Bộ phụ tùng sửa chữa vi sai sau - Kho DN', 4500000, '150e8400-e29b-41d4-a716-44665544006D', '550e8400-e29b-41d4-a716-446655440004', GETDATE());

-- Parts for Front Axle Shaft
-- Category: AXLE.FR.LH
INSERT INTO parts (id, serial_number, part_name, unit_price, category_id, office_location_id, created_at)
VALUES
    (NEWID(), 'AXLEFRLH000000001', N'Trục láp trước bên trái - Kho HN', 3500000, '150e8400-e29b-41d4-a716-446655440070', '550e8400-e29b-41d4-a716-446655440001', GETDATE()),
    (NEWID(), 'AXLEFRLH000000002', N'Trục láp trước bên trái - Kho HP', 3500000, '150e8400-e29b-41d4-a716-446655440070', '550e8400-e29b-41d4-a716-446655440002', GETDATE()),
    (NEWID(), 'AXLEFRLH000000003', N'Trục láp trước bên trái - Kho DN', 3500000, '150e8400-e29b-41d4-a716-446655440070', '550e8400-e29b-41d4-a716-446655440004', GETDATE());

-- Category: AXLE.FR.RH
INSERT INTO parts (id, serial_number, part_name, unit_price, category_id, office_location_id, created_at)
VALUES
    (NEWID(), 'AXLEFRRH000000001', N'Trục láp trước bên phải - Kho HN', 3500000, '150e8400-e29b-41d4-a716-446655440071', '550e8400-e29b-41d4-a716-446655440001', GETDATE()),
    (NEWID(), 'AXLEFRRH000000002', N'Trục láp trước bên phải - Kho HP', 3500000, '150e8400-e29b-41d4-a716-446655440071', '550e8400-e29b-41d4-a716-446655440002', GETDATE()),
    (NEWID(), 'AXLEFRRH000000003', N'Trục láp trước bên phải - Kho DN', 3500000, '150e8400-e29b-41d4-a716-446655440071', '550e8400-e29b-41d4-a716-446655440004', GETDATE());

-- Category: AXLE.FR.ASM.SVC
INSERT INTO parts (id, serial_number, part_name, unit_price, category_id, office_location_id, created_at)
VALUES
    (NEWID(), 'AXLEFRASMSVC00001', N'Bộ trục láp trước thay thế - Kho HN', 6800000, '150e8400-e29b-41d4-a716-446655440072', '550e8400-e29b-41d4-a716-446655440001', GETDATE()),
    (NEWID(), 'AXLEFRASMSVC00002', N'Bộ trục láp trước thay thế - Kho HP', 6800000, '150e8400-e29b-41d4-a716-446655440072', '550e8400-e29b-41d4-a716-446655440002', GETDATE()),
    (NEWID(), 'AXLEFRASMSVC00003', N'Bộ trục láp trước thay thế - Kho DN', 6800000, '150e8400-e29b-41d4-a716-446655440072', '550e8400-e29b-41d4-a716-446655440004', GETDATE());

-- Parts for Rear Axle Shaft
-- Category: AXLE.RR.LH
INSERT INTO parts (id, serial_number, part_name, unit_price, category_id, office_location_id, created_at)
VALUES
    (NEWID(), 'AXLERRLH000000001', N'Trục láp sau bên trái - Kho HN', 3500000, '150e8400-e29b-41d4-a716-446655440073', '550e8400-e29b-41d4-a716-446655440001', GETDATE()),
    (NEWID(), 'AXLERRLH000000002', N'Trục láp sau bên trái - Kho HP', 3500000, '150e8400-e29b-41d4-a716-446655440073', '550e8400-e29b-41d4-a716-446655440002', GETDATE()),
    (NEWID(), 'AXLERRLH000000003', N'Trục láp sau bên trái - Kho DN', 3500000, '150e8400-e29b-41d4-a716-446655440073', '550e8400-e29b-41d4-a716-446655440004', GETDATE());

-- Category: AXLE.RR.RH
INSERT INTO parts (id, serial_number, part_name, unit_price, category_id, office_location_id, created_at)
VALUES
    (NEWID(), 'AXLERRRH000000001', N'Trục láp sau bên phải - Kho HN', 3500000, '150e8400-e29b-41d4-a716-446655440074', '550e8400-e29b-41d4-a716-446655440001', GETDATE()),
    (NEWID(), 'AXLERRRH000000002', N'Trục láp sau bên phải - Kho HP', 3500000, '150e8400-e29b-41d4-a716-446655440074', '550e8400-e29b-41d4-a716-446655440002', GETDATE()),
    (NEWID(), 'AXLERRRH000000003', N'Trục láp sau bên phải - Kho DN', 3500000, '150e8400-e29b-41d4-a716-446655440074', '550e8400-e29b-41d4-a716-446655440004', GETDATE());

-- Category: AXLE.RR.ASM.SVC
INSERT INTO parts (id, serial_number, part_name, unit_price, category_id, office_location_id, created_at)
VALUES
    (NEWID(), 'AXLERRASMSVC00001', N'Bộ trục láp sau thay thế - Kho HN', 6800000, '150e8400-e29b-41d4-a716-446655440075', '550e8400-e29b-41d4-a716-446655440001', GETDATE()),
    (NEWID(), 'AXLERRASMSVC00002', N'Bộ trục láp sau thay thế - Kho HP', 6800000, '150e8400-e29b-41d4-a716-446655440075', '550e8400-e29b-41d4-a716-446655440002', GETDATE()),
    (NEWID(), 'AXLERRASMSVC00003', N'Bộ trục láp sau thay thế - Kho DN', 6800000, '150e8400-e29b-41d4-a716-446655440075', '550e8400-e29b-41d4-a716-446655440004', GETDATE());

-- Parts for Electric Coolant Pump
-- Category: PUMP.BAT.12V.BOSCH
INSERT INTO parts (id, serial_number, part_name, unit_price, category_id, office_location_id, created_at)
VALUES
    (NEWID(), 'PUMPBAT12VBOSCH01', N'Bơm làm mát pin, 12V, Bosch - Kho HN', 2200000, '150e8400-e29b-41d4-a716-44665544007B', '550e8400-e29b-41d4-a716-446655440001', GETDATE()),
    (NEWID(), 'PUMPBAT12VBOSCH02', N'Bơm làm mát pin, 12V, Bosch - Kho HP', 2200000, '150e8400-e29b-41d4-a716-44665544007B', '550e8400-e29b-41d4-a716-446655440002', GETDATE()),
    (NEWID(), 'PUMPBAT12VBOSCH03', N'Bơm làm mát pin, 12V, Bosch - Kho DN', 2200000, '150e8400-e29b-41d4-a716-44665544007B', '550e8400-e29b-41d4-a716-446655440004', GETDATE());

-- Category: PUMP.BAT.12V.VALEO
INSERT INTO parts (id, serial_number, part_name, unit_price, category_id, office_location_id, created_at)
VALUES
    (NEWID(), 'PUMPBAT12VVALEO01', N'Bơm làm mát pin, 12V, Valeo - Kho HN', 2100000, '150e8400-e29b-41d4-a716-44665544007C', '550e8400-e29b-41d4-a716-446655440001', GETDATE()),
    (NEWID(), 'PUMPBAT12VVALEO02', N'Bơm làm mát pin, 12V, Valeo - Kho HP', 2100000, '150e8400-e29b-41d4-a716-44665544007C', '550e8400-e29b-41d4-a716-446655440002', GETDATE()),
    (NEWID(), 'PUMPBAT12VVALEO03', N'Bơm làm mát pin, 12V, Valeo - Kho DN', 2100000, '150e8400-e29b-41d4-a716-44665544007C', '550e8400-e29b-41d4-a716-446655440004', GETDATE());

-- Category: PUMP.BAT.12V.HP.GATES
INSERT INTO parts (id, serial_number, part_name, unit_price, category_id, office_location_id, created_at)
VALUES
    (NEWID(), 'PUMPBAT12VHPGATES1', N'Bơm làm mát pin, 12V, HP, Gates - Kho HN', 2800000, '150e8400-e29b-41d4-a716-44665544007D', '550e8400-e29b-41d4-a716-446655440001', GETDATE()),
    (NEWID(), 'PUMPBAT12VHPGATES2', N'Bơm làm mát pin, 12V, HP, Gates - Kho HP', 2800000, '150e8400-e29b-41d4-a716-44665544007D', '550e8400-e29b-41d4-a716-446655440002', GETDATE()),
    (NEWID(), 'PUMPBAT12VHPGATES3', N'Bơm làm mát pin, 12V, HP, Gates - Kho DN', 2800000, '150e8400-e29b-41d4-a716-44665544007D', '550e8400-e29b-41d4-a716-446655440004', GETDATE());

-- Parts for Battery Chiller/Heat Exchanger
-- Category: CHILLER.BAT.V1
INSERT INTO parts (id, serial_number, part_name, unit_price, category_id, office_location_id, created_at)
VALUES
    (NEWID(), 'CHILLERBATV1000001', N'Bộ làm lạnh pin, V1 - Kho HN', 6000000, '150e8400-e29b-41d4-a716-44665544007E', '550e8400-e29b-41d4-a716-446655440001', GETDATE()),
    (NEWID(), 'CHILLERBATV1000002', N'Bộ làm lạnh pin, V1 - Kho HP', 6000000, '150e8400-e29b-41d4-a716-44665544007E', '550e8400-e29b-41d4-a716-446655440002', GETDATE()),
    (NEWID(), 'CHILLERBATV1000003', N'Bộ làm lạnh pin, V1 - Kho DN', 6000000, '150e8400-e29b-41d4-a716-44665544007E', '550e8400-e29b-41d4-a716-446655440004', GETDATE());

-- Category: CHILLER.BAT.V2.ENH
INSERT INTO parts (id, serial_number, part_name, unit_price, category_id, office_location_id, created_at)
VALUES
    (NEWID(), 'CHILLERBATV2ENH001', N'Bộ làm lạnh pin, V2, hiệu suất tăng cường - Kho HN', 7500000, '150e8400-e29b-41d4-a716-44665544007F', '550e8400-e29b-41d4-a716-446655440001', GETDATE()),
    (NEWID(), 'CHILLERBATV2ENH002', N'Bộ làm lạnh pin, V2, hiệu suất tăng cường - Kho HP', 7500000, '150e8400-e29b-41d4-a716-44665544007F', '550e8400-e29b-41d4-a716-446655440002', GETDATE()),
    (NEWID(), 'CHILLERBATV2ENH003', N'Bộ làm lạnh pin, V2, hiệu suất tăng cường - Kho DN', 7500000, '150e8400-e29b-41d4-a716-44665544007F', '550e8400-e29b-41d4-a716-446655440004', GETDATE());

-- Category: CHILLER.BAT.V1.SVC
INSERT INTO parts (id, serial_number, part_name, unit_price, category_id, office_location_id, created_at)
VALUES
    (NEWID(), 'CHILLERBATV1SVC001', N'Phụ tùng thay thế bộ làm lạnh pin V1 - Kho HN', 6200000, '150e8400-e29b-41d4-a716-446655440080', '550e8400-e29b-41d4-a716-446655440001', GETDATE()),
    (NEWID(), 'CHILLERBATV1SVC002', N'Phụ tùng thay thế bộ làm lạnh pin V1 - Kho HP', 6200000, '150e8400-e29b-41d4-a716-446655440080', '550e8400-e29b-41d4-a716-446655440002', GETDATE()),
    (NEWID(), 'CHILLERBATV1SVC003', N'Phụ tùng thay thế bộ làm lạnh pin V1 - Kho DN', 6200000, '150e8400-e29b-41d4-a716-446655440080', '550e8400-e29b-41d4-a716-446655440004', GETDATE());

-- Parts for Cooling Plates
-- Category: PLATE.COOL.82KWH
INSERT INTO parts (id, serial_number, part_name, unit_price, category_id, office_location_id, created_at)
VALUES
    (NEWID(), 'PLATECOOL82KWH001', N'Bộ tấm làm mát cho pin 82kWh - Kho HN', 3200000, '150e8400-e29b-41d4-a716-446655440081', '550e8400-e29b-41d4-a716-446655440001', GETDATE()),
    (NEWID(), 'PLATECOOL82KWH002', N'Bộ tấm làm mát cho pin 82kWh - Kho HP', 3200000, '150e8400-e29b-41d4-a716-446655440081', '550e8400-e29b-41d4-a716-446655440002', GETDATE()),
    (NEWID(), 'PLATECOOL82KWH003', N'Bộ tấm làm mát cho pin 82kWh - Kho DN', 3200000, '150e8400-e29b-41d4-a716-446655440081', '550e8400-e29b-41d4-a716-446655440004', GETDATE());

-- Category: PLATE.COOL.87KWH
INSERT INTO parts (id, serial_number, part_name, unit_price, category_id, office_location_id, created_at)
VALUES
    (NEWID(), 'PLATECOOL87KWH001', N'Bộ tấm làm mát cho pin 87.7kWh - Kho HN', 3500000, '150e8400-e29b-41d4-a716-446655440082', '550e8400-e29b-41d4-a716-446655440001', GETDATE()),
    (NEWID(), 'PLATECOOL87KWH002', N'Bộ tấm làm mát cho pin 87.7kWh - Kho HP', 3500000, '150e8400-e29b-41d4-a716-446655440082', '550e8400-e29b-41d4-a716-446655440002', GETDATE()),
    (NEWID(), 'PLATECOOL87KWH003', N'Bộ tấm làm mát cho pin 87.7kWh - Kho DN', 3500000, '150e8400-e29b-41d4-a716-446655440082', '550e8400-e29b-41d4-a716-446655440004', GETDATE());

-- Category: PLATE.COOL.SVC.KIT
INSERT INTO parts (id, serial_number, part_name, unit_price, category_id, office_location_id, created_at)
VALUES
    (NEWID(), 'PLATECOOLSVCKIT01', N'Bộ phụ tùng tấm làm mát - Kho HN', 3800000, '150e8400-e29b-41d4-a716-446655440083', '550e8400-e29b-41d4-a716-446655440001', GETDATE()),
    (NEWID(), 'PLATECOOLSVCKIT02', N'Bộ phụ tùng tấm làm mát - Kho HP', 3800000, '150e8400-e29b-41d4-a716-446655440083', '550e8400-e29b-41d4-a716-446655440002', GETDATE()),
    (NEWID(), 'PLATECOOLSVCKIT03', N'Bộ phụ tùng tấm làm mát - Kho DN', 3800000, '150e8400-e29b-41d4-a716-446655440083', '550e8400-e29b-41d4-a716-446655440004', GETDATE());

-- Parts for Drivetrain Radiator
-- Category: RAD.DRV.STD
INSERT INTO parts (id, serial_number, part_name, unit_price, category_id, office_location_id, created_at)
VALUES
    (NEWID(), 'RADDRVSTD00000001', N'Két làm mát tiêu chuẩn - Kho HN', 4000000, '150e8400-e29b-41d4-a716-446655440087', '550e8400-e29b-41d4-a716-446655440001', GETDATE()),
    (NEWID(), 'RADDRVSTD00000002', N'Két làm mát tiêu chuẩn - Kho HP', 4000000, '150e8400-e29b-41d4-a716-446655440087', '550e8400-e29b-41d4-a716-446655440002', GETDATE()),
    (NEWID(), 'RADDRVSTD00000003', N'Két làm mát tiêu chuẩn - Kho DN', 4000000, '150e8400-e29b-41d4-a716-446655440087', '550e8400-e29b-41d4-a716-446655440004', GETDATE());

-- Category: RAD.DRV.TROPICAL
INSERT INTO parts (id, serial_number, part_name, unit_price, category_id, office_location_id, created_at)
VALUES
    (NEWID(), 'RADDRVTROPICAL001', N'Két làm mát cho thị trường nhiệt đới - Kho HN', 4800000, '150e8400-e29b-41d4-a716-446655440088', '550e8400-e29b-41d4-a716-446655440001', GETDATE()),
    (NEWID(), 'RADDRVTROPICAL002', N'Két làm mát cho thị trường nhiệt đới - Kho HP', 4800000, '150e8400-e29b-41d4-a716-446655440088', '550e8400-e29b-41d4-a716-446655440002', GETDATE()),
    (NEWID(), 'RADDRVTROPICAL003', N'Két làm mát cho thị trường nhiệt đới - Kho DN', 4800000, '150e8400-e29b-41d4-a716-446655440088', '550e8400-e29b-41d4-a716-446655440004', GETDATE());

-- Category: RAD.DRV.SVC
INSERT INTO parts (id, serial_number, part_name, unit_price, category_id, office_location_id, created_at)
VALUES
    (NEWID(), 'RADDRVSVC00000001', N'Két làm mát thay thế - Kho HN', 4200000, '150e8400-e29b-41d4-a716-446655440089', '550e8400-e29b-41d4-a716-446655440001', GETDATE()),
    (NEWID(), 'RADDRVSVC00000002', N'Két làm mát thay thế - Kho HP', 4200000, '150e8400-e29b-41d4-a716-446655440089', '550e8400-e29b-41d4-a716-446655440002', GETDATE()),
    (NEWID(), 'RADDRVSVC00000003', N'Két làm mát thay thế - Kho DN', 4200000, '150e8400-e29b-41d4-a716-446655440089', '550e8400-e29b-41d4-a716-446655440004', GETDATE());

-- Parts for Radiator Fan Assembly
-- Category: FAN.RAD.500W.V1
INSERT INTO parts (id, serial_number, part_name, unit_price, category_id, office_location_id, created_at)
VALUES
    (NEWID(), 'FANRAD500WV100001', N'Cụm quạt 500W, V1 - Kho HN', 2500000, '150e8400-e29b-41d4-a716-44665544008A', '550e8400-e29b-41d4-a716-446655440001', GETDATE()),
    (NEWID(), 'FANRAD500WV100002', N'Cụm quạt 500W, V1 - Kho HP', 2500000, '150e8400-e29b-41d4-a716-44665544008A', '550e8400-e29b-41d4-a716-446655440002', GETDATE()),
    (NEWID(), 'FANRAD500WV100003', N'Cụm quạt 500W, V1 - Kho DN', 2500000, '150e8400-e29b-41d4-a716-44665544008A', '550e8400-e29b-41d4-a716-446655440004', GETDATE());

-- Category: FAN.RAD.600W.HP
INSERT INTO parts (id, serial_number, part_name, unit_price, category_id, office_location_id, created_at)
VALUES
    (NEWID(), 'FANRAD600WHP00001', N'Cụm quạt 600W, hiệu suất cao - Kho HN', 3200000, '150e8400-e29b-41d4-a716-44665544008B', '550e8400-e29b-41d4-a716-446655440001', GETDATE()),
    (NEWID(), 'FANRAD600WHP00002', N'Cụm quạt 600W, hiệu suất cao - Kho HP', 3200000, '150e8400-e29b-41d4-a716-44665544008B', '550e8400-e29b-41d4-a716-446655440002', GETDATE()),
    (NEWID(), 'FANRAD600WHP00003', N'Cụm quạt 600W, hiệu suất cao - Kho DN', 3200000, '150e8400-e29b-41d4-a716-44665544008B', '550e8400-e29b-41d4-a716-446655440004', GETDATE());

-- Category: FAN.RAD.500W.SVC
INSERT INTO parts (id, serial_number, part_name, unit_price, category_id, office_location_id, created_at)
VALUES
    (NEWID(), 'FANRAD500WSVC0001', N'Phụ tùng thay thế quạt 500W - Kho HN', 2600000, '150e8400-e29b-41d4-a716-44665544008C', '550e8400-e29b-41d4-a716-446655440001', GETDATE()),
    (NEWID(), 'FANRAD500WSVC0002', N'Phụ tùng thay thế quạt 500W - Kho HP', 2600000, '150e8400-e29b-41d4-a716-44665544008C', '550e8400-e29b-41d4-a716-446655440002', GETDATE()),
    (NEWID(), 'FANRAD500WSVC0003', N'Phụ tùng thay thế quạt 500W - Kho DN', 2600000, '150e8400-e29b-41d4-a716-44665544008C', '550e8400-e29b-41d4-a716-446655440004', GETDATE());

-- Parts for Three-Way Coolant Valve
-- Category: VALVE.3WAY.V1
INSERT INTO parts (id, serial_number, part_name, unit_price, category_id, office_location_id, created_at)
VALUES
    (NEWID(), 'VALVE3WAYV1000001', N'Van 3 ngả, V1 - Kho HN', 950000, '150e8400-e29b-41d4-a716-44665544008D', '550e8400-e29b-41d4-a716-446655440001', GETDATE()),
    (NEWID(), 'VALVE3WAYV1000002', N'Van 3 ngả, V1 - Kho HP', 950000, '150e8400-e29b-41d4-a716-44665544008D', '550e8400-e29b-41d4-a716-446655440002', GETDATE()),
    (NEWID(), 'VALVE3WAYV1000003', N'Van 3 ngả, V1 - Kho DN', 950000, '150e8400-e29b-41d4-a716-44665544008D', '550e8400-e29b-41d4-a716-446655440004', GETDATE());

-- Category: VALVE.3WAY.V2.ENH
INSERT INTO parts (id, serial_number, part_name, unit_price, category_id, office_location_id, created_at)
VALUES
    (NEWID(), 'VALVE3WAYV2ENH001', N'Van 3 ngả, V2, cải tiến - Kho HN', 1200000, '150e8400-e29b-41d4-a716-44665544008E', '550e8400-e29b-41d4-a716-446655440001', GETDATE()),
    (NEWID(), 'VALVE3WAYV2ENH002', N'Van 3 ngả, V2, cải tiến - Kho HP', 1200000, '150e8400-e29b-41d4-a716-44665544008E', '550e8400-e29b-41d4-a716-446655440002', GETDATE()),
    (NEWID(), 'VALVE3WAYV2ENH003', N'Van 3 ngả, V2, cải tiến - Kho DN', 1200000, '150e8400-e29b-41d4-a716-44665544008E', '550e8400-e29b-41d4-a716-446655440004', GETDATE());

-- Category: VALVE.3WAY.V1.SVC
INSERT INTO parts (id, serial_number, part_name, unit_price, category_id, office_location_id, created_at)
VALUES
    (NEWID(), 'VALVE3WAYV1SVC001', N'Phụ tùng thay thế van 3 ngả V1 - Kho HN', 1000000, '150e8400-e29b-41d4-a716-44665544008F', '550e8400-e29b-41d4-a716-446655440001', GETDATE()),
    (NEWID(), 'VALVE3WAYV1SVC002', N'Phụ tùng thay thế van 3 ngả V1 - Kho HP', 1000000, '150e8400-e29b-41d4-a716-44665544008F', '550e8400-e29b-41d4-a716-446655440002', GETDATE()),
    (NEWID(), 'VALVE3WAYV1SVC003', N'Phụ tùng thay thế van 3 ngả V1 - Kho DN', 1000000, '150e8400-e29b-41d4-a716-44665544008F', '550e8400-e29b-41d4-a716-446655440004', GETDATE());

-- Chính sách bảo hành cho các mục con của Front Motor Unit
INSERT INTO policy_coverage_parts (id, policy_id, part_category_id, coverage_conditions, created_at)
VALUES
    (NEWID(), '050e8400-e29b-41d4-a716-446655440003', '150e8400-e29b-41d4-a716-44665544001A', N'Bảo hành 10 năm hoặc 200,000 km cho lỗi từ nhà sản xuất.', GETDATE()),
    (NEWID(), '050e8400-e29b-41d4-a716-446655440003', '150e8400-e29b-41d4-a716-44665544001B', N'Bảo hành 10 năm hoặc 200,000 km cho lỗi từ nhà sản xuất.', GETDATE()),
    (NEWID(), '050e8400-e29b-41d4-a716-446655440003', '150e8400-e29b-41d4-a716-44665544001C', N'Bảo hành 12 tháng hoặc 20,000 km cho phụ tùng thay thế mới.', GETDATE()),
    (NEWID(), '050e8400-e29b-41d4-a716-446655440003', '150e8400-e29b-41d4-a716-44665544001D', N'Bảo hành 6 tháng hoặc 10,000 km cho phụ tùng tái sản xuất.', GETDATE());

-- Chính sách bảo hành cho các mục con của Rear Motor Unit
INSERT INTO policy_coverage_parts (id, policy_id, part_category_id, coverage_conditions, created_at)
VALUES
    (NEWID(), '050e8400-e29b-41d4-a716-446655440003', '150e8400-e29b-41d4-a716-44665544001E', N'Bảo hành 10 năm hoặc 200,000 km cho lỗi từ nhà sản xuất.', GETDATE()),
    (NEWID(), '050e8400-e29b-41d4-a716-446655440003', '150e8400-e29b-41d4-a716-44665544001F', N'Bảo hành 10 năm hoặc 200,000 km cho lỗi từ nhà sản xuất.', GETDATE()),
    (NEWID(), '050e8400-e29b-41d4-a716-446655440003', '150e8400-e29b-41d4-a716-446655440020', N'Bảo hành 12 tháng hoặc 20,000 km cho phụ tùng thay thế mới.', GETDATE()),
    (NEWID(), '050e8400-e29b-41d4-a716-446655440003', '150e8400-e29b-41d4-a716-446655440021', N'Bảo hành 6 tháng hoặc 10,000 km cho phụ tùng tái sản xuất.', GETDATE());

-- Chính sách bảo hành cho các mục con của Inverter/Power Electronics Module
INSERT INTO policy_coverage_parts (id, policy_id, part_category_id, coverage_conditions, created_at)
VALUES
    (NEWID(), '050e8400-e29b-41d4-a716-446655440003', '150e8400-e29b-41d4-a716-446655440026', N'Bảo hành 8 năm hoặc 160,000 km cho lỗi phần cứng.', GETDATE()),
    (NEWID(), '050e8400-e29b-41d4-a716-446655440003', '150e8400-e29b-41d4-a716-446655440027', N'Bảo hành 8 năm hoặc 160,000 km cho lỗi phần cứng.', GETDATE()),
    (NEWID(), '050e8400-e29b-41d4-a716-446655440003', '150e8400-e29b-41d4-a716-446655440028', N'Bảo hành 12 tháng cho phụ tùng thay thế.', GETDATE());

-- Chính sách bảo hành cho các mục con của Rotor Assembly
INSERT INTO policy_coverage_parts (id, policy_id, part_category_id, coverage_conditions, created_at)
VALUES
    (NEWID(), '050e8400-e29b-41d4-a716-446655440003', '150e8400-e29b-41d4-a716-446655440029', N'Bảo hành 10 năm hoặc 200,000 km cho lỗi cơ khí và vật liệu.', GETDATE()),
    (NEWID(), '050e8400-e29b-41d4-a716-446655440003', '150e8400-e29b-41d4-a716-44665544002A', N'Bảo hành 10 năm hoặc 200,000 km cho lỗi cơ khí và vật liệu.', GETDATE()),
    (NEWID(), '050e8400-e29b-41d4-a716-446655440003', '150e8400-e29b-41d4-a716-44665544002B', N'Bảo hành 12 tháng cho phụ tùng thay thế.', GETDATE());

-- Chính sách bảo hành cho các mục con của Stator Assembly
INSERT INTO policy_coverage_parts (id, policy_id, part_category_id, coverage_conditions, created_at)
VALUES
    (NEWID(), '050e8400-e29b-41d4-a716-446655440003', '150e8400-e29b-41d4-a716-44665544002C', N'Bảo hành 10 năm hoặc 200,000 km cho lỗi lớp cách điện và cuộn dây.', GETDATE()),
    (NEWID(), '050e8400-e29b-41d4-a716-446655440003', '150e8400-e29b-41d4-a716-44665544002D', N'Bảo hành 10 năm hoặc 200,000 km cho lỗi lớp cách điện và cuộn dây.', GETDATE()),
    (NEWID(), '050e8400-e29b-41d4-a716-446655440003', '150e8400-e29b-41d4-a716-44665544002E', N'Bảo hành 12 tháng cho phụ tùng thay thế.', GETDATE());

-- Chính sách bảo hành cho các mục con của Motor Position Sensor
INSERT INTO policy_coverage_parts (id, policy_id, part_category_id, coverage_conditions, created_at)
VALUES
    (NEWID(), '050e8400-e29b-41d4-a716-446655440003', '150e8400-e29b-41d4-a716-44665544002F', N'Bảo hành 5 năm hoặc 100,000 km cho lỗi cảm biến.', GETDATE()),
    (NEWID(), '050e8400-e29b-41d4-a716-446655440003', '150e8400-e29b-41d4-a716-446655440030', N'Bảo hành 5 năm hoặc 100,000 km cho lỗi cảm biến.', GETDATE()),
    (NEWID(), '050e8400-e29b-41d4-a716-446655440003', '150e8400-e29b-41d4-a716-446655440031', N'Bảo hành 5 năm hoặc 100,000 km cho lỗi cảm biến.', GETDATE());

-- Chính sách bảo hành cho các mục con của Battery Cell
INSERT INTO policy_coverage_parts (id, policy_id, part_category_id, coverage_conditions, created_at)
VALUES
    (NEWID(), '050e8400-e29b-41d4-a716-446655440003', '150e8400-e29b-41d4-a716-446655440038', N'Bảo hành 10 năm không giới hạn km, đảm bảo dung lượng pin còn lại trên 70%.', GETDATE()),
    (NEWID(), '050e8400-e29b-41d4-a716-446655440003', '150e8400-e29b-41d4-a716-446655440039', N'Bảo hành 10 năm không giới hạn km, đảm bảo dung lượng pin còn lại trên 70%.', GETDATE()),
    (NEWID(), '050e8400-e29b-41d4-a716-446655440003', '150e8400-e29b-41d4-a716-44665544003A', N'Bảo hành 10 năm không giới hạn km, đảm bảo dung lượng pin còn lại trên 70%.', GETDATE()),
    (NEWID(), '050e8400-e29b-41d4-a716-446655440003', '150e8400-e29b-41d4-a716-44665544003B', N'Bảo hành 10 năm không giới hạn km, đảm bảo dung lượng pin còn lại trên 70%.', GETDATE());

-- Chính sách bảo hành cho các mục con của Module Interconnect Board
INSERT INTO policy_coverage_parts (id, policy_id, part_category_id, coverage_conditions, created_at)
VALUES
    (NEWID(), '050e8400-e29b-41d4-a716-446655440003', '150e8400-e29b-41d4-a716-44665544003C', N'Bảo hành 8 năm hoặc 160,000 km cho lỗi mạch điện.', GETDATE()),
    (NEWID(), '050e8400-e29b-41d4-a716-446655440003', '150e8400-e29b-41d4-a716-44665544003D', N'Bảo hành 8 năm hoặc 160,000 km cho lỗi mạch điện.', GETDATE()),
    (NEWID(), '050e8400-e29b-41d4-a716-446655440003', '150e8400-e29b-41d4-a716-44665544003E', N'Bảo hành 8 năm hoặc 160,000 km cho lỗi mạch điện.', GETDATE());

-- Chính sách bảo hành cho các mục con của Master Control Unit (BMS MCU)
INSERT INTO policy_coverage_parts (id, policy_id, part_category_id, coverage_conditions, created_at)
VALUES
    (NEWID(), '050e8400-e29b-41d4-a716-446655440003', '150e8400-e29b-41d4-a716-446655440041', N'Bảo hành 8 năm hoặc 160,000 km cho lỗi phần cứng và phần mềm.', GETDATE()),
    (NEWID(), '050e8400-e29b-41d4-a716-446655440003', '150e8400-e29b-41d4-a716-446655440042', N'Bảo hành 8 năm hoặc 160,000 km cho lỗi phần cứng và phần mềm.', GETDATE()),
    (NEWID(), '050e8400-e29b-41d4-a716-446655440003', '150e8400-e29b-41d4-a716-446655440043', N'Bảo hành 8 năm hoặc 160,000 km cho lỗi phần cứng và phần mềm.', GETDATE());

-- Chính sách bảo hành cho các mục con của Slave/Cell Monitoring Unit (BMS CMU)
INSERT INTO policy_coverage_parts (id, policy_id, part_category_id, coverage_conditions, created_at)
VALUES
    (NEWID(), '050e8400-e29b-41d4-a716-446655440003', '150e8400-e29b-41d4-a716-446655440044', N'Bảo hành 8 năm hoặc 160,000 km cho lỗi phần cứng.', GETDATE()),
    (NEWID(), '050e8400-e29b-41d4-a716-446655440003', '150e8400-e29b-41d4-a716-446655440045', N'Bảo hành 8 năm hoặc 160,000 km cho lỗi phần cứng.', GETDATE()),
    (NEWID(), '050e8400-e29b-41d4-a716-446655440003', '150e8400-e29b-41d4-a716-446655440046', N'Bảo hành 12 tháng cho phụ tùng thay thế.', GETDATE());

-- Chính sách bảo hành cho các mục con của Upper Battery Casing
INSERT INTO policy_coverage_parts (id, policy_id, part_category_id, coverage_conditions, created_at)
VALUES
    (NEWID(), '050e8400-e29b-41d4-a716-446655440003', '150e8400-e29b-41d4-a716-44665544004A', N'Bảo hành 10 năm cho lỗi kết cấu và chống ăn mòn.', GETDATE()),
    (NEWID(), '050e8400-e29b-41d4-a716-446655440003', '150e8400-e29b-41d4-a716-44665544004B', N'Bảo hành 10 năm cho lỗi kết cấu và chống ăn mòn.', GETDATE()),
    (NEWID(), '050e8400-e29b-41d4-a716-446655440003', '150e8400-e29b-41d4-a716-44665544004C', N'Bảo hành 10 năm cho lỗi kết cấu và chống ăn mòn.', GETDATE());

-- Chính sách bảo hành cho các mục con của Lower Battery Tray
INSERT INTO policy_coverage_parts (id, policy_id, part_category_id, coverage_conditions, created_at)
VALUES
    (NEWID(), '050e8400-e29b-41d4-a716-446655440003', '150e8400-e29b-41d4-a716-44665544004D', N'Bảo hành 10 năm cho lỗi kết cấu và chống ăn mòn.', GETDATE()),
    (NEWID(), '050e8400-e29b-41d4-a716-446655440003', '150e8400-e29b-41d4-a716-44665544004E', N'Bảo hành 10 năm cho lỗi kết cấu và chống ăn mòn.', GETDATE()),
    (NEWID(), '050e8400-e29b-41d4-a716-446655440003', '150e8400-e29b-41d4-a716-44665544004F', N'Bảo hành 10 năm cho lỗi kết cấu và chống ăn mòn.', GETDATE());

-- Chính sách bảo hành cho các mục con của Sealing Gaskets
INSERT INTO policy_coverage_parts (id, policy_id, part_category_id, coverage_conditions, created_at)
VALUES
    (NEWID(), '050e8400-e29b-41d4-a716-446655440003', '150e8400-e29b-41d4-a716-446655440050', N'Bảo hành 8 năm hoặc 160,000 km cho khả năng làm kín.', GETDATE()),
    (NEWID(), '050e8400-e29b-41d4-a716-446655440003', '150e8400-e29b-41d4-a716-446655440051', N'Bảo hành 8 năm hoặc 160,000 km cho khả năng làm kín.', GETDATE()),
    (NEWID(), '050e8400-e29b-41d4-a716-446655440003', '150e8400-e29b-41d4-a716-446655440052', N'Bảo hành 12 tháng cho bộ phụ tùng thay thế.', GETDATE());

-- Chính sách bảo hành cho các mục con của Contactors/Relays
INSERT INTO policy_coverage_parts (id, policy_id, part_category_id, coverage_conditions, created_at)
VALUES
    (NEWID(), '050e8400-e29b-41d4-a716-446655440003', '150e8400-e29b-41d4-a716-446655440055', N'Bảo hành 5 năm hoặc 100,000 km cho lỗi cơ điện.', GETDATE()),
    (NEWID(), '050e8400-e29b-41d4-a716-446655440003', '150e8400-e29b-41d4-a716-446655440056', N'Bảo hành 5 năm hoặc 100,000 km cho lỗi cơ điện.', GETDATE()),
    (NEWID(), '050e8400-e29b-41d4-a716-446655440003', '150e8400-e29b-41d4-a716-446655440057', N'Bảo hành 5 năm hoặc 100,000 km cho lỗi cơ điện.', GETDATE());

-- Chính sách bảo hành cho các mục con của High-Voltage Fuses
INSERT INTO policy_coverage_parts (id, policy_id, part_category_id, coverage_conditions, created_at)
VALUES
    (NEWID(), '050e8400-e29b-41d4-a716-446655440003', '150e8400-e29b-41d4-a716-446655440058', N'Bảo hành 12 tháng cho lỗi từ nhà sản xuất, không bao gồm trường hợp cháy do quá tải.', GETDATE()),
    (NEWID(), '050e8400-e29b-41d4-a716-446655440003', '150e8400-e29b-41d4-a716-446655440059', N'Bảo hành 12 tháng cho lỗi từ nhà sản xuất, không bao gồm trường hợp cháy do quá tải.', GETDATE()),
    (NEWID(), '050e8400-e29b-41d4-a716-446655440003', '150e8400-e29b-41d4-a716-44665544005A', N'Bảo hành 12 tháng cho bộ phụ tùng thay thế.', GETDATE());

-- Chính sách bảo hành cho các mục con của Front Single-Speed Gearbox
INSERT INTO policy_coverage_parts (id, policy_id, part_category_id, coverage_conditions, created_at)
VALUES
    (NEWID(), '050e8400-e29b-41d4-a716-446655440003', '150e8400-e29b-41d4-a716-446655440060', N'Bảo hành 10 năm hoặc 200,000 km cho lỗi cơ khí.', GETDATE()),
    (NEWID(), '050e8400-e29b-41d4-a716-446655440003', '150e8400-e29b-41d4-a716-446655440061', N'Bảo hành 10 năm hoặc 200,000 km cho lỗi cơ khí.', GETDATE()),
    (NEWID(), '050e8400-e29b-41d4-a716-446655440003', '150e8400-e29b-41d4-a716-446655440062', N'Bảo hành 12 tháng cho phụ tùng thay thế.', GETDATE());

-- Chính sách bảo hành cho các mục con của Rear Single-Speed Gearbox
INSERT INTO policy_coverage_parts (id, policy_id, part_category_id, coverage_conditions, created_at)
VALUES
    (NEWID(), '050e8400-e29b-41d4-a716-446655440003', '150e8400-e29b-41d4-a716-446655440063', N'Bảo hành 10 năm hoặc 200,000 km cho lỗi cơ khí.', GETDATE()),
    (NEWID(), '050e8400-e29b-41d4-a716-446655440003', '150e8400-e29b-41d4-a716-446655440064', N'Bảo hành 10 năm hoặc 200,000 km cho lỗi cơ khí.', GETDATE()),
    (NEWID(), '050e8400-e29b-41d4-a716-446655440003', '150e8400-e29b-41d4-a716-446655440065', N'Bảo hành 12 tháng cho phụ tùng thay thế.', GETDATE());

-- Chính sách bảo hành cho các mục con của Front Differential Unit
INSERT INTO policy_coverage_parts (id, policy_id, part_category_id, coverage_conditions, created_at)
VALUES
    (NEWID(), '050e8400-e29b-41d4-a716-446655440003', '150e8400-e29b-41d4-a716-446655440068', N'Bảo hành 10 năm hoặc 200,000 km cho lỗi cơ khí.', GETDATE()),
    (NEWID(), '050e8400-e29b-41d4-a716-446655440003', '150e8400-e29b-41d4-a716-446655440069', N'Bảo hành 10 năm hoặc 200,000 km cho lỗi cơ khí.', GETDATE()),
    (NEWID(), '050e8400-e29b-41d4-a716-446655440003', '150e8400-e29b-41d4-a716-44665544006A', N'Bảo hành 12 tháng cho bộ phụ tùng sửa chữa.', GETDATE());

-- Chính sách bảo hành cho các mục con của Rear Differential Unit
INSERT INTO policy_coverage_parts (id, policy_id, part_category_id, coverage_conditions, created_at)
VALUES
    (NEWID(), '050e8400-e29b-41d4-a716-446655440003', '150e8400-e29b-41d4-a716-44665544006B', N'Bảo hành 10 năm hoặc 200,000 km cho lỗi cơ khí.', GETDATE()),
    (NEWID(), '050e8400-e29b-41d4-a716-446655440003', '150e8400-e29b-41d4-a716-44665544006C', N'Bảo hành 10 năm hoặc 200,000 km cho lỗi cơ khí.', GETDATE()),
    (NEWID(), '050e8400-e29b-41d4-a716-446655440003', '150e8400-e29b-41d4-a716-44665544006D', N'Bảo hành 12 tháng cho bộ phụ tùng sửa chữa.', GETDATE());

-- Chính sách bảo hành cho các mục con của Front Axle Shaft
INSERT INTO policy_coverage_parts (id, policy_id, part_category_id, coverage_conditions, created_at)
VALUES
    (NEWID(), '050e8400-e29b-41d4-a716-446655440003', '150e8400-e29b-41d4-a716-446655440070', N'Bảo hành 5 năm hoặc 100,000 km cho lỗi vật liệu.', GETDATE()),
    (NEWID(), '050e8400-e29b-41d4-a716-446655440003', '150e8400-e29b-41d4-a716-446655440071', N'Bảo hành 5 năm hoặc 100,000 km cho lỗi vật liệu.', GETDATE()),
    (NEWID(), '050e8400-e29b-41d4-a716-446655440003', '150e8400-e29b-41d4-a716-446655440072', N'Bảo hành 12 tháng cho bộ phụ tùng thay thế.', GETDATE());

-- Chính sách bảo hành cho các mục con của Rear Axle Shaft
INSERT INTO policy_coverage_parts (id, policy_id, part_category_id, coverage_conditions, created_at)
VALUES
    (NEWID(), '050e8400-e29b-41d4-a716-446655440003', '150e8400-e29b-41d4-a716-446655440073', N'Bảo hành 5 năm hoặc 100,000 km cho lỗi vật liệu.', GETDATE()),
    (NEWID(), '050e8400-e29b-41d4-a716-446655440003', '150e8400-e29b-41d4-a716-446655440074', N'Bảo hành 5 năm hoặc 100,000 km cho lỗi vật liệu.', GETDATE()),
    (NEWID(), '050e8400-e29b-41d4-a716-446655440003', '150e8400-e29b-41d4-a716-446655440075', N'Bảo hành 12 tháng cho bộ phụ tùng thay thế.', GETDATE());

-- Chính sách bảo hành cho các mục con của Electric Coolant Pump
INSERT INTO policy_coverage_parts (id, policy_id, part_category_id, coverage_conditions, created_at)
VALUES
    (NEWID(), '050e8400-e29b-41d4-a716-446655440003', '150e8400-e29b-41d4-a716-44665544007B', N'Bảo hành 3 năm hoặc 60,000 km.', GETDATE()),
    (NEWID(), '050e8400-e29b-41d4-a716-446655440003', '150e8400-e29b-41d4-a716-44665544007C', N'Bảo hành 3 năm hoặc 60,000 km.', GETDATE()),
    (NEWID(), '050e8400-e29b-41d4-a716-446655440003', '150e8400-e29b-41d4-a716-44665544007D', N'Bảo hành 4 năm hoặc 80,000 km.', GETDATE());

-- Chính sách bảo hành cho các mục con của Battery Chiller/Heat Exchanger
INSERT INTO policy_coverage_parts (id, policy_id, part_category_id, coverage_conditions, created_at)
VALUES
    (NEWID(), '050e8400-e29b-41d4-a716-446655440003', '150e8400-e29b-41d4-a716-44665544007E', N'Bảo hành 8 năm hoặc 160,000 km cho lỗi rò rỉ hoặc giảm hiệu suất.', GETDATE()),
    (NEWID(), '050e8400-e29b-41d4-a716-446655440003', '150e8400-e29b-41d4-a716-44665544007F', N'Bảo hành 8 năm hoặc 160,000 km cho lỗi rò rỉ hoặc giảm hiệu suất.', GETDATE()),
    (NEWID(), '050e8400-e29b-41d4-a716-446655440003', '150e8400-e29b-41d4-a716-446655440080', N'Bảo hành 12 tháng cho phụ tùng thay thế.', GETDATE());

-- Chính sách bảo hành cho các mục con của Cooling Plates
INSERT INTO policy_coverage_parts (id, policy_id, part_category_id, coverage_conditions, created_at)
VALUES
    (NEWID(), '050e8400-e29b-41d4-a716-446655440003', '150e8400-e29b-41d4-a716-446655440081', N'Bảo hành 8 năm hoặc 160,000 km cho lỗi rò rỉ.', GETDATE()),
    (NEWID(), '050e8400-e29b-41d4-a716-446655440003', '150e8400-e29b-41d4-a716-446655440082', N'Bảo hành 8 năm hoặc 160,000 km cho lỗi rò rỉ.', GETDATE()),
    (NEWID(), '050e8400-e29b-41d4-a716-446655440003', '150e8400-e29b-41d4-a716-446655440083', N'Bảo hành 12 tháng cho bộ phụ tùng thay thế.', GETDATE());

-- Chính sách bảo hành cho các mục con của Drivetrain Radiator
INSERT INTO policy_coverage_parts (id, policy_id, part_category_id, coverage_conditions, created_at)
VALUES
    (NEWID(), '050e8400-e29b-41d4-a716-446655440003', '150e8400-e29b-41d4-a716-446655440087', N'Bảo hành 5 năm hoặc 100,000 km cho lỗi rò rỉ.', GETDATE()),
    (NEWID(), '050e8400-e29b-41d4-a716-446655440003', '150e8400-e29b-41d4-a716-446655440088', N'Bảo hành 5 năm hoặc 100,000 km cho lỗi rò rỉ.', GETDATE()),
    (NEWID(), '050e8400-e29b-41d4-a716-446655440003', '150e8400-e29b-41d4-a716-446655440089', N'Bảo hành 12 tháng cho phụ tùng thay thế.', GETDATE());

-- Chính sách bảo hành cho các mục con của Radiator Fan Assembly
INSERT INTO policy_coverage_parts (id, policy_id, part_category_id, coverage_conditions, created_at)
VALUES
    (NEWID(), '050e8400-e29b-41d4-a716-446655440003', '150e8400-e29b-41d4-a716-44665544008A', N'Bảo hành 3 năm hoặc 60,000 km cho lỗi động cơ quạt.', GETDATE()),
    (NEWID(), '050e8400-e29b-41d4-a716-446655440003', '150e8400-e29b-41d4-a716-44665544008B', N'Bảo hành 3 năm hoặc 60,000 km cho lỗi động cơ quạt.', GETDATE()),
    (NEWID(), '050e8400-e29b-41d4-a716-446655440003', '150e8400-e29b-41d4-a716-44665544008C', N'Bảo hành 12 tháng cho phụ tùng thay thế.', GETDATE());

-- Chính sách bảo hành cho các mục con của Three-Way Coolant Valve
INSERT INTO policy_coverage_parts (id, policy_id, part_category_id, coverage_conditions, created_at)
VALUES
    (NEWID(), '050e8400-e29b-41d4-a716-446655440003', '150e8400-e29b-41d4-a716-44665544008D', N'Bảo hành 3 năm hoặc 60,000 km cho lỗi kẹt van hoặc rò rỉ.', GETDATE()),
    (NEWID(), '050e8400-e29b-41d4-a716-446655440003', '150e8400-e29b-41d4-a716-44665544008E', N'Bảo hành 3 năm hoặc 60,000 km cho lỗi kẹt van hoặc rò rỉ.', GETDATE()),
    (NEWID(), '050e8400-e29b-41d4-a716-446655440003', '150e8400-e29b-41d4-a716-44665544008F', N'Bảo hành 12 tháng cho phụ tùng thay thế.', GETDATE());

COMMIT TRANSACTION;

PRINT 'Seed data inserted successfully!';