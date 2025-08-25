package mqtt

import (
	"log"
	"testing"

	"github.com/labstack/echo/v4"
	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/hooks/auth"
	"github.com/mochi-mqtt/server/v2/listeners"
)

func Test_mochi(t *testing.T) {
	// MQTT 서버 초기화
	server := mqtt.New(&mqtt.Options{
		InlineClient: true, // 인라인 클라이언트 활성화 (선택사항)
	})

	// 인증 설정 (모든 연결 허용, 필요 시 커스텀 인증 추가)
	_ = server.AddHook(new(auth.AllowHook), nil)

	// TCP 리스너 추가 (기본 MQTT 포트: 1883)
	tcp := listeners.NewTCP(listeners.Config{ ID: "t1", Address: ":1883" })
	err := server.AddListener(tcp)
	if err != nil {
		log.Fatal("Failed to add TCP listener:", err)
	}

	// WebSocket 리스너 추가 (포트는 HTTP 서버에서 관리)
	ws := listeners.NewWebsocket(listeners.Config{ ID: "t2", Address: ":1884" })
	err = server.AddListener(ws)
	if err != nil {
		log.Fatal("Failed to add WebSocket listener:", err)
	}

	e := echo.New()

	e.GET("/ws", func(c echo.Context) error {
		ws.Handler(c.Response(), c.Request())
		return	nil
	})

	go func() {
		server.Serve()
	}()

	// MQTT 서버 시작
	log.Println("Starting MQTT server")
	e.Logger.Fatal(e.Start(":8002"))
}