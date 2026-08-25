package storage

const schemaSQL = `
CREATE TABLE IF NOT EXISTS clients (
 id TEXT PRIMARY KEY, name TEXT NOT NULL, phone TEXT NOT NULL, address TEXT NOT NULL,
 preferred_channel TEXT NOT NULL, active INTEGER NOT NULL, created_at TEXT NOT NULL
);
CREATE TABLE IF NOT EXISTS service_types (
 id TEXT PRIMARY KEY, name TEXT NOT NULL, description TEXT NOT NULL, default_days INTEGER NOT NULL, active INTEGER NOT NULL
);
CREATE TABLE IF NOT EXISTS caregivers (
 id TEXT PRIMARY KEY, name TEXT NOT NULL, phone TEXT NOT NULL, skill_tags TEXT NOT NULL, available INTEGER NOT NULL
);
CREATE TABLE IF NOT EXISTS follow_up_records (
 id TEXT PRIMARY KEY, client_id TEXT NOT NULL, client_name TEXT NOT NULL, service_type_id TEXT NOT NULL,
 service_type_name TEXT NOT NULL, caregiver_id TEXT NOT NULL, caregiver_name TEXT NOT NULL,
 visit_date TEXT NOT NULL, next_follow_up TEXT NOT NULL, score INTEGER NOT NULL,
 comment TEXT NOT NULL, improvement TEXT NOT NULL, status TEXT NOT NULL, notes TEXT NOT NULL, updated_at TEXT NOT NULL
);
CREATE TABLE IF NOT EXISTS reminder_settings (
 id TEXT PRIMARY KEY, days_before INTEGER NOT NULL, channel TEXT NOT NULL, enabled INTEGER NOT NULL,
 quiet_start INTEGER NOT NULL, quiet_end INTEGER NOT NULL
);
CREATE TABLE IF NOT EXISTS audit_entries (
 id INTEGER PRIMARY KEY AUTOINCREMENT, record_id TEXT NOT NULL, action TEXT NOT NULL,
 actor TEXT NOT NULL, detail TEXT NOT NULL, created_at TEXT NOT NULL
);`
