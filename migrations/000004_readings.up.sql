CREATE TABLE app.readings (
    id              TEXT PRIMARY KEY,
    property_id     TEXT NOT NULL REFERENCES app.properties(id) ON DELETE CASCADE,
    date            TEXT NOT NULL,
    gvs             REAL NOT NULL,
    hvs             REAL NOT NULL,
    el1             REAL NOT NULL,
    el2             REAL,
    is_accounted    INTEGER NOT NULL DEFAULT 0
);
