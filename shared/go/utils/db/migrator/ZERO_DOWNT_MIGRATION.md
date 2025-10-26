# Zero-Downtime Production Migrations Guide

## 🎯 Das Problem

**Development:**
```bash
make dev         # Services laufen
make db-migrate  # Migration ändert Schema
# ✅ Kein Problem - Services starten neu
```

**Production:**
```bash
# Services laufen (1000 User online)
make prod-db-migrate  # DROP TABLE users
# ❌ BOOM! Services crashen!
# ❌ Gap: 5-10 Minuten bis Deployment
```

---

## 📊 Production Migration Strategien

### Übersicht

| Strategie | Downtime | Komplexität | Use Case |
|-----------|----------|-------------|----------|
| **Backward-Compatible** | 0 min | Medium | Additive Changes |
| **Maintenance Window** | 5-10 min | Low | Breaking Changes |
| **Blue-Green** | 0 min | High | Critical Systems |

---

## ✅ Strategie 1: Backward-Compatible Migrations (Empfohlen)

### Prinzip: Migrations dürfen alte Services NICHT brechen

### Regel-Set

| ✅ Erlaubt (Safe) | ❌ Verboten (Breaking) |
|------------------|----------------------|
| ADD COLUMN | DROP COLUMN |
| CREATE TABLE | DROP TABLE |
| CREATE INDEX | RENAME COLUMN |
| ADD CONSTRAINT (nullable) | ALTER COLUMN TYPE |
| INSERT data | DELETE data |

---

### 🎓 Beispiel 1: Spalte hinzufügen

#### ❌ FALSCH:
```sql
-- Migration
ALTER TABLE users ADD COLUMN phone VARCHAR(20) NOT NULL;
-- ❌ Alte Services schreiben NULL → Fehler!
```

#### ✅ RICHTIG:
```sql
-- Migration
ALTER TABLE users ADD COLUMN phone VARCHAR(20);
-- ✅ NULL ist erlaubt → Alte Services funktionieren!

-- Optional: Default-Wert
ALTER TABLE users ADD COLUMN phone VARCHAR(20) DEFAULT '';
```

---

### 🎓 Beispiel 2: Spalte umbenennen (3-Phasen)

**Anforderung:** `email` → `email_address`

#### Phase 1: Neue Spalte hinzufügen (Migration 001)

```sql
-- 001_add_email_address.up.sql
BEGIN;

-- Neue Spalte hinzufügen
ALTER TABLE users ADD COLUMN email_address VARCHAR(255);

-- Daten kopieren
UPDATE users SET email_address = email WHERE email_address IS NULL;

-- Index für Performance
CREATE INDEX idx_users_email_address ON users(email_address);

COMMIT;
```

```sql
-- 001_add_email_address.down.sql
DROP INDEX IF EXISTS idx_users_email_address;
ALTER TABLE users DROP COLUMN IF EXISTS email_address;
```

**Code-Änderung (backward-compatible):**
```go
// models/user.go
type User struct {
    ID           int    `db:"id"`
    Email        string `db:"email"`          // Alt (noch verwendet)
    EmailAddress string `db:"email_address"`  // Neu
}

// GetEmail() - Hilfs-Methode während Transition
func (u *User) GetEmail() string {
    if u.EmailAddress != "" {
        return u.EmailAddress
    }
    return u.Email  // Fallback
}

// SetEmail() - Schreibt in beide Spalten
func (u *User) SetEmail(email string) {
    u.Email = email          // Für alte Services
    u.EmailAddress = email   // Für neue Services
}
```

**Deployment:**
```bash
# 1. Migration deployen
make prod-db-migrate-safe

# 2. Code deployen (nutzt beide Spalten)
make prod-deploy

# 3. Warten (24-48h) bis alle alten Pods tot sind
# Alte Services nutzen weiter 'email'
# Neue Services nutzen 'email_address'
```

#### Phase 2: Warten (1-2 Tage)

**Warum warten?**
- Alte Pods sterben langsam ab
- Rolling-Updates brauchen Zeit
- Rollback-Option offen halten

**Monitoring:**
```bash
# Prüfe ob alte Services noch laufen
kubectl get pods -l version=old

# Prüfe Logs für Fehler
make logs | grep -i "email"
```

#### Phase 3: Alte Spalte entfernen (Migration 002)

**Erst wenn:**
- ✅ Alle Services updated sind
- ✅ Keine Fehler in Logs
- ✅ Mindestens 24h vergangen

```sql
-- 002_remove_old_email.up.sql
BEGIN;

-- Alte Spalte entfernen
ALTER TABLE users DROP COLUMN email;

-- Index cleanup (falls vorhanden)
DROP INDEX IF EXISTS idx_users_email;

COMMIT;
```

```sql
-- 002_remove_old_email.down.sql
ALTER TABLE users ADD COLUMN email VARCHAR(255);
UPDATE users SET email = email_address;
CREATE INDEX idx_users_email ON users(email);
```

**Code-Cleanup:**
```go
type User struct {
    ID           int    `db:"id"`
    EmailAddress string `db:"email_address"`  // Nur noch diese
}

func (u *User) GetEmail() string {
    return u.EmailAddress
}
```

---

### 🎓 Beispiel 3: Tabelle löschen (3-Phasen)

**Anforderung:** `old_logs` Tabelle nicht mehr benötigt

#### Phase 1: Schreibzugriffe entfernen

```go
// Code: Entferne alle INSERT/UPDATE Statements
// func SaveLog() {} → DELETE
```

**Deployment:**
```bash
make prod-deploy
# ✅ Keine neuen Daten werden mehr geschrieben
```

#### Phase 2: Lesezugriffe entfernen (nach 1 Woche)

```go
// Code: Entferne alle SELECT Statements
// func GetLogs() {} → DELETE
```

**Deployment:**
```bash
make prod-deploy
# ✅ Tabelle wird nicht mehr genutzt
```

#### Phase 3: Tabelle löschen (nach 2 Wochen)

```sql
-- 003_drop_old_logs.up.sql
DROP TABLE IF EXISTS old_logs CASCADE;
```

```bash
# Erst nach Sicherheits-Periode!
make prod-db-migrate-safe
```

---

## ✅ Strategie 2: Maintenance Window

### Wann nutzen?

- ❌ Breaking Changes unvermeidbar
- ❌ Komplexe Schema-Änderungen
- ❌ Daten-Migration mit Logik
- ✅ Niedrige Traffic-Zeiten verfügbar

### Prozess

```bash
# 1. Ankündigung (1 Woche vorher)
echo "Maintenance: Sunday 02:00-02:15 UTC"

# 2. Pre-Deployment (1 Tag vorher)
make prod-db-backup
# Test in Staging!

# 3. Maintenance Window Start
make prod-db-maintenance
# → Stoppt Services
# → Backup
# → Migration
# → Startet Services

# 4. Post-Deployment
# Monitor für 1 Stunde
make prod-db-health
make logs
```

---

## ✅ Strategie 3: Blue-Green Deployment

### Setup

```yaml
# Production Environment
Blue (Current):
  - gateway-blue-1, gateway-blue-2, gateway-blue-3
  - postgres-blue (DB Version 1)
  
Green (New):
  - gateway-green-1, gateway-green-2, gateway-green-3
  - postgres-green (DB Version 2)
```

### Prozess

```bash
# 1. Green Environment aufsetzen
make prod-deploy-green

# 2. Migration auf Green
make prod-db-migrate ENV=green

# 3. Smoke Tests
curl https://green.api.example.com/health
curl https://green.api.example.com/users

# 4. Traffic umleiten (0 Downtime!)
# Load Balancer: 100% Blue → 100% Green

# 5. Monitor
# Bei Problemen: Zurück zu Blue!

# 6. Cleanup (nach 24h)
make prod-down-blue
```

---

## 📋 Production Deployment Checklist

### Pre-Deployment

```
☐ Migrations in Staging getestet
☐ Backward-compatibility geprüft
☐ Breaking Changes dokumentiert
☐ Rollback-Plan vorhanden
☐ Team notified (bei Maintenance)
☐ Monitoring Dashboard bereit
☐ Backup-Strategie definiert
```

### Deployment

```
☐ Backup erstellt
☐ Breaking Changes geprüft (make prod-db-check-breaking)
☐ Migration durchgeführt
☐ Schema verifiziert (make prod-db-verify)
☐ Health Checks erfolgreich
☐ Logs geprüft (keine Errors)
```

### Post-Deployment

```
☐ Monitoring für 1h aktiv
☐ Error-Rate normal
☐ Performance normal
☐ Backup für 7 Tage behalten
☐ Team notifiziert (Success)
```

---

## 🆘 Rollback Strategien

### Scenario 1: Migration fehlgeschlagen

```bash
# Migration crashed
make prod-db-rollback-safe

# Verify
make prod-db-status
make prod-db-health
```

### Scenario 2: Services crashen nach Migration

```bash
# 1. Rollback Migration
make prod-db-rollback-safe

# 2. Rollback Deployment
kubectl rollout undo deployment/gateway

# 3. Verify
make prod-db-health
curl http://api.example.com/health
```

### Scenario 3: Daten-Korruption

```bash
# 1. Stop Services (prevent further damage)
make prod-down

# 2. Restore from Backup
make db-restore FILE=backups/prod_backup_20251026_020000.sql

# 3. Verify Data
make db-connect
SELECT COUNT(*) FROM users;

# 4. Restart Services
make prod-up
```

---

## 🔧 Best Practices

### DO ✅

1. **Immer backward-compatible migrations**
   ```sql
   -- ✅ Additive
   ALTER TABLE users ADD COLUMN phone VARCHAR(20);
   ```

2. **Defaults für neue Spalten**
   ```sql
   -- ✅ Alte Services setzen nichts → Default wird genutzt
   ALTER TABLE users ADD COLUMN status VARCHAR(20) DEFAULT 'active';
   ```

3. **NULLable neue Spalten**
   ```sql
   -- ✅ Alte Services setzen NULL → OK
   ALTER TABLE users ADD COLUMN bio TEXT;
   ```

4. **3-Phasen für Breaking Changes**
   - Phase 1: Add new
   - Phase 2: Use both
   - Phase 3: Remove old

5. **Backup vor JEDER Production-Migration**
   ```bash
   make prod-db-backup
   ```

### DON'T ❌

1. **Breaking Changes ohne Plan**
   ```sql
   -- ❌ Sofortiger Crash!
   DROP TABLE users;
   ```

2. **NOT NULL ohne Default**
   ```sql
   -- ❌ Alte Services schreiben NULL → Error!
   ALTER TABLE users ADD COLUMN email VARCHAR(255) NOT NULL;
   ```

3. **Rename ohne Transition**
   ```sql
   -- ❌ Alte Services kennen neuen Namen nicht!
   ALTER TABLE users RENAME COLUMN email TO email_address;
   ```

4. **Data-Migration mit DOWN-Queries**
   ```sql
   -- ❌ Data-Migration ist nicht rollbar!
   UPDATE users SET status = 'active';
   ```

5. **Production-Migration ohne Test**
   ```bash
   # ❌ NEVER!
   make prod-db-migrate  # ohne Staging-Test
   ```

---

## 📊 Migration Patterns

### Pattern 1: Add Column

```sql
-- ✅ Safe
ALTER TABLE users ADD COLUMN phone VARCHAR(20);

-- ✅ Even safer
ALTER TABLE users ADD COLUMN phone VARCHAR(20) DEFAULT '';
```

### Pattern 2: Add Index (Non-Blocking)

```sql
-- ❌ Blocking (locks table!)
CREATE INDEX idx_users_email ON users(email);

-- ✅ Non-blocking (PostgreSQL 11+)
CREATE INDEX CONCURRENTLY idx_users_email ON users(email);
```

### Pattern 3: Change Column Type (3-Phase)

**Phase 1:** Add new column
```sql
ALTER TABLE users ADD COLUMN age_new INTEGER;
UPDATE users SET age_new = age::INTEGER WHERE age_new IS NULL;
```

**Phase 2:** Use both (Code update)
```go
// Read from both
age := user.AgeNew
if age == 0 {
    age = user.Age
}
```

**Phase 3:** Remove old
```sql
ALTER TABLE users DROP COLUMN age;
ALTER TABLE users RENAME COLUMN age_new TO age;
```

---

## 🎯 Quick Reference

### Development
```bash
make dev         # Start
make db-migrate  # Migrate
make db-rollback # Rollback
```

### Production (Safe)
```bash
make prod-db-backup              # Backup first!
make prod-db-check-breaking      # Check for breaking changes
make prod-db-migrate-safe        # Migrate with checks
make prod-db-verify              # Verify schema
make prod-db-health              # Health check
```

### Production (With Downtime)
```bash
make prod-db-maintenance         # Full maintenance window
# → Stops services
# → Backup
# → Migrate
# → Start services
```

### Rollback
```bash
make prod-db-rollback-safe       # Rollback last migration
make db-restore FILE=backup.sql  # Restore from backup
```

---

## 📚 Resources

- [PostgreSQL Concurrent Index Creation](https://www.postgresql.org/docs/current/sql-createindex.html#SQL-CREATEINDEX-CONCURRENTLY)
- [Blue-Green Deployments](https://martinfowler.com/bliki/BlueGreenDeployment.html)
- [Database Refactoring](https://www.martinfowler.com/articles/evodb.html)
- [Zero-Downtime Migrations](https://www.braintreepayments.com/blog/safe-operations-for-high-volume-postgresql/)