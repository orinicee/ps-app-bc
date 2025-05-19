-- Crear tabla de tipos de contenido
CREATE TABLE content_types (
    type_name VARCHAR(10) PRIMARY KEY,
    description VARCHAR(100) NOT NULL
);

-- Insertar tipos válidos
INSERT INTO content_types (type_name, description) VALUES
    ('vlog', 'Contenido tipo vlog'),
    ('video', 'Contenido tipo video'),
    ('audio', 'Contenido tipo audio'),
    ('live', 'Contenido tipo live');

-- Crear tabla principal de contenidos
CREATE TABLE contents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    url VARCHAR(2048) NOT NULL,
    type VARCHAR(10) NOT NULL,
    is_free BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by UUID NOT NULL,
    FOREIGN KEY (type) REFERENCES content_types(type_name),
    FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE CASCADE
);

-- Crear índices para optimizar búsquedas
CREATE INDEX idx_contents_type ON contents(type);
CREATE INDEX idx_contents_is_free ON contents(is_free);
CREATE INDEX idx_contents_created_by ON contents(created_by);