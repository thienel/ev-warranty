using Backend.Dotnet.Application.DTOs;
using Backend.Dotnet.Application.Interfaces;
using Microsoft.AspNetCore.Mvc;
using static Backend.Dotnet.Application.DTOs.CustomerDto;

namespace Backend.Dotnet.API.Controllers
{
    [ApiController]
    [Route("[controller]")]
    [Produces("application/json")]
    public class CustomersController : ControllerBase
    {
        private readonly ICustomerService _customerService;

        public CustomersController(ICustomerService customerService)
        {
            _customerService = customerService;
        }

        [HttpGet]
        [ProducesResponseType(StatusCodes.Status200OK)]
        [ProducesResponseType(StatusCodes.Status400BadRequest)]
        public async Task<IActionResult> GetAll()
        {
            var result = await _customerService.GetAllAsync();
            if (!result.IsSuccess)
                return BadRequest(result);

            return Ok(result);
        }

        [HttpGet("{id}")]
        [ProducesResponseType(typeof(BaseResponseDto<CustomerResponse>), StatusCodes.Status200OK)]
        [ProducesResponseType(typeof(BaseResponseDto<CustomerResponse>), StatusCodes.Status404NotFound)]
        public async Task<IActionResult> GetById(Guid id)
        {
            var result = await _customerService.GetByIdAsync(id);
            if (!result.IsSuccess)
                return NotFound(result);

            return Ok(result);
        }

        /*
        [HttpGet("{id}/vehicles")]
        public async Task<IActionResult> GetWithVehicles(Guid id)
        {
            var result = await _customerService.GetWithVehiclesAsync(id);
            if (!result.IsSuccess)
                return NotFound(result);

            return Ok(result);
        }
        */
        [HttpGet("by-email")]
        [ProducesResponseType(typeof(BaseResponseDto<CustomerResponse>), StatusCodes.Status200OK)]
        [ProducesResponseType(typeof(BaseResponseDto<CustomerResponse>), StatusCodes.Status400BadRequest)]
        [ProducesResponseType(typeof(BaseResponseDto<CustomerResponse>), StatusCodes.Status404NotFound)]
        public async Task<IActionResult> GetByEmail([FromQuery] string email)
        {
            if (string.IsNullOrWhiteSpace(email))
                return BadRequest(new { message = "Email is required" });

            var result = await _customerService.GetByEmailAsync(email);
            if (!result.IsSuccess)
                return NotFound(result);

            return Ok(result);
        }

        [HttpGet("search")]
        [ProducesResponseType(typeof(BaseResponseDto<IEnumerable<CustomerResponse>>), StatusCodes.Status200OK)]
        [ProducesResponseType(typeof(BaseResponseDto<IEnumerable<CustomerResponse>>), StatusCodes.Status400BadRequest)]
        public async Task<IActionResult> Search([FromQuery] string searchTerm)
        {
            var result = await _customerService.SearchAsync(searchTerm);
            if (!result.IsSuccess)
                return BadRequest(result);

            return Ok(result);
        }

        [HttpPost]
        [ProducesResponseType(StatusCodes.Status201Created)]
        [ProducesResponseType(StatusCodes.Status400BadRequest)]
        public async Task<IActionResult> Create([FromBody] CreateCustomerRequest request)
        {
            if (!ModelState.IsValid)
                return BadRequest(ModelState);

            var result = await _customerService.CreateAsync(request);
            if (!result.IsSuccess)
                return BadRequest(result);

            return CreatedAtAction(nameof(GetById), new { id = result.Data.Id }, result);
        }

        [HttpPut("{id}")]
        [ProducesResponseType(typeof(BaseResponseDto<CustomerResponse>), StatusCodes.Status200OK)]
        [ProducesResponseType(typeof(BaseResponseDto<CustomerResponse>), StatusCodes.Status400BadRequest)]
        [ProducesResponseType(typeof(BaseResponseDto<CustomerResponse>), StatusCodes.Status404NotFound)]
        public async Task<IActionResult> Update(Guid id, [FromBody] UpdateCustomerRequest request)
        {
            if (!ModelState.IsValid)
                return BadRequest(ModelState);

            var result = await _customerService.UpdateAsync(id, request);
            if (!result.IsSuccess)
                return BadRequest(result);

            return Ok(result);
        }

        //[HttpPut("{id}/email")]
        //public async Task<IActionResult> UpdateEmail(Guid id, [FromBody] string email)
        //{
        //    if (string.IsNullOrWhiteSpace(email))
        //        return BadRequest(new { message = "Email is required" });

        //    var result = await _customerService.UpdateEmailAsync(id, email);
        //    if (!result.IsSuccess)
        //        return BadRequest(result);

        //    return Ok(result);
        //}

        //[HttpPatch("{id}/phone")]
        //public async Task<IActionResult> UpdatePhoneNumber(Guid id, [FromBody] string phoneNumber)
        //{
        //    if (string.IsNullOrWhiteSpace(phoneNumber))
        //        return BadRequest(new { message = "Phone number is required" });

        //    var result = await _customerService.UpdatePhoneNumberAsync(id, phoneNumber);
        //    if (!result.IsSuccess)
        //        return BadRequest(result);

        //    return Ok(result);
        //}

        //[HttpPatch("{id}/address")]
        //public async Task<IActionResult> UpdateAddress(Guid id, [FromBody] string address)
        //{
        //    var result = await _customerService.UpdateAddressAsync(id, address);
        //    if (!result.IsSuccess)
        //        return BadRequest(result);

        //    return Ok(result);
        //}

        // Soft delete Endpoints
        [HttpDelete("{id}")]
        [ProducesResponseType(StatusCodes.Status200OK)]
        [ProducesResponseType(StatusCodes.Status400BadRequest)]
        public async Task<IActionResult> SoftDelete(Guid id)
        {
            var result = await _customerService.SoftDeleteAsync(id);
            if (!result.IsSuccess)
                return BadRequest(result);

            return Ok(result);
        }

        [HttpPost("{id}/restore")]
        [ProducesResponseType(StatusCodes.Status200OK)]
        [ProducesResponseType(StatusCodes.Status400BadRequest)]
        public async Task<IActionResult> Restore(Guid id)
        {
            var result = await _customerService.RestoreAsync(id);
            if (!result.IsSuccess)
                return BadRequest(result);

            return Ok(result);
        }

    }
}
