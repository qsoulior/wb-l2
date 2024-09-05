package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)

/*
=== Утилита telnet ===

Реализовать простейший telnet-клиент.

Примеры вызовов:
go-telnet --timeout=10s host port
go-telnet mysite.ru 8080
go-telnet --timeout=3s 1.1.1.1 123

Требования:
- Программа должна подключаться к указанному хосту (ip или доменное имя + порт) по протоколу TCP. После подключения STDIN программы должен записываться в сокет, а данные полученные и сокета должны выводиться в STDOUT
- Опционально в программу можно передать таймаут на подключение к серверу (через аргумент --timeout, по умолчанию 10s)
- При нажатии Ctrl+D программа должна закрывать сокет и завершаться. Если сокет закрывается со стороны сервера, программа должна также завершаться. При подключении к несуществующему сервер, программа должна завершаться через timeout
*/

// RedirectStream перенаправляет поток из r в w, прерывает выполнение если ctx.Done().
func RedirectStream(ctx context.Context, r io.ReadCloser, w io.Writer) error {
	go func() {
		<-ctx.Done()
		r.Close()
	}()

	buf := make([]byte, 4096)
	for {
		n, err := r.Read(buf)
		if err != nil {
			return err
		}
		_, err = fmt.Fprintf(w, "%s", buf[:n])
		if err != nil {
			return err
		}
	}
}

// Listen прослушивает tcp-соединение по адресу address.
// Дублирует полученные из соединения сообщения в это же соединение.
func Listen(ctx context.Context, address string) error {
	var lc net.ListenConfig
	listener, err := lc.Listen(ctx, "tcp", address)
	if err != nil {
		return fmt.Errorf("failed to listen %s: %w", address, err)
	}

	go func() {
		for {
			conn, err := listener.Accept()
			if err == nil {
				go RedirectStream(ctx, conn, conn)
			}
		}
	}()

	<-ctx.Done()
	listener.Close()
	return nil
}

// Dial создает tcp-соединение с сервером по адресу address, перенаправляет входной поток
// из stdin в соединение conn и выходной поток - из conn в stdout.
func Dial(ctx context.Context, timeout time.Duration, address string) error {
	g, ctx := errgroup.WithContext(ctx)

	dialer := net.Dialer{Timeout: timeout}
	conn, err := dialer.DialContext(ctx, "tcp", address)
	if err != nil {
		return fmt.Errorf("failed to connect %s: %w", address, err)
	}
	defer conn.Close()

	g.Go(func() error {
		return RedirectStream(ctx, os.Stdin, conn)
	})

	g.Go(func() error {
		return RedirectStream(ctx, conn, os.Stdout)
	})

	err = g.Wait()
	if !errors.Is(err, io.EOF) {
		return err
	}

	return nil
}

func main() {
	// Парсим опциональные для работы флаги.
	var (
		serverMode bool // флаг для запуска сервера вместо клиента
		timeout    time.Duration
	)
	flag.BoolVar(&serverMode, "s", false, "server mode")
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "connection timeout")
	flag.Parse()

	if len(flag.Args()) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: go-telnet [-s] [--timeout=duration] host port")
		os.Exit(2)
	}

	address := net.JoinHostPort(flag.Arg(0), flag.Arg(1))

	// Слушаем сигналы Ctrl+D и Ctrl+C.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGQUIT, syscall.SIGINT)
	defer stop()

	// Если приложение запущено в режиме сервера, то запускаем сервер.
	if serverMode {
		err := Listen(ctx, address)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		return
	}

	// Иначе запускаем клиент.
	err := Dial(ctx, timeout, address)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
