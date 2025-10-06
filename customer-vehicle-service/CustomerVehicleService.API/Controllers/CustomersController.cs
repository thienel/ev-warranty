using CustomerVehicleService.Application.Interfaces;
using Microsoft.AspNetCore.Mvc;
using static CustomerVehicleService.Application.DTOs.CustomerDto;

namespace CustomerVehicleService.API.Controllers
{
    [ApiController]
    [Route("api/[controller]")]
    [Produces("application/json")]
    public class CustomersController : ControllerBase
    {
        private readonly ICustomerService _customerService;

        public CustomersController(ICustomerService customerService)
        {
            _customerService = customerService;
        }

        [HttpGet]
        public async Task<IActionResult> GetAll()
        {
            var result = await _customerService.GetAllAsync();
            if (!result.IsSuccess)
                return BadRequest(result);

            return Ok(result);
        }

        [HttpGet("{id}")]
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
        public async Task<IActionResult> Search([FromQuery] string searchTerm)
        {
            var result = await _customerService.SearchAsync(searchTerm);
            if (!result.IsSuccess)
                return BadRequest(result);

            return Ok(result);
        }

        [HttpPost]
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
        public async Task<IActionResult> Update(Guid id, [FromBody] UpdateCustomerRequest request)
        {
            if (!ModelState.IsValid)
                return BadRequest(ModelState);

            var result = await _customerService.UpdateAsync(id, request);
            if (!result.IsSuccess)
                return BadRequest(result);

            return Ok(result);
        }

        [HttpPatch("{id}/email")]
        public async Task<IActionResult> UpdateEmail(Guid id, [FromBody] string email)
        {
            if (string.IsNullOrWhiteSpace(email))
                return BadRequest(new { message = "Email is required" });

            var result = await _customerService.UpdateEmailAsync(id, email);
            if (!result.IsSuccess)
                return BadRequest(result);

            return Ok(result);
        }

        [HttpPatch("{id}/phone")]
        public async Task<IActionResult> UpdatePhoneNumber(Guid id, [FromBody] string phoneNumber)
        {
            if (string.IsNullOrWhiteSpace(phoneNumber))
                return BadRequest(new { message = "Phone number is required" });

            var result = await _customerService.UpdatePhoneNumberAsync(id, phoneNumber);
            if (!result.IsSuccess)
                return BadRequest(result);

            return Ok(result);
        }

        [HttpPatch("{id}/address")]
        public async Task<IActionResult> UpdateAddress(Guid id, [FromBody] string address)
        {
            var result = await _customerService.UpdateAddressAsync(id, address);
            if (!result.IsSuccess)
                return BadRequest(result);

            return Ok(result);
        }
    }

}
