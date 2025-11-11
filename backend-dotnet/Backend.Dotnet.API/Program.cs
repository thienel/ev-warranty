using Backend.Dotnet.Infrastructure;
using Backend.Dotnet.Infrastructure.Data.Context;
using Backend.Dotnet.Infrastructure.Handler;
using Microsoft.AspNetCore.Authentication;
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

            builder.Services.AddHttpContextAccessor();

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

                c.AddSecurityDefinition("Bearer", new OpenApiSecurityScheme
                {
                    Name = "Authorization",
                    Type = SecuritySchemeType.ApiKey,
                    Description = "Enter 'Bearer <token>'",
                    In = ParameterLocation.Header,
                });

                c.AddSecurityRequirement(new OpenApiSecurityRequirement
                {
                    {
                        new OpenApiSecurityScheme
                        {
                            Reference = new OpenApiReference
                            {
                                Type = ReferenceType.SecurityScheme,
                                Id = "Bearer"
                            },
                            Scheme = "oauth2",
                            Name = "Bearer",
                            In = ParameterLocation.Header
                        },
                        new List<string>()
                    }
                });

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

            });

            builder.Services.AddAuthentication(options =>
            {
                options.DefaultAuthenticateScheme = HeaderRoleAuthenticationHandler.SchemeName;
                options.DefaultChallengeScheme = HeaderRoleAuthenticationHandler.SchemeName;
            })
            .AddScheme<AuthenticationSchemeOptions, HeaderRoleAuthenticationHandler>(
                HeaderRoleAuthenticationHandler.SchemeName, null);

            builder.Services.AddAuthorization();

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

            app.UseAuthentication();

            app.UseAuthorization();

            app.MapControllers();
            app.MapHealthChecks("/health");
            app.Run();
        }
    }
}
