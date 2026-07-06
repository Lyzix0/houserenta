CREATE SCHEMA app;

CREATE TABLE app.users (
    id            TEXT PRIMARY KEY,
    name          TEXT NOT NULL,
    email         TEXT UNIQUE NOT NULL
                  CHECK (email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$'),
    password_hash TEXT NOT NULL,
    role          TEXT NOT NULL
                  CHECK (role IN ('landlord', 'tenant', 'admin')),
    document      TEXT NOT NULL,
    phone         TEXT NOT NULL,
    payment_card  TEXT
);