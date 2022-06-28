CREATE TABLE IF NOT EXISTS block (
    num NUMERIC,
    block_hash CHAR(66) PRIMARY KEY,
    block_time TIMESTAMPTZ NOT NULL,
    parent_hash CHAR(66),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    stable BOOLEAN DEFAULT TRUE,
    done BOOLEAN DEFAULT FALSE
);

CREATE INDEX num_idx ON block(num);

CREATE TABLE IF NOT EXISTS tx (
  block_hash CHAR(66) NOT NULL,
  tx_hash CHAR(66) PRIMARY KEY,
  sender CHAR(42) NOT NULL,
  receiver CHAR(42) NOT NULL,
  nonce INT,
  tx_data TEXT,
  amount NUMERIC,
  CONSTRAINT fk_block_hash
    FOREIGN KEY(block_hash)
      REFERENCES block(block_hash)
);

CREATE TABLE IF NOT EXISTS log (
  tx_hash CHAR(66),
  log_id int,
  data TEXT,
  CONSTRAINT fk_tx_hash
    FOREIGN KEY(tx_hash)
      REFERENCES tx(tx_hash),
  PRIMARY KEY (tx_hash, log_id)
);
