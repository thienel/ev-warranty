using Backend.Dotnet.API.Middleware;
using Backend.Dotnet.Infrastructure;
using Backend.Dotnet.Infrastructure.Data.Context;
using Microsoft.EntityFrameworkCore;
using Microsoft.OpenApi.Models;

namespace Backend.Dotnet.API
{
    public class Program
    {
       public static async Task Main(string[] args)
       {
           var builder = WebApplication.CreateBuilder(args);

           builder.Services.AddControllers(options => { })
               .ConfigureApiBehaviorOptions(options =>
               {
                   options.SuppressModelStateInvalidFilter = true;
               });

           builder.Services.AddOpenApi();
           builder.Services.AddInfrastructure(builder.Configuration);
           builder.Services.AddHealthChecks();

           builder.Services.AddEndpointsApiExplorer();
           builder.Services.AddSwaggerGen(c =>
           {
               c.SwaggerDoc("v1", new OpenApiInfo
               {
                   Title = "Customer Vehicle Service API",
                   Version = "v1",
                   Description = "RESTful API for managing customers and their vehicles",
                   Contact = new OpenApiContact
                   {
                       Name = "Your Name / Team",
                       Email = "your@email.com"
                   }
               });
           });

           var app = builder.Build();

           using (var scope = app.Services.CreateScope())
           {
               var services = scope.ServiceProvider;
               var logger = services.GetRequiredService<ILogger<Program>>();

               try
               {
                   var context = services.GetRequiredService<AppDbContext>();

                   var canConnect = await context.Database.CanConnectAsync();

                   if (canConnect)
                   {
                       logger.LogInformation("Database connection successful.");
                       var pendingMigrations = await context.Database.GetPendingMigrationsAsync();

                       if (pendingMigrations.Any())
                       {
                           logger.LogInformation($"Applying {pendingMigrations.Count()} pending migrations...");
                           await context.Database.MigrateAsync();
                           logger.LogInformation("Database migration completed successfully.");
                       }
                       else
                       {
                           logger.LogInformation("No pending migrations. Database is up to date.");
                       }
                   }
                   else
                   {
                       logger.LogWarning("Cannot connect to database. Creating database...");
                       await context.Database.MigrateAsync();
                       logger.LogInformation("Database created and migrated successfully.");
                   }
               }
               catch (Exception ex)
               {
                   logger.LogError(ex, "An error occurred while migrating the database.");
                   if (app.Environment.IsDevelopment())
                   {
                       throw;
                   }
               }
           }

           if (app.Environment.IsDevelopment())
           {
               app.MapOpenApi();
               app.UseSwagger();
               app.UseSwaggerUI();
           }

           app.UseHttpsRedirection();
           app.UseAuthorization();
           app.MapControllers();
           app.MapHealthChecks("/health");

           await app.RunAsync();
       }
    }
}