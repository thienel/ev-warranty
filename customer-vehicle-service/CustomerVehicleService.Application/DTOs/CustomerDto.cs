using CustomerVehicleService.Domain.Entities;
using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using static CustomerVehicleService.Application.DTOs.VehicleDto;

namespace CustomerVehicleService.Application.DTOs
{
    public class CustomerDto
    {
        /// <summary>
        /// Request DTO for creating a new customer (vehicle owner)
        /// Used when staff registers a new vehicle owner in the warranty system
        /// </summary>
        public class CreateCustomerRequest
        {
            [Required(ErrorMessage = "First name is required")]
            [StringLength(100, MinimumLength = 1, ErrorMessage = "First name must be between 1 and 100 characters")]
            public string FirstName { get; set; } = string.Empty;

            [Required(ErrorMessage = "Last name is required")]
            [StringLength(100, MinimumLength = 1, ErrorMessage = "Last name must be between 1 and 100 characters")]
            public string LastName { get; set; } = string.Empty;

            [EmailAddress(ErrorMessage = "Invalid email format")]
            [StringLength(255, ErrorMessage = "Email cannot exceed 255 characters")]
            public string? Email { get; set; }

            [Phone(ErrorMessage = "Invalid phone number format")]
            [StringLength(20, ErrorMessage = "Phone number cannot exceed 20 characters")]
            public string? PhoneNumber { get; set; }

            [StringLength(500, ErrorMessage = "Address cannot exceed 500 characters")]
            public string? Address { get; set; }
        }

        /// <summary>
        /// Request DTO for updating existing customer information
        /// Used when staff corrects or updates vehicle owner details
        /// </summary>
        public class UpdateCustomerRequest
        {
            [Required(ErrorMessage = "First name is required")]
            [StringLength(100, MinimumLength = 1, ErrorMessage = "First name must be between 1 and 100 characters")]
            public string FirstName { get; set; } = string.Empty;

            [Required(ErrorMessage = "Last name is required")]
            [StringLength(100, MinimumLength = 1, ErrorMessage = "Last name must be between 1 and 100 characters")]
            public string LastName { get; set; } = string.Empty;

            [EmailAddress(ErrorMessage = "Invalid email format")]
            [StringLength(255, ErrorMessage = "Email cannot exceed 255 characters")]
            public string? Email { get; set; }

            [Phone(ErrorMessage = "Invalid phone number format")]
            [StringLength(20, ErrorMessage = "Phone number cannot exceed 20 characters")]
            public string? PhoneNumber { get; set; }

            [StringLength(500, ErrorMessage = "Address cannot exceed 500 characters")]
            public string? Address { get; set; }
        }

        /// <summary>
        /// Response DTO for customer basic information
        /// Used in list views and simple customer lookups
        /// </summary>
        public class CustomerResponse
        {
            public Guid Id { get; set; }
            public string FirstName { get; set; } = string.Empty;
            public string LastName { get; set; } = string.Empty;
            public string? PhoneNumber { get; set; }
            public string? Email { get; set; }
            public string? Address { get; set; }
            public DateTime CreatedAt { get; set; }
            public DateTime? UpdatedAt { get; set; }
            public DateTime? DeletedAt { get; set; }
            public bool IsDeleted { get; set; }

            // Display helper
            public string FullName => $"{FirstName} {LastName}";
        }

        /// <summary>
        /// Response DTO for customer with related vehicles
        /// Used when staff needs to see all vehicles owned by a customer (warranty context)
        /// </summary>
        public class CustomerWithVehiclesResponse
        {
            public Guid Id { get; set; }
            public string FirstName { get; set; } = string.Empty;
            public string LastName { get; set; } = string.Empty;
            public string? PhoneNumber { get; set; }
            public string? Email { get; set; }
            public string? Address { get; set; }
            public DateTime CreatedAt { get; set; }
            public DateTime? UpdatedAt { get; set; }
            public DateTime? DeletedAt { get; set; }
            public bool IsDeleted { get; set; }

            public List<VehicleResponse> Vehicles { get; set; } = new();

            // Display helpers
            public string FullName => $"{FirstName} {LastName}";
            public int TotalVehicles => Vehicles.Count;
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
            customer.UpdateProfile(
                request.FirstName,
                request.LastName,
                request.Email,
                request.PhoneNumber,
                request.Address
            );
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
