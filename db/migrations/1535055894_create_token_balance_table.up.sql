CREATE TABLE token_balance (
  id                   SERIAL,
  block_id             INTEGER NOT NULL,
  balance              DECIMAL NOT NULL,
  token_address        CHARACTER VARYING(66) NOT NULL,
  token_holder_address CHARACTER VARYING(66) NOT NULL,
  CONSTRAINT blocks_fk FOREIGN KEY (block_id)
  REFERENCES blocks (id)
  ON DELETE CASCADE
)