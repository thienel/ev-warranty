using Backend.Dotnet.Application.Constants;
using Backend.Dotnet.Application.DTOs;
using Backend.Dotnet.Application.Interfaces;
using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Mvc;
using static Backend.Dotnet.Application.DTOs.VehicleModelDto;

namespace Backend.Dotnet.API.Controllers
{
    [ApiController]
    [Route("vehicle-models")]
    [Produces("application/json")]
    public class VehicleModelsController : ControllerBase
    {
        private readonly IVehicleModelService _vehicleModelService;

        public VehicleModelsController(IVehicleModelService vehicleModelService)
        {
            _vehicleModelService = vehicleModelService;
        }

        /// <summary>
        /// Get vehicle models with optional filtering
        /// </summary>
        [HttpGet]
        [ProducesResponseType(typeof(BaseResponseDto<IEnumerable<VehicleModelResponse>>), StatusCodes.Status200OK)]
        [ProducesResponseType(typeof(BaseResponseDto), StatusCodes.Status400BadRequest)]
        public async Task<IActionResult> GetAll(
            [FromQuery] string? brand = null,
            [FromQuery] string? model = null,
            [FromQuery] int? year = null)
        {
            // Brand + model + year - relative
            if (!string.IsNullOrWhiteSpace(brand) &&
                !string.IsNullOrWhiteSpace(model) &&
                year.HasValue)
            {
                var result = await _vehicleModelService.GetByBrandModelYearAsync(brand, model, year.Value);
                return result.IsSuccess ? Ok(result) : NotFound(result);
            }

            // Brand - relative - list
            if (!string.IsNullOrWhiteSpace(brand))
            {
                var result = await _vehicleModelService.GetByBrandAsync(brand);
                return result.IsSuccess ? Ok(result) : NotFound(result);
            }

            // Model - relative - list
            if (!string.IsNullOrWhiteSpace(model))
            {
                var result = await _vehicleModelService.GetByModelNameAsync(model);
                return result.IsSuccess ? Ok(result) : NotFound(result);
            }

            // Year - exact - list
            if (year.HasValue)
            {
                var result = await _vehicleModelService.GetByYearAsync(year.Value);
                return result.IsSuccess ? Ok(result) : NotFound(result);
            }

            // No parameters - get all
            var allResult = await _vehicleModelService.GetAllAsync();
            return allResult.IsSuccess ? Ok(allResult) : BadRequest(allResult);
        }

        [HttpGet("{id}")]
        [ProducesResponseType(typeof(BaseResponseDto<VehicleModelResponse>), StatusCodes.Status200OK)]
        [ProducesResponseType(typeof(BaseResponseDto), StatusCodes.Status404NotFound)]
        public async Task<IActionResult> GetById(Guid id)
        {
            var result = await _vehicleModelService.GetByIdAsync(id);
            if (!result.IsSuccess)
                return NotFound(result);

            return Ok(result);
        }

        [HttpPost]
        [ProducesResponseType(typeof(BaseResponseDto<VehicleModelResponse>), StatusCodes.Status201Created)]
        [ProducesResponseType(typeof(BaseResponseDto), StatusCodes.Status400BadRequest)]
        [Authorize(Roles = SystemRoles.UserRoleAdmin + "," + SystemRoles.UserRoleEvmStaff)]
        public async Task<IActionResult> Create([FromBody] CreateVehicleModelRequest request)
        {
            if (!ModelState.IsValid)
                return BadRequest(ModelState);

            var result = await _vehicleModelService.CreateAsync(request);
            if (!result.IsSuccess)
                return BadRequest(result);

            return CreatedAtAction(nameof(GetById), new { id = result.Data!.Id }, result);
        }

        [HttpPut("{id}")]
        [ProducesResponseType(typeof(BaseResponseDto<VehicleModelResponse>), StatusCodes.Status200OK)]
        [ProducesResponseType(typeof(BaseResponseDto), StatusCodes.Status400BadRequest)]
        [ProducesResponseType(typeof(BaseResponseDto), StatusCodes.Status404NotFound)]
        [Authorize(Roles = SystemRoles.UserRoleAdmin + "," + SystemRoles.UserRoleEvmStaff)]
        public async Task<IActionResult> Update(Guid id, [FromBody] UpdateVehicleModelRequest request)
        {
            if (!ModelState.IsValid)
                return BadRequest(ModelState);

            var result = await _vehicleModelService.UpdateAsync(id, request);
            if (!result.IsSuccess)
                return result.ErrorCode == "NOT_FOUND" ? NotFound(result) : BadRequest(result);

            return Ok(result);
        }

        /// <summary>
        /// Hard delete a model (only if no active vehicles reference it)
        /// </summary>
        [HttpDelete("{id}")]
        [ProducesResponseType(typeof(BaseResponseDto<VehicleModelResponse>), StatusCodes.Status200OK)]
        [ProducesResponseType(typeof(BaseResponseDto), StatusCodes.Status404NotFound)]
        [ProducesResponseType(typeof(BaseResponseDto), StatusCodes.Status400BadRequest)]
        [Authorize(Roles = SystemRoles.UserRoleAdmin + "," + SystemRoles.UserRoleEvmStaff)]
        public async Task<IActionResult> Delete(Guid id)
        {
            var result = await _vehicleModelService.DeleteAsync(id);
            if (!result.IsSuccess)
                return result.ErrorCode == "NOT_FOUND" ? NotFound(result) : BadRequest(result);

            return Ok(result);
        }
    }
}
