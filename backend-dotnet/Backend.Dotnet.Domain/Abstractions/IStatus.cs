using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace Backend.Dotnet.Domain.Abstractions
{
    public interface IStatus<TEnum> where TEnum : struct
    {
        TEnum Status { get; }
        void ChangeStatus(TEnum newStatus);
    }
}
