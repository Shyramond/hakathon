-- Каталог бенефитов
CREATE TABLE IF NOT EXISTS benefits (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name          VARCHAR(255) NOT NULL,
    description   TEXT NOT NULL DEFAULT '',
    type          VARCHAR(50) NOT NULL CHECK (type IN ('promotion', 'skin', 'mascot_skin')),
    price_tokens  INTEGER NOT NULL CHECK (price_tokens >= 0),
    image_url     VARCHAR(512) NOT NULL DEFAULT '',
    is_active     BOOLEAN NOT NULL DEFAULT true,
    metadata      JSONB NOT NULL DEFAULT '{}',
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Кошелёк пользователя
CREATE TABLE IF NOT EXISTS user_wallets (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    xsolla_user_id  VARCHAR(255) UNIQUE NOT NULL,
    username        VARCHAR(255) NOT NULL DEFAULT '',
    email           VARCHAR(255) NOT NULL DEFAULT '',
    token_balance   INTEGER NOT NULL DEFAULT 0 CHECK (token_balance >= 0),
    last_login_at   TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Инвентарь пользователя (купленные бенефиты)
CREATE TABLE IF NOT EXISTS user_inventory (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id       UUID NOT NULL REFERENCES user_wallets(id) ON DELETE CASCADE,
    benefit_id    UUID NOT NULL REFERENCES benefits(id) ON DELETE RESTRICT,
    purchased_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    is_equipped   BOOLEAN NOT NULL DEFAULT false,
    UNIQUE(user_id, benefit_id)
);

-- Лог транзакций токенов
CREATE TABLE IF NOT EXISTS token_transactions (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id       UUID NOT NULL REFERENCES user_wallets(id) ON DELETE CASCADE,
    amount        INTEGER NOT NULL,
    type          VARCHAR(50) NOT NULL CHECK (type IN ('login_reward', 'purchase', 'refund', 'admin')),
    reference_id  UUID,
    description   VARCHAR(500) NOT NULL DEFAULT '',
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Таблица идемпотентности запросов
CREATE TABLE IF NOT EXISTS idempotency_keys (
    key           VARCHAR(255) PRIMARY KEY,
    user_id       UUID NOT NULL,
    response_code INTEGER NOT NULL,
    response_body JSONB NOT NULL DEFAULT '{}',
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at    TIMESTAMPTZ NOT NULL
);

-- Индексы
CREATE INDEX IF NOT EXISTS idx_user_inventory_user ON user_inventory(user_id);
CREATE INDEX IF NOT EXISTS idx_user_inventory_benefit ON user_inventory(benefit_id);
CREATE INDEX IF NOT EXISTS idx_token_transactions_user ON token_transactions(user_id);
CREATE INDEX IF NOT EXISTS idx_token_transactions_created ON token_transactions(created_at);
CREATE INDEX IF NOT EXISTS idx_benefits_active ON benefits(is_active) WHERE is_active = true;
CREATE INDEX IF NOT EXISTS idx_idempotency_expires ON idempotency_keys(expires_at);

-- Seed: базовые бенефиты
INSERT INTO benefits (name, description, type, price_tokens, image_url) VALUES
    ('Весенняя акция', 'Скидка 10% на следующую покупку', 'promotion', 30, '/images/spring_promo.png'),
    ('Скин "Огненный кот"', 'Огненный скин для десктопного маскота', 'mascot_skin', 100, '/images/fire_cat.png'),
    ('Скин "Ледяной пёс"', 'Ледяной скин для десктопного маскота', 'mascot_skin', 80, '/images/ice_dog.png'),
    ('Скин "Космонавт"', 'Космический скин для маскота', 'mascot_skin', 150, '/images/astronaut.png'),
    ('Летняя распродажа', 'Двойные токены на выходных', 'promotion', 50, '/images/summer_sale.png')
ON CONFLICT DO NOTHING;