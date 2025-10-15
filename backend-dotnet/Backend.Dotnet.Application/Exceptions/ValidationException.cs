namespace Backend.Dotnet.Application.Exceptions
{
    public class ValidationException : ApplicationException
    {
        public ValidationException(string message)
            : base(message)
        {
            Errors = new Dictionary<string, string[]>
            {
                { "General", new[] { message } }
            };
        }

        public ValidationException(IDictionary<string, string[]> errors)
            : base(BuildErrorMessage(errors))
        {
            Errors = errors;
        }

        public ValidationException(string field, string error)
            : base($"Validation failed for {field}: {error}")
        {
            Errors = new Dictionary<string, string[]>
            {
                { field, new[] { error } }
            };
        }

        public IDictionary<string, string[]> Errors { get; }

        private static string BuildErrorMessage(IDictionary<string, string[]> errors)
        {
            var errorMessages = errors
                .SelectMany(e => e.Value.Select(v => $"{e.Key}: {v}"))
                .ToList();

            return $"Validation failed: {string.Join("; ", errorMessages)}";
        }
    }
}
