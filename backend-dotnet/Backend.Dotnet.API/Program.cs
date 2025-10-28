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
                /*
                c.AddSecurityDefinition("X-User-Role", new OpenApiSecurityScheme
                {
                    In = ParameterLocation.Header,
                    Name = "X-User-Role",
                    Type = SecuritySchemeType.ApiKey,
                    Description = "User Role header"
                });

                c.AddSecurityRequirement(new OpenApiSecurityRequirement
                {
                    {
                        new OpenApiSecurityScheme
                        {
                            Reference = new OpenApiReference
                            {
                                Type = ReferenceType.SecurityScheme,
                                Id = "X-User-Role"
                            }
                        },
                        Array.Empty<string>()
                    }
                });
                */
            });

            //builder.Services.AddAuthorization();
            var app = builder.Build();

            using (var scope = app.Services.CreateScope())
            {
                var services = scope.ServiceProvider;
                try
                {
                    var context = services.GetRequiredService<AppDbContext>();
                    context.Database.Migrate();
                    var logger = services.GetRequiredService<ILogger<Program>>();
                    logger.LogInformation("Database migration completed successfully.");
                }
                catch (Exception ex)
                {
                    var logger = services.GetRequiredService<ILogger<Program>>();
                    logger.LogError(ex, "An error occurred while migrating the database.");
                    throw;
                }
            }

            if (app.Environment.IsDevelopment())
            {
                app.MapOpenApi();
                app.UseSwagger();
                app.UseSwaggerUI();
            }

            app.UseHttpsRedirection();

            //app.UseMiddleware<HeaderAuthMiddleware>();

            app.UseAuthorization();
            app.MapControllers();
            app.MapHealthChecks("/health");
            app.Run();
        }
    }
}
