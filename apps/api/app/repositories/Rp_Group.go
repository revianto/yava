package repositories

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"math"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/revianto/yava/api/app/models"
	"github.com/revianto/yava/api/exceptions"
	"github.com/revianto/yava/api/helpers"
	"gorm.io/gorm"
)

// =============================================================================
// HELPERS
// =============================================================================

func generateInviteCode() string {
	b := make([]byte, 8)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func groupToMap(g models.YvGroup) map[string]any {
	b, _ := json.Marshal(g)
	var m map[string]any
	json.Unmarshal(b, &m)
	return m
}

func groupMemberToMap(mem models.YvGroupMember) map[string]any {
	b, _ := json.Marshal(mem)
	var m map[string]any
	json.Unmarshal(b, &m)
	return m
}

func groupRecipeToMap(gr models.YvGroupRecipe) map[string]any {
	b, _ := json.Marshal(gr)
	var m map[string]any
	json.Unmarshal(b, &m)
	return m
}

// =============================================================================
// GROUP CRUD
// =============================================================================

func GroupSingle(tx *gorm.DB, data fiber.Map, c *fiber.Ctx, locale string, where func(*gorm.DB) *gorm.DB) (map[string]any, any) {
	var g models.YvGroup
	q := tx.Preload("Founder").Preload("Members.User")
	if where != nil {
		q = where(q)
	}
	if err := q.First(&g).Error; err != nil {
		return nil, exceptions.ErrorException(c, fiber.StatusNotFound, "Grup tidak ditemukan")
	}
	return groupToMap(g), nil
}

func GroupIndex(tx *gorm.DB, data fiber.Map, c *fiber.Ctx, locale string, userId int64) (models.IndexData, any) {
	var memberRows []models.YvGroupMember
	tx.Where("user_id = ?", userId).Find(&memberRows)

	groupIDs := make([]int64, 0, len(memberRows))
	for _, m := range memberRows {
		groupIDs = append(groupIDs, m.GroupId)
	}

	if len(groupIDs) == 0 {
		return models.IndexData{
			Data: []any{},
			Meta: models.Pagination{Page: 1, Limit: 20, Total: 0, TotalPages: 0},
		}, nil
	}

	var groups []models.YvGroup
	var total int64

	q := tx.Model(&models.YvGroup{}).Preload("Founder").Preload("Members.User").
		Where("id IN ?", groupIDs)

	q.Count(&total)

	page := helpers.Conv(data["page"]).Default(c.Query("page", "1")).Int()
	limit := helpers.Conv(data["limit"]).Default(c.Query("limit", "20")).Int()
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	if err := q.Order("yv_group.created_at DESC").Offset(offset).Limit(limit).Find(&groups).Error; err != nil {
		return models.IndexData{}, exceptions.ErrorException(c, fiber.StatusInternalServerError, "Gagal mengambil daftar grup")
	}

	items := make([]any, len(groups))
	for i, g := range groups {
		items[i] = groupToMap(g)
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	return models.IndexData{
		Data: items,
		Meta: models.Pagination{Page: page, Limit: limit, Total: int(total), TotalPages: totalPages},
	}, nil
}

func GroupCreate(tx *gorm.DB, data fiber.Map, c *fiber.Ctx, locale string) (map[string]any, any) {
	createdBy := helpers.Conv(data["created_by"]).Int64()
	name := helpers.Conv(data["name"]).String()

	g := models.YvGroup{
		Name:       name,
		InviteCode: generateInviteCode(),
		CreatedBy:  createdBy,
	}
	if v, ok := data["description"]; ok && v != nil {
		s := helpers.Conv(v).String()
		if s != "" {
			g.Description = &s
		}
	}
	if v, ok := data["avatar_url"]; ok && v != nil {
		s := helpers.Conv(v).String()
		if s != "" {
			g.AvatarUrl = &s
		}
	}

	if txErr := tx.Transaction(func(t *gorm.DB) error {
		if err := t.Create(&g).Error; err != nil {
			return err
		}
		founder := models.YvGroupMember{
			GroupId: g.Id,
			UserId:  createdBy,
			Role:    "founder",
		}
		return t.Create(&founder).Error
	}); txErr != nil {
		return nil, exceptions.ErrorException(c, fiber.StatusNotAcceptable, "Gagal membuat grup")
	}

	return GroupSingle(tx, data, c, locale, func(db *gorm.DB) *gorm.DB {
		return db.Where("yv_group.id = ?", g.Id)
	})
}

func GroupUpdate(tx *gorm.DB, data fiber.Map, c *fiber.Ctx, locale string, id int64) (map[string]any, any) {
	updates := map[string]any{}
	if v := helpers.Conv(data["name"]).String(); v != "" {
		updates["name"] = v
	}
	if v, ok := data["description"]; ok {
		updates["description"] = v
	}
	if v, ok := data["avatar_url"]; ok {
		updates["avatar_url"] = v
	}

	if len(updates) > 0 {
		if txErr := tx.Transaction(func(t *gorm.DB) error {
			return t.Model(&models.YvGroup{}).Where("id = ?", id).Updates(updates).Error
		}); txErr != nil {
			return nil, exceptions.ErrorException(c, fiber.StatusNotAcceptable, "Gagal memperbarui grup")
		}
	}

	return GroupSingle(tx, data, c, locale, func(db *gorm.DB) *gorm.DB {
		return db.Where("yv_group.id = ?", id)
	})
}

func GroupDelete(tx *gorm.DB, data fiber.Map, c *fiber.Ctx, locale string, id int64) (map[string]any, any) {
	if txErr := tx.Transaction(func(t *gorm.DB) error {
		return t.Delete(&models.YvGroup{}, id).Error
	}); txErr != nil {
		return nil, exceptions.ErrorException(c, fiber.StatusNotAcceptable, "Gagal menghapus grup")
	}
	return nil, nil
}

// =============================================================================
// GROUP MEMBERS
// =============================================================================

func GroupMemberIndex(tx *gorm.DB, data fiber.Map, c *fiber.Ctx, locale string, groupId int64) (models.IndexData, any) {
	var members []models.YvGroupMember
	var total int64

	q := tx.Model(&models.YvGroupMember{}).Preload("User").Where("group_id = ?", groupId)
	q.Count(&total)

	page := helpers.Conv(data["page"]).Default(c.Query("page", "1")).Int()
	limit := helpers.Conv(data["limit"]).Default(c.Query("limit", "20")).Int()
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	if err := q.Order("joined_at ASC").Offset(offset).Limit(limit).Find(&members).Error; err != nil {
		return models.IndexData{}, exceptions.ErrorException(c, fiber.StatusInternalServerError, "Gagal mengambil daftar anggota")
	}

	items := make([]any, len(members))
	for i, m := range members {
		items[i] = groupMemberToMap(m)
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	return models.IndexData{
		Data: items,
		Meta: models.Pagination{Page: page, Limit: limit, Total: int(total), TotalPages: totalPages},
	}, nil
}

func GroupMemberAdd(tx *gorm.DB, data fiber.Map, c *fiber.Ctx, locale string, groupId, userId int64, role string) (map[string]any, any) {
	m := models.YvGroupMember{
		GroupId: groupId,
		UserId:  userId,
		Role:    role,
	}
	if txErr := tx.Transaction(func(t *gorm.DB) error {
		return t.Create(&m).Error
	}); txErr != nil {
		return nil, exceptions.ErrorException(c, fiber.StatusNotAcceptable, "Gagal menambahkan anggota")
	}

	var result models.YvGroupMember
	if err := tx.Preload("User").Where("id = ?", m.Id).First(&result).Error; err != nil {
		return nil, exceptions.ErrorException(c, fiber.StatusNotFound, "Anggota tidak ditemukan")
	}
	return groupMemberToMap(result), nil
}

func GroupMemberRemove(tx *gorm.DB, data fiber.Map, c *fiber.Ctx, locale string, groupId, userId int64) (map[string]any, any) {
	if txErr := tx.Transaction(func(t *gorm.DB) error {
		return t.Where("group_id = ? AND user_id = ?", groupId, userId).Delete(&models.YvGroupMember{}).Error
	}); txErr != nil {
		return nil, exceptions.ErrorException(c, fiber.StatusNotAcceptable, "Gagal menghapus anggota")
	}
	return nil, nil
}

func GroupMemberSetRole(tx *gorm.DB, data fiber.Map, c *fiber.Ctx, locale string, groupId, userId int64, role string) (map[string]any, any) {
	if txErr := tx.Transaction(func(t *gorm.DB) error {
		return t.Model(&models.YvGroupMember{}).
			Where("group_id = ? AND user_id = ?", groupId, userId).
			Update("role", role).Error
	}); txErr != nil {
		return nil, exceptions.ErrorException(c, fiber.StatusNotAcceptable, "Gagal mengubah role anggota")
	}

	var result models.YvGroupMember
	if err := tx.Preload("User").Where("group_id = ? AND user_id = ?", groupId, userId).First(&result).Error; err != nil {
		return nil, exceptions.ErrorException(c, fiber.StatusNotFound, "Anggota tidak ditemukan")
	}
	return groupMemberToMap(result), nil
}

// =============================================================================
// GROUP RECIPES
// =============================================================================

func GroupRecipeIndex(tx *gorm.DB, data fiber.Map, c *fiber.Ctx, locale string, groupId int64, status string) (models.IndexData, any) {
	var groupRecipes []models.YvGroupRecipe
	var total int64

	q := tx.Model(&models.YvGroupRecipe{}).
		Preload("Recipe").Preload("Submitter").Preload("Reviewer").
		Where("group_id = ?", groupId)

	if status != "" {
		q = q.Where("status = ?", status)
	}

	q.Count(&total)

	page := helpers.Conv(data["page"]).Default(c.Query("page", "1")).Int()
	limit := helpers.Conv(data["limit"]).Default(c.Query("limit", "20")).Int()
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	if err := q.Order("yv_group_recipe.created_at DESC").Offset(offset).Limit(limit).Find(&groupRecipes).Error; err != nil {
		return models.IndexData{}, exceptions.ErrorException(c, fiber.StatusInternalServerError, "Gagal mengambil daftar resep grup")
	}

	items := make([]any, len(groupRecipes))
	for i, gr := range groupRecipes {
		items[i] = groupRecipeToMap(gr)
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	return models.IndexData{
		Data: items,
		Meta: models.Pagination{Page: page, Limit: limit, Total: int(total), TotalPages: totalPages},
	}, nil
}

func GroupRecipeAdd(tx *gorm.DB, data fiber.Map, c *fiber.Ctx, locale string, groupId, recipeId, submittedBy int64) (map[string]any, any) {
	gr := models.YvGroupRecipe{
		GroupId:     groupId,
		RecipeId:    recipeId,
		SubmittedBy: submittedBy,
		Status:      "pending",
	}
	if txErr := tx.Transaction(func(t *gorm.DB) error {
		return t.Create(&gr).Error
	}); txErr != nil {
		return nil, exceptions.ErrorException(c, fiber.StatusNotAcceptable, "Gagal mengirim resep ke grup")
	}

	var result models.YvGroupRecipe
	if err := tx.Preload("Recipe").Preload("Submitter").Where("id = ?", gr.Id).First(&result).Error; err != nil {
		return nil, exceptions.ErrorException(c, fiber.StatusNotFound, "Resep grup tidak ditemukan")
	}
	return groupRecipeToMap(result), nil
}

func GroupRecipeSetStatus(tx *gorm.DB, data fiber.Map, c *fiber.Ctx, locale string, groupId, recipeId, reviewerId int64, status string) (map[string]any, any) {
	now := time.Now()
	if txErr := tx.Transaction(func(t *gorm.DB) error {
		return t.Model(&models.YvGroupRecipe{}).
			Where("group_id = ? AND recipe_id = ?", groupId, recipeId).
			Updates(map[string]any{
				"status":      status,
				"reviewed_by": reviewerId,
				"reviewed_at": now,
			}).Error
	}); txErr != nil {
		return nil, exceptions.ErrorException(c, fiber.StatusNotAcceptable, "Gagal memperbarui status resep")
	}

	var result models.YvGroupRecipe
	if err := tx.Preload("Recipe").Preload("Submitter").Preload("Reviewer").
		Where("group_id = ? AND recipe_id = ?", groupId, recipeId).First(&result).Error; err != nil {
		return nil, exceptions.ErrorException(c, fiber.StatusNotFound, "Resep grup tidak ditemukan")
	}
	return groupRecipeToMap(result), nil
}

func GroupRecipeRemove(tx *gorm.DB, data fiber.Map, c *fiber.Ctx, locale string, groupId, recipeId int64) (map[string]any, any) {
	if txErr := tx.Transaction(func(t *gorm.DB) error {
		return t.Where("group_id = ? AND recipe_id = ?", groupId, recipeId).Delete(&models.YvGroupRecipe{}).Error
	}); txErr != nil {
		return nil, exceptions.ErrorException(c, fiber.StatusNotAcceptable, "Gagal menghapus resep dari grup")
	}
	return nil, nil
}
