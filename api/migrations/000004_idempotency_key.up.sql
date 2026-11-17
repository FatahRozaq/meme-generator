-- Idempotency key untuk checkout: satu maksud pembayaran = satu key.
ALTER TABLE payments ADD COLUMN idempotency_key TEXT;
CREATE UNIQUE INDEX payments_idempotency_key_idx
  ON payments (idempotency_key)
  WHERE idempotency_key IS NOT NULL;

INSERT INTO schema_migrations (version, dirty) VALUES (4, false)
ON CONFLICT (version) DO NOTHING;
