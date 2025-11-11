BEGIN;

-- Seed data for Vietnamese EV warranty company

-- Insert Offices (10 offices: 1 EVM headquarters + 9 SC service centers)
INSERT INTO offices (id, office_name, office_type, address, is_active, created_at, updated_at) VALUES
-- EVM Office (Electric Vehicle Management)
('550e8400-e29b-41d4-a716-446655440001', 'VinFast Warranty Center - Head Office', 'EVM', '7th Floor, VinFast Building, 7 Bang Lao Street, Dong Da District, Hanoi', true, NOW(), NOW()),

-- SC Offices (Service Centers) across major Vietnamese cities
('550e8400-e29b-41d4-a716-446655440002', 'VinFast Service Center - Ho Chi Minh City District 1', 'SC', '123 Nguyen Hue Boulevard, District 1, Ho Chi Minh City', true, NOW(), NOW()),
('550e8400-e29b-41d4-a716-446655440003', 'VinFast Service Center - Ho Chi Minh City District 7', 'SC', '456 Nguyen Van Linh Street, District 7, Ho Chi Minh City', true, NOW(), NOW()),
('550e8400-e29b-41d4-a716-446655440004', 'VinFast Service Center - Da Nang', 'SC', '789 Vo Nguyen Giap Street, Ngu Hanh Son District, Da Nang', true, NOW(), NOW()),
('550e8400-e29b-41d4-a716-446655440005', 'VinFast Service Center - Can Tho', 'SC', '321 30/4 Street, Xuan Khanh Ward, Ninh Kieu District, Can Tho', true, NOW(), NOW()),
('550e8400-e29b-41d4-a716-446655440006', 'VinFast Service Center - Hai Phong', 'SC', '654 Dien Bien Phu Street, Le Chan District, Hai Phong', true, NOW(), NOW()),
('550e8400-e29b-41d4-a716-446655440007', 'VinFast Service Center - Nha Trang', 'SC', '987 2/4 Street, Vinh Hai Ward, Nha Trang, Khanh Hoa', true, NOW(), NOW()),
('550e8400-e29b-41d4-a716-446655440008', 'VinFast Service Center - Hue', 'SC', '159 Le Duan Street, Vinh Ninh Ward, Hue, Thua Thien Hue', true, NOW(), NOW()),
('550e8400-e29b-41d4-a716-446655440009', 'VinFast Service Center - Vung Tau', 'SC', '753 Truong Cong Dinh Street, Ward 7, Vung Tau, Ba Ria - Vung Tau', true, NOW(), NOW()),
('550e8400-e29b-41d4-a716-44665544000a', 'VinFast Service Center - Buon Ma Thuot', 'SC', '852 Phan Chu Trinh Street, Tan Loi Ward, Buon Ma Thuot, Dak Lak', true, NOW(), NOW());

-- Insert Users (30 users across different roles and offices)

INSERT INTO users (id, office_id, name, email, role, password_hash, is_active, created_at, updated_at) VALUES

-- Head Office (EVM)
('550e8400-e29b-41d4-a716-446655440101', '550e8400-e29b-41d4-a716-446655440001', 'Nguyen Van Admin', 'admin@vinfast.vn', 'ADMIN', '$2b$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.i', true, NOW(), NOW()),
('550e8400-e29b-41d4-a716-446655440102', '550e8400-e29b-41d4-a716-446655440001', 'Le Thi Minh', 'minh.le@vinfast.vn', 'EVM_STAFF', '$2a$10$tPHZ4nC5j9st/u3GRAOZQuKbHT4sdYWKaqKF30Y0dkTi7JiA8tacO', true, NOW(), NOW()),
('550e8400-e29b-41d4-a716-446655440103', '550e8400-e29b-41d4-a716-446655440001', 'Thien EVM Staff', 'lethien.dev@gmail.com', 'EVM_STAFF', '$2a$10$tPHZ4nC5j9st/u3GRAOZQuKbHT4sdYWKaqKF30Y0dkTi7JiA8tacO', true, NOW(), NOW()),

-- Ho Chi Minh City District 1 (SC)
('550e8400-e29b-41d4-a716-446655440104', '550e8400-e29b-41d4-a716-446655440002', 'Pham Thi Lan', 'lan.pham@vinfast.vn', 'SC_STAFF', '$2a$10$tPHZ4nC5j9st/u3GRAOZQuKbHT4sdYWKaqKF30Y0dkTi7JiA8tacO', true, NOW(), NOW()),
('550e8400-e29b-41d4-a716-446655440105', '550e8400-e29b-41d4-a716-446655440002', 'Do Van Hung', 'hung.do@vinfast.vn', 'SC_TECHNICIAN', '$2a$10$tPHZ4nC5j9st/u3GRAOZQuKbHT4sdYWKaqKF30Y0dkTi7JiA8tacO', true, NOW(), NOW()),
('550e8400-e29b-41d4-a716-446655440106', '550e8400-e29b-41d4-a716-446655440002', 'Vo Thi Thao', 'thao.vo@vinfast.vn', 'SC_TECHNICIAN', '$2a$10$tPHZ4nC5j9st/u3GRAOZQuKbHT4sdYWKaqKF30Y0dkTi7JiA8tacO', true, NOW(), NOW()),
('550e8400-e29b-41d4-a716-446655440120', '550e8400-e29b-41d4-a716-446655440002', 'Thien SC Staff', 'neihtel1604@gmail.com', 'SC_STAFF', '$2a$10$tPHZ4nC5j9st/u3GRAOZQuKbHT4sdYWKaqKF30Y0dkTi7JiA8tacO', true, NOW(), NOW()),
('550e8400-e29b-41d4-a716-446655440121', '550e8400-e29b-41d4-a716-446655440002', 'Thien SC Technician', 'lehungthien3@gmail.com', 'SC_TECHNICIAN', '$2a$10$tPHZ4nC5j9st/u3GRAOZQuKbHT4sdYWKaqKF30Y0dkTi7JiA8tacO', true, NOW(), NOW()),
('550e8400-e29b-41d4-a716-446655440122', '550e8400-e29b-41d4-a716-446655440002', 'Luong Van Tuan', 'tuan.luong@vinfast.vn', 'SC_TECHNICIAN', '$2a$10$tPHZ4nC5j9st/u3GRAOZQuKbHT4sdYWKaqKF30Y0dkTi7JiA8tacO', true, NOW(), NOW());

COMMIT;