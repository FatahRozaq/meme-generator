DELETE FROM plans WHERE code IN ('premium_1d', 'premium_3d', 'premium_30d', 'premium_60d', 'premium_90d');

INSERT INTO plans (code, name, price_idr, duration_days, max_memes) VALUES
    ('basic', 'Basic', 25000, 30, 20),
    ('pro', 'Pro', 50000, 30, 100),
    ('unlimited', 'Unlimited', 100000, 30, 0);
