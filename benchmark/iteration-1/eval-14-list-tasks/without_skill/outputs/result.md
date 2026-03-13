Two bugs still need fixing:

1. **001 - Fix login SSO bug** (high priority, **in-progress**)
   - SAML response validation is missing in the callback handler, causing intermittent SSO login failures
   - Remaining work: add SAML validation, handle expired sessions, add error logging

2. **003 - Patch XSS vulnerability in comments** (critical priority, **pending**)
   - Reflected XSS in the comments section — user input isn't sanitized before DOM rendering
   - All work remaining: sanitize input, add CSP headers, write regression tests

The XSS patch (003) is critical priority and hasn't been started yet — likely the most urgent item.
