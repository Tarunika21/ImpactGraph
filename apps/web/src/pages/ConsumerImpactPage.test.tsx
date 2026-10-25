
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import ConsumerImpactPage from './ConsumerImpactPage'

describe('ConsumerImpactPage', () => {
  beforeEach(() => {
    vi.stubGlobal('fetch', vi.fn())
  })

  it('renders the analysis form', () => {
    render(<ConsumerImpactPage />)

    expect(screen.getByDisplayValue('user-service')).toBeInTheDocument()
    expect(screen.getByDisplayValue('1.0.0')).toBeInTheDocument()
    expect(screen.getByDisplayValue('2.0.0')).toBeInTheDocument()
    expect(
      screen.getByRole('button', { name: 'Run Impact Analysis' }),
    ).toBeInTheDocument()
  })

  it('submits an analysis and displays the saved report ID', async () => {
    const user = userEvent.setup()
    const fetchMock = vi.mocked(fetch)

    fetchMock.mockResolvedValue({
      ok: true,
      json: async () => ({
        reportId: 42,
        comparison: {
          totalChanges: 1,
          breakingCount: 1,
          changes: [
            {
              type: 'FIELD_REMOVED',
              endpoint: 'GET /users/{id}',
              field: 'name',
              breaking: true,
            },
          ],
        },
        impact: {
          totalConsumers: 1,
          impactedCount: 1,
          impactedConsumers: [
            {
              consumerId: 'frontend-app',
              consumerName: 'Frontend Application',
              endpoint: 'GET /users/{id}',
              reasons: ['FIELD_REMOVED: name'],
            },
          ],
        },
      }),
    } as Response)

    render(<ConsumerImpactPage />)

    await user.click(
      screen.getByRole('button', { name: 'Run Impact Analysis' }),
    )

    await waitFor(() => {
      expect(screen.getByText('Report #42')).toBeInTheDocument()
    })

    expect(screen.getByText(/Frontend Application/)).toBeInTheDocument()

    expect(fetchMock).toHaveBeenCalledWith(
      '/api/v1/analyze',
      expect.objectContaining({
        method: 'POST',
      }),
    )
  })

  it('displays an API error', async () => {
    const user = userEvent.setup()
    const fetchMock = vi.mocked(fetch)

    fetchMock.mockResolvedValue({
      ok: false,
      status: 500,
      text: async () => 'Analysis unavailable',
    } as Response)

    render(<ConsumerImpactPage />)

    await user.click(
      screen.getByRole('button', { name: 'Run Impact Analysis' }),
    )

    expect(await screen.findByRole('alert')).toHaveTextContent(
      'Analysis failed (500)',
    )
  })
})
