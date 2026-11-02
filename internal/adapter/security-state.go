package adapter

import (
	"cchalop1.com/deploy/internal/domain"
	"database/sql"
	"time"
)

// InitializeAdmin performs the check and writes in one transaction.
func (d *DatabaseAdapter) InitializeAdmin(email, hash, secret, serverDomain string) error {
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	result, err := tx.Exec(`UPDATE settings SET admin_email=?, admin_password_hash=?, jwt_secret=? WHERE id=1 AND admin_email=''`, email, hash, secret)
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return domain.ErrAlreadyConfigured
	}
	result, err = tx.Exec(`UPDATE server SET domain=?`, serverDomain)
	if err != nil {
		return err
	}
	count, err = result.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return sql.ErrNoRows
	}
	return tx.Commit()
}

func (d *DatabaseAdapter) ClaimDelivery(id string) (bool, error) {
	// Bounded retention; GitHub deliveries are deduplicated for seven days.
	if _, err := d.db.Exec(`DELETE FROM webhook_deliveries WHERE created_at < ?`, time.Now().Add(-7*24*time.Hour).Unix()); err != nil {
		return false, err
	}
	r, err := d.db.Exec(`INSERT OR IGNORE INTO webhook_deliveries (id, created_at) VALUES (?, ?)`, id, time.Now().Unix())
	if err != nil {
		return false, err
	}
	n, err := r.RowsAffected()
	return n == 1, err
}
func (d *DatabaseAdapter) ReleaseDelivery(id string) error {
	_, err := d.db.Exec(`DELETE FROM webhook_deliveries WHERE id=?`, id)
	return err
}
func (d *DatabaseAdapter) RotateJWTSecret(secret string) error {
	_, err := d.db.Exec(`UPDATE settings SET jwt_secret=? WHERE id=1`, secret)
	return err
}
func (d *DatabaseAdapter) Close() error { return d.db.Close() }
