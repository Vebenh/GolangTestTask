package db

import (
	"context"
	"fmt"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"strings"
)

type DbConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Dbname   string
}

type PublishInformation struct {
	ID          int    `gorm:"primaryKey;autoIncrement:true"`
	PublishDate string `gorm:"type:date"`
	RecordCount int
	State       int
}

//type SdnList struct {
//	Entries []*SdnEntry `xml:"sdnEntry"`
//}
//type ProgramList struct {
//	Program []*Program `xml:"program"`
//}
//
//type AkaList struct {
//	Aka []*Aka `xml:"aka"`
//}
//type AddressList struct {
//	Address []*Address `xml:"address"`
//}

type SdnEntry struct {
	ID        int        `gorm:"primaryKey;autoIncrement:true"`
	UID       int        `gorm:"index:idx_uid,unique" xml:"uid"`
	FirstName string     `xml:"firstName"`
	LastName  string     `xml:"lastName"`
	SdnType   string     `xml:"sdnType"`
	Programs  []*Program `gorm:"foreignKey:SdnEntryUID" xml:"programList>program"`
	Akas      []*Aka     `gorm:"foreignKey:SdnEntryUID" xml:"akaList>aka"`
	Addresses []*Address `gorm:"foreignKey:SdnEntryUID" xml:"addressList>address"`
}

type Program struct {
	ID          uint `gorm:"primaryKey;autoIncrement:true"`
	SdnEntryUID int  `xml:"-"`
	Program     string
}

type Aka struct {
	ID          int    `gorm:"primaryKey;autoIncrement:true"`
	UID         int    `xml:"uid"`
	SdnEntryUID int    `xml:"-"`
	Type        string `xml:"type"`
	Category    string `xml:"category"`
	LastName    string `xml:"lastName"`
	FirstName   string `xml:"firstName"`
}

type Address struct {
	ID          int    `gorm:"primaryKey;autoIncrement:true"`
	UID         int    `xml:"uid"`
	SdnEntryUID int    `xml:"-"`
	City        string `xml:"city"`
	Country     string `xml:"country"`
}

type Person struct {
	UID       int    `json:"uid"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

const (
	StateEmpty    = 0
	StateUpdating = 1
	StateOk       = 2
)

func InitGorm(config DbConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.Dbname)

	//db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&SdnEntry{}, &Program{}, &Aka{}, &Address{}, &PublishInformation{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func WriteToDB(ctx context.Context, db *gorm.DB, entries <-chan SdnEntry) error {
	for entry := range entries {

		select {
		case <-ctx.Done():
			fmt.Println("Операция записи в базу данных была отменена")
			return ctx.Err()
		default:
		}

		var existingEntry SdnEntry
		// Попытка найти существующую запись по UID
		err := db.Where("uid = ?", entry.UID).First(&existingEntry).Error

		if err == gorm.ErrRecordNotFound {
			// Если запись не найдена, создаем новую
			if err := db.Create(&entry).Error; err != nil {
				fmt.Printf("Ошибка при создании записи: %v\n", err)
				return err
			}
		} else if err != nil {
			fmt.Printf("Ошибка при поиске существующей записи: %v\n", err)
			return err
		} else {
			// Если запись найдена, обновляем ее данными из entry
			entry.ID = existingEntry.ID                   // Убедитесь, что ID обновляемой записи установлен
			if err := db.Save(&entry).Error; err != nil { // Save обновит всю запись на основе ID
				fmt.Printf("Ошибка при обновлении записи: %v\n", err)
				return err
			}
		}
	}
	return nil
}

// Заглушка функции для запроса данных из базы. Необходимо реализовать логику запроса.
func GetPerson(db *gorm.DB, name string, searchType string) ([]*Person, error) {
	var persons []*Person

	nameParts := strings.Fields(name)
	query := db.Model(&SdnEntry{})

	if len(nameParts) == 1 {
		// Поиск по одному слову (имя или фамилия)
		namePart := nameParts[0]
		if searchType == "strong" {
			query = query.Where("first_name = ? OR last_name = ?", namePart, namePart)
		} else {
			query = query.Where("first_name ILIKE ? OR last_name ILIKE ?", "%"+namePart+"%", "%"+namePart+"%")
		}
	} else if len(nameParts) >= 2 {
		// Поиск по имени и фамилии
		firstName, lastName := nameParts[0], nameParts[len(nameParts)-1]
		if searchType == "strong" {
			query = query.Where("first_name = ? AND last_name = ?", firstName, lastName)
		} else {
			query = query.Where("first_name ILIKE ? OR last_name ILIKE ?", "%"+firstName+"%", "%"+lastName+"%")
		}
	}

	result := query.Find(&persons)

	if result.Error != nil {
		return nil, result.Error
	}

	return persons, nil
}

//func getLastPublishInfo() (time.Time, int) {
//	Запрос данных LastPublishInfo из базы
//	Должно возвращать:
//		<Publish_Date>02/08/2024</Publish_Date>
//		<Record_Count>13921</Record_Count>
//
//	return time.Time{}, 1
//}
