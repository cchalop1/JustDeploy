package adapter

type DatabaseAdapter struct {
}

func NewDatabaseAdapter() *DatabaseAdapter {
	return &DatabaseAdapter{}
}

func (d *DatabaseAdapter) GetState() {

}

func (d *DatabaseAdapter) SaveState() {

}
