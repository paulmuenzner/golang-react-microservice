# 1. Teste complete Workflow
make db-migrate
make db-status
make db-rollback
make db-migrate

# 2. Teste with dev running
make dev
make db-migrate  # Nutzt running container
make db-tables

# 3. Backup-Test
make db-backup
make db-restore FILE=backups/...sql

# 4. Production preparation
make prod-db-check-breaking  # Prüft Breaking Changes