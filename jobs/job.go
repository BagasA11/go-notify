package jobs

import (
	"go-notify/dto"
	"go-notify/mail"
	"log"
	"sync"
)

type queue struct {
	Body  dto.Body
	retry uint
}

type WorkerPool struct {
	Queue  chan queue
	Worker uint
}

func NewWorkerPool(worker uint, size int) *WorkerPool {
	return &WorkerPool{
		Queue:  make(chan queue, size),
		Worker: worker,
	}
}

func (wp *WorkerPool) Add(request dto.Body, rty ...uint) {
	r := uint(0)
	if rty != nil {
		r = rty[0]
	}
	q := queue{
		retry: r,
		Body:  request,
	}
	wp.Queue <- q
}

func (wp *WorkerPool) Work(wg *sync.WaitGroup) {
	defer wg.Done()
	log.Println("initializing worker")
	for job := range wp.Queue {
		log.Println("processing job ", job.Body.Receiver)

		err := mail.M.SendMail(job.Body)
		if err != nil {
			log.Println(err)
			// if sender or recepient not found, omit proccess and continue to next queue
			if err == mail.ErrRecipient || err == mail.ErrSender {
				continue
			}
			// retry
			// pass proccess and continue to next queue
			if job.retry > 3 {
				log.Println("this job has been retried more than 3 times")
				continue
			}
			job.retry++
			wp.Add(job.Body, job.retry)
			continue
		}
		log.Println("success send mail to: ", job.Body.Receiver)
	}
}

func (wp *WorkerPool) Do(wg *sync.WaitGroup) {
	for i := 0; i < int(wp.Worker); i++ {
		log.Println("creating worker: ", i+1)
		go wp.Work(wg)
	}
}
