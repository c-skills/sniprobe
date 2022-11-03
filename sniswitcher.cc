/* SNI switcher to work with w3m (sic) to bypass censorship blockings
 *
 * (C) 2022 c->skills research
 */
#include <cstdio>
#include <cstdlib>
#include <cstring>
#include <dlfcn.h>

extern "C" long SSL_ctrl(void *, int, long, void *);

#define SSL_CTRL_SET_TLSEXT_HOSTNAME            55
#define TLSEXT_NAMETYPE_host_name		0

const char *good_sni = GOOD_SNI;
const char *lib = SSL_LIB;

int (*o_SSL_do_handshake)(void *session) = nullptr;
int (*o_SSL_connect)(void *session) = nullptr;

extern "C" void __attribute__ ((constructor)) init()
{
	void *h = dlopen(lib, RTLD_NOW|RTLD_NODELETE);
	if (!h) {
		printf("Failed to open SSL lib. Change $SSL_LIB path in Makefile\n");
		exit(1);
	}
	void *f1 = dlsym(h, "SSL_do_handshake");
	void *f2 = dlsym(h, "SSL_connect");
	dlclose(h);

	if (!f1 || !f2) {
		printf("Failed to resolv SSL symbols.\n");
		exit(1);
	}

	memcpy(&o_SSL_do_handshake, &f1, sizeof(f1));
	memcpy(&o_SSL_connect, &f2, sizeof(f2));
}


extern "C" int SSL_do_handshake(void *session)
{
	SSL_ctrl(session, SSL_CTRL_SET_TLSEXT_HOSTNAME, TLSEXT_NAMETYPE_host_name, (void *)good_sni);
	return o_SSL_do_handshake(session);
}


extern "C" int SSL_connect(void *session)
{
	SSL_ctrl(session, SSL_CTRL_SET_TLSEXT_HOSTNAME, TLSEXT_NAMETYPE_host_name, (void *)good_sni);
	return o_SSL_connect(session);
}

