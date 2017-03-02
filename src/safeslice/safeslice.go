package safeslice

type safeSlice chan commandData

type commandData struct {
	action commandAction
	index int
	item interface{}
	result chan<- interface{}
	data chan<- []interface{}
	updater UpdateFunc
}

type commandAction int

const (
	insert commandAction = iota
	remove
	at
	update
	end
	length
)

type UpdateFunc func(interface{}) interface{}

type SafeSlice interface {
	Append(interface{})
	At(int) interface{}
	Close() []interface{}
	Delete(int)
	Len() int
	Update(int,UpdateFunc)
}

func New() SafeSlice{
	list := make(safeSlice)
	go list.run()
	return list
}

func (slice safeSlice) run() {
	list := make([]interface{}, 0)
	for command := range slice {
		switch command.action {
		case insert:
			list = append(list,command.item)
		case remove:
			if 0 < command.index && command.index < len(list) {
				list = append(list[:command.index],list[command.index+1:])
			}
		case at:
			if 0 < command.index && command.index < len(list) {
				command.result <- list[command.index]
			}else{
				command.result <- nil
			}
		case length:
			command.result <- len(list)
		case update:
			if 0 < command.index && command.index < len(list) {
				list[command.index] = command.updater(list[command.index])
			}
		case end:
			close(slice)
			command.data <- list
		}
	}
}

func (slice safeSlice) Append(item interface{}){
	slice <- commandData{action: insert,item: item}
}

func (slice safeSlice) Delete(index int){
	slice <- commandData{action: remove,index: index}
}

func (slice safeSlice) At(index int) interface{}{
	reply := make(chan interface{})
	slice <- commandData{at,index,reply,nil,nil,nil}
	return <-reply
}

func (slice safeSlice) Len() int{
	reply := make(chan interface{})
	slice <- commandData{action: length, result: reply}
	return (<-reply).(int)
}

func (slice safeSlice) Update(index int,updater UpdateFunc) {
	slice <- commandData{action: update,index: index,updater: updater}
}

func (slice safeSlice) Close() []interface{} {
	reply := make(chan []interface{})
	slice <- commandData{action: end,data: reply}
	return <-reply

}








