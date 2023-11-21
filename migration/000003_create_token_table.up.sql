CREATE TABLE tokens
(
    id         SERIAL PRIMARY KEY,
    chain_id   BIGINT       NOT NULL,
    address    VARCHAR(42)  NOT NULL,
    name       VARCHAR(255) NOT NULL,
    symbol     VARCHAR(10)  NOT NULL,
    decimals   INT          NOT NULL,
    created_at TIMESTAMP DEFAULT now(),
    CONSTRAINT fk_tokens_chains FOREIGN KEY (chain_id) REFERENCES chains(id)
);
