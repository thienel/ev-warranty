-- =========================================================
-- * Debug Work Orders Script
-- * Check work orders and related data for debugging
-- =========================================================

SET QUOTED_IDENTIFIER ON;
SET ANSI_NULLS ON;

PRINT 'Debugging Work Orders...';

-- Check total work orders
SELECT 'Total Work Orders' AS Info, COUNT(*) AS Count FROM work_orders;

-- Show all work orders
SELECT 
    'All Work Orders' AS Section,
    id,
    claim_id,
    assigned_technician_id,
    status,
    scheduled_date,
    completed_date,
    note,
    created_at,
    updated_at
FROM work_orders 
ORDER BY created_at DESC;

-- Check for specific ClaimId that's causing issues
DECLARE @TestClaimId UNIQUEIDENTIFIER = 'a50e8400-e29b-41d4-a716-446655440003';

SELECT 
    'Work Orders for Test ClaimId' AS Section,
    COUNT(*) AS Count
FROM work_orders 
WHERE claim_id = @TestClaimId;

SELECT 
    'Details for Test ClaimId' AS Section,
    *
FROM work_orders 
WHERE claim_id = @TestClaimId;

-- Check for any duplicate ClaimIds
SELECT 
    'Duplicate ClaimIds' AS Section,
    claim_id,
    COUNT(*) as DuplicateCount
FROM work_orders 
GROUP BY claim_id 
HAVING COUNT(*) > 1;

-- Show table structure
SELECT 
    'Work Orders Table Structure' AS Section,
    COLUMN_NAME,
    DATA_TYPE,
    IS_NULLABLE,
    COLUMN_DEFAULT
FROM INFORMATION_SCHEMA.COLUMNS 
WHERE TABLE_NAME = 'work_orders'
ORDER BY ORDINAL_POSITION;

PRINT 'Debug completed.';