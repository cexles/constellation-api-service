CREATE TABLE chains
(
    id       SERIAL PRIMARY KEY,
    chain_id BIGINT       NOT NULL,
    name     VARCHAR(255) NOT NULL
);