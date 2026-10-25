
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import App from './App'

const mockReports = [
  {
    id: 1,
    serviceName: 'user-service',
    oldVersion: '1.0.0',
    newVersion: '2.0.0',
    createdAt: '2026-10-24T10:00:00Z',
    result: {
      comparison: {
        totalChanges: 1,
        breakingCount: 1,
        changes: [],
      },
      impact: {
        totalConsumers: 1,
        impactedCount: 1,
        impactedConsumers: [],
      },
    },
  },
  {
    id: 2,
    serviceName: 'payment-service',
    oldVersion: '1.0.0',
    newVersion: '1.1.0',
    createdAt: '2026-10-25T10:00:00Z',
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
]

describe('ImpactGraph App', () => {
  beforeEach(() => {
    vi.stubGlobal(
      'fetch',
      vi.fn(async (input: RequestInfo | URL) => {
        const url = String(input)

        if (url === '/health') {
          return {
            ok: true,
            json: async () => ({
              status: 'healthy',
              service: 'impactgraph-api',
            }),
          } as Response
        }

        if (url === '/api/v1/reports') {
          return {
            ok: true,
            json: async () => mockReports,
          } as Response
        }

        return {
          ok: false,
          status: 404,
        } as Response
      }),
    )
  })

  afterEach(() => {
    vi.unstubAllGlobals()
  })

  it('renders the Overview dashboard initially', async () => {
    render(<App />)

    expect(
      screen.getByRole('heading', { name: 'Overview' }),
    ).toBeInTheDocument()

    expect(
      screen.getByText('Monitor API contract changes and their downstream impact.'),
    ).toBeInTheDocument()

    expect(await screen.findByText('Go API Connected')).toBeInTheDocument()
  })

  it('displays dashboard statistics from saved reports', async () => {
    render(<App />)

    await waitFor(() => {
      expect(screen.getByText('Saved Reports')).toBeInTheDocument()
    })

    expect(screen.getByText('Services Analyzed')).toBeInTheDocument()
    expect(screen.getByText('Breaking Changes Found')).toBeInTheDocument()
    expect(screen.getAllByText('Impacted Consumers').length).toBeGreaterThan(0)

    expect(screen.getByText('user-service')).toBeInTheDocument()
    expect(screen.getByText('payment-service')).toBeInTheDocument()
  })

  it('opens API Comparison from the sidebar', async () => {
    const user = userEvent.setup()
    render(<App />)

    await user.click(
      screen.getByRole('button', { name: /API Comparison/ }),
    )

    expect(
      screen.getByRole('heading', { name: 'API Comparison' }),
    ).toBeInTheDocument()

    expect(
      screen.getByRole('button', { name: 'Compare Specifications' }),
    ).toBeInTheDocument()
  })

  it('opens API Comparison from New Analysis', async () => {
    const user = userEvent.setup()
    render(<App />)

    await user.click(
      screen.getByRole('button', { name: /New Analysis/ }),
    )

    expect(
      screen.getByRole('heading', { name: 'API Comparison' }),
    ).toBeInTheDocument()
  })

  it('opens Consumer Impact from the sidebar', async () => {
    const user = userEvent.setup()
    render(<App />)

    await user.click(
      screen.getByRole('button', { name: /Consumer Impact/ }),
    )

    expect(
      screen.getByRole('heading', { name: 'Consumer Impact' }),
    ).toBeInTheDocument()

    expect(
      screen.getByRole('button', { name: 'Run Impact Analysis' }),
    ).toBeInTheDocument()
  })

  it('opens Report History from the sidebar', async () => {
    const user = userEvent.setup()
    render(<App />)

    await user.click(
      screen.getByRole('button', { name: /Report History/ }),
    )

    expect(
      screen.getByRole('heading', { name: 'Report History' }),
    ).toBeInTheDocument()

    expect(await screen.findByText('user-service')).toBeInTheDocument()
  })

  it('shows an error when the report API fails', async () => {
    vi.mocked(fetch).mockImplementation(
      async (input: RequestInfo | URL) => {
        if (String(input) === '/health') {
          return { ok: true } as Response
        }

        return {
          ok: false,
          status: 500,
        } as Response
      },
    )

    render(<App />)

    expect(await screen.findByText('Report API returned 500')).toBeInTheDocument()
  })
})
