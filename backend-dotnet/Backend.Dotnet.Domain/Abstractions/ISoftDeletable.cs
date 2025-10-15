namespace Backend.Dotnet.Domain.Abstractions
{
    public interface ISoftDeletable
    {
        DateTime? DeletedAt { get; }
        //bool IsDeleted { get; }
        void Delete();
        void Restore();
    }
}
