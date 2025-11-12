-- =========================================================
-- * EV Warranty Database Clean-up Script
-- * This script removes all data in the correct order to avoid FK constraint violations
-- * Safe to run multiple times
-- =========================================================

-- Set required SQL Server options
SET QUOTED_IDENTIFIER ON;
SET ANSI_NULLS ON;
SET ANSI_PADDING ON;
SET ANSI_WARNINGS ON;
SET CONCAT_NULL_YIELDS_NULL ON;
SET ARITHABORT ON;

BEGIN TRANSACTION;

PRINT 'Starting database cleanup...';

-- =========================================================
-- * Step 1: Delete data with foreign key dependencies first
-- =========================================================

-- Delete work orders (references no other tables as FK)
IF EXISTS (SELECT 1 FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_NAME = 'work_orders')
BEGIN
    DELETE FROM work_orders;
    PRINT 'Deleted all work orders';
END

-- Delete policy coverage parts (references warranty_policies and part_categories)
IF EXISTS (SELECT 1 FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_NAME = 'policy_coverage_parts')
BEGIN
    DELETE FROM policy_coverage_parts;
    PRINT 'Deleted all policy coverage parts';
END

-- Delete parts (references part_categories and office locations)
IF EXISTS (SELECT 1 FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_NAME = 'parts')
BEGIN
    DELETE FROM parts;
    PRINT 'Deleted all parts';
END

-- Delete vehicles (references customers and vehicle_models)
IF EXISTS (SELECT 1 FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_NAME = 'vehicles')
BEGIN
    DELETE FROM vehicles;
    PRINT 'Deleted all vehicles';
END

-- =========================================================
-- * Step 2: Delete parent/lookup tables
-- =========================================================

-- Delete part categories (has self-referencing FK and references from parts)
IF EXISTS (SELECT 1 FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_NAME = 'part_categories')
BEGIN
    -- Delete child categories first (those with parent_category_id)
    DELETE FROM part_categories WHERE parent_category_id IS NOT NULL;
    -- Then delete parent categories
    DELETE FROM part_categories WHERE parent_category_id IS NULL;
    PRINT 'Deleted all part categories';
END

-- Delete vehicle models (references warranty_policies)
IF EXISTS (SELECT 1 FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_NAME = 'vehicle_models')
BEGIN
    DELETE FROM vehicle_models;
    PRINT 'Deleted all vehicle models';
END

-- Delete customers (referenced by vehicles)
IF EXISTS (SELECT 1 FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_NAME = 'customers')
BEGIN
    DELETE FROM customers;
    PRINT 'Deleted all customers';
END

-- Delete warranty policies (referenced by vehicle_models and policy_coverage_parts)
IF EXISTS (SELECT 1 FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_NAME = 'warranty_policies')
BEGIN
    DELETE FROM warranty_policies;
    PRINT 'Deleted all warranty policies';
END

-- =========================================================
-- * Step 3: Reset identity columns (if any)
-- =========================================================

-- Reset identity seeds for tables that might have IDENTITY columns
-- Note: Our current schema uses GUIDs, but this is for future-proofing

IF EXISTS (SELECT 1 FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_NAME = 'work_orders')
    DBCC CHECKIDENT ('work_orders', RESEED, 0);

IF EXISTS (SELECT 1 FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_NAME = 'vehicles')
    DBCC CHECKIDENT ('vehicles', RESEED, 0);

IF EXISTS (SELECT 1 FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_NAME = 'customers')
    DBCC CHECKIDENT ('customers', RESEED, 0);

IF EXISTS (SELECT 1 FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_NAME = 'parts')
    DBCC CHECKIDENT ('parts', RESEED, 0);

IF EXISTS (SELECT 1 FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_NAME = 'part_categories')
    DBCC CHECKIDENT ('part_categories', RESEED, 0);

IF EXISTS (SELECT 1 FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_NAME = 'vehicle_models')
    DBCC CHECKIDENT ('vehicle_models', RESEED, 0);

IF EXISTS (SELECT 1 FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_NAME = 'warranty_policies')
    DBCC CHECKIDENT ('warranty_policies', RESEED, 0);

IF EXISTS (SELECT 1 FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_NAME = 'policy_coverage_parts')
    DBCC CHECKIDENT ('policy_coverage_parts', RESEED, 0);

-- =========================================================
-- * Step 4: Verify cleanup
-- =========================================================

DECLARE @TotalRecords INT = 0;

SELECT @TotalRecords = (
    ISNULL((SELECT COUNT(*) FROM work_orders), 0) +
    ISNULL((SELECT COUNT(*) FROM policy_coverage_parts), 0) +
    ISNULL((SELECT COUNT(*) FROM parts), 0) +
    ISNULL((SELECT COUNT(*) FROM vehicles), 0) +
    ISNULL((SELECT COUNT(*) FROM part_categories), 0) +
    ISNULL((SELECT COUNT(*) FROM vehicle_models), 0) +
    ISNULL((SELECT COUNT(*) FROM customers), 0) +
    ISNULL((SELECT COUNT(*) FROM warranty_policies), 0)
);

IF @TotalRecords = 0
BEGIN
    PRINT 'SUCCESS: Database cleanup completed successfully. All tables are empty.';
END
ELSE
BEGIN
    PRINT 'WARNING: Some records may still exist. Total remaining records: ' + CAST(@TotalRecords AS VARCHAR(10));
END

COMMIT TRANSACTION;

PRINT 'Database cleanup script completed.';

