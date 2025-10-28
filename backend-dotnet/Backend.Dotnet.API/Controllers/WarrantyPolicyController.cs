using Backend.Dotnet.Application.DTOs;
using Backend.Dotnet.Application.Interfaces;
using Microsoft.AspNetCore.Http;
using Microsoft.AspNetCore.Mvc;
using static Backend.Dotnet.Application.DTOs.WarrantyPolicyDto;

namespace Backend.Dotnet.API.Controllers
{
    [ApiController]
    [Route("warranty-policies")]
    [Produces("application/json")]
    public class WarrantyPoliciesController : ControllerBase
    {
        private readonly IWarrantyPolicyService _warrantyPolicyService;

        public WarrantyPoliciesController(IWarrantyPolicyService warrantyPolicyService)
        {
            _warrantyPolicyService = warrantyPolicyService;
        }

        [HttpGet]
        [ProducesResponseType(typeof(BaseResponseDto<IEnumerable<WarrantyPolicyResponse>>), StatusCodes.Status200OK)]
        [ProducesResponseType(typeof(BaseResponseDto), StatusCodes.Status400BadRequest)]
        public async Task<IActionResult> GetAll(
            [FromQuery] string status = null,
            [FromQuery] string policyName = null)
        {
            // Status - exact
            if (!string.IsNullOrWhiteSpace(status))
            {
                var result = await _warrantyPolicyService.GetByStatusAsync(status);
                return result.IsSuccess ? Ok(result) : BadRequest(result);
            }

            // Policy name - exact 
            if (!string.IsNullOrWhiteSpace(policyName))
            {
                var result = await _warrantyPolicyService.GetByPolicyNameAsync(policyName);
                return result.IsSuccess ? Ok(result) : NotFound(result);
            }

            // No parameters - get all
            var allResult = await _warrantyPolicyService.GetAllAsync();
            return allResult.IsSuccess ? Ok(allResult) : BadRequest(allResult);
        }

        [HttpGet("{id}")]
        [ProducesResponseType(typeof(BaseResponseDto<WarrantyPolicyResponse>), StatusCodes.Status200OK)]
        [ProducesResponseType(typeof(BaseResponseDto), StatusCodes.Status404NotFound)]
        public async Task<IActionResult> GetById(Guid id)
        {
            var result = await _warrantyPolicyService.GetByIdAsync(id);
            if (!result.IsSuccess)
                return NotFound(result);

            return Ok(result);
        }

        [HttpGet("{id}/details")]
        [ProducesResponseType(typeof(BaseResponseDto<WarrantyPolicyWithDetailsResponse>), StatusCodes.Status200OK)]
        [ProducesResponseType(typeof(BaseResponseDto), StatusCodes.Status404NotFound)]
        public async Task<IActionResult> GetWithDetails(Guid id)
        {
            var result = await _warrantyPolicyService.GetWithDetailsAsync(id);
            if (!result.IsSuccess)
                return NotFound(result);

            return Ok(result);
        }

        [HttpPost]
        [ProducesResponseType(typeof(BaseResponseDto<WarrantyPolicyResponse>), StatusCodes.Status201Created)]
        [ProducesResponseType(typeof(BaseResponseDto), StatusCodes.Status400BadRequest)]
        public async Task<IActionResult> Create([FromBody] CreateWarrantyPolicyRequest request)
        {
            if (!ModelState.IsValid)
                return BadRequest(ModelState);

            var result = await _warrantyPolicyService.CreateAsync(request);
            if (!result.IsSuccess)
                return BadRequest(result);

            return CreatedAtAction(nameof(GetById), new { id = result.Data.Id }, result);
        }

        [HttpPut("{id}")]
        [ProducesResponseType(typeof(BaseResponseDto<WarrantyPolicyResponse>), StatusCodes.Status200OK)]
        [ProducesResponseType(typeof(BaseResponseDto), StatusCodes.Status400BadRequest)]
        [ProducesResponseType(typeof(BaseResponseDto), StatusCodes.Status404NotFound)]
        public async Task<IActionResult> Update(Guid id, [FromBody] UpdateWarrantyPolicyRequest request)
        {
            if (!ModelState.IsValid)
                return BadRequest(ModelState);

            var result = await _warrantyPolicyService.UpdateAsync(id, request);
            if (!result.IsSuccess)
                return result.ErrorCode == "NOT_FOUND" ? NotFound(result) : BadRequest(result);

            return Ok(result);
        }

        [HttpPut("{id}/status")]
        [ProducesResponseType(typeof(BaseResponseDto<WarrantyPolicyResponse>), StatusCodes.Status200OK)]
        [ProducesResponseType(typeof(BaseResponseDto), StatusCodes.Status400BadRequest)]
        [ProducesResponseType(typeof(BaseResponseDto), StatusCodes.Status404NotFound)]
        public async Task<IActionResult> ChangeStatus(Guid id, [FromBody] ChangeStatusRequest request)
        {
            if (!ModelState.IsValid)
                return BadRequest(ModelState);

            var result = await _warrantyPolicyService.ChangeStatusAsync(id, request);
            if (!result.IsSuccess)
                return result.ErrorCode == "NOT_FOUND" ? NotFound(result) : BadRequest(result);

            return Ok(result);
        }

        [HttpDelete("{id}")]
        [ProducesResponseType(typeof(BaseResponseDto), StatusCodes.Status200OK)]
        [ProducesResponseType(typeof(BaseResponseDto), StatusCodes.Status404NotFound)]
        [ProducesResponseType(typeof(BaseResponseDto), StatusCodes.Status400BadRequest)]
        public async Task<IActionResult> Delete(Guid id)
        {
            var result = await _warrantyPolicyService.DeleteAsync(id);
            if (!result.IsSuccess)
                return result.ErrorCode == "NOT_FOUND" ? NotFound(result) : BadRequest(result);

            return Ok(result);
        }
    }
}
