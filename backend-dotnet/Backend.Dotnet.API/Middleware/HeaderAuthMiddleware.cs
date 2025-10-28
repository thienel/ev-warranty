using System.Security.Claims;

namespace Backend.Dotnet.API.Middleware
{
    public class HeaderAuthMiddleware
    {
        private readonly RequestDelegate _next;

        public HeaderAuthMiddleware(RequestDelegate next)
        {
            _next = next;
        }

        public async Task InvokeAsync(HttpContext context)
        {
            var identity = new ClaimsIdentity("HeaderAuth");

            //if (context.Request.Headers.TryGetValue("X-User-ID", out var userId))
            //    identity.AddClaim(new Claim(ClaimTypes.NameIdentifier, userId!));

            if (context.Request.Headers.TryGetValue("X-User-Role", out var role))
                identity.AddClaim(new Claim(ClaimTypes.Role, role!));

            if (identity.Claims.Any())
            {
                context.User = new ClaimsPrincipal(identity);
            }

            await _next(context);
        }
    }
}
