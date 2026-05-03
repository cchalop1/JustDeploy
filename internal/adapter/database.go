package adapter

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"os"

	"cchalop1.com/deploy/internal"
	"cchalop1.com/deploy/internal/domain"
	_ "modernc.org/sqlite"
)

type DatabaseAdapter struct {
	db *sql.DB
}

func NewDatabaseAdapter() *DatabaseAdapter {
	return &DatabaseAdapter{}
}

const schema = `
CREATE TABLE IF NOT EXISTS server (
	id           TEXT PRIMARY KEY,
	name         TEXT NOT NULL DEFAULT '',
	ip           TEXT NOT NULL DEFAULT '',
	port         TEXT NOT NULL DEFAULT '',
	domain       TEXT NOT NULL DEFAULT '',
	password     TEXT,
	ssh_key      TEXT,
	email        TEXT NOT NULL DEFAULT '',
	status       TEXT NOT NULL DEFAULT '',
	use_https    INTEGER NOT NULL DEFAULT 0,
	created_date TEXT NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS services (
	id              TEXT PRIMARY KEY,
	name            TEXT NOT NULL DEFAULT '',
	full_name       TEXT NOT NULL DEFAULT '',
	type            TEXT NOT NULL DEFAULT '',
	status          TEXT NOT NULL DEFAULT '',
	url             TEXT NOT NULL DEFAULT '',
	image_name      TEXT NOT NULL DEFAULT '',
	image_url       TEXT NOT NULL DEFAULT '',
	docker_hub_url  TEXT NOT NULL DEFAULT '',
	current_path    TEXT NOT NULL DEFAULT '',
	envs            TEXT NOT NULL DEFAULT '[]',
	cmd             TEXT NOT NULL DEFAULT '[]',
	expose_settings TEXT NOT NULL DEFAULT '{}',
	last_commit     TEXT NOT NULL DEFAULT '{}'
);

CREATE TABLE IF NOT EXISTS settings (
	id                  INTEGER PRIMARY KEY CHECK (id = 1),
	admin_email         TEXT NOT NULL DEFAULT '',
	admin_password_hash TEXT NOT NULL DEFAULT '',
	jwt_secret          TEXT NOT NULL DEFAULT '',
	github_token        TEXT NOT NULL DEFAULT '',
	github_app_settings TEXT NOT NULL DEFAULT '{}'
);

CREATE TABLE IF NOT EXISTS logs (
	id         INTEGER PRIMARY KEY AUTOINCREMENT,
	service_id TEXT NOT NULL,
	date       TEXT NOT NULL,
	message    TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_logs_service_id ON logs(service_id);
`

func (d *DatabaseAdapter) Init() {
	if err := os.MkdirAll(internal.JUSTDEPLOY_FOLDER, 0755); err != nil {
		log.Fatalf("failed to create config directory: %v", err)
	}

	db, err := sql.Open("sqlite", internal.DATABASE_SQLITE_PATH)
	if err != nil {
		log.Fatalf("failed to open sqlite database: %v", err)
	}

	db.SetMaxOpenConns(1)

	if _, err := db.Exec("PRAGMA journal_mode=WAL; PRAGMA foreign_keys=ON;"); err != nil {
		log.Fatalf("failed to set sqlite pragmas: %v", err)
	}

	if _, err := db.Exec(schema); err != nil {
		log.Fatalf("failed to create tables: %v", err)
	}

	if _, err := db.Exec(`INSERT OR IGNORE INTO settings (id) VALUES (1)`); err != nil {
		log.Fatalf("failed to initialize settings row: %v", err)
	}

	d.db = db
}

// --- Server ---

func (d *DatabaseAdapter) GetServer() domain.Server {
	row := d.db.QueryRow(`SELECT id, name, ip, port, domain, password, ssh_key, email, status, use_https, created_date FROM server LIMIT 1`)

	var s domain.Server
	var password, sshKey sql.NullString
	var useHttps int
	var createdDate string

	err := row.Scan(&s.Id, &s.Name, &s.Ip, &s.Port, &s.Domain, &password, &sshKey, &s.Email, &s.Status, &useHttps, &createdDate)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Server{}
	}
	if err != nil {
		log.Printf("GetServer scan error: %v", err)
		return domain.Server{}
	}

	if password.Valid {
		s.Password = &password.String
	}
	if sshKey.Valid {
		s.SshKey = &sshKey.String
	}
	s.UseHttps = useHttps == 1

	if createdDate != "" {
		_ = s.CreatedDate.UnmarshalText([]byte(createdDate))
	}

	return s
}

func (d *DatabaseAdapter) SaveServer(s domain.Server) error {
	var createdDate string
	if !s.CreatedDate.IsZero() {
		b, _ := s.CreatedDate.MarshalText()
		createdDate = string(b)
	}

	useHttps := 0
	if s.UseHttps {
		useHttps = 1
	}

	_, err := d.db.Exec(`
		INSERT INTO server (id, name, ip, port, domain, password, ssh_key, email, status, use_https, created_date)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			name=excluded.name, ip=excluded.ip, port=excluded.port,
			domain=excluded.domain, password=excluded.password, ssh_key=excluded.ssh_key,
			email=excluded.email, status=excluded.status, use_https=excluded.use_https,
			created_date=excluded.created_date
	`, s.Id, s.Name, s.Ip, s.Port, s.Domain, s.Password, s.SshKey, s.Email, s.Status, useHttps, createdDate)
	return err
}

func (d *DatabaseAdapter) DeleteServer(s domain.Server) error {
	_, err := d.db.Exec(`DELETE FROM server WHERE id = ?`, s.Id)
	return err
}

// --- Services ---

func (d *DatabaseAdapter) GetServices() []domain.Service {
	rows, err := d.db.Query(`SELECT id, name, full_name, type, status, url, image_name, image_url, docker_hub_url, current_path, envs, cmd, expose_settings, last_commit FROM services`)
	if err != nil {
		log.Printf("GetServices query error: %v", err)
		return nil
	}
	defer rows.Close()

	services := []domain.Service{}
	for rows.Next() {
		s, err := scanService(rows)
		if err != nil {
			log.Printf("GetServices scan error: %v", err)
			continue
		}
		services = append(services, s)
	}
	return services
}

func (d *DatabaseAdapter) GetServiceById(id string) (*domain.Service, error) {
	row := d.db.QueryRow(`SELECT id, name, full_name, type, status, url, image_name, image_url, docker_hub_url, current_path, envs, cmd, expose_settings, last_commit FROM services WHERE id = ?`, id)
	s, err := scanService(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("service not found")
	}
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (d *DatabaseAdapter) SaveService(s domain.Service) error {
	envs, _ := json.Marshal(s.Envs)
	cmd, _ := json.Marshal(s.Cmd)
	expose, _ := json.Marshal(s.ExposeSettings)
	commit, _ := json.Marshal(s.LastCommit)

	_, err := d.db.Exec(`
		INSERT INTO services (id, name, full_name, type, status, url, image_name, image_url, docker_hub_url, current_path, envs, cmd, expose_settings, last_commit)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			name=excluded.name, full_name=excluded.full_name, type=excluded.type,
			status=excluded.status, url=excluded.url, image_name=excluded.image_name,
			image_url=excluded.image_url, docker_hub_url=excluded.docker_hub_url,
			current_path=excluded.current_path, envs=excluded.envs, cmd=excluded.cmd,
			expose_settings=excluded.expose_settings, last_commit=excluded.last_commit
	`, s.Id, s.Name, s.FullName, s.Type, s.Status, s.Url, s.ImageName, s.ImageUrl,
		s.DockerHubUrl, s.CurrentPath, string(envs), string(cmd), string(expose), string(commit))
	return err
}

func (d *DatabaseAdapter) DeleteServiceById(id string) error {
	_, err := d.db.Exec(`DELETE FROM services WHERE id = ?`, id)
	return err
}

// --- Settings ---

func (d *DatabaseAdapter) GetSettings() domain.Settings {
	row := d.db.QueryRow(`SELECT admin_email, admin_password_hash, jwt_secret, github_token, github_app_settings FROM settings WHERE id = 1`)

	var appSettingsJSON string
	var s domain.Settings

	err := row.Scan(&s.AdminEmail, &s.AdminPasswordHash, &s.JwtSecret, &s.GithubToken, &appSettingsJSON)
	if err != nil {
		log.Printf("GetSettings scan error: %v", err)
		return domain.Settings{}
	}

	if appSettingsJSON != "" && appSettingsJSON != "{}" {
		_ = json.Unmarshal([]byte(appSettingsJSON), &s.GithubAppSettings)
	}

	return s
}

func (d *DatabaseAdapter) SaveSettings(s domain.Settings) error {
	appSettings, _ := json.Marshal(s.GithubAppSettings)

	_, err := d.db.Exec(`
		UPDATE settings SET
			admin_email=?, admin_password_hash=?, jwt_secret=?, github_token=?, github_app_settings=?
		WHERE id = 1
	`, s.AdminEmail, s.AdminPasswordHash, s.JwtSecret, s.GithubToken, string(appSettings))
	return err
}

func (d *DatabaseAdapter) SaveInstallationToken(installationId string, token string) error {
	_, err := d.db.Exec(`UPDATE settings SET github_token=? WHERE id = 1`, token)
	return err
}

// --- Logs ---

func (d *DatabaseAdapter) CreateDeployLogsFileIfNot(serviceId string) error {
	return nil
}

func (d *DatabaseAdapter) SaveLogs(serviceId string, logs []domain.Logs) error {
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec(`DELETE FROM logs WHERE service_id = ?`, serviceId); err != nil {
		return err
	}

	stmt, err := tx.Prepare(`INSERT INTO logs (service_id, date, message) VALUES (?, ?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, l := range logs {
		if _, err := stmt.Exec(serviceId, l.Date, l.Message); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (d *DatabaseAdapter) GetLogs(serviceId string) ([]domain.Logs, error) {
	rows, err := d.db.Query(`SELECT date, message FROM logs WHERE service_id = ? ORDER BY id ASC`, serviceId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []domain.Logs
	for rows.Next() {
		var l domain.Logs
		if err := rows.Scan(&l.Date, &l.Message); err != nil {
			return nil, err
		}
		logs = append(logs, l)
	}
	return logs, nil
}

func (d *DatabaseAdapter) DeleteLogFile(serviceId string) error {
	_, err := d.db.Exec(`DELETE FROM logs WHERE service_id = ?`, serviceId)
	return err
}

// --- helpers ---

type scanner interface {
	Scan(dest ...any) error
}

func scanService(row scanner) (domain.Service, error) {
	var s domain.Service
	var envsJSON, cmdJSON, exposeJSON, commitJSON string

	err := row.Scan(
		&s.Id, &s.Name, &s.FullName, &s.Type, &s.Status, &s.Url,
		&s.ImageName, &s.ImageUrl, &s.DockerHubUrl, &s.CurrentPath,
		&envsJSON, &cmdJSON, &exposeJSON, &commitJSON,
	)
	if err != nil {
		return domain.Service{}, err
	}

	_ = json.Unmarshal([]byte(envsJSON), &s.Envs)
	_ = json.Unmarshal([]byte(cmdJSON), &s.Cmd)
	_ = json.Unmarshal([]byte(exposeJSON), &s.ExposeSettings)
	_ = json.Unmarshal([]byte(commitJSON), &s.LastCommit)

	return s, nil
}
