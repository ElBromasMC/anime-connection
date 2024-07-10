package main

import (
	"alc/handler/admin"
	"alc/handler/admin/store"
	"alc/handler/admin/user"
	"alc/handler/public"
	"alc/handler/util"
	middle "alc/middleware"
	"alc/service"
	"context"
	"log"
	"net/http"
	"os"
	_ "time/tzdata"

	"github.com/gorilla/sessions"
	pgxuuid "github.com/jackc/pgx-gofrs-uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/wneessen/go-mail"
)

func main() {
	e := echo.New()
	if os.Getenv("ENV") == "development" {
		e.Debug = true
	}

	// Database connection
	dbconfig, err := pgxpool.ParseConfig(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalln("Failed to parse config:", err)
	}
	dbconfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		// Register uuid type
		pgxuuid.Register(conn.TypeMap())
		return nil
	}
	dbpool, err := pgxpool.NewWithConfig(context.Background(), dbconfig)
	if err != nil {
		log.Fatalln("Failed to connect to the database:", err)
	}
	defer dbpool.Close()

	// Initialize email client
	client, err := mail.NewClient(os.Getenv("SMTP_HOSTNAME"),
		mail.WithSMTPAuth(mail.SMTPAuthPlain), mail.WithTLSPortPolicy(mail.TLSMandatory),
		mail.WithUsername(os.Getenv("SMTP_USER")), mail.WithPassword(os.Getenv("SMTP_PASS")),
	)
	if err != nil {
		log.Fatalln("Failed to create email client:", err)
	}

	// Initialize services
	ps := service.NewPublicService(dbpool)
	as := service.NewAdminService(ps)
	ms := service.NewEmailService(client)
	us := service.NewAuthService(dbpool)
	cs := service.NewCommentService(dbpool)

	// Initialize handlers
	ph := public.Handler{
		PublicService:  ps,
		EmailService:   ms,
		AuthService:    us,
		CommentService: cs,
	}

	ah := admin.Handler{
		AdminService: as,
		AuthService:  us,
	}
	sh := store.Handler(ah)
	uh := user.Handler(ah)

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.RemoveTrailingSlashWithConfig(middleware.TrailingSlashConfig{
		RedirectCode: http.StatusMovedPermanently,
	}))
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{Level: 5}))
	key := os.Getenv("SESSION_KEY")
	if key == "" {
		log.Fatalln("Missing SESSION_KEY env variable")
	}
	e.Use(session.Middleware(sessions.NewCookieStore([]byte(key))))
	authMiddleware := middle.Auth(us)
	cartMiddleware := middle.Cart(ps)

	// Static files
	static(e)

	// Images routes
	e.Static("/images", "images")

	// Page routes
	e.GET("/", ph.HandleIndexShow, authMiddleware, cartMiddleware)
	e.GET("/nosotros", ph.HandleNosotrosShow, authMiddleware, cartMiddleware)
	e.GET("/contacto", ph.HandleContactoShow, authMiddleware, cartMiddleware)

	// Auth routes
	e.GET("/login", ph.HandleLoginShow)
	e.GET("/signup", ph.HandleSignupShow)
	e.POST("/login", ph.HandleLogin)
	e.POST("/signup", ph.HandleSignup)
	e.GET("/logout", ph.HandleLogout)

	// Store routes
	g2 := e.Group("/store")
	g2.Use(authMiddleware, cartMiddleware)
	g2.GET("", func(c echo.Context) error {
		return c.Redirect(http.StatusPermanentRedirect, "/store/categories/all")
	})
	g2.GET("/categories/all", ph.HandleStoreAllShow)
	g2.GET("/categories/:categorySlug", ph.HandleStoreCategoryShow)
	g2.GET("/categories/all/items", ph.HandleStoreAllItemsShow)
	g2.GET("/categories/:categorySlug/items", ph.HandleStoreCategoryItemsShow)
	g2.GET("/categories/:categorySlug/items/:itemSlug", ph.HandleStoreItemShow)
	g2.GET("/categories/:categorySlug/items/:itemSlug/comments", ph.HandleCommentsShow)
	g2.POST("/categories/:categorySlug/items/:itemSlug/comments", ph.HandleCommentInsertion)

	// Cart group
	g4 := e.Group("/cart")
	g4.Use(authMiddleware, cartMiddleware)
	g4.POST("", ph.HandleNewCartItem)
	g4.DELETE("", ph.HandleRemoveCartItem)

	// Checkout group
	g5 := e.Group("/checkout")
	g5.Use(authMiddleware, cartMiddleware)
	g5.GET("", ph.HandleCheckoutShow)
	g5.POST("", ph.HandleCheckoutOrderInsertion)
	g5.GET("/:orderID", ph.HandleCheckoutOrderShow)

	// Admin group
	g3 := e.Group("/admin")
	g3.Use(authMiddleware, middle.Admin)
	g3.GET("", ah.HandleIndexShow)

	// Admin store group
	g31 := g3.Group("/tienda")
	g31.Use(middle.RoleAdmin)
	g31.GET("", sh.HandleIndexShow)

	g31.GET("/type/:typeSlug/categories", sh.HandleCategoriesShow)
	g31.POST("/type/:typeSlug/categories", sh.HandleCategoryInsertion)
	g31.PUT("/type/:typeSlug/categories/:categorySlug", sh.HandleCategoryUpdate)
	g31.DELETE("/type/:typeSlug/categories/:categorySlug", sh.HandleCategoryDeletion)
	g31.GET("/type/:typeSlug/categories/insert", sh.HandleCategoryInsertionFormShow)
	g31.GET("/type/:typeSlug/categories/:categorySlug/update", sh.HandleCategoryUpdateFormShow)
	g31.GET("/type/:typeSlug/categories/:categorySlug/delete", sh.HandleCategoryDeletionFormShow)

	g31.GET("/type/:typeSlug/categories/:categorySlug/items", sh.HandleItemsShow)
	g31.POST("/type/:typeSlug/categories/:categorySlug/items", sh.HandleItemInsertion)
	g31.PUT("/type/:typeSlug/categories/:categorySlug/items/:itemSlug", sh.HandleItemUpdate)
	g31.DELETE("/type/:typeSlug/categories/:categorySlug/items/:itemSlug", sh.HandleItemDeletion)
	g31.GET("/type/:typeSlug/categories/:categorySlug/items/insert", sh.HandleItemInsertionFormShow)
	g31.GET("/type/:typeSlug/categories/:categorySlug/items/:itemSlug/update", sh.HandleItemUpdateFormShow)
	g31.GET("/type/:typeSlug/categories/:categorySlug/items/:itemSlug/delete", sh.HandleItemDeletionFormShow)

	g31.GET("/type/:typeSlug/categories/:categorySlug/items/:itemSlug/products", sh.HandleProductsShow)
	g31.POST("/type/:typeSlug/categories/:categorySlug/items/:itemSlug/products", sh.HandleProductInsertion)
	g31.PUT("/type/:typeSlug/categories/:categorySlug/items/:itemSlug/products/:productSlug", sh.HandleProductUpdate)
	g31.DELETE("/type/:typeSlug/categories/:categorySlug/items/:itemSlug/products/:productSlug", sh.HandleProductDeletion)
	g31.PUT("/type/:typeSlug/categories/:categorySlug/items/:itemSlug/products/:productSlug/stock", sh.HandleProductStockUpdate)
	g31.GET("/type/:typeSlug/categories/:categorySlug/items/:itemSlug/products/insert", sh.HandleProductInsertionFormShow)
	g31.GET("/type/:typeSlug/categories/:categorySlug/items/:itemSlug/products/:productSlug/update", sh.HandleProductUpdateFormShow)
	g31.GET("/type/:typeSlug/categories/:categorySlug/items/:itemSlug/products/:productSlug/delete", sh.HandleProductDeletionFormShow)
	g31.GET("/type/:typeSlug/categories/:categorySlug/items/:itemSlug/products/:productSlug/stock", sh.HandleProductStockUpdateFormShow)

	// Admin user group
	g32 := g3.Group("/usuarios")
	g32.Use(middle.RoleAdmin)
	g32.GET("", uh.HandleIndexShow)

	// Error handler
	e.HTTPErrorHandler = util.HTTPErrorHandler

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatalln(e.Start(":" + port))
}
