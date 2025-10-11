using CustomerVehicleService.Domain.Abstractions;
using CustomerVehicleService.Domain.Exceptions;
using System.Text.RegularExpressions;

namespace CustomerVehicleService.Domain.Entities
{
    public class Customer : BaseEntity, ISoftDeletable
    {
        public string FirstName { get; private set; }
        public string LastName { get; private set; }
        public string PhoneNumber { get; private set; }
        public string Email { get; private set; }
        public string Address { get; private set; }
        public string FullName => $"{FirstName} {LastName}";
        public DateTime? DeletedAt { get; private set; }
        public bool IsDeleted => DeletedAt.HasValue;

        // Navigation property
        public virtual ICollection<Vehicle> Vehicles { get; private set; }

        private Customer()
        {
            Vehicles = new List<Vehicle>();
        }

        public Customer(string firstName, string lastName, string? email = null,
                       string phoneNumber = null, string? address = null) : this()
        {
            SetFirstName(firstName);
            SetLastName(lastName);
            SetEmail(email);
            SetPhoneNumber(phoneNumber);
            SetAddress(address);
        }

        // BEHAVIOUR METHODS
        public void UpdateProfile(string firstName, string lastName, string? email, string phoneNumber, string? address)
        {
            SetFirstName(firstName);
            SetLastName(lastName);
            SetEmail(email);
            SetPhoneNumber(phoneNumber);
            SetAddress(address);
            SetUpdatedAt();
        }

        public void ChangePhoneNumber(string phoneNumber)
        {
            SetPhoneNumber(phoneNumber);
            SetUpdatedAt();
        }

        public void ChangeEmail(string email)
        {
            SetEmail(email);
            SetUpdatedAt();
        }

        public void ChangeAddress(string address)
        {
            SetAddress(address);
            SetUpdatedAt();
        }

        public void AddVehicle(Vehicle vehicle)
        {
            if (vehicle == null) throw new ArgumentNullException(nameof(vehicle));

            if (Vehicles.Any(v => v.Id == vehicle.Id))
                throw new BusinessRuleViolationException("Vehicle already exists for this customer");

            Vehicles.Add(vehicle);
            SetUpdatedAt();
        }

        public void RemoveVehicle(Guid vehicleId)
        {
            var vehicle = Vehicles.FirstOrDefault(v => v.Id == vehicleId);
            if (vehicle == null)
                throw new BusinessRuleViolationException("Vehicle not found for this customer");

            Vehicles.Remove(vehicle);
            SetUpdatedAt();
        }

        // Soft delete mechanism
        public void Delete()
        {
            if (DeletedAt.HasValue)
                throw new BusinessRuleViolationException("Customer is already deleted");

            DeletedAt = DateTime.UtcNow;
            SetUpdatedAt();
        }

        public void Restore()
        {
            if (!DeletedAt.HasValue)
                throw new BusinessRuleViolationException("Customer is not deleted");

            DeletedAt = null;
            SetUpdatedAt();
        }

        // PRIVATE SETTER
        private void SetFirstName(string firstName)
        {
            if (string.IsNullOrWhiteSpace(firstName))
                throw new BusinessRuleViolationException("First name is required");

            if (firstName.Length > 100)
                throw new BusinessRuleViolationException("First name cannot exceed 100 characters");

            FirstName = firstName.Trim();
        }

        private void SetLastName(string lastName)
        {
            if (string.IsNullOrWhiteSpace(lastName))
                throw new BusinessRuleViolationException("Last name is required");

            if (lastName.Length > 100)
                throw new BusinessRuleViolationException("Last name cannot exceed 100 characters");

            LastName = lastName.Trim();
        }

        private void SetPhoneNumber(string phoneNumber)
        {
            if (!string.IsNullOrWhiteSpace(phoneNumber))
            {
                if (phoneNumber.Length > 20)
                    throw new BusinessRuleViolationException("Phone number cannot exceed 20 characters");
            }

            PhoneNumber = phoneNumber.Trim();
        }

        private void SetEmail(string email)
        {
            if (!string.IsNullOrWhiteSpace(email))
            {
                if (!IsValidEmail(email))
                    throw new BusinessRuleViolationException("Invalid email format");

                if (email.Length > 255)
                    throw new BusinessRuleViolationException("Email cannot exceed 255 characters");
            }

            Email = email.Trim().ToLowerInvariant();
        }

        private void SetAddress(string address)
        {
            Address = address.Trim();
        }

        private static bool IsValidEmail(string email)
        {
            var emailRegex = new Regex(@"^[^\s@]+@[^\s@]+\.[^\s@]+$", RegexOptions.Compiled);
            return emailRegex.IsMatch(email);
        }
    }
}
