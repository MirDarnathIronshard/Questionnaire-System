package routes

import (
	"github.com/QBG-P2/Voting-System/internal/handlers"
	"github.com/QBG-P2/Voting-System/internal/repositories"
	"github.com/QBG-P2/Voting-System/internal/services"
)

func accessControlRoute(rs *RouterService) {
	db := rs.db
	app := rs.app
	e := rs.enforcer

	userRepo := repositories.NewUserRepository(db)
	roleRepo := repositories.NewRoleRepository(db)
	permissionRepo := repositories.NewPermissionRepository(db)
	roleService := services.NewRoleService(userRepo, roleRepo, permissionRepo, e)
	userService := services.NewUserService(userRepo, roleRepo, e)
	permissionService := services.NewPermissionService(permissionRepo)
	roleHandler := handlers.NewRoleHandler(roleService)
	userHandler := handlers.NewUserHandler(userService)
	permissionHandler := handlers.NewPermissionHandler(permissionService)

	app.Get("/users/:id", userHandler.GetUser)
	app.Put("/users/:id", userHandler.UpdateUser)
	app.Delete("/users/:id", userHandler.DeleteUser)
	app.Post("/users/:id/roles", userHandler.AssignRole)
	app.Delete("/users/:id/roles", userHandler.RemoveRole)

	app.Post("/roles", roleHandler.CreateRole)
	app.Get("/roles/:id", roleHandler.GetRole)
	app.Put("/roles/:id", roleHandler.UpdateRole)
	app.Delete("/roles/:id", roleHandler.DeleteRole)
	app.Post("/roles/:id/permissions", roleHandler.AssignPermission)
	app.Delete("/roles/:id/permissions", roleHandler.RemovePermission)

	app.Post("/permissions", permissionHandler.CreatePermission)
	app.Get("/permissions/:id", permissionHandler.GetPermission)
	app.Put("/permissions/:id", permissionHandler.UpdatePermission)
	app.Delete("/permissions/:id", permissionHandler.DeletePermission)
}
