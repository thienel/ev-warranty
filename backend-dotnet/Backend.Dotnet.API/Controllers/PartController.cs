using Backend.Dotnet.Application.DTOs;
using Backend.Dotnet.Application.Interfaces;
using Microsoft.AspNetCore.Http;
using Microsoft.AspNetCore.Mvc;
using static Backend.Dotnet.Application.DTOs.PartDto;

namespace Backend.Dotnet.API.Controllers
{
    [ApiController]
    [Route("parts")]
    [Produces("application/json")]
    public class PartsController : ControllerBase
    {
        private readonly IPartService _partService;

        public PartsController(IPartService partService)
        {
            _partService = partService;
        }

        [HttpGet]
        [ProducesResponseType(typeof(BaseResponseDto<IEnumerable<PartResponse>>), StatusCodes.Status200OK)]
        [ProducesResponseType(typeof(BaseResponseDto), StatusCodes.Status400BadRequest)]
        public async Task<IActionResult> GetAll(
            [FromQuery] string serialNumber = null,
            [FromQuery] string status = null,
            [FromQuery] Guid? categoryId = null,
            [FromQuery] Guid? officeLocationId = null,
            [FromQuery] string search = null)
        {
            // Serial number - absolute
            if (!string.IsNullOrWhiteSpace(serialNumber))
            {
                var result = await _partService.GetBySerialNumberAsync(serialNumber);
                return result.IsSuccess ? Ok(result) : NotFound(result);
            }

            // Status - exact
            if (!string.IsNullOrWhiteSpace(status))
            {
                var result = await _partService.GetByStatusAsync(status);
                return result.IsSuccess ? Ok(result) : BadRequest(result);
            }

            // Category ID - exact
            if (categoryId.HasValue)
            {
                var result = await _partService.GetByCategoryIdAsync(categoryId.Value);
                return result.IsSuccess ? Ok(result) : NotFound(result);
            }

            // Office location ID - exact
            if (officeLocationId.HasValue)
            {
                var result = await _partService.GetByOfficeLocationIdAsync(officeLocationId.Value);
                return result.IsSuccess ? Ok(result) : BadRequest(result);
            }

            // Search term - relative
            if (!string.IsNullOrWhiteSpace(search))
            {
                var result = await _partService.SearchAsync(search);
                return result.IsSuccess ? Ok(result) : BadRequest(result);
            }

            // No parameters - get all
            var allResult = await _partService.GetAllAsync();
            return allResult.IsSuccess ? Ok(allResult) : BadRequest(allResult);
        }

        [HttpGet("{id}")]
        [ProducesResponseType(typeof(BaseResponseDto<PartResponse>), StatusCodes.Status200OK)]
        [ProducesResponseType(typeof(BaseResponseDto), StatusCodes.Status404NotFound)]
        public async Task<IActionResult> GetById(Guid id)
        {
            var result = await _partService.GetByIdAsync(id);
            if (!result.IsSuccess)
                return NotFound(result);

            return Ok(result);
        }

        [HttpGet("{id}/details")]
        [ProducesResponseType(typeof(BaseResponseDto<PartWithDetailsResponse>), StatusCodes.Status200OK)]
        [ProducesResponseType(typeof(BaseResponseDto), StatusCodes.Status404NotFound)]
        public async Task<IActionResult> GetWithDetails(Guid id)
        {
            var result = await _partService.GetWithDetailsAsync(id);
            if (!result.IsSuccess)
                return NotFound(result);

            return Ok(result);
        }

        [HttpPost]
        [ProducesResponseType(typeof(BaseResponseDto<PartResponse>), StatusCodes.Status201Created)]
        [ProducesResponseType(typeof(BaseResponseDto), StatusCodes.Status400BadRequest)]
        public async Task<IActionResult> Create([FromBody] CreatePartRequest request)
        {
            if (!ModelState.IsValid)
                return BadRequest(ModelState);

            var result = await _partService.CreateAsync(request);
            if (!result.IsSuccess)
                return BadRequest(result);

            return CreatedAtAction(nameof(GetById), new { id = result.Data.Id }, result);
        }

        [HttpPut("{id}")]
        [ProducesResponseType(typeof(BaseResponseDto<PartResponse>), StatusCodes.Status200OK)]
        [ProducesResponseType(typeof(BaseResponseDto), StatusCodes.Status400BadRequest)]
        [ProducesResponseType(typeof(BaseResponseDto), StatusCodes.Status404NotFound)]
        public async Task<IActionResult> Update(Guid id, [FromBody] UpdatePartRequest request)
        {
            if (!ModelState.IsValid)
                return BadRequest(ModelState);

            var result = await _partService.UpdateAsync(id, request);
            if (!result.IsSuccess)
                return result.ErrorCode == "NOT_FOUND" ? NotFound(result) : BadRequest(result);

            return Ok(result);
        }

        [HttpPut("{id}/category")]
        [ProducesResponseType(typeof(BaseResponseDto<PartResponse>), StatusCodes.Status200OK)]
        [ProducesResponseType(typeof(BaseResponseDto), StatusCodes.Status400BadRequest)]
        [ProducesResponseType(typeof(BaseResponseDto), StatusCodes.Status404NotFound)]
        public async Task<IActionResult> ChangeCategory(Guid id, [FromBody] ChangePartCategoryRequest request)
        {
            if (!ModelState.IsValid)
                return BadRequest(ModelState);

            var result = await _partService.ChangeCategoryAsync(id, request);
            if (!result.IsSuccess)
                return result.ErrorCode == "NOT_FOUND" || result.ErrorCode == "CATEGORY_NOT_FOUND"
                    ? NotFound(result)
                    : BadRequest(result);

            return Ok(result);
        }

        [HttpPut("{id}/status")]
        [ProducesResponseType(typeof(BaseResponseDto<PartResponse>), StatusCodes.Status200OK)]
        [ProducesResponseType(typeof(BaseResponseDto), StatusCodes.Status400BadRequest)]
        [ProducesResponseType(typeof(BaseResponseDto), StatusCodes.Status404NotFound)]
        public async Task<IActionResult> ChangeStatus(Guid id, [FromBody] PartChangeStatusRequest request)
        {
            if (!ModelState.IsValid)
                return BadRequest(ModelState);

            var result = await _partService.ChangeStatusAsync(id, request);
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
            var result = await _partService.DeleteAsync(id);
            if (!result.IsSuccess)
                return result.ErrorCode == "NOT_FOUND" ? NotFound(result) : BadRequest(result);

            return Ok(result);
        }
    }
}
