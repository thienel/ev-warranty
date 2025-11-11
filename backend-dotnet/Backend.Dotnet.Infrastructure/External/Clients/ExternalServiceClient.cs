using Backend.Dotnet.Application.Interfaces.External;
using Backend.Dotnet.Infrastructure.External.Models;
using Microsoft.AspNetCore.Http;
using System.Net.Http.Headers;
using System.Text.Json;
using static Backend.Dotnet.Application.DTOs.IntermediaryDto;
using static Backend.Dotnet.Infrastructure.External.Models.ExternalModel;

namespace Backend.Dotnet.Infrastructure.External.Clients
{
    public class ExternalServiceClient : IExternalServiceClient
    {
        private readonly HttpClient _httpClient;
        private readonly IHttpContextAccessor _httpContextAccessor;

        private const string ExternalServiceBaseUrl = "http://localhost";

        public ExternalServiceClient(HttpClient httpClient, IHttpContextAccessor httpContextAccessor)
        {
            _httpClient = httpClient;
            _httpContextAccessor = httpContextAccessor;
        }

        public async Task<bool> CompleteClaimAsync(Guid claimId)
        {
            try
            {
                var token = _httpContextAccessor.HttpContext?.Request.Headers["Authorization"].ToString();

                if (string.IsNullOrEmpty(token))
                    throw new InvalidOperationException("Missing Authorization header in request");
                // Replace() for clean jwt token
                var jwtToken = token.Replace("Bearer ", "");
                _httpClient.DefaultRequestHeaders.Authorization = new AuthenticationHeaderValue("Bearer", jwtToken);
                
                var response = await _httpClient.PostAsync($"{ExternalServiceBaseUrl}/api/v1/claims/{claimId}/complete", null);

                if (!response.IsSuccessStatusCode)
                {
                    throw new InvalidOperationException($"Failed to complete claim {claimId}. Status: {response.StatusCode}");
                }

                return true;
            }
            catch (Exception ex)
            {

                throw new InvalidOperationException(ex.Message);
            }
        }

        public async Task<ClaimInfo> GetClaimAsync(Guid claimId)
        {
            try
            {
                var response = await _httpClient.GetAsync($"{ExternalServiceBaseUrl}/api/v1/claims/{claimId}");
                if (!response.IsSuccessStatusCode)
                {
                    throw new InvalidOperationException($"Failed to fetch claim {claimId}. Status: {response.StatusCode}");
                }

                var json = await response.Content.ReadAsStringAsync();
                var external = JsonSerializer.Deserialize<ExternalServiceResponse<ClaimInfoExternal>>(json);

                return external.Data.ToInternal();
            }
            catch (Exception ex)
            {
                throw new InvalidOperationException(ex.Message);
            }
        }

        public async Task<List<ClaimItem>> GetClaimItemsAsync(Guid claimId)
        {
            try
            {
                var response = await _httpClient.GetAsync($"{ExternalServiceBaseUrl}/api/v1/claims/{claimId}/items");

                if (!response.IsSuccessStatusCode)
                {
                    throw new InvalidOperationException($"Failed to fetch claim items for {claimId}. Status: {response.StatusCode}");
                }

                var json = await response.Content.ReadAsStringAsync();
                var external = JsonSerializer.Deserialize<ExternalServiceResponse<List<ClaimItemExternal>>>(json);

                return external.Data.ToInternal();
            }
            catch (Exception ex)
            {
                throw new InvalidOperationException(ex.Message);
            }
        }

        public async Task<TechnicianInfo> GetTechnicianAsync(Guid technicianId)
        {
            try
            {
                var response = await _httpClient.GetAsync($"{ExternalServiceBaseUrl}/api/v1/users/{technicianId}");

                if (!response.IsSuccessStatusCode)
                {
                    throw new InvalidOperationException($"Failed to fetch claim items for {technicianId}. Status: {response.StatusCode}");
                }

                var json = await response.Content.ReadAsStringAsync();
                var external = JsonSerializer.Deserialize<ExternalServiceResponse<TechnicianInfoExternal>>(json);

                return external.Data.ToInternal();
            }
            catch (Exception ex)
            {
                throw new InvalidOperationException(ex.Message);
            }
        }
    }
}
