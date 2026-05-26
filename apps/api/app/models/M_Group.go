package models

import "time"

type YvGroup struct {
	Id          int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string    `json:"name" gorm:"not null"`
	Description *string   `json:"description"`
	AvatarUrl   *string   `json:"avatar_url"`
	InviteCode  string    `json:"invite_code" gorm:"uniqueIndex;not null"`
	CreatedBy   int64     `json:"created_by" gorm:"not null"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	Founder *YvUser        `json:"founder,omitempty" gorm:"foreignKey:CreatedBy"`
	Members []YvGroupMember `json:"members,omitempty" gorm:"foreignKey:GroupId"`
}

func (YvGroup) TableName() string { return "yv_group" }

type YvGroupMember struct {
	Id       int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	GroupId  int64     `json:"group_id" gorm:"not null;index"`
	UserId   int64     `json:"user_id" gorm:"not null;index"`
	Role     string    `json:"role" gorm:"default:member;not null"`
	JoinedAt time.Time `json:"joined_at"`

	User *YvUser `json:"user,omitempty" gorm:"foreignKey:UserId"`
}

func (YvGroupMember) TableName() string { return "yv_group_member" }

type YvGroupRecipe struct {
	Id          int64      `json:"id" gorm:"primaryKey;autoIncrement"`
	GroupId     int64      `json:"group_id" gorm:"not null;index"`
	RecipeId    int64      `json:"recipe_id" gorm:"not null;index"`
	SubmittedBy int64      `json:"submitted_by" gorm:"not null"`
	Status      string     `json:"status" gorm:"default:pending;not null"`
	ReviewedBy  *int64     `json:"reviewed_by"`
	ReviewedAt  *time.Time `json:"reviewed_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`

	Recipe    *Recipe `json:"recipe,omitempty" gorm:"foreignKey:RecipeId"`
	Submitter *YvUser `json:"submitter,omitempty" gorm:"foreignKey:SubmittedBy"`
	Reviewer  *YvUser `json:"reviewer,omitempty" gorm:"foreignKey:ReviewedBy"`
}

func (YvGroupRecipe) TableName() string { return "yv_group_recipe" }
