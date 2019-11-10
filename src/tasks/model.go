package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"sync"
	"time"
)

type Task struct {
	ID      string `json:"id"`
	Command string `json:"command"`
	Started bool   `json:"started"`
	Agent   string `json:"agent"`
}

type Agent struct {
	ID            string
	LastHeartbeat time.Time
}

type Manager struct {
	tasks     map[string]*Task
	agents    map[string]*Agent
	directory string
	mutex     *sync.Mutex
}

func NewManager() *Manager {
	return &Manager{
		map[string]*Task{},
		map[string]*Agent{},
		"",
		&sync.Mutex{},
	}
}

func (m *Manager) AddTask(id string, command string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	var err error

	if _, ok := m.tasks[id]; ok {
		return fmt.Errorf("this id is associated with a task already")
	}

	task := &Task{id, command, false, ""}

	err = m.createTaskDirectory(task)
	if err != nil {
		return err
	}

	err = m.saveTaskCommand(task)
	if err != nil {
		return err
	}

	m.tasks[id] = task

	err = m.saveTasks()

	if err != nil {
		delete(m.tasks, id)

		_ = m.deleteTaskDirectory(task)

		return err
	}

	return nil
}

func (m *Manager) AddAgent(id string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if _, ok := m.agents[id]; ok {
		return fmt.Errorf("this id is associated with an agent already")
	}

	m.agents[id] = &Agent{id, time.Now()}

	return nil
}

func (m *Manager) RemoveAgentAndDeallocateTasks(id string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	return m.removeAgentAndDeallocateTasks(id)
}

func (m *Manager) AssignNotStartedTaskToAgent(agent string) (*Task, bool, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	var err error

	for _, task := range m.tasks {
		if task.Started {
			continue
		}

		task.Started = true
		task.Agent = agent

		err = m.saveTasks()
		if err != nil {
			task.Started = true
			task.Agent = ""

			return nil, false, err
		}

		_ = m.resetTaskOutput(task)

		return task, true, nil
	}

	return nil, false, nil
}

func (m *Manager) CompleteTask(id string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	_, ok := m.tasks[id]
	if !ok {
		return fmt.Errorf("this id is not associated with an agent")
	}

	delete(m.tasks, id)

	_ = m.saveTasks()

	return nil
}

func (m *Manager) AppendTaskOutput(id string, text string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	task, ok := m.tasks[id]
	if !ok {
		return fmt.Errorf("this id is not associated with an agent")
	}

	err := m.appendTaskOutput(task, text)
	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) Heartbeat(id string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	agent, ok := m.agents[id]
	if !ok {
		return fmt.Errorf("this id is not associted with an agent")
	}

	agent.LastHeartbeat = time.Now()
	return nil
}

func (m *Manager) WatchAgents() {
	go func() {
		for true {
			time.Sleep(time.Second * 5)

			now := time.Now()

			m.mutex.Lock()

			lost := []*Agent{}
			for _, agent := range m.agents {
				if now.Sub(agent.LastHeartbeat) > 15*time.Second {
					lost = append(lost, agent)
				}
			}

			for _, agent := range lost {
				m.removeAgentAndDeallocateTasks(agent.ID)
			}

			m.mutex.Unlock()
		}
	}()
}

func (m *Manager) removeAgentAndDeallocateTasks(id string) error {
	var err error

	agent, ok := m.agents[id]
	if !ok {
		return fmt.Errorf("this id is not associated with an agent")
	}

	delete(m.agents, id)

	for _, task := range m.tasks {
		if task.Started || task.Agent != agent.ID {
			continue
		}

		task.Agent = ""
		task.Started = false

		_ = m.resetTaskOutput(task)
	}

	err = m.saveTasks()
	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) saveTasks() error {
	var err error

	data, err := json.MarshalIndent(m.tasks, "", "  ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path.Join(m.directory, "tasks.json"), data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) createTaskDirectory(task *Task) error {
	taskDirectoryPath := path.Join(m.directory, task.ID)

	err := os.Mkdir(taskDirectoryPath, 0755)
	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) deleteTaskDirectory(task *Task) error {
	taskDirectoryPath := path.Join(m.directory, task.ID)

	err := os.Remove(taskDirectoryPath)
	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) saveTaskCommand(task *Task) error {
	commandFilePath := path.Join(m.directory, task.ID, "command.txt")

	err := ioutil.WriteFile(commandFilePath, []byte(task.Command), 0644)
	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) resetTaskOutput(task *Task) error {
	outputFilePath := path.Join(m.directory, task.ID, "output.txt")

	if _, err := os.Stat(outputFilePath); err != nil {
		return nil
	}

	err := os.Remove(outputFilePath)
	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) appendTaskOutput(task *Task, text string) error {
	var err error

	outputFilePath := path.Join(m.directory, task.ID, "output.txt")

	file, err := os.OpenFile(outputFilePath, os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = file.WriteString(text)
	if err != nil {
		return err
	}

	return nil
}
