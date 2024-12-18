package services

import (
	"github.com/QBG-P2/Voting-System/internal/models"
	"github.com/QBG-P2/Voting-System/internal/repositories"
	"github.com/casbin/casbin/v2"
)

type RBACInitService struct {
	roleRepo            repositories.RoleRepository
	permissionRepo      repositories.PermissionRepository
	questPermissionRepo repositories.QuestionnairePermissionRepository
	enforcer            *casbin.Enforcer
}

func NewRBACInitService(gpr repositories.QuestionnairePermissionRepository, roleRepo repositories.RoleRepository, permissionRepo repositories.PermissionRepository, enforcer *casbin.Enforcer) *RBACInitService {
	return &RBACInitService{
		roleRepo:            roleRepo,
		permissionRepo:      permissionRepo,
		enforcer:            enforcer,
		questPermissionRepo: gpr,
	}
}

func (s *RBACInitService) InitializeRBAC() error {
	// Define all roles
	roles := []struct {
		name        string
		permissions []models.Permission
	}{
		{
			name: "super_admin",
			permissions: []models.Permission{
				{Name: "all_access", Path: "/*", Method: "*"},
			},
		},
		{
			name: "admin",
			permissions: []models.Permission{
				// Questionnaire Management
				{Name: "view_active_questionnaires", Path: "/api/questionnaire/active", Method: "GET"},
				{Name: "create_questionnaire", Path: "/api/questionnaire", Method: "POST"},
				{Name: "get_questionnaires", Path: "/api/questionnaire", Method: "GET"},
				{Name: "update_questionnaire", Path: "/api/questionnaire/*", Method: "PUT"},
				{Name: "delete_questionnaire", Path: "/api/questionnaire/*", Method: "DELETE"},
				{Name: "view_questionnaire", Path: "/api/questionnaire/*", Method: "GET"},
				{Name: "publish_questionnaire", Path: "/api/questionnaire/publish-questionnaire/*", Method: "GET"},
				{Name: "cancel_questionnaire", Path: "/api/questionnaire/canceled-questionnaire/*", Method: "GET"},
				{Name: "monitor_questionnaire", Path: "/api/questionnaire/monitoring/*", Method: "GET"},

				// Questions Management
				{Name: "create_question", Path: "/api/questions", Method: "POST"},
				{Name: "get_questions", Path: "/api/questions", Method: "GET"},
				{Name: "update_question", Path: "/api/questions/*", Method: "PUT"},
				{Name: "delete_question", Path: "/api/questions/*", Method: "DELETE"},
				{Name: "view_question", Path: "/api/questions/*", Method: "GET"},

				// Options Management
				{Name: "create_option", Path: "/api/options", Method: "POST"},
				{Name: "get_options", Path: "/api/options/question/*", Method: "GET"},
				{Name: "update_option", Path: "/api/options/*", Method: "PUT"},
				{Name: "delete_option", Path: "/api/options/*", Method: "DELETE"},
				{Name: "view_option", Path: "/api/options/*", Method: "GET"},

				// Response Management
				{Name: "view_responses", Path: "/api/responses/questionnaire", Method: "GET"},
				{Name: "manage_responses", Path: "/api/responses/*", Method: "*"},

				// Access Control
				{Name: "assign_role", Path: "/api/questionnaire-access/assign", Method: "POST"},
				{Name: "revoke_access", Path: "/api/questionnaire-access/revoke/*", Method: "DELETE"},
				{Name: "view_questionnaire_users", Path: "/api/questionnaire-access/*", Method: "GET"},

				// Notifications
				{Name: "manage_notifications", Path: "/api/notifications/*", Method: "*"},

				// Chats & Messages
				{Name: "manage_chats", Path: "/api/chat/*", Method: "*"},
				{Name: "manage_messages", Path: "/api/messages/*", Method: "*"},

				// Email Management
				{Name: "send_email", Path: "/api/emails/send", Method: "POST"},

				// Vote Transactions
				{Name: "manage_transactions", Path: "/api/vote-transactions/*", Method: "*"},
			},
		},
		{
			name: "user",
			permissions: []models.Permission{
				// Questionnaire Access
				{Name: "view_user_questionnaires", Path: "/api/questionnaire/user", Method: "GET"},
				{Name: "view_questionnaire_details", Path: "/api/questionnaire/*", Method: "GET"},
				{Name: "view_questionnaire_create", Path: "/api/questionnaire", Method: "POST"},
				{Name: "view_questionnaire_update", Path: "/api/questionnaire/*", Method: "PUT"},
				{Name: "view_questionnaire_delete", Path: "/api/questionnaire/*", Method: "DELETE"},

				// Response Management
				{Name: "submit_response", Path: "/api/responses", Method: "POST"},
				{Name: "update_response", Path: "/api/responses", Method: "PUT"},
				{Name: "view_own_responses", Path: "/api/responses/*", Method: "GET"},

				// Chat & Messages
				{Name: "use_chat", Path: "/api/chat", Method: "POST"},
				{Name: "send_message", Path: "/api/chat/message", Method: "POST"},
				{Name: "view_messages", Path: "/api/chat/*/messages", Method: "GET"},
				{Name: "view_user_chats", Path: "/api/chat/user", Method: "GET"},

				// Notifications
				{Name: "view_notifications", Path: "/api/notifications", Method: "GET"},
				{Name: "mark_notification_read", Path: "/api/notifications/*/read", Method: "PUT"},
				{Name: "get_unread_count", Path: "/api/notifications/unread-count", Method: "GET"},

				// Vote Transactions
				{Name: "create_transaction", Path: "/api/vote-transactions", Method: "POST"},
				{Name: "confirm_transaction", Path: "/api/vote-transactions/*/confirm", Method: "PUT"},

				// Questionnaire Role
				{Name: "view_questionnaire_roles", Path: "/api/questionnaire_role/questionnaire/*", Method: "GET"},
				{Name: "assign_role", Path: "/api/questionnaire-access/assign", Method: "POST"},
				{Name: "revoke_access", Path: "/api/questionnaire-access/revoke/*", Method: "DELETE"},
				{Name: "view_questionnaire_users", Path: "/api/questionnaire-access/*", Method: "GET"},
				{Name: "role_create", Path: "/api/questionnaire_role", Method: "*"},

				//option
				{Name: "option_get_by_id", Path: "/api/options/*", Method: "GET"},
				{Name: "option_create", Path: "/api/options", Method: "POST"},
				{Name: "option_update", Path: "/api/options/*", Method: "PUT"},
				{Name: "option_delete", Path: "/api/options/*", Method: "DELETE"},
				{Name: "option_question_by_id", Path: "/api/options/question/*", Method: "GET"},

				//question
				{Name: "question_create", Path: "/api/question", Method: "POST"},
				{Name: "question_update", Path: "/api/question/*", Method: "PUT"},
				{Name: "question_delete", Path: "/api/question/*", Method: "DELETE"},

				//profile
				{Name: "show_profile", Path: "/api/users/profile", Method: "*"},
				{Name: "wallet", Path: "/api/users/wallet", Method: "*"},

				//permission set to role
				{Name: "questionnaire-role-permissions", Path: "/api/questionnaire-role-permissions/*", Method: "*"},
				{Name: "questionnaire-role-permissions-type", Path: "/api/questionnaire/*/*", Method: "*"},

				{Name: "questionnaire-sequence-v1", Path: "/api/questionnaire-sequence/*", Method: "*"},
				{Name: "questionnaire-sequence-v2", Path: "/api/questionnaire-sequence/*/*", Method: "*"},
				{Name: "questionnaire-sequence-v3", Path: "/api/questionnaire-sequence/*/*/*", Method: "*"},
			},
		},
	}

	// Create roles and their permissions
	for _, roleData := range roles {
		// Create or get role
		role, err := s.roleRepo.GetByName(roleData.name)
		if err != nil {
			role = &models.Role{
				Name: roleData.name,
			}
			err = s.roleRepo.Create(role)
			if err != nil {
				return err
			}
		}

		// Create permissions and assign to role
		for _, perm := range roleData.permissions {
			// Create or get permission
			existingPerm, err := s.permissionRepo.GetByName(perm.Name)
			if err != nil {
				err = s.permissionRepo.Create(&perm)
				if err != nil {
					return err
				}
				existingPerm = &perm
			}

			// Assign permission to role in database
			err = s.roleRepo.AssignPermission(role.ID, existingPerm.ID)
			if err != nil {
				return err
			}

			// Add policy to Casbin enforcer
			_, err = s.enforcer.AddPolicy(role.Name, existingPerm.Path, existingPerm.Method)
			if err != nil {
				return err
			}
		}
	}

	return s.enforcer.SavePolicy()
}

func (s *RBACInitService) InitQuestionnairePermission() error {
	permissions := []models.QuestionnairePermission{
		{Name: "questionnaire-access-assign", Action: "assign", Resource: "access"},
		{Name: "questionnaire-access-revoke", Action: "revoke", Resource: "access"},
		{Name: "questionnaire-access-GetQuestionnaireUsers", Action: "view", Resource: "access"},

		{Name: "questionnaire-role-GetUserQuestionnaireRoles", Action: "view", Resource: "roles"},
		{Name: "questionnaire-role-AssignRole", Action: "assign", Resource: "roles"},
		{Name: "questionnaire-role-CreateQuestionnaireRole", Action: "create", Resource: "roles"},
		{Name: "questionnaire-role-GetQuestionnaireRole", Action: "view", Resource: "roles"},
		{Name: "questionnaire-role-UpdateQuestionnaireRole", Action: "update", Resource: "roles"},
		{Name: "questionnaire-role-DeleteQuestionnaireRole", Action: "delete", Resource: "roles"},

		{Name: "create-Question", Action: "create", Resource: "question"},
		{Name: "getAll-Question", Action: "view", Resource: "question"},
		{Name: "get-Question", Action: "view", Resource: "question"},
		{Name: "update-Question", Action: "update", Resource: "question"},
		{Name: "delete-Question", Action: "delete", Resource: "question"},

		{Name: "create-option", Action: "create", Resource: "option"},
		{Name: "update-option", Action: "update", Resource: "option"},
		{Name: "delete-option", Action: "delete", Resource: "option"},
		{Name: "get-option-with-question-id", Action: "view", Resource: "option"},

		{Name: "showe-questionnaire", Action: "view", Resource: "questionnaire"},
		{Name: "update-questionnaire", Action: "update", Resource: "questionnaire"},
		{Name: "delete-questionnaire", Action: "delete", Resource: "questionnaire"},
		{Name: "monitoring-questionnaire", Action: "monitor", Resource: "questionnaire"},
		{Name: "publish-questionnaire", Action: "publish", Resource: "questionnaire"},
		{Name: "canceled-questionnaire", Action: "cancel", Resource: "questionnaire"},

		{Name: "responses-GetResponsesByQuestionnaireID", Action: "view", Resource: "response"},
		{Name: "responses-GetResponseByID", Action: "view", Resource: "response"},
		{Name: "responses-UpdateResponse", Action: "update", Resource: "response"},
		{Name: "responses-DeleteResponse", Action: "delete", Resource: "response"},

		{Name: "assign-role-permissions", Action: "assign", Resource: "role-permissions"},
		{Name: "remove-role-permissions", Action: "remove", Resource: "role-permissions"},
		{Name: "view-role-permissions", Action: "view", Resource: "role-permissions"},
	}

	for _, permission := range permissions {
		if err := s.questPermissionRepo.Create(&permission); err != nil {
			return err
		}
	}
	return nil
}

func (s *RBACInitService) ValidatePermission(userEmail string, path string, method string) (bool, error) {
	return s.enforcer.Enforce(userEmail, path, method)
}

func (s *RBACInitService) AssignUserToRole(userEmail string, roleName string) error {
	_, err := s.enforcer.AddGroupingPolicy(userEmail, roleName)
	if err != nil {
		return err
	}
	return s.enforcer.SavePolicy()
}

func (s *RBACInitService) RemoveUserFromRole(userEmail string, roleName string) error {
	_, err := s.enforcer.RemoveGroupingPolicy(userEmail, roleName)
	if err != nil {
		return err
	}
	return s.enforcer.SavePolicy()
}

func (s *RBACInitService) GetRolePermissions(roleName string) ([]models.Permission, error) {
	role, err := s.roleRepo.GetByName(roleName)
	if err != nil {
		return nil, err
	}
	return role.Permissions, nil
}

func (s *RBACInitService) GetUserRoles(userEmail string) ([]string, error) {
	return s.enforcer.GetRolesForUser(userEmail)
}
