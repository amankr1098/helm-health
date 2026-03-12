# Contributing to Helm Health

Thank you for your interest in contributing to **Helm Health**! Contributions of all kinds are welcome — bug reports, feature requests, documentation improvements, and code changes.

---

## Ground Rules

- **PRs are welcome**, but only maintainers can merge them.
- All changes must go through a pull request — direct pushes to `main` are not allowed.
- Every PR requires at least **one approval** from a maintainer before it can be merged.
- Keep PRs focused: one feature or fix per pull request.

---

## How to Contribute

### 1. Fork & Clone

```bash
# Fork via GitHub UI, then:
git clone https://github.com/<your-username>/helm-health.git
cd helm-health
```

### 2. Create a Branch

Always branch off `main`:

```bash
git checkout -b feat/my-feature
# or
git checkout -b fix/my-bugfix
```

**Branch naming conventions:**

| Prefix     | Use for                    |
|------------|----------------------------|
| `feat/`    | New features               |
| `fix/`     | Bug fixes                  |
| `docs/`    | Documentation changes      |
| `refactor/`| Code refactoring           |
| `test/`    | Adding or updating tests   |

### 3. Make Your Changes

- Follow existing code style and project structure.
- Ensure the project builds cleanly:

  ```bash
  go build ./...
  ```

- Run existing tests (if any):

  ```bash
  go test ./...
  ```

- Add or update documentation if your change affects behavior.

### 4. Commit Your Changes

Write clear, descriptive commit messages:

```
feat: add health check for CronJob resources

Adds CronJob support with checks for last schedule time
and active/suspended status.
```

### 5. Push & Open a Pull Request

```bash
git push origin feat/my-feature
```

Then open a PR against the `main` branch on GitHub. In your PR description:

- Explain **what** the change does and **why**.
- Reference any related issues (e.g., `Closes #12`).
- Include example output or screenshots if applicable.

---

## What Happens After You Open a PR

1. A maintainer will review your PR.
2. You may be asked to make changes — this is normal and collaborative.
3. Once approved, a **maintainer will merge** the PR. Contributors do not have merge access.

---

## Ideas for Contributions

Not sure where to start? Here are some ideas:

- Add `--watch` mode for continuous health monitoring
- Add configurable thresholds (e.g., restart count limits)
- Add health checks for more resource kinds (CronJob, HPA, ConfigMap, etc.)
- Improve error messages and edge-case handling
- Add unit tests for existing health check functions
- Helm test integration
- Support for custom resource health plugins

---

## Code of Conduct

Be respectful, constructive, and inclusive. We want this to be a welcoming project for everyone.

---

## Questions?

Open an issue or start a discussion on GitHub — happy to help!
