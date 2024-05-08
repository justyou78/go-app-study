package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"golang.org/x/sync/errgroup"
)

func main() {
	if len(os.Args) != 2 {
		log.Printf("need port number\n")
		os.Exit(1)
	}
	p := os.Args[1]

	l, err := net.Listen("tcp", ":"+p)
	if err != nil {
		log.Printf("failed to listen port %s: %v", p, err)
	}
	if err := run(context.Background(), l); err != nil {
		log.Printf("failed to terminate server: %v", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, l net.Listener) error {
	// 인수로 받은 net.LIstener를 이용하므로 Addr 필드는 지정하지 않는다.
	s := &http.Server{Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
	})}

	eg, ctx := errgroup.WithContext(ctx)
	// 다른 고루틴에서 HTTP 서버를 실행한다.
	eg.Go(func() error {
		if err := s.Serve(l); err != nil &&
			// http.ErrServerClosed는 http.Server.Shutdown()이 정상 종료된 것을 나타내므로 이상 처리가 아니다.
			err != http.ErrServerClosed {
			log.Printf("failed to close: %+v", err)
			return err
		}
		return nil
	})

	<-ctx.Done()
	if err := s.Shutdown(context.Background()); err != nil {
		log.Printf("failed to shutdown: %+v", err)
	}

	return eg.Wait()
}

// func run(ctx context.Context) error {
// 	// net/http는 기본 설정에서 병렬 요청을 받을 수 있으므로 별도의 미들웨어를 준비할 필요가 없다.
// 	// 미들웨어: 소프트웨어 컴포넌트로, 주로 서버와 클라이언트 사이의 통신을 처리하거나 데이터를 가공하는 역할. 주로 애플리케이션 수준에서 동작.

// 	err := http.ListenAndServe(
// 		":18080",
// 		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
// 		}),
// 	)

// 	if err != nil {
// 		fmt.Printf("failed to terminate server: %v", err)
// 		os.Exit(1)
// 	}
// 	return err
// }
