Two bugs still need fixing:

1. **[001] Fix login SSO bug** — `in-progress`, **high** priority
   - SAML response validation in the callback handler is broken, causing intermittent SSO login failures
   - Remaining: add SAML validation, handle expired sessions, add error logging

2. **[003] Patch XSS vulnerability in comments** — `pending`, **critical** priority
   - Reflected XSS in the comments section — user input not sanitized before DOM rendering
   - Remaining: sanitize input, add CSP headers, write regression tests

The XSS patch (003) is critical and hasn't been started yet. Want me to work on either of these?
