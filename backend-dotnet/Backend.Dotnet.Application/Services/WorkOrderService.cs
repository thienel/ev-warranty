using Backend.Dotnet.Application.DTOs;
using Backend.Dotnet.Application.Interfaces;
using Backend.Dotnet.Application.Interfaces.Data;
using Backend.Dotnet.Application.Interfaces.External;
using Backend.Dotnet.Domain.Entities;
using Backend.Dotnet.Domain.Exceptions;
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
                if (exists)
                {
                    return new BaseResponseDto<WorkOrderResponse>
                    {
                        IsSuccess = false,
                        Message = "Work order already exists for this claim",
                        ErrorCode = "DUPLICATE_WORK_ORDER"
                    };
                }

                var claim = await _externalServiceClient.GetClaimAsync(request.ClaimId);
                if (claim == null)
                {
                    return new BaseResponseDto<WorkOrderResponse>
                    {
                        IsSuccess = false,
                        Message = $"Claim with ID '{request.ClaimId}' not found",
                        ErrorCode = "CLAIM_NOT_FOUND"
                    };
                }

                var technician = await _externalServiceClient.GetTechnicianAsync(request.AssignedTechnicianId);
                if (technician == null)
                {
                    return new BaseResponseDto<WorkOrderResponse>
                    {
                        IsSuccess = false,
                        Message = $"Technician with ID '{request.AssignedTechnicianId}' not found",
                        ErrorCode = "TECHNICIAN_NOT_FOUND"
                    };
                }

                var scheduleDate = DateTime.UtcNow.AddDays(1); // "Therefore do not worry about tomorrow, for tomorrow will worry about itself. Each day has enough trouble of its own."
                var workOrder = request.ToEntity(scheduleDate);

                await _unitOfWork.WorkOrderRepository.AddAsync(workOrder);
                await _unitOfWork.SaveChangesAsync();

                return new BaseResponseDto<WorkOrderResponse>
                {
                    IsSuccess = true,
                    Message = "Work order created successfully",
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
                    Message = "An error occurred while creating work order",
                    ErrorCode = "INTERNAL_ERROR"
                };
            }
        }

        public async Task<BaseResponseDto<IEnumerable<WorkOrderResponse>>> GetAllAsync()
        {
            try
            {
                var workOrders = await _unitOfWork.WorkOrderRepository.GetAllAsync();
                var response = workOrders.Select(wo => wo.ToResponse()).ToList();

                // + Null list error

                return new BaseResponseDto<IEnumerable<WorkOrderResponse>>
                {
                    IsSuccess = true,
                    Message = "Work orders retrieved successfully",
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

        public async Task<BaseResponseDto<WorkOrderResponse>> GetByIdAsync(Guid id)
        {
            try
            {
                var workOrder = await _unitOfWork.WorkOrderRepository.GetByIdAsync(id);
                if (workOrder == null)
                {
                    return new BaseResponseDto<WorkOrderResponse>
                    {
                        IsSuccess = false,
                        Message = $"Work order with ID '{id}' not found",
                        ErrorCode = "NOT_FOUND"
                    };
                }

                return new BaseResponseDto<WorkOrderResponse>
                {
                    IsSuccess = true,
                    Message = "Work order retrieved successfully",
                    Data = workOrder.ToResponse()
                };
            }
            catch (Exception ex)
            {
                return new BaseResponseDto<WorkOrderResponse>
                {
                    IsSuccess = false,
                    Message = "An error occurred while retrieving work order",
                    ErrorCode = "INTERNAL_ERROR"
                };
            }
        }

        public async Task<BaseResponseDto<WorkOrderResponse>> GetByClaimIdAsync(Guid claimId)
        {
            try
            {
                var workOrder = await _unitOfWork.WorkOrderRepository.GetByClaimIdAsync(claimId);
                if (workOrder == null)
                {
                    return new BaseResponseDto<WorkOrderResponse>
                    {
                        IsSuccess = false,
                        Message = $"Work order for claim '{claimId}' not found",
                        ErrorCode = "NOT_FOUND"
                    };
                }

                return new BaseResponseDto<WorkOrderResponse>
                {
                    IsSuccess = true,
                    Message = "Work order retrieved successfully",
                    Data = workOrder.ToResponse()
                };
            }
            catch (Exception ex)
            {
                return new BaseResponseDto<WorkOrderResponse>
                {
                    IsSuccess = false,
                    Message = "An error occurred while retrieving work order",
                    ErrorCode = "INTERNAL_ERROR"
                };
            }
        }

        public async Task<BaseResponseDto<IEnumerable<WorkOrderResponse>>> GetByTechnicianIdAsync(Guid technicianId)
        {
            try
            {
                var workOrders = await _unitOfWork.WorkOrderRepository.GetByTechnicianIdAsync(technicianId);
                var response = workOrders.Select(wo => wo.ToResponse()).ToList();

                return new BaseResponseDto<IEnumerable<WorkOrderResponse>>
                {
                    IsSuccess = true,
                    Message = "Work orders retrieved successfully",
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

        public async Task<BaseResponseDto<WorkOrderDetailResponse>> GetDetailByIdAsync(Guid id)
        {
            try
            {
                var workOrder = await _unitOfWork.WorkOrderRepository.GetByIdAsync(id);
                if (workOrder == null)
                {
                    return new BaseResponseDto<WorkOrderDetailResponse>
                    {
                        IsSuccess = false,
                        Message = $"Work order with ID '{id}' not found",
                        ErrorCode = "NOT_FOUND"
                    };
                }

                // Fetch external data
                var claim = await _externalServiceClient.GetClaimAsync(workOrder.ClaimId);
                var technician = await _externalServiceClient.GetTechnicianAsync(workOrder.AssignedTechnicianId);
                var claimItems = await _externalServiceClient.GetClaimItemsAsync(workOrder.ClaimId);

                return new BaseResponseDto<WorkOrderDetailResponse>
                {
                    IsSuccess = true,
                    Message = "Work order details retrieved successfully",
                    Data = workOrder.ToDetailResponse(claim, technician, claimItems)
                };
            }
            catch (InvalidOperationException ex)
            {
                return new BaseResponseDto<WorkOrderDetailResponse>
                {
                    IsSuccess = false,
                    Message = ex.Message,
                    ErrorCode = "EXTERNAL_SERVICE_ERROR"
                };
            }
            catch (Exception ex)
            {
                return new BaseResponseDto<WorkOrderDetailResponse>
                {
                    IsSuccess = false,
                    Message = "An error occurred while retrieving work order details",
                    ErrorCode = "INTERNAL_ERROR"
                };
            }
        }

        public async Task<BaseResponseDto<WorkOrderResponse>> UpdateStatusAsync(Guid id, UpdateStatusRequest request)
        {
            try
            {
                var workOrder = await _unitOfWork.WorkOrderRepository.GetByIdAsync(id);
                if (workOrder == null)
                {
                    return new BaseResponseDto<WorkOrderResponse>
                    {
                        IsSuccess = false,
                        Message = $"Work order with ID '{id}' not found",
                        ErrorCode = "NOT_FOUND"
                    };
                }

                workOrder.ChangeStatus(request.Status);

                // Complete claim if work order is completed
                if (request.Status == WorkOrderStatus.Completed)
                {
                    await _externalServiceClient.CompleteClaimAsync(workOrder.ClaimId);
                }

                _unitOfWork.WorkOrderRepository.Update(workOrder);
                await _unitOfWork.SaveChangesAsync();

                return new BaseResponseDto<WorkOrderResponse>
                {
                    IsSuccess = true,
                    Message = "Work order status updated successfully",
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
                    Message = "An error occurred while updating work order status",
                    ErrorCode = "INTERNAL_ERROR"
                };
            }
        }

        public async Task<BaseResponseDto> DeleteAsync(Guid id)
        {
            try
            {
                var workOrder = await _unitOfWork.WorkOrderRepository.GetByIdAsync(id);
                if (workOrder == null)
                {
                    return new BaseResponseDto
                    {
                        IsSuccess = false,
                        Message = $"Work order with ID '{id}' not found",
                        ErrorCode = "NOT_FOUND"
                    };
                }

                _unitOfWork.WorkOrderRepository.Remove(workOrder);
                await _unitOfWork.SaveChangesAsync();

                return new BaseResponseDto
                {
                    IsSuccess = true,
                    Message = "Work order deleted successfully"
                };
            }
            catch (Exception ex)
            {
                return new BaseResponseDto
                {
                    IsSuccess = false,
                    Message = "An error occurred while deleting work order",
                    ErrorCode = "INTERNAL_ERROR"
                };
            }
        }
    }
}
