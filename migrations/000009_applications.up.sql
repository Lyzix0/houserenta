CREATE TABLE app.applications (
    id              TEXT PRIMARY KEY,
    property_id     TEXT NOT NULL REFERENCES app.properties(id) ON DELETE CASCADE,
    tenant_user_id  TEXT NOT NULL REFERENCES app.users(id) ON DELETE CASCADE,
    date            TEXT NOT NULL,
    UNIQUE (property_id, tenant_user_id)
);
