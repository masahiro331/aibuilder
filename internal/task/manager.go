package task

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	openai "github.com/sashabaranov/go-openai"
	yaml "gopkg.in/yaml.v2"

	"aibuilder/internal/ai"
	"aibuilder/internal/config"
	"aibuilder/internal/logger"
	"aibuilder/internal/personas"
	"aibuilder/internal/utils"

	"github.com/sirupsen/logrus"
)

type Manager struct {
	aiClient  ai.Client
	logger    *logger.Logger
	debugMode bool
}

func NewManager(cfg *config.Config, logger *logger.Logger, debugMode bool) *Manager {
	aiClient := ai.NewClient(cfg, logger)
	return &Manager{
		aiClient:  aiClient,
		logger:    logger,
		debugMode: debugMode,
	}
}

type DeveloperInput struct {
	RequestText  string
	Requirements RequirementsJSON
	Designs      DesignsJSON
}

type RequirementsJSON struct {
	Requirements []struct {
		ID          string `json:"id"`
		Description string `json:"description"`
	} `json:"requirements"`
}

type DesignsJSON struct {
	Designs []struct {
		RequirementID string `json:"requirement_id"`
		Design        struct {
			Architecture string   `json:"architecture"`
			DataFlow     string   `json:"data_flow"`
			Components   []string `json:"components"`
		} `json:"design"`
	} `json:"designs"`
}

func (d DeveloperInput) String() string {
	reqirements := map[string]string{}
	for _, v := range d.Requirements.Requirements {
		reqirements[v.ID] = v.Description
	}
	requirementsText := []string{}
	for _, v := range d.Designs.Designs {
		requirement := reqirements[v.RequirementID]
		requirementText := fmt.Sprintf(`
        Requirement: %s
        Architecture: %s
        DataFlow: %s
        Components: %s
        `, requirement, v.Design.Architecture, v.Design.DataFlow, strings.Join(v.Design.Components, ","))
		requirementsText = append(requirementsText, requirementText)
	}
	return fmt.Sprintf(`
    ユーザの要望: %s
    プログラムの設計: %s
    `, d.RequestText, strings.Join(requirementsText, "---\n"))
}

func (m *Manager) Start(toolDescription string) error {
	// 要件定義者との対話
	requirementsJSON, err := m.communicateWithPersona(personas.RequirementsEngineer, toolDescription)
	if err != nil {
		return fmt.Errorf("Failed to get requirements: %v", err)
	}
	m.logger.Info("Requirements received")

	requirements := RequirementsJSON{}
	json.NewDecoder(strings.NewReader(requirementsJSON)).Decode(&requirements)

	// 設計者との対話
	designJSON, err := m.communicateWithPersona(personas.Designer, requirementsJSON)
	if err != nil {
		return fmt.Errorf("Failed to get design: %v", err)
	}
	m.logger.Info("Design received")

	designs := DesignsJSON{}
	json.NewDecoder(strings.NewReader(designJSON)).Decode(&designs)

	input := DeveloperInput{
		RequestText:  toolDescription,
		Requirements: requirements,
		Designs:      designs,
	}

	// 開発者との対話
	os.WriteFile("prompt.txt", []byte(input.String()), 0600)
	tasksJSON, err := m.communicateWithPersona(personas.Developer, input.String())
	if err != nil {
		return fmt.Errorf("Failed to get tasks: %v", err)
	}
	m.logger.Info("Tasks received")

	tasksJSON = strings.TrimSuffix(strings.TrimPrefix(tasksJSON, "```yaml"), "```")
	var taskList TaskList
	err = yaml.Unmarshal([]byte(tasksJSON), &taskList)
	if err != nil {
		return fmt.Errorf("Failed to parse tasks: %v", err)
	}

	// タスクの実行
	for _, task := range taskList.Tasks {
		if m.logger.Level >= logrus.DebugLevel {
			m.logger.Debugf("Preparing to execute task: %s", task.Name)
		}

		err := m.executeTask(task)
		if err != nil {
			m.logger.Errorf("Task %s failed: %v", task.Name, err)
			continue
		}
		m.logger.Infof("Task %s completed successfully", task.Name)
	}

	return nil
}

func (m *Manager) communicateWithPersona(persona string, input string) (string, error) {
	messages := []openai.ChatCompletionMessage{
		personas.GetPersona(persona),
		{
			Role:    openai.ChatMessageRoleUser,
			Content: input,
		},
	}

	response, err := m.aiClient.SendMessage(messages)
	if err != nil {
		return "", err
	}

	return response, nil
}

func (m *Manager) executeTask(task Task) error {
	switch task.Type {
	case "execute_command":
		if m.debugMode {
			m.logger.Debugf("Debug Mode: Would execute command: %s", task.Command)
		} else {
			if utils.IsSafeCommand(task.Command) {
				output, err := utils.ExecuteCommand(task.Command)
				if err != nil {
					return err
				}
				m.logger.Infof("Command Output: %s", output)
			} else {
				return fmt.Errorf("Unsafe command detected: %s", task.Command)
			}
		}
	case "ai_communication":
		if m.debugMode {
			m.logger.Debugf("Debug Mode: Would send AI prompt: %s", task.Command)
		} else {
			response, err := m.aiClient.SendMessage([]openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: task.Command,
				},
			})
			if err != nil {
				return err
			}
			m.logger.Infof("AI Response: %s", response)
		}
	default:
		return fmt.Errorf("Unknown task type: %s", task.Type)
	}
	return nil
}
