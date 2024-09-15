package app

// TwoFA represents a 2FA record.
type TwoFA struct {
	ID       int
	Priority int
	Logo     string
	Name     string
	Secret   string
	Domain   string
}

// AddTwoFA adds a new 2FA record to the database and returns the ID of the new record.
func (d *Database) AddTwoFA(priority int, logo, name, secret, domain string) (int64, error) {
	insertSQL := `INSERT INTO two_fa (priority, logo, name, secret,domain) VALUES (?, ?, ?, ?,?)`
	result, err := d.db.Exec(insertSQL, priority, logo, name, secret, domain)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

// GetTwoFAs retrieves all 2FA records from the database.
func (d *Database) GetTwoFAs() ([]TwoFA, error) {
	rows, err := d.db.Query("SELECT id, priority, logo, name, secret, domain FROM two_fa")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var twoFAs []TwoFA
	for rows.Next() {
		var twoFA TwoFA
		if err := rows.Scan(&twoFA.ID, &twoFA.Priority, &twoFA.Logo, &twoFA.Name, &twoFA.Secret, &twoFA.Domain); err != nil {
			return nil, err
		}
		twoFAs = append(twoFAs, twoFA)
	}

	return twoFAs, nil
}

func (d *Database) UpdateTwoFAFields(id int, fields map[string]interface{}) error {
	if len(fields) == 0 {
		return nil
	}

	updateSQL := "UPDATE two_fa SET "
	args := []interface{}{}
	for column, value := range fields {
		updateSQL += column + " = ?, "
		args = append(args, value)
	}

	updateSQL = updateSQL[:len(updateSQL)-2] + " WHERE id = ?"
	args = append(args, id)

	_, err := d.db.Exec(updateSQL, args...)
	return err
}

// DeleteTwoFA deletes a 2FA record from the database.
func (d *Database) DeleteTwoFA(id int) error {
	deleteSQL := `DELETE FROM two_fa WHERE id = ?`
	_, err := d.db.Exec(deleteSQL, id)
	return err
}

// IsSecretExists checks if a given secret already exists in the database.
func (d *Database) IsSecretExists(secret string) (bool, error) {
	var count int
	query := "SELECT COUNT(*) FROM two_fa WHERE secret = ?"
	err := d.db.QueryRow(query, secret).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (d *Database) SearchTwoFAByName(name string) ([]TwoFA, error) {
	query := "SELECT id, priority, logo, name, secret, domain FROM two_fa WHERE name LIKE ?"
	rows, err := d.db.Query(query, "%"+name+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var twoFAs []TwoFA
	for rows.Next() {
		var twoFA TwoFA
		if err := rows.Scan(&twoFA.ID, &twoFA.Priority, &twoFA.Logo, &twoFA.Name, &twoFA.Secret, &twoFA.Domain); err != nil {
			return nil, err
		}
		twoFAs = append(twoFAs, twoFA)
	}

	return twoFAs, nil
}
