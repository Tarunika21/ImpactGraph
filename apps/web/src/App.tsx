
import { useEffect, useState } from 'react'
import './App.css'

type Report = {
  id: number
  serviceName: string
  oldVersion: string
  newVersion: string
  createdAt: string
  result: {
    comparison?: {
      breakingCount?: number
      totalChanges?: number
    }
    impact?: {
      impactedCount?: number
      totalConsumers?: number
    }
  }
}

type Section = 'overview' | 'comparison' | 'impact' | 'reports'

const navigation: { id: Section; label: string; icon: string }[] = [
  { id: 'overview', label: 'Overview', icon: '◫' },
  { id: 'comparison', label: 'API Comparison', icon: '⇄' },
  { id: 'impact', label: 'Consumer Impact', icon: '◇' },
  { id: 'reports', label: 'Report History', icon: '▤' },
]

function App() {
  const [section, setSection] = useState<Section>('overview')
  const [reports, setReports] = useState<Report[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [apiOnline, setApiOnline] = useState(false)

  useEffect(() => {
    async function loadData() {
      try {
        const [healthResponse, reportsResponse] = await Promise.all([
          fetch('/health'),
          fetch('/api/v1/reports'),
        ])

        setApiOnline(healthResponse.ok)

        if (!reportsResponse.ok) {
          throw new Error(`Report API returned ${reportsResponse.status}`)
        }

        const data: Report[] = await reportsResponse.json()
        setReports(Array.isArray(data) ? data : [])
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Unable to load dashboard')
      } finally {
        setLoading(false)
      }
    }

    void loadData()
  }, [])

  const totalBreaking = reports.reduce(
    (sum, report) => sum + (report.result?.comparison?.breakingCount ?? 0),
    0,
  )

  const totalImpacted = reports.reduce(
    (sum, report) => sum + (report.result?.impact?.impactedCount ?? 0),
    0,
  )

  const uniqueServices = new Set(reports.map((report) => report.serviceName)).size

  return (
    <div className="app-shell">
      <aside className="sidebar">
        <div className="brand">
          <span className="brand-mark">IG</span>
          <div>
            <strong>ImpactGraph</strong>
            <small>API Intelligence</small>
          </div>
        </div>

        <div className="sidebar-label">WORKSPACE</div>

        <nav className="navigation">
          {navigation.map((item) => (
            <button
              key={item.id}
              className={`nav-item ${section === item.id ? 'active' : ''}`}
              onClick={() => setSection(item.id)}
            >
              <span className="nav-icon">{item.icon}</span>
              {item.label}
            </button>
          ))}
        </nav>

        <div className="sidebar-footer">
          <span className={`status-dot ${apiOnline ? 'online' : ''}`} />
          <span>Go API {apiOnline ? 'Connected' : 'Disconnected'}</span>
        </div>
      </aside>

      <main className="main-content">
        <header className="topbar">
          <div>
            <span className="breadcrumb">Workspace / {navigation.find((item) => item.id === section)?.label}</span>
          </div>
          <span className="environment">Development</span>
        </header>

        <div className="page-content">
          {section === 'overview' && (
            <>
              <div className="page-heading">
                <div>
                  <h1>Overview</h1>
                  <p>Monitor API contract changes and their downstream impact.</p>
                </div>
                <button className="primary-button" onClick={() => setSection('comparison')}>
                  + New Analysis
                </button>
              </div>

              {error && <div className="error-banner">{error}</div>}

              <div className="stats-grid">
                <StatCard title="Saved Reports" value={reports.length} symbol="▤" />
                <StatCard title="Services Analyzed" value={uniqueServices} symbol="◫" />
                <StatCard title="Breaking Changes Found" value={totalBreaking} symbol="⚡" />
                <StatCard title="Impacted Consumers" value={totalImpacted} symbol="◇" />
              </div>

              <div className="panel">
                <div className="panel-header">
                  <div>
                    <h2>Recent Analyses</h2>
                    <p>Latest saved API impact reports</p>
                  </div>
                  <button className="text-button" onClick={() => setSection('reports')}>
                    View all →
                  </button>
                </div>
                <ReportTable reports={reports.slice(0, 5)} loading={loading} />
              </div>
            </>
          )}

          {section === 'reports' && (
            <>
              <div className="page-heading">
                <div>
                  <h1>Report History</h1>
                  <p>Explore previously generated impact analysis reports.</p>
                </div>
              </div>
              {error && <div className="error-banner">{error}</div>}
              <div className="panel">
                <div className="panel-header">
                  <h2>All Reports</h2>
                </div>
                <ReportTable reports={reports} loading={loading} />
              </div>
            </>
          )}

          {section === 'comparison' && (
            <Placeholder
              title="API Comparison"
              description="Compare two API contract versions to detect breaking changes."
            />
          )}

          {section === 'impact' && (
            <Placeholder
              title="Consumer Impact"
              description="Identify which applications are affected by API changes."
            />
          )}
        </div>
      </main>
    </div>
  )
}

function StatCard({ title, value, symbol }: { title: string; value: number; symbol: string }) {
  return (
    <div className="stat-card">
      <div className="stat-top">
        <span>{title}</span>
        <span className="stat-symbol">{symbol}</span>
      </div>
      <strong>{value}</strong>
      <small>Based on saved analyses</small>
    </div>
  )
}

function ReportTable({ reports, loading }: { reports: Report[]; loading: boolean }) {
  if (loading) return <div className="empty-state">Loading reports...</div>
  if (reports.length === 0) return <div className="empty-state">No reports available yet.</div>

  return (
    <div className="table-scroll">
      <table>
        <thead>
          <tr>
            <th>Service</th>
            <th>Version Change</th>
            <th>Breaking Changes</th>
            <th>Impacted Consumers</th>
            <th>Created</th>
          </tr>
        </thead>
        <tbody>
          {reports.map((report) => (
            <tr key={report.id}>
              <td className="service-name">{report.serviceName}</td>
              <td>{report.oldVersion} → {report.newVersion}</td>
              <td>
                <span className="metric-pill danger">
                  {report.result?.comparison?.breakingCount ?? 0}
                </span>
              </td>
              <td>
                <span className="metric-pill">
                  {report.result?.impact?.impactedCount ?? 0}
                </span>
              </td>
              <td>{new Date(report.createdAt).toLocaleDateString()}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}

function Placeholder({ title, description }: { title: string; description: string }) {
  return (
    <div className="placeholder">
      <div className="placeholder-icon">◇</div>
      <h1>{title}</h1>
      <p>{description}</p>
      <span>We'll implement this workflow in the next milestone.</span>
    </div>
  )
}

export default App
