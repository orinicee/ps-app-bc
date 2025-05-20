#!/bin/bash

# Crear la base de datos de test
PGPASSWORD=postgres123 psql -h localhost -U postgres -c "DROP DATABASE IF EXISTS ps_app_test;"
PGPASSWORD=postgres123 psql -h localhost -U postgres -c "CREATE DATABASE ps_app_test;"

# Aplicar las migraciones a la base de datos de test
# TODO: Agregar comandos de migración cuando estén disponibles 