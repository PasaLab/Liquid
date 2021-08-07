package main

import "sync"

type GroupManager struct {
	groups map[string]Group
	mu     sync.Mutex
}

var groupManagerInstance *GroupManager
var groupManagerInstanceLock sync.Mutex

func InstanceOfGroupManager() *GroupManager {
	defer groupManagerInstanceLock.Unlock()
	groupManagerInstanceLock.Lock()

	if groupManagerInstance == nil {
		groupManagerInstance = &GroupManager{groups: map[string]Group{}}
		groupManagerInstance.groups["default"] = Group{Name: "default", Weight: 10, Reserved: false}
	}
	return groupManagerInstance
}

func (gm *GroupManager) Start() {

}

func (gm *GroupManager) Add(group Group) MsgGroupCreate {
	defer gm.mu.Unlock()
	gm.mu.Lock()
	if _, ok := gm.groups[group.Name]; ok {
		return MsgGroupCreate{Code: 1, Error: "Name already exists!"}
	}
	gm.groups[group.Name] = group
	return MsgGroupCreate{}
}

func (gm *GroupManager) Update(group Group) MsgGroupCreate {
	defer gm.mu.Unlock()
	gm.mu.Lock()
	if _, ok := gm.groups[group.Name]; !ok {
		return MsgGroupCreate{Code: 1, Error: "Group not exists!"}
	}
	gm.groups[group.Name] = group
	return MsgGroupCreate{}
}

func (gm *GroupManager) Remove(group Group) MsgGroupCreate {
	defer gm.mu.Unlock()
	gm.mu.Lock()
	if _, ok := gm.groups[group.Name]; !ok {
		return MsgGroupCreate{Code: 1, Error: "Group not exists!"}
	}
	delete(gm.groups, group.Name)
	return MsgGroupCreate{}
}

func (gm *GroupManager) List() MsgGroupList {
	defer gm.mu.Unlock()
	gm.mu.Lock()
	var result []Group
	for _, v := range gm.groups {
		result = append(result, v)
	}
	return MsgGroupList{Groups: result}
}

func (gm *GroupManager) get(name string) *Group {
	defer gm.mu.Unlock()
	gm.mu.Lock()

	for _, v := range gm.groups {
		if v.Name == name {
			return &v
		}
	}
	return nil
}
