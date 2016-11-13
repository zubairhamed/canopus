package canopus

/*
#cgo LDFLAGS: -L${SRCDIR}/openssl -lssl -lcrypto -ldl
#cgo CFLAGS: -g -Wno-deprecated -Wno-error -I${SRCDIR}/openssl/include

#include "dtls.h"

extern int go_session_bio_write(BIO* bio, char* buf, int num);
extern int go_session_bio_read(BIO* bio, char* buf, int num);
extern int go_session_bio_free(BIO* bio);
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
	"strings"
	"sync/atomic"
	"syscall"
	"time"
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

// Context used for creating a DTLS connection.
// where you configure global parameters
type DTLSCtx struct {
	ctx *C.SSL_CTX
}

func NewDTLSContext() *DTLSCtx {
	ctx := C.SSL_CTX_new(C.DTLSv1_2_client_method())
	if ctx == nil {
		panic("error creating SSL context")
	}

	C.set_proto_1_2(ctx)
	C.init_ctx(ctx)

	self := DTLSCtx{ctx}
	return &self
}

func (ctx *DTLSCtx) SetCipherList(ciphers string) bool {
	ret := int(C.SSL_CTX_set_cipher_list(ctx.ctx, C.CString(ciphers)))
	return ret == 1
}

type DTLSContext struct {
	sslCtx *C.SSL_CTX
}

func createSslContext() (dtlsCtx *DTLSContext, err error) {
	sslCtx := C.SSL_CTX_new(C.DTLSv1_2_server_method())
	if sslCtx == nil {
		err = errors.New("Error creating SSL context")
		return
	}

	fmt.Println("call >> init_server_ctx")
	C.init_server_ctx(sslCtx)

	ret := int(C.SSL_CTX_set_cipher_list(sslCtx, C.CString("PSK-AES256-CCM8:PSK-AES128-CCM8")))
	if ret != 1 {
		err = errors.New("Unable to set CipherList")
		return
	}

	dtlsCtx = &DTLSContext{
		sslCtx: sslCtx,
	}

	return
}

//export go_session_bio_read
func go_session_bio_read(bio *C.BIO, buf *C.char, num C.int) C.int {
	fmt.Println("cgo :-- go_session_bio_read")
	biodata := *(*string)(C.BIO_get_data(bio))
	// sess := sessions[*(*int32)(C.BIO_get_data(bio))]

	datas := strings.Split(biodata, ",")
	serverId := datas[0]
	addr := datas[1]
	server := getServer(serverId)
	session := server.GetSession(addr).(*DTLSServerSession)

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
	fmt.Println("cgo :-- go_session_bio_write")
	biodata := *(*string)(C.BIO_get_data(bio))
	datas := strings.Split(biodata, ",")
	serverId := datas[0]
	addr := datas[1]
	server := getServer(serverId)
	session := server.GetSession(addr)
	data := goSliceFromCString(buf, int(num))

	n, err := session.Write(data)
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
	fmt.Println("cgo :-- go_session_bio_free")
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
	biodata := *(*string)(C.BIO_get_data(bio))
	datas := strings.Split(biodata, ",")
	serverId := datas[0]

	server := getServer(serverId).(*DefaultCoapServer)
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
	fmt.Println("generate_cookie_callback")
	bio := C.SSL_get_rbio(ssl)
	biodata := *(*string)(C.BIO_get_data(bio))
	datas := strings.Split(biodata, ",")
	serverId := datas[0]
	addr := datas[1]

	server := getServer(serverId).(*DefaultCoapServer)

	mac := hmac.New(sha256.New, server.cookieSecret)
	mac.Write([]byte(addr))
	cookieValue := mac.Sum(nil)

	if len(cookieValue) >= int(*cookie_len) {
		fmt.Println("no enough cookie space (should not happen..)")
		return 0
	}

	data := goSliceFromUCString(cookie, int(*cookie_len))

	*cookie_len = C.uint(copy(data, cookieValue))
	return 1

}

//export verify_cookie_callback
func verify_cookie_callback(ssl *C.SSL, cookie *C.uchar, cookie_len C.uint) C.int {
	fmt.Println("verify_cookie_callback")
	bio := C.SSL_get_rbio(ssl)
	biodata := *(*string)(C.BIO_get_data(bio))
	datas := strings.Split(biodata, ",")
	serverId := datas[0]
	addr := datas[1]

	server := getServer(serverId).(*DefaultCoapServer)

	mac := hmac.New(sha256.New, server.cookieSecret)
	mac.Write([]byte(addr))
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

type DTLSClient struct {
	closed    bool
	connected int32 // connection handshake was done, atomic (0 false, 1 true)
	bio       *C.BIO
	ctx       *C.SSL_CTX
	ssl       *C.SSL
	conn      net.Conn
	pskId     *string
	psk       []byte
}

var nextId int32 = 0

var clients = make(map[int32]*DTLSClient)

// Create a DTLSClient implementing the net.Conn interface
func NewDTLSClient(dtlsCtx *DTLSCtx, conn net.Conn) *DTLSClient {
	ssl := C.SSL_new(dtlsCtx.ctx)

	id := atomic.AddInt32(&nextId, 1)

	self := DTLSClient{false, 0, C.BIO_new(C.BIO_go()), dtlsCtx.ctx, ssl, conn, nil, nil}
	clients[id] = &self

	C.SSL_set_bio(self.ssl, self.bio, self.bio)

	// the ID is used as link between the Go and C lang since sharing Go pointers is
	// so the C is going to own the pointer to the id value
	C.setGoClientId(self.bio, C.uint(id))
	return &self
}

func (c *DTLSClient) connect() error {
	ret := C.SSL_connect(c.ssl)
	if err := c.getError(ret); err != nil {
		return err
	}
	return nil
}

func (c *DTLSClient) SetPSK(identity string, psk []byte) {
	c.psk = psk
	c.pskId = &identity
}

func (c *DTLSClient) Read(b []byte) (n int, err error) {
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

func (c *DTLSClient) Write(b []byte) (int, error) {
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

func (c *DTLSClient) LocalAddr() net.Addr {
	return c.LocalAddr()
}

func (c *DTLSClient) RemoteAddr() net.Addr {
	return c.RemoteAddr()
}

func (c *DTLSClient) SetDeadline(t time.Time) error {
	return nil
}
func (c *DTLSClient) SetReadDeadline(t time.Time) error {
	return nil
}
func (c *DTLSClient) SetWriteDeadline(t time.Time) error {
	return nil
}

func (c *DTLSClient) Close() error {
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

func (c *DTLSClient) getError(ret C.int) error {
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

//export go_conn_bio_write
func go_conn_bio_write(bio *C.BIO, buf *C.char, num C.int) C.int {

	client := clients[*(*int32)(C.BIO_get_data(bio))]
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
	client := clients[*(*int32)(C.BIO_get_data(bio))]
	data := goSliceFromCString(buf, int(num))

	fmt.Println("CLIENT---", client)
	n, err := client.conn.Read(data)
	if err == nil {
		return C.int(n)
	}
	if err == io.EOF || err == io.ErrUnexpectedEOF {
		return 0
	}
	//We expect either a syscall error
	//or a netOp error wrapping a syscall error
	fmt.Println(err)
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
	client := clients[*(*int32)(C.BIO_get_data(bio))]
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
	client := clients[*(*int32)(C.BIO_get_data(bio))]

	if client.pskId == nil || client.psk == nil {
		return 0
	}

	if len(*client.pskId) >= int(max_identity_len) || len(client.psk) >= int(max_psk_len) {
		fmt.Println("PSKID or PSK too large")
		return 0
	}
	targetId := goSliceFromCString(identity, int(max_identity_len))
	copy(targetId, *client.pskId)
	targetPsk := goSliceFromCString(psk, int(max_psk_len))
	return C.uint(copy(targetPsk, client.psk))
}
