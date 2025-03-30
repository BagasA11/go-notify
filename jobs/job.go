package jobs

import (
	"go-notify/dto"
	"go-notify/mail"
	"log"
	"sync"
)

type WorkerPool struct {
	Queue  chan dto.Body
	Worker uint
}

func NewWorkerPool(worker uint, size int) *WorkerPool {
	return &WorkerPool{
		Queue:  make(chan dto.Body, size),
		Worker: worker,
	}
}

func (wp *WorkerPool) Add(request dto.Body) {
	wp.Queue <- request
}

func (wp *WorkerPool) Work(wg *sync.WaitGroup) {
	defer wg.Done()
	log.Println("creating worker")
	for job := range wp.Queue {
		log.Println("processing job ", job.Receiver)

		err := mail.M.SendMail(job)
		if err != nil && (err == mail.ErrRecipient || err == mail.ErrSender) {
			continue
		} else {
			wp.Add(job)
		}
		log.Println("success send mail to: ", job.Receiver)
	}
}

func (wp *WorkerPool) Do(wg *sync.WaitGroup) {
	for i := 0; i < int(wp.Worker); i++ {
		go wp.Work(wg)
	}
}
