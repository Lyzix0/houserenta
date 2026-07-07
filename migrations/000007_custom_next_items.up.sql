CREATE TABLE app.custom_next_items (
    id              TEXT PRIMARY KEY,
    property_id     TEXT NOT NULL REFERENCES app.properties(id) ON DELETE CASCADE,
    description     TEXT NOT NULL,
    amount          REAL NOT NULL
);
