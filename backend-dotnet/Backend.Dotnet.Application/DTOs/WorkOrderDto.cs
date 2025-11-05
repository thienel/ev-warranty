using Backend.Dotnet.Domain.Entities;
using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;
using System.Text.Json.Serialization;
using System.Threading.Tasks;
using static Backend.Dotnet.Application.DTOs.IntermediaryDto;

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
            [EnumDataType(typeof(WorkOrderStatus))]
            public WorkOrderStatus Status { get; set; }
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

    }
}
