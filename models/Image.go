package models

import (
	"errors"
	"html"
	"log"
	"strings"
	"time"

	"github.com/badoux/checkmail"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

// Image object struct
type Image struct {
	ID        uint32    `gorm:"primary_key;auto_increment" json:"id"`
	Type	  string    `gorm:"size:255;not null;" json:"nickname"`
	Size	  string    `gorm:"size:255;not null;" json:"size"`
	Path	  string	`gorm:"size:5000;not null;" json:"path"`
	Owner     User      `json:"owner"`
	OwnerID   uint32     `gorm:"not null" json:"owner_id"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}


// Prepare Image object
func (i *Image) Prepare() {
	i.ID = 0
	i.Type = html.EscapeString(strings.TrimSpace(i.Type))
	i.Size = html.EscapeString(strings.TrimSpace(i.Size))
	i.Path = html.EscapeString(strings.TrimSpace(i.Path))
	i.Owner = User{}
	i.CreatedAt = time.Now()
	i.UpdatedAt = time.Now()
}

// Validate Image before save 
func (i *Image) Validate() error {

	if i.Type == "" {
		return errors.New("Required Type")
	}
	if i.OwnerID < 1 {
		return errors.New("Required Owner")
	}
	return nil
}

// SaveImage saving Image
func (i *Image) SaveImage(db *gorm.DB) (*Image, error) {

	var err error
	err = db.Debug().Model(&Image{}).Create(&i).Error

	if err != nil {
		return &Image{}, err
	}

	if i.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", i.OwnerID).Take(&i.Owner).Error

		if err != nil {
			return &Image{}, err
		}
	}

	return i, err
}

// FindAllImages finds all
func (i *Image) FindAllImages(db *gorm.DB) (*[]Image, error) {
	var err error
	images := []Image{}
	err = db.Debug().Model(&Image{}).Limit(100).Find(&images).Error
	if err != nil {
		return &[]Image{}, err
	}
	if len(images) > 0 {
		for i := range images {
			err := db.Debug().Model(&User{}).Where("id = ?", images[i].OwnerID).Take(&images[i].Owner).Error
			if err != nil {
				return &[]Image{}, err
			}
		}
	}

	return &images, nil
}

