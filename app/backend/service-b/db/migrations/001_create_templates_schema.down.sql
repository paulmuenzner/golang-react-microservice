-- Beispielinhalt für service-b/db/migrations/001_create_templates_schema.down.sql
DROP TABLE IF EXISTS templates.userstemps CASCADE;

DROP SCHEMA IF EXISTS templates CASCADE;