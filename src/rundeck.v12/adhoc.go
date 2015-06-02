package rundeck

type ExecutionId struct {
        ID              string            `xml:"id,attr"`
}

func (c *RundeckClient) RunAdhoc(projectId string, exec string) (ExecutionId, error) {
	options := make(map[string]string)
	options["project"] = projectId
	options["exec"] = exec
	var data ExecutionId
	err := c.Get(&data, "run/command", options)
	return data, err
}
