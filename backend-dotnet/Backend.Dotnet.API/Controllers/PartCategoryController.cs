using Backend.Dotnet.Application.Constants;
using Backend.Dotnet.Application.DTOs;
using Backend.Dotnet.Application.Interfaces;
using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Mvc;
using static Backend.Dotnet.Application.DTOs.PartCategoryDto;

namespace Backend.Dotnet.API.Controllers
{
    [ApiController]
    [Route("part-categories")]
    [Produces("application/json")]
    public class PartCategoriesController : ControllerBase
    {
        private readonly IPartCategoryService _partCategoryService;

        public PartCategoriesController(IPartCategoryService partCategoryService)
        {
            _partCategoryService = partCategoryService;
        }

        [HttpGet]
        [ProducesResponseType(typeof(BaseResponseDto<IEnumerable<PartCategoryResponse>>), StatusCodes.Status200OK)]
        [ProducesResponseType(typeof(BaseResponseDto), StatusCodes.Status400BadRequest)]
        [ProducesResponseType(typeof(BaseResponseDto), StatusCodes.Status404NotFound)]
        public async Task<IActionResult> GetAll(
            [FromQuery] string? name = null,
            [FromQuery] Guid? parentId = null)
        {
            if (!string.IsNullOrWhiteSpace(name))
            {
                var result = await _partCategoryService.GetByCategoryNameAsync(name);
                return result.IsSuccess ? Ok(result) : NotFound(result);
            }

            if (parentId.HasValue)
            {
                var result = await _partCategoryService.GetByParentIdAsync(parentId.Value);
                return result.IsSuccess ? Ok(result) : NotFound(result);
            }

            var allResult = await _partCategoryService.GetAllAsync();
            return allResult.IsSuccess ? Ok(allResult) : BadRequest(allResult);
        }

        [HttpGet("{id}")]
        [ProducesResponseType(typeof(BaseResponseDto<PartCategoryResponse>), StatusCodes.Status200OK)]
        [ProducesResponseType(typeof(BaseResponseDto), StatusCodes.Status404NotFound)]
        public async Task<IActionResult> GetById(Guid id)
        {
            var result = await _partCategoryService.GetByIdAsync(id);
            if (!result.IsSuccess)
                return NotFound(result);

            return Ok(result);
        }

        [HttpGet("{id}/hierarchy")]
        [ProducesResponseType(typeof(BaseResponseDto<PartCategoryWithHierarchyResponse>), StatusCodes.Status200OK)]
        [ProducesResponseType(typeof(BaseResponseDto), StatusCodes.Status404NotFound)]
        public async Task<IActionResult> GetWithHierarchy(Guid id)
        {
            var result = await _partCategoryService.GetWithHierarchyAsync(id);
            if (!result.IsSuccess)
                return NotFound(result);

            return Ok(result);
        }

        [HttpGet("hierarchy")]
        [ProducesResponseType(typeof(BaseResponseDto<IEnumerable<PartCategoryWithHierarchyResponse>>), StatusCodes.Status200OK)]
        [ProducesResponseType(typeof(BaseResponseDto), StatusCodes.Status400BadRequest)]
        public async Task<IActionResult> GetFullHierarchy()
        {
            var result = await _partCategoryService.GetFullHierarchyAsync();
            return result.IsSuccess ? Ok(result) : BadRequest(result);
        }

        [HttpPost]
        [ProducesResponseType(typeof(BaseResponseDto<PartCategoryResponse>), StatusCodes.Status201Created)]
        [ProducesResponseType(typeof(BaseResponseDto), StatusCodes.Status400BadRequest)]
        [Authorize(Roles = SystemRoles.UserRoleAdmin + "," + SystemRoles.UserRoleEvmStaff)]
        public async Task<IActionResult> Create([FromBody] CreatePartCategoryRequest request)
        {
            if (!ModelState.IsValid)
                return BadRequest(ModelState);

            var result = await _partCategoryService.CreateAsync(request);
            if (!result.IsSuccess)
                return BadRequest(result);

            return CreatedAtAction(nameof(GetById), new { id = result.Data!.Id }, result);
        }

        [HttpPut("{id}")]
        [ProducesResponseType(typeof(BaseResponseDto<PartCategoryResponse>), StatusCodes.Status200OK)]
        [ProducesResponseType(typeof(BaseResponseDto), StatusCodes.Status400BadRequest)]
        [ProducesResponseType(typeof(BaseResponseDto), StatusCodes.Status404NotFound)]
        [Authorize(Roles = SystemRoles.UserRoleAdmin + "," + SystemRoles.UserRoleEvmStaff)]
        public async Task<IActionResult> Update(Guid id, [FromBody] UpdatePartCategoryRequest request)
        {
            if (!ModelState.IsValid)
                return BadRequest(ModelState);

            var result = await _partCategoryService.UpdateAsync(id, request);
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
            var result = await _partCategoryService.DeleteAsync(id);
            if (!result.IsSuccess)
                return result.ErrorCode == "NOT_FOUND" ? NotFound(result) : BadRequest(result);

            return Ok(result);
        }
    }
}
