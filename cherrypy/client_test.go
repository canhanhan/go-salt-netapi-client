package cherrypy

import "testing"

func TestSomethingElse(t *testing.T) {
	client, err := NewClient("https://192.168.50.10:8000", "test", "test", "pam")
	if err != nil {
		t.Fatal(err)
	}

	result, err := client.GetKeys()
	if err != nil {
		t.Fatal(err)
	}

	print(result)
}

// func TestSomethingElse(t *testing.T) {
// 	client, err := NewClient("https://192.168.50.10:8000", "test", "test", "pam")
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	result, err := client.GetMinions()
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	print(result)
// }

// func TestSomething(t *testing.T) {
// 	client, err := NewCherryPyClient("https://192.168.50.10:8000", "test", "test", "pam")
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	cmd := CherryPyCommand{}
// 	cmd.Client = "local"
// 	cmd.Target = "minion1"
// 	cmd.Function = "state.highstate"

// 	cmds := make([]CherryPyCommand, 1)
// 	cmds[0] = cmd

// 	result, err := client.Run(cmds)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	print(result)
// }
