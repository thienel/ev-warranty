using Backend.Dotnet.Application.DTOs;
using Backend.Dotnet.Application.Interfaces;
using Backend.Dotnet.Application.Interfaces.Data;
using Backend.Dotnet.Domain.Exceptions;
using static Backend.Dotnet.Application.DTOs.VehicleModelDto;

namespace Backend.Dotnet.Application.Services
{
    public class VehicleModelService : IVehicleModelService
    {
        private readonly IUnitOfWork _unitOfWork;

        public VehicleModelService(IUnitOfWork unitOfWork)
        {
            _unitOfWork = unitOfWork ?? throw new ArgumentNullException(nameof(unitOfWork));
        }

        public async Task<BaseResponseDto<VehicleModelResponse>> CreateAsync(CreateVehicleModelRequest request)
        {
            try
            {
                // Check if combination already exists
                var exists = await _unitOfWork.VehicleModels
                    .ExistsByBrandModelYearAsync(request.Brand, request.ModelName, request.Year);

                if (exists)
                {
                    return new BaseResponseDto<VehicleModelResponse>
                    {
                        IsSuccess = false,
                        Message = $"Vehicle model '{request.Brand} {request.ModelName} {request.Year}' already exists",
                        ErrorCode = "DUPLICATE_MODEL"
                    };
                }

                var model = request.ToEntity();
                await _unitOfWork.VehicleModels.AddAsync(model);
                await _unitOfWork.SaveChangesAsync();

                return new BaseResponseDto<VehicleModelResponse>
                {
                    IsSuccess = true,
                    Message = "Vehicle model created successfully",
                    Data = model.ToResponse()
                };
            }
            catch (BusinessRuleViolationException ex)
            {
                return new BaseResponseDto<VehicleModelResponse>
                {
                    IsSuccess = false,
                    Message = ex.Message,
                    ErrorCode = ex.ErrorCode
                };
            }
            catch (Exception ex)
            {
                return new BaseResponseDto<VehicleModelResponse>
                {
                    IsSuccess = false,
                    Message = "An error occurred while creating vehicle model",
                    ErrorCode = "INTERNAL_ERROR"
                };
            }
        }

        public async Task<BaseResponseDto<VehicleModelResponse>> GetByIdAsync(Guid id)
        {
            try
            {
                var model = await _unitOfWork.VehicleModels.GetByIdAsync(id);
                if (model == null)
                {
                    return new BaseResponseDto<VehicleModelResponse>
                    {
                        IsSuccess = false,
                        Message = $"Vehicle model with ID '{id}' not found",
                        ErrorCode = "NOT_FOUND"
                    };
                }

                return new BaseResponseDto<VehicleModelResponse>
                {
                    IsSuccess = true,
                    Message = "Vehicle model retrieved successfully",
                    Data = model.ToResponse()
                };
            }
            catch (Exception ex)
            {
                return new BaseResponseDto<VehicleModelResponse>
                {
                    IsSuccess = false,
                    Message = "An error occurred while retrieving vehicle model",
                    ErrorCode = "INTERNAL_ERROR"
                };
            }
        }

        public async Task<BaseResponseDto<IEnumerable<VehicleModelResponse>>> GetAllAsync()
        {
            try
            {
                var models = await _unitOfWork.VehicleModels.GetAllAsync();
                var response = models.Select(m => m.ToResponse()).ToList();

                return new BaseResponseDto<IEnumerable<VehicleModelResponse>>
                {
                    IsSuccess = true,
                    Message = "Vehicle models retrieved successfully",
                    Data = response
                };
            }
            catch (Exception ex)
            {
                return new BaseResponseDto<IEnumerable<VehicleModelResponse>>
                {
                    IsSuccess = false,
                    Message = "An error occurred while retrieving vehicle models",
                    ErrorCode = "INTERNAL_ERROR"
                };
            }
        }

        public async Task<BaseResponseDto<VehicleModelResponse>> GetByBrandModelYearAsync(string brand, string modelName, int year)
        {
            try
            {
                var model = await _unitOfWork.VehicleModels.GetByBrandModelYearAsync(brand, modelName, year);
                if (model == null)
                {
                    return new BaseResponseDto<VehicleModelResponse>
                    {
                        IsSuccess = false,
                        Message = $"Vehicle model '{brand} {modelName} {year}' not found",
                        ErrorCode = "NOT_FOUND"
                    };
                }

                return new BaseResponseDto<VehicleModelResponse>
                {
                    IsSuccess = true,
                    Message = "Vehicle model retrieved successfully",
                    Data = model.ToResponse()
                };
            }
            catch (Exception ex)
            {
                return new BaseResponseDto<VehicleModelResponse>
                {
                    IsSuccess = false,
                    Message = "An error occurred while retrieving vehicle model",
                    ErrorCode = "INTERNAL_ERROR"
                };
            }
        }

        public async Task<BaseResponseDto<IEnumerable<VehicleModelResponse>>> GetByBrandAsync(string brand)
        {
            try
            {
                var models = await _unitOfWork.VehicleModels.GetByBrandAsync(brand);
                if (!models.Any())
                {
                    return new BaseResponseDto<IEnumerable<VehicleModelResponse>>
                    {
                        IsSuccess = false,
                        Message = $"No vehicle models found for brand '{brand}'",
                        ErrorCode = "NOT_FOUND"
                    };
                }

                var response = models.Select(m => m.ToResponse()).ToList();

                return new BaseResponseDto<IEnumerable<VehicleModelResponse>>
                {
                    IsSuccess = true,
                    Message = "Vehicle models retrieved successfully",
                    Data = response
                };
            }
            catch (Exception ex)
            {
                return new BaseResponseDto<IEnumerable<VehicleModelResponse>>
                {
                    IsSuccess = false,
                    Message = "An error occurred while retrieving vehicle models",
                    ErrorCode = "INTERNAL_ERROR"
                };
            }
        }

        public async Task<BaseResponseDto<IEnumerable<VehicleModelResponse>>> GetByModelNameAsync(string modelName)
        {
            try
            {
                var models = await _unitOfWork.VehicleModels.GetByModelNameAsync(modelName);
                if (!models.Any())
                {
                    return new BaseResponseDto<IEnumerable<VehicleModelResponse>>
                    {
                        IsSuccess = false,
                        Message = $"No vehicle models found with name '{modelName}'",
                        ErrorCode = "NOT_FOUND"
                    };
                }

                var response = models.Select(m => m.ToResponse()).ToList();

                return new BaseResponseDto<IEnumerable<VehicleModelResponse>>
                {
                    IsSuccess = true,
                    Message = "Vehicle models retrieved successfully",
                    Data = response
                };
            }
            catch (Exception)
            {
                return new BaseResponseDto<IEnumerable<VehicleModelResponse>>
                {
                    IsSuccess = false,
                    Message = "An error occurred while retrieving vehicle models",
                    ErrorCode = "INTERNAL_ERROR"
                };
            }
        }

        public async Task<BaseResponseDto<IEnumerable<VehicleModelResponse>>> GetByYearAsync(int year)
        {
            try
            {
                var models = await _unitOfWork.VehicleModels.GetByYearAsync(year);
                if (!models.Any())
                {
                    return new BaseResponseDto<IEnumerable<VehicleModelResponse>>
                    {
                        IsSuccess = false,
                        Message = $"No vehicle models found for year {year}",
                        ErrorCode = "NOT_FOUND"
                    };
                }

                var response = models.Select(m => m.ToResponse()).ToList();

                return new BaseResponseDto<IEnumerable<VehicleModelResponse>>
                {
                    IsSuccess = true,
                    Message = "Vehicle models retrieved successfully",
                    Data = response
                };
            }
            catch (Exception)
            {
                return new BaseResponseDto<IEnumerable<VehicleModelResponse>>
                {
                    IsSuccess = false,
                    Message = "An error occurred while retrieving vehicle models",
                    ErrorCode = "INTERNAL_ERROR"
                };
            }
        }

        public async Task<BaseResponseDto<VehicleModelResponse>> UpdateAsync(Guid id, UpdateVehicleModelRequest request)
        {
            try
            {
                var model = await _unitOfWork.VehicleModels.GetByIdAsync(id);
                if (model == null)
                {
                    return new BaseResponseDto<VehicleModelResponse>
                    {
                        IsSuccess = false,
                        Message = $"Vehicle model with ID '{id}' not found",
                        ErrorCode = "NOT_FOUND"
                    };
                }

                // Check if combination already exists for another model
                var exists = await _unitOfWork.VehicleModels
                    .ExistsByBrandModelYearAsync(request.Brand, request.ModelName, request.Year, id);

                if (exists)
                {
                    return new BaseResponseDto<VehicleModelResponse>
                    {
                        IsSuccess = false,
                        Message = $"Vehicle model '{request.Brand} {request.ModelName} {request.Year}' already exists",
                        ErrorCode = "DUPLICATE_MODEL"
                    };
                }

                request.ApplyToEntity(model);
                _unitOfWork.VehicleModels.Update(model);
                await _unitOfWork.SaveChangesAsync();

                return new BaseResponseDto<VehicleModelResponse>
                {
                    IsSuccess = true,
                    Message = "Vehicle model updated successfully",
                    Data = model.ToResponse()
                };
            }
            catch (BusinessRuleViolationException ex)
            {
                return new BaseResponseDto<VehicleModelResponse>
                {
                    IsSuccess = false,
                    Message = ex.Message,
                    ErrorCode = ex.ErrorCode
                };
            }
            catch (Exception ex)
            {
                return new BaseResponseDto<VehicleModelResponse>
                {
                    IsSuccess = false,
                    Message = "An error occurred while updating vehicle model",
                    ErrorCode = "INTERNAL_ERROR"
                };
            }
        }

        // HARD DELETE OPERATION
        public async Task<BaseResponseDto> DeleteAsync(Guid id)
        {
            try
            {
                var model = await _unitOfWork.VehicleModels.GetByIdAsync(id);
                if (model == null)
                {
                    return new BaseResponseDto
                    {
                        IsSuccess = false,
                        Message = $"Vehicle model with ID '{id}' not found",
                        ErrorCode = "NOT_FOUND"
                    };
                }

                // Check if any active vehicles reference this model
                var hasVehicles = await _unitOfWork.VehicleModels.HasActiveVehiclesAsync(id);
                if (hasVehicles)
                {
                    var count = await _unitOfWork.VehicleModels.GetActiveVehicleCountAsync(id);
                    return new BaseResponseDto
                    {
                        IsSuccess = false,
                        Message = $"Cannot delete vehicle model. {count} active vehicle(s) are using this model.",
                        ErrorCode = "MODEL_IN_USE"
                    };
                }

                _unitOfWork.VehicleModels.Remove(model);
                await _unitOfWork.SaveChangesAsync();

                return new BaseResponseDto
                {
                    IsSuccess = true,
                    Message = "Vehicle model deleted successfully",
                    ErrorCode = "null"
                    //Data = true //should show which have been permanently deleted
                };
            }
            catch (Exception ex)
            {
                return new BaseResponseDto
                {
                    IsSuccess = false,
                    Message = "An error occurred while deleting vehicle model",
                    ErrorCode = "INTERNAL_ERROR"
                };
            }
        }
    }
}
