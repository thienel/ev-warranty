using Backend.Dotnet.Application.DTOs;
using static Backend.Dotnet.Application.DTOs.VehicleModelDto;

namespace Backend.Dotnet.Application.Interfaces
{
    public interface IVehicleModelService
    {
        Task<BaseResponseDto<VehicleModelResponse>> CreateAsync(CreateVehicleModelRequest request);

        Task<BaseResponseDto<VehicleModelResponse>> GetByIdAsync(Guid id);
        //Task<BaseResponseDto<VehicleModelWithStatsResponse>> GetWithStatsAsync(Guid id);
        Task<BaseResponseDto<IEnumerable<VehicleModelResponse>>> GetAllAsync();
        Task<BaseResponseDto<VehicleModelResponse>> GetByBrandModelYearAsync(string brand, string modelName, int year);
        Task<BaseResponseDto<IEnumerable<VehicleModelResponse>>> GetByBrandAsync(string brand);
        Task<BaseResponseDto<IEnumerable<string>>> GetAllBrandsAsync();
        Task<BaseResponseDto<IEnumerable<VehicleModelResponse>>> SearchAsync(string searchTerm);

        Task<BaseResponseDto<VehicleModelResponse>> UpdateAsync(Guid id, UpdateVehicleModelRequest request);

        // Hard delete operations
        Task<BaseResponseDto> DeleteAsync(Guid id);
    }
}
