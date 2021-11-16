package main

var users map[int64]*UserControlBlock

func AddUser(uid uint64, name string) []byte {
	newControlBlock := new(UserControlBlock)
	newControlBlock.vars["$name"] = name
	newControlBlock.uid = uid
	return []byte("success")
}
