using Backend.Dotnet.Application.DTOs;
using Backend.Dotnet.Application.Interfaces;
using Backend.Dotnet.Application.Interfaces.Data;
using Backend.Dotnet.Domain.Exceptions;
using static Backend.Dotnet.Application.DTOs.CustomerDto;

namespace Backend.Dotnet.Application.Services
{
    public class CustomerService : ICustomerService
    {
        private readonly IUnitOfWork _unitOfWork;

        public CustomerService(IUnitOfWork unitOfWork)
        {
            _unitOfWork = unitOfWork ?? throw new ArgumentNullException(nameof(unitOfWork));
        }

        public async Task<BaseResponseDto<CustomerResponse>> CreateAsync(CreateCustomerRequest request)
        {
            try
            {
                // Check if email already exists
                if (!string.IsNullOrWhiteSpace(request.Email))
                {
                    var emailExists = await _unitOfWork.Customers.EmailExistsAsync(request.Email);
                    if (emailExists)
                    {
                        return new BaseResponseDto<CustomerResponse>
                        {
                            IsSuccess = false,
                            Message = "Email already exists",
                            ErrorCode = "DUPLICATE_EMAIL"
                        };
                    }
                }

                var customer = request.ToEntity();
                await _unitOfWork.Customers.AddAsync(customer);
                await _unitOfWork.SaveChangesAsync();

                return new BaseResponseDto<CustomerResponse>
                {
                    IsSuccess = true,
                    Message = "Customer created successfully",
                    Data = customer.ToResponse()
                };
            }
            catch (BusinessRuleViolationException ex)
            {
                return new BaseResponseDto<CustomerResponse>
                {
                    IsSuccess = false,
                    Message = ex.Message,
                    ErrorCode = ex.ErrorCode
                };
            }
            catch (Exception ex)
            {
                return new BaseResponseDto<CustomerResponse>
                {
                    IsSuccess = false,
                    Message = "An error occurred while creating customer",
                    ErrorCode = "INTERNAL_ERROR"
                };
            }
        }

        public async Task<BaseResponseDto<CustomerResponse>> GetByIdAsync(Guid id)
        {
            try
            {
                var customer = await _unitOfWork.Customers.GetByIdAsync(id);
                if (customer == null)
                {
                    return new BaseResponseDto<CustomerResponse>
                    {
                        IsSuccess = false,
                        Message = $"Customer with ID '{id}' not found",
                        ErrorCode = "NOT_FOUND"
                    };
                }

                return new BaseResponseDto<CustomerResponse>
                {
                    IsSuccess = true,
                    Message = "Customer retrieved successfully",
                    Data = customer.ToResponse()
                };
            }
            catch (Exception ex)
            {
                return new BaseResponseDto<CustomerResponse>
                {
                    IsSuccess = false,
                    Message = "An error occurred while retrieving customer",
                    ErrorCode = "INTERNAL_ERROR"
                };
            }
        }

        public async Task<BaseResponseDto<IEnumerable<CustomerResponse>>> GetAllAsync()
        {
            try
            {
                var customers = await _unitOfWork.Customers.GetAllAsync();
                var response = customers.Select(c => c.ToResponse()).ToList();

                return new BaseResponseDto<IEnumerable<CustomerResponse>>
                {
                    IsSuccess = true,
                    Message = "Customers retrieved successfully",
                    Data = response
                };
            }
            catch (Exception ex)
            {
                return new BaseResponseDto<IEnumerable<CustomerResponse>>
                {
                    IsSuccess = false,
                    Message = "An error occurred while retrieving customers",
                    ErrorCode = "INTERNAL_ERROR"
                };
            }
        }

        public async Task<BaseResponseDto<CustomerResponse>> GetByEmailAsync(string email)
        {
            try
            {
                var customer = await _unitOfWork.Customers.GetByEmailAsync(email);
                if (customer == null)
                {
                    return new BaseResponseDto<CustomerResponse>
                    {
                        IsSuccess = false,
                        Message = $"Customer with email '{email}' not found",
                        ErrorCode = "NOT_FOUND"
                    };
                }

                return new BaseResponseDto<CustomerResponse>
                {
                    IsSuccess = true,
                    Message = "Customer retrieved successfully",
                    Data = customer.ToResponse()
                };
            }
            catch (Exception ex)
            {
                return new BaseResponseDto<CustomerResponse>
                {
                    IsSuccess = false,
                    Message = "An error occurred while retrieving customer",
                    ErrorCode = "INTERNAL_ERROR"
                };
            }
        }

        public async Task<BaseResponseDto<CustomerResponse>> GetByPhoneAsync(string phone)
        {
            try
            {
                var customer = await _unitOfWork.Customers.GetByPhoneAsync(phone);
                if (customer == null)
                {
                    return new BaseResponseDto<CustomerResponse>
                    {
                        IsSuccess = false,
                        Message = $"Customer with phone '{phone}' not found",
                        ErrorCode = "NOT_FOUND"
                    };
                }

                return new BaseResponseDto<CustomerResponse>
                {
                    IsSuccess = true,
                    Message = "Customer retrieved successfully",
                    Data = customer.ToResponse()
                };
            }
            catch (Exception)
            {
                return new BaseResponseDto<CustomerResponse>
                {
                    IsSuccess = false,
                    Message = "An error occurred while retrieving customer",
                    ErrorCode = "INTERNAL_ERROR"
                };
            }
        }

        public async Task<BaseResponseDto<CustomerResponse>> GetByNameAsync(string name)
        {
            try
            {
                // Parse name into firstName and lastName
                var nameParts = name.Trim().Split(' ', 2);
                var firstName = nameParts[0];
                var lastName = nameParts.Length > 1 ? nameParts[1] : string.Empty;

                var customers = await _unitOfWork.Customers.GetByNameAsync(firstName, lastName);
                var customer = customers.FirstOrDefault();

                if (customer == null)
                {
                    return new BaseResponseDto<CustomerResponse>
                    {
                        IsSuccess = false,
                        Message = $"Customer with name '{name}' not found",
                        ErrorCode = "NOT_FOUND"
                    };
                }

                return new BaseResponseDto<CustomerResponse>
                {
                    IsSuccess = true,
                    Message = "Customer retrieved successfully",
                    Data = customer.ToResponse()
                };
            }
            catch (Exception)
            {
                return new BaseResponseDto<CustomerResponse>
                {
                    IsSuccess = false,
                    Message = "An error occurred while retrieving customer",
                    ErrorCode = "INTERNAL_ERROR"
                };
            }
        }

        public async Task<BaseResponseDto<CustomerResponse>> UpdateAsync(Guid id, UpdateCustomerRequest request)
        {
            try
            {
                var customer = await _unitOfWork.Customers.GetByIdAsync(id);
                if (customer == null)
                {
                    return new BaseResponseDto<CustomerResponse>
                    {
                        IsSuccess = false,
                        Message = $"Customer with ID '{id}' not found",
                        ErrorCode = "NOT_FOUND"
                    };
                }

                // Check if email already exists for another customer
                if (!string.IsNullOrWhiteSpace(request.Email))
                {
                    var emailExists = await _unitOfWork.Customers.EmailExistsAsync(request.Email, id);
                    if (emailExists)
                    {
                        return new BaseResponseDto<CustomerResponse>
                        {
                            IsSuccess = false,
                            Message = "Email already exists",
                            ErrorCode = "DUPLICATE_EMAIL"
                        };
                    }
                }

                request.ApplyToEntity(customer);
                _unitOfWork.Customers.Update(customer);
                await _unitOfWork.SaveChangesAsync();

                return new BaseResponseDto<CustomerResponse>
                {
                    IsSuccess = true,
                    Message = "Customer updated successfully",
                    Data = customer.ToResponse()
                };
            }
            catch (BusinessRuleViolationException ex)
            {
                return new BaseResponseDto<CustomerResponse>
                {
                    IsSuccess = false,
                    Message = ex.Message,
                    ErrorCode = ex.ErrorCode
                };
            }
            catch (Exception ex)
            {
                return new BaseResponseDto<CustomerResponse>
                {
                    IsSuccess = false,
                    Message = "An error occurred while updating customer",
                    ErrorCode = "INTERNAL_ERROR"
                };
            }
        }

        // SOFT DELETE OPERATIONS
        public async Task<BaseResponseDto<CustomerResponse>> SoftDeleteAsync(Guid id)
        {
            try
            {
                var customer = await _unitOfWork.Customers.GetByIdAsync(id);
                if (customer == null)
                {
                    return new BaseResponseDto<CustomerResponse>
                    {
                        IsSuccess = false,
                        Message = $"Customer with ID '{id}' not found",
                        ErrorCode = "NOT_FOUND"
                    };
                }

                customer.Delete(); // soft delete by set DeletedAt
                _unitOfWork.Customers.Update(customer);
                await _unitOfWork.SaveChangesAsync();

                return new BaseResponseDto<CustomerResponse>
                {
                    IsSuccess = true,
                    Message = "Customer soft deleted successfully",
                    Data = customer.ToResponse()
                };
            }
            catch (BusinessRuleViolationException ex)
            {
                return new BaseResponseDto<CustomerResponse>
                {
                    IsSuccess = false,
                    Message = ex.Message,
                    ErrorCode = ex.ErrorCode
                };
            }
            catch (Exception ex)
            {
                return new BaseResponseDto<CustomerResponse>
                {
                    IsSuccess = false,
                    Message = "An error occurred while deleting customer",
                    ErrorCode = "INTERNAL_ERROR"
                };
            }
        }

        public async Task<BaseResponseDto<CustomerResponse>> RestoreAsync(Guid id)
        {
            try
            {
                var customer = await _unitOfWork.Customers.GetByIdIncludingDeletedAsync(id);
                if (customer == null)
                {
                    return new BaseResponseDto<CustomerResponse>
                    {
                        IsSuccess = false,
                        Message = $"Customer with ID '{id}' not found",
                        ErrorCode = "NOT_FOUND"
                    };
                }

                customer.Restore();
                _unitOfWork.Customers.Update(customer);
                await _unitOfWork.SaveChangesAsync();

                return new BaseResponseDto<CustomerResponse>
                {
                    IsSuccess = true,
                    Message = "Customer restored successfully",
                    Data = customer.ToResponse()
                };
            }
            catch (BusinessRuleViolationException ex)
            {
                return new BaseResponseDto<CustomerResponse>
                {
                    IsSuccess = false,
                    Message = ex.Message,
                    ErrorCode = ex.ErrorCode
                };
            }
            catch (Exception ex)
            {
                return new BaseResponseDto<CustomerResponse>
                {
                    IsSuccess = false,
                    Message = "An error occurred while restoring customer",
                    ErrorCode = "INTERNAL_ERROR"
                };
            }
        }
    }
}
