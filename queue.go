package main

import (
	"log"
	"os"
	"time"
)

var queue_task_id int64
var tasks_db map[int64]*DownloadRequest

func (DC *DownloadCenter) newQueue(config *QueueConfig) *Queue {
	queue_task_id = 1
	queue := Queue{
		logger:        log.New(os.Stdout, "", log.LstdFlags),
		checkInterval: 3 * time.Second, /// 100 * time.Millisecond
		done:          make(chan bool),
		task:          make(chan int64),
		waiting:       make(map[int64]*WaitingQueueTask),
		active:        make(map[int64]*RunningQueueTask),
		groups:        make(map[string]*QueueStats_Group),
		sites:         make(map[string]*QueueStats_Site),
		users:         make(map[int64]*QueueStats_User),
	}

	queue.updateConfig(config)

	tasks_db = make(map[int64]*DownloadRequest)

	return &queue
}

func (q *Queue) startQueue() {
	DC.logger.Println("Queue: Starting")
	go q.queueRunner()
}

func (q *Queue) stopQueue() {
	DC.logger.Println("Queue: Stopping")
	DC.done <- true
}

func (q *Queue) updateConfig(config *QueueConfig) {
	q.statsGroups(config.Groups)
	q.statsSites(config.Sites)
}

func (q *Queue) queueRunner() {
	for {
		for task := range q.waiting {
			q.maybeStartTask(task)
		}

		select {
		case <-q.done:
			return
		case done_task := <-q.task:
			q.doneTask(done_task)
			continue
		case <-time.After(q.checkInterval):
			continue
		}
	}
}

func (q *Queue) statsGroups(groups map[string]*QueueConfig_Group) {
	for group_name := range groups {
		group := groups[group_name]
		_, ok := q.groups[group_name]
		if !ok {
			q.groups[group_name] = &QueueStats_Group{
				Name:           group_name,
				PerUser:        group.PerUser,
				Simultaneously: group.Simultaneously,
				Running:        0,
			}
		} else {
			qg := q.groups[group_name]
			qg.PerUser = group.PerUser
			qg.Simultaneously = group.Simultaneously
		}
	}
}

func (q *Queue) statsSites(sites map[string]*QueueConfig_Site) {
	for site_name := range sites {
		site := sites[site_name]
		_, ok := q.sites[site_name]
		if !ok {
			q.sites[site_name] = &QueueStats_Site{
				Name:           site_name,
				Group:          site.Group,
				PerUser:        site.PerUser,
				Simultaneously: site.Simultaneously,
				Running:        0,
			}
		} else {
			qg := q.sites[site_name]
			qg.Group = site.Group
			qg.PerUser = site.PerUser
			qg.Simultaneously = site.Simultaneously
		}
	}
}

func (q *Queue) maybeStartTask(taskId int64) {
	task, ok := q.waiting[taskId]
	if !ok {
		return
	}
	// delete(q.waiting, taskId)

	q.logger.Println("maybeStartTask")
	q.logger.Printf("Task: %+v\n", task)
	group := q.groups[task.Group]
	// q.logger.Printf("Group: %+v\n", group)

	if !group.CanStart() {
		return
	}
	site := q.sites[task.Site]
	// q.logger.Printf("Site: %+v\n", site)
	if !site.CanStart() {
		return
	}

	user, ok := q.users[task.UserId]
	if ok {
		user_site, ok := user.BySite[task.Site]
		if ok {
			if user_site >= site.PerUser {
				return
			}
		}
	}
	q.startTask(task)
}

func (q *Queue) save_to_db(request *DownloadRequest) bool {
	queue_task_id++
	tasks_db[queue_task_id] = request
	return true
}
func (q *Queue) addTask(request DownloadRequest) bool {
	if q.save_to_db(&request) {
		_, exists := q.sites[request.Site]
		if exists {
			wqt := &WaitingQueueTask{
				TaskId: queue_task_id,
				UserId: request.UserId,
				Site:   request.Site,
				Group:  q.sites[request.Site].Group,
			}
			q.waiting[wqt.TaskId] = wqt
			return true
		}
		return true
	}
	return false
}

func (q *Queue) startTask(task *WaitingQueueTask) { // *RunningQueueTask
	q.logger.Println("startTask")
	q.logger.Printf("\n%+v\n\n", task)
	_, ok := q.users[task.UserId]
	if !ok {
		user := &QueueStats_User{
			Total:  0,
			BySite: make(map[string]int),
		}
		q.users[task.UserId] = user
	}
	// d := &Downloader{}
	// rqt := &RunningQueueTask{
	// 	TaskId:     queue_task_id,
	// 	Downloader: d,
	// 	done:       &q.task,
	// }
	// q.active[queue_task_id] = rqt
	// go rqt.Downloader.Start()
	// return rqt
}

func (q *Queue) stopTask(taskId int64) {

}

func (q *Queue) doneTask(taskId int64) {

}
