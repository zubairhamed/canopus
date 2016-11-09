package canopus

/*
#cgo LDFLAGS: -L${SRCDIR}/openssl -lssl -lcrypto -ldl
#cgo CFLAGS: -g -Wno-deprecated -Wno-error -I${SRCDIR}/openssl/include

#include <stdlib.h>
#include <string.h>
#include <openssl/err.h>
#include <openssl/ssl.h>
#include <openssl/bio.h>


extern int go_session_bio_write(BIO* bio, char* buf, int num);
extern int go_session_bio_read(BIO* bio, char* buf, int num);
extern int go_session_bio_free(BIO* bio);

extern unsigned int go_server_psk_callback(SSL *ssl, char *identity, char *psk, unsigned int max_psk_len);

extern int generate_cookie_callback(SSL* ssl, unsigned char* cookie, unsigned int *cookie_len);
extern int verify_cookie_callback(SSL* ssl, unsigned char* cookie, unsigned int cookie_len);

extern int get_errno(void);
extern void set_errno(int e);

static long go_session_bio_ctrl(BIO *bp,int cmd,long larg,void *parg) {
	//always return operation not supported
	//http://www.openssl.org/docs/crypto/BIO_ctrl.html
	//printf("go_bio_ctrl %d\n", cmd);
	return 1;
}

static int write_wrapper(BIO* bio,const char* data, int n) {
	return go_session_bio_write(bio,data,n);
}

static int go_session_bio_create( BIO *b ) {
	BIO_set_init(b,1);
	//BIO_set_num(b,-1);
	//BIO_set_ptr(b,NULL);
	BIO_set_flags(b, BIO_FLAGS_READ | BIO_FLAGS_WRITE);
	printf("bio created\n");
	return 1;
}

// a BIO for a client conencted to our server
static BIO_METHOD* go_session_bio_method;

static int init_session_bio_method() {
	go_session_bio_method = BIO_meth_new(BIO_TYPE_SOURCE_SINK,"go session dtls");
	BIO_meth_set_write(go_session_bio_method,write_wrapper);
	BIO_meth_set_read(go_session_bio_method,go_session_bio_read);
	BIO_meth_set_ctrl(go_session_bio_method,go_session_bio_ctrl);
	BIO_meth_set_create(go_session_bio_method,go_session_bio_create);
	BIO_meth_set_destroy(go_session_bio_method,go_session_bio_free);

}

//{
//	BIO_TYPE_SOURCE_SINK,
//	"go session dtls",
//	(int (*)(BIO *, const char *, int))go_session_bio_write,
//	go_session_bio_read,
//	NULL,
//	NULL,
//	go_session_bio_ctrl, // ctrl
//	go_session_bio_create, // new
//	go_session_bio_free // delete
//};

static void init_server_ctx(SSL_CTX *ctx) {
	SSL_CTX_set_min_proto_version(ctx, 0xFEFD); // 1.2
	SSL_CTX_set_max_proto_version(ctx, 0xFEFD); // 1.2
	SSL_CTX_set_read_ahead(ctx, 1);
	SSL_CTX_set_cookie_generate_cb(ctx, &generate_cookie_callback);
	SSL_CTX_set_cookie_verify_cb(ctx, &verify_cookie_callback);

}

static BIO_METHOD* BIO_go_session() {
	return go_session_bio_method;
}

static void setGoSessionId(BIO* bio, unsigned int clientId) {
	unsigned int * pId = malloc(sizeof(unsigned int));
	*pId = clientId;
	BIO_set_data(bio,pId);
}
static unsigned int server_psk_callback(SSL *ssl, const char *identity, unsigned char *psk, unsigned int max_psk_len) {
	return go_server_psk_callback(ssl,identity,(char*)psk,max_psk_len);
}

static void set_psk_callback(SSL *ssl) {
	SSL_set_psk_server_callback(ssl,&server_psk_callback);
}

static void set_cookie_option(SSL *ssl) {
	SSL_set_options(ssl, SSL_OP_COOKIE_EXCHANGE);
}

*/
import "C"
import (
	"errors"
	"net"
	"sync/atomic"
)

func createSslContext() (ctx *DTLSContext, err error) {
	sslCtx := C.SSL_CTX_new(C.DTLSv1_2_server_method())
	ret := int(C.SSL_CTX_set_cipher_list(sslCtx, C.CString("PSK-AES256-CCM8:PSK-AES128-CCM8")))

	if ret != 1 {
		err = errors.New("Unable to set CipherList")
		return
	}

	ctx = &DTLSContext{
		sslCtx: sslCtx,
	}

	return
}

type DTLSContext struct {
	sslCtx *C.SSL_CTX
}

var nextSessionId int32 = 0

func createSslSession(addr net.Addr, ctx *DTLSContext, pskCallback FnHandlePsk) (sslSession *SslSession, err error) {
	ssl := C.SSL_new(ctx.sslCtx)
	id := atomic.AddInt32(&nextSessionId, 1)

	if pskCallback != nil {
		C.set_psk_callback(ssl)
	}

	bio := C.BIO_new(C.BIO_go_session())

	if bio == nil {
		err = errors.New("Error creating session: Bio is nil")
		return
	}

	C.SSL_set_bio(ssl, bio, bio)
	C.setGoSessionId(bio, C.uint(id))
	C.set_cookie_option(ssl)
	C.SSL_set_accept_state(ssl)
	C.DTLSv1_listen

	sslSession = &SslSession{
		addr: addr,
		ssl:  ssl,
		bio:  bio,
	}

	return
}

type SslSession struct {
	addr net.Addr
	ssl  *C.SSL
	bio  *C.BIO
}

type DTLSServerSession struct {
	UDPServerSession
	sslSession *SslSession
}

func (s *DTLSServerSession) GetConnection() ServerConnection {
	return nil
}

type UDPServerSession struct {
	buf    []byte
	addr   net.Addr
	conn   ServerConnection
	server CoapServer
}

func (s *UDPServerSession) GetConnection() ServerConnection {
	return s.conn
}

func (s *UDPServerSession) GetAddress() net.Addr {
	return s.addr
}

func (s *UDPServerSession) Write(b []byte) {
	s.buf = append(s.buf, b...)
}

func (s *UDPServerSession) FlushBuffer() {
	s.buf = nil
}

func (s *UDPServerSession) Read() []byte {
	return s.buf
}

func (s *UDPServerSession) GetServer() CoapServer {
	return s.server
}
