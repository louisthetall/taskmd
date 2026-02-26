# @taskmd/web

The taskmd web frontend, built with React, TypeScript, and Vite.

## Development

```bash
pnpm dev       # Start dev server on localhost:5173
pnpm build     # Type-check and build for production
```

## Testing

```bash
pnpm test              # Run test suite
pnpm test:coverage     # Run tests with coverage reporting
```

After running `pnpm test:coverage`, open `coverage/index.html` in a browser to view the HTML coverage report.

Coverage thresholds are configured in `vitest.config.ts` and enforced in CI. If coverage drops below the thresholds, both the local command and CI will fail.
