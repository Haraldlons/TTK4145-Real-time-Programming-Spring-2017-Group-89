#include <pthread.h>
#include <stdio.h>

/*To compile the file
gcc -pthread -o oppg4 oppg4.c
*/

int i = 0;

void *threadFunc1(){
	int j;
	for (j = 0; j < 1000000; ++j){
		++i;
	}
}

void *threadFunc2(){
	int j;
	for (j = 0; j < 1000000; ++j){
		--i;
	}
}

int main(void){
	//Initialize threads
	pthread_t thread1;
	pthread_t thread2;

	//Create threads
	pthread_create(&thread1, NULL, threadFunc1, "Creating thread 1"); //Create thread 1
	pthread_create(&thread2, NULL, threadFunc2, "Creating thread 2"); //Create thread 2

	//Join threads
	pthread_join(thread1, NULL);
	pthread_join(thread2, NULL);

	printf("%i\n", i);
	
	return 0;
}