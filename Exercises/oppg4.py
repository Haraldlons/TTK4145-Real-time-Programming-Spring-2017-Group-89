from threading import Thread

i = 0

def threadFunc1():
	global i
	for j in range(0, 1000000):
		i+=1

def threadFunc2():
	global i
	for j in range(0, 1000000):
		i-=1

def main():
	#Initialize threads
	thread1 = Thread(target = threadFunc1, args = (), )
	thread2 = Thread(target = threadFunc2, args = (), )

	#Start threads
	thread1.start()
	thread2.start()

	#Join threads
	thread1.join()
	thread2.join()

	print(i)

main()