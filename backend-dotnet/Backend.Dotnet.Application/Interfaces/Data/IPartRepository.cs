using Backend.Dotnet.Domain.Entities;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace Backend.Dotnet.Application.Interfaces.Data
{
    public interface IPartRepository : IRepository<Part>
    {
        Task<Part?> GetBySerialNumberAsync(string serialNumber);
        Task<bool> SerialNumberExistsAsync(string serialNumber, Guid? excludePartId = null);

        Task<IEnumerable<Part>> GetByCategoryIdAsync(Guid categoryId, PartStatus? status = null);
        Task<IEnumerable<Part>> GetByOfficeLocationIdAsync(Guid officeLocationId, PartStatus? status = null);

        Task<IEnumerable<Part>> GetByStatusAsync(PartStatus status);

        Task<Part?> GetWithDetailsAsync(Guid partId);
        Task<Part?> GetByOfficeIdAndCategoryId(Guid officeId, Guid categoryId);
        Task<IEnumerable<Part>> GetAllWithDetailsAsync();
        Task<IEnumerable<Part>> SearchAsync(string searchTerm);

        Task<int> GetAvailableCountByCategoryAsync(Guid categoryId);
        Task<int> GetAvailableCountByLocationAsync(Guid officeLocationId);

        //Task<decimal> GetAveragePriceByCategoryAsync(Guid categoryId); // Further
    }
}
