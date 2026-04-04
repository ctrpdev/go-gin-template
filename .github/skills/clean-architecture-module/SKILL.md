---
name: clean-architecture-module
description: 'Guarantees standard module implementation for Clean Architecture in Go and Gin. Use when creating or completing new domains, like User or Note.'
argument-hint: 'Module name to scaffold or review'
user-invocable: true
disable-model-invocation: false
---

# Clean Architecture Module Scaffolding (Go + Gin)

## When to Use

- When the user creates a new Domain (e.g., `Note`, `Product`).
- When a service or handler has missing/empty methods and needs to be standardized.
- When scaffolding DTOs or modifying routes for an existing domain to match standard architecture.

## Architecture Guidelines (Based on `User` module)

When building or updating an endpoint, follow these standard layers across the `internal/` directory:

1. **Domain (`domain/`)**: Ensure the basic interfaces (`<Entity>Repository`, `<Entity>Service`) and struct (`<Entity>`, embedding `BaseModel`) exist.
2. **DTO (`handler/http/dto/`)**:
   - Every input payload must have a request DTO (e.g. `CreateNoteRequest`) with proper validation tags (`binding:"required"`).
   - Use `omitempty` for partial updates in PUT/PATCH requests.
   - Always map the database output to a response DTO where appropriate.
3. **Handler (`handler/http/`)**:
   - `New<Entity>Handler(service domain.<Entity>Service) *<Entity>Handler`
   - Use `c.GetInt64("userID")` (or equivalent from Auth context) to link resources to the authenticated owner.
   - Validate payload via `c.ShouldBindJSON(&req)`.
   - Propagate errors to standardized HTTP formatting using `errors.MapDomainError(c, err)`.
4. **Service (`service/`)**:
   - Implements `domain.<Entity>Service`.
   - Injects `domain.<Entity>Repository`.
   - Orchestrates logic and returns either domain models or domain errors.
5. **Router (`routes/`)**:
   - Define a dedicated file `<entity>_routes.go`.
   - Encapsulate the group in `Setup<Entity>Routes(rg *gin.RouterGroup, handler *http.<Entity>Handler, auth *middleware.AuthMiddleware)`.
   - Route standard protections via `authMiddleware.RequireAuth()`.
6. **Server Injection (`server/server.go`)**:
   - Wire up the pieces and provide them to `routes.SetupRouter(..., <entity>Handler, ...)`.

## Workflow Process

When the user asks you to apply this skill (e.g., `/clean-architecture-module Note`):

1. **Verify Domain**: Check the `.go` domain file inside `internal/domain`. Read the declared interfaces.
2. **Check the Layers**: Look into `dto`, `handler`, `service`, `repository/postgres`, and `routes`. Identify which implementations are empty or missing.
3. **Fill DTO**: Create the necessary Request/Response JSON bindings.
4. **Fill Controller (Handler)**: Write out Gin methods using the service and mapping errors using the `errors` package mapping standard.
5. **Establish Route**: Create or update the route grouping file `internal/routes/<entity>_routes.go`.
6. **Wire Dependencies**: Update `internal/routes/router.go` and `internal/server/server.go` with the new instances.
7. **Report**: Return a bullet list of files touched and why.

## Code Quality Check
- Is `binding:"required"` present in POST/create DTOs?
- Did you handle `err != nil` directly with `errors.MapDomainError(c, err)`?
- Are IDs received by param appropriately cast with `strconv.ParseInt(idParam, 10, 64)`?
