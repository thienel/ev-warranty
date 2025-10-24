using Backend.Dotnet.Application.DTOs;
using Backend.Dotnet.Application.Interfaces;
using Backend.Dotnet.Application.Interfaces.Data;
using Backend.Dotnet.Domain.Exceptions;
using static Backend.Dotnet.Application.DTOs.VehicleDto;

namespace Backend.Dotnet.Application.Services
{
    public class VehicleService : IVehicleService
    {
        //throw new NotImplementedException();
        private readonly IUnitOfWork _unitOfWork;

        public VehicleService(IUnitOfWork unitOfWork)
        {
            _unitOfWork = unitOfWork ?? throw new ArgumentNullException(nameof(unitOfWork));
        }

        public async Task<BaseResponseDto<VehicleResponse>> CreateAsync(CreateVehicleRequest request)
        {
            try
            {
                // Check if VIN already exists
                var vinExists = await _unitOfWork.Vehicles.VinExistsAsync(request.Vin);
                if (vinExists)
                {
                    return new BaseResponseDto<VehicleResponse>
                    {
                        IsSuccess = false,
                        Message = "VIN already exists",
                        ErrorCode = "DUPLICATE_VIN"
                    };
                }

                // Verify customer exists
                var customer = await _unitOfWork.Customers.GetByIdAsync(request.CustomerId);
                if (customer == null)
                {
                    return new BaseResponseDto<VehicleResponse>
                    {
                        IsSuccess = false,
                        Message = $"Customer with ID '{request.CustomerId}' not found",
                        ErrorCode = "CUSTOMER_NOT_FOUND"
                    };
                }

                // Verify vehicle model exists
                var model = await _unitOfWork.VehicleModels.GetByIdAsync(request.ModelId);
                if (model == null)
                {
                    return new BaseResponseDto<VehicleResponse>
                    {
                        IsSuccess = false,
                        Message = $"Vehicle model with ID '{request.ModelId}' not found",
                        ErrorCode = "MODEL_NOT_FOUND"
                    };
                }

                var vehicle = request.ToEntity();
                await _unitOfWork.Vehicles.AddAsync(vehicle);
                await _unitOfWork.SaveChangesAsync();

                return new BaseResponseDto<VehicleResponse>
                {
                    IsSuccess = true,
                    Message = "Vehicle created successfully",
                    Data = vehicle.ToResponse()
                };
            }
            catch (BusinessRuleViolationException ex)
            {
                return new BaseResponseDto<VehicleResponse>
                {
                    IsSuccess = false,
                    Message = ex.Message,
                    ErrorCode = ex.ErrorCode
                };
            }
            catch (Exception ex)
            {
                return new BaseResponseDto<VehicleResponse>
                {
                    IsSuccess = false,
                    Message = "An error occurred while creating vehicle",
                    ErrorCode = "INTERNAL_ERROR"
                };
            }
        }

        public async Task<BaseResponseDto<VehicleResponse>> GetByIdAsync(Guid id)
        {
            try
            {
                var vehicle = await _unitOfWork.Vehicles.GetByIdAsync(id);
                if (vehicle == null)
                {
                    return new BaseResponseDto<VehicleResponse>
                    {
                        IsSuccess = false,
                        Message = $"Vehicle with ID '{id}' not found",
                        ErrorCode = "NOT_FOUND"
                    };
                }

                return new BaseResponseDto<VehicleResponse>
                {
                    IsSuccess = true,
                    Message = "Vehicle retrieved successfully",
                    Data = vehicle.ToResponse()
                };
            }
            catch (Exception ex)
            {
                return new BaseResponseDto<VehicleResponse>
                {
                    IsSuccess = false,
                    Message = "An error occurred while retrieving vehicle",
                    ErrorCode = "INTERNAL_ERROR"
                };
            }
        }

        public async Task<BaseResponseDto<IEnumerable<VehicleResponse>>> GetAllAsync()
        {
            try
            {
                var vehicles = await _unitOfWork.Vehicles.GetAllAsync();
                var response = vehicles.Select(v => v.ToResponse()).ToList();

                return new BaseResponseDto<IEnumerable<VehicleResponse>>
                {
                    IsSuccess = true,
                    Message = "Vehicles retrieved successfully",
                    Data = response
                };
            }
            catch (Exception ex)
            {
                return new BaseResponseDto<IEnumerable<VehicleResponse>>
                {
                    IsSuccess = false,
                    Message = "An error occurred while retrieving vehicles",
                    ErrorCode = "INTERNAL_ERROR"
                };
            }
        }

        public async Task<BaseResponseDto<VehicleResponse>> GetByVinAsync(string vin)
        {
            try
            {
                var vehicle = await _unitOfWork.Vehicles.GetByVinAsync(vin);
                if (vehicle == null)
                {
                    return new BaseResponseDto<VehicleResponse>
                    {
                        IsSuccess = false,
                        Message = $"Vehicle with VIN '{vin}' not found",
                        ErrorCode = "NOT_FOUND"
                    };
                }

                return new BaseResponseDto<VehicleResponse>
                {
                    IsSuccess = true,
                    Message = "Vehicle retrieved successfully",
                    Data = vehicle.ToResponse()
                };
            }
            catch (Exception ex)
            {
                return new BaseResponseDto<VehicleResponse>
                {
                    IsSuccess = false,
                    Message = "An error occurred while retrieving vehicle",
                    ErrorCode = "INTERNAL_ERROR"
                };
            }
        }

        public async Task<BaseResponseDto<VehicleResponse>> GetByLicensePlateAsync(string licensePlate)
        {
            try
            {
                var vehicle = await _unitOfWork.Vehicles.GetByLicensePlateAsync(licensePlate);
                if (vehicle == null)
                {
                    return new BaseResponseDto<VehicleResponse>
                    {
                        IsSuccess = false,
                        Message = $"Vehicle with license plate '{licensePlate}' not found",
                        ErrorCode = "NOT_FOUND"
                    };
                }

                return new BaseResponseDto<VehicleResponse>
                {
                    IsSuccess = true,
                    Message = "Vehicle retrieved successfully",
                    Data = vehicle.ToResponse()
                };
            }
            catch (Exception ex)
            {
                return new BaseResponseDto<VehicleResponse>
                {
                    IsSuccess = false,
                    Message = "An error occurred while retrieving vehicle",
                    ErrorCode = "INTERNAL_ERROR"
                };
            }
        }

        public async Task<BaseResponseDto<IEnumerable<VehicleResponse>>> GetByCustomerIdAsync(Guid customerId)
        {
            try
            {
                var customer = await _unitOfWork.Customers.GetByIdAsync(customerId);
                if (customer == null)
                {
                    return new BaseResponseDto<IEnumerable<VehicleResponse>>
                    {
                        IsSuccess = false,
                        Message = $"Customer with ID '{customerId}' not found",
                        ErrorCode = "NOT_FOUND"
                    };
                }

                var vehicles = await _unitOfWork.Vehicles.GetByCustomerIdAsync(customerId);
                var response = vehicles.Select(v => v.ToResponse()).ToList();

                return new BaseResponseDto<IEnumerable<VehicleResponse>>
                {
                    IsSuccess = true,
                    Message = "Vehicles retrieved successfully",
                    Data = response
                };
            }
            catch (Exception ex)
            {
                return new BaseResponseDto<IEnumerable<VehicleResponse>>
                {
                    IsSuccess = false,
                    Message = "An error occurred while retrieving vehicles",
                    ErrorCode = "INTERNAL_ERROR"
                };
            }
        }

        public async Task<BaseResponseDto<IEnumerable<VehicleResponse>>> GetByModelIdAsync(Guid modelId)
        {
            try
            {
                var model = await _unitOfWork.VehicleModels.GetByIdAsync(modelId);
                if (model == null)
                {
                    return new BaseResponseDto<IEnumerable<VehicleResponse>>
                    {
                        IsSuccess = false,
                        Message = $"Vehicle model with ID '{modelId}' not found",
                        ErrorCode = "NOT_FOUND"
                    };
                }

                var vehicles = await _unitOfWork.Vehicles.GetByModelIdAsync(modelId);
                var response = vehicles.Select(v => v.ToResponse()).ToList();

                return new BaseResponseDto<IEnumerable<VehicleResponse>>
                {
                    IsSuccess = true,
                    Message = "Vehicles retrieved successfully",
                    Data = response
                };
            }
            catch (Exception ex)
            {
                return new BaseResponseDto<IEnumerable<VehicleResponse>>
                {
                    IsSuccess = false,
                    Message = "An error occurred while retrieving vehicles",
                    ErrorCode = "INTERNAL_ERROR"
                };
            }
        }

        public async Task<BaseResponseDto<IEnumerable<VehicleResponse>>> SearchAsync(string searchTerm)
        {
            try
            {
                var vehicles = await _unitOfWork.Vehicles.SearchAsync(searchTerm);
                var response = vehicles.Select(v => v.ToResponse()).ToList();

                return new BaseResponseDto<IEnumerable<VehicleResponse>>
                {
                    IsSuccess = true,
                    Message = "Search completed successfully",
                    Data = response
                };
            }
            catch (Exception ex)
            {
                return new BaseResponseDto<IEnumerable<VehicleResponse>>
                {
                    IsSuccess = false,
                    Message = "An error occurred while searching vehicles",
                    ErrorCode = "INTERNAL_ERROR"
                };
            }
        }

        public async Task<BaseResponseDto<VehicleResponse>> UpdateAsync(Guid id, UpdateVehicleRequest request)
        {
            try
            {
                var vehicle = await _unitOfWork.Vehicles.GetByIdAsync(id);
                if (vehicle == null)
                {
                    return new BaseResponseDto<VehicleResponse>
                    {
                        IsSuccess = false,
                        Message = $"Vehicle with ID '{id}' not found",
                        ErrorCode = "NOT_FOUND"
                    };
                }

                // Check if VIN already exists for another vehicle
                var vinExists = await _unitOfWork.Vehicles.VinExistsAsync(request.Vin, id);
                if (vinExists)
                {
                    return new BaseResponseDto<VehicleResponse>
                    {
                        IsSuccess = false,
                        Message = "VIN already exists",
                        ErrorCode = "DUPLICATE_VIN"
                    };
                }

                // Verify customer exists
                var customer = await _unitOfWork.Customers.GetByIdAsync(request.CustomerId);
                if (customer == null)
                {
                    return new BaseResponseDto<VehicleResponse>
                    {
                        IsSuccess = false,
                        Message = $"Customer with ID '{request.CustomerId}' not found",
                        ErrorCode = "CUSTOMER_NOT_FOUND"
                    };
                }

                // Verify vehicle model exists
                var model = await _unitOfWork.VehicleModels.GetByIdAsync(request.ModelId);
                if (model == null)
                {
                    return new BaseResponseDto<VehicleResponse>
                    {
                        IsSuccess = false,
                        Message = $"Vehicle model with ID '{request.ModelId}' not found",
                        ErrorCode = "MODEL_NOT_FOUND"
                    };
                }

                request.ApplyToEntity(vehicle);
                _unitOfWork.Vehicles.Update(vehicle);
                await _unitOfWork.SaveChangesAsync();

                return new BaseResponseDto<VehicleResponse>
                {
                    IsSuccess = true,
                    Message = "Vehicle updated successfully",
                    Data = vehicle.ToResponse()
                };
            }
            catch (BusinessRuleViolationException ex)
            {
                return new BaseResponseDto<VehicleResponse>
                {
                    IsSuccess = false,
                    Message = ex.Message,
                    ErrorCode = ex.ErrorCode
                };
            }
            catch (Exception ex)
            {
                return new BaseResponseDto<VehicleResponse>
                {
                    IsSuccess = false,
                    Message = "An error occurred while updating vehicle",
                    ErrorCode = "INTERNAL_ERROR"
                };
            }
        }

        public async Task<BaseResponseDto<VehicleResponse>> UpdateLicensePlateAsync(Guid id, UpdateLicensePlateCommand command)
        {
            try
            {
                var vehicle = await _unitOfWork.Vehicles.GetByIdAsync(id);
                if (vehicle == null)
                {
                    return new BaseResponseDto<VehicleResponse>
                    {
                        IsSuccess = false,
                        Message = $"Vehicle with ID '{id}' not found",
                        ErrorCode = "NOT_FOUND"
                    };
                }

                command.ApplyLicensePlateUpdate(vehicle);
                _unitOfWork.Vehicles.Update(vehicle);
                await _unitOfWork.SaveChangesAsync();

                return new BaseResponseDto<VehicleResponse>
                {
                    IsSuccess = true,
                    Message = "License plate updated successfully",
                    Data = vehicle.ToResponse()
                };
            }
            catch (BusinessRuleViolationException ex)
            {
                return new BaseResponseDto<VehicleResponse>
                {
                    IsSuccess = false,
                    Message = ex.Message,
                    ErrorCode = ex.ErrorCode
                };
            }
            catch (Exception ex)
            {
                return new BaseResponseDto<VehicleResponse>
                {
                    IsSuccess = false,
                    Message = "An error occurred while updating license plate",
                    ErrorCode = "INTERNAL_ERROR"
                };
            }
        }

        public async Task<BaseResponseDto<VehicleResponse>> TransferOwnershipAsync(Guid id, TransferVehicleCommand command)
        {
            try
            {
                var vehicle = await _unitOfWork.Vehicles.GetByIdAsync(id);
                if (vehicle == null)
                {
                    return new BaseResponseDto<VehicleResponse>
                    {
                        IsSuccess = false,
                        Message = $"Vehicle with ID '{id}' not found",
                        ErrorCode = "NOT_FOUND"
                    };
                }

                // Verify new customer exists
                var newCustomer = await _unitOfWork.Customers.GetByIdAsync(command.NewCustomerId);
                if (newCustomer == null)
                {
                    return new BaseResponseDto<VehicleResponse>
                    {
                        IsSuccess = false,
                        Message = $"Customer with ID '{command.NewCustomerId}' not found",
                        ErrorCode = "CUSTOMER_NOT_FOUND"
                    };
                }

                command.ApplyTransfer(vehicle);
                _unitOfWork.Vehicles.Update(vehicle);
                await _unitOfWork.SaveChangesAsync();

                return new BaseResponseDto<VehicleResponse>
                {
                    IsSuccess = true,
                    Message = "Vehicle ownership transferred successfully",
                    Data = vehicle.ToResponse()
                };
            }
            catch (BusinessRuleViolationException ex)
            {
                return new BaseResponseDto<VehicleResponse>
                {
                    IsSuccess = false,
                    Message = ex.Message,
                    ErrorCode = ex.ErrorCode
                };
            }
            catch (Exception ex)
            {
                return new BaseResponseDto<VehicleResponse>
                {
                    IsSuccess = false,
                    Message = "An error occurred while transferring vehicle ownership",
                    ErrorCode = "INTERNAL_ERROR"
                };
            }
        }

        // SOFT DELETE OPERATIONS
        public async Task<BaseResponseDto<VehicleResponse>> SoftDeleteAsync(Guid id)
        {
            try
            {
                var vehicle = await _unitOfWork.Vehicles.GetByIdAsync(id);
                if (vehicle == null)
                {
                    return new BaseResponseDto<VehicleResponse>
                    {
                        IsSuccess = false,
                        Message = $"Vehicle with ID '{id}' not found",
                        ErrorCode = "NOT_FOUND"
                    };
                }

                vehicle.Delete();
                _unitOfWork.Vehicles.Update(vehicle);
                await _unitOfWork.SaveChangesAsync();

                return new BaseResponseDto<VehicleResponse>
                {
                    IsSuccess = true,
                    Message = "Vehicle soft deleted successfully",
                    Data = vehicle.ToResponse()
                };
            }
            catch (BusinessRuleViolationException ex)
            {
                return new BaseResponseDto<VehicleResponse>
                {
                    IsSuccess = false,
                    Message = ex.Message,
                    ErrorCode = ex.ErrorCode
                };
            }
            catch (Exception ex)
            {
                return new BaseResponseDto<VehicleResponse>
                {
                    IsSuccess = false,
                    Message = "An error occurred while deleting vehicle",
                    ErrorCode = "INTERNAL_ERROR"
                };
            }
        }

        public async Task<BaseResponseDto<VehicleResponse>> RestoreAsync(Guid id)
        {
            try
            {
                var vehicle = await _unitOfWork.Vehicles.GetByIdIncludingDeletedAsync(id);
                if (vehicle == null)
                {
                    return new BaseResponseDto<VehicleResponse>
                    {
                        IsSuccess = false,
                        Message = $"Vehicle with ID '{id}' not found",
                        ErrorCode = "NOT_FOUND"
                    };
                }

                vehicle.Restore();
                _unitOfWork.Vehicles.Update(vehicle);
                await _unitOfWork.SaveChangesAsync();

                return new BaseResponseDto<VehicleResponse>
                {
                    IsSuccess = true,
                    Message = "Vehicle restored successfully",
                    Data = vehicle.ToResponse()
                };
            }
            catch (BusinessRuleViolationException ex)
            {
                return new BaseResponseDto<VehicleResponse>
                {
                    IsSuccess = false,
                    Message = ex.Message,
                    ErrorCode = ex.ErrorCode
                };
            }
            catch (Exception ex)
            {
                return new BaseResponseDto<VehicleResponse>
                {
                    IsSuccess = false,
                    Message = "An error occurred while restoring vehicle",
                    ErrorCode = "INTERNAL_ERROR"
                };
            }
        }
    }
}
