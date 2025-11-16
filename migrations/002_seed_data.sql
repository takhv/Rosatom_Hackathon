-- Тестовые города
INSERT OR IGNORE INTO cities (name, region) VALUES
('Ангарск', 'Иркутская область'),
('Иркутск', 'Иркутская область'),
('Братск', 'Иркутская область');

-- Тестовые НКО
INSERT OR IGNORE INTO ngos (name, category, description, volunteer_description, phone, address, city_id) VALUES
('Фонд Помощи Детям', 'социальная', 'Помощь детям из малообеспеченных семей', 'Нужны волонтеры для раздачи продуктов', '+7 123 456-78-90', 'ул. Ленина, 1', 1),
('ЭкоЗащита', 'экология', 'Защита окружающей среды', 'Требуются волонтеры для уборки парков', '+7 123 456-78-91', 'ул. Мира, 15', 2);

INSERT OR IGNORE INTO users (email, password_hash, full_name, role) VALUES
(
    'admin@example.com',
    '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', -- password: admin123
    'Администратор Системы',
    'admin'
);

-- Обычный пользователь
INSERT OR IGNORE INTO users (email, password_hash, full_name, role) VALUES
(
    'user@example.com', 
    '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', -- password: user123
    'Иван Петров',
    'user'
);

-- Пользователь привязанный к НКО "Фонд Помощи Детям"
INSERT OR IGNORE INTO users (email, password_hash, full_name, ngo_id, role) VALUES
(
    'ngo_user@example.com',
    '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', -- password: nko_user123
    'Мария Сидорова',
    1, -- ID НКО "Фонд Помощи Детям"
    'user'
);

-- Еще один пользователь привязанный к НКО "ЭкоЗащита"
INSERT OR IGNORE INTO users (email, password_hash, full_name, ngo_id, role) VALUES
(
    'eco_user@example.com',
    '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', -- password: nko_user123
    'Алексей Зеленов',
    2, -- ID НКО "ЭкоЗащита" 
    'user'
);