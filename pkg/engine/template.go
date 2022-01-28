package engine

import (
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	logs "github.com/bhojpur/logger/pkg/engine"
	"github.com/bhojpur/web/pkg/core/utils"
)

var (
	webTplFuncMap             = make(template.FuncMap)
	webViewPathTemplateLocked = false
	// webViewPathTemplates caching map and supported template file extensions per view
	webViewPathTemplates = make(map[string]map[string]*template.Template)
	templatesLock        sync.RWMutex
	// webTemplateExt stores the template extension which will build
	webTemplateExt = []string{"tpl", "html", "gohtml"}
	// webTemplatePreprocessors stores associations of extension -> preprocessor handler
	webTemplateEngines = map[string]templatePreProcessor{}
	webTemplateFS      = defaultFSFunc
)

// ExecuteTemplate applies the template with name  to the specified data object,
// writing the output to wr.
// A template will be executed safely in parallel.
func ExecuteTemplate(wr io.Writer, name string, data interface{}) error {
	return ExecuteViewPathTemplate(wr, name, BasConfig.WebConfig.ViewsPath, data)
}

// ExecuteViewPathTemplate applies the template with name and from specific viewPath to the specified data object,
// writing the output to wr.
// A template will be executed safely in parallel.
func ExecuteViewPathTemplate(wr io.Writer, name string, viewPath string, data interface{}) error {
	if BasConfig.RunMode == DEV {
		templatesLock.RLock()
		defer templatesLock.RUnlock()
	}
	if webTemplates, ok := webViewPathTemplates[viewPath]; ok {
		if t, ok := webTemplates[name]; ok {
			var err error
			if t.Lookup(name) != nil {
				err = t.ExecuteTemplate(wr, name, data)
			} else {
				err = t.Execute(wr, data)
			}
			if err != nil {
				logs.Trace("template Execute err:", err)
			}
			return err
		}
		panic("can't find templatefile in the path:" + viewPath + "/" + name)
	}
	panic("Unknown view path:" + viewPath)
}

func init() {
	webTplFuncMap["dateformat"] = DateFormat
	webTplFuncMap["date"] = Date
	webTplFuncMap["compare"] = Compare
	webTplFuncMap["compare_not"] = CompareNot
	webTplFuncMap["not_nil"] = NotNil
	webTplFuncMap["not_null"] = NotNil
	webTplFuncMap["substr"] = Substr
	webTplFuncMap["html2str"] = HTML2str
	webTplFuncMap["str2html"] = Str2html
	webTplFuncMap["htmlquote"] = Htmlquote
	webTplFuncMap["htmlunquote"] = Htmlunquote
	webTplFuncMap["renderform"] = RenderForm
	webTplFuncMap["assets_js"] = AssetsJs
	webTplFuncMap["assets_css"] = AssetsCSS
	webTplFuncMap["config"] = GetConfig
	webTplFuncMap["map_get"] = MapGet

	// Comparisons
	webTplFuncMap["eq"] = eq // ==
	webTplFuncMap["ge"] = ge // >=
	webTplFuncMap["gt"] = gt // >
	webTplFuncMap["le"] = le // <=
	webTplFuncMap["lt"] = lt // <
	webTplFuncMap["ne"] = ne // !=

	webTplFuncMap["urlfor"] = URLFor // build a URL to match a Controller and it's method
}

// AddFuncMap let user to register a func in the template.
func AddFuncMap(key string, fn interface{}) error {
	webTplFuncMap[key] = fn
	return nil
}

type templatePreProcessor func(root, path string, funcs template.FuncMap) (*template.Template, error)

type templateFile struct {
	root  string
	files map[string][]string
}

// visit will make the paths into two part,the first is subDir (without tf.root),the second is full path(without tf.root).
// if tf.root="views" and
// paths is "views/errors/404.html",the subDir will be "errors",the file will be "errors/404.html"
// paths is "views/admin/errors/404.html",the subDir will be "admin/errors",the file will be "admin/errors/404.html"
func (tf *templateFile) visit(paths string, f os.FileInfo, err error) error {
	if f == nil {
		return err
	}
	if f.IsDir() || (f.Mode()&os.ModeSymlink) > 0 {
		return nil
	}
	if !HasTemplateExt(paths) {
		return nil
	}

	replace := strings.NewReplacer("\\", "/")
	file := strings.TrimLeft(replace.Replace(paths[len(tf.root):]), "/")
	subDir := filepath.Dir(file)

	tf.files[subDir] = append(tf.files[subDir], file)
	return nil
}

// HasTemplateExt return this path contains supported template extension of Bhojpur.NET Platform application or not.
func HasTemplateExt(paths string) bool {
	for _, v := range webTemplateExt {
		if strings.HasSuffix(paths, "."+v) {
			return true
		}
	}
	return false
}

// AddTemplateExt add new extension for template.
func AddTemplateExt(ext string) {
	for _, v := range webTemplateExt {
		if v == ext {
			return
		}
	}
	webTemplateExt = append(webTemplateExt, ext)
}

// AddViewPath adds a new path to the supported view paths.
//Can later be used by setting a controller ViewPath to this folder
//will panic if called after bhojpur.Run()
func AddViewPath(viewPath string) error {
	if webViewPathTemplateLocked {
		if _, exist := webViewPathTemplates[viewPath]; exist {
			return nil //Ignore if viewpath already exists
		}
		panic("Can not add new view paths after bhojpur.Run()")
	}
	webViewPathTemplates[viewPath] = make(map[string]*template.Template)
	return BuildTemplate(viewPath)
}

func lockViewPaths() {
	webViewPathTemplateLocked = true
}

// BuildTemplate will build all template files in a directory.
// it makes Bhojpur.NET Platform can render any template file in view directory.
func BuildTemplate(dir string, files ...string) error {
	var err error
	fs := webTemplateFS()
	f, err := fs.Open(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return errors.New("dir open err")
	}
	defer f.Close()

	webTemplates, ok := webViewPathTemplates[dir]
	if !ok {
		panic("Unknown view path: " + dir)
	}
	self := &templateFile{
		root:  dir,
		files: make(map[string][]string),
	}
	err = Walk(fs, dir, func(path string, f os.FileInfo, err error) error {
		return self.visit(path, f, err)
	})
	if err != nil {
		fmt.Printf("Walk() returned %v\n", err)
		return err
	}
	buildAllFiles := len(files) == 0
	for _, v := range self.files {
		for _, file := range v {
			if buildAllFiles || utils.InSlice(file, files) {
				templatesLock.Lock()
				ext := filepath.Ext(file)
				var t *template.Template
				if len(ext) == 0 {
					t, err = getTemplate(self.root, fs, file, v...)
				} else if fn, ok := webTemplateEngines[ext[1:]]; ok {
					t, err = fn(self.root, file, webTplFuncMap)
				} else {
					t, err = getTemplate(self.root, fs, file, v...)
				}
				if err != nil {
					logs.Error("parse template err:", file, err)
					templatesLock.Unlock()
					return err
				}
				webTemplates[file] = t
				templatesLock.Unlock()
			}
		}
	}
	return nil
}

func getTplDeep(root string, fs http.FileSystem, file string, parent string, t *template.Template) (*template.Template, [][]string, error) {
	var fileAbsPath string
	var rParent string
	var err error
	if strings.HasPrefix(file, "../") {
		rParent = filepath.Join(filepath.Dir(parent), file)
		fileAbsPath = filepath.Join(root, filepath.Dir(parent), file)
	} else {
		rParent = file
		fileAbsPath = filepath.Join(root, file)
	}
	f, err := fs.Open(fileAbsPath)
	if err != nil {
		panic("can't find template file:" + file)
	}
	defer f.Close()
	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, [][]string{}, err
	}
	t, err = t.New(file).Parse(string(data))
	if err != nil {
		return nil, [][]string{}, err
	}
	reg := regexp.MustCompile(BasConfig.WebConfig.TemplateLeft + "[ ]*template[ ]+\"([^\"]+)\"")
	allSub := reg.FindAllStringSubmatch(string(data), -1)
	for _, m := range allSub {
		if len(m) == 2 {
			tl := t.Lookup(m[1])
			if tl != nil {
				continue
			}
			if !HasTemplateExt(m[1]) {
				continue
			}
			_, _, err = getTplDeep(root, fs, m[1], rParent, t)
			if err != nil {
				return nil, [][]string{}, err
			}
		}
	}
	return t, allSub, nil
}

func getTemplate(root string, fs http.FileSystem, file string, others ...string) (t *template.Template, err error) {
	t = template.New(file).Delims(BasConfig.WebConfig.TemplateLeft, BasConfig.WebConfig.TemplateRight).Funcs(webTplFuncMap)
	var subMods [][]string
	t, subMods, err = getTplDeep(root, fs, file, "", t)
	if err != nil {
		return nil, err
	}
	t, err = _getTemplate(t, root, fs, subMods, others...)

	if err != nil {
		return nil, err
	}
	return
}

func _getTemplate(t0 *template.Template, root string, fs http.FileSystem, subMods [][]string, others ...string) (t *template.Template, err error) {
	t = t0
	for _, m := range subMods {
		if len(m) == 2 {
			tpl := t.Lookup(m[1])
			if tpl != nil {
				continue
			}
			//first check filename
			for _, otherFile := range others {
				if otherFile == m[1] {
					var subMods1 [][]string
					t, subMods1, err = getTplDeep(root, fs, otherFile, "", t)
					if err != nil {
						logs.Trace("template parse file err:", err)
					} else if len(subMods1) > 0 {
						t, err = _getTemplate(t, root, fs, subMods1, others...)
					}
					break
				}
			}
			//second check define
			for _, otherFile := range others {
				var data []byte
				fileAbsPath := filepath.Join(root, otherFile)
				f, err := fs.Open(fileAbsPath)
				if err != nil {
					f.Close()
					logs.Trace("template file parse error, not success open file:", err)
					continue
				}
				data, err = ioutil.ReadAll(f)
				f.Close()
				if err != nil {
					logs.Trace("template file parse error, not success read file:", err)
					continue
				}
				reg := regexp.MustCompile(BasConfig.WebConfig.TemplateLeft + "[ ]*define[ ]+\"([^\"]+)\"")
				allSub := reg.FindAllStringSubmatch(string(data), -1)
				for _, sub := range allSub {
					if len(sub) == 2 && sub[1] == m[1] {
						var subMods1 [][]string
						t, subMods1, err = getTplDeep(root, fs, otherFile, "", t)
						if err != nil {
							logs.Trace("template parse file err:", err)
						} else if len(subMods1) > 0 {
							t, err = _getTemplate(t, root, fs, subMods1, others...)
							if err != nil {
								logs.Trace("template parse file err:", err)
							}
						}
						break
					}
				}
			}
		}

	}
	return
}

type templateFSFunc func() http.FileSystem

func defaultFSFunc() http.FileSystem {
	return FileSystem{}
}

// SetTemplateFSFunc set default filesystem function
func SetTemplateFSFunc(fnt templateFSFunc) {
	beeTemplateFS = fnt
}

// SetViewsPath sets view directory path in Bhojpur.NET Platform application.
func SetViewsPath(path string) *HttpServer {
	BasConfig.WebConfig.ViewsPath = path
	return BhojpurApp
}

// SetStaticPath sets static directory path and proper url pattern in Bhojpur.NET Platform application.
// if bhojpur.SetStaticPath("static","public"), visit /static/* to load static file in folder "public".
func SetStaticPath(url string, path string) *HttpServer {
	if !strings.HasPrefix(url, "/") {
		url = "/" + url
	}
	if url != "/" {
		url = strings.TrimRight(url, "/")
	}
	BasConfig.WebConfig.StaticDir[url] = path
	return BhojpurApp
}

// DelStaticPath removes the static folder setting in this url pattern in Bhojpur.NET Platform application.
func DelStaticPath(url string) *HttpServer {
	if !strings.HasPrefix(url, "/") {
		url = "/" + url
	}
	if url != "/" {
		url = strings.TrimRight(url, "/")
	}
	delete(BasConfig.WebConfig.StaticDir, url)
	return BhojpurApp
}

// AddTemplateEngine add a new templatePreProcessor which support extension
func AddTemplateEngine(extension string, fn templatePreProcessor) *HttpServer {
	AddTemplateExt(extension)
	webTemplateEngines[extension] = fn
	return BhojpurApp
}
