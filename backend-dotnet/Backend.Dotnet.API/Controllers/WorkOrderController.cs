using Backend.Dotnet.Application.Constants;
using Backend.Dotnet.Application.Interfaces;
using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Mvc;

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
    }
}
