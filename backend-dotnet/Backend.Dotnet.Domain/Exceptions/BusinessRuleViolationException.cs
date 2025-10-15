namespace Backend.Dotnet.Domain.Exceptions
{
    public class BusinessRuleViolationException : DomainException
    {
        public const string Code = "VIOLATION_000";

        public BusinessRuleViolationException(string message) : base(message, Code)
        {
        }
    }
    public class InvalidOwnerException : DomainException
    {
        public const string Code = "VIOLATION_001";

        public InvalidOwnerException(string Vin) : base($"Vehicle with VIN {Vin} already exists for this customer", Code)
        {
        }
    }
}
