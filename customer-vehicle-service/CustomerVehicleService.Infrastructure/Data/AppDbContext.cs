using Microsoft.EntityFrameworkCore;
using Microsoft.AspNetCore.Identity.EntityFrameworkCore;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using CustomerVehicleService.Domain.Entities;

namespace CustomerVehicleService.Infrastructure.Data
{
    public class AppDbContext// : IdentityDbContext<Customer>
    {
        public  DbSet<Customer> Customers { get; set; }
        public DbSet<VehicleModel> VehicleModels { get; set; }
        public DbSet<Vehicle> Vehicle { get; set; }

    }
}
