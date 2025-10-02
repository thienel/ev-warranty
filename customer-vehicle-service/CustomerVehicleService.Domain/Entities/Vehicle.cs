using CustomerVehicleService.Domain.Abstractions;
using CustomerVehicleService.Domain.Exceptions;
using System.ComponentModel.DataAnnotations.Schema;
using System.Text.RegularExpressions;

namespace CustomerVehicleService.Domain.Entities
{
    public class Vehicle : BaseEntity
    {
        public string Vin { get; private set; }
        public string LicensePlate { get; private set; }

        public Guid CustomerId { get; private set; }
        public Guid ModelId { get; private set; }

        public DateTime? PurchaseDate { get; private set; }

        // Navigation properties
        [ForeignKey("CustomerId")]
        public Customer Customer { get; private set; }
        public VehicleModel Model { get; private set; }


        // EF Core constructor
        private Vehicle() { }

        public Vehicle(string vin, Guid customerId, Guid modelId, string licensePlate = null,
            DateTime? purchaseDate = null)
        {
            SetVin(vin);
            CustomerId = customerId;
            ModelId = modelId;
            SetLicensePlate(licensePlate);
            SetPurchaseDate(purchaseDate);
        }

        private void SetVin(string vin)
        {
            if (string.IsNullOrWhiteSpace(vin))
                throw new BusinessRuleViolationException("VIN is required");

            if (!IsValidVin(vin))
                throw new BusinessRuleViolationException("Invalid VIN format");

            Vin = vin.Trim().ToUpperInvariant();
        }

        private void SetLicensePlate(string licensePlate)
        {
            if (!string.IsNullOrWhiteSpace(licensePlate))
            {
                if (licensePlate.Length > 20)
                    throw new BusinessRuleViolationException("License plate cannot exceed 20 characters");
            }

            LicensePlate = licensePlate?.Trim().ToUpperInvariant();
        }

        private void SetPurchaseDate(DateTime? purchaseDate)
        {
            if (purchaseDate.HasValue && purchaseDate.Value.Date > DateTime.Now.Date)
                throw new BusinessRuleViolationException("Purchase date cannot be in the future");

            PurchaseDate = purchaseDate?.Date;
        }

        private static bool IsValidVin(string vin)
        {
            if (vin.Length != 17)
                return false;

            var vinRegex = new Regex(@"^[A-HJ-NPR-Z0-9]{17}$", RegexOptions.Compiled);
            return vinRegex.IsMatch(vin);
        }

    }
}
