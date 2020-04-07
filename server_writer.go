package pubsub

import (
	"time"

	"github.com/crosstalkio/log"
	"github.com/gorilla/websocket"
)

const (
	Deadline   = 5 * time.Second
	PingPeriod = 30 * time.Second
)

type serverWriteJob struct {
	messageType int
	data        []byte
}

type serverWriter struct {
	log.Sugar
	conn   *websocket.Conn
	ticker *time.Ticker
	jobCh  chan *serverWriteJob
	exitCh chan bool
}

func newServerWriter(logger log.Logger, conn *websocket.Conn) *serverWriter {
	return &serverWriter{
		Sugar:  log.NewSugar(logger),
		ticker: time.NewTicker(PingPeriod),
		conn:   conn,
		jobCh:  make(chan *serverWriteJob),
		exitCh: make(chan bool, 1),
	}
}

func (w *serverWriter) exit() {
	w.exitCh <- true
}

func (w *serverWriter) loop() {
	defer func() {
		w.ticker.Stop()
		w.conn.Close()
		w.Debugf("Exiting writer: %s", w.conn.RemoteAddr().String())
	}()
	for {
		select {
		case <-w.exitCh:
			return
		case <-w.ticker.C:
			w.Debugf("Pinging: %s", w.conn.RemoteAddr().String())
			err := w.conn.WriteControl(websocket.PingMessage, nil, time.Now().Add(Deadline))
			if err != nil {
				w.Errorf("Failed to ping %s: %s", w.conn.RemoteAddr().String(), err.Error())
				return
			}
		case job := <-w.jobCh:
			err := w.conn.SetWriteDeadline(time.Now().Add(Deadline))
			if err != nil {
				w.Errorf("Failed to set websocket write deadline: %s", err.Error())
				return
			}
			w.Debugf("Writing %d bytes websocket message: %s", len(job.data), w.conn.RemoteAddr().String())
			err = w.conn.WriteMessage(job.messageType, job.data)
			if err != nil {
				w.Errorf("Failed to write websocket message: %s", err.Error())
				return
			}
			w.Debugf("Written %d bytes websocket message: %s", len(job.data), w.conn.RemoteAddr().String())
		}
	}
}

func (w *serverWriter) write(messageType int, data []byte) {
	req := &serverWriteJob{
		messageType: messageType,
		data:        data,
	}
	w.jobCh <- req
}
