package activities

import "testing"

func TestValidateDittoTopic_Positive(t *testing.T) {
	validTopics := []string{
		"org.example/test-thing/things/twin/commands/create",
		"org.example/test-thing/things/twin/commands/modify",
		"org.example/test-thing/things/twin/commands/delete",
		"org.example/test-thing/things/live/commands/create",
		"org.example/test-thing/things/live/commands/modify",
		"org.example/test-thing/things/live/commands/delete",
		"org.example/test-thing/things/twin/events/update",
		"org.example/test-thing/things/live/messages/modify",
		// Positive policy topics
		"org.example/test-policy/policies/commands/create",
		"org.example/test-policy/policies/commands/modify",
		"org.example/test-policy/policies/commands/delete",
		"org.example/test-policy/policies/commands/update",
	}
	for _, topic := range validTopics {
		if err := validateDittoTopicRegex(topic); err != nil {
			t.Errorf("expected topic '%s' to be valid, got error: %v", topic, err)
		}
	}
}

func TestValidateDittoTopic_Negative(t *testing.T) {
	invalidTopics := []string{
		"org.example/test-thing/thing/twin/commands/create",      // typo: thing
		"org.example/test-thing/things/other/commands/create",    // invalid 4th part
		"org.example/test-thing/things/twin/command/create",      // typo: command
		"org.example/test-thing/things/twin/commands/unknown",    // invalid 6th part
		"org.example/test-thing/things/twin/commands",            // too short
		"org.example/test-thing/things/twin/events",              // too short
		"org.example/test-thing/things/live/messages",            // too short
		"org.example/test-thing/things/live/commands/invalidcmd", // invalid command
		// Negative policy topics
		"org.example/test-policy/policy/commands/create",        // typo: policy
		"org.example/test-policy/policies/command/create",       // typo: command
		"org.example/test-policy/policies/commands/unknown",     // invalid command
		"org.example/test-policy/policies/commands",             // too short
		"org.example/test-policy/policies/events/create",        // invalid 5th part for policies
		"org.example/test-policy/policies/live/commands/create", // use live for policies
		"org.example/test-policy/policies/twin/commands/modify", // use twin for policies
		"org.example/test-policy/policies/live/commands/modify",
	}
	for _, topic := range invalidTopics {
		if err := validateDittoTopicRegex(topic); err == nil {
			t.Errorf("expected topic '%s' to be invalid, but got no error", topic)
		}
	}
}
