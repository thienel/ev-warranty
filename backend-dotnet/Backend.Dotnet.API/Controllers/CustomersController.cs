using Backend.Dotnet.Application.Constants;
using Backend.Dotnet.Application.DTOs;
using Backend.Dotnet.Application.Interfaces;
using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Mvc;
using static Backend.Dotnet.Application.DTOs.CustomerDto;

namespace Backend.Dotnet.API.Controllers
{
    [ApiController]
    [Route("customers")]
    [Produces("application/json")]
    //[Authorize(Roles = SystemRoles.UserRoleScStaff)]
    public class CustomersController : ControllerBase
    {
        private readonly ICustomerService _customerService;

        public CustomersController(ICustomerService customerService)
        {
            _customerService = customerService;
        }

        /// <summary>
        /// Get customers with optional filtering
        /// </summary>
        [HttpGet]
        //[Authorize(Roles = SystemRoles.UserRoleScTechnician)]
        [ProducesResponseType(typeof(BaseResponseDto<IEnumerable<CustomerResponse>>), StatusCodes.Status200OK)]
        [ProducesResponseType(typeof(BaseResponseDto), StatusCodes.Status400BadRequest)]
        public async Task<IActionResult> GetAll(
            [FromQuery] string name = null,
            [FromQuery] string phone = null,
            [FromQuery] string email = null)
        {

            // Email - absolute
            if (!string.IsNullOrWhiteSpace(email))
            {
                var result = await _customerService.GetByEmailAsync(email);
                return result.IsSuccess ? Ok(result) : NotFound(result);
            }

            // Phone - absolute
            if (!string.IsNullOrWhiteSpace(phone))
            {
                var result = await _customerService.GetByPhoneAsync(phone);
                return result.IsSuccess ? Ok(result) : NotFound(result);
            }

            // Name - absolute - first last combined
            if (!string.IsNullOrWhiteSpace(name))
            {
                var result = await _customerService.GetByNameAsync(name);
                return result.IsSuccess ? Ok(result) : NotFound(result);
            }

            // No parameters - get all
            var allResult = await _customerService.GetAllAsync();
            return allResult.IsSuccess ? Ok(allResult) : BadRequest(allResult);
        }

        [HttpGet("{id}")]
        [ProducesResponseType(typeof(BaseResponseDto<CustomerResponse>), StatusCodes.Status200OK)]
        [ProducesResponseType(typeof(BaseResponseDto), StatusCodes.Status404NotFound)]
        public async Task<IActionResult> GetById(Guid id)
        {
            var result = await _customerService.GetByIdAsync(id);
            if (!result.IsSuccess)
                return NotFound(result);

            return Ok(result);
        }

        [HttpPost]
        [ProducesResponseType(typeof(BaseResponseDto<CustomerResponse>), StatusCodes.Status201Created)]
        [ProducesResponseType(typeof(BaseResponseDto), StatusCodes.Status400BadRequest)]
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
        [ProducesResponseType(typeof(BaseResponseDto), StatusCodes.Status400BadRequest)]
        [ProducesResponseType(typeof(BaseResponseDto), StatusCodes.Status404NotFound)]
        public async Task<IActionResult> Update(Guid id, [FromBody] UpdateCustomerRequest request)
        {
            if (!ModelState.IsValid)
                return BadRequest(ModelState);

            var result = await _customerService.UpdateAsync(id, request);
            if (!result.IsSuccess)
                return result.ErrorCode == "NOT_FOUND" ? NotFound(result) : BadRequest(result);

            return Ok(result);
        }


        /// <summary>
        /// Soft delete a customer
        /// </summary>
        [HttpDelete("{id}")]
        [ProducesResponseType(typeof(BaseResponseDto<CustomerResponse>), StatusCodes.Status200OK)]
        [ProducesResponseType(typeof(BaseResponseDto), StatusCodes.Status404NotFound)]
        [ProducesResponseType(typeof(BaseResponseDto), StatusCodes.Status400BadRequest)]
        public async Task<IActionResult> SoftDelete(Guid id)
        {
            var result = await _customerService.SoftDeleteAsync(id);
            if (!result.IsSuccess)
                return result.ErrorCode == "NOT_FOUND" ? NotFound(result) : BadRequest(result);

            return Ok(result);
        }

        [HttpPost("{id}/restore")]
        [ProducesResponseType(typeof(BaseResponseDto<CustomerResponse>), StatusCodes.Status200OK)]
        [ProducesResponseType(typeof(BaseResponseDto), StatusCodes.Status404NotFound)]
        [ProducesResponseType(typeof(BaseResponseDto), StatusCodes.Status400BadRequest)]
        public async Task<IActionResult> Restore(Guid id)
        {
            var result = await _customerService.RestoreAsync(id);
            if (!result.IsSuccess)
                return result.ErrorCode == "NOT_FOUND" ? NotFound(result) : BadRequest(result);

            return Ok(result);
        }

    }
}
