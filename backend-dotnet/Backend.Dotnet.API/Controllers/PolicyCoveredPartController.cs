using Backend.Dotnet.Application.Constants;
using Backend.Dotnet.Application.DTOs;
using Backend.Dotnet.Application.Interfaces;
using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Mvc;
using static Backend.Dotnet.Application.DTOs.PolicyCoveragePartDto;

namespace Backend.Dotnet.API.Controllers
{
    [ApiController]
    [Route("policy-coverage-parts")]
    [Produces("application/json")]
    public class PolicyCoveragePartsController : ControllerBase
    {
        private readonly IPolicyCoveragePartService _policyCoveragePartService;

        public PolicyCoveragePartsController(IPolicyCoveragePartService policyCoveragePartService)
        {
            _policyCoveragePartService = policyCoveragePartService;
        }

        [HttpGet]
        [ProducesResponseType(typeof(BaseResponseDto<IEnumerable<PolicyCoveragePartResponse>>), StatusCodes.Status200OK)]
        [ProducesResponseType(typeof(BaseResponseDto), StatusCodes.Status404NotFound)]
        [ProducesResponseType(typeof(BaseResponseDto), StatusCodes.Status400BadRequest)]
        public async Task<IActionResult> GetAll(
            [FromQuery] Guid? policyId = null,
            [FromQuery] Guid? partCategoryId = null)
        {
            // Filter by policy and category - exact match
            if (policyId.HasValue && partCategoryId.HasValue)
            {
                var result = await _policyCoveragePartService
                    .GetByPolicyAndCategoryAsync(policyId.Value, partCategoryId.Value);

                if (!result.IsSuccess)
                    return NotFound(result);

                return Ok(new BaseResponseDto<IEnumerable<PolicyCoveragePartResponse>>
                {
                    IsSuccess = true,
                    Message = result.Message,
                    Data = new List<PolicyCoveragePartResponse> { result.Data }
                });
            }

            // Filter by policy - exact match
            if (policyId.HasValue)
            {
                var result = await _policyCoveragePartService.GetByPolicyIdAsync(policyId.Value);
                return result.IsSuccess ? Ok(result) : NotFound(result);
            }

            // Filter by part category - exact match
            if (partCategoryId.HasValue)
            {
                var result = await _policyCoveragePartService.GetByPartCategoryIdAsync(partCategoryId.Value);
                return result.IsSuccess ? Ok(result) : NotFound(result);
            }

            // No parameters - get all
            var allResult = await _policyCoveragePartService.GetAllAsync();
            return allResult.IsSuccess ? Ok(allResult) : BadRequest(allResult);
        }

        [HttpGet("{id}")]
        [ProducesResponseType(typeof(BaseResponseDto<PolicyCoveragePartResponse>), StatusCodes.Status200OK)]
        [ProducesResponseType(typeof(BaseResponseDto), StatusCodes.Status404NotFound)]
        public async Task<IActionResult> GetById(Guid id)
        {
            var result = await _policyCoveragePartService.GetByIdAsync(id);
            if (!result.IsSuccess)
                return NotFound(result);

            return Ok(result);
        }

        [HttpGet("coverage-details")]
        [ProducesResponseType(typeof(BaseResponseDto<CoverageDetailsResponse>), StatusCodes.Status200OK)]
        [ProducesResponseType(typeof(BaseResponseDto), StatusCodes.Status404NotFound)]
        [ProducesResponseType(typeof(BaseResponseDto), StatusCodes.Status400BadRequest)]
        public async Task<IActionResult> GetCoverageDetails(
            [FromQuery] Guid policyId,
            [FromQuery] Guid partCategoryId)
        {
            if (policyId == Guid.Empty || partCategoryId == Guid.Empty)
            {
                return BadRequest(new BaseResponseDto
                {
                    IsSuccess = false,
                    Message = "Both policyId and partCategoryId are required",
                    ErrorCode = "INVALID_PARAMETERS"
                });
            }

            var result = await _policyCoveragePartService.GetCoverageDetailsAsync(policyId, partCategoryId);
            return result.IsSuccess ? Ok(result) : NotFound(result);
        }

        // Remove later due to low usage - useless api
        [HttpGet("{id}/details")]
        [ProducesResponseType(typeof(BaseResponseDto<PolicyCoveragePartDetailResponse>), StatusCodes.Status200OK)]
        [ProducesResponseType(typeof(BaseResponseDto), StatusCodes.Status404NotFound)]
        public async Task<IActionResult> GetWithDetails(Guid id)
        {
            var result = await _policyCoveragePartService.GetWithDetailsAsync(id);
            if (!result.IsSuccess)
                return NotFound(result);

            return Ok(result);
        }

        [HttpPost]
        [ProducesResponseType(typeof(BaseResponseDto<PolicyCoveragePartResponse>), StatusCodes.Status201Created)]
        [ProducesResponseType(typeof(BaseResponseDto), StatusCodes.Status400BadRequest)]
        [Authorize(Roles = SystemRoles.UserRoleAdmin + "," + SystemRoles.UserRoleEvmStaff)]
        public async Task<IActionResult> Create([FromBody] CreatePolicyCoveragePartRequest request)
        {
            if (!ModelState.IsValid)
                return BadRequest(ModelState);

            var result = await _policyCoveragePartService.CreateAsync(request);
            if (!result.IsSuccess)
                return BadRequest(result);

            return CreatedAtAction(nameof(GetById), new { id = result.Data.Id }, result);
        }

        [HttpPut("{id}")]
        [ProducesResponseType(typeof(BaseResponseDto<PolicyCoveragePartResponse>), StatusCodes.Status200OK)]
        [ProducesResponseType(typeof(BaseResponseDto), StatusCodes.Status400BadRequest)]
        [ProducesResponseType(typeof(BaseResponseDto), StatusCodes.Status404NotFound)]
        [Authorize(Roles = SystemRoles.UserRoleAdmin + "," + SystemRoles.UserRoleEvmStaff)]
        public async Task<IActionResult> Update(Guid id, [FromBody] UpdatePolicyCoveragePartRequest request)
        {
            if (!ModelState.IsValid)
                return BadRequest(ModelState);

            var result = await _policyCoveragePartService.UpdateAsync(id, request);
            if (!result.IsSuccess)
                return result.ErrorCode == "NOT_FOUND" ? NotFound(result) : BadRequest(result);

            return Ok(result);
        }

        [HttpDelete("{id}")]
        [ProducesResponseType(typeof(BaseResponseDto), StatusCodes.Status200OK)]
        [ProducesResponseType(typeof(BaseResponseDto), StatusCodes.Status404NotFound)]
        [ProducesResponseType(typeof(BaseResponseDto), StatusCodes.Status400BadRequest)]
        [Authorize(Roles = SystemRoles.UserRoleAdmin + "," + SystemRoles.UserRoleEvmStaff)]
        public async Task<IActionResult> Delete(Guid id)
        {
            var result = await _policyCoveragePartService.DeleteAsync(id);
            if (!result.IsSuccess)
                return result.ErrorCode == "NOT_FOUND" ? NotFound(result) : BadRequest(result);

            return Ok(result);
        }
    }
}
