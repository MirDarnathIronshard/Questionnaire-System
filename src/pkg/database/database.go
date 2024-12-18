package database

import (
	"fmt"
	"github.com/QBG-P2/Voting-System/internal/repositories"
	"github.com/QBG-P2/Voting-System/internal/services"
	"github.com/QBG-P2/Voting-System/pkg/security"
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"log"
	"strings"
	"time"

	"github.com/QBG-P2/Voting-System/config"
	"github.com/QBG-P2/Voting-System/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var dbClient *gorm.DB

func InitDb(cfg *config.Config) error {
	var err error
	cnn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=Asia/Tehran",
		cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.User, cfg.Postgres.Password,
		cfg.Postgres.DbName, cfg.Postgres.SSLMode)

	dbClient, err = gorm.Open(postgres.Open(cnn), &gorm.Config{})
	if err != nil {
		return err
	}

	sqlDb, _ := dbClient.DB()
	err = sqlDb.Ping()
	if err != nil {
		return err
	}

	sqlDb.SetMaxIdleConns(cfg.Postgres.MaxIdleConns)
	sqlDb.SetMaxOpenConns(cfg.Postgres.MaxOpenConns)
	sqlDb.SetConnMaxLifetime(cfg.Postgres.ConnMaxLifetime * time.Minute)

	log.Println("Db connection established")
	return nil
}

func GetDb() *gorm.DB {
	return dbClient
}

func CloseDb() {
	con, _ := dbClient.DB()
	con.Close()
}

func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(&models.User{}, &models.Questionnaire{}, &models.Chat{}, &models.Message{}, &models.Question{},
		&models.QuestionnairePermission{}, &models.UserQuestionnaireAccess{},
		&models.QuestionnaireRole{}, &models.QuestionnaireRolePermission{},
		&models.Option{}, &models.VoteTransaction{}, &models.Response{}, &models.Notification{}, &models.Role{}, &models.Permission{})
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
	log.Println("Migration completed successfully")
}

func SetupCasbin(db *gorm.DB) *casbin.Enforcer {
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		log.Fatalf("Failed to create Casbin adapter: %v", err)
	}

	enforcer, err := casbin.NewEnforcer("config/rbac_model.conf", adapter)
	if err != nil {
		log.Fatalf("Failed to create Casbin enforcer: %v", err)
	}

	err = enforcer.LoadPolicy()
	if err != nil {
		log.Fatalf("Failed to load Casbin policy: %v", err)
	}

	return enforcer
}

func SuperAdminInit(enforcer *casbin.Enforcer, db *gorm.DB, cfg *config.Config) {
	hasPolicy, err := enforcer.HasPolicy("super_admin", "*", "*")
	if err != nil {
		log.Fatalf("Failed to check if super admin policy exists: %v", err)
	}
	if !hasPolicy {
		_, err = enforcer.AddPolicy("super_admin", "*", "*")
		if err != nil {
			log.Fatalf("Failed to add super admin policy: %v", err)
		}

		err = enforcer.SavePolicy()
		if err != nil {
			log.Fatalf("Failed to save Casbin policy: %v", err)
		}
	}

	userRepo := repositories.NewUserRepository(db)
	roleRepo := repositories.NewRoleRepository(db)
	permissionRepo := repositories.NewPermissionRepository(db)
	userService := services.NewUserService(userRepo, roleRepo, enforcer)
	//authService := services.NewAuthService(userRepo, roleRepo, enforcer, cfg.JWT.Secret)
	roleService := services.NewRoleService(userRepo, roleRepo, permissionRepo, enforcer)

	superAdmin, err := userService.GetByEmail("super_admin@exampel.com")
	if err != nil {

		password, _ := security.HashPassword("admin")

		superAdmin = &models.User{
			Email:    "super_admin@exampel.com",
			Password: password,
		}
		err = userService.Create(superAdmin)
		if err != nil {
			log.Fatalln(err)
			return
		}

		if err != nil {
			log.Fatalf("Failed to create super admin: %v", err)
		}

		role, err := roleService.GetByName("super_admin")
		if err != nil {
			role = &models.Role{
				Name: "super_admin",
			}
			err = roleService.Create(role)
			if err != nil {
				log.Fatalf("Failed to create super admin role: %v", err)
			}
		}

		err = userService.AssignRole(superAdmin.ID, role.ID)
		if err != nil {
			log.Fatalf("Failed to assign super admin role to user: %v", err)
		}
	}

	go func() {
		time.Sleep(time.Second)
		printSuperAdminInfo(&models.User{
			Email:    "super_admin@exampel.com",
			Password: "admin",
			Role:     "super_admin",
		})

	}()
}

func printSuperAdminInfo(user *models.User) {
	separator := strings.Repeat("-", 60)
	fmt.Println(separator)
	fmt.Printf("| %-20s | %-30s |\n", "Field", "Value")
	fmt.Println(separator)
	fmt.Printf("| %-20s | %-30s |\n", "Email - Username", user.Email)
	fmt.Printf("| %-20s | %-30s |\n", "Role", user.Role)
	fmt.Printf("| %-20s | %-30s |\n", "Password", user.Password)
	fmt.Println(separator)
	fmt.Println("user has been successfully initialized.")
}
