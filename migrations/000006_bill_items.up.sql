CREATE TABLE app.bill_items (
    id              TEXT PRIMARY KEY,
    bill_id         TEXT NOT NULL REFERENCES app.bills(id) ON DELETE CASCADE,
    description     TEXT NOT NULL,
    amount          REAL NOT NULL
);
