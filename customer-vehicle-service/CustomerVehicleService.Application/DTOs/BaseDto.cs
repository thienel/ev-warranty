using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Text.Json.Serialization;
using System.Threading.Tasks;

namespace CustomerVehicleService.Application.DTOs
{
    public class BaseResponseDto
    {
        [JsonPropertyName("is_success")]
        public bool IsSuccess { get; set; }

        [JsonPropertyName("message")]
        public string Message { get; set; } = string.Empty;

        [JsonPropertyName("error_code")]
        public string? ErrorCode { get; set; }
    }

    public class BaseResponseDto<T> : BaseResponseDto where T : class
    {
        public T? Data { get; set; } = null;
    }
}
