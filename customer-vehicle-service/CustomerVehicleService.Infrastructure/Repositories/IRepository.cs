using System;
using System.Collections.Generic;
using System.Linq;
using System.Linq.Expressions;
using System.Text;
using System.Threading.Tasks;

namespace CustomerVehicleService.Infrastructure.Repositories
{
    public interface IRepository<TEntity> where TEntity : class
    {
        // Basic query operations
        Task<TEntity> GetByIdAsync(Guid id);
        Task<IEnumerable<TEntity>> GetAllAsync();
        Task<IEnumerable<TEntity>> FindAsync(Expression<Func<TEntity, bool>> predicate);
        Task<TEntity> FirstOrDefaultAsync(Expression<Func<TEntity, bool>> predicate);
        Task<bool> ExistsAsync(Expression<Func<TEntity, bool>> predicate);

        // Query access for flexibility
        IQueryable<TEntity> Query(); // AI them chu chua biet lam gi =)?

        // Command operations
        Task AddAsync(TEntity entity);
        Task AddRangeAsync(IEnumerable<TEntity> entities);
        void Update(TEntity entity);
        void UpdateRange(IEnumerable<TEntity> entities); //???
        void Remove(TEntity entity);
        void RemoveRange(IEnumerable<TEntity> entities); //???
    }
}
