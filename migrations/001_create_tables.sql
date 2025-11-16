-- === ТАБЛИЦА ГОРОДОВ ===
CREATE TABLE IF NOT EXISTS cities (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    name        VARCHAR(255) NOT NULL,          -- "Ангарск"
    region      VARCHAR(255)                    -- "Иркутская область"
);

-- === ТАБЛИЦА НКО ===
CREATE TABLE IF NOT EXISTS ngos (
    id                      INTEGER PRIMARY KEY AUTOINCREMENT,
    name                    VARCHAR(255) NOT NULL,      -- name (форма НКО)
    category                VARCHAR(100) NOT NULL,      -- category
    description             TEXT,                       -- description
    volunteer_description   TEXT,                       -- volunteerDescription
    phone                   VARCHAR(50),                -- phone
    address                 TEXT,                       -- address
    logo_url                TEXT,                       -- опционально: фото/логотип
    website_url             TEXT,                       -- ссылка на сайт
    social_links            TEXT,                       -- соцсети (можно хранить через запятую или JSON)
    city_id                 INTEGER NOT NULL REFERENCES cities(id),
    status                  VARCHAR(20) NOT NULL DEFAULT 'pending',  -- для модерации: pending/approved/rejected
    created_at              DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at              DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- === ТАБЛИЦА ПОЛЬЗОВАТЕЛЕЙ ===
CREATE TABLE IF NOT EXISTS users (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    email           VARCHAR(255) UNIQUE NOT NULL,   -- login/email (login-form, register-form)
    password_hash   TEXT NOT NULL,                  -- ХРАНИМ ТОЛЬКО ХЭШ, НЕ ПАРОЛЬ!
    full_name       VARCHAR(255),                   -- fullname из формы
    ngo_id          INTEGER UNIQUE REFERENCES ngos(id),
    role            VARCHAR(20) NOT NULL DEFAULT 'user',  -- user / admin
    created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
