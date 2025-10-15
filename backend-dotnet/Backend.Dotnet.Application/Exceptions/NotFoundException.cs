namespace Backend.Dotnet.Application.Exceptions
{
    public class NotFoundException : ApplicationException
    {
        public NotFoundException(string entityName, Guid id)
            : base($"{entityName} with ID '{id}' was not found.")
        {
            EntityName = entityName;
            EntityId = id;
        }

        public NotFoundException(string entityName, string key, string value)
            : base($"{entityName} with {key} '{value}' was not found.")
        {
            EntityName = entityName;
            Key = key;
            Value = value;
        }

        public string EntityName { get; }
        public Guid? EntityId { get; }
        public string Key { get; }
        public string Value { get; }
    }
}
