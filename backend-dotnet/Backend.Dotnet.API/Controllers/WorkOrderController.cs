using Backend.Dotnet.Application.Constants;
using Backend.Dotnet.Application.DTOs;
using Backend.Dotnet.Application.Interfaces;
using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Mvc;
using static Backend.Dotnet.Application.DTOs.WorkOrderDto;

namespace Backend.Dotnet.API.Controllers
{
    [ApiController]
    [Route("work-orders")]
    [Produces("application/json")]
    [Authorize(Roles = SystemRoles.UserRoleScTechnician + "," + SystemRoles.UserRoleScStaff)]

    public class WorkOrderController : Controller
    {
        private readonly IWorkOrderService _workOrderService;

        public WorkOrderController(IWorkOrderService workOrderService)
        {
            _workOrderService = workOrderService;
        }

        [HttpGet]
        [ProducesResponseType(typeof(BaseResponseDto<IEnumerable<WorkOrderResponse>>), StatusCodes.Status200OK)]
        [ProducesResponseType(typeof(BaseResponseDto), StatusCodes.Status400BadRequest)]
        public async Task<IActionResult> GetAll(
            [FromQuery] Guid? claimId = null,
            [FromQuery] Guid? technicianId = null)
        {
            // Claim ID - exact - OtO
            if (claimId.HasValue)
            {
                var result = await _workOrderService.GetByIdAsync(claimId.Value);
                return result.IsSuccess ? Ok(result) : NotFound(result);
            }

            // Tech ID - list work
            if (technicianId.HasValue)
            {
                var result = await _workOrderService.GetByTechnicianIdAsync(technicianId.Value);
                return result.IsSuccess ? Ok(result) : NotFound(result);
            }

            // No params - get all 
            var allResult = await _workOrderService.GetAllAsync();
            return allResult.IsSuccess ? Ok(allResult) : BadRequest(allResult);
        }

        [HttpGet("{id}")]
        [ProducesResponseType(typeof(BaseResponseDto<WorkOrderResponse>), StatusCodes.Status200OK)]
        [ProducesResponseType(typeof(BaseResponseDto), StatusCodes.Status404NotFound)]
        public async Task<IActionResult> GetById(Guid id)
        {
            var result = await _workOrderService.GetByIdAsync(id);
            if (!result.IsSuccess)
                return NotFound(result);

            return Ok(result);
        }


        [HttpGet("{id}/details")]
        [ProducesResponseType(typeof(BaseResponseDto<WorkOrderDetailResponse>), StatusCodes.Status200OK)]
        [ProducesResponseType(typeof(BaseResponseDto), StatusCodes.Status404NotFound)]
        public async Task<IActionResult> GetDetailById(Guid id)
        {
            var result = await _workOrderService.GetDetailByIdAsync(id);
            if (!result.IsSuccess)
                return NotFound(result);

            return Ok(result);
        }

        [HttpPost]
        [ProducesResponseType(typeof(BaseResponseDto<WorkOrderResponse>), StatusCodes.Status201Created)]
        [ProducesResponseType(typeof(BaseResponseDto), StatusCodes.Status400BadRequest)]
        [Authorize(Roles = SystemRoles.UserRoleEvmStaff)]
        public async Task<IActionResult> Create([FromBody] CreateWorkOrderRequest request)
        {
            if (!ModelState.IsValid)
                return BadRequest(ModelState);

            var result = await _workOrderService.CreateAsync(request);
            if (!result.IsSuccess)
                return BadRequest(result);

            return CreatedAtAction(nameof(GetById), new { id = result.Data.Id }, result);
        }

        [HttpPut("{id}/status")]
        [ProducesResponseType(typeof(BaseResponseDto<WorkOrderResponse>), StatusCodes.Status200OK)]
        [ProducesResponseType(typeof(BaseResponseDto), StatusCodes.Status400BadRequest)]
        [ProducesResponseType(typeof(BaseResponseDto), StatusCodes.Status404NotFound)]
        public async Task<IActionResult> UpdateStatus(Guid id, [FromBody] UpdateStatusRequest request)
        {
            if (!ModelState.IsValid)
                return BadRequest(ModelState);

            var result = await _workOrderService.UpdateStatusAsync(id, request);
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
            var result = await _workOrderService.DeleteAsync(id);
            if (!result.IsSuccess)
                return result.ErrorCode == "NOT_FOUND" ? NotFound(result) : BadRequest(result);

            return Ok(result);
        }
    }
}
