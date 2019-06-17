/*
 * Copyright © 2019 Hedzr Yeh.
 */

package restful

import (
	"fmt"
	"github.com/hedzr/cmdr"
	"github.com/hedzr/cmdr/conf"
	voxr_common "github.com/hedzr/voxr-common"
	"github.com/hedzr/voxr-common/db/dbi"
	"github.com/hedzr/voxr-common/tool"
	"github.com/hedzr/voxr-common/vxconf"
	"github.com/hedzr/voxr-common/xs"
	"github.com/hedzr/voxr-common/xs/mjwt"
	"github.com/hedzr/voxr-lite/internal/scheduler"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/sirupsen/logrus"
	"html/template"
	"io"
	"os"
	"strings"
)

type (
	myService struct {
		// es *xs.echoServerImpl
		// dbConfig *dbi.Config
		handlers Handlers
	}

	Handlers interface {
		OnGetBanner() string
		InitRoutes(e *echo.Echo, s vxconf.CoolServer) (ready bool)
		InitWebSocket(e *echo.Echo, s vxconf.CoolServer) (ready bool)
	}
)

var (
	myXsServerCoreImpl vxconf.CoolServer
	xsServer           vxconf.XsServer
)

func NewXsServer(cmd *cmdr.Command, args []string, stopCh, doneCh chan struct{}, h Handlers) (err error) {
	logrus.Infof("    Starting XS Server...")
	myXsServerCoreImpl = newXsService(h)
	xsServer = xs.New(myXsServerCoreImpl)

	xs.SetRegistrarOkCallback(scheduler.RegistrarOkCallback)
	xs.SetRegistrarChangesHandler(scheduler.RegistrarChangesHandler)

	xsServer.Start(stopCh, doneCh)

	return
}

// New() return the implementation of XsServer, a RESTful Echo server
func newXsService(h Handlers) *myService {
	return &myService{handlers: h}
}

func (s *myService) GetApiPrefix() string {
	return voxr_common.GetApiPrefix()
}

func (s *myService) GetApiVersion() string {
	return voxr_common.GetApiVersion()
}

func (s *myService) OnPrintBanner(addr string) {
	fmt.Printf("%s\n%s %s\n%s\n%s\n@@ %s running at %s, pid=%d. runmode=%v\n%s\n",
		s.handlers.OnGetBanner(),
		conf.AppName, "XS-SERVER",
		"__________________________________________O/_______",
		"                                          O\\",
		vxconf.GetStringR("server.serviceName", "voxr-lite"),
		addr,
		os.Getpid(),
		vxconf.GetStringR("runmode", "devel"),
		"Press CONTROL-C to exit.")
}

func (s *myService) OnInitRegistrar() { // never used
}

func (s *myService) OnShutdownRegistrar() {
}

func (s *myService) OnInitDB(e *echo.Echo) (dbiConfig *dbi.Config) {
	// s.dbConfig = webui.InitDB(e)
	// return s.dbConfig
	return
}

func (s *myService) OnShutdownDB() {
	// db.Close()
}

func (s *myService) onGetFirstPages() []string {
	return []string{
		"/index.html",
		"/",
		"", // generally equals "/"
	}
}
func (s *myService) onGetEveryoneUrls() []string {
	return []string{
		"/favicon.ico",
		"/login", "/signup", "/login.html", "/signup.html",
		"/random.html",
		"/crypt", "/crypt.c",
		"/refresh-token", "/refresh_token",
	}
}
func (s *myService) onGetEveryonePrefixes() []string {
	return []string{
		"/css/", "/js/", "/img/",
		"/images/", "/stylesheets/", "/javascripts/",
		"/health",
		"/routes",
		"/public",
		"/upload",
		"/v2/ws",
		"/v1/ws",                    // ws entry
		"/send",                     // ws test page
		"/bal", "/bal-rnd", "/time", // balancer entry
	}
}

// OnInitCORS https://echo.labstack.com/middleware/cors
func (s *myService) OnInitCORS(config *middleware.CORSConfig) echo.MiddlewareFunc {
	config.AllowHeaders = []string{
		echo.HeaderOrigin, echo.HeaderContentType,
		echo.HeaderAccept, echo.HeaderAuthorization,
	}
	return middleware.CORSWithConfig(*config)
}

// OnInitJWT give a chance to hook the Claims Got Event, or you could modify the `config` and return it to xs-server framework.
func (s *myService) OnInitJWT(e *echo.Echo, config mjwt.JWTConfig) mjwt.JWTConfig {
	return config
}

func (s *myService) JWTSkipper(c echo.Context) (skip bool) {
	skip = false
	path := c.Path()
	if strings.HasPrefix(path, s.GetApiPrefix()) {
		path = path[len(s.GetApiPrefix()):]
	}

	for _, s := range s.onGetFirstPages() {
		if strings.EqualFold(path, s) {
			return true
		}
	}

	// Skip authentication for and signup login requests
	for _, s := range s.onGetEveryoneUrls() {
		if strings.EqualFold(path, s) {
			return true
		}
	}

	if vxconf.GetBoolR("server.webui.mustLogin", false) {
		// TODO c.Redirect(http.StatusTemporaryRedirect, webui.LoginUrl())
		return
	}

	for _, s := range s.onGetEveryonePrefixes() {
		if strings.HasPrefix(path, s) {
			return true
		}
	}

	urlPrefix := vxconf.GetStringR("server.static.index", "index.html") // webui.StaticPagesUrlPrefix()
	if strings.HasPrefix(path, urlPrefix) {
		return true
	}

	return
}

// func (s *myService) OnInitForwarders(e *echo.Echo) {
// 	// forwarder.Init(e)
// }
//
// func (s *myService) OnShutdownForwarders() {
// 	// forwarder.Shutdown()
// }

func (s *myService) OnInitStats(e *echo.Echo) (ready bool) {
	ready = false // use default stats configurations
	return
}

func (s *myService) OnInitStatic(e *echo.Echo) (ready bool) {
	// e.Static("/public", tool.GetCurrentDir()+"/public")

	e.File("/favicon.ico", tool.GetCurrentDir()+"/public/images/favicon.ico")

	// logrus.Println(tool.GetCurrentDir() + "/public : static root.")

	ready = false // use default static pages configurations
	return
}

// TemplateRenderer is a custom html/template renderer for Echo framework
type TemplateRenderer struct {
	templates *template.Template
}

// Render renders a template document
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {

	// Add global methods if data is a map
	if viewContext, isMap := data.(map[string]interface{}); isMap {
		viewContext["reverse"] = c.Echo().Reverse
	}

	return t.templates.ExecuteTemplate(w, name, data)
}

func (s *myService) OnInitTemplates(e *echo.Echo) (ready bool) {
	ready = false // use default static pages configurations
	// renderer := &TemplateRenderer{
	// 	templates: template.Must(template.ParseGlob(webui.GetTemplatesPath("*.html"))),
	// }
	// e.Renderer = renderer
	ready = true
	return
}

func (s *myService) OnInitRoutes(e *echo.Echo) (ready bool) {
	ready = false

	ready = s.handlers.InitRoutes(e, s)

	// /home -> static SPA app -> admin web ui

	if tool.FileExists(tool.GetCurrentDir() + "/public/index.html") {
		e.File("/send", tool.GetCurrentDir()+"/public/index.html")
	}

	return
}

func (s *myService) OnInitWebSocket(e *echo.Echo) (ready bool) {
	s.handlers.InitWebSocket(e, s)
	return
}

func (s *myService) OnPreStart(e *echo.Echo) (err error) {
	logrus.Debugf("OnPreStart")
	return
}

func (s *myService) OnPreShutdown() {
	logrus.Debugf("OnPreShutdown")
}

func (s *myService) OnPostShutdown() {
	logrus.Debugf("OnPostShutdown")
}
