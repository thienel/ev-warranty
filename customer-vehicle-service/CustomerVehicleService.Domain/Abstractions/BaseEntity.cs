﻿using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace CustomerVehicleService.Domain.Abstractions
{
    public class BaseEntity
    {
        public Guid Id { get; protected set; } = Guid.NewGuid();
        public DateTime CreatedAt { get; protected set; } = DateTime.UtcNow;
        public DateTime? UpdatedAt { get; protected set; } = DateTime.UtcNow;

        protected void SetUpdatedAt()
        {
            UpdatedAt = DateTime.UtcNow;
        }
    }
}
