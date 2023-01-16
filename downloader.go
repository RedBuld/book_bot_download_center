package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"sync"
)

var wg sync.WaitGroup

var args_help bool
var args_url string
var args_format string
var args_save string
var args_log string
var args_start int
var args_end int
var args_proxy string
var args_timeout int
var args_cover bool
var args_no_image bool
var args_login string
var args_password string

var sites_regexp = regexp.MustCompile(`https?:\/\/(www\.)*(?P<site>[^\/]*)\/`)
var allowed_formats = []string{"epub", "fb2", "cbz"}
var e2e_sites = []string{"acomics.ru", "author.today", "bigliba.com", "bookinbook.ru", "bookinist.pw", "booknet.com", "booknet.ua", "bookriver.ru", "bookstab.ru", "dark-novels.ru", "dreame.com", "eznovels.com", "fb2.top", "ficbook.net", "fictionbook.ru", "hentailib.me", "hogwartsnet.ru", "hotnovelpub.com", "hub-book.com", "ifreedom.su", "jaomix.ru", "ladylib.top", "lanovels.com", "libbox.ru", "libst.ru", "lightnoveldaily.com", "i-gram.ru", "litgorod.ru", "litmarket.ru", "litmir.me", "litnet.com", "litres.ru", "manga.ovh", "mangalib.me", "mir-knig.com", "mybook.ru", "noveltranslate.com", "novelxo.com", "online-knigi.com.ua", "prodaman.ru", "ranobe-novels.ru", "ranobe.ovh", "ranobehub.org", "ranobelib.me", "ranobes.com", "readli.net", "readmanga.live", "remanga.org", "renovels.org", "royalroad.com", "ru.novelxo.com", "samlib.ru", "tl.rulate.ru", "topliba.com", "twilightrussia.ru", "wattpad.com", "wuxiaworld.ru", "xn--80ac9aeh6f.xn--p1ai"}

func _inArray(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func e2e__download(site string) {
	defer wg.Done()

	_exec, _cmd, _cwd := e2e__prepare_command()

	cmd := exec.Command(_exec, _cmd...)
	cmd.Dir = _cwd

	if args_log != "" {
		_log, logErr := os.Create(args_log)
		if logErr != nil {
			fmt.Println(logErr)
			return
		}
		defer _log.Close()
		cmd.Stdout = _log
		cmd.Stderr = _log
		startErr := cmd.Run()
		if startErr != nil {
			fmt.Println(startErr)
			return
		}
	} else {
		r, w := io.Pipe()
		defer w.Close()
		cmd.Stdout = w
		cmd.Stderr = w
		scanner := bufio.NewScanner(r)
		go func() {
			for scanner.Scan() {
				fmt.Println(scanner.Text())
			}
		}()
		startErr := cmd.Run()
		if startErr != nil {
			fmt.Println(startErr)
			return
		}
	}
}

func e2e__prepare_command() (_exec string, _cmd []string, _cwd string) {
	_cwd, cwdErr := filepath.Abs("_Elib2Ebook")
	if cwdErr != nil {
		panic(cwdErr)
	}

	if runtime.GOOS == "windows" {
		_exec = "./Elib2Ebook.exe"
	} else {
		_exec = "./Elib2Ebook"
	}

	if args_url != "" {
		_cmd = append(_cmd, "--url", args_url)
	}
	if args_format != "" {
		_cmd = append(_cmd, "--format", args_format+",json_lite")
	}
	if args_save != "" {
		args_save, savePathErr := filepath.Abs(args_save)
		if savePathErr != nil {
			panic(savePathErr)
		}
		_path := filepath.Dir(args_save)
		_ = os.Mkdir(_path, os.ModePerm)

		_cmd = append(_cmd, "--save", args_save)
	}
	if args_log != "" {
		args_log, logPathErr := filepath.Abs(args_log)
		if logPathErr != nil {
			panic(logPathErr)
		}
		_path := filepath.Dir(args_log)
		_ = os.Mkdir(_path, os.ModePerm)
	}
	if args_start != 0 {
		_cmd = append(_cmd, "--start", strconv.Itoa(args_start))
	}
	if args_end != 0 {
		_cmd = append(_cmd, "--end", strconv.Itoa(args_end))
	}
	if args_proxy != "" {
		_cmd = append(_cmd, "--proxy", args_proxy)
		args_timeout = 120
	}
	if args_timeout > 0 {
		_cmd = append(_cmd, "--timeout", strconv.Itoa(args_timeout))
	}
	if args_cover {
		_cmd = append(_cmd, "--cover")
	}
	if args_no_image {
		_cmd = append(_cmd, "--no-image")
	}
	if args_login != "" && args_password != "" {
		_cmd = append(_cmd, "--login", args_login)
		_cmd = append(_cmd, "--password", args_password)
	}
	return _exec, _cmd, _cwd
}

func download() {
	defer wg.Done()

	match := sites_regexp.FindStringSubmatch(args_url)
	if len(match) < 3 {
		fmt.Println("Ссылка не найдена")
		os.Exit(1)
	}

	site := match[2]

	if _inArray(site, e2e_sites) {
		wg.Add(1)
		go e2e__download(site)
		return
	}
	fmt.Println("Ссылка не поддерживается")
	os.Exit(1)
}

func NewDownloader() {
	var args_error bool = false

	if args_url != "" || args_format != "" || args_save != "" || args_log != "" || args_start != 0 || args_end != 0 || args_proxy != "" || args_timeout != 60 || args_cover || args_no_image || args_login != "" || args_password != "" {
		args_help = false
	}

	if args_help {
		flag.Usage()
		os.Exit(0)
	}

	if len(args_url) == 0 {
		args_error = true
		fmt.Println("Не задан url")
		fmt.Println("")
	}
	if len(args_format) == 0 {
		args_error = true
		fmt.Println("Не задан format")
		fmt.Println("")
	}
	if !_inArray(args_format, allowed_formats) {
		args_error = true
		if args_format != "" {
			fmt.Println("Формат \"" + args_format + "\" неизвестен")
			fmt.Println("")
		}
	}
	if args_error {
		flag.Usage()
		os.Exit(1)
	}

	wg.Add(1)

	download()

	wg.Wait()
}
