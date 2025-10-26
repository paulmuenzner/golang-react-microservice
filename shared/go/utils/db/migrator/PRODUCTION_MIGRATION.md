# Production Migrations Cheat Sheet

## 🚨 Golden Rules

```
1. NEVER break backward compatibility in Phase 1
2. ALWAYS backup before migrations
3. ADD first, REMOVE later (3-phase strategy)
4. TEST in staging before production
5. MONITOR for 24-48h between phases
```

---

## ✅ Safe Operations (Zero Downtime)

```sql
-- ✅ Add column (nullable)
ALTER TABLE users ADD COLUMN phone VARCHAR(20);

-- ✅ Add column with default
ALTER TABLE users ADD COLUMN status VARCHAR(20) DEFAULT 'active';

-- ✅ Create table
CREATE TABLE logs (...);

-- ✅ Create index (non-blocking)
CREATE INDEX CONCURRENTLY idx_users_email ON users(email);

-- ✅ Add constraint (nullable)
ALTER TABLE users ADD CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(id);

-- ✅ Insert data
INSERT INTO config (key, value) VALUES ('version', '2.0');
```

---

## ❌ Unsafe Operations (Causes Downtime)

```sql
-- ❌ Drop column (services crash!)
ALTER TABLE users DROP COLUMN email;

-- ❌ Drop table (services crash!)
DROP TABLE users;

-- ❌ Rename column (services can't find it!)
ALTER TABLE users RENAME COLUMN email TO email_address;

-- ❌ Change column type (data conversion issues!)
ALTER TABLE users ALTER COLUMN age TYPE INTEGER;

-- ❌ Add NOT NULL without default (old inserts fail!)
ALTER TABLE users ADD COLUMN email VARCHAR(255) NOT NULL;

-- ❌ Create blocking index (locks table!)
CREATE INDEX idx_users_email ON users(email);  -- Missing CONCURRENTLY
```

---

## 🔧 Quick Commands

### Development
```bash
make dev                    # Start development
make db-migrate            # Run migrations
make db-rollback           # Rollback last
make db-tables             # Show tables
make db-connect            # Connect to DB
```

### Production (Safe)
```bash
make prod-db-backup                # Backup first!
make prod-db-check-breaking        # Check for breaking changes
make prod-db-migrate-safe          # Migrate with safety checks
make prod-db-verify                # Verify schema
make prod-db-status                # Show migration history
make prod-db-health                # Health check
```

### Production (Maintenance Window)
```bash
make prod-db-maintenance           # Full maintenance
# Stops services → Backup → Migrate → Start services
```

### Rollback
```bash
make prod-db-rollback-safe         # Safe rollback
make db-restore FILE=backup.sql    # Restore from backup
```

---

## 📋 3-Phase Pattern (Rename Column)

### Phase 1: ADD (Week 1)
```sql
-- Migration: ADD new column
ALTER TABLE users ADD COLUMN email_address VARCHAR(255);
UPDATE users SET email_address = email;
CREATE INDEX CONCURRENTLY idx_users_email_address ON users(email_address);
```
```go
// Code: Use BOTH columns
type User struct {
    Email        string `db:"email"`          // Old
    EmailAddress string `db:"email_address"`  // New
}
func (u *User) GetEmail() string {
    if u.EmailAddress != "" { return u.EmailAddress }
    return u.Email  // Fallback
}
```
```bash
make prod-db-migrate-safe
make prod-deploy
```

### Phase 2: MONITOR (Week 2)
```bash
# Wait for all pods to update
# Monitor logs for errors
# Verify both columns have same data
```

### Phase 3: REMOVE (Week 3)
```sql
-- Migration: DROP old column
DROP INDEX IF EXISTS idx_users_email;
ALTER TABLE users DROP COLUMN email;
ALTER TABLE users ALTER COLUMN email_address SET NOT NULL;
```
```go
// Code: Use ONLY new column
type User struct {
    EmailAddress string `db:"email_address"`
}
```
```bash
make prod-db-migrate-safe
make prod-deploy
```

---

## 🎯 Decision Tree

```
Need to change schema?
│
├─ Is it additive? (ADD COLUMN, CREATE TABLE)
│  └─ ✅ Safe! Deploy immediately
│      make prod-db-migrate-safe
│
├─ Is it removable? (DROP COLUMN, DROP TABLE)
│  └─ ⚠️  3-Phase migration required
│      1. Add new (Week 1)
│      2. Monitor (Week 2)
│      3. Remove old (Week 3)
│
└─ Is it breaking? (RENAME, ALTER TYPE)
   ├─ Can I use 3-phase?
   │  └─ ✅ Yes → Use 3-phase
   │
   └─ ❌ No → Maintenance window
       make prod-db-maintenance
```

---

## 📊 Deployment Checklist

### Pre-Deployment
```
☐ Tested in staging
☐ Backup created
☐ Breaking changes checked (make prod-db-check-breaking)
☐ Team notified
☐ Rollback plan documented
☐ Monitoring dashboard ready
```

### During Deployment
```
☐ Backup verified
☐ Migration executed
☐ Schema verified (make prod-db-verify)
☐ Health checks pass
☐ No errors in logs
```

### Post-Deployment
```
☐ Monitor for 1 hour
☐ Error rate normal
☐ Performance normal
☐ Backup kept for 7 days
```

---

## 🆘 Emergency Procedures

### Migration Failed
```bash
# Check status
make prod-db-status

# Check logs
make db-logs

# Rollback if needed
make prod-db-rollback-safe
```

### Services Crashing After Migration
```bash
# 1. Rollback migration immediately
make prod-db-rollback-safe

# 2. Rollback deployment
kubectl rollout undo deployment/gateway

# 3. Verify
make prod-db-health
curl http://api.example.com/health
```

### Data Corruption
```bash
# 1. Stop all services (prevent further damage)
make prod-down

# 2. Restore from backup
make db-restore FILE=backups/prod_backup_TIMESTAMP.sql

# 3. Verify data
make db-connect
SELECT COUNT(*) FROM users;

# 4. Restart services
make prod-up
```

---

## 💡 Pro Tips

### Tip 1: Use CONCURRENTLY for indexes
```sql
-- ❌ Locks table (avoid!)
CREATE INDEX idx_users_email ON users(email);

-- ✅ Non-blocking
CREATE INDEX CONCURRENTLY idx_users_email ON users(email);
```

### Tip 2: Add columns with defaults in 2 steps
```sql
-- ❌ Slow on large tables (rewrites entire table)
ALTER TABLE users ADD COLUMN status VARCHAR(20) NOT NULL DEFAULT 'active';

-- ✅ Fast (PostgreSQL 11+)
ALTER TABLE users ADD COLUMN status VARCHAR(20);
ALTER TABLE users ALTER COLUMN status SET DEFAULT 'active';
UPDATE users SET status = 'active' WHERE status IS NULL;  -- Backfill
ALTER TABLE users ALTER COLUMN status SET NOT NULL;
```

### Tip 3: Test rollback in staging
```bash
# In staging
make db-migrate
make db-rollback
make db-migrate

# Verify it works before production!
```

### Tip 4: Monitor between phases
```sql
-- Check data consistency
SELECT 
    COUNT(*) as total,
    COUNT(email) as has_old,
    COUNT(email_address) as has_new,
    COUNT(*) FILTER (WHERE email != email_address) as mismatched
FROM users;
```

### Tip 5: Document WHY not just WHAT
```sql
-- ❌ Bad comment
-- Add email_address column

-- ✅ Good comment
-- Add email_address column (Phase 1 of 3)
-- Replacing 'email' column for better naming consistency
-- Related ticket: PROJ-123
-- Phase 2 (remove old column) scheduled for 2025-11-10
```

---

## 📚 Common Patterns

### Pattern: Add Column
```sql
ALTER TABLE users ADD COLUMN phone VARCHAR(20);
```

### Pattern: Add Column with Default
```sql
ALTER TABLE users ADD COLUMN status VARCHAR(20) DEFAULT 'active';
```

### Pattern: Add Index (Safe)
```sql
CREATE INDEX CONCURRENTLY idx_users_email ON users(email);
```

### Pattern: Rename Column (3-Phase)
```sql
-- Phase 1: Add
ALTER TABLE users ADD COLUMN new_name TYPE;
UPDATE users SET new_name = old_name;

-- Phase 2: Wait & Monitor

-- Phase 3: Remove
ALTER TABLE users DROP COLUMN old_name;
```

### Pattern: Change Type (3-Phase)
```sql
-- Phase 1: Add
ALTER TABLE users ADD COLUMN age_new INTEGER;
UPDATE users SET age_new = age::INTEGER;

-- Phase 2: Wait & Monitor

-- Phase 3: Remove
ALTER TABLE users DROP COLUMN age;
ALTER TABLE users RENAME COLUMN age_new TO age;
```

### Pattern: Drop Table (3-Phase)
```sql
-- Phase 1: Stop writes (remove INSERT/UPDATE in code)
-- Deploy code

-- Phase 2: Stop reads (remove SELECT in code)
-- Deploy code
-- Wait 1-2 weeks

-- Phase 3: Drop table
DROP TABLE old_table;
```

---

## 🎯 When to Use Each Strategy

| Strategy | Use When | Downtime | Complexity |
|----------|----------|----------|------------|
| **Backward-Compatible** | Additive changes | 0 min | Low |
| **3-Phase** | Removals, Renames | 0 min | Medium |
| **Maintenance Window** | Complex changes | 5-10 min | Low |
| **Blue-Green** | Critical systems | 0 min | High |

---

## 📖 Further Reading

- [ZERO_DOWNTIME_MIGRATIONS.md](./ZERO_DOWNTIME_MIGRATIONS.md) - Full guide
- [MIGRATION_EXAMPLE_3_PHASE.md](./MIGRATION_EXAMPLE_3_PHASE.md) - Detailed example
- [db-prod.mk](./db-prod.mk) - Production Makefile targets