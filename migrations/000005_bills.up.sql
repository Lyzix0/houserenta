CREATE TABLE app.bills (
    id              TEXT PRIMARY KEY,
    property_id     TEXT NOT NULL REFERENCES app.properties(id) ON DELETE CASCADE,
    date            TEXT NOT NULL,
    due_date        TEXT NOT NULL,
    status          TEXT NOT NULL CHECK (status IN ('paid', 'unpaid')),
    total           REAL NOT NULL
);
