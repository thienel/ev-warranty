using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Text.Json.Serialization;
using System.Threading.Tasks;
using static Backend.Dotnet.Application.DTOs.IntermediaryDto;
using static Backend.Dotnet.Infrastructure.External.Models.ExternalModel;

namespace Backend.Dotnet.Infrastructure.External.Models
{
    public class ExternalModel
    {
        public class ExternalServiceResponse<T>
        {
            [JsonPropertyName("is_success")]
            public bool IsSuccess { get; set; }

            [JsonPropertyName("message")]
            public string Message { get; set; }

            [JsonPropertyName("data")]
            public T Data { get; set; }

            [JsonPropertyName("error")]
            public string ErrorCode { get; set; }
        }

        public class ClaimResponseExternal
        {
            [JsonPropertyName("data")]
            public ClaimInfoExternal Data { get; set; }
        }

        public class ClaimInfoExternal
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

        public class ClaimItemsResponseExternal
        {
            [JsonPropertyName("data")]
            public List<ClaimItemExternal> Data { get; set; }
        }

        public class ClaimItemExternal
        {
            [JsonPropertyName("id")]
            public string Id { get; set; }

            [JsonPropertyName("issue_description")]
            public string IssueDescription { get; set; }

            [JsonPropertyName("part_category_id")]
            public string PartCategoryId { get; set; }

            [JsonPropertyName("faulty_part_id")]
            public string FaultyPartId { get; set; }

            [JsonPropertyName("replacement_part_id")]
            public string ReplacementPartId { get; set; }

            [JsonPropertyName("status")]
            public string Status { get; set; }
        }

        public class TechnicianResponseExternal
        {
            [JsonPropertyName("data")]
            public TechnicianInfoExternal Data { get; set; }
        }
        public class TechnicianInfoExternal
        {
            [JsonPropertyName("ID")]
            public Guid Id { get; set; }

            [JsonPropertyName("Name")]
            public Guid Name { get; set; }

            [JsonPropertyName("Role")]
            public string Role { get; set; }
        }
    }

    public static class ExternalToInternalMapper
    {
        public static ClaimInfo ToInternal(this ClaimInfoExternal external)
        {
            if (external == null) return null;

            return new ClaimInfo
            {
                Id = external.Id,
                VehicleId = external.VehicleId,
                Status = external.Status,
                Description = external.Description
            };
        }

        public static TechnicianInfo ToInternal(this TechnicianInfoExternal external)
        {
            if (external == null) return null;

            return new TechnicianInfo
            {
                Id = external.Id,
                Name = external.Name,
                Role = external.Role
            };
        }

        public static ClaimItem ToInternal(this ClaimItemExternal external)
        {
            if (external == null) return null;

            return new ClaimItem
            {
                Id = external.Id,
                IssueDescription = external.IssueDescription,
                PartCategoryId = external.PartCategoryId,
                FaultyPartId = external.FaultyPartId,
                ReplacementPartId = external.ReplacementPartId,
                Status = external.Status
            };
        }

        public static List<ClaimItem> ToInternal(this List<ClaimItemExternal> externals)
        {
            if (externals == null) return new();
            return externals.Select(e => e.ToInternal()).ToList();
        }
    }
}
