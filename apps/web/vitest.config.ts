import { defineConfig } from "vitest/config";
import react from "@vitejs/plugin-react";

export default defineConfig({
  plugins: [react()],
  test: {
    environment: "jsdom",
    setupFiles: ["./src/test-setup.ts"],
    coverage: {
      provider: "v8",
      reporter: ["text", "html", "json-summary", "lcov"],
      include: ["src/**/*.{ts,tsx}"],
      exclude: ["src/test-setup.ts", "src/**/*.test.{ts,tsx}", "src/main.tsx", "src/test-utils/**"],
      thresholds: {
        lines: 5,
        branches: 25,
        functions: 20,
        statements: 5,
      },
    },
  },
});
