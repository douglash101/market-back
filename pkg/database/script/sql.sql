CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE states (
    id VARCHAR(2) NOT NULL PRIMARY KEY, -- SP, RJ, MG, etc.
    name VARCHAR(50) NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE cities (
    id UUID NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
    state_id VARCHAR(2) NOT NULL,
    name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (state_id) REFERENCES states(id) ON DELETE CASCADE,
    UNIQUE(state_id, name)
);

-- Insert some sample states
INSERT INTO states (id, name) VALUES 
('PR', 'Paraná'),
('SP', 'São Paulo'),
('RJ', 'Rio de Janeiro'),
('MG', 'Minas Gerais'),
('RS', 'Rio Grande do Sul'),
('SC', 'Santa Catarina'),
('BA', 'Bahia'),
('PE', 'Pernambuco'),
('CE', 'Ceará'),
('GO', 'Goiás'),
('DF', 'Distrito Federal'),
('ES', 'Espírito Santo'),
('MT', 'Mato Grosso'),
('MS', 'Mato Grosso do Sul'),
('AM', 'Amazonas'),
('PA', 'Pará'),
('MA', 'Maranhão'),
('PI', 'Piauí'),
('AL', 'Alagoas'),
('SE', 'Sergipe'),
('RN', 'Rio Grande do Norte'),
('PB', 'Paraíba'),
('TO', 'Tocantins'),
('AC', 'Acre'),
('AP', 'Amapá'),
('RO', 'Rondônia'),
('RR', 'Roraima');

-- Insert some sample cities
INSERT INTO cities (state_id, name) VALUES ('PR', 'Maringá');

CREATE TABLE markets (
    id UUID NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(120) NOT NULL,
    description VARCHAR(255) NOT NULL,
    status VARCHAR(20) DEFAULT 'active', -- active, inactive, deleted
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE market_locations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    market_id UUID NOT NULL,
    state_id VARCHAR(2) NOT NULL,
    city_id UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (market_id) REFERENCES markets(id) ON DELETE CASCADE,
    FOREIGN KEY (state_id) REFERENCES states(id) ON DELETE CASCADE,
    FOREIGN KEY (city_id) REFERENCES cities(id) ON DELETE CASCADE,
    UNIQUE(market_id, state_id, city_id)
);

CREATE TABLE attachments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    url VARCHAR(255) NOT NULL,
    type VARCHAR(50), -- image, document, etc.
    description VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(80) UNIQUE NOT NULL,
    password VARCHAR(120) NOT NULL,
    name VARCHAR(80) NOT NULL,
    status VARCHAR(20) DEFAULT 'active', -- active, inactive, deleted
    email_verified BOOLEAN DEFAULT FALSE,
    email_verification_token VARCHAR(255),
    password_reset_token VARCHAR(255),
    password_reset_expires TIMESTAMP WITH TIME ZONE,
    last_login TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_users_email ON users(email);

CREATE TABLE categories (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    status VARCHAR(20) DEFAULT 'active', -- active, inactive, deleted
    category_id UUID REFERENCES categories(id) ON DELETE SET NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE products (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    image_url VARCHAR(180),
    name VARCHAR(100) NOT NULL,
    category_id UUID REFERENCES categories(id) ON DELETE SET NULL,
    unit VARCHAR(20), -- kg, unit, liter, etc.
    status VARCHAR(20) DEFAULT 'active', -- active, inactive, deleted
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE product_markets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    provider_id VARCHAR(100),
    product_id UUID NOT NULL,
    market_id UUID NOT NULL,
    price NUMERIC(10,2) NOT NULL,
    promotional_price NUMERIC(10,2),
    status VARCHAR(20) DEFAULT 'active', -- active, inactive, deleted
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE,
    FOREIGN KEY (market_id) REFERENCES markets(id) ON DELETE CASCADE,
    UNIQUE(provider_id, product_id, market_id)
);
CREATE INDEX idx_product_markets_product_id ON product_markets(product_id);
CREATE INDEX idx_product_markets_market_id ON product_markets(market_id);
CREATE INDEX idx_product_markets_provider_id ON product_markets(provider_id);

INSERT INTO public.markets
(id, "name", description, status, created_at, updated_at)
VALUES('65dcfe06-0381-47fa-8fee-64aa45fa30b4'::uuid, 'Muffato', 'Muffato', 'active', '2025-11-02 22:18:54.834', '2025-11-02 22:18:54.834');
INSERT INTO public.markets
(id, "name", description, status, created_at, updated_at)
VALUES('f7c82abd-bd7b-4bf6-a0fc-811e2d589b89'::uuid, 'Amigão', 'Amigão', 'active', '2025-11-02 22:18:54.834', '2025-11-02 22:18:54.834');