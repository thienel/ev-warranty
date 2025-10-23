using Backend.Dotnet.Application.DTOs;
using Backend.Dotnet.Application.Interfaces;
using Microsoft.AspNetCore.Mvc;
using static Backend.Dotnet.Application.DTOs.VehicleDto;

namespace Backend.Dotnet.API.Controllers
{

    [ApiController]
    [Route("vehicles")]
    [Produces("application/json")]
    public class VehiclesController : ControllerBase
    {
        private readonly IVehicleService _vehicleService;

        public VehiclesController(IVehicleService vehicleService)
        {
            _vehicleService = vehicleService;
        }

        [HttpGet]
        [ProducesResponseType(typeof(BaseResponseDto<IEnumerable<VehicleResponse>>), StatusCodes.Status200OK)]
        [ProducesResponseType(typeof(BaseResponseDto), StatusCodes.Status400BadRequest)]
        public async Task<IActionResult> GetAll(
        [FromQuery] string vin = null,
        [FromQuery] string licensePlate = null,
        [FromQuery] Guid? customerId = null,
        [FromQuery] Guid? modelId = null,
        [FromQuery] string search = null)
        {
            if (!string.IsNullOrWhiteSpace(vin))
            {
                var result = await _vehicleService.GetByVinAsync(vin);
                return result.IsSuccess ? Ok(result) : NotFound(result);
            }

            if (!string.IsNullOrWhiteSpace(licensePlate))
            {
                var result = await _vehicleService.GetByLicensePlateAsync(licensePlate);
                return result.IsSuccess ? Ok(result) : NotFound(result);
            }

            if (customerId.HasValue)
            {
                var result = await _vehicleService.GetByCustomerIdAsync(customerId.Value);
                return result.IsSuccess ? Ok(result) : BadRequest(result);
            }

            if (modelId.HasValue)
            {
                var result = await _vehicleService.GetByModelIdAsync(modelId.Value);
                return result.IsSuccess ? Ok(result) : BadRequest(result);
            }

            if (!string.IsNullOrWhiteSpace(search))
            {
                var result = await _vehicleService.SearchAsync(search);
                return result.IsSuccess ? Ok(result) : BadRequest(result);
            }

            var allResult = await _vehicleService.GetAllAsync();
            return allResult.IsSuccess ? Ok(allResult) : BadRequest(allResult);
        }

        [HttpGet("{id}")]
        [ProducesResponseType(typeof(BaseResponseDto<VehicleResponse>), StatusCodes.Status200OK)]
        [ProducesResponseType(typeof(BaseResponseDto), StatusCodes.Status404NotFound)]
        public async Task<IActionResult> GetById(Guid id)
        {
            var result = await _vehicleService.GetByIdAsync(id);
            return result.IsSuccess ? Ok(result) : NotFound(result);
        }

        [HttpPost]
        [ProducesResponseType(typeof(BaseResponseDto<VehicleResponse>), StatusCodes.Status201Created)]
        [ProducesResponseType(typeof(BaseResponseDto), StatusCodes.Status400BadRequest)]
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
        [ProducesResponseType(typeof(BaseResponseDto<VehicleResponse>), StatusCodes.Status200OK)]
        [ProducesResponseType(typeof(BaseResponseDto), StatusCodes.Status400BadRequest)]
        public async Task<IActionResult> Update(Guid id, [FromBody] UpdateVehicleRequest request)
        {
            if (!ModelState.IsValid)
                return BadRequest(ModelState);

            var result = await _vehicleService.UpdateAsync(id, request);
            if (!result.IsSuccess)
                return BadRequest(result);

            return Ok(result);
        }

        [HttpDelete("{id}")]
        [ProducesResponseType(typeof(BaseResponseDto<VehicleResponse>), StatusCodes.Status200OK)]
        [ProducesResponseType(typeof(BaseResponseDto), StatusCodes.Status400BadRequest)]
        public async Task<IActionResult> SoftDelete(Guid id)
        {
            var result = await _vehicleService.SoftDeleteAsync(id);
            if (!result.IsSuccess)
                return BadRequest(result);

            return Ok(result);
        }

        [HttpPost("{id}/restore")]
        [ProducesResponseType(typeof(BaseResponseDto<VehicleResponse>), StatusCodes.Status200OK)]
        [ProducesResponseType(typeof(BaseResponseDto), StatusCodes.Status400BadRequest)]
        public async Task<IActionResult> Restore(Guid id)
        {
            var result = await _vehicleService.RestoreAsync(id);
            if (!result.IsSuccess)
                return BadRequest(result);

            return Ok(result);
        }
    }
}
