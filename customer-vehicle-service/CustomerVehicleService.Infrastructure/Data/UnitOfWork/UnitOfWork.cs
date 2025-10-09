using CustomerVehicleService.Application.Interfaces.Data;
using CustomerVehicleService.Infrastructure.Data.Context;
using CustomerVehicleService.Infrastructure.Data.Repositories;
using Microsoft.EntityFrameworkCore;
using Microsoft.EntityFrameworkCore.Storage;
using Microsoft.Extensions.Logging;

namespace CustomerVehicleService.Infrastructure.Data.UnitOfWork
{
    public class UnitOfWork : IUnitOfWork
    {
        private readonly AppDbContext _context;
        private readonly ILogger<UnitOfWork>? _logger;
        private IDbContextTransaction? _transaction;
        private bool _disposed;

        // Lazy-loaded repositories
        private ICustomerRepository? _customers;
        private IVehicleRepository? _vehicles;
        private IVehicleModelRepository? _vehicleModels;

        public UnitOfWork(
            AppDbContext context,
            ILogger<UnitOfWork>? logger = null)
        {
            _context = context ?? throw new ArgumentNullException(nameof(context));
            _logger = logger;
        }

        public ICustomerRepository Customers
        {
            get
            {
                _customers ??= new CustomerRepository(_context);
                return _customers;
            }
        }

        public IVehicleRepository Vehicles
        {
            get
            {
                _vehicles ??= new VehicleRepository(_context);
                return _vehicles;
            }
        }

        public IVehicleModelRepository VehicleModels
        {
            get
            {
                _vehicleModels ??= new VehicleModelRepository(_context);
                return _vehicleModels;
            }
        }

        public async Task<int> SaveChangesAsync()
        {
            try
            {
                return await _context.SaveChangesAsync();
            }
            catch (DbUpdateConcurrencyException ex)
            {
                _logger?.LogError(ex, "Concurrency conflict while saving changes");
                throw;
            }
            catch (DbUpdateException ex)
            {
                _logger?.LogError(ex, "Database error while saving changes");
                throw;
            }
        }

        public int SaveChanges()
        {
            try
            {
                return _context.SaveChanges();
            }
            catch (DbUpdateConcurrencyException ex)
            {
                _logger?.LogError(ex, "Concurrency conflict while saving changes");
                throw;
            }
            catch (DbUpdateException ex)
            {
                _logger?.LogError(ex, "Database error while saving changes");
                throw;
            }
        }

        public async Task BeginTransactionAsync()
        {
            if (_transaction != null)
            {
                throw new InvalidOperationException("Transaction already started");
            }

            _transaction = await _context.Database.BeginTransactionAsync();
            _logger?.LogDebug("Transaction started");
        }

        public async Task CommitTransactionAsync()
        {
            if (_transaction == null)
            {
                throw new InvalidOperationException("No active transaction");
            }

            try
            {
                await _transaction.CommitAsync();
                _logger?.LogDebug("Transaction committed");
            }
            catch (Exception ex)
            {
                _logger?.LogError(ex, "Error committing transaction");
                await RollbackTransactionAsync();
                throw;
            }
            finally
            {
                await DisposeTransactionAsync();
            }
        }

        public async Task RollbackTransactionAsync()
        {
            if (_transaction == null)
            {
                return;
            }

            try
            {
                await _transaction.RollbackAsync();
                _logger?.LogDebug("Transaction rolled back");
            }
            catch (Exception ex)
            {
                _logger?.LogError(ex, "Error rolling back transaction");
                throw;
            }
            finally
            {
                await DisposeTransactionAsync();
            }
        }

        private async Task DisposeTransactionAsync()
        {
            if (_transaction != null)
            {
                await _transaction.DisposeAsync();
                _transaction = null;
            }
        }

        public void Dispose()
        {
            if (_disposed)
            {
                return;
            }

            _transaction?.Dispose();
            _transaction = null;
            _disposed = true;

            GC.SuppressFinalize(this);
        }
    }
}
