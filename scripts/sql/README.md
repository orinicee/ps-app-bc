# Migraciones SQL

Este directorio contiene los scripts SQL para la gestión de la base de datos.

## Estructura de Archivos

- `XXX_nombre.sql`: Scripts de migración forward (crear/modificar)
- `XXX_nombre_rollback.sql`: Scripts de rollback (revertir cambios)

## Convenciones de Nombrado

- Los archivos se numeran con un prefijo de 3 dígitos (001, 002, etc.)
- El nombre debe ser descriptivo de los cambios que realiza
- Cada script forward debe tener su correspondiente script de rollback

## Ejecutar Migraciones

Para ejecutar las migraciones, usa el cliente psql:

```bash
# Conectar a la base de datos
psql -h localhost -U usuario -d nombre_db

# Ejecutar un script
\i path/to/script.sql
```

## Orden de Ejecución

1. 001_create_contents_table.sql

## Rollback

Para revertir los cambios, ejecuta los scripts rollback en orden inverso:

1. 001_create_contents_table_rollback.sql 