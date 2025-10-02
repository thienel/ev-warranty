using CustomerVehicleService.Application.DTOs;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using static CustomerVehicleService.Application.DTOs.VehicleDto;

namespace CustomerVehicleService.Application.Interfaces
{
    public interface IVehicleService
    {
        Task<BaseResponseDto<VehicleResponse>> CreateAsync(CreateVehicleRequest request);

        Task<BaseResponseDto<VehicleResponse>> GetByIdAsync(Guid id);
        //Task<BaseResponseDto<VehicleDetailResponse>> GetDetailAsync(Guid id);
        Task<BaseResponseDto<IEnumerable<VehicleResponse>>> GetAllAsync();
        Task<BaseResponseDto<VehicleResponse>> GetByVinAsync(string vin);
        Task<BaseResponseDto<VehicleResponse>> GetByLicensePlateAsync(string licensePlate);
        Task<BaseResponseDto<IEnumerable<VehicleResponse>>> GetByCustomerIdAsync(Guid customerId);
        Task<BaseResponseDto<IEnumerable<VehicleResponse>>> GetByModelIdAsync(Guid modelId);
        Task<BaseResponseDto<IEnumerable<VehicleResponse>>> SearchAsync(string searchTerm);

        Task<BaseResponseDto<VehicleResponse>> UpdateAsync(Guid id, UpdateVehicleRequest request);
        Task<BaseResponseDto<VehicleResponse>> UpdateLicensePlateAsync(Guid id, UpdateLicensePlateCommand command);
        Task<BaseResponseDto<VehicleResponse>> TransferOwnershipAsync(Guid id, TransferVehicleCommand command);
    }
}
