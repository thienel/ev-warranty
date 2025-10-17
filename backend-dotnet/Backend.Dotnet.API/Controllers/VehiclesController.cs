using Backend.Dotnet.Application.Interfaces;
using Microsoft.AspNetCore.Mvc;
using static Backend.Dotnet.Application.DTOs.VehicleDto;

namespace Backend.Dotnet.API.Controllers
{

    [ApiController]
    [Route("[controller]")]
    [Produces("application/json")]
    public class VehiclesController : ControllerBase
    {
        private readonly IVehicleService _vehicleService;

        public VehiclesController(IVehicleService vehicleService)
        {
            _vehicleService = vehicleService;
        }

        [HttpGet]
        public async Task<IActionResult> GetAll()
        {
            var result = await _vehicleService.GetAllAsync();
            if (!result.IsSuccess)
                return BadRequest(result);

            return Ok(result);
        }

        [HttpGet("{id}")]
        public async Task<IActionResult> GetById(Guid id)
        {
            var result = await _vehicleService.GetByIdAsync(id);
            if (!result.IsSuccess)
                return NotFound(result);

            return Ok(result);
        }
        /*
        [HttpGet("{id}/details")]
        public async Task<IActionResult> GetDetail(Guid id)
        {
            var result = await _vehicleService.GetDetailAsync(id);
            if (!result.IsSuccess)
                return NotFound(result);

            return Ok(result);
        }
        */
        [HttpGet("by-vin/{vin}")]
        public async Task<IActionResult> GetByVin(string vin)
        {
            if (string.IsNullOrWhiteSpace(vin))
                return BadRequest(new { message = "VIN is required" });

            var result = await _vehicleService.GetByVinAsync(vin);
            if (!result.IsSuccess)
                return NotFound(result);

            return Ok(result);
        }

        [HttpGet("by-plate/{licensePlate}")]
        public async Task<IActionResult> GetByLicensePlate(string licensePlate)
        {
            if (string.IsNullOrWhiteSpace(licensePlate))
                return BadRequest(new { message = "License plate is required" });

            var result = await _vehicleService.GetByLicensePlateAsync(licensePlate);
            if (!result.IsSuccess)
                return NotFound(result);

            return Ok(result);
        }

        [HttpGet("customer/{customerId}")]
        public async Task<IActionResult> GetByCustomerId(Guid customerId)
        {
            var result = await _vehicleService.GetByCustomerIdAsync(customerId);
            if (!result.IsSuccess)
                return BadRequest(result);

            return Ok(result);
        }

        [HttpGet("model/{modelId}")]
        public async Task<IActionResult> GetByModelId(Guid modelId)
        {
            var result = await _vehicleService.GetByModelIdAsync(modelId);
            if (!result.IsSuccess)
                return BadRequest(result);

            return Ok(result);
        }

        [HttpGet("search")]
        public async Task<IActionResult> Search([FromQuery] string searchTerm)
        {
            var result = await _vehicleService.SearchAsync(searchTerm);
            if (!result.IsSuccess)
                return BadRequest(result);

            return Ok(result);
        }

        [HttpPost]
        public async Task<IActionResult> Create([FromBody] CreateVehicleRequest request)
        {
            if (!ModelState.IsValid)
                return BadRequest(ModelState);

            var result = await _vehicleService.CreateAsync(request);
            if (!result.IsSuccess)
                return BadRequest(result);

            return CreatedAtAction(nameof(GetById), new { id = result.Data.Id }, result);
        }

        [HttpPut("{id}")]
        public async Task<IActionResult> Update(Guid id, [FromBody] UpdateVehicleRequest request)
        {
            if (!ModelState.IsValid)
                return BadRequest(ModelState);

            var result = await _vehicleService.UpdateAsync(id, request);
            if (!result.IsSuccess)
                return BadRequest(result);

            return Ok(result);
        }

        [HttpPatch("{id}/license-plate")]
        public async Task<IActionResult> UpdateLicensePlate(Guid id, [FromBody] UpdateLicensePlateCommand command)
        {
            if (!ModelState.IsValid)
                return BadRequest(ModelState);

            var result = await _vehicleService.UpdateLicensePlateAsync(id, command);
            if (!result.IsSuccess)
                return BadRequest(result);

            return Ok(result);
        }

        [HttpPatch("{id}/transfer")]
        public async Task<IActionResult> TransferOwnership(Guid id, [FromBody] TransferVehicleCommand command)
        {
            if (!ModelState.IsValid)
                return BadRequest(ModelState);

            var result = await _vehicleService.TransferOwnershipAsync(id, command);
            if (!result.IsSuccess)
                return BadRequest(result);

            return Ok(result);
        }

        // Soft delete Endpoints
        [HttpDelete("{id}")]
        public async Task<IActionResult> SoftDelete(Guid id)
        {
            var result = await _vehicleService.SoftDeleteAsync(id);
            if (!result.IsSuccess)
                return BadRequest(result);

            return Ok(result);
        }

        [HttpPost("{id}/restore")]
        public async Task<IActionResult> Restore(Guid id)
        {
            var result = await _vehicleService.RestoreAsync(id);
            if (!result.IsSuccess)
                return BadRequest(result);

            return Ok(result);
        }
    }
}
