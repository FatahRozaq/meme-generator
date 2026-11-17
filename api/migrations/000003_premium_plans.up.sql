-- Paket premium berbasis durasi: 1, 3, 30, 60, 90 hari (meme unlimited).
INSERT INTO plans (code, name, price_idr, duration_days, max_memes) VALUES
    ('premium_1d',  'Premium 1 Hari',  5000,  1,  0),
    ('premium_3d',  'Premium 3 Hari',  12000, 3,  0),
    ('premium_30d', 'Premium 30 Hari', 25000, 30, 0),
    ('premium_60d', 'Premium 60 Hari', 45000, 60, 0),
    ('premium_90d', 'Premium 90 Hari', 60000, 90, 0)
ON CONFLICT (code) DO UPDATE
SET name = EXCLUDED.name,
    price_idr = EXCLUDED.price_idr,
    duration_days = EXCLUDED.duration_days,
    max_memes = EXCLUDED.max_memes;

-- Hapus paket lama (aman bila belum ada subscription yang merujuknya).
DELETE FROM plans WHERE code IN ('basic', 'pro', 'unlimited');
