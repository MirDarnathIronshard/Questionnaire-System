package routes

import (
	"github.com/QBG-P2/Voting-System/internal/handlers"
	"github.com/QBG-P2/Voting-System/internal/repositories"
	"github.com/QBG-P2/Voting-System/internal/services"
)

func roleRoutes(rs *RouterService) {
	db := rs.db
	app := rs.app
	enforcer := rs.enforcer

	roleRepo := repositories.NewRoleRepository(db)
	userRepo := repositories.NewUserRepository(db)
	permissionRepo := repositories.NewPermissionRepository(db)

	roleService := services.NewRoleService(userRepo, roleRepo, permissionRepo, enforcer)

	roleHandler := handlers.NewRoleHandler(roleService)

	roles := app.Group("/api/roles")

	roles.Post("/", roleHandler.CreateRole)
	roles.Get("/", roleHandler.GetAllRoles)
	roles.Get("/:id", roleHandler.GetRole)
	roles.Put("/:id", roleHandler.UpdateRole)
	roles.Delete("/:id", roleHandler.DeleteRole)

	roles.Get("/:id/permissions", roleHandler.GetRolePermissions)
	roles.Post("/:id/permissions", roleHandler.AssignPermission)
	roles.Delete("/:id/permissions/:permissionId", roleHandler.RemovePermission)

	roles.Get("/users/:userId", roleHandler.GetUserRoles)
	roles.Post("/users/:userId", roleHandler.AssignRoleToUser)
	roles.Delete("/users/:userId/:roleId", roleHandler.RemoveRoleFromUser)

	// Role validation routes
	roles.Post("/validate", roleHandler.ValidateUserRole)
}
