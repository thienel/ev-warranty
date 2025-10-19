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

COPY --from=build /app/publish .

USER $APP_UID

ENTRYPOINT ["dotnet", "Backend.Dotnet.API.dll"]