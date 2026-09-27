---
name: receiving-code-review
description: Use when receiving code review feedback before changing implementation. Verify suggestions against the current codebase, requirements, tests, and architecture instead of applying feedback blindly.
---

# Receiving Code Review

Code review feedback is input to evaluate, not an automatic instruction to modify code.

The goal is to understand the intended outcome, verify the suggestion against the actual codebase, and apply only changes that improve correctness, maintainability, or required behavior.

## Workflow

When receiving review feedback:

1. Read the entire feedback before changing code.
2. Understand the intended outcome behind each comment.
3. Inspect the relevant code and surrounding context.
4. Verify the suggestion against:
   - current behavior;
   - requirements;
   - tests;
   - existing architecture;
   - existing abstractions and conventions.
5. Decide whether the proposed change is appropriate.
6. Implement the smallest correct change.
7. Run the relevant tests after implementation.

Do not modify code solely because a reviewer suggested a particular implementation.

## Verify before implementing

Before accepting a suggestion, check whether:

- the issue actually exists;
- the proposed change solves the real problem;
- the change breaks existing behavior;
- the change duplicates existing functionality;
- the current implementation exists for a reason not visible in the review comment;
- the suggestion conflicts with requirements, tests, or architectural decisions.

Search the repository when necessary instead of assuming the reviewer's local context is complete.

## Feedback from the user or project developer

Treat feedback from the user or another developer as a strong signal about the intended outcome, but not as proof that the proposed implementation is technically correct.

Before changing code:

1. Understand what behavior or problem they want addressed.
2. Verify the proposed implementation against the current codebase.
3. Prefer the requested implementation when it is correct and appropriate.
4. If the requested implementation introduces unnecessary complexity, breaks existing behavior, or conflicts with established architecture, explain the issue and use a simpler correct solution when possible.
5. Ask for clarification only when the intended behavior itself is materially ambiguous.

Do not silently reinterpret clear product requirements.

## External review feedback

Treat feedback from external reviewers, automated review tools, or other agents as hypotheses to verify.

Do not assume that the reviewer:

- understands the entire codebase;
- knows all project requirements;
- saw every relevant dependency;
- understands why the current implementation exists.

Validate the feedback against the repository before changing production code.

## Complexity check

Do not accept a suggestion solely because it sounds more:

- architectural;
- generic;
- extensible;
- defensive;
- production-ready;
- abstract;
- future-proof.

Before introducing a new:

- interface;
- abstraction;
- layer;
- configuration option;
- fallback;
- validation branch;
- helper;
- dependency;
- design pattern;

identify the concrete current problem it solves.

If the current implementation is simpler and already satisfies the required behavior, tests, and architecture, prefer keeping it.

Do not build infrastructure for hypothetical future requirements.

## YAGNI check

If review feedback suggests supporting behavior that is not currently required:

1. Search the repository for actual usage or requirements.
2. Check business requirements or system design when relevant.
3. If no current requirement exists, do not implement speculative support unless explicitly requested.

Prefer solving the present problem.

## Unclear feedback

Do not guess when a review comment is materially ambiguous.

If unclear feedback affects other review items or changes the intended behavior, clarify it before implementing dependent changes.

If review items are independent:

- proceed with the clear items;
- leave the unclear item unchanged;
- identify what remains unresolved.

Do not block unrelated fixes unnecessarily.

## Applying feedback

When a suggestion is valid:

- make the smallest change that addresses the issue;
- preserve existing behavior unless the review explicitly requires changing it;
- avoid unrelated refactoring;
- reuse existing code where possible;
- follow existing project conventions;
- add or update tests when behavior changes.

If the change modifies behavior, follow the `test-driven-development` skill.

## When pushing back

Push back when the proposed change:

- is technically incorrect;
- contradicts project requirements;
- breaks existing behavior;
- duplicates existing functionality;
- introduces unnecessary complexity;
- solves a hypothetical rather than current problem;
- weakens testability or dependency boundaries;
- conflicts with an established architectural decision without a clear benefit.

Base pushback on concrete technical reasoning.

Prefer statements such as:

- "The current implementation already handles this case because..."
- "This would introduce an additional abstraction without a current second use case."
- "This conflicts with the existing requirement that..."
- "The suggested change would break..."
- "A smaller change addresses the issue without introducing..."

Avoid performative agreement. Respond with technical reasoning, clarification, or the implemented change.

## Review order

For multiple review comments:

1. Understand all comments first.
2. Identify dependencies between them.
3. Handle correctness and behavioral issues before stylistic improvements.
4. Apply independent changes separately where practical.
5. Run relevant tests after meaningful changes.

Do not blindly implement comments one by one without understanding how they interact.

## Before finishing

Check:

- Was every implemented review comment verified against the codebase?
- Did any change introduce unnecessary complexity?
- Did we preserve unrelated behavior?
- Did we avoid duplicating existing functionality?
- Were behavior changes tested?
- Did we leave unclear or questionable feedback unresolved rather than guessing?
- Is the resulting code simpler or at least no more complex than necessary?
