using Backend.Dotnet.Domain.Entities;
using System.ComponentModel.DataAnnotations;
using System.Text.Json.Serialization;
using static Backend.Dotnet.Application.DTOs.VehicleDto;

namespace Backend.Dotnet.Application.DTOs
{
    public class CustomerDto
    {
        /// <summary>
        /// Request DTO for creating a new customer (vehicle owner)
        /// Used when staff registers a new vehicle owner in the warranty system
        /// </summary>
        public class CreateCustomerRequest
        {
            [JsonPropertyName("first_name")]
            [Required(ErrorMessage = "First name is required")]
            [StringLength(100, MinimumLength = 1, ErrorMessage = "First name must be between 1 and 100 characters")]
            public string FirstName { get; set; } = string.Empty;

            [JsonPropertyName("last_name")]
            [Required(ErrorMessage = "Last name is required")]
            [StringLength(100, MinimumLength = 1, ErrorMessage = "Last name must be between 1 and 100 characters")]
            public string LastName { get; set; } = string.Empty;

            [JsonPropertyName("email")]
            [EmailAddress(ErrorMessage = "Invalid email format")]
            [StringLength(255, ErrorMessage = "Email cannot exceed 255 characters")]
            public string? Email { get; set; }

            [JsonPropertyName("phone_number")]
            [Phone(ErrorMessage = "Invalid phone number format")]
            [StringLength(20, ErrorMessage = "Phone number cannot exceed 20 characters")]
            public string? PhoneNumber { get; set; }

            [JsonPropertyName("address")]
            [StringLength(500, ErrorMessage = "Address cannot exceed 500 characters")]
            public string? Address { get; set; }
        }

        /// <summary>
        /// Request DTO for updating existing customer information
        /// Used when staff corrects or updates vehicle owner details
        /// </summary>
        public class UpdateCustomerRequest
        {
            [JsonPropertyName("first_name")]
            [Required(ErrorMessage = "First name is required")]
            [StringLength(100, MinimumLength = 1, ErrorMessage = "First name must be between 1 and 100 characters")]
            public string FirstName { get; set; } = string.Empty;

            [JsonPropertyName("last_name")]
            [Required(ErrorMessage = "Last name is required")]
            [StringLength(100, MinimumLength = 1, ErrorMessage = "Last name must be between 1 and 100 characters")]
            public string LastName { get; set; } = string.Empty;

            [JsonPropertyName("email")]
            [EmailAddress(ErrorMessage = "Invalid email format")]
            [StringLength(255, ErrorMessage = "Email cannot exceed 255 characters")]
            public string? Email { get; set; }

            [JsonPropertyName("phone_number")]
            [Phone(ErrorMessage = "Invalid phone number format")]
            [StringLength(20, ErrorMessage = "Phone number cannot exceed 20 characters")]
            public string? PhoneNumber { get; set; }

            [JsonPropertyName("address")]
            [StringLength(500, ErrorMessage = "Address cannot exceed 500 characters")]
            public string? Address { get; set; }

            [JsonPropertyName("update_mask")]
            [Required(ErrorMessage = "updateMask is required")]
            public List<string> UpdateMask { get; set; } = new();
        }

        /// <summary>
        /// Response DTO for customer basic information
        /// Used in list views and simple customer lookups
        /// </summary>
        public class CustomerResponse
        {
            [JsonPropertyName("id")]
            public Guid Id { get; set; }

            [JsonPropertyName("first_name")]
            public string FirstName { get; set; } = string.Empty;

            [JsonPropertyName("last_name")]
            public string LastName { get; set; } = string.Empty;

            [JsonPropertyName("phone_number")]
            public string? PhoneNumber { get; set; }

            [JsonPropertyName("email")]
            public string? Email { get; set; }

            [JsonPropertyName("address")]
            public string? Address { get; set; }

            [JsonPropertyName("created_at")]
            public DateTime CreatedAt { get; set; }

            [JsonPropertyName("updated_at")]
            public DateTime? UpdatedAt { get; set; }

            [JsonPropertyName("deleted_at")]
            public DateTime? DeletedAt { get; set; }

            [JsonPropertyName("is_deleted")]
            public bool IsDeleted { get; set; }

            // Display helper
            [JsonPropertyName("full_name")]
            public string FullName => $"{FirstName} {LastName}";
        }

        /// <summary>
        /// Response DTO for customer with related vehicles
        /// Used when staff needs to see all vehicles owned by a customer (warranty context)
        /// </summary>
        public class CustomerWithVehiclesResponse
        {
            [JsonPropertyName("id")]
            public Guid Id { get; set; }

            [JsonPropertyName("first_name")]
            public string FirstName { get; set; } = string.Empty;

            [JsonPropertyName("last_name")]
            public string LastName { get; set; } = string.Empty;

            [JsonPropertyName("phone_number")]
            public string? PhoneNumber { get; set; }

            [JsonPropertyName("email")]
            public string? Email { get; set; }

            [JsonPropertyName("address")]
            public string? Address { get; set; }

            [JsonPropertyName("created_at")]
            public DateTime CreatedAt { get; set; }

            [JsonPropertyName("updated_at")]
            public DateTime? UpdatedAt { get; set; }

            [JsonPropertyName("deleted_at")]
            public DateTime? DeletedAt { get; set; }

            [JsonPropertyName("is_deleted")]
            public bool IsDeleted { get; set; }

            [JsonPropertyName("properties")]
            public List<VehicleResponse> Vehicles { get; set; } = new();

            // Display helpers
            [JsonPropertyName("full_name")]
            public string FullName => $"{FirstName} {LastName}";
        }
    }

    public static class CustomerMapper
    {
        public static Customer ToEntity(this CustomerDto.CreateCustomerRequest request)
        {
            return new Customer(
                request.FirstName,
                request.LastName,
                request.Email,
                request.PhoneNumber,
                request.Address
            );
        }

        public static void ApplyToEntity(this CustomerDto.UpdateCustomerRequest request, Customer customer)
        {
            if (request.UpdateMask == null || !request.UpdateMask.Any())
            {
                customer.UpdateProfile(
                    request.FirstName,
                    request.LastName,
                    request.Email,
                    request.PhoneNumber,
                    request.Address);
                return;
            }
            
            if (request.UpdateMask.Contains("first_name", StringComparer.OrdinalIgnoreCase))
                customer.ChangeFirstName(request.FirstName ?? string.Empty);

            if (request.UpdateMask.Contains("last_name", StringComparer.OrdinalIgnoreCase))
                customer.ChangeLastName(request.LastName ?? string.Empty);

            if (request.UpdateMask.Contains("email", StringComparer.OrdinalIgnoreCase))
                customer.ChangeEmail(request.Email ?? string.Empty);

            if (request.UpdateMask.Contains("phone_number", StringComparer.OrdinalIgnoreCase))
                customer.ChangePhoneNumber(request.PhoneNumber ?? string.Empty);

            if (request.UpdateMask.Contains("address", StringComparer.OrdinalIgnoreCase))
                customer.ChangeAddress(request.Address ?? string.Empty);

            //updatedAt occur during change method
        }

        public static CustomerDto.CustomerResponse ToResponse(this Customer customer)
        {
            return new CustomerDto.CustomerResponse
            {
                Id = customer.Id,
                FirstName = customer.FirstName,
                LastName = customer.LastName,
                PhoneNumber = customer.PhoneNumber,
                Email = customer.Email,
                Address = customer.Address,
                CreatedAt = customer.CreatedAt,
                UpdatedAt = customer.UpdatedAt,
                DeletedAt = customer.DeletedAt,
                IsDeleted = customer.IsDeleted
            };
        }

        public static CustomerDto.CustomerWithVehiclesResponse ToWithVehiclesResponse(this Customer customer)
        {
            return new CustomerDto.CustomerWithVehiclesResponse
            {
                Id = customer.Id,
                FirstName = customer.FirstName,
                LastName = customer.LastName,
                PhoneNumber = customer.PhoneNumber,
                Email = customer.Email,
                Address = customer.Address,
                CreatedAt = customer.CreatedAt,
                UpdatedAt = customer.UpdatedAt,
                DeletedAt = customer.DeletedAt,
                IsDeleted = customer.IsDeleted,
                Vehicles = customer.Vehicles.Select(VehicleMapper.ToResponse).ToList()
            };
        }
    }
}
