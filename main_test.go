package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"

	"golang.org/x/sync/errgroup"
)

/*
	go test -v ./...
		./...: 현재 디렉토리부터 하위 모든 디렉토리를 테스트한다.
		-v (verbose): 테스트 실행의 상세한 출력을 제공
*/
/*
	t.Fatal: 테스트르를 실패로 표시하고 테스트를 즉시 중단.
	t.Error: 테스트를 실패료 표시. 테스트를 중단하지 않는다.
*/
/*
	fmt.print: 서식 지정이나 형식 변환이 없는 출력.
	fmt.Print("Name: ", name, ", Age: ", age) => Name: Alice, Age: 30
	fmt.printf: C 스타일의 형식 지정 문자열 지원.
	fmt.Printf("Name: %s, Age: %d", name, age) => Name: Bob, Age: 25
*/
/*
	%v: 구조체 출력 (필드명 포함하지 않음)
	%+v 구조체 출력 (필드명 포함)
	%q: 문자열을 따옴표로 둘러싼 문자열로 리턴
	문자열 내에 따옴표나 특수 문자가 있는 경우에도 올바르게 출력
		str := "Hello, World!"
		fmt.Printf("%q\n", str) => "Hello, World!"

		"Hello, "Gopher"!" => "Hello, \"Gopher\"!"
*/

func TestRun(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		return run(ctx)
	})
	in := "message"
	rsp, err := http.Get("http://localhost:18080/" + in)
	if err != nil {
		t.Errorf("failed to get: %+v", err)
	}
	// Body 필드는 io.ReadClose 인터페이스를 따르며, Close 수행하지 않을 경우, 메모리 누수 발생.
	defer rsp.Body.Close()
	got, err := io.ReadAll(rsp.Body)
	if err != nil {
		t.Fatalf("failed to read body: %v", err)
	}

	// HTTP 서버의 반환 값을 검증.
	want := fmt.Sprintf("Hello, %s!", in)
	if string(got) != want {
		t.Errorf("want %q, but ot %q", want, got)
	}

	// run 함수에 종료 알림을 전송한다.
	cancel()

	// run 함수의 반환값을 검증한다.
	if err := eg.Wait(); err != nil {
		t.Fatal(err)
	}
}
