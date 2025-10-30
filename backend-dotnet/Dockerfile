FROM mcr.microsoft.com/dotnet/sdk:9.0 AS build
WORKDIR /src

COPY Backend.Dotnet.sln ./
COPY Backend.Dotnet.API/Backend.Dotnet.API.csproj Backend.Dotnet.API/
COPY Backend.Dotnet.Application/Backend.Dotnet.Application.csproj Backend.Dotnet.Application/
COPY Backend.Dotnet.Domain/Backend.Dotnet.Domain.csproj Backend.Dotnet.Domain/
COPY Backend.Dotnet.Infrastructure/Backend.Dotnet.Infrastructure.csproj Backend.Dotnet.Infrastructure/

RUN dotnet restore

COPY . .

WORKDIR /src/Backend.Dotnet.API
RUN dotnet publish -c Release -o /app/publish /p:UseAppHost=false

FROM mcr.microsoft.com/dotnet/aspnet:9.0 AS final
WORKDIR /app
EXPOSE 8080
EXPOSE 8081

USER root


RUN apt-get update && \
    apt-get install -y curl gnupg2 ca-certificates && \
    curl -fsSL https://packages.microsoft.com/keys/microsoft.asc | gpg --dearmor -o /usr/share/keyrings/microsoft-prod.gpg && \
    curl https://packages.microsoft.com/config/debian/12/prod.list | tee /etc/apt/sources.list.d/mssql-release.list && \
    apt-get update && \
    ACCEPT_EULA=Y apt-get install -y mssql-tools18 unixodbc-dev && \
    echo 'export PATH="$PATH:/opt/mssql-tools18/bin"' >> ~/.bashrc && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

COPY --from=build /app/publish .

# Copy seed data files
COPY seed_data_up.sql /app/seed_data_up.sql
COPY seed_data_down.sql /app/seed_data_down.sql

# Copy entrypoint script
COPY entrypoint.sh /app/entrypoint.sh
RUN chmod +x /app/entrypoint.sh

USER $APP_UID

ENTRYPOINT ["/app/entrypoint.sh"]