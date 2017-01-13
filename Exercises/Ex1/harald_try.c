#include <pthread.h>

//Define working function
void *foo(void *args){}

int main(){
	// Initialize pthread_attr_t
	pthread_attr_t attr;
	pthread_attr_init(attr);

	// Create a thread
	pthread_t thread
	pthread_create(&thread, &attr, worker_functino, arg);

	// Exit current thread

}