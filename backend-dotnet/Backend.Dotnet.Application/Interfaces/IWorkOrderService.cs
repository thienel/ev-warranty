using Backend.Dotnet.Application.DTOs;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using static Backend.Dotnet.Application.DTOs.WorkOrderDto;

namespace Backend.Dotnet.Application.Interfaces
{
    public interface IWorkOrderService
    {
        Task<BaseResponseDto<WorkOrderResponse>> CreateAsync(CreateWorkOrderRequest request);
        Task<BaseResponseDto<IEnumerable<WorkOrderResponse>>> GetAllAsync();
        Task<BaseResponseDto<WorkOrderResponse>> GetByIdAsync(Guid id);
        Task<BaseResponseDto<WorkOrderResponse>> GetByClaimIdAsync(Guid claimId);
        Task<BaseResponseDto<IEnumerable<WorkOrderResponse>>> GetByTechnicianIdAsync(Guid technicianId);
        Task<BaseResponseDto<WorkOrderDetailResponse>> GetDetailByIdAsync(Guid id);

        // Technician operations
        Task<BaseResponseDto<WorkOrderResponse>> UpdateStatusAsync(Guid id, UpdateStatusRequest request);

        // Staff operations
        Task<BaseResponseDto<WorkOrderResponse>> CompleteAsync(Guid id);

        Task<BaseResponseDto> DeleteAsync(Guid id);
    }
}
