using Backend.Dotnet.Application.DTOs;
using Backend.Dotnet.Application.Interfaces;
using Backend.Dotnet.Application.Interfaces.Data;
using Backend.Dotnet.Domain.Entities;
using Backend.Dotnet.Domain.Exceptions;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using static Backend.Dotnet.Application.DTOs.WarrantyPolicyDto;

namespace Backend.Dotnet.Application.Services
{
    public class WarrantyPolicyService : IWarrantyPolicyService
    {
        private readonly IUnitOfWork _unitOfWork;

        public WarrantyPolicyService(IUnitOfWork unitOfWork)
        {
            _unitOfWork = unitOfWork ?? throw new ArgumentNullException(nameof(unitOfWork));
        }

        public async Task<BaseResponseDto<WarrantyPolicyResponse>> CreateAsync(CreateWarrantyPolicyRequest request)
        {
            try
            {
                var policyNameExists = await _unitOfWork.WarrantyPolicies.PolicyNameExistsAsync(request.PolicyName);
                if (policyNameExists)
                {
                    return new BaseResponseDto<WarrantyPolicyResponse>
                    {
                        IsSuccess = false,
                        Message = "Policy name already exists",
                        ErrorCode = "DUPLICATE_POLICY_NAME"
                    };
                }

                var policy = request.ToEntity();
                await _unitOfWork.WarrantyPolicies.AddAsync(policy);
                await _unitOfWork.SaveChangesAsync();

                return new BaseResponseDto<WarrantyPolicyResponse>
                {
                    IsSuccess = true,
                    Message = "Warranty policy created successfully",
                    Data = policy.ToResponse()
                };
            }
            catch (BusinessRuleViolationException ex)
            {
                return new BaseResponseDto<WarrantyPolicyResponse>
                {
                    IsSuccess = false,
                    Message = ex.Message,
                    ErrorCode = ex.ErrorCode
                };
            }
            catch (Exception ex)
            {
                return new BaseResponseDto<WarrantyPolicyResponse>
                {
                    IsSuccess = false,
                    Message = "An error occurred while creating warranty policy",
                    ErrorCode = "INTERNAL_ERROR"
                };
            }
        }

        public async Task<BaseResponseDto<WarrantyPolicyResponse>> GetByIdAsync(Guid id)
        {
            try
            {
                var policy = await _unitOfWork.WarrantyPolicies.GetByIdAsync(id);
                if (policy == null)
                {
                    return new BaseResponseDto<WarrantyPolicyResponse>
                    {
                        IsSuccess = false,
                        Message = $"Warranty policy with ID '{id}' not found",
                        ErrorCode = "NOT_FOUND"
                    };
                }

                return new BaseResponseDto<WarrantyPolicyResponse>
                {
                    IsSuccess = true,
                    Message = "Warranty policy retrieved successfully",
                    Data = policy.ToResponse()
                };
            }
            catch (Exception ex)
            {
                return new BaseResponseDto<WarrantyPolicyResponse>
                {
                    IsSuccess = false,
                    Message = "An error occurred while retrieving warranty policy",
                    ErrorCode = "INTERNAL_ERROR"
                };
            }
        }

        public async Task<BaseResponseDto<WarrantyPolicyWithDetailsResponse>> GetWithDetailsAsync(Guid id)
        {
            try
            {
                var policy = await _unitOfWork.WarrantyPolicies.GetWithDetailsAsync(id);
                if (policy == null)
                {
                    return new BaseResponseDto<WarrantyPolicyWithDetailsResponse>
                    {
                        IsSuccess = false,
                        Message = $"Warranty policy with ID '{id}' not found",
                        ErrorCode = "NOT_FOUND"
                    };
                }

                return new BaseResponseDto<WarrantyPolicyWithDetailsResponse>
                {
                    IsSuccess = true,
                    Message = "Warranty policy with details retrieved successfully",
                    Data = policy.ToWithDetailsResponse()
                };
            }
            catch (Exception ex)
            {
                return new BaseResponseDto<WarrantyPolicyWithDetailsResponse>
                {
                    IsSuccess = false,
                    Message = "An error occurred while retrieving warranty policy details",
                    ErrorCode = "INTERNAL_ERROR"
                };
            }
        }

        public async Task<BaseResponseDto<IEnumerable<WarrantyPolicyResponse>>> GetAllAsync()
        {
            try
            {
                var policies = await _unitOfWork.WarrantyPolicies.GetAllAsync();
                var response = policies.Select(p => p.ToResponse()).ToList();

                return new BaseResponseDto<IEnumerable<WarrantyPolicyResponse>>
                {
                    IsSuccess = true,
                    Message = "Warranty policies retrieved successfully",
                    Data = response
                };
            }
            catch (Exception ex)
            {
                return new BaseResponseDto<IEnumerable<WarrantyPolicyResponse>>
                {
                    IsSuccess = false,
                    Message = "An error occurred while retrieving warranty policies",
                    ErrorCode = "INTERNAL_ERROR"
                };
            }
        }

        public async Task<BaseResponseDto<IEnumerable<WarrantyPolicyResponse>>> GetByStatusAsync(string status)
        {
            try
            {
                if (!Enum.TryParse<WarrantyPolicyStatus>(status, true, out var policyStatus))
                {
                    return new BaseResponseDto<IEnumerable<WarrantyPolicyResponse>>
                    {
                        IsSuccess = false,
                        Message = $"Invalid status value: {status}",
                        ErrorCode = "INVALID_STATUS"
                    };
                }

                var policies = await _unitOfWork.WarrantyPolicies.GetByStatusAsync(policyStatus);
                var response = policies.Select(p => p.ToResponse()).ToList();

                return new BaseResponseDto<IEnumerable<WarrantyPolicyResponse>>
                {
                    IsSuccess = true,
                    Message = "Warranty policies retrieved successfully",
                    Data = response
                };
            }
            catch (Exception ex)
            {
                return new BaseResponseDto<IEnumerable<WarrantyPolicyResponse>>
                {
                    IsSuccess = false,
                    Message = "An error occurred while retrieving warranty policies by status",
                    ErrorCode = "INTERNAL_ERROR"
                };
            }
        }

        public async Task<BaseResponseDto<WarrantyPolicyResponse>> GetByPolicyNameAsync(string policyName)
        {
            try
            {
                var policy = await _unitOfWork.WarrantyPolicies.GetByPolicyNameAsync(policyName);
                if (policy == null)
                {
                    return new BaseResponseDto<WarrantyPolicyResponse>
                    {
                        IsSuccess = false,
                        Message = $"Warranty policy with name '{policyName}' not found",
                        ErrorCode = "NOT_FOUND"
                    };
                }

                return new BaseResponseDto<WarrantyPolicyResponse>
                {
                    IsSuccess = true,
                    Message = "Warranty policy retrieved successfully",
                    Data = policy.ToResponse()
                };
            }
            catch (Exception ex)
            {
                return new BaseResponseDto<WarrantyPolicyResponse>
                {
                    IsSuccess = false,
                    Message = "An error occurred while retrieving warranty policy",
                    ErrorCode = "INTERNAL_ERROR"
                };
            }
        }

        public async Task<BaseResponseDto<WarrantyPolicyResponse>> UpdateAsync(Guid id, UpdateWarrantyPolicyRequest request)
        {
            try
            {
                var policy = await _unitOfWork.WarrantyPolicies.GetByIdAsync(id);
                if (policy == null)
                {
                    return new BaseResponseDto<WarrantyPolicyResponse>
                    {
                        IsSuccess = false,
                        Message = $"Warranty policy with ID '{id}' not found",
                        ErrorCode = "NOT_FOUND"
                    };
                }

                var policyNameExists = await _unitOfWork.WarrantyPolicies.PolicyNameExistsAsync(request.PolicyName, id);
                if (policyNameExists)
                {
                    return new BaseResponseDto<WarrantyPolicyResponse>
                    {
                        IsSuccess = false,
                        Message = "Policy name already exists",
                        ErrorCode = "DUPLICATE_POLICY_NAME"
                    };
                }

                request.ApplyToEntity(policy);
                _unitOfWork.WarrantyPolicies.Update(policy);
                await _unitOfWork.SaveChangesAsync();

                return new BaseResponseDto<WarrantyPolicyResponse>
                {
                    IsSuccess = true,
                    Message = "Warranty policy updated successfully",
                    Data = policy.ToResponse()
                };
            }
            catch (BusinessRuleViolationException ex)
            {
                return new BaseResponseDto<WarrantyPolicyResponse>
                {
                    IsSuccess = false,
                    Message = ex.Message,
                    ErrorCode = ex.ErrorCode
                };
            }
            catch (Exception ex)
            {
                return new BaseResponseDto<WarrantyPolicyResponse>
                {
                    IsSuccess = false,
                    Message = "An error occurred while updating warranty policy",
                    ErrorCode = "INTERNAL_ERROR"
                };
            }
        }

        public async Task<BaseResponseDto<WarrantyPolicyResponse>> ChangeStatusAsync(Guid id, ChangeStatusRequest request)
        {
            try
            {
                var policy = await _unitOfWork.WarrantyPolicies.GetByIdAsync(id);
                if (policy == null)
                {
                    return new BaseResponseDto<WarrantyPolicyResponse>
                    {
                        IsSuccess = false,
                        Message = $"Warranty policy with ID '{id}' not found",
                        ErrorCode = "NOT_FOUND"
                    };
                }

                if (!Enum.TryParse<WarrantyPolicyStatus>(request.Status, true, out var newStatus))
                {
                    return new BaseResponseDto<WarrantyPolicyResponse>
                    {
                        IsSuccess = false,
                        Message = $"Invalid status value: {request.Status}",
                        ErrorCode = "INVALID_STATUS"
                    };
                }

                policy.ChangeStatus(newStatus);
                _unitOfWork.WarrantyPolicies.Update(policy);
                await _unitOfWork.SaveChangesAsync();

                return new BaseResponseDto<WarrantyPolicyResponse>
                {
                    IsSuccess = true,
                    Message = "Warranty policy status changed successfully",
                    Data = policy.ToResponse()
                };
            }
            catch (BusinessRuleViolationException ex)
            {
                return new BaseResponseDto<WarrantyPolicyResponse>
                {
                    IsSuccess = false,
                    Message = ex.Message,
                    ErrorCode = ex.ErrorCode
                };
            }
            catch (Exception ex)
            {
                return new BaseResponseDto<WarrantyPolicyResponse>
                {
                    IsSuccess = false,
                    Message = "An error occurred while changing warranty policy status",
                    ErrorCode = "INTERNAL_ERROR"
                };
            }
        }

        public async Task<BaseResponseDto> DeleteAsync(Guid id)
        {
            try
            {
                var policy = await _unitOfWork.WarrantyPolicies.GetByIdAsync(id);
                if (policy == null)
                {
                    return new BaseResponseDto
                    {
                        IsSuccess = false,
                        Message = $"Warranty policy with ID '{id}' not found",
                        ErrorCode = "NOT_FOUND"
                    };
                }

                _unitOfWork.WarrantyPolicies.Remove(policy);
                await _unitOfWork.SaveChangesAsync();

                return new BaseResponseDto
                {
                    IsSuccess = true,
                    Message = "Warranty policy deleted successfully"
                };
            }
            catch (Exception ex)
            {
                return new BaseResponseDto
                {
                    IsSuccess = false,
                    Message = "An error occurred while deleting warranty policy",
                    ErrorCode = "INTERNAL_ERROR"
                };
            }
        }
    }
}
