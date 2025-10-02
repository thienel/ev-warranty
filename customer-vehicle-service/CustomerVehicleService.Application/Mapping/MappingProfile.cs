using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using AutoMapper;
using CustomerVehicleService.Application.DTOs;
using CustomerVehicleService.Domain.Entities;

namespace CustomerVehicleService.Application.Mapping
{
    public class MappingProfile : Profile
    {
        // ====================================================================
        // CUSTOMER MAPPINGS
        // ====================================================================

        // Entity -> Response DTO (when returning data to API)
        CreateMap<Customer, CustomerResponse>
            .ForMember(dest => dest.FullName, opt => opt.Ignore()); // Computed property, don't map

        // Entity -> CustomerWithVehiclesResponse (includes navigation)
        CreateMap<Customer, CustomerWithVehiclesResponse>()
            .ForMember(dest => dest.Vehicles, opt => opt.MapFrom(src => src.Vehicles))
            .ForMember(dest => dest.FullName, opt => opt.Ignore())
            .ForMember(dest => dest.TotalVehicles, opt => opt.Ignore());

        // Request DTO -> Entity (when creating new customer)
        CreateMap<CreateCustomerRequest, Customer>()
            .ConstructUsing(src => new Customer(
                src.FirstName,
                src.LastName,
                src.Email,
                src.PhoneNumber,
                src.Address
            ));
        
        // Update Request -> Entity (manual update in service, but useful for reference)
        // Note: We won't use this directly, but it documents the mapping intent
        CreateMap<UpdateCustomerRequest, Customer>()
            .ForMember(dest => dest.Id, opt => opt.Ignore())
            .ForMember(dest => dest.CreatedAt, opt => opt.Ignore())
            .ForMember(dest => dest.UpdatedAt, opt => opt.Ignore())
            .ForMember(dest => dest.DeletedAt, opt => opt.Ignore())
            .ForMember(dest => dest.Vehicles, opt => opt.Ignore());

        // ====================================================================
        // VEHICLE MODEL MAPPINGS
        // ====================================================================

        // Entity -> Response DTO
        CreateMap<VehicleModel, VehicleModelResponse>()
            .ForMember(dest => dest.DisplayName, opt => opt.Ignore());

        // Entity -> VehicleModelWithStatsResponse
        CreateMap<VehicleModel, VehicleModelWithStatsResponse>()
            .ForMember(dest => dest.VehicleCount, opt => opt.MapFrom(src => src.Vehicles.Count))
            .ForMember(dest => dest.DisplayName, opt => opt.Ignore())
            .ForMember(dest => dest.CanBeDeleted, opt => opt.Ignore());

        // Request DTO -> Entity (when creating new model)
        CreateMap<CreateVehicleModelRequest, VehicleModel>()
            .ConstructUsing(src => new VehicleModel(
                src.Brand,
                src.ModelName,
                src.Year
            ));
        
        // Update Request -> Entity (for reference)
        CreateMap<UpdateVehicleModelRequest, VehicleModel>()
            .ForMember(dest => dest.Id, opt => opt.Ignore())
            .ForMember(dest => dest.CreatedAt, opt => opt.Ignore())
            .ForMember(dest => dest.UpdatedAt, opt => opt.Ignore())
            .ForMember(dest => dest.Vehicles, opt => opt.Ignore());

        // ====================================================================
        // VEHICLE MAPPINGS
        // ====================================================================

        // Entity -> Basic Response DTO
        CreateMap<Vehicle, VehicleResponse>();
        
        // Entity -> Detail Response DTO (includes navigation)
        CreateMap<Vehicle, VehicleDetailResponse>()
            .ForMember(dest => dest.Owner, opt => opt.MapFrom(src => src.Customer))
            .ForMember(dest => dest.Model, opt => opt.MapFrom(src => src.Model))
            .ForMember(dest => dest.DisplayName, opt => opt.Ignore())
            .ForMember(dest => dest.OwnerName, opt => opt.Ignore())
            .ForMember(dest => dest.VehicleAgeYears, opt => opt.Ignore());

        // Request DTO -> Entity (when creating new vehicle)
        CreateMap<CreateVehicleRequest, Vehicle>()
            .ConstructUsing(src => new Vehicle(
                src.Vin,
                src.CustomerId,
                src.ModelId,
                src.LicensePlate,
                src.PurchaseDate
            ));
        
        // Update Request -> Entity (for reference)
        CreateMap<UpdateVehicleRequest, Vehicle>()
            .ForMember(dest => dest.Id, opt => opt.Ignore())
            .ForMember(dest => dest.CreatedAt, opt => opt.Ignore())
            .ForMember(dest => dest.UpdatedAt, opt => opt.Ignore())
            .ForMember(dest => dest.Customer, opt => opt.Ignore())
            .ForMember(dest => dest.Model, opt => opt.Ignore());
    }
}
