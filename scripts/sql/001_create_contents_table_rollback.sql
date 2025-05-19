-- Eliminar índices
DROP INDEX IF EXISTS idx_contents_type;
DROP INDEX IF EXISTS idx_contents_is_free;
DROP INDEX IF EXISTS idx_contents_created_by;

-- Eliminar tablas
DROP TABLE IF EXISTS contents;
DROP TABLE IF EXISTS content_types; 