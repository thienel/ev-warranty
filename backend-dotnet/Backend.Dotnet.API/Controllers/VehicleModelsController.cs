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
        [ProducesResponseType(typeof(BaseResponseDto<VehicleModelResponse>), StatusCodes.Status200OK)]
        [ProducesResponseType(typeof(BaseResponseDto), StatusCodes.Status400BadRequest)]
        public async Task<IActionResult> GetAll()
        {
            var result = await _vehicleModelService.GetAllAsync();
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
        /*
        [HttpGet("{id}/stats")]
        public async Task<IActionResult> GetWithStats(Guid id)
        {
            var result = await _vehicleModelService.GetWithStatsAsync(id);
            if (!result.IsSuccess)
                return NotFound(result);

            return Ok(result);
        }
        */
        [HttpGet("by-details")]
        [ProducesResponseType(typeof(BaseResponseDto<VehicleModelResponse>), StatusCodes.Status200OK)]
        [ProducesResponseType(typeof(BaseResponseDto), StatusCodes.Status400BadRequest)]
        [ProducesResponseType(typeof(BaseResponseDto), StatusCodes.Status404NotFound)]
        public async Task<IActionResult> GetByBrandModelYear([FromQuery] string brand, [FromQuery] string modelName, [FromQuery] int year)
        {
            if (string.IsNullOrWhiteSpace(brand) || string.IsNullOrWhiteSpace(modelName) || year <= 0)
                return BadRequest(new { message = "Brand, model name, and year are required" });

            var result = await _vehicleModelService.GetByBrandModelYearAsync(brand, modelName, year);
            if (!result.IsSuccess)
                return NotFound(result);

            return Ok(result);
        }

        [HttpGet("brand/{brand}")]
        [ProducesResponseType(typeof(BaseResponseDto<VehicleModelResponse>), StatusCodes.Status200OK)]
        [ProducesResponseType(typeof(BaseResponseDto), StatusCodes.Status400BadRequest)]
        public async Task<IActionResult> GetByBrand(string brand)
        {
            if (string.IsNullOrWhiteSpace(brand))
                return BadRequest(new { message = "Brand is required" });

            var result = await _vehicleModelService.GetByBrandAsync(brand);
            if (!result.IsSuccess)
                return BadRequest(result);

            return Ok(result);
        }

        [HttpGet("brands")]
        [ProducesResponseType(typeof(BaseResponseDto<VehicleModelResponse>), StatusCodes.Status200OK)]
        [ProducesResponseType(typeof(BaseResponseDto), StatusCodes.Status400BadRequest)]
        public async Task<IActionResult> GetAllBrands()
        {
            var result = await _vehicleModelService.GetAllBrandsAsync();
            if (!result.IsSuccess)
                return BadRequest(result);

            return Ok(result);
        }

        [HttpGet("search")]
        [ProducesResponseType(typeof(BaseResponseDto<VehicleModelResponse>), StatusCodes.Status200OK)]
        [ProducesResponseType(typeof(BaseResponseDto), StatusCodes.Status400BadRequest)]
        public async Task<IActionResult> Search([FromQuery] string searchTerm)
        {
            var result = await _vehicleModelService.SearchAsync(searchTerm);
            if (!result.IsSuccess)
                return BadRequest(result);

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

        // Hard delete Endpoints
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
