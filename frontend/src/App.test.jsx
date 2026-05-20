import { render, screen, waitFor } from "@testing-library/react";
import { expect, test, vi } from "vitest";
import App from "./App.jsx";

test("renders the StartTech dashboard and API status", async () => {
  global.fetch = vi.fn(() =>
    Promise.resolve({
      ok: true,
      json: () =>
        Promise.resolve({
          service: "starttech-api",
          environment: "test",
          uptime_seconds: 1,
        }),
    }),
  );

  render(<App />);

  expect(screen.getByText("Full-stack delivery pipeline")).toBeInTheDocument();

  await waitFor(() => {
    expect(screen.getByText("API online")).toBeInTheDocument();
  });
});
