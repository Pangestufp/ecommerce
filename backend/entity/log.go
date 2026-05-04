package entity

import "time"


type Log struct {
	LogID         string `gorm:"primaryKey;"`
	ReferenceType string    
	ReferenceID   string    
	ReferenceName string
	SourceID      string    
    SourceName    string    
    SourceType    string    
	CreatedAt     time.Time 
	CreatedBy     string    
	Note          string    
	CreatedName   string    
}
	
	
