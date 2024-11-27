package task

type Task struct {
	Name    string `yaml:"name"`
	Type    string `yaml:"type"`
	Command string `yaml:"command"`
}

type TaskList struct {
	Tasks []Task `yaml:"tasks"`
}
