using Backend.Dotnet.Application.DTOs;
using Backend.Dotnet.Application.Interfaces;
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

        [HttpGet]
        [ProducesResponseType(typeof(BaseResponseDto<IEnumerable<VehicleModelResponse>>), StatusCodes.Status200OK)]
        [ProducesResponseType(typeof(BaseResponseDto), StatusCodes.Status400BadRequest)]
        public async Task<IActionResult> GetAll(
        [FromQuery] string brand = null,
        [FromQuery] string modelName = null,
        [FromQuery] int? year = null,
        [FromQuery] string search = null)
        {
            // Filter by brand + model + year (returns single or not found)
            if (!string.IsNullOrWhiteSpace(brand) &&
                !string.IsNullOrWhiteSpace(modelName) &&
                year.HasValue)
            {
                var result = await _vehicleModelService.GetByBrandModelYearAsync(brand, modelName, year.Value);
                return result.IsSuccess ? Ok(result) : NotFound(result);
            }

            // Filter by brand only (returns list)
            if (!string.IsNullOrWhiteSpace(brand))
            {
                var result = await _vehicleModelService.GetByBrandAsync(brand);
                return result.IsSuccess ? Ok(result) : BadRequest(result);
            }

            // Search/filter by term (returns list)
            if (!string.IsNullOrWhiteSpace(search))
            {
                var result = await _vehicleModelService.SearchAsync(search);
                return result.IsSuccess ? Ok(result) : BadRequest(result);
            }

            // No parameters = get all
            var allResult = await _vehicleModelService.GetAllAsync();
            return allResult.IsSuccess ? Ok(allResult) : BadRequest(allResult);
        }

        [HttpGet("brands")]
        [ProducesResponseType(typeof(BaseResponseDto<IEnumerable<string>>), StatusCodes.Status200OK)]
        [ProducesResponseType(typeof(BaseResponseDto), StatusCodes.Status400BadRequest)]
        public async Task<IActionResult> GetAllBrands()
        {
            var result = await _vehicleModelService.GetAllBrandsAsync();
            if (!result.IsSuccess)
                return BadRequest(result);

            return Ok(result);
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
        public async Task<IActionResult> Create([FromBody] CreateVehicleModelRequest request)
        {
            if (!ModelState.IsValid)
                return BadRequest(ModelState);

            var result = await _vehicleModelService.CreateAsync(request);
            if (!result.IsSuccess)
                return BadRequest(result);

            return CreatedAtAction(nameof(GetById), new { id = result.Data.Id }, result);
        }

        [HttpPut("{id}")]
        [ProducesResponseType(typeof(BaseResponseDto<VehicleModelResponse>), StatusCodes.Status200OK)]
        [ProducesResponseType(typeof(BaseResponseDto), StatusCodes.Status400BadRequest)]
        public async Task<IActionResult> Update(Guid id, [FromBody] UpdateVehicleModelRequest request)
        {
            if (!ModelState.IsValid)
                return BadRequest(ModelState);

            var result = await _vehicleModelService.UpdateAsync(id, request);
            if (!result.IsSuccess)
                return BadRequest(result);

            return Ok(result);
        }

        [HttpDelete("{id}")]
        [ProducesResponseType(typeof(BaseResponseDto<VehicleModelResponse>), StatusCodes.Status200OK)]
        [ProducesResponseType(typeof(BaseResponseDto), StatusCodes.Status400BadRequest)]
        public async Task<IActionResult> Delete(Guid id)
        {
            var result = await _vehicleModelService.DeleteAsync(id);
            if (!result.IsSuccess)
                return BadRequest(result);

            return Ok(result);
        }
    }
}
