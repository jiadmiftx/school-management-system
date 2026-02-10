package routes

import (
	"fmt"
	"log"

	"sekolah-madrasah/app/controller/activity_controller"
	"sekolah-madrasah/app/controller/auth_controller"
	"sekolah-madrasah/app/controller/class_controller"
	"sekolah-madrasah/app/controller/class_enrollment_controller"
	"sekolah-madrasah/app/controller/organization_controller"
	"sekolah-madrasah/app/controller/permission_controller"
	"sekolah-madrasah/app/controller/post_controller"
	"sekolah-madrasah/app/controller/role_controller"
	"sekolah-madrasah/app/controller/student_profile_controller"
	"sekolah-madrasah/app/controller/subject_controller"
	"sekolah-madrasah/app/controller/teacher_profile_controller"
	"sekolah-madrasah/app/controller/unit_controller"
	"sekolah-madrasah/app/controller/unit_member_controller"
	"sekolah-madrasah/app/controller/unit_settings_controller"
	"sekolah-madrasah/app/controller/user_controller"
	"sekolah-madrasah/app/repository/activity_repository"
	"sekolah-madrasah/app/repository/class_enrollment_repository"
	"sekolah-madrasah/app/repository/class_repository"
	"sekolah-madrasah/app/repository/org_member_repository"
	"sekolah-madrasah/app/repository/organization_repository"
	"sekolah-madrasah/app/repository/permission_repository"
	"sekolah-madrasah/app/repository/post_repository"
	"sekolah-madrasah/app/repository/role_repository"
	"sekolah-madrasah/app/repository/student_profile_repository"
	"sekolah-madrasah/app/repository/subject_repository"
	"sekolah-madrasah/app/repository/teacher_profile_repository"
	"sekolah-madrasah/app/repository/unit_member_repository"
	"sekolah-madrasah/app/repository/unit_repository"
	"sekolah-madrasah/app/repository/user_repository"
	"sekolah-madrasah/app/service/membership_service"
	"sekolah-madrasah/app/use_case/activity_use_case"
	"sekolah-madrasah/app/use_case/auth_use_case"
	"sekolah-madrasah/app/use_case/class_enrollment_use_case"
	"sekolah-madrasah/app/use_case/class_use_case"
	"sekolah-madrasah/app/use_case/organization_use_case"
	"sekolah-madrasah/app/use_case/permission_use_case"
	"sekolah-madrasah/app/use_case/post_use_case"
	"sekolah-madrasah/app/use_case/role_use_case"
	"sekolah-madrasah/app/use_case/student_profile_use_case"
	"sekolah-madrasah/app/use_case/subject_use_case"
	"sekolah-madrasah/app/use_case/teacher_profile_use_case"
	"sekolah-madrasah/app/use_case/unit_member_use_case"
	"sekolah-madrasah/app/use_case/unit_use_case"
	"sekolah-madrasah/app/use_case/user_use_case"
	"sekolah-madrasah/config"
	"sekolah-madrasah/database"
	_ "sekolah-madrasah/docs" // swagger docs
	"sekolah-madrasah/pkg/http_middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

type Container struct {
	AuthController            auth_controller.AuthController
	UserController            user_controller.UserController
	RoleController            role_controller.RoleController
	PermissionController      permission_controller.PermissionController
	OrganizationController    organization_controller.OrganizationController
	UnitController            unit_controller.GinUnitController
	UnitMemberController      unit_member_controller.GinUnitMemberController
	UnitSettingsController    unit_settings_controller.GinUnitSettingsController
	PostController            *post_controller.PostController
	TeacherProfileController  *teacher_profile_controller.TeacherProfileController
	StudentProfileController  *student_profile_controller.StudentProfileController
	ClassController           *class_controller.ClassController
	ClassEnrollmentController *class_enrollment_controller.ClassEnrollmentController
	SubjectController         *subject_controller.SubjectController
	ActivityController        *activity_controller.ActivityController
}

func NewContainer(db *gorm.DB) *Container {
	userRepo := user_repository.NewUserRepository(db)
	roleRepo := role_repository.NewRoleRepository(db)
	permissionRepo := permission_repository.NewPermissionRepository(db)
	orgRepo := organization_repository.NewOrganizationRepository(db)
	orgMemberRepo := org_member_repository.NewOrgMemberRepository(db)
	unitRepo := unit_repository.NewUnitRepository(db)
	unitMemberRepo := unit_member_repository.NewUnitMemberRepository(db)
	postRepo := post_repository.NewPostRepository(db)
	teacherProfileRepo := teacher_profile_repository.NewTeacherProfileRepository(db)
	studentProfileRepo := student_profile_repository.NewStudentProfileRepository(db)
	classRepo := class_repository.NewClassRepository(db)
	classEnrollmentRepo := class_enrollment_repository.NewClassEnrollmentRepository(db)
	subjectRepo := subject_repository.NewSubjectRepository(db)
	activityRepo := activity_repository.NewActivityRepository(db)

	authUseCase := auth_use_case.NewAuthUseCase(userRepo)
	userUseCase := user_use_case.NewUserUseCase(userRepo)
	roleUseCase := role_use_case.NewRoleUseCase(roleRepo, permissionRepo)
	permissionUseCase := permission_use_case.NewPermissionUseCase(permissionRepo)
	orgUseCase := organization_use_case.NewOrganizationUseCase(orgRepo, orgMemberRepo)
	unitUseCase := unit_use_case.NewUnitUseCase(unitRepo)
	unitMemberUseCase := unit_member_use_case.NewUnitMemberUseCase(unitMemberRepo)
	postUseCase := post_use_case.NewPostUseCase(postRepo, userRepo)
	teacherProfileUseCase := teacher_profile_use_case.NewTeacherProfileUseCase(teacherProfileRepo)
	studentProfileUseCase := student_profile_use_case.NewStudentProfileUseCase(studentProfileRepo)
	classUseCase := class_use_case.NewClassUseCase(classRepo)
	classEnrollmentUseCase := class_enrollment_use_case.NewClassEnrollmentUseCase(classEnrollmentRepo)
	subjectUseCase := subject_use_case.NewSubjectUseCase(subjectRepo)
	activityUseCase := activity_use_case.NewActivityUseCase(activityRepo)
	membershipService := membership_service.NewMembershipService(db)

	authController := auth_controller.NewAuthController(authUseCase)
	userController := user_controller.NewUserController(userUseCase, membershipService)
	roleController := role_controller.NewRoleController(roleUseCase)
	permissionController := permission_controller.NewPermissionController(permissionUseCase)
	orgController := organization_controller.NewOrganizationController(orgUseCase)
	unitController := unit_controller.NewUnitController(unitUseCase)
	unitMemberController := unit_member_controller.NewUnitMemberController(unitMemberUseCase)
	unitSettingsController := unit_settings_controller.NewUnitSettingsController(db)
	postCtrl := post_controller.NewPostController(postUseCase)
	teacherProfileCtrl := teacher_profile_controller.NewTeacherProfileController(teacherProfileUseCase, userUseCase)
	studentProfileCtrl := student_profile_controller.NewStudentProfileController(studentProfileUseCase, userUseCase)
	classCtrl := class_controller.NewClassController(classUseCase)
	classEnrollmentCtrl := class_enrollment_controller.NewClassEnrollmentController(classEnrollmentUseCase)
	subjectCtrl := subject_controller.NewSubjectController(subjectUseCase)
	activityCtrl := activity_controller.NewActivityController(activityUseCase)

	return &Container{
		AuthController:            authController,
		UserController:            userController,
		RoleController:            roleController,
		PermissionController:      permissionController,
		OrganizationController:    orgController,
		UnitController:            unitController,
		UnitMemberController:      unitMemberController,
		UnitSettingsController:    unitSettingsController,
		PostController:            postCtrl,
		TeacherProfileController:  teacherProfileCtrl,
		StudentProfileController:  studentProfileCtrl,
		ClassController:           classCtrl,
		ClassEnrollmentController: classEnrollmentCtrl,
		SubjectController:         subjectCtrl,
		ActivityController:        activityCtrl,
	}
}

func InitRoutes() *gin.Engine {
	mainDB, err := database.NewDBConnection(config.APP.MainDB)
	if err != nil {
		log.Fatalf("❌ Failed to connect to main database: %v", err)
	}
	log.Println("✅ Main database connected")

	router := gin.Default()

	router.Use(http_middleware.CORS)

	container := NewContainer(mainDB)

	// Swagger docs
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := router.Group("/api/v1")
	{
		v1.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "pong"})
		})

		auth := v1.Group("/auth")
		{
			auth.POST("/login", container.AuthController.Login)
			auth.POST("/register", container.AuthController.Register)
			auth.POST("/refresh", container.AuthController.RefreshToken)
		}

		users := v1.Group("/users")
		users.Use(http_middleware.JWTAuthentication)
		{
			users.GET("", container.UserController.GetUsers)
			users.GET("/me", container.UserController.GetCurrentUser)
			users.GET("/me/memberships", container.UserController.GetMyMemberships)
			users.GET("/:id", container.UserController.GetUser)
			users.POST("", container.UserController.CreateUser)
			users.PUT("/:id", container.UserController.UpdateUser)
			users.DELETE("/:id", container.UserController.DeleteUser)
		}

		roles := v1.Group("/roles")
		roles.Use(http_middleware.JWTAuthentication)
		{
			roles.GET("", container.RoleController.GetRoles)
			roles.GET("/:id", container.RoleController.GetRole)
			roles.POST("", container.RoleController.CreateRole)
			roles.PUT("/:id", container.RoleController.UpdateRole)
			roles.DELETE("/:id", container.RoleController.DeleteRole)
		}

		permissions := v1.Group("/permissions")
		permissions.Use(http_middleware.JWTAuthentication)
		{
			permissions.GET("", container.PermissionController.GetPermissions)
			permissions.GET("/:id", container.PermissionController.GetPermission)
			permissions.POST("", container.PermissionController.CreatePermission)
			permissions.DELETE("/:id", container.PermissionController.DeletePermission)
		}

		organizations := v1.Group("/organizations")
		organizations.Use(http_middleware.JWTAuthentication)
		{
			organizations.GET("", container.OrganizationController.GetOrganizations)
			organizations.GET("/:id", container.OrganizationController.GetOrganization)
			organizations.POST("", container.OrganizationController.CreateOrganization)
			organizations.PUT("/:id", container.OrganizationController.UpdateOrganization)
			organizations.DELETE("/:id", container.OrganizationController.DeleteOrganization)

			organizations.GET("/:id/members", container.OrganizationController.GetMembers)
			organizations.POST("/:id/members", container.OrganizationController.AddMember)
			organizations.PUT("/:id/members/:userId", container.OrganizationController.UpdateMember)
			organizations.DELETE("/:id/members/:userId", container.OrganizationController.RemoveMember)
		}

		units := v1.Group("/units")
		units.Use(http_middleware.JWTAuthentication)
		{
			units.GET("", container.UnitController.GetUnits)
			units.GET("/:id", container.UnitController.GetUnit)
			units.POST("", container.UnitController.CreateUnit)
			units.PUT("/:id", container.UnitController.UpdateUnit)
			units.DELETE("/:id", container.UnitController.DeleteUnit)

			units.GET("/:id/members", container.UnitMemberController.GetMembers)
			units.GET("/:id/members/:memberId", container.UnitMemberController.GetMember)
			units.POST("/:id/members", container.UnitMemberController.AddMember)
			units.PUT("/:id/members/:memberId", container.UnitMemberController.UpdateMember)
			units.DELETE("/:id/members/:memberId", container.UnitMemberController.RemoveMember)

			units.GET("/:id/settings", container.UnitSettingsController.GetSettings)
			units.PUT("/:id/settings", container.UnitSettingsController.UpdateSettings)

			// Teacher profiles
			units.GET("/:id/teachers", container.TeacherProfileController.GetAll)
			units.GET("/:id/teachers/:teacherId", container.TeacherProfileController.GetById)
			units.POST("/:id/teachers", container.TeacherProfileController.Create)
			units.POST("/:id/teachers/with-user", container.TeacherProfileController.CreateWithUser)
			units.PUT("/:id/teachers/:teacherId", container.TeacherProfileController.Update)
			units.DELETE("/:id/teachers/:teacherId", container.TeacherProfileController.Delete)

			// Student profiles
			units.GET("/:id/students", container.StudentProfileController.GetAll)
			units.GET("/:id/students/:studentId", container.StudentProfileController.GetById)
			units.POST("/:id/students", container.StudentProfileController.Create)
			units.POST("/:id/students/with-user", container.StudentProfileController.CreateWithUser)
			units.PUT("/:id/students/:studentId", container.StudentProfileController.Update)
			units.DELETE("/:id/students/:studentId", container.StudentProfileController.Delete)

			// Classes
			units.GET("/:id/classes", container.ClassController.GetAll)
			units.GET("/:id/classes/:classId", container.ClassController.GetById)
			units.POST("/:id/classes", container.ClassController.Create)
			units.PUT("/:id/classes/:classId", container.ClassController.Update)
			units.DELETE("/:id/classes/:classId", container.ClassController.Delete)

			// Class enrollments
			units.GET("/:id/classes/:classId/students", container.ClassEnrollmentController.GetByClass)
			units.POST("/:id/classes/:classId/enroll", container.ClassEnrollmentController.Enroll)

			// Subjects
			units.GET("/:id/subjects", container.SubjectController.GetAll)
			units.GET("/:id/subjects/:subjectId", container.SubjectController.GetById)
			units.POST("/:id/subjects", container.SubjectController.Create)
			units.PUT("/:id/subjects/:subjectId", container.SubjectController.Update)
			units.DELETE("/:id/subjects/:subjectId", container.SubjectController.Delete)

			// Activities
			units.GET("/:id/activities", container.ActivityController.GetAll)
			units.POST("/:id/activities", container.ActivityController.Create)
		}

		// Class enrollment management (outside unit scope)
		classEnrollments := v1.Group("/class-enrollments")
		classEnrollments.Use(http_middleware.JWTAuthentication)
		{
			classEnrollments.PUT("/:enrollmentId", container.ClassEnrollmentController.UpdateStatus)
			classEnrollments.POST("/:enrollmentId/transfer", container.ClassEnrollmentController.Transfer)
			classEnrollments.DELETE("/:enrollmentId", container.ClassEnrollmentController.Remove)
		}

		posts := v1.Group("/posts")
		posts.Use(http_middleware.JWTAuthentication)
		{
			posts.GET("", container.PostController.GetPosts)
			posts.GET("/:id", container.PostController.GetPost)
			posts.POST("", container.PostController.CreatePost)
			posts.PUT("/:id", container.PostController.UpdatePost)
			posts.DELETE("/:id", container.PostController.DeletePost)
			posts.GET("/:id/comments", container.PostController.GetComments)
			posts.POST("/:id/comments", container.PostController.CreateComment)
			posts.DELETE("/:id/comments/:commentId", container.PostController.DeleteComment)
			posts.POST("/:id/vote", container.PostController.VotePoll)
		}

		// Subject-Teacher assignments
		subjects := v1.Group("/subjects")
		subjects.Use(http_middleware.JWTAuthentication)
		{
			subjects.POST("/:subjectId/teachers", container.SubjectController.AssignTeacher)
			subjects.DELETE("/:subjectId/teachers/:teacherId", container.SubjectController.RemoveTeacher)
		}

		// Teacher subjects (get subjects for a teacher)
		teachers := v1.Group("/teachers")
		teachers.Use(http_middleware.JWTAuthentication)
		{
			teachers.GET("/:teacherId/subjects", container.SubjectController.GetByTeacher)
		}

		// Activities management
		activities := v1.Group("/activities")
		activities.Use(http_middleware.JWTAuthentication)
		{
			activities.GET("/:activityId", container.ActivityController.GetById)
			activities.PUT("/:activityId", container.ActivityController.Update)
			activities.DELETE("/:activityId", container.ActivityController.Delete)
			// Teacher assignments
			activities.GET("/:activityId/teachers", container.ActivityController.GetTeachers)
			activities.POST("/:activityId/teachers", container.ActivityController.AssignTeacher)
			activities.DELETE("/:activityId/teachers/:teacherId", container.ActivityController.RemoveTeacher)
			// Student enrollments
			activities.GET("/:activityId/students", container.ActivityController.GetStudents)
			activities.POST("/:activityId/students", container.ActivityController.EnrollStudent)
			activities.DELETE("/:activityId/students/:studentId", container.ActivityController.RemoveStudent)
		}
	}

	log.Println("✅ All routes configured")
	return router
}

func StartServer(engine *gin.Engine) {
	address := fmt.Sprintf("0.0.0.0:%d", config.APP.Rest.Port)
	log.Printf("🚀 Server starting on %s", address)

	if err := engine.Run(address); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}
