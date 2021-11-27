// 1. Доработать программу из практической части так, чтобы при отправке ей сигнала SIGUSR1
// она увеличивала глубину поиска на 2.
// 2. Добавить общий таймаут на выполнение следующих операций: работа парсера, получений ссылок со страницы,
// формирование заголовка.

package main

import (
	"context"
	"github.com/t0pep0/GB_best_go1/hw1/config"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	crw "github.com/t0pep0/GB_best_go1/hw1/crawler"
)

func main() {
	cfg := config.NewConfig()
	r := crw.NewRequester(crw.NewHttpClient(time.Duration(cfg.Timeout) * time.Second))
	cr := crw.NewCrawler(r, cfg.MaxDepth)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.Timeout)*time.Second)
	go cr.Scan(ctx, cfg.Url, 1)             //Запускаем краулер в отдельной рутине
	go processResult(ctx, cancel, cr, *cfg) //Обрабатываем результаты в отдельной рутине

	sigCh := make(chan os.Signal)                         //Создаем канал для приема сигналов
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGUSR1) //Подписываемся на сигнал SIGINT, SIGUSR1
	for {
		select {
		case <-ctx.Done(): //Если всё завершили - выходим
			return
		case sig := <-sigCh:
			switch sig {
			case syscall.SIGINT:
				cancel() //Если пришёл сигнал SigInt - завершаем контекст
				return
			case syscall.SIGUSR1:
				depth := uint64(2)
				// Increment depth while catch SIGUSR1
				cr.IncreaseMaxDepth(depth)
				log.Printf("Depth increment set to: %d\n", depth)
			}
		}
	}
}

func processResult(ctx context.Context, cancel func(), cr crw.Crawler, cfg config.Config) {
	var maxResult, maxErrors = cfg.MaxResults, cfg.MaxErrors
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-cr.ChanResult():
			if msg.Err != nil {
				maxErrors--
				log.Printf("crawler result return err: %s\n", msg.Err.Error())
				if maxErrors <= 0 {
					cancel()
					return
				}
			} else {
				maxResult--
				log.Printf("crawler result: [url: %s] Title: %s\n", msg.Url, msg.Title)
				if maxResult <= 0 {
					cancel()
					return
				}
			}
		}
	}
}
