using CustomerVehicleService.Domain.Abstractions;
using CustomerVehicleService.Domain.Exceptions;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace CustomerVehicleService.Domain.Entities
{
    public class VehicleModel : BaseEntity
    {
        public string Brand { get; private set; }
        public string ModelName { get; private set; }
        public int Year { get; private set; }

        // Navigation property
        public virtual ICollection<Vehicle> Vehicles { get; private set; }

        private VehicleModel()
        {
            Vehicles = new List<Vehicle>();
        }

        public VehicleModel(string brand, string modelName, int year)
        {
            SetBrand(brand);
            SetModelName(modelName);
            SetYear(year);
        }

        // BEHAVIOUR METHODS
        public void UpdateModel(string brand, string modelName, int year)
        {
            SetBrand(brand);
            SetModelName(modelName);
            SetYear(year);
            SetUpdatedAt();
        }

        public void ChangeBrand(string brand)
        {
            SetBrand(brand);
            SetUpdatedAt();
        }

        public void ChangeModelName(string modelName)
        {
            SetModelName(modelName);
            SetUpdatedAt();
        }

        public void ChangeYear(int year)
        {
            SetYear(year);
            SetUpdatedAt();
        }

        // PRIVATE SETTERS
        private void SetBrand(string brand)
        {
            if (string.IsNullOrWhiteSpace(brand))
                throw new BusinessRuleViolationException("Brand is required");

            if (brand.Length > 100)
                throw new BusinessRuleViolationException("Brand cannot exceed 100 characters");

            Brand = brand.Trim();
        }

        private void SetModelName(string modelName)
        {
            if (string.IsNullOrWhiteSpace(modelName))
                throw new BusinessRuleViolationException("Model name is required");

            if (modelName.Length > 100)
                throw new BusinessRuleViolationException("Model name cannot exceed 100 characters");

            ModelName = modelName.Trim();
        }

        private void SetYear(int year)
        {
            var currentYear = DateTime.Now.Year;
            if (year < 1900 || year > currentYear + 2)
                throw new BusinessRuleViolationException($"Year must be between 1900 and {currentYear + 2}");

            Year = year;
        }
    }
}
