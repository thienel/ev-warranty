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
    ('050e8400-e29b-41d4-a716-446655440003', 'VF8_2025_WARRANTY', 120, 200000, N'Bảo hành 10 năm hoặc 200,000 km (tùy điều kiện nào đến trước). Pin và động cơ điện được bảo hành trọn đời theo điều kiện sử dụng. Bảo dưỡng định kỳ theo khuyến cáo nhà sản xuất. Không áp dụng cho xe sử dụng vào mục đích thương mại.', GETDATE()),
    ('050e8400-e29b-41d4-a716-446655440018', 'VF3_2024_WARRANTY', 96, 160000, N'Bảo hành 8 năm hoặc 160,000 km (tùy điều kiện nào đến trước). Pin được bảo hành 10 năm. Bảo dưỡng định kỳ miễn phí 2 năm đầu. Không áp dụng cho xe sử dụng vào mục đích thương mại.', GETDATE()),

    -- BYD Atto 3
    ('050e8400-e29b-41d4-a716-446655440020', 'BYD_ATTO3_2022_WARRANTY', 72, 150000, N'Bảo hành 6 năm hoặc 150,000 km (tùy điều kiện nào đến trước). Pin và hệ thống điện được bảo hành 8 năm hoặc 150,000 km. Bảo dưỡng định kỳ theo khuyến cáo. Áp dụng cho xe sử dụng cá nhân.', GETDATE());
    
-- Insert Vehicle Models (EV models from VinFast, BYD, and Mercedes-Benz available in Vietnam)

-- =========================================================
-- * DECLARE GUIDs for Vehicle Models
-- =========================================================
DECLARE
    -- VinFast 4401xx
   
    @vf8_2025 UNIQUEIDENTIFIER = '650e8400-e29b-41d4-a716-446655440103',
    @vf3_2024 UNIQUEIDENTIFIER = '650e8400-e29b-41d4-a716-446655440118',

    -- BYD 4402xx
    @byd_atto3_2022 UNIQUEIDENTIFIER = '650e8400-e29b-41d4-a716-446655440201';
    

-- =========================================================
-- * INSERT DATA (VinFast / BYD / Mercedes-Benz)
-- =========================================================

INSERT INTO vehicle_models (id, brand, model_name, year, created_at, policy_id)
VALUES
    -- VinFast
    
    (@vf8_2025, 'VinFast', 'VF8', 2025, GETDATE(), '050e8400-e29b-41d4-a716-446655440003'),
    (@vf3_2024, 'VinFast', 'VF3', 2024, GETDATE(), '050e8400-e29b-41d4-a716-446655440018'),

    -- BYD
    (@byd_atto3_2022, 'BYD', 'Atto 3', 2022, GETDATE(), '050e8400-e29b-41d4-a716-446655440020');
   

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
    ('850e8400-e29b-41d4-a716-44665544000a', N'Thảo', N'Đinh Thị', '0900123456', 'thao.dt@gmail.com', N'246 Âu Cơ, Quận Tân Phú, TP.HCM', GETDATE());

INSERT INTO vehicles (id, vin, license_plate, customer_id, model_id, purchase_date, created_at)
VALUES
    ('750e8400-e29b-41d4-a716-446655440001', 'VF8ABC123XYZ45678', '51K-55555', '850e8400-e29b-41d4-a716-446655440001', @vf8_2025, '2023-03-15', GETDATE()),
    ('750e8400-e29b-41d4-a716-446655440002', 'RL4VF8S23RF123456', '51K-12345', '850e8400-e29b-41d4-a716-446655440001', @vf8_2025, '2025-01-20', GETDATE()),
    ('750e8400-e29b-41d4-a716-446655440003', 'RL4VF8S23RF123457', '51K-23456', '850e8400-e29b-41d4-a716-446655440002', @vf8_2025, '2025-02-10', GETDATE()),
    ('750e8400-e29b-41d4-a716-446655440004', 'RL4VF8S23RF123458', '51K-34567', '850e8400-e29b-41d4-a716-446655440003', @vf8_2025, '2025-03-05', GETDATE()),
    ('750e8400-e29b-41d4-a716-446655440005', 'RL4VF8S23RF123459', '51K-45678', '850e8400-e29b-41d4-a716-446655440004', @vf8_2025, '2025-04-18', GETDATE()),
    ('750e8400-e29b-41d4-a716-446655440006', 'RL4VF8S23RF123460', '51K-56789', '850e8400-e29b-41d4-a716-446655440005', @vf8_2025, '2025-05-22', GETDATE()),
    ('750e8400-e29b-41d4-a716-446655440007', 'RL4VF8S23RF123461', '51L-67890', '850e8400-e29b-41d4-a716-446655440006', @vf8_2025, '2025-06-11', GETDATE()),
    ('750e8400-e29b-41d4-a716-446655440008', 'RL4VF3S22RE123462', '51L-78901', '850e8400-e29b-41d4-a716-446655440006', @vf3_2024, '2024-07-15', GETDATE()),
    ('750e8400-e29b-41d4-a716-446655440009', 'L1BAT3S20RD123463', '51L-89012', '850e8400-e29b-41d4-a716-446655440006', @byd_atto3_2022, '2022-08-25', GETDATE());

-- Declare parent category IDs for readability
DECLARE @electric_drivetrain_id UNIQUEIDENTIFIER = '150e8400-e29b-41d4-a716-446655440001';
DECLARE @hv_battery_system_id UNIQUEIDENTIFIER = '150e8400-e29b-41d4-a716-446655440002';
DECLARE @chassis_suspension_id UNIQUEIDENTIFIER = '150e8400-e29b-41d4-a716-446655440003';
DECLARE @steering_system_id UNIQUEIDENTIFIER = '150e8400-e29b-41d4-a716-446655440004';
DECLARE @braking_system_id UNIQUEIDENTIFIER = '150e8400-e29b-41d4-a716-446655440005';
DECLARE @body_doors_id UNIQUEIDENTIFIER = '150e8400-e29b-41d4-a716-446655440006';
DECLARE @body_exterior_panels_id UNIQUEIDENTIFIER = '150e8400-e29b-41d4-a716-446655440007';
DECLARE @lighting_system_id UNIQUEIDENTIFIER = '150e8400-e29b-41d4-a716-446655440008';
DECLARE @glazing_mirrors_id UNIQUEIDENTIFIER = '150e8400-e29b-41d4-a716-446655440009';
DECLARE @interior_seating_restraints_id UNIQUEIDENTIFIER = '150e8400-e29b-41d4-a716-44665544000a';
DECLARE @interior_cockpit_id UNIQUEIDENTIFIER = '150e8400-e29b-41d4-a716-44665544000b';
DECLARE @electronics_sensors_id UNIQUEIDENTIFIER = '150e8400-e29b-41d4-a716-44665544000c';
DECLARE @charging_system_id UNIQUEIDENTIFIER = '150e8400-e29b-41d4-a716-44665544000d';


INSERT INTO part_categories (id, category_name, description, parent_category_id, created_at)
VALUES
    -- Cấp 1 Categories
    (@electric_drivetrain_id, N'Hệ thống Truyền động Điện', N'Electric Drivetrain', NULL, GETDATE()),
    (@hv_battery_system_id, N'Hệ thống Pin Cao áp', N'High-Voltage Battery System', NULL, GETDATE()),
    (@chassis_suspension_id, N'Hệ thống Khung gầm & Treo', N'Chassis & Suspension', NULL, GETDATE()),
    (@steering_system_id, N'Hệ thống Lái', N'Steering System', NULL, GETDATE()),
    (@braking_system_id, N'Hệ thống Phanh', N'Braking System', NULL, GETDATE()),
    (@body_doors_id, N'Thân vỏ - Cửa', N'Body - Doors', NULL, GETDATE()),
    (@body_exterior_panels_id, N'Thân vỏ - Tấm ốp Ngoài', N'Body - Exterior Panels', NULL, GETDATE()),
    (@lighting_system_id, N'Hệ thống Chiếu sáng', N'Lighting System', NULL, GETDATE()),
    (@glazing_mirrors_id, N'Kính & Gương', N'Glazing & Mirrors', NULL, GETDATE()),
    (@interior_seating_restraints_id, N'Nội thất - Ghế & An toàn', N'Interior - Seating & Restraints', NULL, GETDATE()),
    (@interior_cockpit_id, N'Nội thất - Khoang lái', N'Interior - Cockpit', NULL, GETDATE()),
    (@electronics_sensors_id, N'Điện tử & Cảm biến', N'Electronics & Sensors', NULL, GETDATE()),
    (@charging_system_id, N'Hệ thống Sạc', N'Charging System', NULL, GETDATE()),

    -- Cấp 2 Categories: 1. Hệ thống Truyền động Điện
    ('150e8400-e29b-41d4-a716-446655440011', N'Động cơ Đồng bộ Nam châm Vĩnh cửu (PMSM)', N'Permanent Magnet Synchronous Motor', @electric_drivetrain_id, GETDATE()),
    ('150e8400-e29b-41d4-a716-446655440012', N'Động cơ Cảm ứng (IM - Induction Motor)', N'Induction Motor', @electric_drivetrain_id, GETDATE()),
    ('150e8400-e29b-41d4-a716-446655440013', N'Động cơ Từ trở Đồng bộ (SynRM)', N'Synchronous Reluctance Motor', @electric_drivetrain_id, GETDATE()),
    ('150e8400-e29b-41d4-a716-446655440014', N'Bộ Biến tần / Điện tử Công suất (Inverter / PEM)', N'Inverter / Power Electronics Module', @electric_drivetrain_id, GETDATE()),
    ('150e8400-e29b-41d4-a716-446655440015', N'Hộp Giảm tốc Một cấp (Single-Speed Gear Reducer)', N'Single-Speed Gear Reducer', @electric_drivetrain_id, GETDATE()),

    -- Cấp 2 Categories: 2. Hệ thống Pin Cao áp
    ('150e8400-e29b-41d4-a716-446655440021', N'Pin Lithium-ion (NMC - Nickel Manganese Cobalt)', N'Lithium-ion Battery (NMC)', @hv_battery_system_id, GETDATE()),
    ('150e8400-e29b-41d4-a716-446655440022', N'Pin LFP (Lithium Iron Phosphate)', N'LFP Battery (Lithium Iron Phosphate)', @hv_battery_system_id, GETDATE()),
    ('150e8400-e29b-41d4-a716-446655440023', N'Pin Blade (Công nghệ của BYD)', N'Blade Battery (BYD Technology)', @hv_battery_system_id, GETDATE()),
    ('150e8400-e29b-41d4-a716-446655440024', N'Mô-đun Pin Cao áp (High-Voltage Battery Module)', N'High-Voltage Battery Module', @hv_battery_system_id, GETDATE()),
    ('150e8400-e29b-41d4-a716-446655440025', N'Hệ thống Quản lý Pin (BMS - Battery Management System)', N'Battery Management System', @hv_battery_system_id, GETDATE()),

    -- Cấp 2 Categories: 3. Hệ thống Khung gầm & Treo
    ('150e8400-e29b-41d4-a716-446655440031', N'Hệ thống Treo Lò xo Thép', N'Steel Spring Suspension', @chassis_suspension_id, GETDATE()),
    ('150e8400-e29b-41d4-a716-446655440032', N'Hệ thống Treo Khí nén', N'Air Suspension System', @chassis_suspension_id, GETDATE()),
    ('150e8400-e29b-41d4-a716-446655440033', N'Giảm xóc Thủy lực', N'Hydraulic Damper', @chassis_suspension_id, GETDATE()),
    ('150e8400-e29b-41d4-a716-446655440034', N'Giảm xóc Thích ứng', N'Adaptive Damper', @chassis_suspension_id, GETDATE()),
    ('150e8400-e29b-41d4-a716-446655440035', N'Khung Phụ Trước', N'Front Subframe', @chassis_suspension_id, GETDATE()),
    ('150e8400-e29b-41d4-a716-446655440036', N'Khung Phụ Sau', N'Rear Subframe', @chassis_suspension_id, GETDATE()),

    -- Cấp 2 Categories: 4. Hệ thống Lái
    ('150e8400-e29b-41d4-a716-446655440041', N'Thước lái Trợ lực Điện (EPS - Electric Power Steering Rack)', N'Electric Power Steering Rack', @steering_system_id, GETDATE()),
    ('150e8400-e29b-41d4-a716-446655440042', N'Trục Lái', N'Steering Column', @steering_system_id, GETDATE()),
    ('150e8400-e29b-41d4-a716-446655440043', N'Vô-lăng', N'Steering Wheel', @steering_system_id, GETDATE()),

    -- Cấp 2 Categories: 5. Hệ thống Phanh
    ('150e8400-e29b-41d4-a716-446655440051', N'Ngàm Phanh', N'Brake Caliper', @braking_system_id, GETDATE()),
    ('150e8400-e29b-41d4-a716-446655440052', N'Đĩa Phanh', N'Brake Disc / Rotor', @braking_system_id, GETDATE()),
    ('150e8400-e29b-41d4-a716-446655440053', N'Má Phanh', N'Brake Pad', @braking_system_id, GETDATE()),
    ('150e8400-e29b-41d4-a716-446655440054', N'Hệ thống Phanh Tái tạo Năng lượng', N'Regenerative Braking System', @braking_system_id, GETDATE()),
    ('150e8400-e29b-41d4-a716-446655440055', N'Phanh Tay Điện tử', N'Electric Parking Brake', @braking_system_id, GETDATE()),

    -- Cấp 2 Categories: 6. Thân vỏ - Cửa
    ('150e8400-e29b-41d4-a716-446655440061', N'Cửa Trước - Trái', N'Front Door - Left', @body_doors_id, GETDATE()),
    ('150e8400-e29b-41d4-a716-446655440062', N'Cửa Trước - Phải', N'Front Door - Right', @body_doors_id, GETDATE()),
    ('150e8400-e29b-41d4-a716-446655440063', N'Cửa Sau - Trái', N'Rear Door - Left', @body_doors_id, GETDATE()),
    ('150e8400-e29b-41d4-a716-446655440064', N'Cửa Sau - Phải', N'Rear Door - Right', @body_doors_id, GETDATE()),
    ('150e8400-e29b-41d4-a716-446655440065', N'Cửa Cốp / Cửa Hậu', N'Trunk Lid / Liftgate', @body_doors_id, GETDATE()),

    -- Cấp 2 Categories: 7. Thân vỏ - Tấm ốp Ngoài
    ('150e8400-e29b-41d4-a716-446655440071', N'Cản Trước', N'Front Bumper', @body_exterior_panels_id, GETDATE()),
    ('150e8400-e29b-41d4-a716-446655440072', N'Cản Sau', N'Rear Bumper', @body_exterior_panels_id, GETDATE()),
    ('150e8400-e29b-41d4-a716-446655440073', N'Nắp Ca-pô', N'Hood Panel', @body_exterior_panels_id, GETDATE()),
    ('150e8400-e29b-41d4-a716-446655440074', N'Vè Trước', N'Front Fender', @body_exterior_panels_id, GETDATE()),
    ('150e8400-e29b-41d4-a716-446655440075', N'Tấm Mái xe', N'Roof Panel', @body_exterior_panels_id, GETDATE()),

    -- Cấp 2 Categories: 8. Hệ thống Chiếu sáng
    ('150e8400-e29b-41d4-a716-446655440081', N'Đèn Pha LED', N'LED Headlight', @lighting_system_id, GETDATE()),
    ('150e8400-e29b-41d4-a716-446655440082', N'Đèn Pha LED Matrix / Pixel', N'Matrix / Pixel LED Headlight', @lighting_system_id, GETDATE()),
    ('150e8400-e29b-41d4-a716-446655440083', N'Đèn Hậu LED', N'LED Taillight', @lighting_system_id, GETDATE()),
    ('150e8400-e29b-41d4-a716-446655440084', N'Dải Đèn LED Ban ngày', N'Daytime Running Light Bar', @lighting_system_id, GETDATE()),

    -- Cấp 2 Categories: 9. Kính & Gương
    ('150e8400-e29b-41d4-a716-446655440091', N'Kính Chắn gió', N'Windshield', @glazing_mirrors_id, GETDATE()),
    ('150e8400-e29b-41d4-a716-446655440092', N'Kính Cửa', N'Door Glass', @glazing_mirrors_id, GETDATE()),
    ('150e8400-e29b-41d4-a716-446655440093', N'Gương Chiếu hậu', N'Side Mirror', @glazing_mirrors_id, GETDATE()),
    ('150e8400-e29b-41d4-a716-446655440094', N'Mái Kính Toàn cảnh', N'Panoramic Glass Roof', @glazing_mirrors_id, GETDATE()),

    -- Cấp 2 Categories: 10. Nội thất - Ghế & An toàn
    ('150e8400-e29b-41d4-a716-4466554400a1', N'Ghế Trước', N'Front Seat', @interior_seating_restraints_id, GETDATE()),
    ('150e8400-e29b-41d4-a716-4466554400a2', N'Ghế Sau', N'Rear Seat', @interior_seating_restraints_id, GETDATE()),
    ('150e8400-e29b-41d4-a716-4466554400a3', N'Túi khí', N'Airbag', @interior_seating_restraints_id, GETDATE()),
    ('150e8400-e29b-41d4-a716-4466554400a4', N'Dây đai An toàn', N'Seatbelt', @interior_seating_restraints_id, GETDATE()),

    -- Cấp 2 Categories: 11. Nội thất - Khoang lái
    ('150e8400-e29b-41d4-a716-4466554400b1', N'Bảng Đồng hồ Kỹ thuật số', N'Digital Instrument Cluster', @interior_cockpit_id, GETDATE()),
    ('150e8400-e29b-41d4-a716-4466554400b2', N'Màn hình Thông tin Giải trí Trung tâm', N'Central Infotainment Display', @interior_cockpit_id, GETDATE()),
    ('150e8400-e29b-41d4-a716-4466554400b3', N'Hiển thị Thông tin trên Kính lái (HUD - Head-Up Display)', N'Head-Up Display', @interior_cockpit_id, GETDATE()),
    ('150e8400-e29b-41d4-a716-4466554400b4', N'Bảng Điều khiển Trung tâm', N'Center Console', @interior_cockpit_id, GETDATE()),

    -- Cấp 2 Categories: 12. Điện tử & Cảm biến
    ('150e8400-e29b-41d4-a716-4466554400c1', N'Hộp Điều khiển Thân xe (BCM - Body Control Module)', N'Body Control Module', @electronics_sensors_id, GETDATE()),
    ('150e8400-e29b-41d4-a716-4466554400c2', N'Camera Hỗ trợ Lái (ADAS Camera)', N'ADAS Camera', @electronics_sensors_id, GETDATE()),
    ('150e8400-e29b-41d4-a716-4466554400c3', N'Cảm biến Radar', N'Radar Sensor', @electronics_sensors_id, GETDATE()),
    ('150e8400-e29b-41d4-a716-4466554400c4', N'Cảm biến Siêu âm', N'Ultrasonic Sensor', @electronics_sensors_id, GETDATE()),
    ('150e8400-e29b-41d4-a716-4466554400c5', N'Hộp Điều khiển Viễn thông (TCU - Telematics Control Unit)', N'Telematics Control Unit', @electronics_sensors_id, GETDATE()),

    -- Cấp 2 Categories: 13. Hệ thống Sạc
    ('150e8400-e29b-41d4-a716-4466554400d1', N'Bộ Sạc embarqué (On-Board Charger)', N'On-Board Charger', @charging_system_id, GETDATE()),
    ('150e8400-e29b-41d4-a716-4466554400d2', N'Cổng Sạc', N'Charging Port', @charging_system_id, GETDATE()),
    ('150e8400-e29b-41d4-a716-4466554400d3', N'Dây Sạc Di động', N'Portable Charging Cable', @charging_system_id, GETDATE());

-- Khai báo biến cho các Office ID để dễ quản lý
-- Ghi chú: Giả định rằng bạn có hai văn phòng khác nhau để phân phối phụ tùng.
-- phân phối nhầm về số lượng | DUNG > of1: 1 part - of2: 3 parts | nên tạm thời đổi guid
DECLARE @office_id_1 UNIQUEIDENTIFIER = '550e8400-e29b-41d4-a716-446655440002';
DECLARE @office_id_2 UNIQUEIDENTIFIER = '550e8400-e29b-41d4-a716-446655440001';


-- Khai báo ID của các danh mục để dễ đọc và quản lý
-- Cấp 2
DECLARE @pmsm_motor_cat_id UNIQUEIDENTIFIER = '150e8400-e29b-41d4-a716-446655440011';
DECLARE @im_motor_cat_id UNIQUEIDENTIFIER = '150e8400-e29b-41d4-a716-446655440012';
DECLARE @synrm_motor_cat_id UNIQUEIDENTIFIER = '150e8400-e29b-41d4-a716-446655440013';
DECLARE @inverter_cat_id UNIQUEIDENTIFIER = '150e8400-e29b-41d4-a716-446655440014';
DECLARE @gear_reducer_cat_id UNIQUEIDENTIFIER = '150e8400-e29b-41d4-a716-446655440015';
DECLARE @nmc_battery_cat_id UNIQUEIDENTIFIER = '150e8400-e29b-41d4-a716-446655440021';
DECLARE @lfp_battery_cat_id UNIQUEIDENTIFIER = '150e8400-e29b-41d4-a716-446655440022';
DECLARE @blade_battery_cat_id UNIQUEIDENTIFIER = '150e8400-e29b-41d4-a716-446655440023';
DECLARE @battery_module_cat_id UNIQUEIDENTIFIER = '150e8400-e29b-41d4-a716-446655440024';
DECLARE @bms_cat_id UNIQUEIDENTIFIER = '150e8400-e29b-41d4-a716-446655440025';
DECLARE @steel_susp_cat_id UNIQUEIDENTIFIER = '150e8400-e29b-41d4-a716-446655440031';
DECLARE @air_susp_cat_id UNIQUEIDENTIFIER = '150e8400-e29b-41d4-a716-446655440032';
DECLARE @hydraulic_damper_cat_id UNIQUEIDENTIFIER = '150e8400-e29b-41d4-a716-446655440033';
DECLARE @adaptive_damper_cat_id UNIQUEIDENTIFIER = '150e8400-e29b-41d4-a716-446655440034';
DECLARE @front_subframe_cat_id UNIQUEIDENTIFIER = '150e8400-e29b-41d4-a716-446655440035';
DECLARE @rear_subframe_cat_id UNIQUEIDENTIFIER = '150e8400-e29b-41d4-a716-446655440036';
DECLARE @eps_rack_cat_id UNIQUEIDENTIFIER = '150e8400-e29b-41d4-a716-446655440041';
DECLARE @steering_column_cat_id UNIQUEIDENTIFIER = '150e8400-e29b-41d4-a716-446655440042';
DECLARE @steering_wheel_cat_id UNIQUEIDENTIFIER = '150e8400-e29b-41d4-a716-446655440043';
DECLARE @brake_caliper_cat_id UNIQUEIDENTIFIER = '150e8400-e29b-41d4-a716-446655440051';
DECLARE @brake_disc_cat_id UNIQUEIDENTIFIER = '150e8400-e29b-41d4-a716-446655440052';
DECLARE @brake_pad_cat_id UNIQUEIDENTIFIER = '150e8400-e29b-41d4-a716-446655440053';
DECLARE @regen_braking_cat_id UNIQUEIDENTIFIER = '150e8400-e29b-41d4-a716-446655440054';
DECLARE @epb_cat_id UNIQUEIDENTIFIER = '150e8400-e29b-41d4-a716-446655440055';
DECLARE @front_door_left_cat_id UNIQUEIDENTIFIER = '150e8400-e29b-41d4-a716-446655440061';
DECLARE @front_door_right_cat_id UNIQUEIDENTIFIER = '150e8400-e29b-41d4-a716-446655440062';
DECLARE @rear_door_left_cat_id UNIQUEIDENTIFIER = '150e8400-e29b-41d4-a716-446655440063';
DECLARE @rear_door_right_cat_id UNIQUEIDENTIFIER = '150e8400-e29b-41d4-a716-446655440064';
DECLARE @liftgate_cat_id UNIQUEIDENTIFIER = '150e8400-e29b-41d4-a716-446655440065';
DECLARE @front_bumper_cat_id UNIQUEIDENTIFIER = '150e8400-e29b-41d4-a716-446655440071';
DECLARE @rear_bumper_cat_id UNIQUEIDENTIFIER = '150e8400-e29b-41d4-a716-446655440072';
DECLARE @hood_panel_cat_id UNIQUEIDENTIFIER = '150e8400-e29b-41d4-a716-446655440073';
DECLARE @front_fender_cat_id UNIQUEIDENTIFIER = '150e8400-e29b-41d4-a716-446655440074';
DECLARE @roof_panel_cat_id UNIQUEIDENTIFIER = '150e8400-e29b-41d4-a716-446655440075';
DECLARE @led_headlight_cat_id UNIQUEIDENTIFIER = '150e8400-e29b-41d4-a716-446655440081';
DECLARE @matrix_led_headlight_cat_id UNIQUEIDENTIFIER = '150e8400-e29b-41d4-a716-446655440082';
DECLARE @led_taillight_cat_id UNIQUEIDENTIFIER = '150e8400-e29b-41d4-a716-446655440083';
DECLARE @drl_bar_cat_id UNIQUEIDENTIFIER = '150e8400-e29b-41d4-a716-446655440084';
DECLARE @windshield_cat_id UNIQUEIDENTIFIER = '150e8400-e29b-41d4-a716-446655440091';
DECLARE @door_glass_cat_id UNIQUEIDENTIFIER = '150e8400-e29b-41d4-a716-446655440092';
DECLARE @side_mirror_cat_id UNIQUEIDENTIFIER = '150e8400-e29b-41d4-a716-446655440093';
DECLARE @panoramic_roof_cat_id UNIQUEIDENTIFIER = '150e8400-e29b-41d4-a716-446655440094';
DECLARE @front_seat_cat_id UNIQUEIDENTIFIER = '150e8400-e29b-41d4-a716-4466554400a1';
DECLARE @rear_seat_cat_id UNIQUEIDENTIFIER = '150e8400-e29b-41d4-a716-4466554400a2';
DECLARE @airbag_cat_id UNIQUEIDENTIFIER = '150e8400-e29b-41d4-a716-4466554400a3';
DECLARE @seatbelt_cat_id UNIQUEIDENTIFIER = '150e8400-e29b-41d4-a716-4466554400a4';
DECLARE @instrument_cluster_cat_id UNIQUEIDENTIFIER = '150e8400-e29b-41d4-a716-4466554400b1';
DECLARE @infotainment_display_cat_id UNIQUEIDENTIFIER = '150e8400-e29b-41d4-a716-4466554400b2';
DECLARE @hud_cat_id UNIQUEIDENTIFIER = '150e8400-e29b-41d4-a716-4466554400b3';
DECLARE @center_console_cat_id UNIQUEIDENTIFIER = '150e8400-e29b-41d4-a716-4466554400b4';
DECLARE @bcm_cat_id UNIQUEIDENTIFIER = '150e8400-e29b-41d4-a716-4466554400c1';
DECLARE @adas_camera_cat_id UNIQUEIDENTIFIER = '150e8400-e29b-41d4-a716-4466554400c2';
DECLARE @radar_sensor_cat_id UNIQUEIDENTIFIER = '150e8400-e29b-41d4-a716-4466554400c3';
DECLARE @ultrasonic_sensor_cat_id UNIQUEIDENTIFIER = '150e8400-e29b-41d4-a716-4466554400c4';
DECLARE @tcu_cat_id UNIQUEIDENTIFIER = '150e8400-e29b-41d4-a716-4466554400c5';
DECLARE @onboard_charger_cat_id UNIQUEIDENTIFIER = '150e8400-e29b-41d4-a716-4466554400d1';
DECLARE @charging_port_cat_id UNIQUEIDENTIFIER = '150e8400-e29b-41d4-a716-4466554400d2';
DECLARE @portable_cable_cat_id UNIQUEIDENTIFIER = '150e8400-e29b-41d4-a716-4466554400d3';


INSERT INTO parts (id, serial_number, part_name, unit_price, category_id, office_location_id, created_at)
VALUES
    -- 1. Hệ thống Truyền động Điện -> Động cơ PMSM
    --    Biến thể 1: Loại tiêu chuẩn
    (NEWID(), 'VF8-PMSM-STD-001', N'Động cơ PMSM VF8 2025 - Loại Tiêu chuẩn', 45000000, @pmsm_motor_cat_id, @office_id_1, GETDATE()),
    (NEWID(), 'VF8-PMSM-STD-002', N'Động cơ PMSM VF8 2025 - Loại Tiêu chuẩn', 45000000, @pmsm_motor_cat_id, @office_id_1, GETDATE()),
    (NEWID(), 'VF8-PMSM-STD-003', N'Động cơ PMSM VF8 2025 - Loại Tiêu chuẩn', 45000000, @pmsm_motor_cat_id, @office_id_1, GETDATE()),
    (NEWID(), 'VF8-PMSM-STD-004', N'Động cơ PMSM VF8 2025 - Loại Tiêu chuẩn', 45000000, @pmsm_motor_cat_id, @office_id_2, GETDATE()),
    --    Biến thể 2: Hiệu suất cao
    (NEWID(), 'VF8-PMSM-PERF-001', N'Động cơ PMSM VF8 2025 - Loại Hiệu suất cao', 65000000, @pmsm_motor_cat_id, @office_id_1, GETDATE()),
    (NEWID(), 'VF8-PMSM-PERF-002', N'Động cơ PMSM VF8 2025 - Loại Hiệu suất cao', 65000000, @pmsm_motor_cat_id, @office_id_1, GETDATE()),
    (NEWID(), 'VF8-PMSM-PERF-003', N'Động cơ PMSM VF8 2025 - Loại Hiệu suất cao', 65000000, @pmsm_motor_cat_id, @office_id_1, GETDATE()),
    (NEWID(), 'VF8-PMSM-PERF-004', N'Động cơ PMSM VF8 2025 - Loại Hiệu suất cao', 65000000, @pmsm_motor_cat_id, @office_id_2, GETDATE()),
    --    Biến thể 3: Nhà cung cấp B
    (NEWID(), 'VF8-PMSM-BOSCH-001', N'Động cơ PMSM VF8 2025 - Nhà cung cấp Bosch', 48000000, @pmsm_motor_cat_id, @office_id_1, GETDATE()),
    (NEWID(), 'VF8-PMSM-BOSCH-002', N'Động cơ PMSM VF8 2025 - Nhà cung cấp Bosch', 48000000, @pmsm_motor_cat_id, @office_id_1, GETDATE()),
    (NEWID(), 'VF8-PMSM-BOSCH-003', N'Động cơ PMSM VF8 2025 - Nhà cung cấp Bosch', 48000000, @pmsm_motor_cat_id, @office_id_1, GETDATE()),
    (NEWID(), 'VF8-PMSM-BOSCH-004', N'Động cơ PMSM VF8 2025 - Nhà cung cấp Bosch', 48000000, @pmsm_motor_cat_id, @office_id_2, GETDATE()),

    -- 2. Hệ thống Pin Cao áp -> Pin Lithium-ion (NMC)
    --    Biến thể 1: 92kWh
    (NEWID(), 'VF8-BAT-NMC92-001', N'Khối pin NMC VF8 2025 - 92kWh', 250000000, @nmc_battery_cat_id, @office_id_1, GETDATE()),
    (NEWID(), 'VF8-BAT-NMC92-002', N'Khối pin NMC VF8 2025 - 92kWh', 250000000, @nmc_battery_cat_id, @office_id_1, GETDATE()),
    (NEWID(), 'VF8-BAT-NMC92-003', N'Khối pin NMC VF8 2025 - 92kWh', 250000000, @nmc_battery_cat_id, @office_id_1, GETDATE()),
    (NEWID(), 'VF8-BAT-NMC92-004', N'Khối pin NMC VF8 2025 - 92kWh', 250000000, @nmc_battery_cat_id, @office_id_2, GETDATE()),
    --    Biến thể 2: 105kWh
    (NEWID(), 'VF8-BAT-NMC105-001', N'Khối pin NMC VF8 2025 - 105kWh (Tầm xa)', 320000000, @nmc_battery_cat_id, @office_id_1, GETDATE()),
    (NEWID(), 'VF8-BAT-NMC105-002', N'Khối pin NMC VF8 2025 - 105kWh (Tầm xa)', 320000000, @nmc_battery_cat_id, @office_id_1, GETDATE()),
    (NEWID(), 'VF8-BAT-NMC105-003', N'Khối pin NMC VF8 2025 - 105kWh (Tầm xa)', 320000000, @nmc_battery_cat_id, @office_id_1, GETDATE()),
    (NEWID(), 'VF8-BAT-NMC105-004', N'Khối pin NMC VF8 2025 - 105kWh (Tầm xa)', 320000000, @nmc_battery_cat_id, @office_id_2, GETDATE()),
    --    Biến thể 3: Nhà cung cấp CATL
    (NEWID(), 'VF8-BAT-CATL-001', N'Khối pin NMC VF8 2025 - Nhà cung cấp CATL', 260000000, @nmc_battery_cat_id, @office_id_1, GETDATE()),
    (NEWID(), 'VF8-BAT-CATL-002', N'Khối pin NMC VF8 2025 - Nhà cung cấp CATL', 260000000, @nmc_battery_cat_id, @office_id_1, GETDATE()),
    (NEWID(), 'VF8-BAT-CATL-003', N'Khối pin NMC VF8 2025 - Nhà cung cấp CATL', 260000000, @nmc_battery_cat_id, @office_id_1, GETDATE()),
    (NEWID(), 'VF8-BAT-CATL-004', N'Khối pin NMC VF8 2025 - Nhà cung cấp CATL', 260000000, @nmc_battery_cat_id, @office_id_2, GETDATE()),

    -- 3. Hệ thống Khung gầm & Treo -> Giảm xóc Thích ứng
    --    Biến thể 1: Loại tiêu chuẩn
    (NEWID(), 'VF8-ADAPTIVEDMP-STD-001', N'Giảm xóc Thích ứng VF8 2025 - Loại Tiêu chuẩn', 7500000, @adaptive_damper_cat_id, @office_id_1, GETDATE()),
    (NEWID(), 'VF8-ADAPTIVEDMP-STD-002', N'Giảm xóc Thích ứng VF8 2025 - Loại Tiêu chuẩn', 7500000, @adaptive_damper_cat_id, @office_id_1, GETDATE()),
    (NEWID(), 'VF8-ADAPTIVEDMP-STD-003', N'Giảm xóc Thích ứng VF8 2025 - Loại Tiêu chuẩn', 7500000, @adaptive_damper_cat_id, @office_id_1, GETDATE()),
    (NEWID(), 'VF8-ADAPTIVEDMP-STD-004', N'Giảm xóc Thích ứng VF8 2025 - Loại Tiêu chuẩn', 7500000, @adaptive_damper_cat_id, @office_id_2, GETDATE()),
    --    Biến thể 2: Loại thể thao
    (NEWID(), 'VF8-ADAPTIVEDMP-SPORT-001', N'Giảm xóc Thích ứng VF8 2025 - Loại Thể thao', 9800000, @adaptive_damper_cat_id, @office_id_1, GETDATE()),
    (NEWID(), 'VF8-ADAPTIVEDMP-SPORT-002', N'Giảm xóc Thích ứng VF8 2025 - Loại Thể thao', 9800000, @adaptive_damper_cat_id, @office_id_1, GETDATE()),
    (NEWID(), 'VF8-ADAPTIVEDMP-SPORT-003', N'Giảm xóc Thích ứng VF8 2025 - Loại Thể thao', 9800000, @adaptive_damper_cat_id, @office_id_1, GETDATE()),
    (NEWID(), 'VF8-ADAPTIVEDMP-SPORT-004', N'Giảm xóc Thích ứng VF8 2025 - Loại Thể thao', 9800000, @adaptive_damper_cat_id, @office_id_2, GETDATE()),
    --    Biến thể 3: Nhà cung cấp ZF
    (NEWID(), 'VF8-ADAPTIVEDMP-ZF-001', N'Giảm xóc Thích ứng VF8 2025 - Nhà cung cấp ZF', 8200000, @adaptive_damper_cat_id, @office_id_1, GETDATE()),
    (NEWID(), 'VF8-ADAPTIVEDMP-ZF-002', N'Giảm xóc Thích ứng VF8 2025 - Nhà cung cấp ZF', 8200000, @adaptive_damper_cat_id, @office_id_1, GETDATE()),
    (NEWID(), 'VF8-ADAPTIVEDMP-ZF-003', N'Giảm xóc Thích ứng VF8 2025 - Nhà cung cấp ZF', 8200000, @adaptive_damper_cat_id, @office_id_1, GETDATE()),
    (NEWID(), 'VF8-ADAPTIVEDMP-ZF-004', N'Giảm xóc Thích ứng VF8 2025 - Nhà cung cấp ZF', 8200000, @adaptive_damper_cat_id, @office_id_2, GETDATE()),

    -- 5. Hệ thống Phanh -> Ngàm Phanh
    --    Biến thể 1: 2 Piston
    (NEWID(), 'VF8-CALIPER-2P-001', N'Ngàm Phanh VF8 2025 - 2 Piston', 3200000, @brake_caliper_cat_id, @office_id_1, GETDATE()),
    (NEWID(), 'VF8-CALIPER-2P-002', N'Ngàm Phanh VF8 2025 - 2 Piston', 3200000, @brake_caliper_cat_id, @office_id_1, GETDATE()),
    (NEWID(), 'VF8-CALIPER-2P-003', N'Ngàm Phanh VF8 2025 - 2 Piston', 3200000, @brake_caliper_cat_id, @office_id_1, GETDATE()),
    (NEWID(), 'VF8-CALIPER-2P-004', N'Ngàm Phanh VF8 2025 - 2 Piston', 3200000, @brake_caliper_cat_id, @office_id_2, GETDATE()),
    --    Biến thể 2: 4 Piston
    (NEWID(), 'VF8-CALIPER-4P-001', N'Ngàm Phanh VF8 2025 - 4 Piston (Hiệu suất cao)', 5500000, @brake_caliper_cat_id, @office_id_1, GETDATE()),
    (NEWID(), 'VF8-CALIPER-4P-002', N'Ngàm Phanh VF8 2025 - 4 Piston (Hiệu suất cao)', 5500000, @brake_caliper_cat_id, @office_id_1, GETDATE()),
    (NEWID(), 'VF8-CALIPER-4P-003', N'Ngàm Phanh VF8 2025 - 4 Piston (Hiệu suất cao)', 5500000, @brake_caliper_cat_id, @office_id_1, GETDATE()),
    (NEWID(), 'VF8-CALIPER-4P-004', N'Ngàm Phanh VF8 2025 - 4 Piston (Hiệu suất cao)', 5500000, @brake_caliper_cat_id, @office_id_2, GETDATE()),
    --    Biến thể 3: Nhà cung cấp Brembo
    (NEWID(), 'VF8-CALIPER-BREMBO-001', N'Ngàm Phanh VF8 2025 - Nhà cung cấp Brembo', 6100000, @brake_caliper_cat_id, @office_id_1, GETDATE()),
    (NEWID(), 'VF8-CALIPER-BREMBO-002', N'Ngàm Phanh VF8 2025 - Nhà cung cấp Brembo', 6100000, @brake_caliper_cat_id, @office_id_1, GETDATE()),
    (NEWID(), 'VF8-CALIPER-BREMBO-003', N'Ngàm Phanh VF8 2025 - Nhà cung cấp Brembo', 6100000, @brake_caliper_cat_id, @office_id_1, GETDATE()),
    (NEWID(), 'VF8-CALIPER-BREMBO-004', N'Ngàm Phanh VF8 2025 - Nhà cung cấp Brembo', 6100000, @brake_caliper_cat_id, @office_id_2, GETDATE()),
    
    -- 7. Thân vỏ - Tấm ốp Ngoài -> Cản Trước
    --    Biến thể 1: Tiêu chuẩn
    (NEWID(), 'VF8-FBUMPER-STD-001', N'Cản Trước VF8 2025 - Tiêu chuẩn', 4500000, @front_bumper_cat_id, @office_id_1, GETDATE()),
    (NEWID(), 'VF8-FBUMPER-STD-002', N'Cản Trước VF8 2025 - Tiêu chuẩn', 4500000, @front_bumper_cat_id, @office_id_1, GETDATE()),
    (NEWID(), 'VF8-FBUMPER-STD-003', N'Cản Trước VF8 2025 - Tiêu chuẩn', 4500000, @front_bumper_cat_id, @office_id_1, GETDATE()),
    (NEWID(), 'VF8-FBUMPER-STD-004', N'Cản Trước VF8 2025 - Tiêu chuẩn', 4500000, @front_bumper_cat_id, @office_id_2, GETDATE()),
    --    Biến thể 2: Plus
    (NEWID(), 'VF8-FBUMPER-PLUS-001', N'Cản Trước VF8 2025 - Plus (Có cảm biến)', 6200000, @front_bumper_cat_id, @office_id_1, GETDATE()),
    (NEWID(), 'VF8-FBUMPER-PLUS-002', N'Cản Trước VF8 2025 - Plus (Có cảm biến)', 6200000, @front_bumper_cat_id, @office_id_1, GETDATE()),
    (NEWID(), 'VF8-FBUMPER-PLUS-003', N'Cản Trước VF8 2025 - Plus (Có cảm biến)', 6200000, @front_bumper_cat_id, @office_id_1, GETDATE()),
    (NEWID(), 'VF8-FBUMPER-PLUS-004', N'Cản Trước VF8 2025 - Plus (Có cảm biến)', 6200000, @front_bumper_cat_id, @office_id_2, GETDATE()),
    --    Biến thể 3: Thể thao
    (NEWID(), 'VF8-FBUMPER-SPORT-001', N'Cản Trước VF8 2025 - Gói Thể thao', 7100000, @front_bumper_cat_id, @office_id_1, GETDATE()),
    (NEWID(), 'VF8-FBUMPER-SPORT-002', N'Cản Trước VF8 2025 - Gói Thể thao', 7100000, @front_bumper_cat_id, @office_id_1, GETDATE()),
    (NEWID(), 'VF8-FBUMPER-SPORT-003', N'Cản Trước VF8 2025 - Gói Thể thao', 7100000, @front_bumper_cat_id, @office_id_1, GETDATE()),
    (NEWID(), 'VF8-FBUMPER-SPORT-004', N'Cản Trước VF8 2025 - Gói Thể thao', 7100000, @front_bumper_cat_id, @office_id_2, GETDATE()),
    
    -- 12. Điện tử & Cảm biến -> Cảm biến Radar
    --    Biến thể 1: Tầm ngắn
    (NEWID(), 'VF8-RADAR-SR-001', N'Cảm biến Radar VF8 2025 - Tầm ngắn (Góc)', 1800000, @radar_sensor_cat_id, @office_id_1, GETDATE()),
    (NEWID(), 'VF8-RADAR-SR-002', N'Cảm biến Radar VF8 2025 - Tầm ngắn (Góc)', 1800000, @radar_sensor_cat_id, @office_id_1, GETDATE()),
    (NEWID(), 'VF8-RADAR-SR-003', N'Cảm biến Radar VF8 2025 - Tầm ngắn (Góc)', 1800000, @radar_sensor_cat_id, @office_id_1, GETDATE()),
    (NEWID(), 'VF8-RADAR-SR-004', N'Cảm biến Radar VF8 2025 - Tầm ngắn (Góc)', 1800000, @radar_sensor_cat_id, @office_id_2, GETDATE()),
    --    Biến thể 2: Tầm xa
    (NEWID(), 'VF8-RADAR-LR-001', N'Cảm biến Radar VF8 2025 - Tầm xa (Phía trước)', 4500000, @radar_sensor_cat_id, @office_id_1, GETDATE()),
    (NEWID(), 'VF8-RADAR-LR-002', N'Cảm biến Radar VF8 2025 - Tầm xa (Phía trước)', 4500000, @radar_sensor_cat_id, @office_id_1, GETDATE()),
    (NEWID(), 'VF8-RADAR-LR-003', N'Cảm biến Radar VF8 2025 - Tầm xa (Phía trước)', 4500000, @radar_sensor_cat_id, @office_id_1, GETDATE()),
    (NEWID(), 'VF8-RADAR-LR-004', N'Cảm biến Radar VF8 2025 - Tầm xa (Phía trước)', 4500000, @radar_sensor_cat_id, @office_id_2, GETDATE()),
    --    Biến thể 3: Nhà cung cấp Continental
    (NEWID(), 'VF8-RADAR-CONTI-001', N'Cảm biến Radar VF8 2025 - Nhà cung cấp Continental', 4900000, @radar_sensor_cat_id, @office_id_1, GETDATE()),
    (NEWID(), 'VF8-RADAR-CONTI-002', N'Cảm biến Radar VF8 2025 - Nhà cung cấp Continental', 4900000, @radar_sensor_cat_id, @office_id_1, GETDATE()),
    (NEWID(), 'VF8-RADAR-CONTI-003', N'Cảm biến Radar VF8 2025 - Nhà cung cấp Continental', 4900000, @radar_sensor_cat_id, @office_id_1, GETDATE()),
    (NEWID(), 'VF8-RADAR-CONTI-004', N'Cảm biến Radar VF8 2025 - Nhà cung cấp Continental', 4900000, @radar_sensor_cat_id, @office_id_2, GETDATE());
    -- Khai báo Policy ID chung để sử dụng lại
DECLARE @vf8_policy_id UNIQUEIDENTIFIER = '050e8400-e29b-41d4-a716-446655440003';



INSERT INTO policy_coverage_parts (id, policy_id, part_category_id, coverage_conditions, created_at)
VALUES
    -- Cấp 1: General Systems
    ('050e8400-e29b-41d4-a716-446655440101', @vf8_policy_id, @electric_drivetrain_id, N'Lỗi vận hành gây mất công suất đột ngột, tiếng kêu lạ từ động cơ hoặc hộp số, không thể vào số (D/R).', GETDATE()),
    ('050e8400-e29b-41d4-a716-446655440102', @vf8_policy_id, @hv_battery_system_id, N'Dung lượng pin thực tế (SOH) giảm xuống dưới 70% so với dung lượng thiết kế, hoặc lỗi cách điện nghiêm trọng báo lỗi lên hệ thống.', GETDATE()),
    ('050e8400-e29b-41d4-a716-446655440103', @vf8_policy_id, @chassis_suspension_id, N'Biến dạng, nứt, gãy các thành phần khung gầm không do va chạm. Hệ thống treo mất khả năng giảm chấn, xe bị nghiêng/lệch.', GETDATE()),
    ('050e8400-e29b-41d4-a716-446655440104', @vf8_policy_id, @steering_system_id, N'Mất trợ lực lái, có độ rơ/lắc bất thường trên vô-lăng, phát ra tiếng kêu khi đánh lái.', GETDATE()),
    ('050e8400-e29b-41d4-a716-446655440105', @vf8_policy_id, @braking_system_id, N'Mất áp suất dầu phanh, lỗi hệ thống phanh tái tạo, lỗi chức năng phanh tay điện tử. Không bao gồm hao mòn tự nhiên.', GETDATE()),
    ('050e8400-e29b-41d4-a716-446655440106', @vf8_policy_id, @body_doors_id, N'Thân vỏ bị han gỉ xuyên thủng từ trong ra ngoài. Lỗi cơ cấu khóa, lên xuống kính, tay nắm cửa không hoạt động.', GETDATE()),
    ('050e8400-e29b-41d4-a716-446655440107', @vf8_policy_id, @body_exterior_panels_id, N'Lớp sơn bị bong tróc, phồng rộp, hoặc bay màu không đồng đều do lỗi sản xuất, không phải do tác động từ môi trường hoặc hóa chất.', GETDATE()),
    ('050e8400-e29b-41d4-a716-446655440108', @vf8_policy_id, @lighting_system_id, N'Đèn không sáng, ánh sáng chập chờn, hoặc bị hơi nước ngưng tụ bên trong do lỗi sản xuất.', GETDATE()),
    ('050e8400-e29b-41d4-a716-446655440109', @vf8_policy_id, @glazing_mirrors_id, N'Kính bị biến dạng hình ảnh, tự nứt vỡ không do va đập. Gương mất chức năng chỉnh điện, gập điện, sấy.', GETDATE()),
    ('050e8400-e29b-41d4-a716-446655440110', @vf8_policy_id, @interior_seating_restraints_id, N'Hệ thống túi khí hoặc dây đai an toàn báo lỗi. Cơ cấu ghế (trượt, ngả, nâng) bị kẹt hoặc không hoạt động.', GETDATE()),
    ('050e8400-e29b-41d4-a716-446655440111', @vf8_policy_id, @interior_cockpit_id, N'Màn hình có điểm chết, bị sọc, hoặc mất hiển thị. Các nút bấm vật lý bị kẹt hoặc không nhận tín hiệu.', GETDATE()),
    ('050e8400-e29b-41d4-a716-446655440112', @vf8_policy_id, @electronics_sensors_id, N'Các mô-đun điều khiển (BCM, TCU) hoặc cảm biến (Radar, Camera) không hoạt động, gây mất các tính năng thông minh và an toàn.', GETDATE()),
    ('050e8400-e29b-41d4-a716-446655440113', @vf8_policy_id, @charging_system_id, N'Xe không nhận sạc hoặc sạc không đúng công suất. Lỗi cơ cấu khóa cổng sạc.', GETDATE()),

    -- Cấp 2: Specific Components
    -- 1. Hệ thống Truyền động Điện
    ('050e8400-e29b-41d4-a716-446655440114', @vf8_policy_id, @pmsm_motor_cat_id, N'Lỗi chập cuộn dây stator, hỏng vòng bi cơ khí gây tiếng kêu bất thường, mất mô-men xoắn.', GETDATE()),
    ('050e8400-e29b-41d4-a716-446655440115', @vf8_policy_id, @im_motor_cat_id, N'Lỗi chập cuộn dây stator/rotor, hỏng vòng bi cơ khí gây tiếng kêu bất thường, mất mô-men xoắn.', GETDATE()),
    ('050e8400-e29b-41d4-a716-446655440116', @vf8_policy_id, @synrm_motor_cat_id, N'Lỗi cuộn dây stator, hỏng vòng bi cơ khí, mất khả năng đồng bộ từ trường.', GETDATE()),
    ('050e8400-e29b-41d4-a716-446655440117', @vf8_policy_id, @inverter_cat_id, N'Lỗi mạch công suất (IGBT/MOSFET), không thể chuyển đổi điện DC-AC, gây mất khả năng vận hành của động cơ.', GETDATE()),
    ('050e8400-e29b-41d4-a716-446655440118', @vf8_policy_id, @gear_reducer_cat_id, N'Mòn/vỡ bánh răng, rò rỉ dầu, hỏng vòng bi gây tiếng hú/gào khi vận hành.', GETDATE()),

    -- 2. Hệ thống Pin Cao áp
    ('050e8400-e29b-41d4-a716-446655440119', @vf8_policy_id, @nmc_battery_cat_id, N'Dung lượng thực tế (SOH) của khối pin giảm dưới 70%.', GETDATE()),
    ('050e8400-e29b-41d4-a716-446655440120', @vf8_policy_id, @lfp_battery_cat_id, N'Dung lượng thực tế (SOH) của khối pin giảm dưới 70%.', GETDATE()),
    ('050e8400-e29b-41d4-a716-446655440121', @vf8_policy_id, @blade_battery_cat_id, N'Dung lượng thực tế (SOH) của khối pin giảm dưới 70%.', GETDATE()),
    ('050e8400-e29b-41d4-a716-446655440122', @vf8_policy_id, @battery_module_cat_id, N'Lỗi một hoặc nhiều cell trong mô-đun dẫn đến mất cân bằng điện áp nghiêm trọng, không thể sạc/xả.', GETDATE()),
    ('050e8400-e29b-41d4-a716-446655440123', @vf8_policy_id, @bms_cat_id, N'Lỗi phần cứng không thể đọc/quản lý điện áp, nhiệt độ cell, hoặc không thể giao tiếp với các bộ điều khiển khác.', GETDATE()),

    -- 3. Hệ thống Khung gầm & Treo
    ('050e8400-e29b-41d4-a716-446655440124', @vf8_policy_id, @steel_susp_cat_id, N'Lò xo bị gãy, lún hoặc mất độ đàn hồi tiêu chuẩn.', GETDATE()),
    ('050e8400-e29b-41d4-a716-446655440125', @vf8_policy_id, @air_susp_cat_id, N'Rò rỉ khí ở bóng hơi, lỗi van, lỗi máy nén khí làm xe không thể nâng/hạ.', GETDATE()),
    ('050e8400-e29b-41d4-a716-446655440126', @vf8_policy_id, @hydraulic_damper_cat_id, N'Giảm xóc bị chảy dầu hoặc mất hoàn toàn khả năng dập tắt dao động.', GETDATE()),
    ('050e8400-e29b-41d4-a716-446655440127', @vf8_policy_id, @adaptive_damper_cat_id, N'Lỗi hệ thống điện tử điều khiển van biến thiên, giảm xóc không thể thay đổi độ cứng.', GETDATE()),
    ('050e8400-e29b-41d4-a716-446655440128', @vf8_policy_id, @front_subframe_cat_id, N'Nứt, gãy tại các điểm liên kết với thân xe hoặc hệ thống treo không do va đập.', GETDATE()),
    ('050e8400-e29b-41d4-a716-446655440129', @vf8_policy_id, @rear_subframe_cat_id, N'Nứt, gãy tại các điểm liên kết với thân xe hoặc hệ thống treo không do va đập.', GETDATE()),

    -- 4. Hệ thống Lái
    ('050e8400-e29b-41d4-a716-446655440130', @vf8_policy_id, @eps_rack_cat_id, N'Lỗi mô-tơ trợ lực điện hoặc cảm biến mô-men xoắn, gây nặng lái hoặc mất trợ lực.', GETDATE()),
    ('050e8400-e29b-41d4-a716-446655440131', @vf8_policy_id, @steering_column_cat_id, N'Hỏng các khớp các-đăng, trục trung gian gây độ rơ, kẹt khi xoay vô-lăng.', GETDATE()),
    ('050e8400-e29b-41d4-a716-446655440132', @vf8_policy_id, @steering_wheel_cat_id, N'Bong tróc, rộp da tự nhiên (không do hóa chất), lỗi các nút bấm chức năng.', GETDATE()),

    -- 5. Hệ thống Phanh
    ('050e8400-e29b-41d4-a716-446655440133', @vf8_policy_id, @brake_caliper_cat_id, N'Piston bị kẹt không thể hồi vị, hoặc rò rỉ dầu phanh tại ngàm phanh.', GETDATE()),
    ('050e8400-e29b-41d4-a716-446655440134', @vf8_policy_id, @brake_disc_cat_id, N'Đĩa phanh bị cong vênh, nứt vỡ do lỗi vật liệu (không bao gồm hao mòn).', GETDATE()),
    ('050e8400-e29b-41d4-a716-446655440135', @vf8_policy_id, @brake_pad_cat_id, N'Vật liệu má phanh bị bong, vỡ do lỗi sản xuất (không bao gồm hao mòn).', GETDATE()),
    ('050e8400-e29b-41d4-a716-446655440136', @vf8_policy_id, @regen_braking_cat_id, N'Hệ thống không thể tái tạo năng lượng hoặc tạo ra lực phanh không phù hợp.', GETDATE()),
    ('050e8400-e29b-41d4-a716-446655440137', @vf8_policy_id, @epb_cat_id, N'Mô-tơ hoặc cơ cấu chấp hành không thể siết/nhả phanh.', GETDATE()),

    -- 6. Thân vỏ - Cửa
    ('050e8400-e29b-41d4-a716-446655440138', @vf8_policy_id, @front_door_left_cat_id, N'Áp dụng theo điều kiện của hệ thống Thân vỏ - Cửa.', GETDATE()),
    ('050e8400-e29b-41d4-a716-446655440139', @vf8_policy_id, @front_door_right_cat_id, N'Áp dụng theo điều kiện của hệ thống Thân vỏ - Cửa.', GETDATE()),
    ('050e8400-e29b-41d4-a716-446655440140', @vf8_policy_id, @rear_door_left_cat_id, N'Áp dụng theo điều kiện của hệ thống Thân vỏ - Cửa.', GETDATE()),
    ('050e8400-e29b-41d4-a716-446655440141', @vf8_policy_id, @rear_door_right_cat_id, N'Áp dụng theo điều kiện của hệ thống Thân vỏ - Cửa.', GETDATE()),
    ('050e8400-e29b-41d4-a716-446655440142', @vf8_policy_id, @liftgate_cat_id, N'Hệ thống đóng mở điện (đá cốp, nút bấm) không hoạt động, ty chống thủy lực bị lỗi.', GETDATE()),

    -- 7. Thân vỏ - Tấm ốp Ngoài
    ('050e8400-e29b-41d4-a716-446655440143', @vf8_policy_id, @front_bumper_cat_id, N'Lỗi sơn từ nhà sản xuất, hoặc nhựa bị biến dạng tự nhiên không do nhiệt độ cao.', GETDATE()),
    ('050e8400-e29b-41d4-a716-446655440144', @vf8_policy_id, @rear_bumper_cat_id, N'Lỗi sơn từ nhà sản xuất, hoặc nhựa bị biến dạng tự nhiên không do nhiệt độ cao.', GETDATE()),
    ('050e8400-e29b-41d4-a716-446655440145', @vf8_policy_id, @hood_panel_cat_id, N'Han gỉ xuyên thủng, lỗi lớp sơn từ nhà sản xuất.', GETDATE()),
    ('050e8400-e29b-41d4-a716-446655440146', @vf8_policy_id, @front_fender_cat_id, N'Han gỉ xuyên thủng, lỗi lớp sơn từ nhà sản xuất.', GETDATE()),
    ('050e8400-e29b-41d4-a716-446655440147', @vf8_policy_id, @roof_panel_cat_id, N'Han gỉ xuyên thủng, lỗi lớp sơn từ nhà sản xuất.', GETDATE()),

    -- 8. Hệ thống Chiếu sáng
    ('050e8400-e29b-41d4-a716-446655440148', @vf8_policy_id, @led_headlight_cat_id, N'Một hoặc nhiều đi-ốt LED không sáng, lỗi bo mạch điều khiển gây chập chờn.', GETDATE()),
    ('050e8400-e29b-41d4-a716-446655440149', @vf8_policy_id, @matrix_led_headlight_cat_id, N'Lỗi chức năng liếc theo góc lái, hoặc không thể tạo vùng tối thích ứng.', GETDATE()),
    ('050e8400-e29b-41d4-a716-446655440150', @vf8_policy_id, @led_taillight_cat_id, N'Một hoặc nhiều đi-ốt LED không sáng, lỗi bo mạch điều khiển gây chập chờn.', GETDATE()),
    ('050e8400-e29b-41d4-a716-446655440151', @vf8_policy_id, @drl_bar_cat_id, N'Một hoặc nhiều đi-ốt LED không sáng, ánh sáng không đồng đều.', GETDATE()),

    -- 9. Kính & Gương
    ('050e8400-e29b-41d4-a716-446655440152', @vf8_policy_id, @windshield_cat_id, N'Lỗi tách lớp phim an toàn, biến dạng hình ảnh gây khó quan sát.', GETDATE()),
    ('050e8400-e29b-41d4-a716-446655440153', @vf8_policy_id, @door_glass_cat_id, N'Lỗi tách lớp, biến dạng hình ảnh gây khó quan sát.', GETDATE()),
    ('050e8400-e29b-41d4-a716-446655440154', @vf8_policy_id, @side_mirror_cat_id, N'Mô-tơ gập/chỉnh điện, chức năng sấy, đèn báo rẽ không hoạt động.', GETDATE()),
    ('050e8400-e29b-41d4-a716-446655440155', @vf8_policy_id, @panoramic_roof_cat_id, N'Cơ cấu trượt bị kẹt, rèm che nắng bị lỗi, gioăng cao su bị rò nước.', GETDATE()),

    -- 10. Nội thất - Ghế & An toàn
    ('050e8400-e29b-41d4-a716-446655440156', @vf8_policy_id, @front_seat_cat_id, N'Chức năng chỉnh điện, sưởi, thông gió không hoạt động.', GETDATE()),
    ('050e8400-e29b-41d4-a716-446655440157', @vf8_policy_id, @rear_seat_cat_id, N'Cơ cấu ngả, gập ghế bị kẹt, không thể khóa vào vị trí.', GETDATE()),
    ('050e8400-e29b-41d4-a716-446655440158', @vf8_policy_id, @airbag_cat_id, N'Hệ thống báo lỗi túi khí trên bảng đồng hồ, lỗi ngòi nổ hoặc hộp điều khiển.', GETDATE()),
    ('050e8400-e29b-41d4-a716-446655440159', @vf8_policy_id, @seatbelt_cat_id, N'Cơ cấu khóa không giữ chặt, bộ căng đai không hoạt động khi có tín hiệu va chạm.', GETDATE()),

    -- 11. Nội thất - Khoang lái
    ('050e8400-e29b-41d4-a716-446655440160', @vf8_policy_id, @instrument_cluster_cat_id, N'Màn hình có từ 3 điểm chết trở lên, bị sọc, hoặc mất hoàn toàn hiển thị.', GETDATE()),
    ('050e8400-e29b-41d4-a716-446655440161', @vf8_policy_id, @infotainment_display_cat_id, N'Màn hình có từ 3 điểm chết trở lên, liệt cảm ứng, treo/lag liên tục không do phần mềm.', GETDATE()),
    ('050e8400-e29b-41d4-a716-446655440162', @vf8_policy_id, @hud_cat_id, N'Mất hiển thị, hình ảnh bị mờ hoặc méo không thể đọc được.', GETDATE()),
    ('050e8400-e29b-41d4-a716-446655440163', @vf8_policy_id, @center_console_cat_id, N'Lỗi cơ cấu đóng mở hộc đồ, bệ tỳ tay, cửa gió điều hòa.', GETDATE()),

    -- 12. Điện tử & Cảm biến
    ('050e8400-e29b-41d4-a716-446655440164', @vf8_policy_id, @bcm_cat_id, N'Lỗi phần cứng gây mất điều khiển các chức năng cơ bản của xe (đèn, khóa cửa, gạt mưa...).', GETDATE()),
    ('050e8400-e29b-41d4-a716-446655440165', @vf8_policy_id, @adas_camera_cat_id, N'Lỗi phần cứng camera, không thể hiệu chỉnh (calibration), hình ảnh thu về bị mờ/sọc.', GETDATE()),
    ('050e8400-e29b-41d4-a716-446655440166', @vf8_policy_id, @radar_sensor_cat_id, N'Lỗi phần cứng cảm biến, không thể phát hoặc nhận sóng radar, gây mất chức năng ADAS.', GETDATE()),
    ('050e8400-e29b-41d4-a716-446655440167', @vf8_policy_id, @ultrasonic_sensor_cat_id, N'Cảm biến không phát hiện được vật cản hoặc báo lỗi liên tục.', GETDATE()),
    ('050e8400-e29b-41d4-a716-446655440168', @vf8_policy_id, @tcu_cat_id, N'Lỗi phần cứng không thể kết nối mạng di động (LTE/5G), mất các tính năng kết nối từ xa.', GETDATE()),

    -- 13. Hệ thống Sạc
    ('050e8400-e29b-41d4-a716-446655440169', @vf8_policy_id, @onboard_charger_cat_id, N'Lỗi phần cứng không thể chuyển đổi dòng điện AC thành DC để sạc pin.', GETDATE()),
    ('050e8400-e29b-41d4-a716-446655440170', @vf8_policy_id, @charging_port_cat_id, N'Lỗi cơ cấu khóa chốt sạc, hoặc lỗi chân tín hiệu giao tiếp với trụ sạc.', GETDATE()),
    ('050e8400-e29b-41d4-a716-446655440171', @vf8_policy_id, @portable_cable_cat_id, N'Bộ sạc báo lỗi, không thể cấp nguồn cho xe dù nguồn điện dân dụng ổn định.', GETDATE());

COMMIT TRANSACTION;

PRINT 'Seed data inserted successfully!';