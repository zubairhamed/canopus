package canopus

/*
#cgo LDFLAGS: -L${SRCDIR}/openssl -lssl -lcrypto -ldl
#cgo CFLAGS: -g -Wno-deprecated -Wno-error -I${SRCDIR}/openssl/include

#include <stdlib.h>
#include <string.h>
#include <openssl/err.h>
#include <openssl/ssl.h>
#include <openssl/bio.h>
#include <internal/bio.h>

extern int go_session_bio_write(BIO* bio, char* buf, int num);
extern int go_session_bio_read(BIO* bio, char* buf, int num);
extern int go_session_bio_free(BIO* bio);
extern unsigned int go_server_psk_callback(SSL *ssl, char *identity, char *psk, unsigned int max_psk_len);
extern int generate_cookie_callback(SSL* ssl, unsigned char* cookie, unsigned int *cookie_len);
extern int verify_cookie_callback(SSL* ssl, unsigned char* cookie, unsigned int cookie_len);

static long go_session_bio_ctrl(BIO *bp,int cmd,long larg,void *parg) {
	return 1;
}

static int write_wrapper(BIO* bio, char* data, int n) {
	return go_session_bio_write(bio,data,n);
}

static int go_session_bio_create( BIO *b ) {
	BIO_set_init(b,1);
	BIO_set_flags(b, BIO_FLAGS_READ | BIO_FLAGS_WRITE);
	return 1;
}

// a BIO for a client conencted to our server
static BIO_METHOD* go_session_bio_method;

static BIO_METHOD* BIO_go_session() {
	return go_session_bio_method;
}

static void set_errno(int e) {
	errno = e;
}

static char *getGoData(BIO* bio) {
	return BIO_get_data(bio);
}

static unsigned int server_psk_callback(SSL *ssl, char *identity, unsigned char *psk, unsigned int max_psk_len) {
	return go_server_psk_callback(ssl,identity,(char*)psk,max_psk_len);
}

static void init_lib() {
	setvbuf(stdout, NULL, _IOLBF, 0);
	SSL_library_init();
	ERR_load_BIO_strings();
	SSL_load_error_strings();
}

static int init_session_bio_method() {
	go_session_bio_method = BIO_meth_new(BIO_TYPE_SOURCE_SINK,"go session dtls");
	BIO_meth_set_write(go_session_bio_method,write_wrapper);
	BIO_meth_set_read(go_session_bio_method,go_session_bio_read);
	BIO_meth_set_ctrl(go_session_bio_method,go_session_bio_ctrl);
	BIO_meth_set_create(go_session_bio_method,go_session_bio_create);
	BIO_meth_set_destroy(go_session_bio_method,go_session_bio_free);
}

static void init_server_ctx(SSL_CTX *ctx) {
	SSL_CTX_set_min_proto_version(ctx, 0xFEFD); // 1.2
	SSL_CTX_set_max_proto_version(ctx, 0xFEFD); // 1.2
	SSL_CTX_set_read_ahead(ctx, 1);
	SSL_CTX_set_cookie_generate_cb(ctx, &generate_cookie_callback);
	SSL_CTX_set_cookie_verify_cb(ctx, &verify_cookie_callback);
}

static int get_errno(void) {
	return errno;
}

static void setGoData(BIO* bio, char *data) {
	BIO_set_data(bio, data);
}

static void set_cookie_option(SSL *ssl) {
	SSL_set_options(ssl, SSL_OP_COOKIE_EXCHANGE);
}

static void set_psk_callback(SSL *ssl) {
	SSL_set_psk_server_callback(ssl, &server_psk_callback);
}

static void setGoSessionId(BIO* bio, unsigned int clientId) {
	unsigned int * pId = malloc(sizeof(unsigned int));
	*pId = clientId;
	BIO_set_data(bio,pId);
}

// Client
extern int go_conn_bio_write(BIO* bio, char* buf, int num);
extern int go_conn_bio_read(BIO* bio, char* buf, int num);
extern int go_conn_bio_free(BIO* bio);
extern unsigned int go_psk_callback(SSL *ssl, char *hint, char *identity, unsigned int max_identity_len, char *psk, unsigned int max_psk_len);

static long go_bio_ctrl(BIO *bp,int cmd,long larg,void *parg) {
	//always return operation not supported
	//http://www.openssl.org/docs/crypto/BIO_ctrl.html
	//printf("go_bio_ctrl %d\n", cmd);
	return 1;
}

static int go_bio_create( BIO *b ) {
	BIO_set_init(b,1);
	//BIO_set_num(b,-1);
	//BIO_set_ptr(b,NULL);
	BIO_set_flags(b, BIO_FLAGS_READ | BIO_FLAGS_WRITE);
	return 1;
}

static BIO_METHOD go_bio_method = {
	BIO_TYPE_SOURCE_SINK,
	"go dtls",
	(int (*)(BIO *, const char *, int))go_conn_bio_write,
	go_conn_bio_read,
	NULL,
	NULL,
	go_bio_ctrl, // ctrl
	go_bio_create, // new
	go_conn_bio_free // delete
};

static BIO_METHOD* BIO_go() {
	return &go_bio_method;
}

static void set_proto_1_2(SSL_CTX *ctx) {
	SSL_CTX_set_min_proto_version(ctx, 0xFEFD); // 1.2
	SSL_CTX_set_max_proto_version(ctx, 0xFEFD); // 1.2
}

static unsigned int psk_callback(SSL *ssl, const char *hint,
        char *identity, unsigned int max_identity_len,
        unsigned char *psk, unsigned int max_psk_len) {
	return go_psk_callback(ssl,hint,identity,max_identity_len,(char*)psk,max_psk_len);
}

static void init_ctx(SSL_CTX *ctx) {
	SSL_CTX_set_read_ahead(ctx, 1);
	SSL_CTX_set_psk_client_callback(ctx,&psk_callback);
}

static void setGoClientId(BIO* bio, unsigned int clientId) {
	unsigned int * pId = malloc(sizeof(unsigned int));
	*pId = clientId;
	BIO_set_data(bio,pId);
}

*/
import "C"
import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"net"
	"reflect"
	"sync"
	"sync/atomic"
	"syscall"
	"unsafe"
)

func init() {
	// low level init of OpenSSL
	C.init_lib()

	// init server BIO
	server_bio_method_init()
}

func server_bio_method_init() {
	C.init_session_bio_method()
}

func NewServerDtlsContext() (ctx *ServerDtlsContext, err error) {
	sslCtx := C.SSL_CTX_new(C.DTLSv1_2_server_method())

	if sslCtx == nil {
		err = errors.New("error creating SSL context")
		return
	}

	C.init_server_ctx(sslCtx)

	ret := int(C.SSL_CTX_set_cipher_list(sslCtx, C.CString("PSK-AES256-CCM8:PSK-AES128-CCM8")))
	if ret != 1 {
		err = errors.New("impossible to set cipherlist")
	}

	ctx = &ServerDtlsContext{
		sslCtx: sslCtx,
	}

	return
}

type ServerDtlsContext struct {
	sslCtx *C.SSL_CTX
}

//export go_session_bio_read
func go_session_bio_read(bio *C.BIO, buf *C.char, num C.int) C.int {
	session := DTLS_SERVER_SESSIONS[*(*int32)(C.BIO_get_data(bio))]
	socketData := <-session.rcvd

	data := goSliceFromCString(buf, int(num))
	if data == nil {
		return 0
	}

	wrote := copy(data, socketData)
	return C.int(wrote)
}

//export go_session_bio_write
func go_session_bio_write(bio *C.BIO, buf *C.char, num C.int) C.int {
	session := DTLS_SERVER_SESSIONS[*(*int32)(C.BIO_get_data(bio))]
	data := goSliceFromCString(buf, int(num))

	n, err := session.GetConnection().WriteTo(data, session.GetAddress())
	if err != nil && err != io.EOF {
		//We expect either a syscall error
		//or a netOp error wrapping a syscall error
	TESTERR:
		switch err.(type) {
		case syscall.Errno:
			C.set_errno(C.int(err.(syscall.Errno)))
		case *net.OpError:
			err = err.(*net.OpError).Err
			break TESTERR
		}
		return C.int(-1)
	}
	return C.int(n)
}

//export go_session_bio_free
func go_session_bio_free(bio *C.BIO) C.int {
	// TODO

	// some flags magic
	if C.int(C.BIO_get_shutdown(bio)) != 0 {
		C.BIO_set_data(bio, nil)
		C.BIO_set_flags(bio, 0)
		C.BIO_set_init(bio, 0)
	}
	return C.int(1)
}

//export go_server_psk_callback
func go_server_psk_callback(ssl *C.SSL, identity *C.char, psk *C.char, max_psk_len C.uint) C.uint {
	bio := C.SSL_get_rbio(ssl)
	session := DTLS_SERVER_SESSIONS[*(*int32)(C.BIO_get_data(bio))]
	server := session.GetServer().(*DefaultCoapServer)

	goPskID := C.GoString(identity)

	serverPsk := server.fnPskHandler(goPskID)

	if serverPsk == nil {
		return 0
	}

	if len(serverPsk) >= int(max_psk_len) {
		return 0
	}

	targetPsk := goSliceFromCString(psk, int(max_psk_len))
	return C.uint(copy(targetPsk, serverPsk))
}

//export generate_cookie_callback
func generate_cookie_callback(ssl *C.SSL, cookie *C.uchar, cookie_len *C.uint) C.int {
	bio := C.SSL_get_rbio(ssl)
	session := DTLS_SERVER_SESSIONS[*(*int32)(C.BIO_get_data(bio))]

	mac := hmac.New(sha256.New, session.GetServer().GetCookieSecret())
	mac.Write([]byte(session.GetAddress().String()))
	cookieValue := mac.Sum(nil)

	if len(cookieValue) >= int(*cookie_len) {
		logMsg("Not enough cookie space (should not happen..)")
		return 0
	}

	data := goSliceFromUCString(cookie, int(*cookie_len))

	*cookie_len = C.uint(copy(data, cookieValue))
	return 1

}

//export verify_cookie_callback
func verify_cookie_callback(ssl *C.SSL, cookie *C.uchar, cookie_len C.uint) C.int {
	bio := C.SSL_get_rbio(ssl)
	session := DTLS_SERVER_SESSIONS[*(*int32)(C.BIO_get_data(bio))]

	mac := hmac.New(sha256.New, session.GetServer().GetCookieSecret())
	mac.Write([]byte(session.GetAddress().String()))
	cookieValue := mac.Sum(nil)

	if len(cookieValue) != int(cookie_len) {
		return 0
	}

	data := goSliceFromUCString(cookie, int(cookie_len))

	if bytes.Equal(data, cookieValue) {
		return 1
	} else {
		return 0
	}
}

// Provides a zero copy interface for returning a go slice backed by a c array.
func goSliceFromCString(cArray *C.char, size int) (cslice []byte) {
	//See http://code.google.com/p/go-wiki/wiki/cgo
	//It turns out it's really easy to
	//make a string from a *C.char and vise versa.
	//not so easy to write to a c array.
	sliceHeader := (*reflect.SliceHeader)((unsafe.Pointer(&cslice)))
	sliceHeader.Cap = size
	sliceHeader.Len = size
	sliceHeader.Data = uintptr(unsafe.Pointer(cArray))
	return
}

// Provides a zero copy interface for returning a go slice backed by a c array.
func goSliceFromUCString(cArray *C.uchar, size int) (cslice []byte) {
	//See http://code.google.com/p/go-wiki/wiki/cgo
	//It turns out it's really easy to
	//make a string from a *C.char and vise versa.
	//not so easy to write to a c array.
	sliceHeader := (*reflect.SliceHeader)((unsafe.Pointer(&cslice)))
	sliceHeader.Cap = size
	sliceHeader.Len = size
	sliceHeader.Data = uintptr(unsafe.Pointer(cArray))
	return
}

func getErrorString(code C.ulong) string {
	if code == 0 {
		return ""
	}
	msg := fmt.Sprintf("%s:%s:%s\n",
		C.GoString(C.ERR_lib_error_string(code)),
		C.GoString(C.ERR_func_error_string(code)),
		C.GoString(C.ERR_reason_error_string(code)))
	if len(msg) == 4 { //being lazy here, all the strings were empty
		return ""
	}
	//Check for extra line data
	var file *C.char
	var line C.int
	var data *C.char
	var flags C.int
	if int(C.ERR_get_error_line_data(&file, &line, &data, &flags)) != 0 {
		msg += fmt.Sprintf("%s:%s", C.GoString(file), int(line))
		if flags&C.ERR_TXT_STRING != 0 {
			msg += ":" + C.GoString(data)
		}
		if flags&C.ERR_TXT_MALLOCED != 0 {
			C.CRYPTO_free(unsafe.Pointer(data), C.CString(""), 0)
		}
	}
	return msg
}

func newSslSession(session *DTLSServerSession, ctx *ServerDtlsContext, pskCallback func(id string) []byte) (err error) {
	ssl := C.SSL_new(ctx.sslCtx)

	id := atomic.AddInt32(&NEXT_SESSION_ID, 1)

	if pskCallback != nil {
		C.set_psk_callback(ssl)
	}

	bio := C.BIO_new(C.BIO_go_session())

	if bio == nil {
		err = errors.New("Error creating session: Bio is nil")
		return
	}
	C.SSL_set_bio(ssl, bio, bio)

	session.ssl = ssl
	session.bio = bio

	DTLS_SERVER_SESSIONS[id] = session

	C.setGoSessionId(bio, C.uint(id))

	C.set_cookie_option(ssl)
	C.SSL_set_accept_state(ssl)
	C.DTLSv1_listen

	return
}

type DTLSServerSession struct {
	UDPServerSession
	ssl *C.SSL
	bio *C.BIO
}

func (s *DTLSServerSession) GetConnection() ServerConnection {
	return s.conn
}

func (s *DTLSServerSession) Write(b []byte) (int, error) {
	// TODO test is connected ?
	length := len(b)
	ret := C.SSL_write(s.ssl, unsafe.Pointer(&b[0]), C.int(length))
	if err := s.getError(ret); err != nil {
		return 0, err
	}
	return int(ret), nil
}

func (s *DTLSServerSession) Read(b []byte) (n int, err error) {
	// TODO test if closed?
	length := len(b)
	// s.rcvd <- s.buf

	ret := C.SSL_read(s.ssl, unsafe.Pointer(&b[0]), C.int(length))
	if err = s.getError(ret); err != nil {
		n = 0
		return
	}
	// if there's no error, but a return value of 0
	// let's say it's an EOF
	if ret == 0 {
		n = 0
		err = io.EOF
		return
	}
	n = int(ret)
	return
}

func (s *DTLSServerSession) getError(ret C.int) error {
	err := C.SSL_get_error(s.ssl, ret)
	switch err {
	case C.SSL_ERROR_NONE:
		return nil
	case C.SSL_ERROR_ZERO_RETURN:
		return io.EOF
	case C.SSL_ERROR_SYSCALL:
		if int(C.ERR_peek_error()) != 0 {
			return syscall.Errno(C.get_errno())
		}

	default:
		msg := ""
		for {
			errCode := C.ERR_get_error()
			if errCode == 0 {
				break
			}
			msg += getErrorString(errCode)
		}
		C.ERR_clear_error()
		return errors.New(msg)
	}
	return nil
}

// Client DTLS
func NewDTLSConnection(c net.Conn, identity, psk string) (conn Connection, err error) {
	sslCtx := C.SSL_CTX_new(C.DTLSv1_2_client_method())

	C.set_proto_1_2(sslCtx)
	C.init_ctx(sslCtx)

	ret := int(C.SSL_CTX_set_cipher_list(sslCtx, C.CString("PSK-AES256-CCM8:PSK-AES128-CCM8")))
	if ret != 1 {
		err = errors.New("impossible to set cipherlist")
		return
	}

	ssl := C.SSL_new(sslCtx)

	// self := DTLSClient{false, 0, C.BIO_new(C.BIO_go()), dtlsCtx.ctx, ssl, conn, nil, nil}
	bio := C.BIO_new(C.BIO_go())

	conn = &DTLSConnection{
		UDPConnection: UDPConnection{
			conn: c,
		},
		sslCtx: sslCtx,
		ssl:    ssl,
		bio:    bio,
		psk:    []byte(psk),
		pskId:  &identity,
	}

	C.SSL_set_bio(ssl, bio, bio)

	id := atomic.AddInt32(&NEXT_SESSION_ID, 1)
	C.setGoClientId(bio, C.uint(id))
	DTLS_CLIENT_CONNECTIONS[id] = conn.(*DTLSConnection)

	return
}

type DTLSConnection struct {
	UDPConnection
	closed    bool
	connected int32 // connection handshake was done, atomic (0 false, 1 true)
	sslCtx    *C.SSL_CTX
	bio       *C.BIO
	ssl       *C.SSL
	pskId     *string
	psk       []byte
}

func (c *DTLSConnection) ObserveResource(resource string) (tok string, err error) {
	req := NewRequest(MessageConfirmable, Get, GenerateMessageID())
	req.SetRequestURI(resource)
	req.GetMessage().AddOption(OptionObserve, 0)

	resp, err := c.Send(req)
	tok = string(resp.GetMessage().GetToken())

	return
}

func (c *DTLSConnection) CancelObserveResource(resource string, token string) (err error) {
	req := NewRequest(MessageConfirmable, Get, GenerateMessageID())
	req.SetRequestURI(resource)
	req.GetMessage().AddOption(OptionObserve, 1)

	_, err = c.Send(req)
	return
}

func (c *DTLSConnection) StopObserve(ch chan ObserveMessage) {
	close(ch)
}

func (c *DTLSConnection) Observe(ch chan ObserveMessage) {

	readBuf := make([]byte, MaxPacketSize)
	for {
		len, err := c.Read(readBuf)
		if err == nil {
			msgBuf := make([]byte, len)
			copy(msgBuf, readBuf)

			msg, err := BytesToMessage(msgBuf)
			if msg.GetOption(OptionObserve) != nil {
				ch <- NewObserveMessage(msg.GetURIPath(), msg.GetPayload(), msg)
			}
			if err != nil {
				logMsg("Error occured reading UDP", err)
				close(ch)
			}
		} else {
			logMsg("Error occured reading UDP", err)
			close(ch)
		}
	}
}

func (c *DTLSConnection) Send(req Request) (resp Response, err error) {
	msg := req.GetMessage()
	opt := msg.GetOption(OptionBlock1)

	if opt == nil { // Block1 was not set
		if MessageSizeAllowed(req) != true {
			return nil, ErrMessageSizeTooLongBlockOptionValNotSet
		}

	} else { // Block1 was set
		// fmt.Println("Block 1 was set")
	}

	if opt != nil {
		blockOpt := Block1OptionFromOption(opt)

		if blockOpt.Value == nil {
			if MessageSizeAllowed(req) != true {
				err = ErrMessageSizeTooLongBlockOptionValNotSet
				return
			} else {
				// - Block # = one and only block (sz = unspecified), whereas 0 = 16bits
				// - MOre bit = 0
			}
		} else { // BLock transfer request
			payload := msg.GetPayload().GetBytes()
			payloadLen := uint32(len(payload))
			blockSize := blockOpt.BlockSizeLength()
			currSeq := uint32(0)
			totalBlocks := uint32(payloadLen / blockSize)
			completed := false

			var wg sync.WaitGroup
			wg.Add(1)

			for completed == false {
				if currSeq <= totalBlocks {

					var blockPayloadStart uint32
					var blockPayloadEnd uint32
					var blockPayload []byte

					blockPayloadStart = currSeq*uint32(blockSize) + (currSeq * 1)

					more := true
					if currSeq == totalBlocks {
						more = false
						blockPayloadEnd = payloadLen
					} else {
						blockPayloadEnd = blockPayloadStart + uint32(blockSize)
					}

					blockPayload = payload[blockPayloadStart:blockPayloadEnd]

					blockOpt = NewBlock1Option(blockOpt.Size(), more, currSeq)
					msg.ReplaceOptions(blockOpt.Code, []Option{blockOpt})
					modifiedMsg := msg.(*CoapMessage)
					modifiedMsg.SetMessageId(GenerateMessageID())
					modifiedMsg.SetPayload(NewBytesPayload(blockPayload))

					// send message
					_, err2 := c.sendMessage(msg)
					if err2 != nil {
						wg.Done()
						return
					}
					currSeq = currSeq + 1

				} else {
					completed = true
					wg.Done()
				}
			}
		}
	}
	resp, err = c.sendMessage(msg)
	return
}

func (c *DTLSConnection) sendMessage(msg Message) (resp Response, err error) {

	if msg == nil {
		return nil, ErrNilMessage
	}

	b, err := MessageToBytes(msg)
	if err != nil {
		return
	}

	if msg.GetMessageType() == MessageNonConfirmable {
		go c.Write(b)
		resp = NewResponse(NewEmptyMessage(msg.GetMessageId()), nil)
		return
	}

	_, err = c.Write(b)
	if err != nil {
		return
	}

	msgBuf := make([]byte, 1500)
	if msg.GetMessageType() == MessageAcknowledgment {
		resp = NewResponse(NewEmptyMessage(msg.GetMessageId()), nil)
		return
	}

	n, err := c.Read(msgBuf)
	if err != nil {
		return
	}

	respMsg, err := BytesToMessage(msgBuf[:n])
	if err != nil {
		return
	}

	resp = NewResponse(respMsg, nil)

	if msg.GetMessageType() == MessageConfirmable {
		// TODO: Send out message and wait for a confirm. If confirmation not retrieved,
		// resend (taking into account timeouts and back-off transmissions

		// c.Send(NewRequestFromMessage(msg))
	}
	return
}

func (c *DTLSConnection) Write(b []byte) (int, error) {
	if atomic.CompareAndSwapInt32(&c.connected, 0, 1) {
		if err := c.connect(); err != nil {
			return 0, err
		}
	}
	length := len(b)
	ret := C.SSL_write(c.ssl, unsafe.Pointer(&b[0]), C.int(length))
	if err := c.getError(ret); err != nil {
		return 0, err
	}

	return int(ret), nil
}

func (c *DTLSConnection) Read(b []byte) (int, error) {
	if atomic.CompareAndSwapInt32(&c.connected, 0, 1) {
		if err := c.connect(); err != nil {
			return 0, err
		}
	}

	length := len(b)
	ret := C.SSL_read(c.ssl, unsafe.Pointer(&b[0]), C.int(length))
	if err := c.getError(ret); err != nil {
		return 0, err
	}

	// if there's no error, but a return value of 0
	// let's say it's an EOF
	if ret == 0 {
		return 0, io.EOF
	}

	return int(ret), nil
}

func (c *DTLSConnection) Close() error {
	if c.closed {
		return nil
	}
	c.closed = true
	defer func() {
		C.SSL_free(c.ssl)
	}()

	ret := C.SSL_shutdown(c.ssl)
	if int(ret) == 0 {
		ret = C.SSL_shutdown(c.ssl)
		if int(ret) != 1 {
			return c.getError(ret)
		}
	}
	return nil
}

func (c *DTLSConnection) connect() error {
	ret := C.SSL_connect(c.ssl)
	if err := c.getError(ret); err != nil {
		return err
	}
	return nil
}

func (c *DTLSConnection) getError(ret C.int) error {
	err := C.SSL_get_error(c.ssl, ret)
	switch err {
	case C.SSL_ERROR_NONE:
		return nil
	case C.SSL_ERROR_ZERO_RETURN:
		return io.EOF
	case C.SSL_ERROR_SYSCALL:
		if int(C.ERR_peek_error()) != 0 {
			return syscall.Errno(C.get_errno())
		}

	default:
		msg := ""
		for {
			errCode := C.ERR_get_error()
			if errCode == 0 {
				break
			}
			msg += getErrorString(errCode)
		}
		C.ERR_clear_error()
		return errors.New(msg)
	}
	return nil
}

//export go_conn_bio_write
func go_conn_bio_write(bio *C.BIO, buf *C.char, num C.int) C.int {
	client := DTLS_CLIENT_CONNECTIONS[*(*int32)(C.BIO_get_data(bio))]
	data := goSliceFromCString(buf, int(num))
	n, err := client.conn.Write(data)
	if err != nil && err != io.EOF {
		//We expect either a syscall error
		//or a netOp error wrapping a syscall error
	TESTERR:
		switch err.(type) {
		case syscall.Errno:
			C.set_errno(C.int(err.(syscall.Errno)))
		case *net.OpError:
			err = err.(*net.OpError).Err
			break TESTERR
		}
		return C.int(-1)
	}
	return C.int(n)
}

//export go_conn_bio_read
func go_conn_bio_read(bio *C.BIO, buf *C.char, num C.int) C.int {
	client := DTLS_CLIENT_CONNECTIONS[*(*int32)(C.BIO_get_data(bio))]
	data := goSliceFromCString(buf, int(num))
	n, err := client.conn.Read(data)
	if err == nil {
		return C.int(n)
	}

	if err == io.EOF || err == io.ErrUnexpectedEOF {
		return 0
	}
	//We expect either a syscall error
	//or a netOp error wrapping a syscall error

TESTERR:
	switch err.(type) {
	case syscall.Errno:
		C.set_errno(C.int(err.(syscall.Errno)))
	case *net.OpError:
		err = err.(*net.OpError).Err
		break TESTERR
	}
	return C.int(-1)
}

//export go_conn_bio_free
func go_conn_bio_free(bio *C.BIO) C.int {
	client := DTLS_CLIENT_CONNECTIONS[*(*int32)(C.BIO_get_data(bio))]
	client.Close()
	if C.int(C.BIO_get_shutdown(bio)) != 0 {
		C.BIO_set_data(bio, nil)
		C.BIO_set_flags(bio, 0)
		C.BIO_set_init(bio, 0)
	}
	return C.int(1)
}

//export go_psk_callback
func go_psk_callback(ssl *C.SSL, hint *C.char, identity *C.char, max_identity_len C.uint, psk *C.char, max_psk_len C.uint) C.uint {
	bio := C.SSL_get_rbio(ssl)
	client := DTLS_CLIENT_CONNECTIONS[*(*int32)(C.BIO_get_data(bio))]

	if client.pskId == nil || client.psk == nil {
		return 0
	}

	if len(*client.pskId) >= int(max_identity_len) || len(client.psk) >= int(max_psk_len) {
		logMsg("PSKID or PSK too large")
		return 0
	}
	targetId := goSliceFromCString(identity, int(max_identity_len))
	copy(targetId, *client.pskId)
	targetPsk := goSliceFromCString(psk, int(max_psk_len))
	return C.uint(copy(targetPsk, client.psk))
}
