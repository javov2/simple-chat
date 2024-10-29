package client

import (
	"bufio"
	"fmt"
	globals "go-chat/config"
	"log"
	"net"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const UserPrompt = "[%s][%s] >> "

type serverMessage struct {
	msg string
}

type uiModel struct {
	viewport    viewport.Model
	textarea    textarea.Model
	senderStyle lipgloss.Style
	err         error
	messages    []string
	sessionId   string
	conn        net.Conn
}

func (m uiModel) Init() tea.Cmd {
	tea.SetWindowTitle("Bubble Tea Example")
	return textarea.Blink
}

func (m uiModel) View() string {
	return fmt.Sprintf(
		"%s\n\n%s",
		m.viewport.View(),
		m.textarea.View(),
	) + "\n\n"
}

func (m uiModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)

	//m.textarea, tiCmd = m.textarea.Update(msg)
	//m.viewport, vpCmd = m.viewport.Update(msg)

	switch msg := msg.(type) {
	case serverMessage:
		m.messages = append(m.messages, msg.msg)
		m.viewport.SetContent(strings.Join(m.messages, ""))
		m.viewport.GotoBottom()

	case tea.KeyMsg:

		m.textarea, tiCmd = m.textarea.Update(msg)
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			fmt.Println(m.textarea.Value())
			return m, tea.Quit
		case tea.KeyEnter:
			m.conn.Write([]byte(m.textarea.Value() + "\n"))
			m.messages = append(m.messages, m.senderStyle.Render("You: ")+m.textarea.Value()+"\n")
			m.viewport.SetContent(strings.Join(m.messages, ""))
			m.viewport.GotoBottom()
			m.textarea.Reset()
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	return m, tea.Batch(tiCmd, vpCmd)
}

func initChatView(conn net.Conn, sessionId string) uiModel {
	ta := textarea.New()
	ta.Placeholder = "Send a message..."
	ta.Focus()

	ta.Prompt = "┃ "
	ta.CharLimit = 280

	ta.SetWidth(90)
	ta.SetHeight(1)

	// Remove cursor line styling
	//ta.FocusedStyle.CursorLine = lipgloss.NewStyle()

	ta.ShowLineNumbers = false

	vp := viewport.New(90, 15)
	vp.SetContent(`Welcome to the chat room!
Type a message and press Enter to send.`)

	ta.KeyMap.InsertNewline.SetEnabled(false)

	return uiModel{
		textarea:    ta,
		messages:    []string{},
		viewport:    vp,
		senderStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
		err:         nil,
		conn:        conn,
		sessionId:   sessionId,
	}
}

func Client(serverAddress string) {
	fmt.Print("Connecting to: ", serverAddress+"... ")
	c, err := net.Dial("tcp", serverAddress)
	if err != nil {
		fmt.Print("[ERROR]")
		return
	}
	defer c.Close()
	fmt.Print("[OK]")
	fmt.Println()
	// login
	sessionId := login(c)
	app := tea.NewProgram(initChatView(c, sessionId))
	go receiveMessages(c, app)
	if _, err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

func login(conn net.Conn) string {
	fmt.Print("Logging in... \n")
	conn.Write([]byte(globals.Commands.Login + "\n"))
	reader := bufio.NewReader(conn)
	msg, _ := reader.ReadString('\n')
	fmt.Printf("Your Session Id: %s", msg)
	fmt.Println()
	return strings.TrimRight(msg, "\r\n")
}

func receiveMessages(conn net.Conn, program *tea.Program) {
	reader := bufio.NewReader(conn)
	for {
		msg, _ := reader.ReadString('\n')
		program.Send(serverMessage{msg: msg})
		//	fmt.Println("receiveing message: " + msg)
	}
}
