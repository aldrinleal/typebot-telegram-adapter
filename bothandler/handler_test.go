package bothandler

import (
	"encoding/json"
	"testing"
)

func TestSomething(t *testing.T) {
	rawMessage := `{
	         "id": "e7ckog4qfrsuw0pbh62twb8e",
	         "type": "text",
	         "content": {
	            "type": "richText",
	            "richText": [
	               {
	                  "type": "p",
	                  "children": [
	                     {
	                        "text": "Faaala, "
	                     },
	                     {
	                        "type": "inline-variable",
	                        "children": [
	                           {
	                              "type": "p",
	                              "children": [
	                                 {
	                                    "text": "aldrin"
	                                 }
	                              ]
	                           }
	                        ]
	                     },
	                     {
	                        "text": "!"
	                     }
	                  ]
	               }
	            ]
	         }
	      }`

	messages := &TypeBotMessage{}

	err := json.Unmarshal([]byte(rawMessage), messages)

	// Assert err is nil
	if err != nil {
		t.Errorf("Error: %s", err)
	}

	replyText := getRichTextFor(*messages)

	// Assert replyText is "Faaala, aldrin!"
	if replyText != "Faaala, aldrin!" {
		t.Fail()
	}
}
