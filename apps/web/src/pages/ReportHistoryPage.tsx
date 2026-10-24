
import { useEffect, useState } from 'react'
import './ReportHistoryPage.css'

type Change = {
  type?: string
  endpoint?: string
  field?: string
  breaking?: boolean
}

type ImpactedConsumer = {
  consumerId?: string
  consumerName?: string
  endpoint?: string
  reasons?: string[]
}

type ReportResult = {
  comparison?: {
    totalChanges?: number
    breakingCount?: number
    changes?: Change[]
  }
  impact?: {
    totalConsumers?: number
    impactedCount?: number
    impactedConsumers?: ImpactedConsumer[]
  }
}

type Report = {
  id: number
  serviceName: string
  oldVersion: string
  newVersion: string
  createdAt: string
  result: ReportResult
}

function ReportHistoryPage() {
  const [reports, setReports] = useState<Report[]>([])
  const [search, setSearch] = useState('')
  const [selectedId, setSelectedId] = useState<number | null>(null)
  const [selectedReport, setSelectedReport] = useState<Report | null>(null)
  const [loading, setLoading] = useState(true)
  const [detailsLoading, setDetailsLoading] = useState(false)
  const [error, setError] = useState('')
  const [detailsError, setDetailsError] = useState('')

  useEffect(() => {
    const controller = new AbortController()

    async function loadReports() {
      try {
        const response = await fetch('/api/v1/reports', {
          signal: controller.signal,
        })
        if (!response.ok) {
          throw new Error(`Failed to load reports (${response.status})`)
        }

        const data: Report[] = await response.json()
        if (!controller.signal.aborted) {
          setReports(Array.isArray(data) ? data : [])
        }
      } catch (err) {
        if (!controller.signal.aborted) {
          setError(
            err instanceof Error ? err.message : 'Unable to load reports.',
          )
        }
      } finally {
        if (!controller.signal.aborted) setLoading(false)
      }
    }

    void loadReports()
    return () => controller.abort()
  }, [])

  async function openReport(id: number) {
    if (selectedId === id) {
      setSelectedId(null)
      setSelectedReport(null)
      setDetailsError('')
      return
    }

    setSelectedId(id)
    setSelectedReport(null)
    setDetailsError('')
    setDetailsLoading(true)

    try {
      const response = await fetch(`/api/v1/reports/${id}`)
      if (!response.ok) {
        throw new Error(`Failed to load report (${response.status})`)
      }

      const data: Report = await response.json()
      setSelectedReport(data)
    } catch (err) {
      setDetailsError(
        err instanceof Error ? err.message : 'Unable to load report details.',
      )
    } finally {
      setDetailsLoading(false)
    }
  }

  const filteredReports = reports
    .filter((report) =>
      report.serviceName.toLowerCase().includes(search.trim().toLowerCase()),
    )
    .sort(
      (a, b) =>
        new Date(b.createdAt).getTime() -
        new Date(a.createdAt).getTime() || b.id - a.id,
    )

  return (
    <div className="history-page">
      <div className="page-heading">
        <div>
          <h1>Report History</h1>
          <p>Search and inspect previously saved API impact analyses.</p>
        </div>
      </div>

      <div className="history-toolbar">
        <label className="history-search">
          <span>Search by service</span>
          <input
            type="search"
            value={search}
            onChange={(event) => setSearch(event.target.value)}
            placeholder="Search service name..."
          />
        </label>
        <span>{filteredReports.length} reports</span>
      </div>

      {error && <div className="error-banner" role="alert">{error}</div>}

      {loading ? (
        <div className="history-empty">Loading reports...</div>
      ) : filteredReports.length === 0 ? (
        <div className="history-empty">
          {search ? 'No reports match your search.' : 'No saved reports yet.'}
        </div>
      ) : (
        <div className="history-list">
          {filteredReports.map((report) => (
            <article className="history-card" key={report.id}>
              <div className="history-card-top">
                <div>
                  <h2>{report.serviceName}</h2>
                  <p>
                    Report #{report.id} · {report.oldVersion} → {report.newVersion}
                  </p>
                </div>
                <button
                  type="button"
                  className="history-open-button"
                  onClick={() => void openReport(report.id)}
                  aria-expanded={selectedId === report.id}
                >
                  {selectedId === report.id ? 'Close Details' : 'View Details'}
                </button>
              </div>

              <div className="history-card-meta">
                <span>
                  Breaking changes: {report.result?.comparison?.breakingCount ?? 0}
                </span>
                <span>
                  Impacted consumers: {report.result?.impact?.impactedCount ?? 0}
                </span>
                <span>{new Date(report.createdAt).toLocaleString()}</span>
              </div>

              {selectedId === report.id && (
                <div className="history-details">
                  {detailsLoading && <p>Loading report details...</p>}
                  {detailsError && (
                    <div className="error-banner" role="alert">
                      {detailsError}
                    </div>
                  )}
                  {selectedReport && (
                    <ReportDetails report={selectedReport} />
                  )}
                </div>
              )}
            </article>
          ))}
        </div>
      )}
    </div>
  )
}

function ReportDetails({ report }: { report: Report }) {
  const changes = report.result?.comparison?.changes ?? []
  const consumers = report.result?.impact?.impactedConsumers ?? []

  return (
    <>
      <div className="history-details-summary">
        <div>
          <small>Total Changes</small>
          <strong>{report.result?.comparison?.totalChanges ?? 0}</strong>
        </div>
        <div>
          <small>Breaking Changes</small>
          <strong>{report.result?.comparison?.breakingCount ?? 0}</strong>
        </div>
        <div>
          <small>Impacted Consumers</small>
          <strong>{report.result?.impact?.impactedCount ?? 0}</strong>
        </div>
      </div>

      <h3>Detected API Changes</h3>
      {changes.length === 0 ? (
        <p>No API changes recorded.</p>
      ) : (
        <div className="history-detail-list">
          {changes.map((change, index) => (
            <div className="history-detail-item" key={index}>
              <div>
                <strong>{change.type ?? 'API Change'}</strong>
                <p>{change.endpoint ?? 'Unknown endpoint'}</p>
                {change.field && <p>Field: {change.field}</p>}
              </div>
              <span className={change.breaking ? 'history-breaking' : 'history-safe'}>
                {change.breaking ? 'Breaking' : 'Non-breaking'}
              </span>
            </div>
          ))}
        </div>
      )}

      <h3>Impacted Consumers</h3>
      {consumers.length === 0 ? (
        <p>No impacted consumers recorded.</p>
      ) : (
        <div className="history-detail-list">
          {consumers.map((consumer, index) => (
            <div className="history-detail-item" key={index}>
              <div>
                <strong>{consumer.consumerName ?? consumer.consumerId ?? 'Consumer'}</strong>
                <p>{consumer.endpoint ?? 'Unknown endpoint'}</p>
                {consumer.reasons?.map((reason, reasonIndex) => (
                  <p key={reasonIndex}>Reason: {reason}</p>
                ))}
              </div>
            </div>
          ))}
        </div>
      )}
    </>
  )
}

export default ReportHistoryPage
