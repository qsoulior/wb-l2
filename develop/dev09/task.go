package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"slices"

	"golang.org/x/net/html"
)

/*
=== Утилита wget ===

Реализовать утилиту wget с возможностью скачивать сайты целиком.
*/

// Структура-обертка над http-клиентом.
type Wget struct {
	http.Client
}

// Вспомогательный тип множества.
type set[T comparable] map[T]struct{}

// ExtractAttrs читает r и извлекает значения атрибутов для каждого html-тэга из tags.
// Возвращает мапу вида: тэг - множество значений атрибутов.
func (w Wget) ExtractAttrs(r io.Reader, tags map[string]set[string]) map[string]set[string] {
	values := make(map[string]set[string])
	tokenizer := html.NewTokenizer(r)
	for {
		switch tokenizer.Next() {
		case html.ErrorToken:
			return values
		case html.StartTagToken:
			token := tokenizer.Token()
			if attrs, ok := tags[token.Data]; ok {
				for _, attr := range token.Attr {
					if _, ok := attrs[attr.Key]; ok {
						if _, ok := values[token.Data]; !ok {
							values[token.Data] = make(set[string])
						}
						values[token.Data][attr.Val] = struct{}{}
					}
				}
			}
		}
	}
}

// ExtractLinks читает r и извлекает значения атрибутов для каждого html-тэга из tags,
// ожидая увидеть в качестве значений валидные url-адреса.
// Возвращает мапу вида: тэг - множество url-адресов.
func (w Wget) ExtractLinks(r io.Reader, host string, tags map[string]set[string]) (map[string][]string, error) {
	hostURL, err := url.Parse(host)
	if err != nil {
		return nil, err
	}

	links := make(map[string][]string)
	for tag, attrs := range w.ExtractAttrs(r, tags) {
		for attr := range attrs {
			linkURL, err := url.Parse(attr)
			if err != nil {
				continue
			}
			// Если url-адрес относительный, то добавляем scheme и host хост-адреса.
			if linkURL.Scheme == "" {
				linkURL.Scheme = hostURL.Scheme
			}
			if linkURL.Host == "" {
				linkURL.Host = hostURL.Host
			}
			links[tag] = append(links[tag], linkURL.String())
		}
	}

	return links, nil
}

// Путь до папки для хранения результата загрузки.
const RESULT_PATH = "result"

// DownloadFile запрашивает файл из src и создает его в папке RESULT_PATH.
func (w Wget) DownloadFile(src string) (io.Reader, error) {
	fileURL, err := url.Parse(src)
	if err != nil {
		return nil, err
	}

	filePath := filepath.Join(RESULT_PATH, fileURL.Host, fileURL.Path)

	resp, err := w.Get(src)
	if err != nil {
		return nil, err
	}

	extensions, _ := mime.ExtensionsByType(resp.Header.Get("Content-Type"))
	if len(extensions) > 0 && slices.Contains(extensions, ".html") {
		filePath += ".html"
	}

	var buf bytes.Buffer
	r := io.TeeReader(resp.Body, &buf)

	err = w.WriteFile(r, filePath)
	if err != nil {
		return nil, err
	}

	return &buf, nil
}

// WriteFile читает данные из r и записывает в файл по пути path.
// Создает файл и все родительские папки, если они не существует по пути path.
func (w Wget) WriteFile(r io.Reader, path string) error {
	err := os.MkdirAll(filepath.Dir(path), 0666)
	if err != nil {
		return err
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}

	_, err = io.Copy(f, r)
	f.Close()

	return err
}

// BFS обходит в ширину html-страницы с максимальной глубиной depth.
// Для каждой страницы также загружает медиа-файлы.
func (w Wget) BFS(start string, depth int) <-chan error {
	errCh := make(chan error)
	go func() {
		queue := []string{start}
		depths := map[string]int{start: 0}

		for len(queue) > 0 {
			src := queue[0]
			queue = queue[1:]

			body, err := w.DownloadFile(src)
			if err != nil {
				errCh <- err
				continue
			}

			links, err := w.ExtractLinks(body, src, map[string]set[string]{
				"a":   {"href": {}},
				"img": {"src": {}},
			})
			if err != nil {
				errCh <- err
				continue
			}

			for _, link := range links["img"] {
				_, err := w.DownloadFile(link)
				if err != nil {
					errCh <- err
				}
			}

			if depths[src] < depth {
				for _, link := range links["a"] {
					if _, ok := depths[link]; !ok {
						queue = append(queue, link)
						depths[link] = depths[src] + 1
					}
				}
			}
		}
		close(errCh)
	}()
	return errCh
}

func main() {
	// Парсим максимальную глубину обхода, переданную в качестве флага.
	var depth int
	flag.IntVar(&depth, "l", 0, "set the maximum number of subdirectories that Wget will recurse into to depth")
	flag.Parse()

	// Парсим аргумент целевого url-адреса.
	target := flag.Arg(0)
	if target == "" {
		fmt.Fprintln(os.Stderr, "target url is empty")
		os.Exit(2)
	}

	wget := Wget{}
	errCh := wget.BFS(target, depth)

	// Записываем ошибки в отдельный файл.
	file, err := os.Create("result.log")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create log file: %s", err)
		os.Exit(1)
	}
	for err := range errCh {
		fmt.Fprintln(file, err)
	}
	file.Close()
}
