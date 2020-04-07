package pubsub

import (
	"github.com/crosstalkio/log"
	api "github.com/crosstalkio/pubsub/api/pubsub"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

type clientWriteJob struct {
	msg   *api.Message
	resCh chan error
}

type clientWriter struct {
	log.Sugar
	conn   *websocket.Conn
	jobs   chan *clientWriteJob
	exitCh chan bool
}

func newClientWriter(logger log.Logger, conn *websocket.Conn) *clientWriter {
	return &clientWriter{
		Sugar:  log.NewSugar(logger),
		conn:   conn,
		jobs:   make(chan *clientWriteJob),
		exitCh: make(chan bool, 1),
	}
}

func (w *clientWriter) exit() {
	w.exitCh <- true
}

func (w *clientWriter) write(msg *api.Message) error {
	job := &clientWriteJob{
		msg:   msg,
		resCh: make(chan error, 1),
	}
	w.jobs <- job
	return <-job.resCh
}

func (w *clientWriter) loop() {
	defer w.Debugf("Exiting writer: %s", w.conn.LocalAddr().String())
	for {
		select {
		case <-w.exitCh:
			return
		case job := <-w.jobs:
			data, err := proto.Marshal(job.msg)
			if err != nil {
				w.Errorf("Failed to marshal proto: %s", err.Error())
				job.resCh <- err
				break
			}
			err = w.conn.WriteMessage(websocket.BinaryMessage, data)
			if err != nil {
				w.Errorf("Failed to write proto: %s", err.Error())
				job.resCh <- err
				break
			}
			job.resCh <- nil
		}
	}
}
