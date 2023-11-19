CREATE TABLE users
(
    id         SERIAL PRIMARY KEY,
    address    VARCHAR(42) UNIQUE,
    created_at TIMESTAMP(0) NOT NULL DEFAULT now(),
    online_at  TIMESTAMP(0) NOT NULL DEFAULT now()
);