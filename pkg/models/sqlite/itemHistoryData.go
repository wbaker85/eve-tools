package sqlite

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/wbaker85/eve-tools/pkg/models"
)

// ItemHistoryDataModel deals with average volume data for an item
type ItemHistoryDataModel struct {
	DB *sql.DB
}

// LoadData loads a slice of averages into the database based on a regionID
func (i *ItemHistoryDataModel) LoadData(regionID int, averages []models.ItemHistoryData) {
	i.init(regionID)
	i.addMany(regionID, averages)
}

func (i *ItemHistoryDataModel) addMany(regionID int, averages []models.ItemHistoryData) {
	if len(averages) == 0 {
		return
	}

	var b strings.Builder
	stmt := fmt.Sprintf(`INSERT INTO "%d_averages" (item_id, num_days, orders_avg, volume_avg) VALUES `, regionID)
	b.WriteString(stmt)

	for _, item := range averages {
		sqlStr := `(%d, "%d", %d, %d),`
		b.WriteString(fmt.Sprintf(sqlStr, item.ItemID, item.NumDays, item.OrdersAvg, item.VolumeAvg))
	}

	stmt = b.String()
	stmt = stmt[:len(stmt)-1]

	_, err := i.DB.Exec(stmt)
	if err != nil {
		log.Fatal(err)
	}
}

func (i *ItemHistoryDataModel) init(regionID int) {
	stmt := `CREATE TABLE "%d_averages" (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		item_id INT,
		num_days INT,
		orders_avg INT,
		volume_avg INT
	)`

	drop := `DROP TABLE "%d_averages"`

	i.DB.Exec(fmt.Sprintf(drop, regionID))

	_, err := i.DB.Exec(fmt.Sprintf(stmt, regionID))
	if err != nil {
		log.Fatal(err)
	}
}

// GetVolumesForRegion returns a list of all items in the region (in the database)
// with their volumes
func (i *ItemHistoryDataModel) GetVolumesForRegion(regionID int) map[int]models.ItemHistoryData {
	stmt := `SELECT item_id, num_days, orders_avg, volume_avg FROM "%d_averages"`

	rows, err := i.DB.Query(fmt.Sprintf(stmt, regionID))
	if err != nil {
		log.Fatal(err)
	}

	output := make(map[int]models.ItemHistoryData)

	for rows.Next() {
		i := models.ItemHistoryData{}
		err = rows.Scan(
			&i.ItemID,
			&i.NumDays,
			&i.OrdersAvg,
			&i.VolumeAvg,
		)

		if err != nil {
			log.Fatal(err)
		}

		output[i.ItemID] = i
	}

	return output
}
