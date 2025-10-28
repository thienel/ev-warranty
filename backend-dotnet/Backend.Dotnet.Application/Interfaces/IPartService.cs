using Backend.Dotnet.Application.DTOs;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using static Backend.Dotnet.Application.DTOs.PartDto;

namespace Backend.Dotnet.Application.Interfaces
{
    public interface IPartService
    {
        Task<BaseResponseDto<PartResponse>> CreateAsync(CreatePartRequest request);

        Task<BaseResponseDto<PartResponse>> GetByIdAsync(Guid id);
        Task<BaseResponseDto<IEnumerable<PartResponse>>> GetAllAsync();
        Task<BaseResponseDto<PartWithDetailsResponse>> GetWithDetailsAsync(Guid id);

        Task<BaseResponseDto<PartResponse>> GetBySerialNumberAsync(string serialNumber);
        Task<BaseResponseDto<IEnumerable<PartResponse>>> GetByCategoryIdAsync(Guid categoryId);
        Task<BaseResponseDto<IEnumerable<PartResponse>>> GetByOfficeLocationIdAsync(Guid officeLocationId);
        Task<BaseResponseDto<IEnumerable<PartResponse>>> GetByStatusAsync(string status);

        Task<BaseResponseDto<IEnumerable<PartResponse>>> SearchAsync(string searchTerm);

        Task<BaseResponseDto<PartResponse>> UpdateAsync(Guid id, UpdatePartRequest request);
        Task<BaseResponseDto<PartResponse>> ChangeCategoryAsync(
            Guid id, ChangePartCategoryRequest request);

        Task<BaseResponseDto<PartResponse>> ChangeStatusAsync(
            Guid id, PartChangeStatusRequest request);

        // Soft delete
        Task<BaseResponseDto> DeleteAsync(Guid id);
    }
}
