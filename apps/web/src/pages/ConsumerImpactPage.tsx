
import { useState } from 'react'
import './ConsumerImpactPage.css'

type Change = {
  type?: string
  endpoint?: string
  field?: string
  breaking?: boolean
}

type AnalysisResponse = {
  reportId: number
  comparison: {
    totalChanges?: number
    breakingCount?: number
    changes?: Change[]
  }
  impact: Record<string, unknown>
}

const oldExample = {
  paths: {
    '/users/{id}': {
      get: {
        responses: {
          '200': {
            content: {
              'application/json': {
                schema: {
                  type: 'object',
                  properties: {
                    id: { type: 'integer' },
                    name: { type: 'string' },
                    email: { type: 'string' },
                  },
                },
              },
            },
          },
        },
      },
    },
  },
}

const newExample = {
  paths: {
    '/users/{id}': {
      get: {
        responses: {
          '200': {
            content: {
              'application/json': {
                schema: {
                  type: 'object',
                  properties: {
                    id: { type: 'integer' },
                    email: { type: 'string' },
                  },
                },
              },
            },
          },
        },
      },
    },
  },
}

function ConsumerImpactPage() {
  const [serviceName, setServiceName] = useState('user-service')
  const [oldVersion, setOldVersion] = useState('1.0.0')
  const [newVersion, setNewVersion] = useState('2.0.0')
  const [oldSpec, setOldSpec] = useState(JSON.stringify(oldExample, null, 2))
  const [newSpec, setNewSpec] = useState(JSON.stringify(newExample, null, 2))
  const [result, setResult] = useState<AnalysisResponse | null>(null)
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  async function runAnalysis(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault()
    setError('')
    setResult(null)

    if (!serviceName.trim() || !oldVersion.trim() || !newVersion.trim()) {
      setError('Service name and both versions are required.')
      return
    }

    let parsedOld: unknown
    let parsedNew: unknown

    try {
      parsedOld = JSON.parse(oldSpec)
      parsedNew = JSON.parse(newSpec)
    } catch {
      setError('Both API specifications must contain valid JSON.')
      return
    }

    setLoading(true)

    try {
      const response = await fetch('/api/v1/analyze', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          serviceName: serviceName.trim(),
          oldVersion: oldVersion.trim(),
          newVersion: newVersion.trim(),
          oldSpec: parsedOld,
          newSpec: parsedNew,
        }),
      })

      if (!response.ok) {
        const message = await response.text()
        throw new Error(`Analysis failed (${response.status}): ${message}`)
      }

      const data: AnalysisResponse = await response.json()
      setResult(data)
    } catch (err) {
      setError(
        err instanceof Error
          ? err.message
          : 'Unable to complete impact analysis.',
      )
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="impact-page">
      <div className="page-heading">
        <div>
          <h1>Consumer Impact</h1>
          <p>Discover which consumers may be affected by API contract changes.</p>
        </div>
      </div>

      <form onSubmit={(event) => void runAnalysis(event)}>
        <div className="impact-fields">
          <label>
            Service Name
            <input
              value={serviceName}
              onChange={(event) => setServiceName(event.target.value)}
              placeholder="user-service"
              required
              maxLength={150}
            />
          </label>

          <label>
            Old Version
            <input
              value={oldVersion}
              onChange={(event) => setOldVersion(event.target.value)}
              placeholder="1.0.0"
              required
              maxLength={100}
            />
          </label>

          <label>
            New Version
            <input
              value={newVersion}
              onChange={(event) => setNewVersion(event.target.value)}
              placeholder="2.0.0"
              required
              maxLength={100}
            />
          </label>
        </div>

        <div className="impact-editors">
          <label className="impact-editor">
            <span>Old API Specification</span>
            <textarea
              value={oldSpec}
              onChange={(event) => setOldSpec(event.target.value)}
              spellCheck={false}
              required
            />
          </label>

          <label className="impact-editor">
            <span>New API Specification</span>
            <textarea
              value={newSpec}
              onChange={(event) => setNewSpec(event.target.value)}
              spellCheck={false}
              required
            />
          </label>
        </div>

        <div className="impact-actions">
          <button className="primary-button" type="submit" disabled={loading}>
            {loading ? 'Analyzing...' : 'Run Impact Analysis'}
          </button>
          <span>Successful analyses are saved to PostgreSQL.</span>
        </div>
      </form>

      {error && (
        <div className="error-banner" role="alert">
          {error}
        </div>
      )}

      {result && (
        <section className="impact-results">
          <div className="impact-results-header">
            <div>
              <h2>Analysis Results</h2>
              <p>Service: {serviceName}</p>
            </div>
            <span className="impact-report-id">Report #{result.reportId}</span>
          </div>

          <div className="impact-summary">
            <div>
              <small>Total Changes</small>
              <strong>{result.comparison?.totalChanges ?? 0}</strong>
            </div>
            <div>
              <small>Breaking Changes</small>
              <strong>{result.comparison?.breakingCount ?? 0}</strong>
            </div>
            <div>
              <small>Impact Analysis</small>
              <strong>Completed</strong>
            </div>
          </div>

          <h3>Detected Changes</h3>
          <pre className="impact-json">
            {JSON.stringify(result.comparison?.changes ?? [], null, 2)}
          </pre>

          <h3>Consumer Impact Details</h3>
          <pre className="impact-json">
            {JSON.stringify(result.impact, null, 2)}
          </pre>
        </section>
      )}
    </div>
  )
}

export default ConsumerImpactPage
