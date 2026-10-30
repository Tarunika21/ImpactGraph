
import { afterEach, describe, expect, it, vi } from 'vitest'
import { apiRequest, ApiError } from './client'

afterEach(() => {
  vi.unstubAllGlobals()
})

describe('apiRequest', () => {
  it('retrieves JSON from a successful GET request', async () => {
    const fetchMock = vi.fn().mockResolvedValue({
      ok: true,
      status: 200,
      json: async () => ({ status: 'healthy' }),
    })

    vi.stubGlobal('fetch', fetchMock)

    const result = await apiRequest<{ status: string }>('/health')

    expect(result.status).toBe('healthy')
    expect(fetchMock).toHaveBeenCalledWith(
      '/health',
      expect.objectContaining({ headers: {} }),
    )
  })

  it('serializes POST request bodies as JSON', async () => {
    const fetchMock = vi.fn().mockResolvedValue({
      ok: true,
      status: 200,
      json: async () => ({ reportId: 42 }),
    })

    vi.stubGlobal('fetch', fetchMock)

    const result = await apiRequest<{ reportId: number }>('/api/v1/analyze', {
      method: 'POST',
      body: { serviceName: 'user-service' },
    })

    expect(result.reportId).toBe(42)

    expect(fetchMock).toHaveBeenCalledWith(
      '/api/v1/analyze',
      expect.objectContaining({
        method: 'POST',
        body: JSON.stringify({ serviceName: 'user-service' }),
        headers: {
          'Content-Type': 'application/json',
        },
      }),
    )
  })

  it('throws ApiError for unsuccessful requests', async () => {
    vi.stubGlobal(
      'fetch',
      vi.fn().mockResolvedValue({
        ok: false,
        status: 500,
        text: async () => 'Database unavailable',
      }),
    )

    await expect(apiRequest('/api/v1/reports')).rejects.toMatchObject({
      name: 'ApiError',
      status: 500,
      message: 'Database unavailable',
    })
  })

  it('handles empty HTTP 204 responses', async () => {
    vi.stubGlobal(
      'fetch',
      vi.fn().mockResolvedValue({
        ok: true,
        status: 204,
      }),
    )

    await expect(apiRequest<void>('/api/example')).resolves.toBeUndefined()
  })

  it('passes AbortSignal to fetch', async () => {
    const controller = new AbortController()

    const fetchMock = vi.fn().mockResolvedValue({
      ok: true,
      status: 200,
      json: async () => ({ reports: [] }),
    })

    vi.stubGlobal('fetch', fetchMock)

    await apiRequest('/api/v1/reports', {
      signal: controller.signal,
    })

    expect(fetchMock).toHaveBeenCalledWith(
      '/api/v1/reports',
      expect.objectContaining({
        signal: controller.signal,
      }),
    )
  })

  it('creates an ApiError with the supplied status', () => {
    const error = new ApiError(404, 'Report not found')

    expect(error).toBeInstanceOf(Error)
    expect(error.status).toBe(404)
    expect(error.message).toBe('Report not found')
  })
})
