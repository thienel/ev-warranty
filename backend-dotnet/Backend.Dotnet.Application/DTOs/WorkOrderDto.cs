using Backend.Dotnet.Domain.Entities;
using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;
using System.Text.Json.Serialization;
using System.Threading.Tasks;
using static Backend.Dotnet.Application.DTOs.IntermediaryDto;
using static Backend.Dotnet.Application.DTOs.WorkOrderDto;

namespace Backend.Dotnet.Application.DTOs
{
    public class WorkOrderDto
    {
        public class CreateWorkOrderRequest
        {
            [JsonPropertyName("claim_id")]
            [Required(ErrorMessage = "Claim ID is required")]
            public Guid ClaimId { get; set; }

            [JsonPropertyName("assigned_technician_id")]
            [Required(ErrorMessage = "Technician ID is required")]
            public Guid AssignedTechnicianId { get; set; }
        }

        public class UpdateStatusRequest
        {
            [JsonPropertyName("status")]
            [Required(ErrorMessage = "Status is required")]
            [RegularExpression("^(Pending|InProgress|ToVerify)$",
                ErrorMessage = "Status must be Pending, InProgress, ToVerify")]
            public string Status { get; set; }
        }

        public class WorkOrderResponse
        {
            [JsonPropertyName("id")]
            public Guid Id { get; set; }

            [JsonPropertyName("claim_id")]
            public Guid ClaimId { get; set; }

            [JsonPropertyName("assigned_technician_id")]
            public Guid AssignedTechnicianId { get; set; }

            [JsonPropertyName("status")]
            public string Status { get; set; }

            [JsonPropertyName("scheduled_date")]
            public DateTime ScheduledDate { get; set; }

            [JsonPropertyName("completed_date")]
            public DateTime? CompletedDate { get; set; }

            [JsonPropertyName("note")]
            public string Note { get; set; }

            [JsonPropertyName("created_at")]
            public DateTime CreatedAt { get; set; }

            [JsonPropertyName("updated_at")]
            public DateTime? UpdatedAt { get; set; }
        }

        public class WorkOrderDetailResponse
        {
            [JsonPropertyName("id")]
            public Guid Id { get; set; }

            [JsonPropertyName("claim_id")]
            public Guid ClaimId { get; set; }

            [JsonPropertyName("assigned_technician_id")]
            public Guid AssignedTechnicianId { get; set; }

            [JsonPropertyName("status")]
            public string Status { get; set; }

            [JsonPropertyName("scheduled_date")]
            public DateTime ScheduledDate { get; set; }

            [JsonPropertyName("completed_date")]
            public DateTime? CompletedDate { get; set; }

            [JsonPropertyName("note")]
            public string Note { get; set; }

            [JsonPropertyName("created_at")]
            public DateTime CreatedAt { get; set; }

            [JsonPropertyName("updated_at")]
            public DateTime? UpdatedAt { get; set; }

            // External data
            [JsonPropertyName("claim_info")]
            public ClaimInfo Claim { get; set; }

            [JsonPropertyName("technician_info")]
            public TechnicianInfo Technician { get; set; }

            [JsonPropertyName("claim_items")]
            public List<ClaimItem> ClaimItems { get; set; }
        }
    }
    public static class WorkOrderMapper
    {
        public static WorkOrder ToEntity(this CreateWorkOrderRequest request, DateTime scheduledDate)
        {
            return new WorkOrder(
                request.ClaimId,
                request.AssignedTechnicianId,
                scheduledDate
            );
        }

        public static WorkOrderResponse ToResponse(this WorkOrder entity)
        {
            return new WorkOrderResponse
            {
                Id = entity.Id,
                ClaimId = entity.ClaimId,
                AssignedTechnicianId = entity.AssignedTechnicianId,
                Status = entity.Status.ToString(),
                ScheduledDate = entity.ScheduledDate,
                CompletedDate = entity.CompletedDate,
                Note = entity.Note,
                CreatedAt = entity.CreatedAt,
                UpdatedAt = entity.UpdatedAt
            };
        }

        public static WorkOrderDetailResponse ToDetailResponse(
            this WorkOrder entity,
            ClaimInfo claim,
            TechnicianInfo technician,
            List<ClaimItem> claimItems)
        {
            return new WorkOrderDetailResponse
            {
                Id = entity.Id,
                ClaimId = entity.ClaimId,
                AssignedTechnicianId = entity.AssignedTechnicianId,
                Status = entity.Status.ToString(),
                ScheduledDate = entity.ScheduledDate,
                CompletedDate = entity.CompletedDate,
                Note = entity.Note,
                CreatedAt = entity.CreatedAt,
                UpdatedAt = entity.UpdatedAt,
                Claim = claim,
                Technician = technician,
                ClaimItems = claimItems
            };
        }
    }
}
