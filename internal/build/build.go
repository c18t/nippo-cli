package build

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

type RunEFunc func(cmd *cobra.Command, args []string) error

func CreateCmdFunc() RunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		err := downloadNippoData()
		cobra.CheckErr(err)

		err = buildIndexPage()
		cobra.CheckErr(err)

		return nil
	}
}

// page content
type Content struct {
	PageTitle string
	Date      string
	Content   template.HTML
}

func buildIndexPage() error {
	// Find home directory.
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	defaultDataDir := path.Join(home, ".local", "share")
	dataDir := os.Getenv("XDG_DATA_HOME")
	if dataDir == "" || !path.IsAbs(dataDir) {
		dataDir = defaultDataDir
	}
	dataDir = path.Join(dataDir, "nippo")

	t, err := template.ParseGlob(path.Join(dataDir, "templates", "*.html"))
	if err != nil {
		fmt.Println(err)
		return nil
	}
	tmpl := template.Must(template.Must(t.Lookup("layout").Clone()).AddParseTree("content", t.Lookup("index").Tree))

	f, err := os.Create("index.html")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer f.Close()

	// XDG_CACHE_HOMEディレクトリを取得
	defaultCacheDir := path.Join(home, ".cache")
	cacheDir := os.Getenv("XDG_CACHE_HOME")
	if cacheDir == "" || !path.IsAbs(cacheDir) {
		cacheDir = defaultCacheDir
	}
	cacheDir = path.Join(cacheDir, "nippo", "md")

	files, err := os.ReadDir(cacheDir)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	sort.Slice(files, func(i, j int) bool { return files[i].Name() < files[j].Name() })
	var fileName string
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		fileName = file.Name()
		continue
	}

	nippo, err := os.Open(path.Join(cacheDir, fileName))
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer nippo.Close()
	nippoData, err := io.ReadAll(nippo)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	tmpl.ExecuteTemplate(f, "layout", Content{
		Date:    strings.TrimSuffix(fileName, filepath.Ext(fileName)),
		Content: template.HTML(parseMarkdownToHtml(nippoData)),
	})
	return nil
}

func parseMarkdownToHtml(nippoData []byte) []byte {
	extensions := parser.CommonExtensions | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(nippoData)

	// create HTML renderer with extensions
	htmlFlags := html.CommonFlags
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	return markdown.Render(doc, renderer)
}

// download nippo data in google drive
func downloadNippoData() error {
	// Find home directory.
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	defaultDataDir := path.Join(home, ".local", "share")
	dataDir := os.Getenv("XDG_DATA_HOME")
	if dataDir == "" || !path.IsAbs(dataDir) {
		dataDir = defaultDataDir
	}
	dataDir = path.Join(dataDir, "nippo")

	b, err := os.ReadFile(path.Join(dataDir, "credentials.json"))
	if err != nil {
		fmt.Printf("Unable to read client secret file: %v\n", err)
		return nil
	}

	config, err := google.ConfigFromJSON(b, drive.DriveMetadataReadonlyScope, drive.DriveReadonlyScope)
	if err != nil {
		fmt.Printf("Unable to parse client secret file to config: %v\n", err)
		return nil
	}
	client := getClient(config)

	ctx := context.Background()
	srv, err := drive.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		fmt.Printf("Unable to retrieve Drive client: %v\n", err)
	}

	r, err := srv.Files.List().
		Q("parents in '1FZEaqRa8NmuRheHjTiW-_gUP3E5Ddw2T' and fileExtension = 'md'").
		// OrderBy("modifiedTime desc").
		OrderBy("name desc").
		Fields("nextPageToken, files(id, name, fileExtension)").
		PageSize(3).Do()
	if err != nil {
		fmt.Printf("Unable to retrieve files: %v\n", err)
	}
	fmt.Println("Files:")

	// XDG_CACHE_HOMEディレクトリを取得
	defaultCacheDir := path.Join(home, ".cache")
	cacheDir := os.Getenv("XDG_CACHE_HOME")
	if cacheDir == "" || !path.IsAbs(cacheDir) {
		cacheDir = defaultCacheDir
	}
	cacheDir = path.Join(cacheDir, "nippo", "md")
	err = os.MkdirAll(cacheDir, 0755)
	if err != nil && !os.IsExist(err) {
		fmt.Println(err)
		return nil
	}

	if len(r.Files) == 0 {
		fmt.Println("No files found.")
	} else {
		for _, i := range r.Files {
			fmt.Printf("%s (%s)\n", i.Name, i.Id)
			err = downloadFile(srv.Files, i, cacheDir)
			if err != nil {
				fmt.Printf("download failed: %v\n", err)
				continue
			}
		}
	}

	return nil
}

func downloadFile(fsrv *drive.FilesService, i *drive.File, folderPath string) error {
	g, err := fsrv.Get(i.Id).Download()
	if err != nil {
		return err
	}
	defer g.Body.Close()

	f, err := os.Create(filepath.Join(folderPath, i.Name))
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer f.Close()

	_, err = io.Copy(f, g.Body)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return nil
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	// Find home directory.
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)
	defaultDataDir := path.Join(home, ".local", "share")
	dataDir := os.Getenv("XDG_DATA_HOME")
	if dataDir == "" || !path.IsAbs(dataDir) {
		dataDir = defaultDataDir
	}
	dataDir = path.Join(dataDir, "nippo")
	tok, err := tokenFromFile(path.Join(dataDir, "token.json"))
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}
