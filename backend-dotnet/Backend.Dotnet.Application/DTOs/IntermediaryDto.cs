using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Text.Json.Serialization;
using System.Threading.Tasks;

namespace Backend.Dotnet.Application.DTOs
{
    public class IntermediaryDto
    {
        public class ClaimInfo
        {
            [JsonPropertyName("id")]
            public Guid Id { get; set; }

            [JsonPropertyName("vehicle_id")]
            public Guid VehicleId { get; set; }

            [JsonPropertyName("status")]
            public string Status { get; set; }

            [JsonPropertyName("description")]
            public string Description { get; set; }
        }

        public class TechnicianInfo
        {
            [JsonPropertyName("ID")]
            public Guid Id { get; set; }

            [JsonPropertyName("Name")]
            public string Name { get; set; }

            [JsonPropertyName("Role")]
            public string Role { get; set; }
        }
        public class ClaimItem
        {
            [JsonPropertyName("id")]
            public Guid Id { get; set; }

            [JsonPropertyName("issue_description")]
            public string IssueDescription { get; set; }

            [JsonPropertyName("part_category_id")]
            public Guid PartCategoryId { get; set; }

            [JsonPropertyName("faulty_part_id")]
            public Guid FaultyPartId { get; set; }

            [JsonPropertyName("replacement_part_id")]
            public Guid ReplacementPartId { get; set; }

            [JsonPropertyName("status")]
            public string Status { get; set; }
        }
    }
}
