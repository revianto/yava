package services

import (
	"github.com/gofiber/fiber/v2"
	"github.com/revianto/yava/api/app/middlewares"
	"github.com/revianto/yava/api/app/models"
	"github.com/revianto/yava/api/app/repositories"
	"github.com/revianto/yava/api/exceptions"
	"github.com/revianto/yava/api/helpers"
	"gorm.io/gorm"
)

// =============================================================================
// HELPERS
// =============================================================================

func getGroupMemberRole(tx *gorm.DB, groupId, userId int64) (string, bool) {
	var m models.YvGroupMember
	err := tx.Where("group_id = ? AND user_id = ?", groupId, userId).First(&m).Error
	if err != nil {
		return "", false
	}
	return m.Role, true
}

func isAdminOrFounder(role string) bool {
	return role == "founder" || role == "admin"
}

// =============================================================================
// GROUP CRUD
// =============================================================================

func GroupList(tx *gorm.DB, data fiber.Map, c *fiber.Ctx, locale string) (any, any) {
	userId := middlewares.YvUserID(c)
	return repositories.GroupIndex(tx, data, c, locale, userId)
}

func GroupShow(tx *gorm.DB, data fiber.Map, c *fiber.Ctx, locale string, id any) (any, any) {
	gid := helpers.Conv(id).Int64()
	return repositories.GroupSingle(tx, data, c, locale, func(db *gorm.DB) *gorm.DB {
		return db.Where("yv_group.id = ?", gid)
	})
}

type ValidateGroupCreate struct {
	Name string `json:"name" validate:"required"`
}

func GroupCreate(tx *gorm.DB, data fiber.Map, c *fiber.Ctx, locale string) (any, any) {
	if _, err := models.Validate(c, data, new(ValidateGroupCreate), locale); err != nil {
		return nil, err
	}
	data["created_by"] = middlewares.YvUserID(c)
	return repositories.GroupCreate(tx, data, c, locale)
}

func GroupUpdate(tx *gorm.DB, data fiber.Map, c *fiber.Ctx, locale string, id any) (any, any) {
	gid := helpers.Conv(id).Int64()
	userId := middlewares.YvUserID(c)

	role, isMember := getGroupMemberRole(tx, gid, userId)
	if !isMember || !isAdminOrFounder(role) {
		return nil, exceptions.ErrorException(c, fiber.StatusForbidden, "Hanya admin atau founder yang dapat memperbarui grup")
	}

	return repositories.GroupUpdate(tx, data, c, locale, gid)
}

func GroupDelete(tx *gorm.DB, data fiber.Map, c *fiber.Ctx, locale string, id any) (any, any) {
	gid := helpers.Conv(id).Int64()
	userId := middlewares.YvUserID(c)

	role, isMember := getGroupMemberRole(tx, gid, userId)
	if !isMember || role != "founder" {
		return nil, exceptions.ErrorException(c, fiber.StatusForbidden, "Hanya founder yang dapat menghapus grup")
	}

	return repositories.GroupDelete(tx, data, c, locale, gid)
}

// =============================================================================
// GROUP MEMBERS
// =============================================================================

func GroupMemberList(tx *gorm.DB, data fiber.Map, c *fiber.Ctx, locale string, id any) (any, any) {
	gid := helpers.Conv(id).Int64()
	return repositories.GroupMemberIndex(tx, data, c, locale, gid)
}

func GroupJoin(tx *gorm.DB, data fiber.Map, c *fiber.Ctx, locale string, id any) (any, any) {
	gid := helpers.Conv(id).Int64()
	userId := middlewares.YvUserID(c)
	inviteCode := helpers.Conv(data["invite_code"]).String()

	if inviteCode == "" {
		return nil, exceptions.ErrorException(c, fiber.StatusUnprocessableEntity, "invite_code diperlukan")
	}

	// Verify invite code matches group
	var g models.YvGroup
	if err := tx.Where("id = ? AND invite_code = ?", gid, inviteCode).First(&g).Error; err != nil {
		return nil, exceptions.ErrorException(c, fiber.StatusNotFound, "Kode undangan tidak valid atau grup tidak ditemukan")
	}

	// Check not already a member
	if _, already := getGroupMemberRole(tx, gid, userId); already {
		return nil, exceptions.ErrorException(c, fiber.StatusConflict, "Anda sudah menjadi anggota grup ini")
	}

	return repositories.GroupMemberAdd(tx, data, c, locale, gid, userId, "member")
}

func GroupMemberRemove(tx *gorm.DB, data fiber.Map, c *fiber.Ctx, locale string, groupId, userId any) (any, any) {
	gid := helpers.Conv(groupId).Int64()
	uid := helpers.Conv(userId).Int64()
	actorId := middlewares.YvUserID(c)

	actorRole, isMember := getGroupMemberRole(tx, gid, actorId)
	if !isMember || !isAdminOrFounder(actorRole) {
		return nil, exceptions.ErrorException(c, fiber.StatusForbidden, "Hanya admin atau founder yang dapat menghapus anggota")
	}

	targetRole, targetExists := getGroupMemberRole(tx, gid, uid)
	if !targetExists {
		return nil, exceptions.ErrorException(c, fiber.StatusNotFound, "Anggota tidak ditemukan")
	}
	if targetRole == "founder" {
		return nil, exceptions.ErrorException(c, fiber.StatusForbidden, "Founder tidak dapat dihapus dari grup")
	}

	return repositories.GroupMemberRemove(tx, data, c, locale, gid, uid)
}

func GroupMemberSetRole(tx *gorm.DB, data fiber.Map, c *fiber.Ctx, locale string, groupId, userId any) (any, any) {
	gid := helpers.Conv(groupId).Int64()
	uid := helpers.Conv(userId).Int64()
	actorId := middlewares.YvUserID(c)
	newRole := helpers.Conv(data["role"]).String()

	validRoles := map[string]bool{"admin": true, "member": true}
	if !validRoles[newRole] {
		return nil, exceptions.ErrorException(c, fiber.StatusUnprocessableEntity, "Role tidak valid. Gunakan 'admin' atau 'member'")
	}

	actorRole, isMember := getGroupMemberRole(tx, gid, actorId)
	if !isMember || !isAdminOrFounder(actorRole) {
		return nil, exceptions.ErrorException(c, fiber.StatusForbidden, "Hanya admin atau founder yang dapat mengubah role anggota")
	}

	targetRole, targetExists := getGroupMemberRole(tx, gid, uid)
	if !targetExists {
		return nil, exceptions.ErrorException(c, fiber.StatusNotFound, "Anggota tidak ditemukan")
	}
	if targetRole == "founder" {
		return nil, exceptions.ErrorException(c, fiber.StatusForbidden, "Role founder tidak dapat diubah")
	}

	return repositories.GroupMemberSetRole(tx, data, c, locale, gid, uid, newRole)
}

// =============================================================================
// GROUP RECIPES
// =============================================================================

func GroupRecipeList(tx *gorm.DB, data fiber.Map, c *fiber.Ctx, locale string, id any) (any, any) {
	gid := helpers.Conv(id).Int64()
	return repositories.GroupRecipeIndex(tx, data, c, locale, gid, "approved")
}

func GroupRecipePending(tx *gorm.DB, data fiber.Map, c *fiber.Ctx, locale string, id any) (any, any) {
	gid := helpers.Conv(id).Int64()
	userId := middlewares.YvUserID(c)

	role, isMember := getGroupMemberRole(tx, gid, userId)
	if !isMember || !isAdminOrFounder(role) {
		return nil, exceptions.ErrorException(c, fiber.StatusForbidden, "Hanya admin atau founder yang dapat melihat resep pending")
	}

	return repositories.GroupRecipeIndex(tx, data, c, locale, gid, "pending")
}

func GroupRecipeSubmit(tx *gorm.DB, data fiber.Map, c *fiber.Ctx, locale string, id any) (any, any) {
	gid := helpers.Conv(id).Int64()
	userId := middlewares.YvUserID(c)
	recipeId := helpers.Conv(data["recipe_id"]).Int64()

	if recipeId == 0 {
		return nil, exceptions.ErrorException(c, fiber.StatusUnprocessableEntity, "recipe_id diperlukan")
	}

	_, isMember := getGroupMemberRole(tx, gid, userId)
	if !isMember {
		return nil, exceptions.ErrorException(c, fiber.StatusForbidden, "Anda harus menjadi anggota grup untuk mengirim resep")
	}

	return repositories.GroupRecipeAdd(tx, data, c, locale, gid, recipeId, userId)
}

func GroupRecipeApprove(tx *gorm.DB, data fiber.Map, c *fiber.Ctx, locale string, groupId, recipeId any) (any, any) {
	gid := helpers.Conv(groupId).Int64()
	rid := helpers.Conv(recipeId).Int64()
	userId := middlewares.YvUserID(c)

	role, isMember := getGroupMemberRole(tx, gid, userId)
	if !isMember || !isAdminOrFounder(role) {
		return nil, exceptions.ErrorException(c, fiber.StatusForbidden, "Hanya admin atau founder yang dapat menyetujui resep")
	}

	return repositories.GroupRecipeSetStatus(tx, data, c, locale, gid, rid, userId, "approved")
}

func GroupRecipeReject(tx *gorm.DB, data fiber.Map, c *fiber.Ctx, locale string, groupId, recipeId any) (any, any) {
	gid := helpers.Conv(groupId).Int64()
	rid := helpers.Conv(recipeId).Int64()
	userId := middlewares.YvUserID(c)

	role, isMember := getGroupMemberRole(tx, gid, userId)
	if !isMember || !isAdminOrFounder(role) {
		return nil, exceptions.ErrorException(c, fiber.StatusForbidden, "Hanya admin atau founder yang dapat menolak resep")
	}

	return repositories.GroupRecipeSetStatus(tx, data, c, locale, gid, rid, userId, "rejected")
}

func GroupRecipeRemove(tx *gorm.DB, data fiber.Map, c *fiber.Ctx, locale string, groupId, recipeId any) (any, any) {
	gid := helpers.Conv(groupId).Int64()
	rid := helpers.Conv(recipeId).Int64()
	userId := middlewares.YvUserID(c)

	role, isMember := getGroupMemberRole(tx, gid, userId)
	if !isMember || !isAdminOrFounder(role) {
		return nil, exceptions.ErrorException(c, fiber.StatusForbidden, "Hanya admin atau founder yang dapat menghapus resep dari grup")
	}

	return repositories.GroupRecipeRemove(tx, data, c, locale, gid, rid)
}
