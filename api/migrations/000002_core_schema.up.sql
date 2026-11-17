-- Skema inti meme-portal: users, plans, subscriptions, payments, memes.

CREATE TABLE users (
    id            uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    email         citext NOT NULL UNIQUE,
    name          text NOT NULL DEFAULT '',
    password_hash text NOT NULL,
    role          text NOT NULL DEFAULT 'user' CHECK (role IN ('user', 'admin')),
    created_at    timestamptz NOT NULL DEFAULT now(),
    updated_at    timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE plans (
    id            uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    code          text NOT NULL UNIQUE,
    name          text NOT NULL,
    price_idr     bigint NOT NULL CHECK (price_idr >= 0),
    duration_days integer NOT NULL CHECK (duration_days > 0),
    max_memes     integer NOT NULL DEFAULT 0 CHECK (max_memes >= 0), -- 0 = unlimited
    active        boolean NOT NULL DEFAULT true,
    created_at    timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE subscriptions (
    id         uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    plan_id    uuid NOT NULL REFERENCES plans (id),
    status     text NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'active', 'expired', 'canceled')),
    starts_at  timestamptz,
    ends_at    timestamptz,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX idx_subscriptions_user_active ON subscriptions (user_id, status, ends_at DESC);

CREATE TABLE payments (
    id              uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    subscription_id uuid NOT NULL REFERENCES subscriptions (id) ON DELETE CASCADE,
    amount_idr      bigint NOT NULL,
    status          text NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'paid', 'failed')),
    reference       text NOT NULL UNIQUE,
    paid_at         timestamptz,
    created_at      timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE memes (
    id         uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    title      text NOT NULL DEFAULT '',
    base_url   text NOT NULL,
    result_url text,
    top_text   text NOT NULL DEFAULT '',
    bottom_text text NOT NULL DEFAULT '',
    font_size  integer NOT NULL DEFAULT 48,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX idx_memes_user_created ON memes (user_id, created_at DESC);

-- Paket awal (seed).
INSERT INTO plans (code, name, price_idr, duration_days, max_memes) VALUES
    ('basic', 'Basic', 25000, 30, 20),
    ('pro', 'Pro', 50000, 30, 100),
    ('unlimited', 'Unlimited', 100000, 30, 0);
