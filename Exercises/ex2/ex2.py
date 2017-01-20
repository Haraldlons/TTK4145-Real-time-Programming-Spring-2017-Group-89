import threading
# To compile and run
# python oppg4.py


i = 0
lock = threading.Lock()

def threadFunc1():
	global i
	for j in range(0, 1000000):
		lock.acquire(1)
		i+=1
		lock.release()

def threadFunc2():
	global i
	for j in range(0, 1000000):
		lock.acquire(1)
		i-=1
		lock.release()

def main():
	#Initialize threads
	thread1 = threading.Thread(target = threadFunc1, args = (), )
	thread2 = threading.Thread(target = threadFunc2, args = (), )

	#Initialize locks
	#lock1 = thread1.allocate_lock()

	#Start threads
	thread1.start()
	thread2.start()

	#Join threads
	thread1.join()
	thread2.join()

	print(i)

main()