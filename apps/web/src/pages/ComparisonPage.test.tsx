
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import ComparisonPage from './ComparisonPage'

describe('ComparisonPage', () => {
  beforeEach(() => {
    vi.stubGlobal('fetch', vi.fn())
  })

  it('renders both API specification editors', () => {
    render(<ComparisonPage />)

    expect(
      screen.getByRole('heading', { name: 'API Comparison' }),
    ).toBeInTheDocument()

    expect(
      screen.getByRole('textbox', { name: 'Old API specification' }),
    ).toBeInTheDocument()

    expect(
      screen.getByRole('textbox', { name: 'New API specification' }),
    ).toBeInTheDocument()
  })

  it('rejects invalid JSON without calling the API', async () => {
    const user = userEvent.setup()
    const fetchMock = vi.mocked(fetch)

    render(<ComparisonPage />)

    const oldEditor = screen.getByRole('textbox', {
      name: 'Old API specification',
    })

    await user.clear(oldEditor)
    await user.type(oldEditor, 'invalid json')

    await user.click(
      screen.getByRole('button', { name: 'Compare Specifications' }),
    )

    expect(
      screen.getByRole('alert'),
    ).toHaveTextContent('Both API specifications must contain valid JSON.')

    expect(fetchMock).not.toHaveBeenCalled()
  })

  it('displays breaking changes returned by the Go API', async () => {
    const user = userEvent.setup()
    const fetchMock = vi.mocked(fetch)

    fetchMock.mockResolvedValue({
      ok: true,
      json: async () => ({
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
      }),
    } as Response)

    render(<ComparisonPage />)

    await user.click(
      screen.getByRole('button', { name: 'Compare Specifications' }),
    )

    await waitFor(() => {
      expect(
        screen.getByRole('heading', { name: 'Comparison Results' }),
      ).toBeInTheDocument()
    })

    expect(screen.getByText('Breaking')).toBeInTheDocument()
    expect(screen.getByText(/FIELD_REMOVED/)).toBeInTheDocument()

    expect(fetchMock).toHaveBeenCalledWith(
      '/api/v1/compare',
      expect.objectContaining({
        method: 'POST',
      }),
    )
  })

  it('shows an error when the API request fails', async () => {
    const user = userEvent.setup()
    const fetchMock = vi.mocked(fetch)

    fetchMock.mockResolvedValue({
      ok: false,
      status: 500,
      text: async () => 'Internal Server Error',
    } as Response)

    render(<ComparisonPage />)

    await user.click(
      screen.getByRole('button', { name: 'Compare Specifications' }),
    )

    expect(await screen.findByRole('alert')).toHaveTextContent(
      'Comparison failed (500)',
    )
  })
})
