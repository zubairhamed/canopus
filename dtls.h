#ifndef __DTLS_H_
#define __DTLS_H_

#include <stdlib.h>
#include <string.h>
#include <openssl/err.h>
#include <openssl/ssl.h>
#include <openssl/bio.h>
#include <internal/bio.h>

extern int go_conn_bio_write(BIO* bio, char* buf, int num);
extern int go_conn_bio_read(BIO* bio, char* buf, int num);
extern int go_conn_bio_free(BIO* bio);
extern unsigned int go_psk_callback(SSL *ssl, char *hint, char *identity, unsigned int max_identity_len, char *psk, unsigned int max_psk_len);

// a BIO for a client conencted to our server
static BIO_METHOD* go_session_bio_method;

extern int go_session_bio_write(BIO* bio, char* buf, int num);
extern int go_session_bio_read(BIO* bio, char* buf, int num);
extern int go_session_bio_free(BIO* bio);

extern unsigned int go_server_psk_callback(SSL *ssl, char *identity, char *psk, unsigned int max_psk_len);
extern int generate_cookie_callback(SSL* ssl, unsigned char* cookie, unsigned int *cookie_len);
extern int verify_cookie_callback(SSL* ssl, unsigned char* cookie, unsigned int cookie_len);

static BIO_METHOD* BIO_go_session() {
	return go_session_bio_method;
}

static int get_errno(void)
{
	return errno;
}

static void set_errno(int e)
{
	errno = e;
}

static void setGoData(BIO* bio, char *data) {
	char *sId = malloc(sizeof(data));
	BIO_set_data(bio, data);
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

static unsigned int psk_callback(SSL *ssl, const char *hint,
        char *identity, unsigned int max_identity_len,
        unsigned char *psk, unsigned int max_psk_len) {
	return go_psk_callback(ssl,hint,identity,max_identity_len,(char*)psk,max_psk_len);
}

static void init_ctx(SSL_CTX *ctx) {
	SSL_CTX_set_read_ahead(ctx, 1);

	SSL_CTX_set_psk_client_callback(ctx, &psk_callback);
}


static void init_lib() {
	SSL_library_init();
	ERR_load_BIO_strings();
	SSL_load_error_strings();
}


static void set_proto_1_2(SSL_CTX *ctx) {
	SSL_CTX_set_min_proto_version(ctx, 0xFEFD); // 1.2
	SSL_CTX_set_max_proto_version(ctx, 0xFEFD); // 1.2
}

static int write_wrapper(BIO* bio,const char* data, int n) {
	return go_session_bio_write(bio,data,n);
}

static long go_session_bio_ctrl(BIO *bp,int cmd,long larg,void *parg) {
	//always return operation not supported
	//http://www.openssl.org/docs/crypto/BIO_ctrl.html
	//printf("go_bio_ctrl %d\n", cmd);
	return 1;
}

static int go_session_bio_create( BIO *b ) {
	BIO_set_init(b,1);
	//BIO_set_num(b,-1);
	//BIO_set_ptr(b,NULL);
	BIO_set_flags(b, BIO_FLAGS_READ | BIO_FLAGS_WRITE);
	printf("bio created\n");
	return 1;
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

static void setGoClientId(BIO* bio, unsigned int clientId) {
	unsigned int * pId = malloc(sizeof(unsigned int));
	*pId = clientId;
	BIO_set_data(bio,pId);
}

#endif