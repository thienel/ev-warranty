namespace Backend.Dotnet.Domain.Abstractions
{
    public interface IStatus<TEnum> where TEnum : struct
    {
        TEnum Status { get; }
        void ChangeStatus(TEnum newStatus);
    }
}
