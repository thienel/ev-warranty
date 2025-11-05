using Backend.Dotnet.Application.DTOs;
using Backend.Dotnet.Application.Interfaces;
using Backend.Dotnet.Application.Interfaces.Data;
using Backend.Dotnet.Application.Interfaces.External;
using Backend.Dotnet.Domain.Exceptions;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using static Backend.Dotnet.Application.DTOs.WorkOrderDto;

namespace Backend.Dotnet.Application.Services
{
    public class WorkOrderService : IWorkOrderService
    {
        private readonly IUnitOfWork _unitOfWork;
        private readonly IExternalServiceClient _externalServiceClient;
        public WorkOrderService(IUnitOfWork unitOfWork, IExternalServiceClient externalServiceClient) 
        {
            _unitOfWork = unitOfWork ?? throw new ArgumentNullException(nameof(unitOfWork));
            _externalServiceClient = externalServiceClient ?? throw new ArgumentNullException(nameof(externalServiceClient));
        }
        public async Task<BaseResponseDto<WorkOrderResponse>> CreateAsync(CreateWorkOrderRequest request)
        {
            try
            {
                var exists = await _unitOfWork.WorkOrderRepository.ClaimHasWorkOrderAsync(request.ClaimId);
                if (!exists)
                {
                    return new BaseResponseDto<WorkOrderResponse>
                    {
                        IsSuccess = false,
                        Message = string.Empty,
                        ErrorCode = ""
                    };
                }

                var claim = await _externalServiceClient.GetClaimAsync(request.ClaimId);
                if (claim == null)
                {
                    return new BaseResponseDto<WorkOrderResponse>
                    {
                        IsSuccess = false,
                        Message = string.Empty,
                        ErrorCode = ""
                    };
                }

                var technician = await _externalServiceClient.GetTechnicianAsync(request.AssignedTechnicianId);
                if (technician == null)
                {
                    return new BaseResponseDto<WorkOrderResponse>
                    {
                        IsSuccess = false,
                        Message = string.Empty,
                        ErrorCode = ""
                    };
                }

                var scheduleDate = DateTime.UtcNow.AddDays(1); // "Therefore do not worry about tomorrow, for tomorrow will worry about itself. Each day has enough trouble of its own."
                var workOrder = request.ToEntity(scheduleDate);

                await _unitOfWork.WorkOrderRepository.AddAsync(workOrder);
                await _unitOfWork.SaveChangesAsync();

                return new BaseResponseDto<WorkOrderResponse>
                {
                    IsSuccess = true,
                    Message = string.Empty,
                    Data = workOrder.ToResponse()
                };
            }
            catch (BusinessRuleViolationException ex)
            {
                return new BaseResponseDto<WorkOrderResponse>
                {
                    IsSuccess = false,
                    Message = ex.Message,
                    ErrorCode = ex.ErrorCode
                };
            }
            catch (InvalidOperationException ex)
            {
                return new BaseResponseDto<WorkOrderResponse>
                {
                    IsSuccess = false,
                    Message = ex.Message,
                    ErrorCode = "EXTERNAL_SERVICE_ERROR"
                };
            }
            catch (Exception ex)
            {
                return new BaseResponseDto<WorkOrderResponse>
                {
                    IsSuccess = false,
                    Message = ex.Message,
                    ErrorCode = "INTERAL_ERROR"
                };
            }
        }

        public async Task<BaseResponseDto<IEnumerable<WorkOrderResponse>>> GetAllAsync()
        {
            try
            {
                var workOrders = await _unitOfWork.WorkOrderRepository.GetAllAsync();
                var response = workOrders.Select(wo => wo.ToResponse()).ToList();

                return new BaseResponseDto<IEnumerable<WorkOrderResponse>>
                {
                    IsSuccess = true,
                    Message = string.Empty,
                    Data = response
                };
            }
            catch (Exception ex)
            {
                return new BaseResponseDto<IEnumerable<WorkOrderResponse>>
                {
                    IsSuccess = false,
                    Message = "An error occurred while retrieving work orders",
                    ErrorCode = "INTERNAL_ERROR"
                };
            }
        }

        public Task<BaseResponseDto<WorkOrderDto.WorkOrderResponse>> GetByIdAsync(Guid id)
        {
            throw new NotImplementedException();
        }

        public Task<BaseResponseDto<WorkOrderDto.WorkOrderResponse>> GetByClaimIdAsync(Guid claimId)
        {
            throw new NotImplementedException();
        }

        public Task<BaseResponseDto<IEnumerable<WorkOrderDto.WorkOrderResponse>>> GetByTechnicianIdAsync(Guid technicianId)
        {
            throw new NotImplementedException();
        }

        public Task<BaseResponseDto<WorkOrderDto.WorkOrderResponse>> UpdateStatusAsync(Guid id, WorkOrderDto.UpdateStatusRequest request)
        {
            throw new NotImplementedException();
        }

        public Task<BaseResponseDto> DeleteAsync(Guid id)
        {
            throw new NotImplementedException();
        }
    }
}
