# Architecture Decision: Session Handling Strategy

## Context

We are implementing `SessionRepository` and `SessionService` in a Clean Architecture structure.
A key question arose: Should `SessionRepository` preload (JOIN) `User` data when fetching a session?

## Decision

**We decided NOT to preload `User` data in `SessionRepository`.**

Instead, we perform two separate queries: one to fetch the Session, and another to fetch the User (if needed) in the Service layer.

## Trade-offs & Rationale

### Why NO Preload (Selected Approach)?

- **Decoupling (Clean Architecture):** Keeps `SessionRepository` focused solely on Session data. It doesn't need to know about User entity structure or database tables.
- **Scalability:** Allows `User` and `Session` data to potentially live in different databases or microservices in the future. `JOIN` queries would break in a distributed environment.
- **Caching:** Enables easier and more granular caching. User data (which changes less often) can be cached independently of Session data.
- **Performance:** While it incurs a "2nd query" penalty, in high-performance databases like Postgres with proper indexing (Primary Key), this latency is negligible for our current scale.

### Why Preload (Rejected Approach)?

- **Pros:** Reduces network round-trips (1 query vs 2). Convenient for simple use cases.
- **Cons:** Tightly couples Session and User repositories. Harder to refactor or split services later. Can fetch unnecessary data (User info) when only Session validation is needed.

## Future Considerations

If performance becomes a bottleneck or complex reporting is needed (e.g., "List sessions with user details"), we can implement a specific **Read Model** (CQRS pattern) query that performs the JOIN, separate from the core domain repository logic.
