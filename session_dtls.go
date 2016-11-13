package canopus

/*
#cgo LDFLAGS: -L${SRCDIR}/openssl -lssl -lcrypto -ldl
#cgo CFLAGS: -g -Wno-deprecated -Wno-error -I${SRCDIR}/openssl/include

#include "dtls.h"
*/
import "C"
import (
	"errors"
	"fmt"
	"io"
	"log"
	"syscall"
	"unsafe"
)

func newSslSession(session *DTLSServerSession, ctx *DTLSContext, pskCallback FnHandlePsk, id string) (err error) {
	ssl := C.SSL_new(ctx.sslCtx)

	if pskCallback != nil {
		C.set_psk_callback(ssl)
	}

	log.Println("BIO GO SESSION ==", C.BIO_go_session())
	bio := C.BIO_new(C.BIO_go_session())

	if bio == nil {
		err = errors.New("Error creating session: Bio is nil")
		return
	}

	C.SSL_set_bio(ssl, bio, bio)
	C.setGoData(bio, C.CString(id+","+session.addr.String()))
	C.set_cookie_option(ssl)
	C.SSL_set_accept_state(ssl)
	C.DTLSv1_listen

	session.ssl = ssl
	session.bio = bio
	session.rcvd = make(chan []byte)

	return
}

type DTLSServerSession struct {
	UDPServerSession
	ssl  *C.SSL
	bio  *C.BIO
	rcvd chan []byte
}

func (s *DTLSServerSession) GetConnection() ServerConnection {
	return nil
}

func (s *DTLSServerSession) Received(b []byte) (n int) {
	s.rcvd <- b
	return len(b)
}

func (s *DTLSServerSession) Read(b []byte) (n int, err error) {
	// TODO test if closed?
	length := len(b)

	ret := C.SSL_read(s.ssl, unsafe.Pointer(&b[0]), C.int(length))
	fmt.Println("SSL READ done")
	if err := s.getError(ret); err != nil {
		return 0, err
	}
	// if there's no error, but a return value of 0
	// let's say it's an EOF
	if ret == 0 {
		return 0, io.EOF
	}
	return int(ret), nil
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
			msg += s.getErrorString(errCode)
		}
		C.ERR_clear_error()
		return errors.New(msg)
	}
	return nil
}

func (s *DTLSServerSession) getErrorString(code C.ulong) string {
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
