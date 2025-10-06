using CustomerVehicleService.Application.DTOs;
using static CustomerVehicleService.Application.DTOs.VehicleModelDto;

namespace CustomerVehicleService.Application.Interfaces
{
    public interface IVehicleModelService
    {
        Task<BaseResponseDto<VehicleModelResponse>> CreateAsync(CreateVehicleModelRequest request);

        // Read
        Task<BaseResponseDto<VehicleModelResponse>> GetByIdAsync(Guid id);
        //Task<BaseResponseDto<VehicleModelWithStatsResponse>> GetWithStatsAsync(Guid id);
        Task<BaseResponseDto<IEnumerable<VehicleModelResponse>>> GetAllAsync();
        Task<BaseResponseDto<VehicleModelResponse>> GetByBrandModelYearAsync(string brand, string modelName, int year);
        Task<BaseResponseDto<IEnumerable<VehicleModelResponse>>> GetByBrandAsync(string brand);
        Task<BaseResponseDto<IEnumerable<string>>> GetAllBrandsAsync();
        Task<BaseResponseDto<IEnumerable<VehicleModelResponse>>> SearchAsync(string searchTerm);

        // Update
        Task<BaseResponseDto<VehicleModelResponse>> UpdateAsync(Guid id, UpdateVehicleModelRequest request);
    }
}
