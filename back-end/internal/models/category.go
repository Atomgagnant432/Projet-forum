package models

type Category struct {
	ID   int
	Name string
}

func GetCategories(db *sql.DB) ([]Category, error) {
	rows, err := db.Query(`
		SELECT id, name
		FROM category
		ORDER BY name ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []Category

	for rows.Next() {
		var category Category

		if err := rows.Scan(&category.ID, &category.Name); err != nil {
			return nil, err
		}

		categories = append(categories, category)
	}

	return categories, rows.Err()
}