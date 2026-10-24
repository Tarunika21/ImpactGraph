
import { useState } from 'react'
import './ComparisonPage.css'

type ComparisonResult = {
  totalChanges?: number
  breakingCount?: number
  changes?: unknown[]
}

const exampleOldSpec = {
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

const exampleNewSpec = {
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

function ComparisonPage() {
  const [oldSpec, setOldSpec] = useState(
    JSON.stringify(exampleOldSpec, null, 2),
  )
  const [newSpec, setNewSpec] = useState(
    JSON.stringify(exampleNewSpec, null, 2),
  )
  const [result, setResult] = useState<ComparisonResult | null>(null)
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  async function compareSpecs() {
    setError('')
    setResult(null)

    let oldParsed: unknown
    let newParsed: unknown

    try {
      oldParsed = JSON.parse(oldSpec)
      newParsed = JSON.parse(newSpec)
    } catch {
      setError('Both API specifications must contain valid JSON.')
      return
    }

    setLoading(true)

    try {
      const response = await fetch('/api/v1/compare', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          oldSpec: oldParsed,
          newSpec: newParsed,
        }),
      })

      if (!response.ok) {
        const message = await response.text()
        throw new Error(
          `Comparison failed (${response.status}): ${message}`,
        )
      }

      const data: ComparisonResult = await response.json()
      setResult(data)
    } catch (err) {
      setError(
        err instanceof Error
          ? err.message
          : 'Unable to compare API specifications.',
      )
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="comparison-page">
      <div className="page-heading">
        <div>
          <h1>API Comparison</h1>
          <p>Compare API contract versions to identify breaking changes.</p>
        </div>
      </div>

      <div className="comparison-editors">
        <div className="comparison-editor">
          <div className="comparison-editor-header">
            <h2>Old API Specification</h2>
            <span>Previous version</span>
          </div>
          <textarea
            value={oldSpec}
            onChange={(event) => setOldSpec(event.target.value)}
            spellCheck={false}
            aria-label="Old API specification"
          />
        </div>

        <div className="comparison-editor">
          <div className="comparison-editor-header">
            <h2>New API Specification</h2>
            <span>Updated version</span>
          </div>
          <textarea
            value={newSpec}
            onChange={(event) => setNewSpec(event.target.value)}
            spellCheck={false}
            aria-label="New API specification"
          />
        </div>
      </div>

      <div className="comparison-actions">
        <button
          className="primary-button"
          onClick={() => void compareSpecs()}
          disabled={loading}
        >
          {loading ? 'Comparing...' : 'Compare Specifications'}
        </button>
        <span>Example specifications are preloaded for testing.</span>
      </div>

      {error && (
        <div className="error-banner" role="alert">
          {error}
        </div>
      )}

      {result && (
        <section className="comparison-results">
          <div className="comparison-results-heading">
            <h2>Comparison Results</h2>
            <span>Analysis completed</span>
          </div>

          <div className="comparison-summary">
            <div>
              <small>Total Changes</small>
              <strong>{result.totalChanges ?? 0}</strong>
            </div>
            <div>
              <small>Breaking Changes</small>
              <strong>{result.breakingCount ?? 0}</strong>
            </div>
            <div>
              <small>Compatibility</small>
              <strong className={
                (result.breakingCount ?? 0) > 0
                  ? 'comparison-danger'
                  : 'comparison-success'
              }>
                {(result.breakingCount ?? 0) > 0
                  ? 'Breaking'
                  : 'No breaking changes'}
              </strong>
            </div>
          </div>

          <h3>Detailed Changes</h3>
          <pre className="comparison-json">
            {JSON.stringify(result.changes ?? [], null, 2)}
          </pre>
        </section>
      )}
    </div>
  )
}

export default ComparisonPage
