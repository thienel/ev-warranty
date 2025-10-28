using Backend.Dotnet.Application.DTOs;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace Backend.Dotnet.Application.Interfaces
{
    public interface IPartService
    {
        Task<BaseResponseDto<PartDto.PartResponse>> CreateAsync(PartDto.CreatePartRequest request);

        Task<BaseResponseDto<PartDto.PartResponse>> GetByIdAsync(Guid id);
        Task<BaseResponseDto<IEnumerable<PartDto.PartResponse>>> GetAllAsync();
        Task<BaseResponseDto<PartDto.PartWithDetailsResponse>> GetWithDetailsAsync(Guid id);

        Task<BaseResponseDto<PartDto.PartResponse>> GetBySerialNumberAsync(string serialNumber);

        Task<BaseResponseDto<IEnumerable<PartDto.PartResponse>>> GetByCategoryIdAsync(Guid categoryId);
        Task<BaseResponseDto<IEnumerable<PartDto.PartResponse>>> GetAvailableByCategoryIdAsync(Guid categoryId);

        Task<BaseResponseDto<IEnumerable<PartDto.PartResponse>>> GetByOfficeLocationIdAsync(Guid officeLocationId);
        Task<BaseResponseDto<IEnumerable<PartDto.PartResponse>>> GetAvailableByOfficeLocationIdAsync(Guid officeLocationId);

        Task<BaseResponseDto<IEnumerable<PartDto.PartResponse>>> GetByStatusAsync(string status);
        Task<BaseResponseDto<IEnumerable<PartDto.PartResponse>>> GetAvailablePartsAsync();

        Task<BaseResponseDto<IEnumerable<PartDto.PartResponse>>> GetByPartNameAsync(string partName);

        Task<BaseResponseDto<IEnumerable<PartDto.PartResponse>>> SearchAsync(string searchTerm);
        Task<BaseResponseDto<IEnumerable<PartDto.PartResponse>>> SearchAvailableAsync(string searchTerm);

        Task<BaseResponseDto<PartDto.PartResponse>> UpdateAsync(Guid id, PartDto.UpdatePartRequest request);
        Task<BaseResponseDto<PartDto.PartResponse>> ChangeCategoryAsync(
            Guid id, PartDto.ChangePartCategoryRequest request);

        Task<BaseResponseDto<PartDto.PartResponse>> ChangeStatusAsync(
            Guid id, PartDto.ChangeStatusRequest request);
        Task<BaseResponseDto<PartDto.PartResponse>> ReserveAsync(Guid id);
        Task<BaseResponseDto<PartDto.PartResponse>> MarkAsInstalledAsync(Guid id);
        Task<BaseResponseDto<PartDto.PartResponse>> MarkAsDefectiveAsync(Guid id);
        Task<BaseResponseDto<PartDto.PartResponse>> MakeObsoleteAsync(Guid id);
        Task<BaseResponseDto<PartDto.PartResponse>> MakeAvailableAsync(Guid id);
        Task<BaseResponseDto<PartDto.PartResponse>> ArchiveAsync(Guid id);

        // Soft delete
        Task<BaseResponseDto> DeleteAsync(Guid id);
    }
}
