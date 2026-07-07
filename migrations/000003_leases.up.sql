CREATE TABLE app.leases (
    id              TEXT PRIMARY KEY,
    property_id     TEXT UNIQUE NOT NULL REFERENCES app.properties(id) ON DELETE CASCADE,
    tenant_user_id  TEXT NOT NULL REFERENCES app.users(id) ON DELETE CASCADE,
    name            TEXT NOT NULL,
    document        TEXT NOT NULL,
    phone           TEXT NOT NULL,
    months_of_rent  INTEGER NOT NULL,
    price           REAL NOT NULL,
    payment_day     INTEGER NOT NULL CHECK (payment_day BETWEEN 1 AND 28),
    reading_day     INTEGER NOT NULL CHECK (reading_day BETWEEN 1 AND 28),
    start_date      TEXT NOT NULL,
    end_date        TEXT NOT NULL
);
