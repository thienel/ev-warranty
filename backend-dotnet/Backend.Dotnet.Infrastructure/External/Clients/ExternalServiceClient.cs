using Backend.Dotnet.Application.Interfaces.External;
using Backend.Dotnet.Infrastructure.External.Models;
using System.Text.Json;
using static Backend.Dotnet.Application.DTOs.IntermediaryDto;
using static Backend.Dotnet.Infrastructure.External.Models.ExternalModel;

namespace Backend.Dotnet.Infrastructure.External.Clients
{
    public class ExternalServiceClient : IExternalServiceClient
    {
        private readonly HttpClient _httpClient;
        private const string ExternalServiceBaseUrl = "http://localhost:8080";

        public ExternalServiceClient(HttpClient httpClient)
        {
            _httpClient = httpClient;
        }

        public async Task<bool> CompleteClaimAsync(Guid claimId)
        {
            try
            {
                var response = await _httpClient.PostAsync($"{ExternalServiceBaseUrl}/claims/{claimId}/complete", null);

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
                var response = await _httpClient.GetAsync($"{ExternalServiceBaseUrl}/claims/{claimId}");
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
                var response = await _httpClient.GetAsync($"{ExternalServiceBaseUrl}/claims/{claimId}/items");

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
                var response = await _httpClient.GetAsync($"{ExternalServiceBaseUrl}/users/{technicianId}");

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

        public async Task<bool> ValidateClaimStatusAsync(Guid claimId)
        {
            throw new NotImplementedException();
        }

        public Task<bool> ValidateTechnicianRoleAsync(Guid technicianId)
        {
            throw new NotImplementedException();
        }
    }
}
