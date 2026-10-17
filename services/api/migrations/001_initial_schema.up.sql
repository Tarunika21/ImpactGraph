
CREATE TABLE api_contracts (
    id BIGSERIAL PRIMARY KEY,
    service_name VARCHAR(150) NOT NULL,
    version VARCHAR(100) NOT NULL,
    specification JSONB NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(service_name, version)
);

CREATE TABLE consumers (
    id VARCHAR(100) PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE consumer_dependencies (
    id BIGSERIAL PRIMARY KEY,
    consumer_id VARCHAR(100) NOT NULL
        REFERENCES consumers(id) ON DELETE CASCADE,
    service_name VARCHAR(150) NOT NULL,
    endpoint VARCHAR(300) NOT NULL,
    fields JSONB NOT NULL DEFAULT '[]'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(consumer_id, service_name, endpoint)
);

CREATE TABLE impact_reports (
    id BIGSERIAL PRIMARY KEY,
    service_name VARCHAR(150) NOT NULL,
    old_version VARCHAR(100) NOT NULL,
    new_version VARCHAR(100) NOT NULL,
    result JSONB NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_contract_service
ON api_contracts(service_name);

CREATE INDEX idx_dependencies_service
ON consumer_dependencies(service_name);

CREATE INDEX idx_impact_reports_service
ON impact_reports(service_name);
