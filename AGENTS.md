# AGENTS.md

## Project

This repository contains a Go application for collecting and processing marketplace listings and interacting with users through Telegram.

Keep changes consistent with the existing architecture and project conventions.

Use documentation in `docs/` when it is relevant to the task:

- use business requirements when changing product behavior;
- use system design when changing service boundaries, data flow, persistence, concurrency, or integrations.

Do not read unrelated documentation for every small change.

## Engineering principles

- Prefer simple, explicit, readable code over clever abstractions.
- Implement only what is required by the current task.
- Keep production code minimal: do not add defensive branches unless required by actual behavior, requirements, or a test.
- Do not design extension points, abstractions, configuration, or fallback behavior solely for hypothetical future requirements.
- Prefer a small amount of straightforward duplication over premature abstraction.
- Keep changes focused; avoid unrelated refactoring.

## Before implementing

- Inspect the relevant existing code and understand the current pattern.
- Search the repository before creating a new function, type, helper, service, or abstraction.
- Reuse or extend existing functionality when doing so remains simple.
- Do not introduce a second implementation of behavior that already exists.

## Code structure

- Keep functions focused on a clear responsibility.
- Split a function when distinct responsibilities can be meaningfully named, understood, or tested independently.
- Do not split trivial code merely to make functions shorter.
- Do not place unrelated responsibilities in the same file.
- Prefer descriptive names and straightforward control flow over explanatory comments and unnecessary indirection.
- Code should be easy to understand visually and by reading it top-to-bottom.

## Dependencies

- Keep business logic independent from external systems where practical.
- Isolate external boundaries such as HTTP APIs, PostgreSQL, filesystem, clocks, queues, and third-party services.
- Use dependency injection when it provides a concrete benefit for separation or testing.
- Prefer explicit constructor or parameter injection.
- Do not introduce interfaces or dependency-injection infrastructure solely for architectural purity.

## Testing

- Follow TDD when implementing new behavior or fixing behavior-related bugs:
  failing test -> minimal implementation -> refactor.
- Test observable behavior rather than implementation details.
- Mock or fake external boundaries, not internal application logic.
- Prefer a small number of meaningful tests over tests written only to increase coverage.

## Finishing a change

- Run the relevant tests after implementation.
- Run broader tests when the scope of the change makes them appropriate.
- Check that the implementation does not duplicate existing functionality.
- Remove unnecessary complexity introduced during implementation.
- Update documentation only when the behavior, architecture, or usage it describes has actually changed.
