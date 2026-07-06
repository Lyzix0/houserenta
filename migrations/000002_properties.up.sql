CREATE TABLE app.properties (
    id              TEXT PRIMARY KEY,
    landlord_id     TEXT NOT NULL REFERENCES app.users(id) ON DELETE CASCADE,
    name            TEXT NOT NULL,
    coordinates     TEXT NOT NULL,
    country         TEXT,
    region          TEXT NOT NULL,
    city            TEXT NOT NULL,
    street          TEXT NOT NULL,
    house           TEXT NOT NULL,
    apartment       TEXT NOT NULL,
    gvs_tariff      REAL NOT NULL,
    hvs_tariff      REAL NOT NULL,
    el1_tariff      REAL NOT NULL,
    el2_tariff      REAL,
    balance         REAL NOT NULL DEFAULT 0
);