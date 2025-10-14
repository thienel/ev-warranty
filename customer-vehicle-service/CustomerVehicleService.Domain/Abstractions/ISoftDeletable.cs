using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace CustomerVehicleService.Domain.Abstractions
{
    public interface ISoftDeletable
    {
        DateTime? DeletedAt { get; }
        //bool IsDeleted { get; }
        void Delete();
        void Restore();
    }
}
