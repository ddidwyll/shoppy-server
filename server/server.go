package server

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	//"golang.org/x/crypto/acme/autocert"
	"shoppy/cnst"
	"shoppy/model"

	"io"
	"log"
	"os"
	"strconv"
	"time"
)

/* USER */

func insUser(c echo.Context) error {
	return upsert(c, cnst.USER, cnst.INS)
}
func updUser(c echo.Context) error {
	return upsert(c, cnst.USER, cnst.UPD)
}
func getUsers(c echo.Context) error {
	return getAll(c, cnst.USER)
}
func delUser(c echo.Context) error {
	return del(c, cnst.USER)
}

/* FABRIC */

func insProduct(c echo.Context) error {
	return upsert(c, cnst.PRDT, cnst.INS)
}
func updProduct(c echo.Context) error {
	return upsert(c, cnst.PRDT, cnst.UPD)
}
func getFabrics(c echo.Context) error {
	return getAll(c, cnst.PRDT)
}
func delProduct(c echo.Context) error {
	return del(c, cnst.PRDT)
}

/* ORDER */

func insOrder(c echo.Context) error {
	return upsert(c, cnst.ORDR, cnst.INS)
}
func updOrder(c echo.Context) error {
	return upsert(c, cnst.ORDR, cnst.UPD)
}
func getOrders(c echo.Context) error {
	return getAll(c, cnst.ORDR)
}
func delOrder(c echo.Context) error {
	return del(c, cnst.ORDR)
}

/* PASSWORD */

func updPass(c echo.Context) error {
	return upsert(c, cnst.PASS, cnst.UPD)
}

// USERID
func userId(c echo.Context) string {
	user, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return "without_auth"
	}
	claims := user.Claims.(jwt.MapClaims)
	id := claims["id"].(string)
	return id
}

// MODTIMES
func modTimes(c echo.Context) error {
	t := new(model.MTime)
	t.Get()
	return c.JSON(200, t)
}

// LOGIN
func login(c echo.Context) error {
	p := new(model.Password)
	if err := c.Bind(p); err != nil {
		return c.JSONBlob(400,
			[]byte("{\"message\":\""+cnst.INCLOGIN+"\"}"))
	}
	u, err := p.Login()
	if err != nil {
		return c.JSONBlob(400,
			[]byte("{\"message\":\""+err.Error()+"\"}"))
	}
	token := jwt.New(jwt.SigningMethodHS256)
	exp := time.Now().Add(time.Hour * 168).Unix()
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = u.Id
	claims["exp"] = exp
	t, err := token.SignedString([]byte(cnst.SECRET))
	if err != nil {
		return err
	}
	return c.JSON(200, map[string]string{
		"token":  t,
		"id":     u.Id,
		"expire": strconv.FormatInt(exp, 10),
	})
}

// UPSERT
func upsert(c echo.Context, t, action byte) error {
	m := model.Model(t)
	log.Println(m)
	user := userId(c)
	if err := c.Bind(m); err != nil {
		log.Println(m)
		log.Println(err)
		return c.JSONBlob(cnst.Status(err, action+3, t))
	}
	if err := model.Upsert(t, action, m, user); err != nil {
		return c.JSONBlob(cnst.Status(err, action+3, t))
	}
	return c.JSONBlob(cnst.Status(nil, action, t))
}

// GETALL
func getAll(c echo.Context, t byte) error {
	user := userId(c)
	all, err := model.GetAll(t, user)
	if err != nil {
		return c.JSON(400, err.Error())
	}

	return c.JSONBlob(200, []byte(all))
}

// DELONE
func del(c echo.Context, t byte) error {
	user := userId(c)
	err := model.Delete(t, c.Param("id"), user)
	if err != nil {
		return c.JSONBlob(cnst.Status(err, cnst.NTDEL, t))
	}
	return c.JSONBlob(cnst.Status(nil, cnst.DEL, t))
}

func upload(c echo.Context) error {
	name := c.FormValue("name")
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()
	dst, err := os.Create(cnst.STATIC + "/img/" + name)
	if err != nil {
		return err
	}
	defer dst.Close()
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}
	log.Println(file.Header)
	return c.JSONBlob(200, []byte("{\"message\":\"Image successeful uploaded\"}"))
}

func Start() {
	e := echo.New()
	//e.AutoTLSManager.HostPolicy = autocert.HostWhitelist("<OPRIME>")
	//e.AutoTLSManager.Cache = autocert.DirCache("./.cache")
	//e.Pre(middleware.HTTPSNonWWWRedirect())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root:   cnst.STATIC,
		HTML5:  true,
		Browse: false,
	}))
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{
			echo.GET,
			echo.PUT,
			echo.POST,
			echo.DELETE},
	}))
	r := e.Group("/api")
	p := e.Group("/pub")
	p.POST("/login", login)
	p.GET("/mtimes", modTimes)
	r.Use(middleware.JWT([]byte(cnst.SECRET)))
	// PRODUCTS
	r.POST("/products", insProduct)
	p.GET("/products", getFabrics)
	r.PUT("/products", updProduct)
	r.DELETE("/products/:id", delProduct)
	// ORDERS
	r.POST("/orders", insOrder)
	r.GET("/orders", getOrders)
	//r.PUT("/orders", updOrder)
	r.DELETE("/orders/:id", delOrder)
	// USERS
	r.POST("/users", insUser)
	r.GET("/users", getUsers)
	r.PUT("/users", updUser)
	r.DELETE("/users/:id", delUser)
	// PASSWORDS
	r.PUT("/pass", updPass)
	// IMAGES
	r.POST("/upload", upload)

	e.Logger.Fatal(e.Start(":5001"))
	//e.Logger.Fatal(e.StartAutoTLS("<OPRIME>:431"))
}
