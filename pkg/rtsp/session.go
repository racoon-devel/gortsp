package rtsp

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
)

// Session is a low-level representation of RTSP session
type Session struct {
	ctx    context.Context
	cancel context.CancelFunc

	// channel for requests
	reqCh chan *request
	// channel for receiving items such as packets, requests, etc
	recvCh chan interface{}
	// receiving items from connection
	readCh chan interface{}

	seq  uint64
	creq map[uint64]*request
	wg   sync.WaitGroup
}

// NewSession creates new session
func NewSession(conn net.Conn, ctx context.Context) *Session {
	s := &Session{
		reqCh:  make(chan *request),
		recvCh: make(chan interface{}, incomingItemsCapacity),
		readCh: make(chan interface{}, incomingItemsCapacity),
		creq:   map[uint64]*request{},
	}

	s.ctx, s.cancel = context.WithCancel(ctx)

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		s.eventsProcess(conn)
	}()

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		s.parseProcess(conn)
	}()

	return s
}

// Do performs request and wait response
func (s *Session) Do(r *Request) (*Response, error) {
	in := &request{
		req:  r,
		resp: make(chan interface{}),
	}
	s.reqCh <- in

	out, _ := <-in.resp
	switch t := out.(type) {
	case *Response:
		return t, nil
	case error:
		return nil, t
	default:
		panic("unexpected type")
	}
}

// Incoming gets channel which can forward any of item:
// 1) *IncomingRTP - incoming RTP packet
// 2) *IncomingRTCP - incoming RTCP packet
// 3) *Request - incoming RTSP request
// 4) error - in case when error occurs
func (s *Session) Incoming() <-chan interface{} {
	return s.recvCh
}

func (s *Session) Close() {
	s.cancel()
	s.wg.Wait()
}

func (s *Session) eventsProcess(conn net.Conn) {
	var err error
	for err == nil {
		select {
		case data := <-s.readCh:
			err = s.processIncoming(data)
		case req := <-s.reqCh:
			err = s.sendRequest(conn, req)
		case <-s.ctx.Done():
			err = s.ctx.Err()
		}
	}

	s.recvCh <- err

	_ = conn.Close()
	close(s.reqCh)
	close(s.recvCh)
	close(s.readCh)
	for _, r := range s.creq {
		r.resp <- err
		close(r.resp)
	}
}

func (s *Session) sendRequest(conn net.Conn, req *request) error {
	// set sequence number
	s.seq++
	req.seq = s.seq
	req.req.Header.Add("Cseq", fmt.Sprintf("%d", s.seq))

	// serialize and send request
	if err := req.req.Write(conn); err != nil {
		return err
	}

	s.creq[req.seq] = req
	return nil
}

// reads all incoming messages
func (s *Session) parseProcess(conn net.Conn) {
	r := bufio.NewReader(conn)
	for {
		b, err := r.Peek(4)
		if err != nil {
			s.readCh <- fmt.Errorf("receive RTSP data failed: %w", err)
			return
		}
		switch {
		case b[0] == MagicSymbol: // parse interleaved packet
			h := InterleavedHeader{}
			if err = h.Read(r); err != nil {
				s.readCh <- fmt.Errorf("read interleaved header failed: %w", err)
				return
			}
			// todo: mempool
			buf := make([]byte, h.Length)
			if _, err = io.ReadFull(r, buf); err != nil {
				s.readCh <- fmt.Errorf("read packet failed: %w", err)
				return
			}
			if h.Channel%2 == 0 {
				s.readCh <- &IncomingRTP{
					Channel: h.Channel,
					Packet:  buf,
				}
			} else {
				s.readCh <- &IncomingRTCP{
					Channel: h.Channel,
					Packet:  buf,
				}
			}

		case b[0] == 'R' && b[1] == 'T' && b[2] == 'S' && b[3] == 'P': // parse response
			var resp Response
			if err = resp.Read(r); err != nil {
				s.readCh <- fmt.Errorf("read RTSP response failed: %w", err)
				return
			}
			s.readCh <- &resp

		case b[0] >= 'A' && b[0] <= 'Z': // parse request
			var req Request
			if err = req.Read(r); err != nil {
				s.readCh <- fmt.Errorf("read RTSP request failed: %w", err)
				return
			}
			s.readCh <- &req

		default:
			s.readCh <- errors.New("parse RTSP stream failed")
			return
		}
	}
}

func (s *Session) processIncoming(data interface{}) error {
	switch t := data.(type) {
	case *Response:
		seq, err := t.Seq()
		if err != nil {
			return fmt.Errorf("CSeq header is not presented: %w", err)
		}

		req, ok := s.creq[seq]
		if !ok {
			return fmt.Errorf("unknown response: CSeq = %d", seq)
		}

		req.resp <- t
		delete(s.creq, seq)
	case *IncomingRTP:
		s.recvCh <- t
	case *IncomingRTCP:
		s.recvCh <- t
	case *Request:
		s.recvCh <- t
	case error:
		return t
	}

	return nil
}
