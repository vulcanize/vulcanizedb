CREATE TABLE token_burns (
  id                    SERIAL,
  vulcanize_log_id      INTEGER NOT NULL UNIQUE,
  token_name            CHARACTER VARYING(66) NOT NULL,
  token_address         CHARACTER VARYING(66) NOT NULL,
  burner                CHARACTER VARYING(66) NOT NULL,
  tokens                DECIMAL NOT NULL,
  block                 INTEGER NOT NULL,
  tx                    CHARACTER VARYING(66) NOT NULL,
  CONSTRAINT log_index_fk FOREIGN KEY (vulcanize_log_id)
  REFERENCES logs (id)
  ON DELETE CASCADE
)