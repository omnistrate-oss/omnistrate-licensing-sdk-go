# Copilot Instructions for omnistrate-licensing-sdk-go

<!-- security-checklist-managed -->

## Security Checklist

This SDK validates Omnistrate-issued license files. A bypass here lets a customer run unlicensed software, and a panic/DoS here takes down the consuming application. Apply this checklist to every change.

### Signature & Trust
- Signature verification is **mandatory** for every license. There must be no code path — including "debug", "test", or "dev" modes — that accepts an unsigned or invalidly-signed license at runtime.
- Use only vetted crypto from the standard library or well-known libraries (e.g., `crypto/ecdsa`, `crypto/ed25519`, `java.security.Signature`). Do not implement custom signature schemes.
- Pin the trusted public key(s) at build time or load from a verified, integrity-protected source. Do not fetch trust roots over plain HTTP.
- Reject licenses signed with deprecated algorithms (e.g., MD5, SHA-1, RSA < 2048).

### License Validation
- Enforce `notBefore` / `notAfter` (with documented clock-skew tolerance, typically ≤ 5 minutes).
- Validate every claim that gates entitlement (product, edition, feature flags, instance ID, tenant ID). A missing claim must default to deny, not allow.
- Treat license replay across instances/tenants as a failure mode: bind licenses to a tenant or installation identifier and validate it.
- Ensure deterministic rejection: the same invalid license must always be rejected, regardless of caller, locale, or environment.

### Input Handling & Robustness
- Treat the license blob as untrusted input. Cap its size before parsing. Use safe parsers (no arbitrary type construction, no XXE in any XML path).
- Never `panic` / throw an unchecked exception on a malformed license — return a typed error so callers can fail closed gracefully.
- Use the crypto library's signature **verify** API (`ecdsa.Verify`, `ed25519.Verify`, `Signature.verify`) — do not implement signature checking via byte comparison. Reserve constant-time comparison (`subtle.ConstantTimeCompare`, `MessageDigest.isEqual`) for symmetric MAC / token equality checks.

### Secrets Handling
- The SDK validates **public** keys; it should never need a private key. If a private key ever appears in this repo (tests, fixtures, examples), it must be a clearly-marked dummy key not used in any production trust chain.
- No real customer license files in fixtures.

### Logging Hygiene
- Do not log full license contents or signatures at any level. Log only the validation outcome, license ID (if present), and a short reason code for failures.
- Do not include the license blob in returned error messages.

### Dependencies & Supply Chain
- Keep the dependency surface minimal — this SDK ships into customer environments.
- Run `govulncheck` (Go) / OWASP Dependency-Check (Java) on dependency changes.
- Pin direct dependency versions. Review every transitive bump before release.
- Releases must be signed and published with provenance (Maven Central with PGP signature for Java; signed git tags for Go).
- Pin GitHub Actions to commit SHAs.

### API Stability
- Do not silently change validation behavior (which claims are required, which algorithms are accepted, default clock skew). Such changes are breaking and must be called out in release notes.

### What to do when unsure
- If a change touches signature verification, claim validation, trust root loading, or default-deny behavior, call it out explicitly in the PR description and request a security review.
