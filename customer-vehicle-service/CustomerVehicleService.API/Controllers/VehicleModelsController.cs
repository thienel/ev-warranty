using CustomerVehicleService.Application.Interfaces;
using Microsoft.AspNetCore.Mvc;
using static CustomerVehicleService.Application.DTOs.VehicleModelDto;

namespace CustomerVehicleService.API.Controllers
{
    [ApiController]
    [Route("[controller]")]
    [Produces("application/json")]
    public class VehicleModelsController : ControllerBase
    {
        private readonly IVehicleModelService _vehicleModelService;

        public VehicleModelsController(IVehicleModelService vehicleModelService)
        {
            _vehicleModelService = vehicleModelService;
        }

        [HttpGet]
        public async Task<IActionResult> GetAll()
        {
            var result = await _vehicleModelService.GetAllAsync();
            if (!result.IsSuccess)
                return BadRequest(result);

            return Ok(result);
        }

        [HttpGet("{id}")]
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
        public async Task<IActionResult> GetAllBrands()
        {
            var result = await _vehicleModelService.GetAllBrandsAsync();
            if (!result.IsSuccess)
                return BadRequest(result);

            return Ok(result);
        }

        [HttpGet("search")]
        public async Task<IActionResult> Search([FromQuery] string searchTerm)
        {
            var result = await _vehicleModelService.SearchAsync(searchTerm);
            if (!result.IsSuccess)
                return BadRequest(result);

            return Ok(result);
        }

        [HttpPost]
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
        public async Task<IActionResult> Delete(Guid id)
        {
            var result = await _vehicleModelService.DeleteAsync(id);
            if (!result.IsSuccess)
                return BadRequest(result);

            return Ok(result);
        }
    }
}
