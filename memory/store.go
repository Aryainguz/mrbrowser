package memory

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mrbrowser/mrbrowser/telemetry"
)

// Store is the SQLite-backed persistence layer.
type Store struct {
	db     *sql.DB
	log    *telemetry.Logger
	dbPath string
}

// Open opens (or creates) the SQLite database at the given path.
func Open(path string) (*Store, error) {
	db, err := sql.Open("sqlite3", path+"?_journal=WAL&_timeout=5000")
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	db.SetMaxOpenConns(1) // SQLite doesn't support concurrent writes
	db.SetMaxIdleConns(1)

	s := &Store{db: db, log: telemetry.New("memory"), dbPath: path}
	if err := s.migrate(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("migrate db: %w", err)
	}

	s.log.Success("Memory store opened", telemetry.F("path", path))
	return s, nil
}

// DBPath returns the file path to the SQLite database.
func (s *Store) DBPath() string {
	return s.dbPath
}

// Close closes the database connection.
func (s *Store) Close() error {
	return s.db.Close()
}

// migrate creates all required tables if they don't exist.
func (s *Store) migrate() error {
	schema := `
CREATE TABLE IF NOT EXISTS element_fingerprints (
	id           TEXT PRIMARY KEY,
	workflow_id  TEXT NOT NULL,
	step_name    TEXT NOT NULL,
	target       TEXT NOT NULL,
	selector     TEXT,
	text         TEXT,
	role         TEXT,
	type         TEXT,
	section      TEXT,
	nearby_text  TEXT,    -- JSON array
	attributes   TEXT,    -- JSON object
	last_seen    DATETIME,
	use_count    INTEGER DEFAULT 0,
	success_count INTEGER DEFAULT 0,
	success_rate REAL DEFAULT 0.0
);

CREATE INDEX IF NOT EXISTS idx_fingerprints_target ON element_fingerprints(workflow_id, target);

CREATE TABLE IF NOT EXISTS workflow_results (
	id            TEXT PRIMARY KEY,
	workflow_name TEXT NOT NULL,
	started_at    DATETIME,
	completed_at  DATETIME,
	duration_ms   INTEGER,
	success       BOOLEAN,
	steps_total   INTEGER DEFAULT 0,
	steps_passed  INTEGER DEFAULT 0,
	error_message TEXT,
	url           TEXT
);

CREATE INDEX IF NOT EXISTS idx_workflow_results_name ON workflow_results(workflow_name);

CREATE TABLE IF NOT EXISTS action_records (
	id            TEXT PRIMARY KEY,
	workflow_id   TEXT NOT NULL,
	step_name     TEXT NOT NULL,
	action        TEXT NOT NULL,
	target        TEXT,
	selector      TEXT,
	success       BOOLEAN,
	duration_ms   INTEGER,
	error_message TEXT,
	recovered     BOOLEAN DEFAULT FALSE,
	recovery_note TEXT,
	created_at    DATETIME
);

CREATE INDEX IF NOT EXISTS idx_action_records_workflow ON action_records(workflow_id);

CREATE TABLE IF NOT EXISTS recovery_events (
	id           TEXT PRIMARY KEY,
	workflow_id  TEXT NOT NULL,
	target       TEXT,
	old_selector TEXT,
	new_selector TEXT,
	old_text     TEXT,
	new_text     TEXT,
	confidence   REAL,
	created_at   DATETIME
);
`
	_, err := s.db.Exec(schema)
	return err
}

// SaveFingerprint inserts or updates an element fingerprint.
func (s *Store) SaveFingerprint(fp *ElementFingerprint) error {
	nearbyJSON, _ := json.Marshal(fp.NearbyText)
	attrsJSON, _ := json.Marshal(fp.Attributes)

	_, err := s.db.Exec(`
		INSERT INTO element_fingerprints
			(id, workflow_id, step_name, target, selector, text, role, type, section, nearby_text, attributes, last_seen, use_count, success_count, success_rate)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			selector=excluded.selector,
			text=excluded.text,
			role=excluded.role,
			type=excluded.type,
			section=excluded.section,
			nearby_text=excluded.nearby_text,
			attributes=excluded.attributes,
			last_seen=excluded.last_seen,
			use_count=excluded.use_count,
			success_count=excluded.success_count,
			success_rate=excluded.success_rate
	`,
		fp.ID, fp.WorkflowID, fp.StepName, fp.Target,
		fp.Selector, fp.Text, fp.Role, fp.Type, fp.Section,
		string(nearbyJSON), string(attrsJSON),
		fp.LastSeen, fp.UseCount, fp.SuccessCount, fp.SuccessRate,
	)
	return err
}

// GetFingerprints returns all fingerprints for the given workflow + target.
func (s *Store) GetFingerprints(workflowID, target string) ([]*ElementFingerprint, error) {
	rows, err := s.db.Query(`
		SELECT id, workflow_id, step_name, target, selector, text, role, type, section,
		       nearby_text, attributes, last_seen, use_count, success_count, success_rate
		FROM element_fingerprints
		WHERE workflow_id = ? AND target = ?
		ORDER BY success_rate DESC, last_seen DESC
	`, workflowID, target)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanFingerprints(rows)
}

// GetAllFingerprintsForWorkflow returns all fingerprints for a workflow.
func (s *Store) GetAllFingerprintsForWorkflow(workflowID string) ([]*ElementFingerprint, error) {
	rows, err := s.db.Query(`
		SELECT id, workflow_id, step_name, target, selector, text, role, type, section,
		       nearby_text, attributes, last_seen, use_count, success_count, success_rate
		FROM element_fingerprints WHERE workflow_id = ?
	`, workflowID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanFingerprints(rows)
}

func scanFingerprints(rows *sql.Rows) ([]*ElementFingerprint, error) {
	var fps []*ElementFingerprint
	for rows.Next() {
		var fp ElementFingerprint
		var nearbyJSON, attrsJSON string
		var lastSeen string

		if err := rows.Scan(
			&fp.ID, &fp.WorkflowID, &fp.StepName, &fp.Target,
			&fp.Selector, &fp.Text, &fp.Role, &fp.Type, &fp.Section,
			&nearbyJSON, &attrsJSON, &lastSeen,
			&fp.UseCount, &fp.SuccessCount, &fp.SuccessRate,
		); err != nil {
			return nil, err
		}

		_ = json.Unmarshal([]byte(nearbyJSON), &fp.NearbyText)
		_ = json.Unmarshal([]byte(attrsJSON), &fp.Attributes)
		fp.LastSeen, _ = time.Parse(time.RFC3339, lastSeen)

		fps = append(fps, &fp)
	}
	return fps, rows.Err()
}

// SaveWorkflowResult persists the result of a complete workflow run.
func (s *Store) SaveWorkflowResult(r *WorkflowResult) error {
	_, err := s.db.Exec(`
		INSERT OR REPLACE INTO workflow_results
			(id, workflow_name, started_at, completed_at, duration_ms, success, steps_total, steps_passed, error_message, url)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		r.ID, r.WorkflowName,
		r.StartedAt.Format(time.RFC3339),
		r.CompletedAt.Format(time.RFC3339),
		r.Duration.Milliseconds(),
		r.Success, r.StepsTotal, r.StepsPassed,
		r.ErrorMessage, r.URL,
	)
	return err
}

// GetWorkflowHistory returns the last N workflow results for the given workflow name.
func (s *Store) GetWorkflowHistory(workflowName string, limit int) ([]*WorkflowResult, error) {
	rows, err := s.db.Query(`
		SELECT id, workflow_name, started_at, completed_at, duration_ms, success, steps_total, steps_passed, error_message, url
		FROM workflow_results
		WHERE workflow_name = ?
		ORDER BY started_at DESC
		LIMIT ?
	`, workflowName, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*WorkflowResult
	for rows.Next() {
		var r WorkflowResult
		var startedAt, completedAt string
		var durationMS int64
		if err := rows.Scan(
			&r.ID, &r.WorkflowName, &startedAt, &completedAt,
			&durationMS, &r.Success, &r.StepsTotal, &r.StepsPassed,
			&r.ErrorMessage, &r.URL,
		); err != nil {
			return nil, err
		}
		r.StartedAt, _ = time.Parse(time.RFC3339, startedAt)
		r.CompletedAt, _ = time.Parse(time.RFC3339, completedAt)
		r.Duration = time.Duration(durationMS) * time.Millisecond
		results = append(results, &r)
	}
	return results, rows.Err()
}

// SaveActionRecord persists a single action's outcome.
func (s *Store) SaveActionRecord(a *ActionRecord) error {
	_, err := s.db.Exec(`
		INSERT OR REPLACE INTO action_records
			(id, workflow_id, step_name, action, target, selector, success, duration_ms, error_message, recovered, recovery_note, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		a.ID, a.WorkflowID, a.StepName, a.Action, a.Target, a.Selector,
		a.Success, a.Duration.Milliseconds(), a.ErrorMessage,
		a.Recovered, a.RecoveryNote, a.CreatedAt.Format(time.RFC3339),
	)
	return err
}

// SaveRecoveryEvent persists a self-healing recovery event.
func (s *Store) SaveRecoveryEvent(e *RecoveryEvent) error {
	_, err := s.db.Exec(`
		INSERT OR REPLACE INTO recovery_events
			(id, workflow_id, target, old_selector, new_selector, old_text, new_text, confidence, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		e.ID, e.WorkflowID, e.Target, e.OldSelector, e.NewSelector,
		e.OldText, e.NewText, e.Confidence, e.CreatedAt.Format(time.RFC3339),
	)
	return err
}

// Stats returns high-level statistics.
type Stats struct {
	TotalWorkflows  int     `json:"total_workflows"`
	SuccessRate     float64 `json:"success_rate"`
	TotalRecoveries int     `json:"total_recoveries"`
	TotalActions    int     `json:"total_actions"`
}

// GetStats returns overall statistics from the memory store.
func (s *Store) GetStats() (*Stats, error) {
	stats := &Stats{}

	row := s.db.QueryRow(`SELECT COUNT(*), AVG(CASE WHEN success THEN 1.0 ELSE 0.0 END) FROM workflow_results`)
	_ = row.Scan(&stats.TotalWorkflows, &stats.SuccessRate)

	row = s.db.QueryRow(`SELECT COUNT(*) FROM recovery_events`)
	_ = row.Scan(&stats.TotalRecoveries)

	row = s.db.QueryRow(`SELECT COUNT(*) FROM action_records`)
	_ = row.Scan(&stats.TotalActions)

	return stats, nil
}
