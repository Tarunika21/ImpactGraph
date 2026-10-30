import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { render, screen, waitFor, within } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import ReportHistoryPage from "./ReportHistoryPage";

const reports = [
  {
    id: 3,
    serviceName: "user-service",
    oldVersion: "1.0.0",
    newVersion: "2.0.0",
    createdAt: "2026-10-24T15:36:00Z",
    result: {
      comparison: {
        totalChanges: 1,
        breakingCount: 1,
        changes: [
          {
            type: "FIELD_REMOVED",
            endpoint: "GET /users/{id}",
            field: "name",
            breaking: true,
          },
        ],
      },
      impact: {
        totalConsumers: 1,
        impactedCount: 1,
        impactedConsumers: [
          {
            consumerId: "frontend-app",
            consumerName: "Frontend Application",
            endpoint: "GET /users/{id}",
            reasons: ["FIELD_REMOVED: name"],
          },
        ],
      },
    },
  },
  {
    id: 4,
    serviceName: "payment-service",
    oldVersion: "1.0.0",
    newVersion: "1.1.0",
    createdAt: "2026-10-25T10:00:00Z",
    result: {
      comparison: {
        totalChanges: 0,
        breakingCount: 0,
        changes: [],
      },
      impact: {
        totalConsumers: 0,
        impactedCount: 0,
        impactedConsumers: [],
      },
    },
  },
];

describe("ReportHistoryPage", () => {
  beforeEach(() => {
    vi.stubGlobal(
      "fetch",
      vi.fn(async (input: RequestInfo | URL) => {
        const url = String(input);

        if (url === "/api/v1/reports") {
          return {
            ok: true,
            json: async () => reports,
          } as Response;
        }

        if (url === "/api/v1/reports/3") {
          return {
            ok: true,
            json: async () => reports[0],
          } as Response;
        }

        if (url === "/api/v1/reports/4") {
          return {
            ok: true,
            json: async () => reports[1],
          } as Response;
        }

        return {
          ok: false,
          status: 404,
        } as Response;
      }),
    );
  });

  afterEach(() => {
    vi.unstubAllGlobals();
  });

  it("loads and displays saved reports", async () => {
    render(<ReportHistoryPage />);

    expect(await screen.findByText("user-service")).toBeInTheDocument();
    expect(screen.getByText("payment-service")).toBeInTheDocument();
    expect(screen.getByText("2 reports")).toBeInTheDocument();
  });

  it("filters reports by service name", async () => {
    const user = userEvent.setup();
    render(<ReportHistoryPage />);

    await screen.findByText("user-service");

    await user.type(
      screen.getByRole("searchbox", { name: "Search by service" }),
      "payment",
    );

    expect(screen.getByText("payment-service")).toBeInTheDocument();
    expect(screen.queryByText("user-service")).not.toBeInTheDocument();
    expect(screen.getByText("1 reports")).toBeInTheDocument();
  });

  it("opens and closes report details", async () => {
    const user = userEvent.setup();
    render(<ReportHistoryPage />);

    const reportCard = (await screen.findByText("user-service")).closest(
      "article",
    );
    if (!reportCard) throw new Error("Report card not found");

    await user.click(
      within(reportCard).getByRole("button", { name: "View Details" }),
    );

    await waitFor(() => {
      expect(within(reportCard).getByText("FIELD_REMOVED")).toBeInTheDocument();
    });

    expect(
      within(reportCard).getByText("Frontend Application"),
    ).toBeInTheDocument();

    await user.click(
      within(reportCard).getByRole("button", { name: "Close Details" }),
    );

    expect(
      within(reportCard).queryByText("FIELD_REMOVED"),
    ).not.toBeInTheDocument();
  });

  it("shows an error when the report list fails", async () => {
    vi.mocked(fetch).mockResolvedValueOnce({
      ok: false,
      status: 500,
      text: async () => "",
    } as Response);

    render(<ReportHistoryPage />);

    expect(await screen.findByRole("alert")).toHaveTextContent(
      "Request failed with status 500",
    );
  });

  it("shows an error when report details fail", async () => {
    const user = userEvent.setup();
    vi.mocked(fetch).mockImplementation(async (input: RequestInfo | URL) => {
      if (String(input) === "/api/v1/reports") {
        return {
          ok: true,
          json: async () => reports,
        } as Response;
      }

      return {
        ok: false,
        status: 500,
        text: async () => "",
      } as Response;
    });

    render(<ReportHistoryPage />);

    const reportCard = (await screen.findByText("user-service")).closest(
      "article",
    );
    if (!reportCard) throw new Error("Report card not found");

    await user.click(
      within(reportCard).getByRole("button", { name: "View Details" }),
    );

    expect(await within(reportCard).findByRole("alert")).toHaveTextContent(
      "Request failed with status 500",
    );
  });
});
