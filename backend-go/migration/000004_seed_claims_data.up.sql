BEGIN;

-- Sử dụng offices và users đã có từ migration 3
-- Office IDs từ migration 3:
-- Head Office (EVM): '550e8400-e29b-41d4-a716-446655440001' 
-- HCM Service Center: '550e8400-e29b-41d4-a716-446655440002'
-- HCM District 7 Service Center: '550e8400-e29b-41d4-a716-446655440003'
-- Da Nang Service Center: '550e8400-e29b-41d4-a716-446655440004'

-- User IDs từ migration 3:
-- EVM Staff: '550e8400-e29b-41d4-a716-446655440102' (Le Thi Minh), '550e8400-e29b-41d4-a716-446655440103' (Thien EVM Staff)
-- SC Staff: '550e8400-e29b-41d4-a716-446655440104' (Pham Thi Lan), '550e8400-e29b-41d4-a716-446655440120' (Thien SC Staff)
-- SC Technicians: '550e8400-e29b-41d4-a716-446655440105' (Do Van Hung), '550e8400-e29b-41d4-a716-446655440106' (Vo Thi Thao), '550e8400-e29b-41d4-a716-446655440121' (Thien SC Technician), '550e8400-e29b-41d4-a716-446655440122' (Luong Van Tuan)

-- Insert Claims with different statuses
INSERT INTO claims (id, customer_id, vehicle_id, kilometers, description, status, staff_id, technician_id, approved_by, total_cost, created_at, updated_at)
VALUES
    -- 1. DRAFT Status - Customer Minh's first vehicle  
    ('a50e8400-e29b-41d4-a716-446655440001', '850e8400-e29b-41d4-a716-446655440001', '750e8400-e29b-41d4-a716-446655440001', 15000, 
     'Xe báo lỗi hệ thống pin, đèn cảnh báo pin sáng liên tục trên bảng đồng hồ. Khách hàng phản ánh xe giảm tầm di chuyển đáng kể.',
     'DRAFT', '550e8400-e29b-41d4-a716-446655440104', '550e8400-e29b-41d4-a716-446655440105', NULL, 0, '2024-10-15 08:30:00', '2024-10-15 08:30:00'),
    
    -- 2. SUBMITTED Status - Customer Lan's vehicle  
    ('a50e8400-e29b-41d4-a716-446655440002', '850e8400-e29b-41d4-a716-446655440002', '750e8400-e29b-41d4-a716-446655440003', 23500,
     'Động cơ phát ra tiếng kêu bất thường khi tăng tốc, đặc biệt ở vận tốc cao. Mất mô-men xoắn khi leo dốc.',
     'SUBMITTED', '550e8400-e29b-41d4-a716-446655440120', '550e8400-e29b-41d4-a716-446655440106', NULL, 0, '2024-10-20 14:15:00', '2024-10-22 09:45:00'),
    
    -- 3. REVIEWING Status - Customer Hùng's vehicle
    ('a50e8400-e29b-41d4-a716-446655440003', '850e8400-e29b-41d4-a716-446655440003', '750e8400-e29b-41d4-a716-446655440004', 18750,
     'Hệ thống phanh hoạt động bất thường, cảm giác phanh bị mềm và khoảng cách phanh tăng. Nghi ngờ lỗi ngàm phanh.',
     'REVIEWING', '550e8400-e29b-41d4-a716-446655440104', '550e8400-e29b-41d4-a716-446655440105', NULL, 0, '2024-10-25 11:20:00', '2024-10-28 16:30:00'),
    
    -- 4. APPROVED Status - Customer Hoa's vehicle (Technician Do Van Hung has max 3 claims)
    ('a50e8400-e29b-41d4-a716-446655440004', '850e8400-e29b-41d4-a716-446655440004', '750e8400-e29b-41d4-a716-446655440005', 12300,
     'Giảm xóc thích ứng không hoạt động, xe bị xóc mạnh khi đi qua ổ gà. Hệ thống không thể điều chỉnh độ cứng.',
     'APPROVED', '550e8400-e29b-41d4-a716-446655440102', '550e8400-e29b-41d4-a716-446655440121', '550e8400-e29b-41d4-a716-446655440101', 8200000, '2024-11-01 09:10:00', '2024-11-05 15:45:00'),
    
    -- 5. PARTIALLY_APPROVED Status - Customer Tuấn's vehicle  
    ('a50e8400-e29b-41d4-a716-446655440005', '850e8400-e29b-41d4-a716-446655440005', '750e8400-e29b-41d4-a716-446655440006', 28900,
     'Cản trước bị biến dạng và đèn pha LED bị hỏng một phần. Khách hàng cho biết không có va chạm.',
     'PARTIALLY_APPROVED', '550e8400-e29b-41d4-a716-446655440120', '550e8400-e29b-41d4-a716-446655440106', '550e8400-e29b-41d4-a716-446655440102', 4500000, '2024-11-03 13:25:00', '2024-11-07 10:20:00'),
    
    -- 6. REJECTED Status - Customer Linh's vehicle (Technician Vo Thi Thao has max 3 claims)
    ('a50e8400-e29b-41d4-a716-446655440006', '850e8400-e29b-41d4-a716-446655440006', '750e8400-e29b-41d4-a716-446655440007', 45600,
     'Cảm biến radar bị lỗi, hệ thống ADAS không hoạt động. Các tính năng an toàn chủ động bị tắt.',
     'REJECTED', '550e8400-e29b-41d4-a716-446655440104', '550e8400-e29b-41d4-a716-446655440121', '550e8400-e29b-41d4-a716-446655440101', 0, '2024-11-08 10:45:00', '2024-11-10 14:15:00'),
    
    -- 7. CANCELLED Status - Customer Minh's second vehicle (Technician Luong Van Tuan has max 3 claims)
    ('a50e8400-e29b-41d4-a716-446655440007', '850e8400-e29b-41d4-a716-446655440001', '750e8400-e29b-41d4-a716-446655440002', 8950,
     'Khách hàng hủy yêu cầu bảo hành do đã tự khắc phục được sự cố.',
     'CANCELLED', '550e8400-e29b-41d4-a716-446655440102', '550e8400-e29b-41d4-a716-446655440122', NULL, 0, '2024-11-09 16:00:00', '2024-11-09 18:30:00');

-- Insert Claim Items (ensuring MinItemPerClaim = 1)
INSERT INTO claim_items (id, claim_id, part_category_id, faulty_part_id, replacement_part_id, issue_description, status, type, cost, created_at, updated_at)
VALUES
    -- Claim 1 (DRAFT) - Battery Issue
    ('b50e8400-e29b-41d4-a716-446655440001', 'a50e8400-e29b-41d4-a716-446655440001', '150e8400-e29b-41d4-a716-446655440021', 
     'c50e8400-e29b-41d4-a716-446655440001', 'c50e8400-e29b-41d4-a716-446655440002', 
     'Khối pin NMC hiển thị SOH giảm xuống 65%, dưới ngưỡng bảo hành 70%', 'PENDING', 'REPLACEMENT', 250000000, NOW(), NOW()),
    
    -- Claim 2 (SUBMITTED) - Motor Issue  
    ('b50e8400-e29b-41d4-a716-446655440002', 'a50e8400-e29b-41d4-a716-446655440002', '150e8400-e29b-41d4-a716-446655440011',
     'c50e8400-e29b-41d4-a716-446655440011', 'c50e8400-e29b-41d4-a716-446655440012',
     'Động cơ PMSM phát tiếng kêu vòng bi, giảm hiệu suất', 'PENDING', 'REPLACEMENT', 45000000, NOW(), NOW()),
    
    -- Claim 3 (REVIEWING) - Brake Caliper Issue
    ('b50e8400-e29b-41d4-a716-446655440003', 'a50e8400-e29b-41d4-a716-446655440003', '150e8400-e29b-41d4-a716-446655440051',
     'c50e8400-e29b-41d4-a716-446655440021', 'c50e8400-e29b-41d4-a716-446655440022',
     'Ngàm phanh bị rò rỉ dầu, piston kẹt không hồi vị', 'PENDING', 'REPLACEMENT', 3200000, NOW(), NOW()),
    
    -- Claim 4 (APPROVED) - Adaptive Damper Issue  
    ('b50e8400-e29b-41d4-a716-446655440004', 'a50e8400-e29b-41d4-a716-446655440004', '150e8400-e29b-41d4-a716-446655440034',
     'c50e8400-e29b-41d4-a716-446655440031', 'c50e8400-e29b-41d4-a716-446655440032',
     'Giảm xóc thích ứng mất chức năng điện tử, van biến thiên lỗi', 'APPROVED', 'REPLACEMENT', 8200000, NOW(), NOW()),
    
    -- Claim 5 (PARTIALLY_APPROVED) - Front Bumper Issue
    ('b50e8400-e29b-41d4-a716-446655440005', 'a50e8400-e29b-41d4-a716-446655440005', '150e8400-e29b-41d4-a716-446655440071',
     'c50e8400-e29b-41d4-a716-446655440041', 'c50e8400-e29b-41d4-a716-446655440042',
     'Cản trước biến dạng không do va chạm, nghi ngờ lỗi vật liệu', 'APPROVED', 'REPLACEMENT', 4500000, NOW(), NOW()),
    -- Additional item for claim 5 - LED Headlight (REJECTED)
    ('b50e8400-e29b-41d4-a716-446655440006', 'a50e8400-e29b-41d4-a716-446655440005', '150e8400-e29b-41d4-a716-446655440081',
     'c50e8400-e29b-41d4-a716-446655440051', NULL,
     'Đèn pha LED một số điốt không sáng', 'REJECTED', 'REPLACEMENT', 0, NOW(), NOW()),
    
    -- Claim 6 (REJECTED) - Radar Sensor Issue
    ('b50e8400-e29b-41d4-a716-446655440007', 'a50e8400-e29b-41d4-a716-446655440006', '150e8400-e29b-41d4-a716-4466554400c3',
     'c50e8400-e29b-41d4-a716-446655440061', NULL,
     'Cảm biến radar không phát hiện vật cản, khách hàng đã qua km bảo hành', 'REJECTED', 'REPLACEMENT', 0, NOW(), NOW()),
    
    -- Claim 7 (CANCELLED) - No items needed since cancelled
    ('b50e8400-e29b-41d4-a716-446655440008', 'a50e8400-e29b-41d4-a716-446655440007', '150e8400-e29b-41d4-a716-446655440011',
     'c50e8400-e29b-41d4-a716-446655440071', NULL,
     'Khách hàng hủy yêu cầu', 'PENDING', 'REPLACEMENT', 0, NOW(), NOW());

-- Insert Claim Attachments (ensuring MinAttachmentPerClaim = 2)
INSERT INTO claim_attachments (id, claim_id, type, url, created_at)
VALUES
    -- Claim 1 attachments
    ('d50e8400-e29b-41d4-a716-446655440001', 'a50e8400-e29b-41d4-a716-446655440001', 'IMAGE', 
     'https://storage.vinfast.vn/claims/minh_vf8_1_battery_error_dashboard.jpg', NOW()),
    ('d50e8400-e29b-41d4-a716-446655440002', 'a50e8400-e29b-41d4-a716-446655440001', 'VIDEO', 
     'https://storage.vinfast.vn/claims/minh_vf8_1_battery_error_test.mp4', NOW()),
    ('d50e8400-e29b-41d4-a716-446655440003', 'a50e8400-e29b-41d4-a716-446655440001', 'DOCUMENT', 
     'https://storage.vinfast.vn/claims/minh_vf8_1_battery_diagnostic_report.pdf', NOW()),
    
    -- Claim 2 attachments
    ('d50e8400-e29b-41d4-a716-446655440004', 'a50e8400-e29b-41d4-a716-446655440002', 'IMAGE', 
     'https://storage.vinfast.vn/claims/lan_vf8_motor_noise_photo.jpg', NOW()),
    ('d50e8400-e29b-41d4-a716-446655440005', 'a50e8400-e29b-41d4-a716-446655440002', 'AUDIO', 
     'https://storage.vinfast.vn/claims/lan_vf8_motor_noise_recording.mp3', NOW()),
    
    -- Claim 3 attachments  
    ('d50e8400-e29b-41d4-a716-446655440006', 'a50e8400-e29b-41d4-a716-446655440003', 'IMAGE', 
     'https://storage.vinfast.vn/claims/hung_vf8_brake_fluid_leak.jpg', NOW()),
    ('d50e8400-e29b-41d4-a716-446655440007', 'a50e8400-e29b-41d4-a716-446655440003', 'DOCUMENT', 
     'https://storage.vinfast.vn/claims/hung_vf8_brake_test_report.pdf', NOW()),
    
    -- Claim 4 attachments
    ('d50e8400-e29b-41d4-a716-446655440008', 'a50e8400-e29b-41d4-a716-446655440004', 'IMAGE', 
     'https://storage.vinfast.vn/claims/hoa_vf8_damper_malfunction.jpg', NOW()),
    ('d50e8400-e29b-41d4-a716-446655440009', 'a50e8400-e29b-41d4-a716-446655440004', 'VIDEO', 
     'https://storage.vinfast.vn/claims/hoa_vf8_damper_test_drive.mp4', NOW()),
    
    -- Claim 5 attachments
    ('d50e8400-e29b-41d4-a716-446655440010', 'a50e8400-e29b-41d4-a716-446655440005', 'IMAGE', 
     'https://storage.vinfast.vn/claims/tuan_vf8_bumper_deformation.jpg', NOW()),
    ('d50e8400-e29b-41d4-a716-446655440011', 'a50e8400-e29b-41d4-a716-446655440005', 'IMAGE', 
     'https://storage.vinfast.vn/claims/tuan_vf8_led_headlight_failure.jpg', NOW()),
    
    -- Claim 6 attachments
    ('d50e8400-e29b-41d4-a716-446655440012', 'a50e8400-e29b-41d4-a716-446655440006', 'DOCUMENT', 
     'https://storage.vinfast.vn/claims/linh_vf8_radar_diagnostic.pdf', NOW()),
    ('d50e8400-e29b-41d4-a716-446655440013', 'a50e8400-e29b-41d4-a716-446655440006', 'IMAGE', 
     'https://storage.vinfast.vn/claims/linh_vf8_adas_error_screen.jpg', NOW()),
    
    -- Claim 7 attachments
    ('d50e8400-e29b-41d4-a716-446655440014', 'a50e8400-e29b-41d4-a716-446655440007', 'IMAGE', 
     'https://storage.vinfast.vn/claims/minh_vf8_2_initial_report.jpg', NOW()),
    ('d50e8400-e29b-41d4-a716-446655440015', 'a50e8400-e29b-41d4-a716-446655440007', 'DOCUMENT', 
     'https://storage.vinfast.vn/claims/minh_vf8_2_cancellation_request.pdf', NOW());

-- Insert Claim Histories to track status changes
INSERT INTO claim_histories (id, claim_id, status, changed_by, changed_at)
VALUES
    -- Claim 1 history (DRAFT)
    ('e50e8400-e29b-41d4-a716-446655440001', 'a50e8400-e29b-41d4-a716-446655440001', 'DRAFT', '550e8400-e29b-41d4-a716-446655440104', '2024-10-15 08:30:00'),
    
    -- Claim 2 history (DRAFT -> SUBMITTED)  
    ('e50e8400-e29b-41d4-a716-446655440002', 'a50e8400-e29b-41d4-a716-446655440002', 'DRAFT', '550e8400-e29b-41d4-a716-446655440120', '2024-10-20 14:15:00'),
    ('e50e8400-e29b-41d4-a716-446655440003', 'a50e8400-e29b-41d4-a716-446655440002', 'SUBMITTED', '550e8400-e29b-41d4-a716-446655440120', '2024-10-22 09:45:00'),
    
    -- Claim 3 history (DRAFT -> SUBMITTED -> REVIEWING)
    ('e50e8400-e29b-41d4-a716-446655440004', 'a50e8400-e29b-41d4-a716-446655440003', 'DRAFT', '550e8400-e29b-41d4-a716-446655440104', '2024-10-25 11:20:00'),
    ('e50e8400-e29b-41d4-a716-446655440005', 'a50e8400-e29b-41d4-a716-446655440003', 'SUBMITTED', '550e8400-e29b-41d4-a716-446655440104', '2024-10-26 15:30:00'),
    ('e50e8400-e29b-41d4-a716-446655440006', 'a50e8400-e29b-41d4-a716-446655440003', 'REVIEWING', '550e8400-e29b-41d4-a716-446655440101', '2024-10-28 16:30:00'),
    
    -- Claim 4 history (DRAFT -> SUBMITTED -> REVIEWING -> APPROVED)
    ('e50e8400-e29b-41d4-a716-446655440007', 'a50e8400-e29b-41d4-a716-446655440004', 'DRAFT', '550e8400-e29b-41d4-a716-446655440102', '2024-11-01 09:10:00'),
    ('e50e8400-e29b-41d4-a716-446655440008', 'a50e8400-e29b-41d4-a716-446655440004', 'SUBMITTED', '550e8400-e29b-41d4-a716-446655440102', '2024-11-02 14:20:00'),
    ('e50e8400-e29b-41d4-a716-446655440009', 'a50e8400-e29b-41d4-a716-446655440004', 'REVIEWING', '550e8400-e29b-41d4-a716-446655440101', '2024-11-04 10:45:00'),
    ('e50e8400-e29b-41d4-a716-446655440010', 'a50e8400-e29b-41d4-a716-446655440004', 'APPROVED', '550e8400-e29b-41d4-a716-446655440101', '2024-11-05 15:45:00'),
    
    -- Claim 5 history (DRAFT -> SUBMITTED -> REVIEWING -> PARTIALLY_APPROVED)
    ('e50e8400-e29b-41d4-a716-446655440011', 'a50e8400-e29b-41d4-a716-446655440005', 'DRAFT', '550e8400-e29b-41d4-a716-446655440120', '2024-11-03 13:25:00'),
    ('e50e8400-e29b-41d4-a716-446655440012', 'a50e8400-e29b-41d4-a716-446655440005', 'SUBMITTED', '550e8400-e29b-41d4-a716-446655440120', '2024-11-04 16:15:00'),
    ('e50e8400-e29b-41d4-a716-446655440013', 'a50e8400-e29b-41d4-a716-446655440005', 'REVIEWING', '550e8400-e29b-41d4-a716-446655440102', '2024-11-06 11:30:00'),
    ('e50e8400-e29b-41d4-a716-446655440014', 'a50e8400-e29b-41d4-a716-446655440005', 'PARTIALLY_APPROVED', '550e8400-e29b-41d4-a716-446655440102', '2024-11-07 10:20:00'),
    
    -- Claim 6 history (DRAFT -> SUBMITTED -> REVIEWING -> REJECTED)
    ('e50e8400-e29b-41d4-a716-446655440015', 'a50e8400-e29b-41d4-a716-446655440006', 'DRAFT', '550e8400-e29b-41d4-a716-446655440104', '2024-11-08 10:45:00'),
    ('e50e8400-e29b-41d4-a716-446655440016', 'a50e8400-e29b-41d4-a716-446655440006', 'SUBMITTED', '550e8400-e29b-41d4-a716-446655440104', '2024-11-08 14:20:00'),
    ('e50e8400-e29b-41d4-a716-446655440017', 'a50e8400-e29b-41d4-a716-446655440006', 'REVIEWING', '550e8400-e29b-41d4-a716-446655440101', '2024-11-09 09:15:00'),
    ('e50e8400-e29b-41d4-a716-446655440018', 'a50e8400-e29b-41d4-a716-446655440006', 'REJECTED', '550e8400-e29b-41d4-a716-446655440101', '2024-11-10 14:15:00'),
    
    -- Claim 7 history (DRAFT -> CANCELLED) 
    ('e50e8400-e29b-41d4-a716-446655440019', 'a50e8400-e29b-41d4-a716-446655440007', 'DRAFT', '550e8400-e29b-41d4-a716-446655440102', '2024-11-09 16:00:00'),
    ('e50e8400-e29b-41d4-a716-446655440020', 'a50e8400-e29b-41d4-a716-446655440007', 'CANCELLED', '550e8400-e29b-41d4-a716-446655440102', '2024-11-09 18:30:00');

COMMIT;